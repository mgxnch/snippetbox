package main

import (
	"fmt"
	"html/template"
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

	files := []string{
		"./ui/html/base.tmpl",
		"./ui/html/pages/home.tmpl",
		"./ui/html/partials/nav.tmpl",
	}

	// ts stands for Template Set
	ts, err := template.ParseFiles(files...)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	// ExecuteTemplate writes the output of the parsed template into the writer w
	// We have to specify the named template to parse and apply
	err = ts.ExecuteTemplate(w, "base", nil)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
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
