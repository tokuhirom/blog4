package admin

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"time"
)

const (
	sessionCookieName     = "admin_session"
	sessionIDLength       = 32
	defaultSessionTimeout = 24 * time.Hour
)

func generateSessionID() (string, error) {
	b := make([]byte, sessionIDLength)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

func getSessionID(r *http.Request) string {
	cookie, err := r.Cookie(sessionCookieName)
	if err != nil {
		return ""
	}
	return cookie.Value
}
