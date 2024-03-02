package main

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	_ "github.com/google/uuid"
	"github.com/nikitahahua/RssAgg/internal/database"
	"github.com/nikitahahua/RssAgg/models"
	_ "github.com/nikitahahua/RssAgg/models"
	"net/http"
	"time"
)

func handleRequest(w http.ResponseWriter, r *http.Request) {
	jsonResponse(w, http.StatusOK, map[string]string{"status": "ok"})
}

func handlerErr(w http.ResponseWriter, r *http.Request) {
	respondWithError(w, http.StatusInternalServerError, "Internal Gay Server Error")
}

func (apiCfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	type params struct {
		Name string `json:"name"`
	}

	parameters := params{}
	err := json.NewDecoder(r.Body).Decode(&parameters)
	if parameters.Name == "" {
		respondWithError(w, 400, fmt.Sprintf("you should send name attribute. :( "))
		return
	}
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error parsing json :( \n %s", err))
		return
	}

	user, err := apiCfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      parameters.Name,
	})

	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error creating user :(  \n %s", err))
		return
	}
	jsonResponse(w, 201, models.GetUserFromDb(user))
}

func (apiCfg *apiConfig) handleGetUser(w http.ResponseWriter, r *http.Request, user database.User) {
	jsonResponse(w, http.StatusOK, models.GetUserFromDb(user))
}
