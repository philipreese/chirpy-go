package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

func (cfg *apiConfig) handlerCreateUser(writer http.ResponseWriter, req *http.Request) {
	type createUserRequest struct {
		Email string `json:"email"`
	}

	decoder := json.NewDecoder(req.Body)
	var userRequest createUserRequest
	if err := decoder.Decode(&userRequest); err != nil {
		respondWithError(writer, http.StatusInternalServerError, err.Error())
		return
	}

	dbUser, err := cfg.db.CreateUser(req.Context(), userRequest.Email)
	if err != nil {
		respondWithError(writer, http.StatusInternalServerError, err.Error())
		return
	}
	
	user := User{
		ID: dbUser.ID,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
		Email: dbUser.Email,
	}
	respondWithJSON(writer, http.StatusCreated, user)
}