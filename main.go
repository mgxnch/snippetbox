package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
)

// home is a function haldner which writes a byte slice
// as the response body
func home(w http.ResponseWriter, r *http.Request) {
	// Prevent non-existent paths from matching on this handler
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	w.Write([]byte("hello world!"))
}

func snippetView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get(("id")))
	if err != nil || id < 1 {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	// Fprintf can even take in http.ResponseWriter as an io.Writer, how cool!
	fmt.Fprintf(w, "displaying snippet of ID: %d", id)
}

func snippetCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		// ref: https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Allow
		w.Header().Set("Allow", "POST")
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("method not allowed"))
		return
	}
	w.Write([]byte("create a new snippet"))
}

func main() {
	// Initialise mux and declare the routes
	mux := http.NewServeMux()
	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippet/view", snippetView)
	mux.HandleFunc("/snippet/create", snippetCreate)

	log.Println("Starting server on http://localhost:4000")
	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}
