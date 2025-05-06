package handler

import (
	"backend/context"
	"net/http"
)

type HandlerFuncWithContext func(*context.AppContext, http.ResponseWriter, *http.Request) *AppError

type Handler struct {
	ctx       *context.AppContext
	handlerFn HandlerFuncWithContext
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := h.handlerFn(h.ctx, w, r); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(err.code)
		w.Write([]byte(`{"error": "` + err.message + `"}`))
	}
}

func NewHandler(ctx *context.AppContext, handlerFn HandlerFuncWithContext) *Handler {
	return &Handler{ctx: ctx, handlerFn: handlerFn}
}
