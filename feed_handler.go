package main

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/nikitahahua/RssAgg/internal/database"
	"github.com/nikitahahua/RssAgg/models"
	"net/http"
	"time"
)

func (apiCfg *apiConfig) handlerCreateFeed(w http.ResponseWriter, r *http.Request, user database.User) {
	type params struct {
		Name string `json:"name"`
		Url  string `json:"url"`
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

	feed, err := apiCfg.DB.CreateFeed(r.Context(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      parameters.Name,
		Url:       parameters.Url,
		UserID:    user.ID,
	})

	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error creating user :(  \n %s", err))
		return
	}
	jsonResponse(w, 201, models.GetFeedFromDb(feed))
}

func (apiCfg *apiConfig) handlerGetFeeds(w http.ResponseWriter, r *http.Request) {

	feeds, err := apiCfg.DB.GetFeeds(r.Context())

	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("couldn't get a feed :( "))
		return
	}
	jsonResponse(w, 200, models.GetAllFeedsFromDb(feeds))
}
