package router

import (
	"database/sql"
	"net/http"
	"github.com/TechBowl-japan/go-stations/handler"
	"github.com/TechBowl-japan/go-stations/service"
	"github.com/TechBowl-japan/go-stations/handler/middleware" //station1
)

// NewRouter sets up the HTTP router with all necessary endpoints.
//station4
func NewRouter(db *sql.DB, userID, password string) http.Handler {
	mux := http.NewServeMux()

	// Register HealthzHandler
	healthzHandler := handler.NewHealthzHandler()
	mux.Handle("/healthz", healthzHandler)

	// Create TODOService
	todoService := service.NewTODOService(db)

	// Create TODOHandler and register
	todoHandler := handler.NewTODOHandler(todoService)

	//station3, 4 
	// ミドルウェアを適用
	// Order: Recovery -> OSExtractor -> BasicAuth -> LoggingMiddleware -> Handler
	// Basic 認証ミドルウェア
	basicAuthMiddleware := middleware.NewBasicAuthMiddleware(userID, password)
	wrappedTODOHandler := middleware.Recovery(
		middleware.OSExtractor(
			basicAuthMiddleware.Handler(
				middleware.LoggingMiddleware(todoHandler),
			),
		),
	)

	mux.Handle("/todos", wrappedTODOHandler)
	//station3 end

	// 他のエンドポイントの登録もここで行う
	//station1
	// PanicHandler をミドルウェアでラップして登録
    panicHandler := &handler.PanicHandler{}
    mux.Handle("/do-panic", middleware.Recovery(panicHandler))
	//station1 end

	return mux
}
