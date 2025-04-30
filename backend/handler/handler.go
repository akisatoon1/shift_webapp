package handler

import (
	"backend/context"
	"net/http"
)

// TODO
// ログインしているユーザーのIDを取得する

type HandlerFuncWithContext func(*context.AppContext, http.ResponseWriter, *http.Request) *AppError

type Handler struct {
	ctx       *context.AppContext
	handlerFn HandlerFuncWithContext
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := h.handlerFn(h.ctx, w, r); err != nil {
		http.Error(w, err.message, err.code)
	}
}

func NewHandler(ctx *context.AppContext, handlerFn HandlerFuncWithContext) *Handler {
	return &Handler{ctx: ctx, handlerFn: handlerFn}
}
