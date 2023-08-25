package handler_test

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/4epyx/testtask/handler"
	"github.com/4epyx/testtask/repository/mongorepo"
	"github.com/4epyx/testtask/service"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type envData struct {
	DbUrl string

	accessTokenTTL     time.Duration
	refreshTokenTTL    time.Duration
	accessTokenSecret  []byte
	refreshTokenSecret []byte
}

type TokenTest struct {
	suite.Suite
	ctx     context.Context
	handler *handler.TokenHandler

	refToken string
}

func (t *TokenTest) SetupTest() {
	t.ctx = context.Background()

	envVars, err := getTestEnv()
	if err != nil {
		t.T().Fatal(err)
	}

	opts := options.Client().ApplyURI(envVars.DbUrl)
	conn, err := mongo.Connect(t.ctx, opts)
	if err != nil {
		t.T().Fatal(err)
	}

	collection := conn.Database("testjunior").Collection("refresh_tokens")
	repo := mongorepo.NewMongoRTRepo(collection)

	tokenService := service.NewTokenService(repo, envVars.accessTokenTTL, envVars.refreshTokenTTL, envVars.accessTokenSecret, envVars.refreshTokenSecret)

	t.refToken, err = tokenService.GenerateRefreshToken(t.ctx, "f022c263-4796-4bdc-92ac-4c74020183b3")
	t.T().Log(t.refToken)
	if err != nil {
		t.T().Fatal(err)
	}

	t.handler = handler.NewTokenHandler(tokenService)
}

func getTestEnv() (envData, error) {
	data := envData{}
	ok := true

	data.DbUrl, ok = os.LookupEnv("TEST_DB_URL")
	if !ok {
		return data, errors.New("couldn't connect to database")
	}

	strAccessTokenTTL, ok := os.LookupEnv("ACCESS_TOKEN_TTL")
	if !ok {
		strAccessTokenTTL = "1h"
	}
	strRefreshTokenTTL, ok := os.LookupEnv("REFRESH_TOKEN_TTL")
	if !ok {
		strRefreshTokenTTL = "720h"
	}

	var err error
	data.accessTokenTTL, err = time.ParseDuration(strAccessTokenTTL)
	if err != nil {
		return data, err
	}

	data.refreshTokenTTL, err = time.ParseDuration(strRefreshTokenTTL)
	if err != nil {
		return data, err
	}

	strAccessTokenSecret, ok := os.LookupEnv("ACCESS_TOKEN_SECRET")
	if !ok {
		return data, err
	}

	strRefreshTokenSecret, ok := os.LookupEnv("REFRESH_TOKEN_SECRET")
	if !ok {
		return data, err
	}

	data.accessTokenSecret = []byte(strAccessTokenSecret)
	data.refreshTokenSecret = []byte(strRefreshTokenSecret)

	return data, nil
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
