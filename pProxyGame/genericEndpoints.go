package pProxyGame

import (
	"encoding/base32"
	"net/http"
	"pProxy/util"
	"strings"
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

// /game/reguser/<username>/<password>

func RegUserHandler(w http.ResponseWriter, r *http.Request) {

	//Parse URL Values
	values := strings.Split(r.URL.Path, "/")

	//Get username
	encodedUsername := values[len(values)-2]
	newUsernameBytes, err := base32.StdEncoding.DecodeString(strings.ToUpper(encodedUsername))
	if err != nil {
		util.ReturnMessage(w, r, "", err)
		return
	}

	//Get password
	encodedPassword := values[len(values)-1]
	newPasswordBytes, err := base32.StdEncoding.DecodeString(strings.ToUpper(encodedPassword))
	if err != nil {
		util.ReturnMessage(w, r, "", err)
		return
	}

	//Change to string from bytes
	newUsername := string(newUsernameBytes)
	newPassword := string(newPasswordBytes)

	//Attempt to regiser user
	err = RegisterUser(newUsername, newPassword)
	if err != nil {
		util.ReturnMessage(w, r, "", err)
		return
	}

	//Notify client of successful registration
	util.ReturnMessage(w, r, "User: "+newUsername+" successfully registered", nil)
	return
}
