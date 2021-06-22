package response

import (
	"flag"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"runtime"
)

var debug = flag.Bool("debug", false, "enable debug mode")

func BadRequestWithMessage(c echo.Context, message string) error {
	return c.JSON(http.StatusBadRequest, map[string]string{
		"message": message,
	})
}

func BadRequestWithMessagef(c echo.Context, format string, args ...interface{}) error {
	return c.JSON(http.StatusBadRequest, map[string]string{
		"message": fmt.Sprintf(format, args...),
	})
}

func BadRequestWithMessageWithJson(c echo.Context, message string, content interface{}) error {
	return c.JSON(http.StatusBadRequest, map[string]interface{}{
		"message": message,
		"type":    "json",
		"content": content,
	})
}

func StatusOKNoContent(c echo.Context) error {
	return c.String(http.StatusOK, "success")
}

func CustomHTTPErrorHandler(err error, c echo.Context) {
	he, ok := err.(*echo.HTTPError)
	if !ok {
		he = &echo.HTTPError{
			Code:    http.StatusInternalServerError,
			Message: http.StatusText(http.StatusInternalServerError),
		}
	}

	// Issue #1426
	code := he.Code
	message := he.Message
	if Debug() {
		message = err.Error()
	} else if m, ok := message.(string); ok {
		message = map[string]string{"message": m}
	}

	// Send response
	if !c.Response().Committed {
		if code == http.StatusInternalServerError {
			if he.Internal != nil {
				c.Logger().Error(he.Internal.Error())
			} else {
				c.Logger().Error("http error without internal: " + err.Error())
			}
		}

		if c.Request().Method == http.MethodHead { // Issue #608
			err = c.NoContent(he.Code)
		} else {
			err = c.JSON(code, message)
		}
		if err != nil {
			c.Logger().Error(err)
		}
	}
}

func Debug() bool {
	return *debug
}

func GetRuntimeLocation() string {
	_, fn, line, _ := runtime.Caller(1)
	return fmt.Sprintf("%s:%d", fn, line)
}
