package adaptor

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/newnok6/kkp-dime-golang-meetup-2025/backend/domain"
	"github.com/newnok6/kkp-dime-golang-meetup-2025/backend/port"
)

type HTTPHandler struct {
	service port.StockOrderService
}

func NewHTTPHandler(service port.StockOrderService) *HTTPHandler {
	return &HTTPHandler{
		service: service,
	}
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type SuccessResponse struct {
	Message string `json:"message"`
}

func (h *HTTPHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/api/orders", h.CreateOrder).Methods("POST")
	router.HandleFunc("/api/orders", h.ListOrders).Methods("GET")
	router.HandleFunc("/api/orders/{id}", h.GetOrder).Methods("GET")
	router.HandleFunc("/api/orders/{id}/cancel", h.CancelOrder).Methods("POST")
	router.HandleFunc("/health", h.HealthCheck).Methods("GET")
}

func (h *HTTPHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var req domain.CreateOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	order, err := h.service.CreateOrder(r.Context(), req)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondJSON(w, http.StatusCreated, order)
}

func (h *HTTPHandler) GetOrder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	orderID := vars["id"]

	order, err := h.service.GetOrder(r.Context(), orderID)
	if err != nil {
		respondError(w, http.StatusNotFound, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, order)
}

func (h *HTTPHandler) ListOrders(w http.ResponseWriter, r *http.Request) {
	orders, err := h.service.ListOrders(r.Context())
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, orders)
}

func (h *HTTPHandler) CancelOrder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	orderID := vars["id"]

	if err := h.service.CancelOrder(r.Context(), orderID); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, SuccessResponse{Message: "order cancelled successfully"})
}

func (h *HTTPHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, http.StatusOK, map[string]string{"status": "healthy"})
}

func respondJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}

func respondError(w http.ResponseWriter, status int, message string) {
	respondJSON(w, status, ErrorResponse{Error: message})
}
