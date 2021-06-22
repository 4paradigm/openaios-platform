package handler

import (
	"flag"
	"github.com/4paradigm/openaios-platform/src/pineapple/utils"
	"github.com/labstack/echo/v4"
	"net/http"
)

var (
	version = flag.String("version", utils.GetEnvDefault("PINEAPPLE_VERSION", "unknown"),
		"pineapple version, default to unknown")
)

func (handler *Handler) PineappleVersion(c echo.Context) error {
	return c.String(http.StatusOK, *version)
}
