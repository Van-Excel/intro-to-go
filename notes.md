RUNES, CODEPOINTS AND ENCODINGS
-----------------------------------
what is a unicode code point in go?
is it the integer representing a string in a character set?

6:25 PM
Yes exactly. A unicode code point is the integer assigned to a character in the Unicode character set.

In Go a code point is called a rune which is just an alias for int32:

go
type rune = int32
Examples:

go
'A' = 65
'a' = 97
'0' = 48
'😀' = 128512
'€' = 8364
In Go:

go
r := 'A'          // rune, value is 65
fmt.Println(r)    // 65  ← the integer
fmt.Printf("%c", r) // A  ← the character

s := "van"
for _, r := range s {
    fmt.Println(r)  // 118, 97, 110  ← code points
}
```



---

**The full pipeline you built earlier:**
```
character → code point → bytes
  'A'     →    65      →  [65]         (UTF-8, 1 byte)
  '😀'    →  128512    →  [240][159][152][128]  (UTF-8, 4 bytes)
code point = the integer in the character set (Unicode)
encoding = the rule for turning that integer into bytes (UTF-8)
So rune in Go = code point = the integer representing a character. UTF-8 then decides how many bytes to use to store that integer.


Use %b for binary or %x for hex:
gor := 'A'
fmt.Printf("%d\n", r)   // 65          ← decimal
fmt.Printf("%b\n", r)   // 1000001     ← binary (bits)
fmt.Printf("%x\n", r)   // 41          ← hex
fmt.Printf("%c\n", r)   // A           ← character
For the actual UTF-8 bytes of a string:
gos := "A😀"
fmt.Printf("%x\n", []byte(s))   // 41 f09f9880  ← hex bytes
fmt.Printf("%b\n", []byte(s))   // binary bytes

