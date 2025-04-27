package router

import (
	"backend/context"
	"encoding/json"
	"net/http"
	"strconv"
)

type handler struct {
	ctx *context.AppContext
}

func (h *handler) getRequestsRequest(w http.ResponseWriter, r *http.Request) {
	requests, err := h.ctx.DB.GetRequests()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(requests)
}

func (h *handler) getEntriesRequest(w http.ResponseWriter, r *http.Request) {
	requestId := r.PathValue("id")
	requestIdInt, err := strconv.Atoi(requestId)
	if err != nil {
		http.Error(w, "requestidが整数ではありません", http.StatusBadRequest)
		return
	}

	entries, err := h.ctx.DB.GetEntriesByRequestID(requestIdInt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(entries)
}
