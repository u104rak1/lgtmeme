package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/ucho456job/lgtmeme/config"
	resourceDto "github.com/ucho456job/lgtmeme/internal/resource/dto"
)

type ImageService interface {
	Post(c echo.Context, b resourceDto.PostImageReqBody, token string) (respBody *resourceDto.PostImageResp, status int, err error)
	BulkGet(c echo.Context, q resourceDto.GetImagesQuery, token string) (respBody *resourceDto.GetImagesResp, status int, err error)
	Patch(c echo.Context, b resourceDto.PatchImageReqBody, imageID, token string) (status int, err error)
}

type imageService struct{}

func NewImageService() ImageService {
	return &imageService{}
}

func (s *imageService) Post(c echo.Context, b resourceDto.PostImageReqBody, token string) (respBody *resourceDto.PostImageResp, status int, err error) {
	baseURL := os.Getenv("BASE_URL")

	url := fmt.Sprintf("%s%s", baseURL, config.RESOURCE_IMAGES_ENDPOINT)

	body, err := json.Marshal(b)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	req, err := http.NewRequest("POST", url, strings.NewReader(string(body)))
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, http.StatusServiceUnavailable, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return nil, resp.StatusCode, errors.New("failed to create image")
	}

	respBody = new(resourceDto.PostImageResp)
	if err := json.NewDecoder(resp.Body).Decode(respBody); err != nil {
		return nil, http.StatusInternalServerError, err
	}

	return respBody, resp.StatusCode, nil
}

func (s *imageService) BulkGet(c echo.Context, q resourceDto.GetImagesQuery, token string) (respBody *resourceDto.GetImagesResp, status int, err error) {
	var reqDataParts []string
	reqDataParts = append(reqDataParts, fmt.Sprintf("page=%d", q.Page))
	reqDataParts = append(reqDataParts, fmt.Sprintf("keyword=%s", q.Keyword))
	reqDataParts = append(reqDataParts, fmt.Sprintf("sort=%s", q.Sort))
	reqDataParts = append(reqDataParts, fmt.Sprintf("favorite_image_ids=%s", q.FavoriteImageIDs))
	reqDataParts = append(reqDataParts, fmt.Sprintf("auth_check=%t", q.AuthCheck))
	reqData := strings.Join(reqDataParts, "&")

	baseURL := os.Getenv("BASE_URL")

	url := fmt.Sprintf("%s%s?%s", baseURL, config.RESOURCE_IMAGES_ENDPOINT, reqData)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, http.StatusServiceUnavailable, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, resp.StatusCode, errors.New("failed to get images")
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

func (s *imageService) Patch(c echo.Context, b resourceDto.PatchImageReqBody, imageID, token string) (status int, err error) {
	baseURL := os.Getenv("BASE_URL")

	url := fmt.Sprintf("%s%s/%s", baseURL, config.RESOURCE_IMAGES_ENDPOINT, imageID)

	body, err := json.Marshal(b)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	req, err := http.NewRequest("PATCH", url, strings.NewReader(string(body)))
	if err != nil {
		return http.StatusInternalServerError, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		return http.StatusServiceUnavailable, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return resp.StatusCode, errors.New("failed to update image")
	}

	return resp.StatusCode, nil
}
