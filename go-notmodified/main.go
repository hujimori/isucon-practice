package main

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

type Response struct {
	Icon string `json:"icon"`
}

func main() {
	e := echo.New()
	e.Debug = true
	e.Logger.SetLevel(log.DEBUG)

	e.GET("/api/icon", getIsuIcon)
	e.GET("/api/hello", hello)
	serverPort := fmt.Sprintf(":%v", "3000")
	e.Logger.Fatal(e.Start(serverPort))
}

func getIsuIcon(c echo.Context) error {
	c.Response().Header().Set("Cache-Control", "max-age=180")
	res := Response{
		Icon: "icon",
	}
	return c.JSON(http.StatusOK, res)
	// return c.String(http.StatusOK, "Hello, World!")
}

func hello(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}
