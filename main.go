package main

import (
	//"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"pProxy/pProxyGame"
	"pProxy/pProxyWeb"
	//_ "github.com/mattn/go-sqlite3"
)

// //////////////////////////////////////////////////////////////////////////////////////////////////////////
//
//	Generic Endpoint handlers
//
// //////////////////////////////////////////////////////////////////////////////////////////////////////////

func defaultHandler(w http.ResponseWriter, r *http.Request) {
	println(r.RemoteAddr + ": Invalid Endpoint : " + r.URL.String()) // Bots spam this
	fmt.Fprintf(w, "Invalid Endpoint")
}

func getVersionHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, version)
	println(r.RemoteAddr + ": version request")
}

func main() {

	args := os.Args[1:]
	server := new(http.Server)
	mux := new(http.ServeMux)

	// //////////////////////////////////////////////////////////////////////////////////////////////////////////
	//
	//	Database Setup
	//
	// //////////////////////////////////////////////////////////////////////////////////////////////////////////

	pProxyGame.InitDatabase()

	// //////////////////////////////////////////////////////////////////////////////////////////////////////////
	//
	//	Register Endpoints
	//
	// //////////////////////////////////////////////////////////////////////////////////////////////////////////

	///////////////////// Register Generic Endpoints /////////////////////
	//Catch all, return error
	mux.HandleFunc("/", defaultHandler)

	//Utils
	mux.HandleFunc("/version", getVersionHandler)

	///////////////////// Register Web Endpoints /////////////////////
	//Basic http
	//mux.HandleFunc("/http/get", pProxyWeb.HttpGetHandler)

	//Client stuff
	mux.HandleFunc("/http/client/new", pProxyWeb.HttpClientNewHandler)
	mux.HandleFunc("/http/client/getjar/*", pProxyWeb.HttpClientGetJarHandler)

	//Http with client
	mux.HandleFunc("/http/client/dorequest/start/*", pProxyWeb.HttpClientDoRequestStartHandler)
	mux.HandleFunc("/http/client/dorequest/msg/*", pProxyWeb.HttpClientDoRequestMsgHandler)
	mux.HandleFunc("/http/client/dorequest/end/*", pProxyWeb.HttpClientDoRequestEndHandler)

	///////////////////// Register Game Endpoints /////////////////////
	mux.HandleFunc("/game/register/*", pProxyGame.RegUserHandler)
	//mux.HandleFunc("/game/login/*", pProxyGame.Login)
	//http.HandleFunc("/game/lobby/getgames")

	//Apply the mux to the server
	server.Handler = mux

	// //////////////////////////////////////////////////////////////////////////////////////////////////////////
	//
	//	Start server based on args HTTP/HTTPS
	//
	// //////////////////////////////////////////////////////////////////////////////////////////////////////////

	if len(args) == 2 {
		//HTTP, user only specified IP and PORT
		ip := args[0]
		port := args[1]

		server.Addr = ip + ":" + port

		//Start Server
		log.Fatal(server.ListenAndServe())

	} else if len(args) == 4 {
		//HTTPS, user specified key and cert files
		ip := args[0]
		port := args[1]
		certLocation := args[2]
		keyLocation := args[3]

		/*
			// Load the cert and key file
			TLSCert, err := tls.LoadX509KeyPair(certLocation, keyLocation)
			if err != nil {
				log.Fatalf("Error loading certificate and key file: %v", err)
			}

			tlsConfig := &tls.Config{
				Certificates: []tls.Certificate{TLSCert},
			}*/

		server.Addr = ip + ":" + port
		//server.TLSConfig = tlsConfig

		//Start Server with TLS
		log.Fatal(server.ListenAndServeTLS(certLocation, keyLocation))
	}

}
