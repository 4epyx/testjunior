package handler

import (
	"encoding/json"
	"fmt"
	"io"
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
	type bodyData struct {
		UserGuid string `json:"user_guid"`
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"status": "couldn't read body"}`))
		return
	}

	data := &bodyData{}
	if err := json.Unmarshal(body, data); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"status": "couldn't parse body to JSON"}`))
		return
	}

	if _, err := uuid.Parse(data.UserGuid); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf(`{"status": "invalid user GUID %s"}`, data.UserGuid)))
		return
	}

	tokens := Tokens{}

	tokens.AccessToken, err = h.service.GenerateAccessToken(r.Context(), data.UserGuid)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"status": "error occured while generating access token"}`))
		return
	}

	tokens.RefreshToken, err = h.service.GenerateRefreshToken(r.Context(), data.UserGuid)
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
	type bodyData struct {
		RefreshToken string `json:"refresh_token"`
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"status": "couldn't read body"}`))
		return
	}

	data := &bodyData{}
	if err := json.Unmarshal(body, data); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"status": "couldn't parse body to JSON"}`))
		return
	}
	if data.RefreshToken == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"status": "requst body not contains refresh token"}`))
		return
	}

	tokens := Tokens{}

	tokens.AccessToken, tokens.RefreshToken, err = h.service.RefreshToken(r.Context(), data.RefreshToken)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf(`{"status": "error occured while refreshing token: %s"}`, err.Error())))
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
