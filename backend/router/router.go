package router

import (
	"backend/context"
	"backend/handler"
	"backend/middleware"
	"net/http"
	"path/filepath"
)

// ルーティング情報を表す構造体
type route struct {
	method    string
	pattern   string
	handlerFn handler.HandlerFuncWithContext
}

// ミドルウェアを適用してルーティングを設定するヘルパー関数
func applyRoutes(ctx *context.AppContext, mux *http.ServeMux, routes []route) {
	basePath := "/api"
	for _, r := range routes {
		if r.method == "POST" {
			r.handlerFn = middleware.ValidateContentType(r.handlerFn)
		}
		handler := handler.NewHandler(ctx, r.handlerFn)
		path := filepath.Join(basePath, r.pattern)
		mux.Handle(r.method+" "+path, handler)
	}
}

func Routes(mux *http.ServeMux, ctx *context.AppContext) {
	routes := []route{
		{"POST", "/login", handler.LoginRequest},
		{"GET", "/session", handler.GetSessionRequest},
		{"DELETE", "/session", handler.LogoutRequest},
		{"GET", "/requests", handler.GetRequestsRequest},
		{"GET", "/requests/{id}", handler.GetRequestRequest},
		{"POST", "/requests", handler.PostRequestsRequest},
		{"POST", "/requests/{id}/submissions", handler.PostSubmissionsRequest},
	}

	applyRoutes(ctx, mux, routes)
}
