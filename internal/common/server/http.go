package server

import (
	"context"
	"github.com/dejano-with-tie/fantastigo/config"
	openapi "github.com/go-openapi/runtime/middleware"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func Run(cfg config.Config, opts ...func(e *echo.Echo)) {
	e := echo.New()

	// Customize
	e.Validator = NewValidator()
	e.Binder = NewBinder()
	e.HTTPErrorHandler = errHandler

	setMiddlewares(e)

	if cfg.Server.SwaggerUi {
		openApiHandler(cfg.Server.OpenapiSpec, e)
	}
	for _, o := range opts {
		o(e)
	}

	// Sets Log level to debug, nothing fancy
	e.Debug = cfg.Server.Debug

	// Start server in a separate goroutine
	go func() {
		if err := e.Start(cfg.Server.Addr); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("shutting down the server")
		}
	}()

	shutdownGracefully(e)
}

func setMiddlewares(e *echo.Echo) {
	e.Use(
		middleware.Recover(),   // Recover from all panics to always have your server up
		middleware.Logger(),    // Log everything to stdout
		middleware.RequestID(), // Generate a request id on the HTTP response headers for identification
	)
}

func shutdownGracefully(e *echo.Echo) {
	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 10 seconds.
	// Use a buffered channel to avoid missing signals as recommended for signal.Notify
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}

func openApiHandler(_ []string, e *echo.Echo) {
	e.Static("/docs/openapi", "docs/openapi/")
	// TODO use provided specs arg
	opts := openapi.SwaggerUIOpts{SpecURL: "docs/openapi/fleet.yaml"}
	sh := openapi.SwaggerUI(opts, nil)
	e.GET("/docs", echo.WrapHandler(sh))
}
