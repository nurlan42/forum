package internal

import (
	"net/http"

	uuid "github.com/satori/go.uuid"
)

// SessCheker is ...
func SetCookie(w http.ResponseWriter) uuid.UUID {
	sID := uuid.NewV4()

	cookie := &http.Cookie{
		Name:  "session",
		Value: sID.String(),
	}
	http.SetCookie(w, cookie)
	return sID
}
