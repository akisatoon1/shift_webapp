package router

import (
	"backend/handler"
	"net/http"
)

func Routing() {
	http.HandleFunc("GET /requests", handler.GetRequestsRequest)
	http.HandleFunc("GET /requests/{id}/entries", handler.GetEntriesRequest)
}
