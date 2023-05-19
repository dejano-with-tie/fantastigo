package main

import (
	"flag"
	"github.com/dejano-with-tie/fantastigo/config"
	"github.com/dejano-with-tie/fantastigo/internal/common/server"
	"github.com/dejano-with-tie/fantastigo/internal/ector/app"
	"github.com/dejano-with-tie/fantastigo/internal/ector/metrics"
	inner "github.com/dejano-with-tie/fantastigo/internal/ector/server"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

func main() {
	var activeEnv string
	flag.StringVar(&activeEnv, "env", "", "active environment: dev, qa, prod")
	flag.Parse()
	cfg := config.MustLoad(config.NewLoadConfig(activeEnv))
	v := validator.New()

	storer := metrics.NewStore()
	application := app.NewEctor(storer)

	server.Run(cfg, v, func(e *echo.Echo) {
		g := e.Group("/api")
		inner.RegisterHandlers(g, inner.NewEctorHttpHandler(application))
	})
}
