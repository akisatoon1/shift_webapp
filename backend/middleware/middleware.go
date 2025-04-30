package middleware

import (
	"backend/context"
	"backend/handler"
	"net/http"
)

func ValidateContentType(next handler.HandlerFuncWithContext) handler.HandlerFuncWithContext {
	fn := func(ctx *context.AppContext, w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			contentType := r.Header.Get("Content-Type")
			if contentType != "application/json" {
				http.Error(w, "Content-Type must be application/json", http.StatusUnsupportedMediaType)
				return
			}
		}
		next(ctx, w, r)
	}
	return fn
}
