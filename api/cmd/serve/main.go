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
	http.ListenAndServe("localhost:8100", &cgi.Handler{Path: "./serv"})
}
