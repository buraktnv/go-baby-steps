package handlers

import (
    "encoding/json"
    "net/http"
    "strconv"
    "financial-service/internal/services"
    "github.com/gorilla/mux"
)

type BalanceHandler struct {
    balanceService *services.BalanceService
}

func NewBalanceHandler(balanceService *services.BalanceService) *BalanceHandler {
    return &BalanceHandler{
        balanceService: balanceService,
    }
}

func (h *BalanceHandler) GetBalance(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    userID, err := strconv.ParseUint(vars["user_id"], 10, 32)
    if err != nil {
        http.Error(w, "Invalid user ID", http.StatusBadRequest)
        return
    }

    balance, err := h.balanceService.GetBalance(r.Context(), uint(userID))
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(balance)
} 