package sessions

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
)

const headerAuthorization = "Authorization"
const paramAuthorization = "auth"
const schemeBearer = "Bearer "

//Empty struct for deleting sessionState
type session struct{}

//ErrNoSessionID is used when no session ID was found in the Authorization header
var ErrNoSessionID = errors.New("no session ID found in " + headerAuthorization + " header")

//ErrInvalidScheme is used when the authorization scheme is not supported
var ErrInvalidScheme = errors.New("authorization scheme not supported")

//BeginSession creates a new SessionID, saves the `sessionState` to the store, adds an
//Authorization header to the response with the SessionID, and returns the new SessionID
func BeginSession(signingKey string, store Store, sessionState interface{}, w http.ResponseWriter) (SessionID, error) {
	sessionID, err := NewSessionID(signingKey)
	if err != nil {
		return InvalidSessionID, fmt.Errorf("error creating session id %v", err)
	}
	if err := store.Save(sessionID, sessionState); err != nil {
		return InvalidSessionID, fmt.Errorf("error saving session id %v", err)
	}
	w.Header().Add(headerAuthorization, schemeBearer+sessionID.String())
	return sessionID, nil
}

//GetSessionID extracts and validates the SessionID from the request headers
func GetSessionID(r *http.Request, signingKey string) (SessionID, error) {
	authHeader := strings.Split(r.Header.Get(headerAuthorization), " ")
	authParam := strings.Split(r.URL.Query().Get(paramAuthorization), " ")
	sessionID, err := extractSessionID(authHeader, r, signingKey)
	if err != nil {
		sessionID, err = extractSessionID(authParam, r, signingKey)
	}
	//Returns InvalidSessionID if err != nil
	return sessionID, err
}

//GetState extracts the SessionID from the request,
//gets the associated state from the provided store into
//the `sessionState` parameter, and returns the SessionID
func GetState(r *http.Request, signingKey string, store Store, sessionState interface{}) (SessionID, error) {
	sessionID, err := GetSessionID(r, signingKey)
	if err != nil { // Returns InvalidSessionID
		return sessionID, err
	}
	if err := store.Get(sessionID, sessionState); err != nil {
		return InvalidSessionID, err
	}
	return sessionID, nil
}

//EndSession extracts the SessionID from the request,
//and deletes the associated data in the provided store, returning
//the extracted SessionID.
func EndSession(r *http.Request, signingKey string, store Store) (SessionID, error) {
	sessionID, err := GetState(r, signingKey, store, &session{})
	if err != nil {
		return sessionID, err
	}
	if err := store.Delete(sessionID); err != nil {
		return InvalidSessionID, fmt.Errorf("error deleting session id %d", err)
	}
	return sessionID, nil
}

//extractSessionID extracts the sessionID from the request header
func extractSessionID(auth []string, r *http.Request, signingKey string) (SessionID, error) {
	switch {
	case auth[0] == "":
	case len(auth) == 1 && auth[0]+" " == schemeBearer:
	case len(auth) == 2 && auth[0]+" " == schemeBearer:
		SessionID, err := ValidateID(auth[1], signingKey)
		if err != nil { // returns InvalidSessionID
			return SessionID, err
		}
		return SessionID, nil
	}
	return InvalidSessionID, ErrInvalidID
}