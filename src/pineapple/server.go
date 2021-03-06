/*
 * Copyright © 2021 peizhaoyou <peizhaoyou@4paradigm.com>
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import (
	"flag"
	"fmt"
	"github.com/4paradigm/openaios-platform/src/internal/auth"
	"github.com/4paradigm/openaios-platform/src/internal/response"
	"github.com/4paradigm/openaios-platform/src/internal/version"
	"github.com/4paradigm/openaios-platform/src/pineapple/apigen"
	"github.com/4paradigm/openaios-platform/src/pineapple/apigen/internalapigen"
	"github.com/4paradigm/openaios-platform/src/pineapple/handler"
	"github.com/labstack/echo-contrib/prometheus"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"net/http"
	"os"
)

var (
	staticDir = flag.String("static-dir", os.Getenv("PINEAPPLE_STATIC_DIR"),
		"static directory")
)

func main() {
	flag.Parse()
	version.CheckVersionFlag()

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

	internalAPIGroup := e.Group("/internal-api")
	internalapigen.RegisterHandlers(internalAPIGroup, &h)

	e.Logger.Info(e.Start("0.0.0.0:80"))
}
