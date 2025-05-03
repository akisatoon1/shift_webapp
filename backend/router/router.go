package router

import (
	"backend/context"
	"backend/handler"
	"backend/middleware"
	"net/http"
)

// ルーティング情報を表す構造体
type route struct {
	method    string
	pattern   string
	handlerFn handler.HandlerFuncWithContext
}

// ミドルウェアを適用してルーティングを設定するヘルパー関数
func applyRoutes(ctx *context.AppContext, mux *http.ServeMux, routes []route) {
	for _, r := range routes {
		r.handlerFn = middleware.ValidateContentType(r.handlerFn)
		handler := handler.NewHandler(ctx, r.handlerFn)
		mux.Handle(r.method+" "+r.pattern, handler)
	}
}

func Routes(mux *http.ServeMux, ctx *context.AppContext) {
	routes := []route{
		{"POST", "/login", handler.LoginRequest},
		{"DELETE", "/session", handler.LogoutRequest},
		{"GET", "/requests", handler.GetRequestsRequest},
		{"GET", "/requests/{id}/entries", handler.GetEntriesRequest},
		{"POST", "/requests", handler.PostRequestsRequest},
		{"POST", "/requests/{id}/entries", handler.PostEntriesRequest},
	}

	applyRoutes(ctx, mux, routes)
}
