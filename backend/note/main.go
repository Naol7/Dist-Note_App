package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

type Note struct {
	ID      string `json:"id" bson:"id"`
	Content string `json:"content" bson:"content"`
}

func init() {
	var err error
	uri := "mongodb+srv://naolgezahegne:x9fRCT7mxvC2oQMU@cluster0.pblqu.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0"
	clientOptions := options.Client().ApplyURI(uri)

	client, err = mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		panic(err)
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		panic(err)
	}
	fmt.Println("Connected to MongoDB!")
}

func main() {
	http.HandleFunc("/notes", corsMiddleware(createNote))
	http.HandleFunc("/notes/get", corsMiddleware(getNote))
	http.HandleFunc("/notes/update", corsMiddleware(updateNote))
	http.HandleFunc("/notes/delete", corsMiddleware(deleteNote))
	http.HandleFunc("/notes/list", corsMiddleware(listNotes))

	fmt.Println("Server running on port 8081...")
	if err := http.ListenAndServe(":8081", nil); err != nil {
		panic(err)
	}
}

func corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next(w, r)
	}
}

func createNote(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var note Note
	if err := json.NewDecoder(r.Body).Decode(&note); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	collection := client.Database("notes_db").Collection("notes")
	_, err := collection.InsertOne(context.TODO(), note)
	if err != nil {
		http.Error(w, "Failed to save note", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(note)
}

func getNote(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Missing note ID", http.StatusBadRequest)
		return
	}

	collection := client.Database("notes_db").Collection("notes")
	var note Note
	err := collection.FindOne(context.TODO(), bson.M{"id": id}).Decode(&note)
	if err != nil {
		http.Error(w, "Note not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(note)
}

func updateNote(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var updatedNote Note
	if err := json.NewDecoder(r.Body).Decode(&updatedNote); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	collection := client.Database("notes_db").Collection("notes")
	filter := bson.M{"id": updatedNote.ID}
	update := bson.M{
		"$set": bson.M{"content": updatedNote.Content},
	}

	_, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		http.Error(w, "Failed to update note", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updatedNote)
}

func deleteNote(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Missing note ID", http.StatusBadRequest)
		return
	}

	collection := client.Database("notes_db").Collection("notes")
	_, err := collection.DeleteOne(context.TODO(), bson.M{"id": id})
	if err != nil {
		http.Error(w, "Failed to delete note", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func listNotes(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	collection := client.Database("notes_db").Collection("notes")
	cursor, err := collection.Find(context.TODO(), bson.M{})
	if err != nil {
		http.Error(w, "Failed to list notes", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(context.TODO())

	var notes []Note
	for cursor.Next(context.TODO()) {
		var note Note
		if err := cursor.Decode(&note); err != nil {
			http.Error(w, "Failed to decode note", http.StatusInternalServerError)
			return
		}
		notes = append(notes, note)
	}

	if err := cursor.Err(); err != nil {
		http.Error(w, "Cursor error", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(notes)
}
