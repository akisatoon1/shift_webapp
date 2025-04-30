package router

import (
	"backend/context"
	"net/http"
)

func validateContentType(next handlerFunc) handlerFunc {
	fn := func(ctx *context.AppContext, w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			contentType := r.Header.Get("Content-Type")
			if contentType != "application/json" {
				// ここのエラーをトップレベルでキャッチできない
				http.Error(w, "Content-Type must be application/json", http.StatusUnsupportedMediaType)
				return
			}
		}
		next(ctx, w, r)
	}
	return fn
}
