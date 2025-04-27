package router

import (
	"backend/context"
	"net/http"
)

func Routes(ctx *context.AppContext) {
	hdlr := &handler{ctx: ctx}

	http.HandleFunc("GET /requests", hdlr.getRequestsRequest)
	http.HandleFunc("GET /requests/{id}/entries", hdlr.getEntriesRequest)
}
