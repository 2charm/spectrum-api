package handlers

import (
	"github.com/2charm/spectrum-api/pkg/sessions"
	"github.com/2charm/spectrum-api/pkg/users"
)

//HandlerContext represents a receiver for handlers to utilize and gain context
type HandlerContext struct {
	SigningKey   string
	SessionStore sessions.Store `json:"sessionStore,omitempty"`
	UserStore    users.Store    `json:"userStore,omitempty"`
}
