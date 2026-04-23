package main

import (
	"fmt"
	"net/http"
	"sync"
)

type topMuxEntry struct {
	h http.Handler
}

// a server receives a request, parses it and then calls a handler to do some work
// toyServeMux is a router
// its goal is simple. map strings or urls or paths to handlers that can be called to do some work
// we need a way to add these info when it is supplied
// we need a way to look info and return the appropriate handler when a request comes in
type toyServeMux struct {
	mu          sync.Mutex
	routemapper map[string]http.Handler
}

func NewToyServeMux() *toyServeMux {
	return &toyServeMux{
		routemapper: make(map[string]http.Handler),
	}
}

// write to map or data structure for future lookups
func (t *toyServeMux) Handle(pattern string, controller http.Handler) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.routemapper[pattern] = controller

}

// lookup and retrieve data from data store based on some key or data
func (t *toyServeMux) Handler(r *http.Request) http.Handler {
	requestUrl := r.URL.Path
	t.mu.Lock()
	defer t.mu.Unlock()
	requestHandler := t.routemapper[requestUrl]

	return requestHandler

}

// Required: dispatches requests to the right handler
func (t *toyServeMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.mu.Lock()
	handler, exists := t.routemapper[r.URL.Path]
	t.mu.Unlock()

	if !exists {
		http.NotFound(w, r)
		return
	}
	handler.ServeHTTP(w, r)
}

// thin wrapper around my Handler method
// func (t *toyServeMux) ServeHTTP(w http.ResponseWriter, incomingRequest *http.Request) {
// 	r := t.Handler(incomingRequest)
// 	if r != nil {
// 		http.NotFound(w, incomingRequest)
// 	}
// 	r.ServeHTTP(w, incomingRequest)
// }

// create sample handler
type homeHandler struct{}

func (h *homeHandler) ServeHTTP(w http.ResponseWriter, incomingRequest *http.Request) {
	fmt.Fprintf(w, "Welcome Home")

}

func main() {

	// 1. Create mux and register routes FIRS
	serverMux := NewToyServeMux()
	homeController := &homeHandler{}
	serverMux.Handle("/", homeController)

	// 2. Create server with the custom mux as Handler
	server := &http.Server{
		Addr:    "localhost:8000",
		Handler: serverMux,
	}
	// 3. Start server (this blocks, so do it last)
	fmt.Println("Server running on http://localhost:8000")
	if err := server.ListenAndServe(); err != nil {
		fmt.Println("Server error:", err)
	}

}

/*

a server is some object that creates a listening socket, binds it to a port to handle
incoming client connections by parsing these connections and then handing the data over
to some handler or function to service the connection
if this server happens to be a server handling concurrent connections and multiple routes
there is a need for some object or logic which will route or multiplex connections to the
right handler based on some logic
For separation of concerns we use an object called a router or multiplexer
It will need a way or an api to write or persist routing data and an api to also return the
correct handler based on some data provided
It also needs a data structure to hold its data and this structure needs to be thread safe
we can also have structures for storing accounting and meta data ( extend it)



*/

/*
	func (mux *ServeMux) Handle(pattern string, handler Handler)
	Handle registers the handler for the given pattern. If the given pattern
	conflicts with one that is already registered or if the pattern is invalid, Handle panics.


	func (mux *ServeMux) Handler(r *Request) (h Handler, pattern string)
		Handler returns the handler to use for the given request, consulting r.Method,
		r.Host, and r.URL.Path. It always returns a non-nil handler. If the path is not in
		its canonical form,the handler will be an internally-generated handler that redirects
		to the canonical path. If the host contains a port, it is ignored when matching handlers.

*/
