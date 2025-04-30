package handler

import (
	"backend/context"
	"net/http"
)

// TODO
// ログインしているユーザーのIDを取得する

// コンテキストに依存するhandler関数
type HandlerFuncWithContext func(*context.AppContext, http.ResponseWriter, *http.Request)

type Handler struct {
	ctx       *context.AppContext
	handlerFn HandlerFuncWithContext
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.handlerFn(h.ctx, w, r)
}

func NewHandler(ctx *context.AppContext, handlerFn HandlerFuncWithContext) *Handler {
	return &Handler{ctx: ctx, handlerFn: handlerFn}
}