All format specifiers for numbers:
go%d  → decimal    (65)
%b  → binary     (1000001)
%x  → hex        (41)
%o  → octal      (101)
%c  → character  (A)
%T  → type       (int32)
Most common for seeing raw bytes is %x — easier to read than binary, and it's what you'll see in network debugging and memory dumps.    s := "van"
    for _, r := range s {
        fmt.Println(r) // 118, 97, 110  ← code points
        fmt.Printf("unicode code point:%d, character:%c, utf-8: %b, type: %T \n", r, r, []byte(r), r)
getting an error that []byte() cant convert r since its a rune 
ahy cant it convert an int to bytes6:38 PMBecause []byte() expects a string, not a rune/int32.
The confusion is that []byte(s) converts a string to its UTF-8 bytes — it's not a general int-to-bytes converter.
To get the UTF-8 bytes of a rune, convert it to a string first:
gofor _, r := range s {
    fmt.Printf("unicode code point:%d, character:%c, utf-8: %b, type: %T \n", r, r, []byte(string(r)), r)
}
string(r) converts the rune back to a string, then []byte() gets the UTF-8 bytes.

Why []byte(r) fails:
go[]byte(r)        // error - r is int32, []byte expects string
[]byte(string(r)) // works - convert rune → string → bytes
[]byte is specifically designed to get the UTF-8 byte representation of a string — not to convert arbitrary integers to bytes. For that you'd use encoding/binary or struct.Pack like you saw earlier.[]byte() is only for strings? how do i convert an int with binary?6:39 PMYes []byte() is only for strings.
For integers you use encoding/binary:
goimport "encoding/binary"

n := uint32(65)

buf := make([]byte, 4)  // 4 bytes for uint32
binary.LittleEndian.PutUint32(buf, n)
fmt.Printf("%b\n", buf)  // [1000001 0 0 0]

// or big endian
binary.BigEndian.PutUint32(buf, n)
fmt.Printf("%b\n", buf)  // [0 0 0 1000001]

Or using encoding/binary with a buffer:
goimport (
    "bytes"
    "encoding/binary"
)

var buf bytes.Buffer
binary.Write(&buf, binary.LittleEndian, uint32(65))
fmt.Printf("%b\n", buf.Bytes())  // [1000001 0 0 0]

This connects back to your socket/protocol work:
go// remember this from earlier?
struct.pack('>I', 65)   // Python - pack int to bytes big endian

// Go equivalent
binary.BigEndian.PutUint32(buf, 65)  // same thing in Go
Same concept — you're converting an integer to its raw byte representation with a specified endianness. []byte() can't do this becaus


Because binary.PutUint32 writes bytes into an existing memory location — it doesn't allocate memory itself:
gobinary.BigEndian.PutUint32(buf, 65)
//                          ↑
//                 needs somewhere to write the bytes
It's the same pattern you saw in C:
cread(fd, &buf, 4)  // read needs a buffer to write into
binary is just saying "give me a place to put these bytes" — it won't create that place for you.

Why not just return a []byte directly?
Some functions do:
go// this allocates for you
buf := make([]byte, 4)
binary.BigEndian.PutUint32(buf, 65)  // writes into buf

// alternatively append version allocates internally
buf = binary.BigEndian.AppendUint32(buf, 65)  // returns new slice
```

---

**The design philosophy:**

Go's standard library prefers **no hidden allocations** — same reason `copy` doesn't resize, same reason `read()` in C needs a buffer. You control memory, the function just does the work.

So:
```
you create the buffer  → you control memory allocation
binary writes into it  → it just does the byte conversionwhy dont we do that  for string(), []byte() and other type casting functions6:44 PMBecause those are built-in type conversions handled by the compiler, not functions that write to memory.
go[]byte(s)    // compiler handles this, allocates internally
string(r)    // compiler handles this, allocates internally
int(x)       // compiler handles this, no allocation needed
The compiler knows exactly how much memory is needed upfront:
gos := "van"
[]byte(s)    // compiler knows → 3 bytes, allocates 3, copies, done
```

---

**Why `binary` is different:**

`binary.PutUint32` is a **regular function** in a package, not a compiler built-in. Regular functions in Go follow the no-hidden-allocation philosophy.
```
built-in conversions → compiler magic, allocates for you
package functions    → no hidden allocations, you provide buffer

The pattern generalizes:
go// compiler built-ins → just convert, no buffer needed
[]byte(s)
string(b)
int(f)

// package functions → need buffer
binary.PutUint32(buf, n)   // encoding/binary
json.Unmarshal(buf, &v)    // encoding/json
copy(dst, src)             // builtin but needs dst
encoding/json is the same — you provide the buffer, it writes into it. The rule is: if it's a package function, expect to manage your own memory.

explain how scanf works in go
does it read from a terminal(file) and then write to some location for us
does it have to write to a variable or it can write to a file

4:32 PM
Good question. Let me break down exactly what Scanf does.

What Scanf actually is
fmt.Scanf reads formatted input from standard input (stdin) and writes the parsed values into variables you provide via pointers.

go
var name string
var age int

fmt.Scanf("%s %d", &name, &age)
```

Two distinct things are happening:
```
stdin (keyboard/terminal)  →  Scanf reads raw bytes
                           →  parses according to format string
                           →  writes parsed values into &name and &age
Your question — does it have to write to a variable?
It writes to any memory address you give it via a pointer. A variable is the most common case but the mechanism is just pointer writes.

go
// writing to a variable — most common
var amount float64
fmt.Scanf("%f", &amount)

// writing to a struct field — same mechanism
type Investor struct {
    Name   string
    Amount float64
}
inv := Investor{}
fmt.Scanf("%s %f", &inv.Name, &inv.Amount)

// writing to a slice element
amounts := make([]float64, 5)
fmt.Scanf("%f", &amounts[0])
All of these are just pointers. Scanf does not know or care that it is a variable, struct field, or slice element — it just writes to the address you give it.

Can it write to a file?
Not directly — Scanf is hardcoded to read from stdin. But the fmt package has Fscanf which accepts any io.Reader as the source:

go
// Scanf — always reads from stdin
fmt.Scanf("%s", &name)

// Fscanf — reads from any io.Reader
file, _ := os.Open("investors.txt")
fmt.Fscanf(file, "%s %f", &name, &amount)

// Sscanf — reads from a string
data := "John 5000.00"
fmt.Sscanf(data, "%s %f", &name, &amount)
```

The write target is always a pointer you provide — what changes between `Scanf`, `Fscanf`, and `Sscanf` is only the **source** being read from.

---

## The full picture
```
SOURCE                    FUNCTION       TARGET
─────────────────────────────────────────────────
stdin (terminal)      →   Scanf      →   &variable
any io.Reader (file)  →   Fscanf     →   &variable
string                →   Sscanf     →   &variable
The write target is always a pointer to memory — variable, struct field, or slice element. The source is what varies. Scanf itself cannot write to a file — that is what os.File and io.Writer are for, which is a completely separate concern.

Yes, exactly. %s in Scanf reads until it hits whitespace — space, tab, or newline.
What happened
input:  "vab exe\n"
%s      reads  "vab"  → stops at the space
" exe"  is still sitting in stdin buffer, unread

The stdin buffer is the key mental model
Scanf does not consume all of stdin. It reads only what the format string asks for and leaves the rest in the buffer.
before Scanf:   [ v | a | b |   | e | x | e | \n ]
                                ↑
after %s:       [   | e | x | e | \n ]  ← these are still in buffer
                ↑
                next Scanf call will start reading from here
This causes subtle bugs — a second Scanf call would read exe from the leftover buffer without waiting for new input.

How to read a full line including spaces
go// option 1 — bufio.Scanner, most common in Go
scanner := bufio.NewScanner(os.Stdin)
fmt.Print("enter your name: ")
scanner.Scan()
name := scanner.Text()  // reads entire line, strips newline
fmt.Println("user entered:", name)

// option 2 — bufio.Reader
reader := bufio.NewReader(os.Stdin)
fmt.Print("enter your name: ")
name, _ := reader.ReadString('\n')  // reads until newline
name = strings.TrimSpace(name)      // strip the trailing \n
```

---

## The format verbs and what they stop at
```
%s      →  stops at any whitespace (space, tab, newline)
%d      →  stops at any non-digit character
%f      →  stops at any non-numeric character
%v      →  same as %s for strings — stops at whitespace
So Scanf is really designed for whitespace delimited tokens, not full lines. For anything involving spaces — names, addresses, descriptions — bufio.Scanner is the right tool.

Using Scanf in a loop
govar fullName string
var part string

fmt.Print("enter your name: ")
for {
    n, err := fmt.Scanf("%s", &part)
    if n == 0 || err != nil {
        break
    }
    // check if we should stop — but how do we detect newline?
    fullName += part + " "
}
fmt.Println("user entered:", strings.TrimSpace(fullName))
The problem here is fundamental — %s already consumed the whitespace delimiter so you never actually see it. By the time Scanf returns, the space is gone. You cannot check for the newline because Scanf skipped past it to find the next token.
You would have to use %c to read character by character:
govar fullName strings.Builder
var ch rune

fmt.Print("enter your name: ")
for {
    fmt.Scanf("%c", &ch)
    if ch == '\n' {
        break
    }
    fullName.WriteRune(ch)
}
fmt.Println("user entered:", fullName.String())
This works but it is fighting against what Scanf is designed for. Reading character by character with a format verb is the wrong layer of abstraction.

Why Go has better tools for this
Scanf comes from C where it made sense because C had limited standard library options. Go was designed with better primitives from the start.
go// bufio.Scanner — cleanest, most idiomatic
scanner := bufio.NewScanner(os.Stdin)

fmt.Print("enter your name: ")
scanner.Scan()
fullName := scanner.Text()
fmt.Println("user entered:", fullName)
```

Internally `scanner.Scan()` is doing the loop and newline detection for you — it just hides the complexity behind a clean API.

---

## The mental model for when to use each
```
Scanf           →  reading structured whitespace delimited tokens
                   "enter age and amount: " → Scanf("%d %f", &age, &amount)

bufio.Scanner   →  reading full lines of text
                   "enter your name: "      → scanner.Scan()

bufio.Reader    →  reading lines with more control over delimiter
                   reader.ReadString('\n')
The rule of thumb is — if your input has spaces that are meaningful content, do not use Scanf. It was designed for tokens separated by whitespace, not for capturing whitespace as part of the value.

println()- adds spaces between args and appends \n at end of args moving cursor to the next line
print()- does not add space between args or append it with a newline(\n)

https://zetcode.com/golang/inttostring/

difference between Itoa and FormatInt
What Itoa actually is
strconv.Itoa is literally just a wrapper around FormatInt:
go// from Go standard library source
func Itoa(i int) string {
    return FormatInt(int64(i), 10)  // always base 10
}
That is the entire implementation. Itoa is a convenience function for the most common case — converting an integer to a base 10 string.

What FormatInt gives you that Itoa cannot
gonumber := int64(255)

// Itoa — base 10 only
strconv.Itoa(255)                        // "255"

// FormatInt — any base
strconv.FormatInt(int64(number), 2)      // "11111111"  binary
strconv.FormatInt(int64(number), 8)      // "377"       octal
strconv.FormatInt(int64(number), 10)     // "255"       decimal
strconv.FormatInt(int64(number), 16)     // "ff"        hexadecimal
strconv.FormatInt(int64(number), 36)     // "73"        base 36
Any base from 2 to 36 is valid.

The second distinction — type precision
govar bigNumber int64 = 9223372036854775807  // max int64

strconv.Itoa(int(bigNumber))              // unsafe — int is 32 or 64 bit depending on platform
strconv.FormatInt(bigNumber, 10)          // always safe — explicitly int64
```

`Itoa` takes a plain `int` which is platform dependent — 32 bit on 32 bit systems, 64 bit on 64 bit systems. `FormatInt` always takes `int64` so it is explicit and portable.

---

## When to use each
```
Itoa            →  simple int to decimal string, most common case
FormatInt       →  need a specific base (binary, hex, octal)
                   working explicitly with int64
                   building things like colour codes, bitmasks, permissions
A concrete G4M example — loan permission flags stored as bitmasks:
gopermissions := int64(0b10110101)  // binary flags for loan permissions

// store human readable binary representation for debugging
strconv.FormatInt(permissions, 2)   // "10110101"
strconv.FormatInt(permissions, 16)  // "b5"  compact hex for storage
You could not do this with Itoa at all. So Itoa is the shortcut for the 90% case, FormatInt is the full tool when you need control over the base or explicit int64 handling.

number := 255  // stored in memory as 11111111
```
```
FormatInt(255, 2)   →  "11111111"   string of 8 chars, each char is '1' or '0'
FormatInt(255, 10)  →  "255"        string of 3 chars, '2','5','5'
FormatInt(255, 16)  →  "ff"         string of 2 chars, 'f','f'
```

Same number in memory. Three completely different string representations. None of them are the actual bits — they are all human readable text expressing the same value in different bases.

---

## The cleaner way to see it
```
255 in memory is always:   00000000 11111111  (as 16 bit example)

base 2  string  →  "11111111"   how many 1s and 0s do I need to express this value
base 10 string  →  "255"        how many 0-9 digits do I need to express this value
base 16 string  →  "ff"         how many 0-f digits do I need to express this value


so as a string say "255" is each character 1 byte making it 3 bytes instead of the 16bits representing it in binary?

5:38 PM
Yes exactly. That is the real cost of stringification.

The memory comparison
int 255 in memory       →  1 byte (fits in uint8) or 2/4/8 bytes depending on type
string "255" in memory  →  3 bytes, one per character ('2', '5', '5')
Each character in a Go string is 1 byte (UTF-8). So:

"255"       →  3 bytes
"11111111"  →  8 bytes
"ff"        →  2 bytes
All represent the same value 255. The binary string representation is actually the most expensive — 8 bytes to represent a number that fits in 1 byte.

This is why you never store numbers as strings
1 million records, storing an int32 field:
    as int32    →  4 bytes × 1,000,000  =  4MB
    as string   →  10 bytes × 1,000,000 =  10MB  ("1000000" is 7 chars + string overhead)
Databases, binary protocols, and file formats all store numbers in their binary form for exactly this reason.

When stringification makes sense despite the cost
displaying to a user        →  humans read "255" not raw bytes
logging                     →  log files are text
JSON APIs                   →  JSON is text based
CSV files                   →  text format
You convert to string at the boundary between your system and humans or text based systems. Internally you keep numbers as numbers.

The general principle
computation and storage     →  keep as binary int/float
display, logging, transfer  →  convert to string at the boundary
Itoa and FormatInt are boundary tools — you reach for them at the edges of your system, not in the middle of it.




FILES, STREAMS, READALL()
---------------------------------
difference between a stream and file on disk

A file on disk is static — all the data is already there, you know its size upfront, and you can jump to any position (seek forward, backward, random access).

A stream is dynamic — data arrives sequentially over time, you often don't know the total size upfront, and you can only read forward. You get the data as it arrives, chunk by chunk.

So for an HTTP response body:

The server is sending data over the network piece by piece
You don't know how much is coming until it's done
You can't jump to the middle of it
That's why io.ReadAll keeps reading until EOF — it's waiting for the stream to finish delivering all its data before handing it all back to you as a complete []byte.

if we dont know size, how do we know eof? from content length header? and for sockets or just tcp connections when do we know we have reached eof

9:48 AM
Great question. There are a few ways EOF is determined:

For HTTP:

Content-Length header — if the server sends this, the client knows exactly how many bytes to expect and stops after reading that many
Chunked transfer encoding — if no Content-Length, the server sends data in chunks, each prefixed with its size, and a final chunk of size 0 signals EOF
Connection close — the server just closes the connection when done, and that closure signals EOF
For raw TCP:

EOF is signaled when the other side calls close() on the socket, which sends a FIN packet
When you receive the FIN, your Read call returns 0 bytes and an io.EOF error
There's no concept of "message boundaries" in raw TCP — it's just a continuous stream of bytes, so the application layer has to define its own protocol for knowing when a message is complete (like HTTP does with Content-Length or chunked encoding)
So at the TCP level it's always the FIN packet, but higher level protocols like HTTP build their own signaling on top of that.

when do we use chunked transfer encoding

9:49 AM
Chunked transfer encoding is used when the server doesn't know the total size of the response before it starts sending — i.e. the content is being generated dynamically.

Common cases:

Dynamically generated pages — server is building the HTML on the fly from database queries, it doesn't know the final size until it's done
Streaming responses — like server-sent events, video streaming, or LLM responses (like this one) where data is produced and sent progressively
Large file compression on the fly — server is compressing while sending, so it doesn't know the compressed size upfront
Database query results — server starts streaming rows before the full query finishes
Basically the rule is:

Know the size upfront → use Content-Length
Don't know the size upfront → use chunked transfer encoding
A good real world example you'd know — when Claude responds to you, the response is streamed back chunk by chunk as the model generates tokens. The server has no idea how long the response will be when it starts sending, so chunked encoding is used.

so basically readall is a loop for read until it returns 0 or EOF

9:51 AM
Yes exactly. Under the hood io.ReadAll is essentially:

go
for {
    n, err := reader.Read(buf)
    // append buf[:n] to result
    if err == io.EOF {
        break
    }
}
Keep reading chunks into a buffer, accumulating them, until Read signals EOF — then return everything accumulated as a single []byte.




ROUTER
router or ServerMux is a router
we have a server that orchestrates the network layer, reads data and hands it to a handler or 
some function which manipulates the data received from the socket with the logic it contains
if we are building an http server which is essentially a tcp server that speaks http we need
a way to determine which handler to call and pass data to depending on the path in the request line
we could use conditionals like if r.url.Path == '/' call homeHandler() but we can also create 
a component or object to handle that responsibility for cleaner architecture
multiplexer is seen at the tcp layer when matching requests to processes on a machine and at the cpu
level when processor wants to write or read data
- need to revise tcp server, parsing of request, serialization and confirm whether handler
writes response directly to connection or gives it to server to send it
type ServeMux struct{}
func (s *ServeMux) HandleFunc(){}
func (s *ServeMux) Handle(pattern string handler Handler){}
type Handler interface {
	ServeHTTP(ResponseWriter, *Request)
}
func (mux *ServeMux) handler(r *Request) Handler {
mux.mu.RLock()
defer mux.mu.RUnlock()
// Host-specific pattern takes precedence over generic ones
h := mux.match(r.Host + r.URL.Path)
if h == nil {
h = mux.match(r.URL.Path)
}
if h == nil {
h = NotFoundHandler()
}
return h
}

muxServer implements ServeHTTP(w responsewriter r *Request) which also makes it a handler now

type ServeMux struct {
    mu sync.RWMutex // because of concurrency, we have to use a mutex here
    m map[string]muxEntry // router rules, every string mapping to a handler
    }

type muxEntry struct {
    explicit bool // exact match or not
    h Handler
    }
- when the server calls mux.ServeHTTP() what happens?
 - the request that has been parsed is passed to mux.ServeHTTP 
 - it retrieves the url using r.URL.path or the [method] [path] [] and lookup against its 
   internal map (map[string]muxEntry // router rules, every string mapping to a handler)
 - this will retrieve the muxEntry struct which contains the necessary handler if the pattern exists
 - we can call the ServeHttp() on that handler so we run its logic against the request data and
  return a response
-the muxServer is an object that routes requests to specific handlers so it needs some metadata and
 data structures so it can do its job. We use a map or hashtable to store our mappings of 
 routes to handlers and perform lookups when a request arrives
 HandleFunc and Handle methods allow this object to receive information from user and store it in
 its database (map or key value store) for future retrieval when a request arrives
 -There is also a Handler function which can return the handler for a particular request

Sacrifice 4 months
 Data representation, CPU, cache, RAM, DISK,Networks (TCP, UDP, IP, DNS, TLS,security, HTTP, authentication-oauth2, jwt, sessions), Data structures, OS(Linux, 3 easy pieces), OOP(design patterns),DB, System Design (redis, kafka, postgres, patterns, distributed systems, reliability patterns, access patterns and storage engines)
 Accounting and Finance- CFA, FMVA, Financial Analyst
 Trading
 Business

 -deep understanding of CPU internals and Go and C
 - deep understanding of core data structures
 - improved math skills
 - Ability to understand, build and write about redis and database projects



TCP
----
application sends a message
tcp breaks it down into segments or bytes along with header. These segments are numbered and sent
This allow reliable and in order delivery (acks and retransmission on dropped or missing bytes)
Tcp or os layer arranges bytes (lets say http message) and stores in buffer of socket
we can now parse this into a struct and use the data

How does this apply to custom protocols and serialisation and deserialisation in general?

SERVICES
Services are generally long lived processes that have a listening socket. This socket which represents
a network connection listens for client requests to connect to the sockets of clients on other hosts across a network. It also has a read and write buffer or a send and receive buffer. These are like files. The read buffer is for receiving messages from a connected client similar to the concept of a message inbox. the send buffer is for writing messages for the OS to send to the client's socket using the network stack.
Services run on hosts or end systems which operate multiple processes or services so services are mapped to ports to help identify the service or process a message is for (remember the 4 tuples that describe a network connection)
service basically call socket(), bind(), listen(), accept(), parse read data and pass it to handler
service -> handler -> database
-Think about how certain requests may consist of multiples services or processes across multiple hoss and networks communicating together and how it affects latency, nature of service workload, etc
-Think about how important TCP connections and lifecycles, multiprocessing, threads, blocking and non blocking calls, Multiplexed IO are to services


TASKS
- implement api
- practise db connection and execution
-implement router
- learn http parsing
- revise servers (non blocking IO, multiplexing, event loops) and data serialisation( custom protocols, redis protocol)
-implement OOP and services
-resume build your own redis (learn kafka clone fork)
- resume build your own db (learn badgerdb fork)

-networks, servers,loops, connections and sockets
- data serialization(protocols, character sets and encodings), 
- reading data, io, and parsers (strings, numbers, objects)
- http and how to parse a http request and create an http response
- application level protocols: redis, http, custom
- handlers and threads
- routing
- security (TLS)


CONCURRENCY
-----------
Process -> Threads -> user level threads
We start with processes but they are expensive and can be time consuming so we use threads
Threads are lighter or less expensive than processes but they introduce race conditions
We use synchronization methods to tackle race conditions
These sometimes lead to deadlocks and spinlocks?
We also have threading models: 1:N and M:N




4 months: make good judgments