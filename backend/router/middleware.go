package router

import (
	"net/http"
)

func validateContentType(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			contentType := r.Header.Get("Content-Type")
			if contentType != "application/json" {
				http.Error(w, "Content-Type must be application/json", http.StatusUnsupportedMediaType)
				return
			}
		}
		next(w, r)
	}
}
