package controllers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Avnish-oP/opinionz/config"
	"github.com/Avnish-oP/opinionz/models"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
)

func ViewPost(w http.ResponseWriter, r *http.Request) {
    log.Println("[ViewPost] Incoming request")

    // Get post ID from URL parameters
    vars := mux.Vars(r)
    postID := vars["id"]
    log.Printf("[ViewPost] Post ID from URL: %s\n", postID)

    if postID == "" {
        http.Error(w, "Post ID is required", http.StatusBadRequest)
        log.Println("[ViewPost] Error: Post ID is missing")
        return
    }

    // Fetch the post from DB
    collection := config.MongoDB.Collection("posts")
    var post models.Post
    err := collection.FindOne(r.Context(), bson.M{"_id": postID}).Decode(&post)
    if err != nil {
        http.Error(w, "Post not found", http.StatusNotFound)
        log.Printf("[ViewPost] Error: Post not found with ID %s | err: %v\n", postID, err)
        return
    }
    log.Printf("[ViewPost] Post retrieved: %+v\n", post)

    // Ensure upvotes_user_ids is an array (initialize it if nil)
    if post.Upvotes == nil {
        post.Upvotes = []string{}
    }

    // Fetch comments
    commentCollection := config.MongoDB.Collection("comments")
    cursor, err := commentCollection.Find(r.Context(), bson.M{"post_id": postID})
    if err != nil {
        http.Error(w, "Failed to fetch comments", http.StatusInternalServerError)
        log.Printf("[ViewPost] Error: Failed to fetch comments for post %s | err: %v\n", postID, err)
        return
    }
    defer cursor.Close(r.Context())

    var comments []models.Comment
    for cursor.Next(r.Context()) {
        var comment models.Comment
        if err := cursor.Decode(&comment); err == nil {
            comments = append(comments, comment)
            log.Printf("[ViewPost] Comment decoded: %+v\n", comment)
        } else {
            log.Printf("[ViewPost] Error decoding comment: %v\n", err)
        }
    }

    if err := cursor.Err(); err != nil {
        http.Error(w, "Error reading comments", http.StatusInternalServerError)
        log.Printf("[ViewPost] Error iterating comments: %v\n", err)
        return
    }

    log.Println("[ViewPost] Successfully fetched post and comments")

    // Send response
    w.Header().Set("Content-Type", "application/json")
    response := models.Response{
        Message: "Post fetched successfully",
        Success: true,
        Data: map[string]interface{}{
            "post":     post,
            "comments": comments,
        },
    }
    if err := json.NewEncoder(w).Encode(response); err != nil {
        log.Printf("[ViewPost] Error encoding response: %v\n", err)
    } else {
        log.Println("[ViewPost] Response sent successfully")
    }
}
