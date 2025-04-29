package router

import (
	"net/http"
)

func validateContentType(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			contentType := r.Header.Get("Content-Type")
			if contentType != "application/json" {
				// ここのエラーをトップレベルでキャッチできない
				http.Error(w, "Content-Type must be application/json", http.StatusUnsupportedMediaType)
				return
			}
		}
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}
