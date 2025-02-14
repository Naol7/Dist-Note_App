package main

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

func TestCreateNote(t *testing.T) {
	note := Note{ID: "test-id", Content: "This is a test note"}
	reqBody, _ := json.Marshal(note)

	req := httptest.NewRequest(http.MethodPost, "/notes", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	createNote(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", w.Code)
	}

	var createdNote Note
	json.NewDecoder(w.Body).Decode(&createdNote)

	if createdNote.ID != note.ID || createdNote.Content != note.Content {
		t.Errorf("Created note mismatch. Expected %+v, got %+v", note, createdNote)
	}
}

func TestGetNote(t *testing.T) {
	note := Note{ID: "test-id", Content: "This is a test note"}
	collection := client.Database("notes_db").Collection("notes")
	_, _ = collection.InsertOne(context.TODO(), note)

	req := httptest.NewRequest(http.MethodGet, "/notes/get?id=test-id", nil)
	w := httptest.NewRecorder()

	getNote(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var fetchedNote Note
	json.NewDecoder(w.Body).Decode(&fetchedNote)

	if fetchedNote.ID != note.ID || fetchedNote.Content != note.Content {
		t.Errorf("Fetched note mismatch. Expected %+v, got %+v", note, fetchedNote)
	}

	collection.DeleteOne(context.TODO(), bson.M{"id": "test-id"})
}

func TestUpdateNote(t *testing.T) {
	note := Note{ID: "test-id", Content: "Old content"}
	collection := client.Database("notes_db").Collection("notes")
	_, _ = collection.InsertOne(context.TODO(), note)

	updatedNote := Note{ID: "test-id", Content: "Updated content"}
	reqBody, _ := json.Marshal(updatedNote)
	req := httptest.NewRequest(http.MethodPut, "/notes/update", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	updateNote(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var responseNote Note
	json.NewDecoder(w.Body).Decode(&responseNote)

	if responseNote.Content != "Updated content" {
		t.Errorf("Expected content to be 'Updated content', got '%s'", responseNote.Content)
	}

	collection.DeleteOne(context.TODO(), bson.M{"id": "test-id"})
}

func TestDeleteNote(t *testing.T) {
	note := Note{ID: "test-id", Content: "Delete me"}
	collection := client.Database("notes_db").Collection("notes")
	_, _ = collection.InsertOne(context.TODO(), note)

	req := httptest.NewRequest(http.MethodDelete, "/notes/delete?id=test-id", nil)
	w := httptest.NewRecorder()

	deleteNote(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("Expected status 204, got %d", w.Code)
	}

	time.Sleep(100 * time.Millisecond)

	var deletedNote Note
	err := collection.FindOne(context.TODO(), bson.M{"_id": "test-id"}).Decode(&deletedNote)
	if err == nil {
		t.Errorf("Note was not deleted")
	} else {
		t.Logf("Note successfully deleted")
	}
}

func TestListNotes(t *testing.T) {
	collection := client.Database("notes_db").Collection("notes")
	collection.InsertMany(context.TODO(), []interface{}{
		Note{ID: "note-1", Content: "First note"},
		Note{ID: "note-2", Content: "Second note"},
	})

	req := httptest.NewRequest(http.MethodGet, "/notes/list", nil)
	w := httptest.NewRecorder()

	listNotes(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var notes []Note
	json.NewDecoder(w.Body).Decode(&notes)

	if len(notes) < 2 {
		t.Errorf("Expected at least 2 notes, got %d", len(notes))
	}

	collection.DeleteMany(context.TODO(), bson.M{})
}
