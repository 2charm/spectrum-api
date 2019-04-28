package handlers

import (
	"github.com/2charm/spectrum-api/gateway/models/sessions"
	"github.com/2charm/spectrum-api/gateway/models/users"
)

//HandlerContext represents a receiver for handlers to utilize and gain context
type HandlerContext struct {
	APIKey       string
	SigningKey   string
	SessionStore sessions.Store `json:"sessionStore,omitempty"`
	UserStore    users.Store    `json:"userStore,omitempty"`
}
