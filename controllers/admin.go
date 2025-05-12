package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Avnish-oP/opinionz/config"
	"github.com/Avnish-oP/opinionz/models"
	"go.mongodb.org/mongo-driver/bson"
)

func GetAllUsers(w http.ResponseWriter, r *http.Request) {
	userCollection := config.MongoDB.Collection("users")
	var users []models.User
	cursor, err := userCollection.Find(r.Context(), bson.M{})
	if err != nil {
		http.Error(w, "Failed to fetch users", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(r.Context())

	for cursor.Next(r.Context()) {
		var user models.User
		if err := cursor.Decode(&user); err == nil {
			users = append(users, user)
		}
	}

	if err := cursor.Err(); err != nil {
		http.Error(w, "Error iterating users", http.StatusInternalServerError)
		return
	}

	response := models.Response{
		Message: "Users fetched successfully",
		Success: true,
		Data:    users,
	}
	json.NewEncoder(w).Encode(response)
}

func GetAllPosts(w http.ResponseWriter, r *http.Request) {
	postCollection := config.MongoDB.Collection("posts")
	var posts []models.Post
	cursor, err := postCollection.Find(r.Context(), bson.M{})
	if err != nil {
		http.Error(w, "Failed to fetch posts", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(r.Context())

	for cursor.Next(r.Context()) {
		var post models.Post
		if err := cursor.Decode(&post); err == nil {
			posts = append(posts, post)
		}
	}

	if err := cursor.Err(); err != nil {
		http.Error(w, "Error iterating posts", http.StatusInternalServerError)
		return
	}

	response := models.Response{
		Message: "Posts fetched successfully",
		Success: true,
		Data:    posts,
	}
	json.NewEncoder(w).Encode(response)
}

func GetAllComments(w http.ResponseWriter, r *http.Request) {
	commentCollection := config.MongoDB.Collection("comments")
	var comments []models.Comment
	cursor, err := commentCollection.Find(r.Context(), bson.M{})
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

	response := models.Response{
		Message: "Comments fetched successfully",
		Success: true,
		Data:    comments,
	}
	json.NewEncoder(w).Encode(response)
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("id")
	if userID == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}
	fmt.Println(userID)

	userCollection := config.MongoDB.Collection("users")
	_, err := userCollection.DeleteOne(r.Context(), bson.M{"_id": userID})
	if err != nil {
		http.Error(w, "Failed to delete user", http.StatusInternalServerError)
		return
	}

	response := models.Response{
		Message: "User deleted successfully",
		Success: true,
	}
	json.NewEncoder(w).Encode(response)
}

func DeletePost(w http.ResponseWriter, r *http.Request) {
	postID := r.URL.Query().Get("id")
	if postID == "" {
		http.Error(w, "Post ID is required", http.StatusBadRequest)
		return
	}

	postCollection := config.MongoDB.Collection("posts")
	_, err := postCollection.DeleteOne(r.Context(), bson.M{"_id": postID})
	if err != nil {
		http.Error(w, "Failed to delete post", http.StatusInternalServerError)
		return
	}

	response := models.Response{
		Message: "Post deleted successfully",
		Success: true,
	}
	json.NewEncoder(w).Encode(response)
}

func DeleteComment(w http.ResponseWriter, r *http.Request) {
	commentID := r.URL.Query().Get("id")
	if commentID == "" {
		http.Error(w, "Comment ID is required", http.StatusBadRequest)
		return
	}

	commentCollection := config.MongoDB.Collection("comments")
	_, err := commentCollection.DeleteOne(r.Context(), bson.M{"_id": commentID})
	if err != nil {
		http.Error(w, "Failed to delete comment", http.StatusInternalServerError)
		return
	}

	response := models.Response{
		Message: "Comment deleted successfully",
		Success: true,
	}
	json.NewEncoder(w).Encode(response)
}
