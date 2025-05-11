package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/Avnish-oP/opinionz/config"
	"github.com/Avnish-oP/opinionz/models"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
)

func ViewPost(w http.ResponseWriter, r *http.Request) {
	// Get post ID from URL parameters
	vars := mux.Vars(r)
	postID := vars["id"]
	if postID == "" {
		http.Error(w, "Post ID is required", http.StatusBadRequest)
		return
	}

	// Fetch the post
	postCollection := config.MongoDB.Collection("posts")
	var post models.Post
	err := postCollection.FindOne(r.Context(), bson.M{"_id": postID}).Decode(&post)
	if err != nil {
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	}

	// Fetch comments for the post
	commentCollection := config.MongoDB.Collection("comments")
	var comments []models.Comment
	cursor, err := commentCollection.Find(r.Context(), bson.M{"post_id": postID})
	if err != nil {
		http.Error(w, "Failed to fetch comments", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(r.Context())

	for cursor.Next(r.Context()) {
		var comment models.Comment
		if err := cursor.Decode(&comment); err == nil {
			comments = append(comments, comment)
		}
	}

	if err := cursor.Err(); err != nil {
		http.Error(w, "Error iterating comments", http.StatusInternalServerError)
		return
	}

	// Respond with post and comments
	response := models.Response{
		Message: "Post fetched successfully",
		Success: true,
		Data: map[string]interface{}{
			"post":     post,
			"comments": comments,
		},
	}
	json.NewEncoder(w).Encode(response)
}
