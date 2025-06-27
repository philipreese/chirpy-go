package main

import "net/http"

func (cfg *apiConfig) handlerReset(writer http.ResponseWriter, req *http.Request) {
	cfg.fileserverHits.Store(0)
	writer.WriteHeader(http.StatusOK)
	writer.Write([]byte("Hits reset to 0"))
}