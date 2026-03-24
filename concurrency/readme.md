CONCURRENCY
-----------
Process -> Threads -> user level threads
We start with processes but they are expensive and can be time consuming so we use threads
Threads are lighter or less expensive than processes but they introduce race conditions
We use synchronization methods to tackle race conditions
These sometimes lead to deadlocks and spinlocks?
We also have threading models: 1:N and M:N


2. Dealing with threads
When processes do need to communicate and synchronize with each other, we program them to
use operating system tools and other applications such as files, databases, pipes, sockets etc since
processes don't share memory.
Processes use the OS as a communication manager- the objects like files, buffers, pipes, etc used to handle
the communication are all data structures and state managed by the OS.

3.1 Sharing Memory
It is similar to communicating with a friend but instead of exchanging messages you use a whiteboard
to exchange ideas or messages using symbols and abstractions

we allocate a part of a process' memory for example a shared data structure or variable and have different
threads (goroutines) work concurrently on this structure
we then context switch between threads and allow them to work with the memory as they please

Computers may have multiple processors that share a system bus and main memory. Before a processor uses
the bus, it must make sure the bus is free and not in use by another processor before making a request
for a memory location and goes back to listening waiting for a reply
As we scale processors, the bus becomes a bottle neck so we implement various layers
of caches between the cpu and main memory to reduce the load on the system bus
This architecture comes with new problems such has when we have 2 threads executing in parallel on
different processors, if thread 1 reads some value into its cache and updates it, the results in the cache
or value of the variable will be different from value in main memory
If thread 2 decides to use this value and reads it from main memory it will have an outdated value
One solution is to perform a write through or propagate updates made in cache content back to main memory
How ever if the situation is that thread 2 has an outdated copy of the value in another cache then we can
make caches listen to bus memory for update messages and if it notices an update to the memory that it replicates
in the cache, it either applies the update or invalidates the cache content
If we invalidate the content, if the thread needs the variable or data it can fetch the updated value from main memory
# study cache coherency protocols

Processor and main memory have a request response cycle. it asks for data using some form of identifier. Variables
map human readable names to memory addresses. Allows us to index these cells for easy retrieval and updates
We also have bus architectures which affect system performance and latency

Project- 
build a simple project of 2 processes communicating with queues or pipes and another reading from it
to write to a db server




There are 2 ways of communication:
read 43. Interprocess communication
- sharing memory
- message passing
sharing memory is used by threads in a process that share the same memory space
message passing is used by processes with isolated memory space
in distributed systems, applications or processes send messages to eachother  via communication protocols
TCP is a byte stream protocol so the application layer is responsible for the protocol (rules) that 
determine message framing from a stream of bytes


Pipes are the oldest UNIX inter process communication tool
A pipeline is basically a question about how do we get two or more processes running different programs to communicate by allowing the output of one process to become the input of another process.
Think about this when you multiprocessing and fork()
A pipeline basically allows us to move data between two or more processes where the output of one process serves
as an input of another process
Pipes are used for related processes: parent or similar ancestor: fork() and exec() was used in some timeline
process1 ---pipeline or stream---- process2
this is the logical plan of a job: source process ------stream------ transform process

question: fork() copies almost everything aside optimisations like copy on write, etc. can I modify the code of
child process?