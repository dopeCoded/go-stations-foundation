package router

import (
	"database/sql"
	"net/http"
	"github.com/TechBowl-japan/go-stations/handler"
	"github.com/TechBowl-japan/go-stations/service"
)

// NewRouter sets up the HTTP router with all necessary endpoints.
func NewRouter(db *sql.DB) http.Handler {
	mux := http.NewServeMux()

	// Register HealthzHandler
	healthzHandler := handler.NewHealthzHandler()
	mux.Handle("/healthz", healthzHandler)

	// Create TODOService
	todoService := service.NewTODOService(db)

	// Create TODOHandler and register
	todoHandler := handler.NewTODOHandler(todoService)
	mux.Handle("/todos", todoHandler)

	// 他のエンドポイントの登録もここで行う

	return mux
}
