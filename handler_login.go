package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/philipreese/chirpy-go/internal/auth"
)

func (cfg *apiConfig) handlerLogin(writer http.ResponseWriter, req *http.Request) {
	type createLoginRequest struct {
		Email            string `json:"email"`
		Password         string `json:"password"`
		ExpiresInSeconds *int   `json:"expires_in_seconds"`
	}

	decoder := json.NewDecoder(req.Body)
	var loginRequest createLoginRequest
	if err := decoder.Decode(&loginRequest); err != nil {
		respondWithError(writer, http.StatusInternalServerError, "Couldn't decode parameters: " + err.Error())
		return
	}

	expiresIn := 3600
	if loginRequest.ExpiresInSeconds != nil {
		if *loginRequest.ExpiresInSeconds > 0 && *loginRequest.ExpiresInSeconds < 3600 {
			expiresIn = *loginRequest.ExpiresInSeconds
		}
	}

	user, err := cfg.db.GetUserByEmail(req.Context(), loginRequest.Email)
	if err != nil {
		respondWithError(writer, http.StatusUnauthorized, "Incorrect email or password")
		return
	}

	if err := auth.CheckPasswordHash(loginRequest.Password, user.HashedPassword); err != nil {
		respondWithError(writer, http.StatusUnauthorized, "Incorrect email or password")
		return
	}

	tokenString, err := auth.MakeJWT(user.ID, cfg.tokenSecret, time.Duration(expiresIn) * time.Second)
	if err != nil {		
		respondWithError(writer, http.StatusUnauthorized, "Failed to create token: " + err.Error())
	}

	respondWithJSON(writer, http.StatusOK, User{
		ID: user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email: user.Email,
		Token: tokenString,
	})
}