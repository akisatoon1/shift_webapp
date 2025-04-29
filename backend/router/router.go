package router

import (
	"backend/context"
	"net/http"
)

// ルーティング情報を表す構造体
type route struct {
	method  string
	pattern string
	handler http.Handler
}

// ミドルウェアを適用してルーティングを設定するヘルパー関数
func applyRoutes(mux *http.ServeMux, routes []route) {
	for _, r := range routes {
		handler := validateContentType(r.handler)
		mux.Handle(r.method+" "+r.pattern, handler)
	}
}

func Routes(mux *http.ServeMux, ctx *context.AppContext) {
	routes := []route{
		{"GET", "/requests", &handler{ctx: ctx, handlerFn: getRequestsRequest}},
		{"GET", "/requests/{id}/entries", &handler{ctx: ctx, handlerFn: getEntriesRequest}},
		{"POST", "/requests", &handler{ctx: ctx, handlerFn: postRequestsRequest}},
	}

	applyRoutes(mux, routes)
}
