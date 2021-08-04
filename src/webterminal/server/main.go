package main

import (
	"github.com/labstack/echo-contrib/prometheus"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"github.com/4paradigm/openaios-platform/src/internal/auth"
	"github.com/4paradigm/openaios-platform/src/webterminal/server/apigen"
	"github.com/4paradigm/openaios-platform/src/webterminal/server/handler"
)

func main() {
	e := echo.New()
	e.Logger.SetHeader(`${time_rfc3339} ${level} ${short_file}:${line} `)
	e.Logger.SetLevel(log.INFO)

	e.Use(middleware.Logger())

	p := prometheus.NewPrometheus("echo", nil)
	p.Use(e)

	//authVerifier, err := auth.NewAuthVerifier()
	//if err != nil {
	//	e.Logger.Fatalf("auth handler initialize failed, err = " + err.Error())
	//	return
	//}

	h := handler.NewHandler()

	apiGroup := e.Group("/web-terminal")

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

	apigen.RegisterHandlers(apiGroup, &h)

	e.Logger.Info(e.Start("0.0.0.0:80"))

}
