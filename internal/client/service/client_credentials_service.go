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
	authDto "github.com/ucho456job/lgtmeme/internal/auth/dto"
)

type ClientCredentialsService interface {
	GetAccessToken(c echo.Context) (accessToken string, statusCode int, err error)
}

type clientCredentialsService struct{}

func NewClientCredentialsService() ClientCredentialsService {
	return &clientCredentialsService{}
}

func (s *clientCredentialsService) GetAccessToken(c echo.Context) (accessToken string, status int, err error) {
	baseURL := os.Getenv("BASE_URL")
	clientID := os.Getenv("GENERAL_CLIENT_ID")
	clientSecret := os.Getenv("GENERAL_CLIENT_SECRET")
	reqData := fmt.Sprintf("grant_type=client_credentials&client_id=%s&client_secret=%s", clientID, clientSecret)

	req, err := http.NewRequest("POST", baseURL+config.TOKEN_ENDPOINT, bytes.NewBufferString(reqData))
	if err != nil {
		return "", http.StatusInternalServerError, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		return "", http.StatusServiceUnavailable, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", resp.StatusCode, errors.New("failed to get access token")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", http.StatusInternalServerError, err
	}

	respBody := authDto.ClientCredentialsResponse{}
	if err := json.Unmarshal(body, &respBody); err != nil {
		return "", http.StatusInternalServerError, err
	}

	return respBody.AccessToken, resp.StatusCode, nil
}
