package types

import (
	"encoding/json"
	"net/http"
)

type ErrorResponse struct {
	Error   string            `json:"error"`
	Details map[string]string `json:"details,omitempty"`
}

type SuccessResponse struct {
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

func WriteJSON(w http.ResponseWriter, status int, data interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
}

func WriteError(w http.ResponseWriter, status int, message string) {
	WriteJSON(w, status, ErrorResponse{Error: message})
}

func WriteValidationErrors(w http.ResponseWriter, errors map[string]string) {
	WriteJSON(w, http.StatusBadRequest, ErrorResponse{
		Error:   "Validation failed",
		Details: errors,
	})
}

func WriteSuccess(w http.ResponseWriter, message string, data interface{}) {
	WriteJSON(w, http.StatusOK, SuccessResponse{
		Message: message,
		Data:    data,
	})
}

type LoginResponse struct {
	Token string      `json:"token"`
}

type PaginatedResponse struct {
	Message    string      `json:"message,omitempty"`
	Data       interface{} `json:"data"`
	Total      int         `json:"total"`
	Page       int         `json:"page"`
	Limit      int         `json:"limit"`
	TotalPages int         `json:"total_pages"`
}

func WritePaginatedSuccess(w http.ResponseWriter, message string, data interface{}, total, page, limit int) {
	totalPages := (total + limit - 1) / limit
	if totalPages == 0 {
		totalPages = 1
	}

	WriteJSON(w, http.StatusOK, PaginatedResponse{
		Message:    message,
		Data:       data,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	})
}