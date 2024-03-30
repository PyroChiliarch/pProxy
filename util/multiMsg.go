package util

import (
	"strconv"
	"strings"

	"github.com/google/uuid"
)

// Client get s a new message ID with start
// Client sends part of the full message as base32 encoded, indexed by number of part and with the message identifier

// /http/client/dorequest/start
// /http/client/dorequest/msg/7cd344ab-1fb2-4a50-8894-a2a97474f68c/2/3tajjwgustmrjfgiycknzteu3dkjjxgmstmmjfgzccknrvej6x2===
// /http/client/dorequest/end

// /msg/<requestID>/<msgID>/<msgPart>

func StartMsg(requestCache map[uuid.UUID]map[int]string) uuid.UUID {
	//Get new ID
	id := uuid.New()

	//Store ID
	requestCache[id] = make(map[int]string) //Make a new map to store incoming requests

	//Give ID back for future requests
	return id
}

func PartMsg(requestCache map[uuid.UUID]map[int]string, urlPath string) error {

	//Parse URL Values
	values := strings.Split(urlPath, "/")

	//Get the request ID (UUID)
	requestID, err := uuid.Parse(values[len(values)-3])
	if err != nil {
		//fmt.Fprintf(w, err.Error())
		return err
	}

	//Get the Message ID (Int, increments)
	msgID, err := strconv.Atoi(values[len(values)-2])
	if err != nil {
		//fmt.Fprintf(w, err.Error())
		return err
	}

	//Get Data, string, no error checking needed
	data := values[len(values)-1]

	//Store data, Data is part of a Base32 string
	requestCache[requestID][msgID] = data

	//No errors, return nil
	return nil

}

func EndMsg(requestCache map[uuid.UUID]map[int]string, urlPath string) (msg string, err error) {

	values := strings.Split(urlPath, "/")

	//Request ID is the last value, its the only value thats needed
	requestID, err := uuid.Parse(values[len(values)-1])
	if err != nil {
		return //Return just the error
	}

	//Rebuild data from each request
	msg = ""
	for i := 0; i < len(requestCache[requestID]); i++ { //Maps do not have a guaranteed order, need to loop manualy, iterators probably will cause messag to be in wrong order
		msg = msg + requestCache[requestID][i]
	}

	//Remove the request from the cache
	delete(requestCache, requestID)

	//Change it back to upper case
	//msg is still base32 encoded
	msg = strings.ToUpper(msg)

	//Return full msg, its encoded as Base32
	return //Return the msg, err is nil
}
