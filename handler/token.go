package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/4epyx/testtask/service"
	"github.com/4epyx/testtask/util"
	"github.com/google/uuid"
)

type TokenHandler struct {
	service *service.TokenService
}

func NewTokenHandler(service *service.TokenService) *TokenHandler {
	return &TokenHandler{service: service}
}

// GetAccessAndRefreshTokens handles the HTTP-request and returns the access and refresh tokens, generated by the given user GUID
func (h *TokenHandler) GetAccessAndRefreshTokens(w http.ResponseWriter, r *http.Request) {
	type bodyData struct {
		UserGuid string `json:"user_guid"`
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, `{"status": "couldn't read body"}`, http.StatusBadRequest)
		return
	}

	data := &bodyData{}
	if err := json.Unmarshal(body, data); err != nil {
		http.Error(w, `{"status": "couldn't parse body"}`, http.StatusBadRequest)
		return
	}

	if _, err := uuid.Parse(data.UserGuid); err != nil {
		http.Error(w, fmt.Sprintf(`{"status": "invalid user GUID %s"}`, data.UserGuid), http.StatusBadRequest)
		return
	}

	tokens := util.Tokens{}

	tokens.AccessToken, err = h.service.GenerateAccessToken(r.Context(), data.UserGuid)
	if err != nil {
		http.Error(w, `{"status": "error occured while generating access token"}`, http.StatusInternalServerError)
		return
	}

	tokens.RefreshToken, err = h.service.GenerateRefreshToken(r.Context(), data.UserGuid)
	if err != nil {
		http.Error(w, `{"status": "error occured while generating refresh token"}`, http.StatusInternalServerError)
		return
	}

	jsonResp, err := json.Marshal(tokens)
	if err != nil {
		http.Error(w, `{"status": "error occured while marshaling JSON"}`, http.StatusInternalServerError)
		return
	}

	w.Write(jsonResp)
}

// RefreshToken handles the HTTP-request and returns new pair of tokens,
// generated by the refresh token. For work, refresh token must be contained
// in the request's body and must be generated by this server
func (h *TokenHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	type bodyData struct {
		RefreshToken string `json:"refresh_token"`
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, `{"status": "couldn't read body"}`, http.StatusBadRequest)
		return
	}

	data := &bodyData{}
	if err := json.Unmarshal(body, data); err != nil {
		http.Error(w, `{"status": "couldn't parse body to JSON"}`, http.StatusBadRequest)
		return
	}
	if data.RefreshToken == "" {
		http.Error(w, `{"status": "requst body not contains refresh token"}`, http.StatusBadRequest)
		return
	}

	tokens := util.Tokens{}

	tokens, err = h.service.RefreshToken(r.Context(), data.RefreshToken)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"status": "error occured while refreshing token: %s"}`, err.Error()), http.StatusBadRequest)
		return
	}

	jsonResp, err := json.Marshal(tokens)
	if err != nil {
		http.Error(w, `{"status": "error occured while marshaling JSON"}`, http.StatusInternalServerError)
		return
	}

	w.Write(jsonResp)
}
