package main

import (
	"net/http"
)

func routing() {
	http.HandleFunc("GET /requests", handleGetRequestsRequest)
	http.HandleFunc("GET /requests/{id}/entries", handleGetEntriesRequest)
}
