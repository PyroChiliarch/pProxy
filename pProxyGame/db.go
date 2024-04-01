package pProxyGame

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"errors"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/argon2"
)

const dbFile string = "pProxy.db"

var DB *sql.DB
var err error

// ID is same as username, but all caps
const schema string = `
CREATE TABLE IF NOT EXISTS user (
id TEXT NOT NULL PRIMARY KEY,
time_created INT NOT NULL,
username TEXT,
hash TEXT,
salt TEXT
);`

func InitDatabase() {
	//Make database file
	DB, err = sql.Open("sqlite3", dbFile)
	if err != nil {
		println("Error opening " + dbFile + "\n" + err.Error())
	}

	//Apply the schema
	_, err = DB.Exec(schema)

	//
	// row, err := c.db.Query("SELECT * FROM activities WHERE id=?", id)
	// row := c.db.QueryRow("SELECT id, time, description FROM activities WHERE id=?", id)
	/*
			activity := api.Activity{}
		 	var err error
		 	if err = row.Scan(&activity.ID, &activity.Time, &activity.Description); err == sql.ErrNoRows {
		  	log.Printf("Id not found")
		  	return api.Activity{}, ErrIDNotFound
	*/
}

// Generate 16 random bytes
func genSalt() []byte {
	salt := make([]byte, 32)
	rand.Read(salt)
	return salt
}

func calcHash(password string, salt []byte) []byte {
	//These values should stay consistent, changes will mean users can no longer log in as their passwords will not match
	key := argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, 32) //func IDKey(password, salt []byte, time, memory uint32, threads uint8, keyLen uint32) []byte
	return key
}

func RegisterUser(username string, password string) error {

	// Query for existing user
	var result string
	if err := DB.QueryRow("SELECT (username) FROM user WHERE id = ?", strings.ToUpper(username)).Scan(&result); err != sql.ErrNoRows {
		if err == nil {
			return errors.New("Error registering user " + result + ": \n user exists!")
		} else {
			return errors.New("Error registering user " + username + ": " + dbFile + "\n" + err.Error())
		}
	}

	//Query had error sql.ErrNoRows, user does not exist, can register them
	salt := genSalt()
	hash := calcHash(password, salt)

	_, err := DB.Exec("INSERT INTO user VALUES(?,?,?,?,?);",
		strings.ToUpper(username), // ID
		time.Now().UnixNano(),     // Date created, fmt.Printf("%v", timestamp) Gets human readable
		username,
		base64.StdEncoding.EncodeToString(hash), // hash
		base64.StdEncoding.EncodeToString(salt)) // salt

	if err != nil {
		return errors.New("Error registering user: " + dbFile + "\n" + err.Error())
	}
	println("DEBUG: " + username + ":" + password)
	println("User registered: " + username)
	return nil // Return no error
}
