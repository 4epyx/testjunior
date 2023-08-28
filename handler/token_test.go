package handler_test

import (
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
	ctx      context.Context
	handler  *handler.TokenHandler
	userGuid string

	refToken           string
	expiredToken       string
	fakeSignatureToken string
}

func (t *TokenTest) SetupTest() {
	t.userGuid = "f022c263-4796-4bdc-92ac-4c74020183b3"
	t.ctx = context.Background()

	envVars, errs := util.ParseEnv("TEST_DB_URI")
	if len(errs) != 0 {
		t.T().Fatal(errs)
	}

	conn, err := database.SetupMongoConnection(t.ctx, envVars["TEST_DB_URI"])
	if err != nil {
		t.T().Fatal(err)
	}
	collection := conn.Database("testjunior").Collection("refresh_tokens")
	repo := mongorepo.NewMongoRTRepo(collection)

	tokenService := service.NewTokenService(repo, time.Minute, time.Hour, []byte("access-token-secret"), []byte("refresh-token-secret"))

	t.refToken, err = tokenService.GenerateRefreshToken(t.ctx, t.userGuid)
	if err != nil {
		t.T().Fatal(err)
	}
	t.expiredToken, err = generateExpiredToken(t.ctx, repo, t.userGuid)
	if err != nil {
		t.T().Fatal(err)
	}
	t.fakeSignatureToken, err = generateFakeSignatureToken(t.ctx, repo, t.userGuid)
	if err != nil {
		t.T().Fatal(err)
	}

	t.handler = handler.NewTokenHandler(tokenService)
}

func (t *TokenTest) TestGetAccessAndRefreshTokens() {
	tests := []struct {
		name       string
		userGuid   string
		wantStatus int
	}{
		{
			name:       "default test",
			userGuid:   t.userGuid,
			wantStatus: 200,
		},
		{
			name:       "empty user GUID",
			userGuid:   "",
			wantStatus: 400,
		},
		{
			name:       "invalid GUID",
			userGuid:   "qwerqwer",
			wantStatus: 400,
		},
	}

	for _, tt := range tests {
		req, err := http.NewRequest(http.MethodPost, "/token/generate", strings.NewReader(fmt.Sprintf(`{"user_guid":"%s"}`, tt.userGuid)))
		if err != nil {
			t.Fail(err.Error())
		}

		rr := httptest.NewRecorder()
		h := http.HandlerFunc(t.handler.GetAccessAndRefreshTokens)
		h.ServeHTTP(rr, req)

		if !t.Equal(tt.wantStatus, rr.Code) {
			t.T().Log(tt.name, rr.Body)
		}
	}
}

func (t *TokenTest) TestRefreshToken() {
	tests := []struct {
		name       string
		token      string
		wantStatus int
	}{
		{
			name:       "default test",
			token:      t.refToken,
			wantStatus: 200,
		},
		{
			name:       "expired token",
			token:      t.expiredToken,
			wantStatus: 400,
		},
		{
			name:       "fake signature token",
			token:      t.fakeSignatureToken,
			wantStatus: 400,
		},
		{
			name:       "empty token",
			token:      "",
			wantStatus: 400,
		},
	}

	for _, tt := range tests {
		req, err := http.NewRequest(http.MethodPost, "/token/refresh", strings.NewReader(fmt.Sprintf(`{"refresh_token":"%s"}`, tt.token)))
		if err != nil {
			t.Fail(err.Error())
		}

		rr := httptest.NewRecorder()
		h := http.HandlerFunc(t.handler.RefreshToken)
		h.ServeHTTP(rr, req)

		if !t.Equal(tt.wantStatus, rr.Code) {
			t.T().Log(tt.name, rr.Body)
		}
	}

}

func TestTokenSuite(t *testing.T) {
	test := new(TokenTest)
	suite.Run(t, test)
}

func generateExpiredToken(ctx context.Context, repo mongorepo.MongoRefreshTokenRepository, userGuid string) (string, error) {
	ts := service.NewTokenService(repo, time.Second, time.Second*-1, []byte("access-token-secret"), []byte("refresh-token-secret"))
	return ts.GenerateRefreshToken(ctx, userGuid)
}

func generateFakeSignatureToken(ctx context.Context, repo mongorepo.MongoRefreshTokenRepository, userGuid string) (string, error) {
	ts := service.NewTokenService(repo, time.Second, time.Hour, []byte("access-token-secret"), []byte("fake-refresh-token-secret"))
	return ts.GenerateRefreshToken(ctx, userGuid)
}
