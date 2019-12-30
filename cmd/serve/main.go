// Test server for the SERV portal CGI handler.
//
// This program listens on HTTP port 8100, and redirects all requests to the
// SERV portal CGI server (invoked as "./portal").
package main

import (
	"net/http"
	"net/http/cgi"
)

func main() {
	var cgiHandler = cgi.Handler{Path: "./portal"}
	var fsHandler = http.FileServer(http.Dir("."))
	http.Handle("/", &cgiHandler)
	http.Handle("/portal.css", fsHandler)
	http.Handle("/portal.js", fsHandler)
	http.ListenAndServe("localhost:8100", http.DefaultServeMux)
}
