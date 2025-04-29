package router

import (
	"backend/context"
	"net/http"
)

func Routes(mux *http.ServeMux, ctx *context.AppContext) {
	hdlr := &handler{ctx: ctx}

	mux.HandleFunc("GET /requests", validateContentType(hdlr.getRequestsRequest))
	mux.HandleFunc("GET /requests/{id}/entries", validateContentType(hdlr.getEntriesRequest))
	mux.HandleFunc("POST /requests", validateContentType(hdlr.postRequestsRequest))
}
