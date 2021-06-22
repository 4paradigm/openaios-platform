package main

import (
	"flag"
	"fmt"
	"github.com/4paradigm/openaios-platform/src/billing/apigen"
	"github.com/4paradigm/openaios-platform/src/billing/handler"
	"github.com/4paradigm/openaios-platform/src/internal/response"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
)

func main() {
	flag.Parse()
	fmt.Printf("FLAGS:\n")
	flag.VisitAll(func(f *flag.Flag) {
		fmt.Printf("%-25v : %v\n", f.Name, f.Value.String())
	})

	e := echo.New()
	e.Logger.SetHeader(`${time_rfc3339} ${level} ${short_file}:${line} `)
	e.Logger.SetLevel(log.WARN)
	e.HTTPErrorHandler = response.CustomHTTPErrorHandler

	e.Use(middleware.Logger())
	h := handler.NewHandler()
	handler.InitBillingServer()

	apiGroup := e.Group("/api")
	apigen.RegisterHandlers(apiGroup, &h)
	go heartbeat()
	e.Logger.Info(e.Start("0.0.0.0:80"))
}
