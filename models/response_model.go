package models

import (
	"github.com/labstack/echo/v4"
)

type Response struct {
	Status  int       `json:"status"`
	Message string    `json:"message"`
	Body    *echo.Map `json:"body"`
}
