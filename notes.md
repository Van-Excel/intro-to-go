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