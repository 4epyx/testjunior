package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/4epyx/testtask/handler"
	"github.com/4epyx/testtask/repository/mongorepo"
	"github.com/4epyx/testtask/service"
	"github.com/4epyx/testtask/util"
	"github.com/4epyx/testtask/util/database"
)

func main() {
	ctx := context.Background()

	// get required environment variables
	envs, errs := util.ParseEnv("SERVER_PORT", "DB_URI", "DB_NAME", "ACCESS_TOKEN_TTL", "REFRESH_TOKEN_TTL", "ACCESS_TOKEN_SECRET", "REFRESH_TOKEN_SECRET")
	if len(errs) != 0 {
		panic(fmt.Sprintf("%v", errs))
	}

	// parse access and refresh token TTL from string to time.Duration
	accessTokenTTL, err := time.ParseDuration(envs["ACCESS_TOKEN_TTL"])
	if err != nil {
		panic(err)
	}
	refreshTokenTTL, err := time.ParseDuration(envs["REFRESH_TOKEN_TTL"])
	if err != nil {
		panic(err)
	}

	// setup mongodb connection
	db, err := database.SetupMongoConnection(ctx, envs["DB_URI"])
	if err != nil {
		panic(err)
	}
	// disconnect from database
	defer func() {
		dbCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		if err := db.Disconnect(dbCtx); err != nil {
			panic(err)
		}
	}()

	// setup repository, service and handler
	repo := mongorepo.NewMongoRTRepo(db.Database(envs["DB_NAME"]).Collection("refresh_tokens"))
	s := service.NewTokenService(repo, accessTokenTTL, refreshTokenTTL, []byte(envs["ACCESS_TOKEN_SECRET"]), []byte(envs["ERFRESH_TOKEN_SECRET"]))
	h := handler.NewTokenHandler(s)

	mux := setupRoutes(h)

	if err := http.ListenAndServe(fmt.Sprintf(":%s", envs["SERVER_PORT"]), mux); err != nil {
		panic(err)
	}
}

func setupRoutes(h *handler.TokenHandler) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/token/generate", h.GetAccessAndRefreshTokens)
	mux.HandleFunc("/token/refresh", h.RefreshToken)

	return mux
}
