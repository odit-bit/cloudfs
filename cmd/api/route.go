package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

const (
	listHTML     = "/list"
	loginHTML    = "/login"
	registerHTML = "/register"
	uploadHTML   = "/upload"
)

const (
	listAPIEndpoint     = "/api/list"
	loginAPIEndpoint    = "/api/login"
	registerAPIEndpoint = "/api/register"
	uploadAPIEndpoint   = "/api/upload"
)

func httpServer(addr string, app *api) {
	e := echo.New()

	// e.Use(middleware.Logger())
	e.Use(middleware.CORS())
	e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(5)))

	// html page endpoint
	e.GET("/health", func(c echo.Context) error {
		return c.String(200, "health")
	})

	e.GET("/", func(c echo.Context) error {
		app.serveIndex(c.Response(), c.Request())
		return nil
	})

	e.GET(loginHTML, func(c echo.Context) error {
		app.serveLoginPage(loginAPIEndpoint).ServeHTTP(c.Response(), c.Request())
		return nil
	})

	e.GET(listHTML, func(c echo.Context) error {
		app.serveListPage(listAPIEndpoint).ServeHTTP(c.Response(), c.Request())
		return nil
	})

	e.GET(uploadHTML, func(c echo.Context) error {
		app.serveUploadPage(uploadAPIEndpoint).ServeHTTP(c.Response(), c.Request())
		return nil
	})

	e.GET(registerHTML, func(c echo.Context) error {
		app.serveRegisterPage(registerAPIEndpoint).ServeHTTP(c.Response(), c.Request())
		return nil
	})

	//=============================================================
	// service endpoint

	//user related

	e.POST(loginAPIEndpoint, func(c echo.Context) error {
		app.handleLogin(c.Response(), c.Request())
		return nil
	})

	e.POST(registerAPIEndpoint, func(c echo.Context) error {
		app.handleRegister(c.Response(), c.Request())
		return nil
	})

	//storage related

	e.GET(listAPIEndpoint, func(c echo.Context) error {
		app.HandleList(c.Response(), c.Request())
		return nil
	})

	e.POST(uploadAPIEndpoint, func(c echo.Context) error {
		app.handleUpload(c.Response(), c.Request())
		return nil
	})

	e.GET("/api/download", func(c echo.Context) error {
		app.handleDownload(c.Response(), c.Request())
		return nil
	})

	e.Start(addr)
}
