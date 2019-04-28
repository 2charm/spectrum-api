package sessions

import (
	"errors"
	"net/http"
	"strings"
)

const headerAuthorization = "Authorization"
const paramAuthorization = "auth"
const schemeBearer = "Bearer "

//ErrNoSessionID is used when no session ID was found in the Authorization header
var ErrNoSessionID = errors.New("no session ID found in " + headerAuthorization + " header")

//ErrInvalidScheme is used when the authorization scheme is not supported
var ErrInvalidScheme = errors.New("authorization scheme not supported")

//BeginSession creates a new SessionID, saves the `sessionState` to the store, adds an
//Authorization header to the response with the SessionID, and returns the new SessionID
func BeginSession(signingKey string, store Store, sessionState interface{}, w http.ResponseWriter) (SessionID, error) {
	sid, err := NewSessionID(signingKey)
	if err != nil {
		return InvalidSessionID, err
	}

	err = store.Save(sid, sessionState)
	if err != nil {
		return InvalidSessionID, err
	}

	w.Header().Add(headerAuthorization, schemeBearer+sid.String())

	return sid, nil
}

//GetSessionID extracts and validates the SessionID from the request headers
func GetSessionID(r *http.Request, signingKey string) (SessionID, error) {
	//TODO: get the value of the Authorization header,
	//or the "auth" query string parameter if no Authorization header is present,
	//and validate it. If it's valid, return the SessionID. If not
	//return the validation error.
	sid := r.Header.Get(headerAuthorization)
	if sid == "" {
		sid = r.URL.Query().Get(paramAuthorization)
		if sid == "" {
			return InvalidSessionID, ErrNoSessionID
		}
	}
	if !strings.HasPrefix(sid, schemeBearer) {
		return InvalidSessionID, ErrInvalidScheme
	}
	sid = strings.TrimPrefix(sid, schemeBearer)
	// return ValidateID(sid, signingKey)
	sidValid, err := ValidateID(sid, signingKey)
	if err != nil {
		return InvalidSessionID, err
	}
	return sidValid, nil
}

//GetState extracts the SessionID from the request,
//gets the associated state from the provided store into
//the `sessionState` parameter, and returns the SessionID
func GetState(r *http.Request, signingKey string, store Store, sessionState interface{}) (SessionID, error) {
	//TODO: get the SessionID from the request, and get the data
	//associated with that SessionID from the store.
	sid, err := GetSessionID(r, signingKey)
	if err != nil {
		return InvalidSessionID, ErrNoSessionID
	}

	errGet := store.Get(sid, sessionState)
	if errGet != nil {
		return InvalidSessionID, ErrStateNotFound
	}
	return sid, nil
}

//EndSession extracts the SessionID from the request,
//and deletes the associated data in the provided store, returning
//the extracted SessionID.
func EndSession(r *http.Request, signingKey string, store Store) (SessionID, error) {
	//TODO: get the SessionID from the request, and delete the
	//data associated with it in the store.
	sid, err := GetSessionID(r, signingKey)
	if err != nil {
		return InvalidSessionID, ErrNoSessionID
	}
	id := strings.TrimPrefix(sid.String(), schemeBearer)
	err = store.Delete(SessionID(id))
	return SessionID(sid), err
}
