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
func (h *TODOHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.createTODO(w, r)
	case http.MethodPut:
		h.updateTODO(w, r)
	default:
		w.Header().Set("Allow", "POST, PUT")
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

// createTODO handles POST requests to create a new TODO.
func (h *TODOHandler) createTODO(w http.ResponseWriter, r *http.Request) {
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
		TODO: *todo, // model.CreateTODOResponse は TODO を値として期待
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

// updateTODO handles PUT requests to update an existing TODO.
func (h *TODOHandler) updateTODO(w http.ResponseWriter, r *http.Request) {
	// Decode the JSON request body into UpdateTODORequest
	var req model.UpdateTODORequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields() // Reject unknown fields
	if err := decoder.Decode(&req); err != nil {
		http.Error(w, "Bad Request: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Validate that ID is not zero and Subject is not empty
	if req.ID == 0 {
		http.Error(w, "Bad Request: id is required and must be greater than 0", http.StatusBadRequest)
		return
	}
	if req.Subject == "" {
		http.Error(w, "Bad Request: subject is required", http.StatusBadRequest)
		return
	}

	// Call the service layer to update the TODO
	updatedTODO, err := h.service.UpdateTODO(r.Context(), req.ID, req.Subject, req.Description)
	if err != nil {
		// Check if the error is ErrNotFound
		if model.IsErrNotFound(err) {
			http.Error(w, "Not Found: "+err.Error(), http.StatusNotFound)
			return
		}
		// Handle other potential errors (e.g., database constraints)
		log.Println("Error updating TODO:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Prepare the UpdateTODOResponse with the updated TODO
	resp := model.UpdateTODOResponse{
		TODO: updatedTODO, // model.UpdateTODOResponse は *TODO を期待
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
