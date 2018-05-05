package main

import (
	"net/http"

	"github.com/labstack/echo"
	// "github.com/labstack/echo/middleware"
)

func main() {
	// Echo instance
	e := echo.New()

	e.Static("/static", "client/dist/static")
	e.File("/", "client/dist/index.html")

	// Routes
	e.GET("/hello", hello)

	// Start server
	e.Logger.Fatal(e.Start(":80"))
}

// Handler
func hello(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}
