package response

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/ucho456job/lgtmeme/config"
)

type ErrResp struct {
	ErrCode string `json:"errorCode"`
	ErrMsg  string `json:"errorMessage"`
}

func HandleErrResp(c echo.Context, status int, err error) error {
	switch status {
	case http.StatusBadRequest:
		return BadRequest(c, err)
	case http.StatusUnauthorized:
		return Unauthorized(c, err)
	case http.StatusNotFound:
		return NotFound(c, err)
	default:
		return InternalServerError(c, err)
	}
}

func BadRequest(c echo.Context, err error) error {
	config.Logger.Warn("Bad request", "error", err.Error())
	return c.JSON(http.StatusBadRequest, ErrResp{
		ErrCode: "bad_request",
		ErrMsg:  err.Error(),
	})
}

func Unauthorized(c echo.Context, err error) error {
	config.Logger.Warn("Unauthorized", "error", err.Error())
	return c.JSON(http.StatusUnauthorized, ErrResp{
		ErrCode: "unauthorized",
		ErrMsg:  err.Error(),
	})
}

func NotFound(c echo.Context, err error) error {
	config.Logger.Warn("Not found", "error", err.Error())
	return c.JSON(http.StatusNotFound, ErrResp{
		ErrCode: "not_found",
		ErrMsg:  err.Error(),
	})
}

func InternalServerError(c echo.Context, err error) error {
	config.Logger.Error("Internal server error", "error", err.Error())
	return c.JSON(http.StatusInternalServerError, ErrResp{
		ErrCode: "internal_server_error",
		ErrMsg:  err.Error(),
	})
}
