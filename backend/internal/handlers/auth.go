package handlers

import (
	"backend/internal/auth"
	"backend/internal/models"
	"database/sql"
	"encoding/json"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

var DB *sql.DB

func SetDB(d *sql.DB) {
	DB = d
}

func SignUpHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Login string `json:"login"`
		Pass  string `json:"pass"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request payload", http.StatusBadRequest)
		return
	}
	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Pass), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "failed hashing password", http.StatusInternalServerError)
		return
	}
	user := models.User{Login: req.Login, PassHash: string(hashed)}
	if err := user.Create(DB); err != nil {
		http.Error(w, fmt.Sprintf("failed creating user: %v", err), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func LogInHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Login string `json:"login"`
		Pass  string `json:"pass"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request payload", http.StatusBadRequest)
		return
	}
	user := models.User{Login: req.Login}
	if err := user.FindByLogin(DB, req.Login); err != nil {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PassHash), []byte(req.Pass)); err != nil {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}
	tokenStr, err := auth.GenerateToken(user.ID)
	if err != nil {
		http.Error(w, fmt.Sprintf("token generation failed: %v", err), http.StatusInternalServerError)
		return
	}
	resp := map[string]string{"token": tokenStr}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
