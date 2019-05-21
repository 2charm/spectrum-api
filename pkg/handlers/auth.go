package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/2charm/spectrum-api/pkg/sessions"
	"github.com/2charm/spectrum-api/pkg/users"
)

//SessionState represents a session that is started by an authenticated user
type SessionState struct {
	StartTime time.Time   `json:"startTime,omitempty"`
	User      *users.User `json:"user,omitempty"`
}

func (ctx *HandlerContext) UsersHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		contentType := r.Header.Get("Content-Type")
		if !strings.HasPrefix(contentType, "application/json") {
			http.Error(w, "The request body must be of type JSON.", http.StatusUnsupportedMediaType)
			return
		}

		newUser := &users.NewUser{}
		dec := json.NewDecoder(r.Body)
		if err := dec.Decode(newUser); err != nil {
			http.Error(w, "Error decoding json into User.", http.StatusBadRequest)
			return
		}

		user, err := newUser.ToUser()
		if err != nil {
			http.Error(w, "Error creating User from NewUser.", http.StatusBadRequest)
			return
		}

		user, err = ctx.UserStore.Insert(user)
		if err != nil {
			http.Error(w, "Error inserting user into User store.", http.StatusInternalServerError)
			return
		}

		session := &SessionState{
			StartTime: time.Now(),
			User:      user,
		}
		_, err = sessions.BeginSession(ctx.SigningKey, ctx.SessionStore, session, w)
		if err != nil {
			http.Error(w, "Error creating session in server.", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		if user.ID == 0 {
			http.Error(w, "Error with adding user to database", http.StatusInternalServerError)
			return
		}
		buffer, err := json.Marshal(user)
		if err != nil {
			log.Printf("%v", err)
			http.Error(w, "Error marshaling JSON.", http.StatusInternalServerError)
			return
		}
		w.Write(buffer)
	} else {
		http.Error(w, "Incompatible http method.", http.StatusMethodNotAllowed)
		return
	}
}

//SessionsHandler handles requests for sessions
func (ctx *HandlerContext) SessionsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "incompatible http method", http.StatusMethodNotAllowed)
		return
	}

	contentType := r.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "application/json") {
		http.Error(w, "request body must be of type JSON", http.StatusUnsupportedMediaType)
		return
	}

	creds := &users.Credentials{}
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(creds); err != nil {
		http.Error(w, fmt.Sprintf("error decoding JSON: %v", err), http.StatusBadRequest)
		return
	}
	user, err := ctx.UserStore.GetByEmail(creds.Email)
	if err != nil {
		time.Sleep(time.Second * 7)
		http.Error(w, "invalid credentials, email not found", http.StatusUnauthorized)
		return
	}

	err = user.Authenticate(creds.Password)
	if err != nil {
		http.Error(w, "invalid credentials, password is not correct", http.StatusUnauthorized)
		return
	}
	sessState := SessionState{
		StartTime: time.Now(),
		User:      user,
	}
	_, err = sessions.BeginSession(ctx.SigningKey, ctx.SessionStore, sessState, w)
	if err != nil {
		http.Error(w, "error starting new session", http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	buffer, err := json.Marshal(user)
	if err != nil {
		http.Error(w, "error marshaling JSON.", http.StatusInternalServerError)
		return
	}
	w.Write(buffer)
}

const sessionResourcePath = "/v1/sessions/"

func (ctx *HandlerContext) SpecificSessionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "DELETE" {
		http.Error(w, "incompatible http method", http.StatusMethodNotAllowed)
		return
	}
	seg := strings.TrimPrefix(r.RequestURI, sessionResourcePath)
	if seg != "mine" {
		http.Error(w, "user session invalid", http.StatusForbidden)
		return
	}

	_, err := sessions.EndSession(r, ctx.SigningKey, ctx.SessionStore)
	if err != nil {
		http.Error(w, "error ending session", http.StatusInternalServerError)
		return
	}
	w.Write([]byte("signed out"))
}
