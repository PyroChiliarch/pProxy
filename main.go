package main

import (
	"fmt"
	"log"
	"net/http"
)

////////////////////////////////////////////////////////////////////
// Start a new HTTP Server

func defaultHandler(w http.ResponseWriter, r *http.Request) {
	println(r.RemoteAddr + ": Invalid Endpoint : " + r.URL.String())
	fmt.Fprintf(w, "Invalid Endpoint")
}

func main() {
	//Catch all, return error
	http.HandleFunc("/", defaultHandler)

	//Utils
	http.HandleFunc("/version", getVersion)

	//Basic http
	http.HandleFunc("/http/get", httpGet)

	//Client stuff
	http.HandleFunc("/http/client/new", httpClientNew)
	http.HandleFunc("/http/client/getjar/*", httpClientGetJar)

	//Http with client
	http.HandleFunc("/http/client/dorequest/start/*", httpClientDoRequestStart)
	http.HandleFunc("/http/client/dorequest/msg/*", httpClientDoRequestMsg)
	http.HandleFunc("/http/client/dorequest/end/*", httpClientDoRequestEnd)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
