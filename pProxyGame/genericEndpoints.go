package pProxyGame

import (
	"encoding/base32"
	"fmt"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

// //////////////////////////////////////////////////////////////////////////////////////////////////////////
//
//	Structure
//	Main Chat
//	List lobbies
//	Create lobby
//	Join lobbies
//	Lobby Chat
//  Lobby Options
//  Start game (pProxy as relay, or external server)
//  ??? MMO Direct connect to external server
//	Game session communications Relay, send Tables
//  Host? Broadcast?
//
// //////////////////////////////////////////////////////////////////////////////////////////////////////////

// /game/reguser/<username>

func RegUser(w http.ResponseWriter, r *http.Request) {

	//Parse URL Values
	values := strings.Split(r.URL.Path, "/")

	//Get username, string, only value needed
	encodedData := values[len(values)-1]
	newUsernameBytes, err := base32.StdEncoding.DecodeString(strings.ToUpper(encodedData))
	if err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}

	newUsername := string(newUsernameBytes)

	//Check if username exists
	userExists := false
	for _, user := range usersCache {
		if user == newUsername {
			userExists = true
			break
		}
	}

	//Return error to client if username already exists
	if userExists {
		fmt.Fprintf(w, "Registering failed!: "+newUsername+" already registered")
		return
	}

	//Make new client in cache
	id := uuid.New()
	usersCache[id] = newUsername

	//Let Picotron know the id of their new user
	fmt.Fprintf(w, id.String())
	println(r.RemoteAddr + ": new user: " + newUsername)
}
