package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

type TestUser struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func TestRegisterLogin(t *testing.T) {
	// Create a test user
	testUser := TestUser{
		Username: "testuser",
		Password: "password123",
	}

	jsonData, _ := json.Marshal(testUser)

	req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	handler := http.HandlerFunc(register)
	handler.ServeHTTP(resp, req)

	log.Printf("Register response: %s", resp.Body.String())
	
	if resp.Code != http.StatusCreated && resp.Code != http.StatusConflict {
		t.Fatalf("Expected status %d or %d but got %d", http.StatusCreated, http.StatusConflict, resp.Code)
	}

	
	req, _ = http.NewRequest("POST", "/login", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	resp = httptest.NewRecorder()
	handler = http.HandlerFunc(login)
	handler.ServeHTTP(resp, req)

	log.Printf("Login response: %s", resp.Body.String())

	if resp.Code != http.StatusOK {
		t.Fatalf("Expected status %d but got %d", http.StatusOK, resp.Code)
	}

	var responseData map[string]string
	json.Unmarshal(resp.Body.Bytes(), &responseData)
	token, ok := responseData["token"]
	if !ok {
		t.Fatal("Token not found in response")
	}

	req, _ = http.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp = httptest.NewRecorder()
	handler = ValidateToken(getUserProfile)
	handler.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("Expected status %d but got %d", http.StatusOK, resp.Code)
	}
}

func TestUpdateUser(t *testing.T) {
	testUser := TestUser{Username: "testuser", Password: "password123"}
	jsonData, _ := json.Marshal(testUser)
	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	handler := http.HandlerFunc(login)
	handler.ServeHTTP(resp, req)

	log.Printf("Login response: %s", resp.Body.String())

	var responseData map[string]string
	json.Unmarshal(resp.Body.Bytes(), &responseData)
	token := responseData["token"]

	updateData := map[string]string{"password": "newpassword123"}
	updateJson, _ := json.Marshal(updateData)
	req, _ = http.NewRequest("PUT", "/update", bytes.NewBuffer(updateJson))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	resp = httptest.NewRecorder()
	handler = ValidateToken(updateUser)
	handler.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("Expected status %d but got %d", http.StatusOK, resp.Code)
	}
}

func TestDeleteUser(t *testing.T) {
	testUser := TestUser{Username: "testuser", Password: "newpassword123"}
	jsonData, _ := json.Marshal(testUser)
	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	handler := http.HandlerFunc(login)
	handler.ServeHTTP(resp, req)

	var responseData map[string]string
	json.Unmarshal(resp.Body.Bytes(), &responseData)
	token := responseData["token"]

	req, _ = http.NewRequest("DELETE", "/delete", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp = httptest.NewRecorder()
	handler = ValidateToken(deleteUser)
	handler.ServeHTTP(resp, req)

	if resp.Code != http.StatusNoContent {
		t.Fatalf("Expected status %d but got %d", http.StatusNoContent, resp.Code)
	}
}
