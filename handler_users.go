package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/philipreese/chirpy-go/internal/auth"
	"github.com/philipreese/chirpy-go/internal/database"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
}

func (cfg *apiConfig) handlerCreateUser(writer http.ResponseWriter, req *http.Request) {
	type createUserRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(req.Body)
	var userRequest createUserRequest
	if err := decoder.Decode(&userRequest); err != nil {
		respondWithError(writer, http.StatusInternalServerError, "Couldn't decode parameters: " + err.Error())
		return
	}

	hashedPassword, err := auth.HashPassword(userRequest.Password)
	if err != nil {
		respondWithError(writer, http.StatusInternalServerError, "Couldn't hash password: " + err.Error())
		return
	}

	dbUser, err := cfg.db.CreateUser(req.Context(), database.CreateUserParams{
		Email: userRequest.Email,
		HashedPassword: hashedPassword,
	})
	if err != nil {
		respondWithError(writer, http.StatusInternalServerError, "Couldn't create user: " + err.Error())
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