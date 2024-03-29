package main

import (
	"fmt"
	"log"
	"net/http"
	"pProxy/pProxyWeb"
	"strconv"
)

// //////////////////////////////////////////////////////////////////////////////////////////////////////////
//
//	Generic Endpoint handlers
//
// //////////////////////////////////////////////////////////////////////////////////////////////////////////

func defaultHandler(w http.ResponseWriter, r *http.Request) {
	println(r.RemoteAddr + ": Invalid Endpoint : " + r.URL.String())
	fmt.Fprintf(w, "Invalid Endpoint")
}

func getVersion(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, version)
	println(r.RemoteAddr + ": version request")
}

func main() {

	///////////////////// Load Generic Endpoints /////////////////////
	//Catch all, return error
	http.HandleFunc("/", defaultHandler)

	//Utils
	http.HandleFunc("/version", getVersion)

	///////////////////// Load Web Endpoints /////////////////////
	//Basic http
	http.HandleFunc("/http/get", pProxyWeb.HttpGet)

	//Client stuff
	http.HandleFunc("/http/client/new", pProxyWeb.HttpClientNew)
	http.HandleFunc("/http/client/getjar/*", pProxyWeb.HttpClientGetJar)

	//Http with client
	http.HandleFunc("/http/client/dorequest/start/*", pProxyWeb.HttpClientDoRequestStart)
	http.HandleFunc("/http/client/dorequest/msg/*", pProxyWeb.HttpClientDoRequestMsg)
	http.HandleFunc("/http/client/dorequest/end/*", pProxyWeb.HttpClientDoRequestEnd)

	///////////////////// Load Game Endpoints /////////////////////

	///////////////////// Start Server /////////////////////
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(port), nil))

}
