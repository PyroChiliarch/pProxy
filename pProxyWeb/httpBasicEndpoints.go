package pProxyWeb

import (
	"encoding/base32"
	"errors"
	"io"
	"net/http"
	"pProxy/util"
	"strings"
)

// Supports got http and https
// Simple http get
func HttpGet(w http.ResponseWriter, r *http.Request) {

	//Read url parameter from url
	encodedUrl := r.URL.Query().Get("url")
	urlBytes, err := base32.StdEncoding.DecodeString(strings.ToUpper(encodedUrl)) //Pico can only do lower case requests ,need to chagne to upper case
	if err != nil {
		util.ReturnMessage(w, r, "", err) // Catch error, send it to client
		return
	}
	url := string(urlBytes[:])

	// Return nothing if no url
	if url == "" {
		util.ReturnMessage(w, r, "", errors.New("URL is empty")) // Catch error, send it to client
		return
	}

	//Make a proxied get request
	resp, err := http.Get(url)
	if err != nil {
		util.ReturnMessage(w, r, "", err) // Catch error, send it to client
		return
	}

	//Get body of the response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		util.ReturnMessage(w, r, "", err) // Catch error, send it to client
		return
	}

	//Send body of response to picotron client
	util.ReturnMessage(w, r, string(body), nil)

}
