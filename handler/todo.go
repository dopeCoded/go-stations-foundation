package handler

import (
	
	"encoding/json"
	"log"
	"net/http"

	"github.com/TechBowl-japan/go-stations/model"
	"github.com/TechBowl-japan/go-stations/service"
)

// TODOHandler handles HTTP requests for TODO operations.
type TODOHandler struct {
	service *service.TODOService
}

// NewTODOHandler creates a new TODOHandler with the provided TODOService.
func NewTODOHandler(svc *service.TODOService) *TODOHandler {
	return &TODOHandler{
		service: svc,
	}
}

// ServeHTTP implements the http.Handler interface for TODOHandler.
// It handles POST requests to create a new TODO.
func (h *TODOHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Ensure the request method is POST
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Decode the JSON request body into CreateTODORequest
	var req model.CreateTODORequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields() // Reject unknown fields
	if err := decoder.Decode(&req); err != nil {
		http.Error(w, "Bad Request: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Validate that the Subject field is not empty
	if req.Subject == "" {
		http.Error(w, "Bad Request: subject is required", http.StatusBadRequest)
		return
	}

	// Call the service layer to create the TODO
	todo, err := h.service.CreateTODO(r.Context(), req.Subject, req.Description)
	if err != nil {
		log.Println("Error creating TODO:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Prepare the CreateTODOResponse with the created TODO
	resp := model.CreateTODOResponse{
		TODO: *todo,
	}

	// Set the Content-Type header to application/json
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Encode the response as JSON and send it
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Println("Error encoding JSON response:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
