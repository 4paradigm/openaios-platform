package main

import (
	"flag"
	"fmt"
	"github.com/labstack/echo-contrib/prometheus"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"github.com/4paradigm/openaios-platform/src/internal/auth"
	"github.com/4paradigm/openaios-platform/src/internal/response"
	"github.com/4paradigm/openaios-platform/src/pineapple/apigen"
	"github.com/4paradigm/openaios-platform/src/pineapple/apigen/internalapigen"
	"github.com/4paradigm/openaios-platform/src/pineapple/handler"
	"net/http"
	"os"
)

var (
	staticDir = flag.String("static-dir", os.Getenv("PINEAPPLE_STATIC_DIR"),
		"static directory")
)

func main() {
	flag.Parse()
	fmt.Printf("FLAGS:\n")
	flag.VisitAll(func(f *flag.Flag) {
		fmt.Printf("%-25v : %v\n", f.Name, f.Value.String())
	})

	e := echo.New()
	e.Logger.SetHeader(`${time_rfc3339} ${level} ${short_file}:${line} `)
	e.Logger.SetLevel(log.INFO)
	e.HTTPErrorHandler = response.CustomHTTPErrorHandler

	e.Use(middleware.Logger())

	p := prometheus.NewPrometheus("echo", nil)
	p.Use(e)

	h := handler.Handler{}

	apiGroup := e.Group("/api")
	if *staticDir != "" {
		e.Logger.Info("using static dir")
		apiGroup.Static("/static", *staticDir)
	}

	if err := auth.InitAuth(); err != nil {
		panic(err.Error())
	}
	apiGroup.Use(middleware.KeyAuth(func(key string, c echo.Context) (bool, error) {
		c.Set("bearerToken", key)
		idTokenClaim, err := auth.Verify(key)
		if err != nil {
			c.Logger().Warn("auth verify failed, err = " + err.Error())
			return false, err
		}
		c.Set("userName", idTokenClaim.PreferredUserName)
		c.Set("userID", idTokenClaim.Sub)
		return true, nil
	}))
	apiGroup.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if c.Get("userID") == nil {
				return next(c)
			}
			if err := InitUser(c); err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError).SetInternal(err)
			}
			return next(c)
		}
	})

	apiGroup.GET("/userinfo", h.UserInfo)
	apigen.RegisterHandlers(apiGroup, &h)

	internalApiGroup := e.Group("/internal-api")
	internalapigen.RegisterHandlers(internalApiGroup, &h)

	e.Logger.Info(e.Start("0.0.0.0:80"))
}
