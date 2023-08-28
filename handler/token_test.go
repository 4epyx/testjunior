package handler_test

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/4epyx/testtask/handler"
	"github.com/4epyx/testtask/repository/mongorepo"
	"github.com/4epyx/testtask/service"
	"github.com/4epyx/testtask/util"
	"github.com/4epyx/testtask/util/database"
	"github.com/stretchr/testify/suite"
)

type TokenTest struct {
	suite.Suite
	ctx     context.Context
	handler *handler.TokenHandler

	refToken string
}

func (t *TokenTest) SetupTest() {
	t.ctx = context.Background()

	envVars, errs := util.ParseEnv("TEST_DB_URI", "ACCESS_TOKEN_TTL", "REFRESH_TOKEN_TTL", "ACCESS_TOKEN_SECRET", "REFRESH_TOKEN_SECRET")
	if len(errs) != 0 {
		t.T().Fatal(errs)
	}

	accessTokenTTL, err := time.ParseDuration(envVars["ACCESS_TOKEN_TTL"])
	if err != nil {
		t.T().Fatal(err)
	}
	refreshTokenTTL, err := time.ParseDuration(envVars["REFRESH_TOKEN_TTL"])
	if err != nil {
		t.T().Fatal(err)
	}

	conn, err := database.SetupMongoConnection(t.ctx, envVars["TEST_DB_URI"])

	collection := conn.Database("testjunior").Collection("refresh_tokens")
	repo := mongorepo.NewMongoRTRepo(collection)

	tokenService := service.NewTokenService(repo, accessTokenTTL, refreshTokenTTL, []byte(envVars["ACCESS_TOKEN_SECRET"]), []byte(envVars["REFRESH_TOKEN_SECRET"]))

	t.refToken, err = tokenService.GenerateRefreshToken(t.ctx, "f022c263-4796-4bdc-92ac-4c74020183b3")
	t.T().Log(t.refToken)
	if err != nil {
		t.T().Fatal(err)
	}

	t.handler = handler.NewTokenHandler(tokenService)
}

func (t *TokenTest) TestGetAccessAndRefreshTokens() {
	req, err := http.NewRequest(http.MethodGet, "/tokens", bytes.NewReader([]byte(`{"user_guid": "f022c263-4796-4bdc-92ac-4c74020183b3"}`)))
	if err != nil {
		t.Fail(err.Error())
	}

	rr := httptest.NewRecorder()
	h := http.HandlerFunc(t.handler.GetAccessAndRefreshTokens)
	h.ServeHTTP(rr, req)

	if !t.Equal(http.StatusOK, rr.Code) {
		t.T().Log(rr.Body)
	}
}

func (t *TokenTest) TestRefreshToken() {
	req, err := http.NewRequest(http.MethodPost, "/refresh", strings.NewReader(fmt.Sprintf(`{"refresh_token":"%s"}`, t.refToken)))
	if err != nil {
		t.Fail(err.Error())
	}

	rr := httptest.NewRecorder()
	h := http.HandlerFunc(t.handler.RefreshToken)
	h.ServeHTTP(rr, req)

	if !t.Equal(http.StatusOK, rr.Code) {
		t.T().Log(rr.Body)
	}
}

func TestTokenSuite(t *testing.T) {
	test := new(TokenTest)
	suite.Run(t, test)
}
