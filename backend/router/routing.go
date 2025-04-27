package router

import (
	"backend/context"
	"backend/handler"
	"net/http"
)

func Routing(ctx *context.AppContext) {
	http.HandleFunc("GET /requests", handler.GetRequestsRequest)
	http.HandleFunc("GET /requests/{id}/entries", handler.GetEntriesRequest)
}
