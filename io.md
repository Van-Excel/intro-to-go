
analogy
you are essentially asking is there data in the kernel buffer?

sync blocking
you ask and if there is none you wait

sync non blocking- returns immediately so you keep polling which burns cycles.Why? because when the thread is contexted switched it doesn't go into the wait queue or sleep but into the ready queue where it is continuously scheduled till the syscall(read()) returns a read value.
you can instead use epoll() which will make thread sleep but wake it up when there is data in kernel buffer
So here you still call read() again to trigger the copy from kernel buffer to user space buffer

async non blocking- no polling you are notified after kernel copies from NIC into kernel buffer
you submit your request. kernel will handle copy from NIC -> kernel buffer -> userspace buffer. it will only
notify after the data is in your buffer so you can use it immediately
async blocking - rarely used


Great question. These are often confused but they are distinct concepts.

Two separate axes:
                    blocking        non-blocking
                 ┌─────────────┬─────────────────┐
synchronous      │   normal    │  spin/poll loop  │
                 │   read()    │  O_NONBLOCK      │
                 ├─────────────┼─────────────────┤
asynchronous     │   rarely    │   aio_read()     │
                 │   used      │   io_uring       │
                 └─────────────┴─────────────────┘

Synchronous — you drive the operation:
You call read(), you wait for or check the result, you handle it. You are in control of when and how you ask.
c// synchronous blocking
n = read(fd, buf, size); // you wait here until done

// synchronous non-blocking
n = read(fd, buf, size); // you check immediately
if (n == -1 && errno == EAGAIN) {
    // not ready, you decide what to do next
}
In both cases you initiated the call and you check the result. Synchronous just means the operation and your code move in step.

Asynchronous — OS drives the operation:
You say "start this read, tell me when done" and move on. You don't check — the OS notifies you.
c// asynchronous non-blocking (aio)
struct aiocb cb;
cb.aio_fildes = fd;
cb.aio_buf = buf;
cb.aio_nbytes = size;

aio_read(&cb); // returns immediately, OS does read in background

// you do other work here...
// OS signals you when read is complete
You never polled or checked. OS came back to you.

The key difference:
who checks for completionwhen does caller returnsync blockingyou, by waitingafter operation donesync non-blockingyou, by pollingimmediately, you retryasyncOS notifies youimmediately, forget it

Concrete analogy:
You're hungry and waiting for food.

sync blocking → you stand at the counter staring until food is ready
sync non-blocking → you keep walking back to counter asking "is it ready yet?"
async → you give your number, go sit down, they call you when ready


So epoll — which category is it?
Interesting — epoll is actually synchronous non-blocking:
cepoll_wait(); // you call this, you block until fd is ready
read(fd, buf, size); // then you do the read yourself
You are still driving everything. You just sleep efficiently instead of spinning. The OS isn't doing the read for you — it's just telling you when to do it yourself.
True async on Linux is io_uring — you submit read requests to a queue, OS completes them and posts results back. You never call read() yourself.

This is also what JavaScript async/await actually is:
javascriptconst data = await read(fd); // looks sync, is it?
It's asynchronous non-blocking under the hood — libuv submits the I/O, event loop does other work, callback fires when done. await just makes it look synchronous in your code.but why does sync non blocking return and ask for data immediately? i thought reads always take time or that isnt the case?7:53 AMGood observation. You're right that reads take time — but where that time is spent is the key.

The read has two phases:
client sends data
    │
    └── data arrives in kernel buffer (socket receive buffer)
        │                                    ▲
        this happens independently           │
        of your application                  │
                                        time spent
                                        here, not in
                                        your read() call
        │
        your read() copies from kernel buffer → your buf
        │
        this part is fast, almost instant

The "time" in a read is waiting for data to arrive in kernel buffer:
blocking read()
    │
    ├── data not in kernel buffer yet
    │       thread sleeps here waiting
    │       ...
    │       data arrives in kernel buffer
    │       thread wakes up
    │
    └── copies kernel buffer → your buf (fast)
        returns n bytes
non-blocking read()
    │
    ├── data not in kernel buffer yet?
    │       return -1 EAGAIN immediately
    │       (no sleeping, no waiting)
    │
    └── data already in kernel buffer?
        copies kernel buffer → your buf (fast)
        returns n bytes immediately

So non-blocking read() is basically asking:
"is there data sitting in the kernel buffer RIGHT NOW?
 yes → give it to me
 no  → don't wait, just tell me nothing is there"

