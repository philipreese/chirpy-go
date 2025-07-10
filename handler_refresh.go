package main

import (
	"net/http"
	"time"

	"github.com/philipreese/chirpy-go/internal/auth"
)

func (cfg *apiConfig) handlerRefresh(writer http.ResponseWriter, req *http.Request) {
	type refreshResponse struct {
		Token string `json:"token"`
	}

	tokenString, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(writer, http.StatusBadRequest, "Couldn't get bearer token: " + err.Error())
		return
	}

	refreshToken, err := cfg.db.GetRefreshToken(req.Context(), tokenString)
	if err != nil {
		respondWithError(writer, http.StatusUnauthorized, "Couldn't get refresh token: " + err.Error())
		return
	}

	if refreshToken.ExpiresAt.Before(time.Now()) {
		respondWithError(writer, http.StatusUnauthorized, "Refresh token expired")
		return
	}

	if refreshToken.RevokedAt.Valid {
		respondWithError(writer, http.StatusUnauthorized, "Refresh token revoked")
		return
	}

	token, err := auth.MakeJWT(refreshToken.UserID, cfg.tokenSecret, time.Hour)
	if err != nil {
		respondWithError(writer, http.StatusUnauthorized, "Failed to create token: " + err.Error())
		return
	}

	respondWithJSON(writer, http.StatusOK, refreshResponse{Token: token})
}

func (cfg *apiConfig) handlerRevoke(writer http.ResponseWriter, req *http.Request) {
	tokenString, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(writer, http.StatusBadRequest, "Couldn't get bearer token: " + err.Error())
		return
	}

	refreshToken, err := cfg.db.GetRefreshToken(req.Context(), tokenString)
	if err != nil {
		respondWithError(writer, http.StatusInternalServerError, "Couldn't get refresh token: " + err.Error())
		return
	}

	if err := cfg.db.RevokeRefreshToken(req.Context(), refreshToken.Token); err != nil {
		respondWithError(writer, http.StatusInternalServerError, "Couldn't revoke refresh token: " + err.Error())
		return
	}

	writer.WriteHeader(http.StatusNoContent)
}