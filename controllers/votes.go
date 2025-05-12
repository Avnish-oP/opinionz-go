package controllers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Avnish-oP/opinionz/config"
	"github.com/Avnish-oP/opinionz/middlewares"
	"github.com/Avnish-oP/opinionz/models"
	"go.mongodb.org/mongo-driver/bson"
)

func HandleVote(w http.ResponseWriter, r *http.Request) {
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

	if voteRequest.PostID == "" || (voteRequest.Vote != "upvote" && voteRequest.Vote != "downvote") {
		http.Error(w, "Invalid vote input", http.StatusBadRequest)
		return
	}

	log.Printf("[HandleVote] userID=%s | postID=%s | vote=%s", userID, voteRequest.PostID, voteRequest.Vote)

	collection := config.MongoDB.Collection("posts")
	filter := bson.M{"_id": voteRequest.PostID}

	// Step 0: Fetch the post to check if upvotes_user_ids and downvotes_user_ids are arrays
	var post bson.M
	err := collection.FindOne(r.Context(), filter).Decode(&post)
	if err != nil {
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	}

	// Step 1: Ensure vote fields are arrays (initialize if not)
	updateFields := bson.M{}
	if _, ok := post["upvotes_user_ids"].(bson.A); !ok {
		updateFields["upvotes_user_ids"] = bson.A{} // Initialize as empty array if not present
	}
	if _, ok := post["downvotes_user_ids"].(bson.A); !ok {
		updateFields["downvotes_user_ids"] = bson.A{} // Initialize as empty array if not present
	}

	// Step 2: Update the post document if needed (to initialize the vote fields)
	if len(updateFields) > 0 {
		_, err = collection.UpdateOne(r.Context(), filter, bson.M{"$set": updateFields})
		if err != nil {
			log.Printf("[HandleVote] Failed to normalize vote fields: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	}

	// Step 3: Pull the user from both upvotes_user_ids and downvotes_user_ids (if the user has voted already)
	_, err = collection.UpdateOne(r.Context(), filter, bson.M{
		"$pull": bson.M{
			"upvotes_user_ids":   userID,
			"downvotes_user_ids": userID,
		},
	})
	if err != nil {
		log.Printf("[HandleVote] Pull failed: %v", err)
		http.Error(w, "Failed to remove previous vote", http.StatusInternalServerError)
		return
	}

	// Step 4: Add the user to the appropriate array (either upvotes or downvotes)
	targetField := "upvotes_user_ids"
	if voteRequest.Vote == "downvote" {
		targetField = "downvotes_user_ids"
	}

	_, err = collection.UpdateOne(r.Context(), filter, bson.M{
		"$addToSet": bson.M{
			targetField: userID,
		},
	})
	if err != nil {
		log.Printf("[HandleVote] Add to set failed: %v", err)
		http.Error(w, "Failed to apply vote", http.StatusInternalServerError)
		return
	}

	// Step 5: Return success message
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(models.Response{
		Message: "Vote updated successfully",
		Success: true,
	})
}
