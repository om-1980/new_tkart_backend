package utils

import (
    "net/http"
    "github.com/labstack/echo/v4"
)

func InternalServerError(c echo.Context, err error) error {
    return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
}

func BadRequest(c echo.Context, msg string) error {
    return c.JSON(http.StatusBadRequest, echo.Map{"error": msg})
}
