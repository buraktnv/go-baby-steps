package handlers

import (
    "encoding/json"
    "net/http"
    "financial-service/internal/services"
)

type TransactionHandler struct {
    service *services.TransactionService
}

func NewTransactionHandler(service *services.TransactionService) *TransactionHandler {
    return &TransactionHandler{
        service: service,
    }
}

type TransactionRequest struct {
    UserID uint    `json:"user_id"`
    Amount float64 `json:"amount"`
}

func (h *TransactionHandler) Credit(w http.ResponseWriter, r *http.Request) {
    var req TransactionRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    tx, err := h.service.Credit(r.Context(), req.UserID, req.Amount)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(tx)
}

func (h *TransactionHandler) Debit(w http.ResponseWriter, r *http.Request) {
    var req TransactionRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    tx, err := h.service.Debit(r.Context(), req.UserID, req.Amount)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(tx)
}

type TransferRequest struct {
    FromUserID uint    `json:"from_user_id"`
    ToUserID   uint    `json:"to_user_id"`
    Amount     float64 `json:"amount"`
}

func (h *TransactionHandler) Transfer(w http.ResponseWriter, r *http.Request) {
    var req TransferRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    tx, err := h.service.Transfer(r.Context(), req.FromUserID, req.ToUserID, req.Amount)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(tx)
} 