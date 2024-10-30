package model

import (
    "time"
)

type (
	// A TODO expresses ...
	TODO struct {
		ID          int       `json:"id"`
		Subject     string    `json:"subject"`
		Description string    `json:"description"`
		CreatedAt   time.Time `json:"created_at"`
		UpdatedAt   time.Time `json:"updated_at"`
	}

	// A CreateTODORequest expresses ...
	CreateTODORequest struct {
		Subject     string `json:"subject"`
		Description string `json:"description,omitempty"`
	}

	// A CreateTODOResponse expresses ...
	CreateTODOResponse struct {
		TODO TODO `json:"todo"`
	}

	// A ReadTODORequest expresses ...
	ReadTODORequest struct {
		ID int `json:"id"`
	}
	
	// A ReadTODOResponse expresses ...
	ReadTODOResponse struct {
		TODO TODO `json:"todo"`
	}

	// A UpdateTODORequest expresses ...
	UpdateTODORequest struct {
		TODO TODO `json:"todo"`
	}

	// A UpdateTODOResponse expresses ...
	UpdateTODOResponse struct {
		TODO TODO `json:"todo"`
	}

	// A DeleteTODORequest expresses ...
	DeleteTODORequest struct {
		TODO TODO `json:"todo"`
	}

	// A DeleteTODOResponse expresses ...
	DeleteTODOResponse struct {
		TODO TODO `json:"todo"`
	}
)