This is why epoll is powerful:
Instead of you asking "is data there yet?" in a loop:
c// you asking repeatedly - wasteful
while (1) {
    n = read(fd, buf, size); // is it there yet?
    if (n == -1 && errno == EAGAIN) continue; // no, ask again
}
You let the OS tell you exactly when data lands in the kernel buffer:
cepoll_wait(); // sleep until OS says "data is in buffer"
read(fd, buf, size); // now guaranteed to be fast, data is there
```

---

**The mental model:**
```
network
    │
    │  data travelling...
    ▼
kernel receive buffer  ← data lands here first
    │
    │  read() copies from here
    ▼
your buf
The network journey (milliseconds) happens in kernel space independently. Your read() only deals with the last step — copying from kernel buffer to your buffer — which is essentially instant.so blocking is basically asking is there data in the kernel buffer and if its not you just go and wait in a queue till its ready (wait in a queue)
for non blocking you keep on asking intermittently till its ready (polling and sometimes you wait in a queue when your time slice on processor is up)7:59 AMExactly, perfect summary. Let me just sharpen it slightly.

Blocking:
read() called
    │
    is data in kernel buffer?
    │
    NO → thread moves from RUNNING to WAITING queue
         off the CPU completely
         not consuming any cycles
         │
         data arrives in kernel buffer
         │
         OS moves thread from WAITING → READY queue
         scheduler picks it up
         thread resumes, copies data

Non-blocking:
read() called
    │
    is data in kernel buffer?
    │
    NO → returns EAGAIN immediately
         thread stays RUNNING
         your code decides what to do
         │
         loop back, call read() again
         │
         ⏰ time slice expires → OS preempts
         thread goes to READY queue (not WAITING)
         gets CPU again soon
         loop continues asking...
         │
         eventually data arrives
         read() returns data

The critical difference in which queue:
no data scenarioCPU consumptionblockingWAITING queue, off CPUzero, sleepingnon-blockingREADY queue, keeps getting CPU turnsburns cycles polling

Blocking thread is truly asleep — not consuming anything. Non-blocking thread is restless — keeps getting scheduled, keeps checking, wastes cycles.
Which is exactly why raw non-blocking polling alone is wasteful, and epoll exists to give you the best of both — non-blocking reads but sleeping efficiently like blocking when there's nothing to do.thread resumes, copies data
where does thread copy data? i thought kernel does all copying8:01 AMGood catch. Let me be precise.
The kernel does the copying — but on behalf of your thread.

