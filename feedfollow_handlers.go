package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/nikitahahua/RssAgg/internal/database"
	"github.com/nikitahahua/RssAgg/models"
	"net/http"
	"time"
)

func (cfg *apiConfig) handlerFeedFollowsGet(w http.ResponseWriter, r *http.Request, user database.User) {
	feedFollows, err := cfg.DB.GetFeedFollowsForUser(r.Context(), user.ID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create feed follow")
		return
	}

	jsonResponse(w, http.StatusOK, models.GetAllFollowsFromDb(feedFollows))
}

func (apiCfg *apiConfig) handlerCreateFeedFollows(w http.ResponseWriter, r *http.Request, user database.User) {
	type params struct {
		FeedID uuid.UUID `json:"feed_id"`
	}

	parameters := params{}
	err := json.NewDecoder(r.Body).Decode(&parameters)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error parsing json :("))
		return
	}

	feed, err := apiCfg.DB.CreateFeedFollow(r.Context(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    user.ID,
		FeedID:    parameters.FeedID,
	})

	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error creating user :( "))
		return
	}

	jsonResponse(w, 201, models.GetFeedFollowsFromDb(feed))
}

func (cfg *apiConfig) handlerDeleteFeedFollow(w http.ResponseWriter, r *http.Request, user database.User) {
	feedFollowIDStr := chi.URLParam(r, "feedFollowID")
	feedFollowID, err := uuid.Parse(feedFollowIDStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid feed follow ID")
		return
	}

	err = cfg.DB.DeleteFeedFollow(r.Context(), database.DeleteFeedFollowParams{
		UserID: user.ID,
		ID:     feedFollowID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create feed follow")
		return
	}

	jsonResponse(w, http.StatusOK, struct{}{})
}
