package main

import (
	"encoding/json"
	"net/http"
)

func handlerValidateChirp(writer http.ResponseWriter, req *http.Request) {
	type chirpRequest struct {
		Body string `json:"body"`
	}
	type chirpValidResponse struct {
		Valid bool `json:"valid"`
	}

	writer.Header().Set("Content-Type", "application/json")

	decoder := json.NewDecoder(req.Body)
	var chirpReq chirpRequest
	if err := decoder.Decode(&chirpReq); err != nil {
		respondWithError(writer, http.StatusInternalServerError, err.Error())
		return
	}

	var chirpResp chirpValidResponse
	chirpResp.Valid = len(chirpReq.Body) <= 400
	if !chirpResp.Valid {
		respondWithError(writer, http.StatusBadRequest, "Chirp is too long")
		return
	}

	respondWithJSON(writer, http.StatusOK, chirpResp)
}