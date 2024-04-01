package util

import (
	"crypto/rand"
	"encoding/base32"
	"strings"
)

// Used for Session tokens etc, should be long enough to be randomly unique, but short enough to keep urls short
func GenToken() string {
	// Req:
	// Must be short, but long enough to be unbruteforcible over web requests
	// Must be properly random (crypto)
	// Must be safe to send over picotron fetch (a-z0-9) no uppercase

	// 5 bytes = 8 base32 chars (Clean, no padding)
	// ~1bil combinations i think?

	randData := make([]byte, 5)
	rand.Read(randData)

	encoded := base32.StdEncoding.EncodeToString(randData)

	return strings.ToLower(encoded)

}
