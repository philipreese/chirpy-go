package main

import "net/http"

func (cfg *apiConfig) handlerReset(writer http.ResponseWriter, req *http.Request) {
	if cfg.platform != "dev" {
		respondWithError(writer, http.StatusForbidden, "Reset is only allowed in dev environment")
		return
	}
	err := cfg.db.Reset(req.Context()); if err != nil {
		respondWithError(writer, http.StatusInternalServerError, err.Error())
		return
	}

	cfg.fileserverHits.Store(0)
	writer.WriteHeader(http.StatusOK)
	writer.Write([]byte("Hits reset to 0 and database reset to initial state"))
}