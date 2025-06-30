package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

func handlerValidateChirp(writer http.ResponseWriter, req *http.Request) {
	type chirpRequest struct {
		Body string `json:"body"`
	}
	type chirpValidResponse struct {
		CleanedBody string `json:"cleaned_body"`
	}

	decoder := json.NewDecoder(req.Body)
	var chirpReq chirpRequest
	if err := decoder.Decode(&chirpReq); err != nil {
		respondWithError(writer, http.StatusInternalServerError, err.Error())
		return
	}

	var chirpResp chirpValidResponse
	lengthValid := len(chirpReq.Body) <= 400
	if !lengthValid {
		respondWithError(writer, http.StatusBadRequest, "Chirp is too long")
		return
	}
	chirpResp.CleanedBody = getCleanedBody(chirpReq.Body)
	respondWithJSON(writer, http.StatusOK, chirpResp)
}

func getCleanedBody(text string) string {
	profanityList := []string{"kerfuffle", "sharbert", "fornax"}
	words := strings.Split(text, " ")
	for i, word := range(words) {
		for _, profanity := range profanityList {
			if strings.ToLower(word) == profanity {
				words[i] = "****"
			}
		}
	}
	return strings.Join(words, " ")
}