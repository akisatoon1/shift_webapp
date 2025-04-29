package router

import (
	"backend/context"
	"net/http"
)

// ルーティング情報を表す構造体
type route struct {
	method    string
	pattern   string
	handlerFn http.HandlerFunc
}

// ミドルウェアを適用してルーティングを設定するヘルパー関数
func applyRoutes(mux *http.ServeMux, routes []route) {
	for _, r := range routes {
		handler := validateContentType(r.handlerFn)
		mux.HandleFunc(r.method+" "+r.pattern, handler)
	}
}

func Routes(mux *http.ServeMux, ctx *context.AppContext) {
	hdlr := &handler{ctx: ctx}

	routes := []route{
		{"GET", "/requests", hdlr.getRequestsRequest},
		{"GET", "/requests/{id}/entries", hdlr.getEntriesRequest},
		{"POST", "/requests", hdlr.postRequestsRequest},
	}

	applyRoutes(mux, routes)
}
