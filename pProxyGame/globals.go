package pProxyGame

import "github.com/google/uuid"

// So we can use util/multiMsg
// Holds all in progress requests
var multiRequestsCache = make(map[uuid.UUID]map[int]string)
var usersCache = make(map[uuid.UUID]string)
