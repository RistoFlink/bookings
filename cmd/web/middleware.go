package main

import (
	"net/http"

	"github.com/justinas/nosurf"
)

// all middleware must take and return a handler!

// NoSurf adds CSRF protection to all POST requests
func NoSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)

	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   app.InProduction,
		SameSite: http.SameSiteLaxMode,
	})

	return csrfHandler
}

// SessionLoad loads and saves the session data for the current request
func SessionLoad(next http.Handler) http.Handler {
	return session.LoadAndSave(next)
}
