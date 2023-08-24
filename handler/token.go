package handler

import (
	"encoding/json"
	"net/http"

	"github.com/4epyx/testtask/service"
	"github.com/google/uuid"
)

type Tokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type TokenHandler struct {
	service *service.TokenService
}

func NewTokenHandler(service *service.TokenService) *TokenHandler {
	return &TokenHandler{service: service}
}

func (h *TokenHandler) GetAccessAndRefreshTokens(w http.ResponseWriter, r *http.Request) {
	guid := r.URL.Query().Get("guid")

	if _, err := uuid.Parse(guid); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"status": "invalid user GUID"}`))
		return
	}

	tokens := Tokens{}
	var err error

	tokens.AccessToken, err = h.service.GenerateAccessToken(r.Context(), guid)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"status": "error occured while generating access token"}`))
		return
	}

	tokens.RefreshToken, err = h.service.GenerateRefreshToken(r.Context(), guid)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"status": "error occured while generating refresh token"}`))
		return
	}

	jsonResp, err := json.Marshal(tokens)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"status": "error occured while marshaling JSON"}`))
		return
	}

	w.Write(jsonResp)
}

func (h *TokenHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"status": "invalid request body"}`))
		return
	}

	oldRefreshToken := r.Form.Get("refresh_token")
	if oldRefreshToken == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"status": "requst body not contains refresh token"}`))
		return
	}

	tokens := Tokens{}
	var err error

	tokens.AccessToken, tokens.RefreshToken, err = h.service.RefreshToken(r.Context(), oldRefreshToken)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"status": "error occured while refreshing token"}`))
		return
	}

	jsonResp, err := json.Marshal(tokens)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"status": "error occured while marshaling JSON"}`))
		return
	}

	w.Write(jsonResp)
}
