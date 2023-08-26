package web

import (
	"github.com/labstack/echo/v4"
)

type Route struct {
	Method, Url  string
	Handler      echo.HandlerFunc
	IsRestricted bool
}

type Routes []Route
