package handlers

import (
    "encoding/json"
    "net/http"
    "financial-service/internal/services"
    "github.com/rs/zerolog/log"
)

type UserHandler struct {
    userService *services.UserService
}

func NewUserHandler(userService *services.UserService) *UserHandler {
    return &UserHandler{
        userService: userService,
    }
}

type RegisterUserRequest struct {
    Username string `json:"username"`
    Email    string `json:"email"`
    Password string `json:"password"`
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
    var req RegisterUserRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    user, err := h.userService.RegisterUser(r.Context(), req.Username, req.Email, req.Password)

    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(user)
}

type LoginUserRequest struct {
    Email    string `json:"email"`
    Password string `json:"password"`
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
    var req LoginUserRequest

    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    log.Printf("Login request: %+v", req)

    token, err := h.userService.LoginUser(r.Context(), req.Email, req.Password)

    log.Printf("Token: %+v", token)
    
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(token)
}