What actually happens:
data arrives from network
    │
    └── kernel copies: NIC buffer → socket receive buffer
        (this happens via DMA, independent of your thread)
        │
        your thread calls read()
        │
        CPU switches to kernel mode (syscall)
        │
        kernel copies: socket receive buffer → your buf
        (kernel code running, but on your thread's behalf)
        │
        CPU switches back to user mode
        │
        your thread continues with data in buf

There are actually two copies:
NIC → kernel socket buffer → your userspace buf
      ───────────────────    ──────────────────
      copy 1                 copy 2
      kernel does this       kernel does this too
      independently          but triggered by your read()

So when I said "thread copies data" I was imprecise:
More accurately:

your thread triggers the copy by calling read()
kernel executes the copy on your thread's behalf
your thread is in kernel mode during this, not user mode


This is actually a known performance problem:
Two copies is expensive for high performance systems. This is why zero-copy techniques exist:
normal read()
NIC → kernel buffer → your buf        (2 copies)

sendfile() zero-copy
NIC → kernel buffer → NIC/socket      (kernel to kernel, 1 copy)
never touches userspace at all

io_uring with zero-copy
NIC → your buf directly               (1 copy, no kernel buffer)
This matters a lot for things like video streaming servers sending gigabytes of data — eliminating that extra copy saves significant CPU and memory bandwidth.got it now
but still some inconsistencies with async model
what does it mean by you dont do reads yourself but you are just notified. how is it different from sync non blocking aside the polling8:03 AMGood question. Let me make the distinction concrete.

Sync non-blocking (epoll model):
You are still the one doing the read. epoll just tells you when to do it.
cepoll_wait(); // sleep until data is in kernel buffer
// OS wakes you: "fd 42 is ready"

// YOU still have to do the read yourself
n = read(fd, buf, size); // you trigger the kernel copy
                         // kernel buffer → your buf
```

The sequence:
```
data arrives → kernel buffer
                    │
               epoll wakes you
                    │
               you call read()
                    │
               kernel copies: kernel buffer → your buf
                    │
               you handle data
You are still involved in triggering the copy.

Async (io_uring / aio):
You tell the OS "read this fd into this buf when data arrives". You walk away. OS does everything including the copy.
c// you submit a read request
struct io_uring_sqe *sqe = io_uring_get_sqe(&ring);
io_uring_prep_read(sqe, fd, buf, size, 0);
io_uring_submit(&ring); // submitted, forget it

// do completely other work...
// OS is watching fd AND will do the copy itself

// later, check completion queue
struct io_uring_cqe *cqe;
io_uring_wait_cqe(&ring, &cqe); // data is already in buf
int n = cqe->res; // just the result, copy already done
```

The sequence:
```
data arrives → kernel buffer
                    │
               OS does copy itself: kernel buffer → your buf
               (no involvement from you)
                    │
               OS posts completion to queue
                    │
               you check queue, data already in buf
```

---

**The fundamental difference:**

| | epoll (sync nb) | io_uring (async) |
|---|---|---|
| who watches fd | OS | OS |
| who triggers copy | **you** (via read()) | **OS** (automatically) |
| when you're involved | before copy | after copy |
| what you get notified | "fd is ready, go read" | "data is already in your buf" |

---

**Analogy:**

- **epoll** → restaurant buzzer goes off, you walk to counter and pick up your food yourself
- **io_uring** → waiter brings food directly to your table, you just eat

In both cases you were notified. The difference is whether the food is already at your table or you still have to go get it.

---

**Why does this matter?**

With epoll at very high loads you still have overhead:
```
wake up → call read() → syscall → kernel mode → copy → user mode → handle
```

That syscall + mode switch happens for every single read. At 1 million requests/second that adds up.

With io_uring:
```
check completion queue → handle
No syscall per read. OS batched everything. This is why io_uring is significantly faster for high throughput servers — fewer syscalls, fewer mode switches, OS does copies in background.





--------------------------------------------------------------------------------------------------------------------
FORK(), PROCESSES, THREADS

why are the parent process and chil dprocess running
after fork() a chil dprocess is created by os
does it schedule and run immediately?

12:09 AM
Good question. The short answer is: it depends on the scheduler — there's no guarantee.

What actually happens:

When fork() is called, the OS creates the child process and puts it in the ready queue. From that point, the scheduler decides who runs next — parent or child. It could be either.

fork() called
     │
     ▼
OS creates child ──→ both go into ready queue
                          │
                    scheduler picks one
                    (could be parent or child)
Why you see both running:

After fork(), you have two independent processes. The OS will eventually schedule and run both of them — not necessarily at the same time, but both will get CPU time. That's why you see output from both.

The order is non-deterministic:

Run your program a few times — sometimes parent prints first, sometimes child. You can't rely on order unless you explicitly synchronize using wait().

c
// in parent block
wait(NULL); // parent waits for child to finish before continuing
Key mental model:

fork() doesn't just "create" a child and move on. After the call:

the parent continues executing from after fork()
the child also continues executing from after fork() — it's a full copy
They are now two separate processes. The OS owes both of them CPU time. How it interleaves them is the scheduler's job.

different from go where go routine is created  but if you dont use sleep or wait() main thread will run and execute

12:11 AM
Exactly, good connection to make. Here's the difference:

Go goroutines: The main goroutine is special — when it exits, the entire program dies, taking all goroutines with it.

go
func main() {
    go doSomething() // spawned but...
    // main exits here → program dead → goroutine never runs
}
So you need time.Sleep or sync.WaitGroup to keep main alive long enough.

Unix fork(): The child is a fully independent OS process — it doesn't die when the parent exits. The OS owes it CPU time regardless.

c
fork();
// parent exits → child keeps running
// OS adopts child under init/systemd (PID 1)
This is called an orphan process — parent died but child lives on, adopted by PID 1.

The fundamental difference:

Go goroutine	Unix process (fork)
Unit	lightweight thread	full OS process
Managed by	Go runtime	OS kernel
If parent dies	goroutine dies	child lives on
Scheduling	Go runtime scheduler	OS scheduler
Memory	shared heap	separate copy
Goroutines are more like threads than processes. A closer Unix equivalent to goroutines would be pthread_create() — where if the main thread exits, all threads die, same as Go.

so for pthreads if the parent thread finishes execution does the progam close like in go or the child thread still executes like we are seeing with processes

12:12 AM
Same as Go — program dies.

If the main thread finishes, the entire process exits, killing all pthreads with it.

c
#include <pthread.h>
#include <stdio.h>

void* worker(void* arg) {
    sleep(2);
    printf("worker done\n"); // may never print
    return NULL;
}

int main() {
    pthread_t t;
    pthread_create(&t, NULL, worker, NULL);
    // main returns here → process exits → worker thread dies
    return 0;
}
To keep main alive you use:

c
pthread_join(t, NULL); // main waits for thread to finish
The pattern is consistent across all three:

pthreads	Go goroutines	fork()
If main/parent exits	all threads die	all goroutines die	child lives on
Why	threads share the process	goroutines share the runtime	child is a separate process
Solution	pthread_join()	wg.Wait() / time.Sleep	wait() (optional)
The core reason:

Threads (pthreads, goroutines) live inside a process — they share its memory and lifetime. When the process dies, everything inside dies.

fork() creates a whole new process — completely independent, its own memory space, its own lifetime. The OS tracks it separately.

trying to udnerstand fork and exec()
so if i just use fork() i create a child process that runs the same program as parent

but if i want to run a new program, i first create a new process and then load the binary of the new program i want to run into the child process?

12:28 AM
Exactly right. You've got the mental model perfectly.

c
pid_t pid = fork();

if (pid == 0) {
    // we are in the child process
    // still running the same program as parent
    // now replace this process image with a new program
    execv("/bin/ls", args);
    // if exec succeeds, nothing below this line ever runs
    // the child is now running "ls"
}
```

---

**What exec() does:**

It doesn't create a new process. It **replaces** the current process's memory image — code, stack, heap — with a new program's binary. Same PID, new program.
```
fork()                          exec()
─────────────────────────────────────────────
parent (PID 27)                 parent (PID 27)
    │                               │
    └── child (PID 28)              └── child (PID 28)
        running same code               now running /bin/ls
        as parent                       completely replaced
Why fork first, then exec?

Because you want the parent to keep running. If you called exec() directly in the parent, the parent program would be replaced and gone.

So the pattern is always:

fork() → create child
exec() in child → replace child with new program
parent continues, optionally wait()s for child to finish
This is literally how your shell works every time you type a command.

if (!fork ())
return 1;
what does this mean?
if fork() returns 0

12:36 AM
!fork() means "if fork() returns 0" — which means we are in the child process.

!0 = 1 → condition is true → return 1 executes.

So it's just a compact way of writing:

c
pid_t pid = fork();
if (pid == 0) {
    // we are in the child
    return 1;
}
What happens to the parent?

fork() returns the child's PID (e.g. 28) in the parent. !28 = 0 → condition is false → parent skips the return 1 and continues executing.

Quick truth table:

who	fork() returns	!fork()	executes?
child	0	1 (true)	yes, returns 1
parent	28 (child PID)	0 (false)	no, continues
error	-1	1 (true)	yes — bug!
The subtle bug in this pattern is that -1 (fork failure) also makes !fork() true, so error case and child case are treated the same. More robust code checks explicitly:

c
pid_t pid = fork();
if (pid == -1) perror("fork failed");
else if (pid == 0) return 1; // child
int main (void)
{
int status;
pid_t pid;
if (!fork ())
return 1;
pid = wait (&status);
if (pid == −1)
perror ("wait");
printf ("pid=%d\n", pid);
if (WIFEXITED (status))
printf ("Normal termination with exit status=%d\n",
WEXITSTATUS (status));
if (WIFSIGNALED (status))
printf ("Killed by signal=%d%s\n",
WTERMSIG (status),
WCOREDUMP (status) ? " (dumped core)" : "");
if (WIFSTOPPED (status))
printf ("Stopped by signal=%d\n",
WSTOPSIG (status));
if (WIFCONTINUED (status))
printf ("Continued\n");
return 0;
}

with this fork()
does it mean while main process is running it is contexted switched so child process runs then child process returns before we continue with main()

because of the wait()

Show more
12:38 AM
Not exactly. Let me break it down carefully.

What actually happens:

main() starts
    │
    fork() called
    │
    ├── child process created
    │       !fork() = true → return 1 (child exits immediately)
    │
    └── parent continues
            wait(&status) ← parent blocks here
            waiting for ANY child to finish
            child already exited (or exits soon)
            wait() returns with child's PID and exit info
            parent continues printing...
To your specific question — is it a context switch?

Not necessarily. The scheduler decides. Three possible scenarios:

Scenario A — child runs first:

fork() → child runs → child returns 1 → child dead
       → parent runs → wait() returns immediately (child already done)
Scenario B — parent runs first:

fork() → parent runs → hits wait() → blocks
       → scheduler runs child → child returns 1
       → parent unblocks → wait() returns
Scenario C — they interleave:

fork() → parent runs a bit → context switch → child runs a bit → etc
All three are valid. wait() handles all of them — it blocks until child finishes regardless of order.

What the status macros tell you:

c
WIFEXITED(status)    // did child exit normally? (return or exit())
WEXITSTATUS(status)  // what was the exit code? → 1 (from return 1)

WIFSIGNALED(status)  // did child get killed by a signal? → no
WIFSTOPPED(status)   // was child stopped? → no
WIFCONTINUED(status) // was child continued? → no
```

So the output will be:
```
pid=28
Normal termination with exit status=1
Key insight:

wait() isn't about causing a context switch — it's about the parent declaring "I won't continue until my child is done". The OS handles the scheduling. The parent just blocks until the child's exit status is ready.

what does blocking and non blocking do while thread is using processor
if non blocking say read() or some syscall does the thread stay on the processor it is always context swtiched?

12:39 AM
Great question. Let me break it down.

Blocking call:

Thread makes a syscall (e.g. read()) that can't complete immediately.

thread calls read()
    │
    └── OS puts thread in WAITING state
        thread is removed from CPU
        CPU given to another thread/process
        
        ... disk fetches data ...
        
        data ready → OS moves thread back to READY queue
        scheduler picks it up → thread resumes
Thread is off the CPU the entire time it waits. It's not wasting cycles — it's just parked.

Non-blocking call:

Thread makes the same syscall but with non-blocking flag.

thread calls read() with O_NONBLOCK
    │
    └── data not ready? OS returns immediately with EAGAIN/EWOULDBLOCK
        thread stays on CPU
        thread decides what to do next
        (check again later, do other work, etc.)
Thread stays on CPU and keeps running. But now it's responsible for checking again later.

The tradeoff:

blocking	non-blocking
Data not ready	thread sleeps, off CPU	returns error immediately
CPU usage	efficient, not wasting cycles	thread must poll/retry
Complexity	simple	more complex logic
Good for	simple programs	event loops, servers
To your exact question — does the thread stay on CPU?

Blocking → no, always context switched off
Non-blocking → yes, stays on CPU, but gets back an error it must handle
This is why Node.js / async patterns exist:

A server handling 10,000 connections can't afford to block a thread per connection — that's 10,000 threads all sleeping waiting for I/O. Instead:

single thread
    │
    ├── read() non-blocking on socket 1 → not ready, move on
    ├── read() non-blocking on socket 2 → data ready, handle it
    ├── read() non-blocking on socket 3 → not ready, move on
    └── loop back...
This is the event loop model — one thread, never blocking, handling many connections. Go's goroutines and Node.js both work this way under the hood.

which means that thread will keep consumeing processor
does it get switched due to time slice or jsut stays indefinitely till it decides to leave

12:41 AM
It gets switched by the time slice — the thread doesn't get to stay indefinitely.

The OS is always in control:

Even if your thread never blocks, the scheduler has a timer (called a preemption timer). When the time slice expires, the OS forcibly kicks the thread off the CPU regardless of what it's doing.

thread running (non-blocking)
    │
    ├── doing work...
    ├── doing work...
    ├── doing work...
    │
    ⏰ timer interrupt fires (e.g. every 10ms)
    │
    └── OS preempts thread → saves its state → picks next thread
        ... other threads run ...
        eventually scheduler gives thread its turn again
        thread resumes exactly where it left off
Two ways a thread leaves the CPU:

how	who decides
Voluntary	blocking call, sleep(), yield()	the thread itself
Involuntary	time slice expires	the OS (preemption)
Non-blocking threads still get preempted involuntarily. They just never volunteer to leave.

This is the difference between cooperative and preemptive scheduling:

Cooperative (old Mac OS, early Windows) — threads run until they voluntarily yield. One greedy thread could freeze the whole system.
Preemptive (Linux, modern OS) — OS forcibly preempts threads. No thread can hog the CPU forever.
Linux is preemptive. Your non-blocking thread will get kicked off eventually.

So the polling pattern in practice:

c
while (1) {
    int n = read(fd, buf, size); // non-blocking
    if (n == -1 && errno == EAGAIN) {
        // not ready, do other work or just continue loop
        continue;
    }
    // handle data
}
This loop will keep getting time slices, getting preempted, getting time slices again — burning CPU checking over and over. That's why raw polling is wasteful, and why epoll/select exist — they let you say "wake me up when data is ready" without burning CPU in a loop.

so then if i have a server which is long running and we are using the event loop or non blocking like node foe example, how does that work?

12:42 AM
Good question. Let me walk through it.

The naive approach (blocking, one thread per connection):

connection 1 → thread 1 → read() → blocks → sleeping
connection 2 → thread 2 → read() → blocks → sleeping
connection 3 → thread 3 → read() → blocks → sleeping
...
connection 10000 → thread 10000 → sleeping
10,000 threads all sleeping, wasting memory, context switching overhead. Doesn't scale.

The event loop approach (Node.js):

One thread. Never blocks. Uses the OS to watch many file descriptors at once.

single thread running event loop
    │
    └── epoll_wait() ← give OS a list of fds to watch
        "wake me up when ANY of these have data"
        │
        thread blocks here (voluntarily, just once)
        │
        ... OS watches all sockets in kernel space ...
        │
        socket 3 has data → OS wakes thread
        │
        thread handles socket 3's data
        runs the callback
        goes back to epoll_wait()
The key syscall is epoll (Linux):

c
// register sockets you care about
epoll_ctl(epfd, EPOLL_CTL_ADD, socket_fd, &event);

// block until one of them is ready
int n = epoll_wait(epfd, events, MAX_EVENTS, -1);

// handle whichever ones are ready
for (int i = 0; i < n; i++) {
    handle(events[i].data.fd);
}
```

The OS does the watching in kernel space — efficiently. Your thread only wakes up when there's actual work to do.

---

**Node.js internals (libuv):**
```
your JS code
    │
    └── libuv (C library)
            │
            └── epoll (Linux) / kqueue (Mac) / IOCP (Windows)
                    │
                    └── kernel watching all sockets
When you write:

javascript
server.on('connection', (socket) => {
    socket.on('data', (data) => {
        // handle data
    });
});
Under the hood libuv registered those sockets with epoll. Your callbacks only fire when the OS says data is ready.

So how does the single thread never get overwhelmed?

It works well as long as callbacks are fast. The model breaks if you do:

javascript
socket.on('data', (data) => {
    // CPU heavy work - blocks the event loop
    const result = heavyComputation(data); // bad
});
That one callback holds the thread hostage. All other connections wait. This is the famous "don't block the event loop" rule in Node.

For CPU heavy work Node uses worker threads — offload computation, keep event loop free.

Summary:

blocking per thread	event loop
threads	one per connection	one (or few)
waiting	thread sleeps	OS watches via epoll
CPU	wasted on sleeping threads	only runs when work exists
weakness	doesn't scale	CPU heavy tasks block everything
examples	Apache (old), Java blocking IO	Node.js, Nginx, Go (hybrid)
Go is a interesting middle ground — it uses goroutines which look like blocking code but the runtime converts them to non-blocking under the hood, with its own scheduler on top of OS threads. Best of both worlds — simple code, scales well.

how does the worker threads for node work? are they assigned to other processors or what?

12:45 AM
Yes, worker threads can run on other CPU cores — that's the whole point.

What libuv does under the hood:

Node ships with a thread pool (libuv's thread pool, default 4 threads) that runs alongside the event loop.

process (one Node.js process)
    │
    ├── main thread (event loop)
    │       handles I/O, callbacks, JS execution
    │
    ├── worker thread 1  ← libuv thread pool
    ├── worker thread 2
    ├── worker thread 3
    └── worker thread 4
These are real OS threads — the OS scheduler can assign them to different CPU cores.

How work gets offloaded:

javascript
const { Worker } = require('worker_threads');

// main thread - event loop keeps running
const worker = new Worker('./heavy-task.js');

worker.on('message', (result) => {
    // result comes back here when done
    console.log(result);
});
```
```
main thread                    worker thread
    │                               │
    ├── receives connection         ├── running heavyComputation()
    ├── handles other requests      │       on separate core
    ├── sends more responses        │
    │                               └── posts result back
    └── callback fires ←──────────────── via message passing
```

---

**The OS and cores:**
```
CPU core 1          CPU core 2          CPU core 3
    │                   │                   │
main thread         worker thread 1     worker thread 2
(event loop)        (heavy task)        (heavy task)
The OS treats each thread independently and schedules them across available cores. True parallelism.

But there's a catch — no shared memory by default:

Worker threads in Node don't share the same JS heap. They communicate by message passing (copying data), not shared variables.

javascript
// main thread
worker.postMessage({ data: bigArray }); // copied, not shared

// worker thread
parentPort.on('message', (msg) => {
    const result = process(msg.data);
    parentPort.postMessage(result); // copied back
});
```

Copying large data is expensive. For that Node has `SharedArrayBuffer` — actual shared memory between threads, but then you have to handle race conditions yourself.

---

**libuv thread pool vs worker_threads:**

These are actually two different things:

| | libuv thread pool | worker_threads |
|---|---|---|
| who uses it | Node internals | your application code |
| what for | file I/O, DNS, crypto | CPU heavy JS work |
| you control it | not directly | yes, explicitly |
| default size | 4 threads | you create them |

When you call `fs.readFile()` — that goes to libuv's pool automatically. You never see it. `worker_threads` is the explicit API for your own CPU heavy work.

---

**The full picture:**
```
Node.js process
    │
    ├── V8 (JS engine) + event loop  ← your JS runs here
    │       │
    │       └── epoll watching all sockets
    │
    ├── libuv thread pool (4 threads)
    │       └── file I/O, DNS lookups etc
    │
    └── your worker threads (worker_threads module)
            └── CPU heavy JS tasks
Go's goroutine model essentially automates all of this — the runtime figures out which goroutine goes on which OS thread and core, transparently. Node makes you manage it more explicitly.

sof or event loops instead of assigning a thread to each connection which are basically files we can read and write
we just let the os monitor these connections and when data say a client send sth to the read buffer of the connection, it will awake the main thread to handle it?

12:48 AM
Exactly. That's the perfect mental model.

You nailed the key insight — connections are just files.

Every connection is a file descriptor (fd). Reading from a socket = reading from a file. The OS doesn't care.

The flow exactly as you described:

client sends data
    │
    └── data lands in socket's receive buffer (kernel space)
        │
        OS marks that fd as "readable"
        │
        epoll_wait() wakes up
        │
        main thread: "oh fd 42 has data, run its callback"
        │
        read(42, buf, size) → copies data from kernel buffer to userspace
        │
        your handler runs
        │
        back to epoll_wait() → sleeping again
The OS is doing the heavy lifting:

main thread                    kernel
    │                               │
    epoll_wait() ─── "watch       │
    sleeping        these fds" ──→ │ monitoring 10,000 fds
                                   │ efficiently in kernel space
    client 4872 sends data ──────→ │
                                   │ marks fd as ready
    wakes up ←─────────────────── │
    handles fd                     │
    back to sleep                  │
Your thread is only awake when there is actual work to do. Everything else the kernel handles silently.

Why this is efficient:

The old model — one thread per connection:

10,000 connections = 10,000 threads
most sleeping, wasting memory
OS context switching between them constantly
Event loop model:

10,000 connections = 1 thread + kernel watching 10,000 fds
thread only wakes up when data arrives
zero wasted context switches
The only job of the event loop is:

tell kernel what to watch → epoll_ctl()
sleep until something is ready → epoll_wait()
handle what's ready → run callback
repeat
That's it. Simple loop, scales to hundreds of thousands of connections on a single thread.

so there is no need for thread context switchign and constanting savign and rereading of data from thread stack and tcb

12:49 AM
Exactly. That's the core efficiency win.

One thread, no context switching overhead at all.

What you eliminate:

With one thread there is nothing to switch between. No:

saving registers to TCB
saving stack pointer
flushing/reloading CPU cache
kernel deciding who runs next
TLB invalidation
The CPU cache stays warm with your event loop's data. Same thread, same memory, running continuously.

The cost of context switching is real:

context switch cost
    │
    ├── save registers (~20-30 registers)
    ├── save program counter
    ├── save stack pointer
    ├── update TCB
    ├── pick next thread (scheduler runs)
    ├── load new TCB
    ├── restore registers
    ├── TLB flush (if different process)
    └── cache is now cold → cache misses on next memory access
Each switch is roughly 1-10 microseconds. Tiny alone, but at 10,000 threads switching constantly it adds up significantly.

The tradeoff though:

You trade context switching overhead for callback complexity:

blocking (simple to reason about)        event loop (complex to reason about)
─────────────────────────────────────    ────────────────────────────────────
read(fd)  ← just wait here              epoll_wait()
// data is here now                      callback fires
do work                                  but where is your state?
                                         you have to manage it manually
With blocking threads your stack is your state — local variables just sit there while you sleep. Natural.

With event loop your stack is gone between callbacks — you have to explicitly carry state forward via closures, promises, or state machines. That's why async/await was invented — to make event loop code look like blocking code.

This is also why Go feels elegant:

Go gives you the simple mental model of blocking code:

go
// looks blocking, simple to reason about
data := read(conn)
process(data)
But under the hood the Go runtime converts it to non-blocking I/O + epoll, switches goroutines instead of OS threads, no kernel context switch overhead.

You get both — simple code AND event loop efficiency. That's the real innovation of the Go runtime.

on the issue of non blocking what if my code is like
read() // non blokcing

do somethinf with data you read

what happens if read returns nothing or an error

7:09 AM
Good question. That's exactly the bug most beginners write.

What actually happens:

c
int n = read(fd, buf, size); // non-blocking

// n == -1, errno == EAGAIN → no data yet
// but you just fall through and use buf anyway
doSomething(buf); // bug — buf has garbage or stale data
You're using uninitialized or stale data. Silent bug, hard to track down.

You must always check the return value:

c
int n = read(fd, buf, size);

if (n == -1 && errno == EAGAIN) {
    // no data yet, come back later
    return;
}

if (n == -1) {
    // actual error
    perror("read failed");
    return;
}

if (n == 0) {
    // connection closed by client
    close(fd);
    return;
}

// only here is buf actually valid
doSomething(buf, n);
The three return values you must handle:

return value	meaning
-1 + EAGAIN	no data yet, try later
-1 + other errno	real error
0	connection closed
> 0	n bytes read, buf is valid
This is exactly why event loop model is better than manual polling:

Instead of you writing:

c
while (1) {
    n = read(fd, buf, size); // keep trying
    if (n == -1 && errno == EAGAIN) continue; // burning CPU
}
You let epoll tell you when data is ready:

c
epoll_wait(epfd, events, MAX, -1); // sleep until ready
// only wake up when read() will actually return data
n = read(fd, buf, size); // now guaranteed to have data
epoll eliminates the EAGAIN case almost entirely — by the time you call read(), the OS already confirmed data is there.

while (1) {
    n = read(fd, buf, size); // keep trying
    if (n == -1 && errno == EAGAIN) continue; // burning CPU
}

got it. so in this scenario if read returns the number of bytes read when do we break out of the loop?

7:11 AM
Good catch. You need to break out when data is successfully read:

c
while (1) {
    n = read(fd, buf, size);
    
    if (n == -1 && errno == EAGAIN) {
        continue; // no data yet, keep trying
    }
    
    if (n == -1) {
        perror("read failed"); // real error
        break;
    }
    
    if (n == 0) {
        printf("connection closed\n"); // client disconnected
        break;
    }
    
    // n > 0, we have data
    doSomething(buf, n);
    break; // got what we needed, exit loop
}
But it's actually more nuanced than that.

read() doesn't guarantee it reads everything in one call. Say you're expecting 1024 bytes — it might return 200, then 300, then 524. This is called partial reads.

So in practice you loop until you've read everything you expect:

c
int total = 0;
int expected = 1024;

while (total < expected) {
    n = read(fd, buf + total, expected - total);
    
    if (n == -1 && errno == EAGAIN) {
        continue; // no data yet
    }
    
    if (n == -1) {
        perror("read failed");
        break;
    }
    
    if (n == 0) {
        printf("connection closed\n");
        break;
    }
    
    total += n; // accumulate bytes read
}

// now buf has all 1024 bytes
doSomething(buf, total);
This is why protocols define message boundaries:

Without knowing how many bytes to expect, you don't know when to stop reading. So protocols solve this explicitly:

protocol	how it signals end of message
HTTP	Content-Length header tells you exact bytes
HTTP chunked	each chunk prefixed with its size
Redis	\r\n terminates each command
custom	first 4 bytes = message length, then read that many
Your loop reads until that boundary is hit — not just until read() returns something.

so it means either you use the protocols way to find out how many bytes to read or you implement readall() which loops continuously till EOF is returned especially for sockets
but for files on disk the file length is always known

7:14 AM
Exactly right on all three points.

For network sockets — you have two options:

Option 1: Protocol defined boundary

c
// read header first (fixed 4 bytes = message length)
read(fd, header, 4);
int msg_len = parse_length(header);

// now you know exactly how much to read
while (total < msg_len) {
    n = read(fd, buf + total, msg_len - total);
    total += n;
}
Option 2: Read until connection closes (EOF)

c
while (1) {
    n = read(fd, buf + total, size);
    
    if (n == 0) break; // EOF → client closed connection
    if (n == -1) { perror("error"); break; }
    
    total += n;
}
// now buf has everything
But option 2 has a problem — you have to wait for client to close the connection to know message is done. Not practical for long-lived connections like HTTP keep-alive where client stays connected and sends multiple requests.

That's exactly why every serious protocol defines message boundaries explicitly.

For files on disk — you're right, length is always known:

c
// stat() gives you file size upfront
struct stat st;
stat("file.txt", &st);
int file_size = st.st_size; // exact bytes

// now read exactly that much
while (total < file_size) {
    n = read(fd, buf + total, file_size - total);
    total += n;
}
Or you just read until n == 0 and it works fine because a file has a definite end unlike a socket which stays open.

The fundamental difference:

file on disk	network socket
size known upfront	yes, via stat()	no
EOF means	end of file, definite	client closed connection
message boundary	the file itself	protocol must define it
partial reads	still happen	happen more frequently
This is why HTTP/1.1 made Content-Length mandatory for responses with a body — so the client knows exactly when the message ends without closing the connection, allowing keep-alive and pipelining multiple requests over one socket.