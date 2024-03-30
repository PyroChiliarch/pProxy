package util

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

// //////////////////////////////////////////////////////////////////////
// Return Message packaging
// Messages to Picotron get packaged in json so we can pass error as a bit of metadata
type returnMessage struct {
	Content string
	Err     bool
}

// Return a message to the client
func ReturnMessage(_writer http.ResponseWriter, _request *http.Request, _content string, _err error) {
	msg := new(returnMessage)

	// Handle Error messages
	// Set err to true, put error value in content, client will see error
	if _err != nil {

		msg.Err = true
		msg.Content = _err.Error()
		println(_request.RemoteAddr + ": ERROR: " + _err.Error())

	} else {
		//No error, err is false, content is message
		msg.Err = false
		msg.Content = _content
		println(_request.RemoteAddr + ": RESPONED TO REQUEST: " + _request.RequestURI)
	}

	//Transmit struct as a json string to be reconstructed at the other end
	jsonMsgBytes, err := json.Marshal(msg)
	if err != nil {
		println("Something went really wrong: ReturnMessage() marshalling json")
		println("content: " + msg.Content)
		println("error: " + strconv.FormatBool(msg.Err))
		return
	}

	//Return crafted message to client
	fmt.Fprintf(_writer, string(jsonMsgBytes))
}
