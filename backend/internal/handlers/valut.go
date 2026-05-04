package handlers

import (
	"backend/internal/crypto"
	"backend/internal/midleware"
	"backend/internal/models"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
)

func CreateEntryHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Sign string `json:"sign"`
		Date string `json:"date"`
		Note string `json:"note"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request payload", http.StatusBadRequest)
		return
	}
	parsedDate, err := time.Parse(time.RFC3339, req.Date)
	if err != nil {
		http.Error(w, "invalid timestamp format, use RFC3339", http.StatusBadRequest)
		return
	}
	rawKey := os.Getenv("ENCRYPTION_KEY")
	if rawKey == "" {
		http.Error(w, "encryption key not set", http.StatusInternalServerError)
		return
	}
	keySum := sha256.Sum256([]byte(rawKey))
	key := keySum[:] // 32-byte AES key
	ciphertext, err := crypto.Encrypt([]byte(req.Note), key)
	if err != nil {
		http.Error(w, fmt.Sprintf("encryption failed: %v", err), http.StatusInternalServerError)
		return
	}
	// get userID from context
	userVal := r.Context().Value(midleware.CtxUserIDKey)
	userID, ok := userVal.(int)
	if !ok {
		http.Error(w, "user not authenticated", http.StatusUnauthorized)
		return
	}
	entry := models.VaultEntry{
		UserID:     userID,
		Sign:       req.Sign,
		Date:       parsedDate,
		Ciphertext: ciphertext,
	}
	if err := entry.Create(DB); err != nil {
		http.Error(w, fmt.Sprintf("failed saving entry: %v", err), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func GetEntryHandler(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 4 {
		http.Error(w, "invalid path", http.StatusBadRequest)
		return
	}
	sign := parts[2]
	dateStr := parts[3]
	// Parse date
	parsedDate, err := time.Parse(time.RFC3339, dateStr)
	if err != nil {
		http.Error(w, "invalid timestamp format", http.StatusBadRequest)
		return
	}
	entry, err := (&models.VaultEntry{}).GetEntry(DB, sign, parsedDate)
	if err != nil {
		http.Error(w, fmt.Sprintf("entry not found: %v", err), http.StatusNotFound)
		return
	}
	rawKey := os.Getenv("ENCRYPTION_KEY")
	if rawKey == "" {
		http.Error(w, "encryption key not set", http.StatusInternalServerError)
		return
	}
	keySum := sha256.Sum256([]byte(rawKey))
	key := keySum[:]
	plaintext, err := crypto.Decrypt(entry.Ciphertext, key)
	if err != nil {
		http.Error(w, fmt.Sprintf("decryption failed: %v", err), http.StatusInternalServerError)
		return
	}
	resp := map[string]string{"sign": sign, "date": dateStr, "note": string(plaintext)}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
