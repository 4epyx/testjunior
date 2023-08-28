package router

import (
	"net/http"

	"github.com/4epyx/testtask/handler"
	"github.com/4epyx/testtask/middleware"
	"github.com/rs/zerolog"
)

func SetupRoutes(h *handler.TokenHandler, logger *zerolog.Logger) *http.ServeMux {
	logMiddleware := middleware.NewLogMiddleware(logger)
	mux := http.NewServeMux()
	mux.Handle("/token/generate", logMiddleware.Log(http.HandlerFunc(h.GetAccessAndRefreshTokens)))
	mux.Handle("/token/refresh", logMiddleware.Log(http.HandlerFunc(h.RefreshToken)))

	return mux
}
