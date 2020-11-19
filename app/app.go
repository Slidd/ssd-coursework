package app

import (
	"encoding/gob"

	"github.com/gorilla/sessions"
)

var (
	// Store session cookie
	Store *sessions.FilesystemStore
)

// Init user session
func Init() error {
	Store = sessions.NewFilesystemStore("", []byte("something-very-secret")) // should this be random?
	gob.Register(map[string]interface{}{})
	return nil
}
