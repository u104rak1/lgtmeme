package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/ucho456job/lgtmeme/config"
	"github.com/ucho456job/lgtmeme/internal/dto"
)

type GeneralAccessTokenService interface {
	CallToken(c echo.Context) (respBody *dto.ClientCredentialsResponse, status int, err error)
}

type generalAccessTokenService struct{}

func NewGeneralAccessTokenService() GeneralAccessTokenService {
	return &generalAccessTokenService{}
}

func (s *generalAccessTokenService) CallToken(c echo.Context) (respBody *dto.ClientCredentialsResponse, status int, err error) {
	baseURL := os.Getenv("BASE_URL")
	clientID := os.Getenv("GENERAL_CLIENT_ID")
	clientSecret := os.Getenv("GENERAL_CLIENT_SECRET")
	reqData := fmt.Sprintf("grant_type=client_credentials&client_id=%s&client_secret=%s", clientID, clientSecret)

	req, err := http.NewRequest("POST", baseURL+config.TOKEN_ENDPOINT, bytes.NewBufferString(reqData))
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, http.StatusServiceUnavailable, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, resp.StatusCode, errors.New("failed to get access token")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	if err := json.Unmarshal(body, &respBody); err != nil {
		return nil, http.StatusInternalServerError, err
	}

	return respBody, resp.StatusCode, nil
}
