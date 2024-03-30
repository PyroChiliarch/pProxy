package pProxyWeb

import (
	"encoding/base32"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"pProxy/util"

	"github.com/google/uuid"
)

func HttpClientNew(w http.ResponseWriter, r *http.Request) {
	//Create new client
	jar := NewJar()
	c := http.Client{Jar: jar} // Give the client a cookie jar
	id := uuid.New()

	//Store the client
	httpClients[id] = c

	//Let Picotron know the id of their new client
	util.ReturnMessage(w, r, id.String(), nil)
	//util.ReturnMessage(w, r, id.String(), errors.New("Unknown error"))
}

// //////////////////////////////////////////////////////////////////////////////////////////////////////////
//
//	Work With Cookies
//
// //////////////////////////////////////////////////////////////////////////////////////////////////////////

func HttpClientGetJar(w http.ResponseWriter, r *http.Request) {

	values := strings.Split(r.URL.Path, "/")

	// Client ID is the second last value
	clientID, err := uuid.Parse(values[len(values)-2])
	if err != nil {
		util.ReturnMessage(w, r, "", err)
		return
	}

	// Url is the last value, decode it into its original string (in bytes)
	urlString, err := base32.StdEncoding.DecodeString(strings.ToUpper(values[len(values)-1]))
	if err != nil {
		util.ReturnMessage(w, r, "", err)
		return
	}

	// Change url string into a url object
	myUrl, err := url.Parse(string(urlString))
	if err != nil {
		util.ReturnMessage(w, r, "", err)
		return
	}

	//Get Client
	client := httpClients[clientID]

	//Get cookies from client
	cookies := client.Jar.Cookies(myUrl)
	jsonBytes, err := json.Marshal(cookies)
	if err != nil {
		util.ReturnMessage(w, r, "", err)
		return
	}

	///////////////////////////////////////////////////////////////////////////////////////////////////////////////////
	//Send cookie jar back
	util.ReturnMessage(w, r, string(jsonBytes), nil)
}

// //////////////////////////////////////////////////////////////////////////////////////////////////////////
//
//	Handle Adv Requests (Uses Clients and Requests to do custom requests)
//
// //////////////////////////////////////////////////////////////////////////////////////////////////////////

func HttpClientDoRequestStart(w http.ResponseWriter, r *http.Request) {

	//Eg /http/client/dorequest/start

	//Start the request, return request ID to client
	id := util.StartMsg(httpMultiMsgsCache)

	//Send new request ID to client for subsequent requests
	util.ReturnMessage(w, r, id.String(), nil)

}

func HttpClientDoRequestMsg(w http.ResponseWriter, r *http.Request) {

	//Eg /http/client/dorequest/msg/7cd344ab-1fb2-4a50-8894-a2a97474f68c/2/3tajjwgustmrjfgiycknzteu3dkjjxgmstmmjfgzccknrvej6x2===

	//Add the message part to the full request
	err := util.PartMsg(httpMultiMsgsCache, r.URL.Path)
	if err != nil {
		util.ReturnMessage(w, r, "", err)
		return
	}

	util.ReturnMessage(w, r, "OK", nil)
}

func HttpClientDoRequestEnd(w http.ResponseWriter, r *http.Request) {

	//Client has sent all their data and wants to do the request and get their data
	// Eg: /http/client/dorequest/end/7cd344ab-1fb2-4a50-8894-a2a97474f68c
	encodedData, err := util.EndMsg(httpMultiMsgsCache, r.URL.Path)

	////////////////////////////////////////////// Before this gets put in util

	//Data is still encoded in base32, and in json
	msg, err := base32.StdEncoding.DecodeString(encodedData)
	if err != nil {
		util.ReturnMessage(w, r, "", err)
		return
	}

	//////////// Unmarshal json
	// Make vars needed

	byteData := []byte(msg)         // Data as byte[]
	var data map[string]interface{} // Empty map

	// Unmarshal json, return on error
	if err := json.Unmarshal(byteData, &data); err != nil {
		util.ReturnMessage(w, r, "", err)
		return
	}

	//Get request from json
	clientID, err := uuid.Parse(data["client"].(string))
	if err != nil {
		util.ReturnMessage(w, r, "", err)
		return
	}

	//Unpack the request from json
	request := data["request"].(map[string]interface{})
	method := request["method"].(string)
	url := request["url"].(string)
	body := strings.NewReader(request["body"].(string))
	headers := make(map[string]interface{})
	if _, ok := request["headers"].(map[string]interface{}); ok { // Handle headers being empty, will be a different interface type ([]interface {}, not map[string]//interface {})
		headers = request["headers"].(map[string]interface{})
	}

	client := httpClients[clientID]

	//Craft new Request object
	newRequest, err := http.NewRequest(strings.ToUpper(method), url, body)

	//Set headers, values are in format "name:data"
	for _, v := range headers {
		values := strings.Split(v.(string), ":")
		newRequest.Header.Set(values[0], values[1])
	}

	//Do the newly reconstructed request
	httpResponse, err := client.Do(newRequest)
	if err != nil {
		util.ReturnMessage(w, r, "", err)
		return
	}

	//Get body of response
	resBodyStr, err := io.ReadAll(httpResponse.Body)
	if err != nil {
		util.ReturnMessage(w, r, "", err)
		return
	}

	//Build a response to send to Picotron
	picoResponse := make(map[string]interface{})
	picoResponse["body"] = string(resBodyStr) //Passing directly will give base64 on pico side, json.Marshall sends []byte as base64 because []byte != string
	picoResponse["status"] = httpResponse.StatusCode

	picoResponse["headers"] = make(map[string]string)

	//Headers are setup to have one name, multiple values in a weird nested array structure
	//Theres usually (in all my limited testing) only one value, so I pull it out
	for name, valueSlice := range httpResponse.Header {
		for _, value := range valueSlice {

			_, ok := picoResponse["headers"].(map[string]string)[name]
			if !ok {
				picoResponse["headers"].(map[string]string)[name] = value // No duplicates, add value
			} else {
				// Everything below here is for duplicate values
				// Its supposed to be possible to have more than one value for a header, never seen it though
				// Maybe this happens when the header is specified twice in the response?
				// If there are somehow duplicates, Just continually append a number
				// 50 sounds like a reasonable maximum number
				for i := 0; i < 50; i++ {
					_, ok := picoResponse["headers"].(map[string]string)[name+"_"+strconv.Itoa(i)]
					if !ok {
						// Up to a non duplicate value! yay!, set and exit loop
						//Never actually tested this code, may not work
						picoResponse["headers"].(map[string]string)[name] = value
						break
					}

				}
			}

		}
	}

	//Format response as json
	jsonString, err := json.Marshal(picoResponse)
	if err != nil {
		util.ReturnMessage(w, r, "", err)
		return
	}
	println(r.RemoteAddr + ": Client Request end")

	//Send response to client
	util.ReturnMessage(w, r, string(jsonString), nil)
}
