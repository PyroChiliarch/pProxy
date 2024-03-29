package pProxyWeb

import (
	"encoding/base32"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// Supports got http and https
// Simple http get
func HttpGet(w http.ResponseWriter, r *http.Request) {

	//Print to log
	println(r.RemoteAddr + ": get request")

	//Read url parameter from url
	encodedUrl := r.URL.Query().Get("url")
	urlBytes, err := base32.StdEncoding.DecodeString(strings.ToUpper(encodedUrl)) //Pico can only do lower case requests ,need to chagne to upper case
	if err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}
	url := string(urlBytes[:])

	// Return nothing if no url
	if url == "" {
		fmt.Fprintf(w, "ERROR-NOURL")
		return
	}

	//Make a proxied get request
	resp, err := http.Get(url)
	if err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}

	encodedResponse := base32.StdEncoding.EncodeToString(body)

	//Pass the response back
	fmt.Fprintf(w, encodedResponse)
}
