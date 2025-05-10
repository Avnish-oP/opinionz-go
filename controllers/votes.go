package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Avnish-oP/opinionz/config"
	"github.com/Avnish-oP/opinionz/middlewares"
	"github.com/Avnish-oP/opinionz/models"
	"go.mongodb.org/mongo-driver/bson"
)

func HandleVote(w http.ResponseWriter, r *http.Request) {
	// Extract userID from the context
	userID, ok := r.Context().Value(middlewares.UserIDKey).(string)
	if !ok || userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var voteRequest struct {
		PostID string `json:"post_id"`
		Vote   string `json:"vote"`
	}
	if err := json.NewDecoder(r.Body).Decode(&voteRequest); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if voteRequest.Vote != "upvote" && voteRequest.Vote != "downvote" {
		http.Error(w, "Invalid vote type", http.StatusBadRequest)
		return
	}

	fmt.Println("Received PostID:", voteRequest.PostID)

	update := bson.M{}
	if voteRequest.Vote == "upvote" {
		update = bson.M{
			"$addToSet": bson.M{"upvotes_user_ids": userID},
			"$pull":     bson.M{"downvotes_user_ids": userID},
		}
	} else if voteRequest.Vote == "downvote" {
		update = bson.M{
			"$addToSet": bson.M{"downvotes_user_ids": userID},
			"$pull":     bson.M{"upvotes_user_ids": userID},
		}
	}

	collection := config.MongoDB.Collection("posts")
	filter := bson.M{"_id": voteRequest.PostID} // Use PostID as a string
	_, err := collection.UpdateOne(r.Context(), filter, update)
	if err != nil {
		http.Error(w, "Failed to update vote", http.StatusInternalServerError)
		return
	}

	response := models.Response{
		Message: "Vote updated successfully",
		Success: true,
	}
	json.NewEncoder(w).Encode(response)
}
