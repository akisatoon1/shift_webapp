package router

import (
	"backend/context"
	"net/http"
)

// ルーティング情報を表す構造体
type route struct {
	method    string
	pattern   string
	handlerFn handlerFunc
}

// ミドルウェアを適用してルーティングを設定するヘルパー関数
func applyRoutes(ctx *context.AppContext, mux *http.ServeMux, routes []route) {
	for _, r := range routes {
		r.handlerFn = validateContentType(r.handlerFn)
		handler := &handler{ctx: ctx, handlerFn: r.handlerFn}
		mux.Handle(r.method+" "+r.pattern, handler)
	}
}

func Routes(mux *http.ServeMux, ctx *context.AppContext) {
	routes := []route{
		{"GET", "/requests", getRequestsRequest},
		{"GET", "/requests/{id}/entries", getEntriesRequest},
		{"POST", "/requests", postRequestsRequest},
	}

	applyRoutes(ctx, mux, routes)
}
