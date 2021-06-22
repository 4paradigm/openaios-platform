package handler

import (
	"github.com/4paradigm/openaios-platform/src/webterminal/server/apigen"
	"github.com/4paradigm/openaios-platform/src/webterminal/server/controller/webterminal"
	"github.com/labstack/echo/v4"
	"net/http"
)

func (h *Handler) GetTerminal(ctx echo.Context, params apigen.GetTerminalParams) error {
	userId := ctx.Get("userID").(string)
	bearerToken := ctx.Get("bearerToken").(string)

	podName := params.Pod
	containerName := params.Container

	terminal := webterminal.Webterminal{
		Namespace:     userId,
		PodName:       podName,
		ContainerName: containerName,
		BearerToken:   bearerToken,
	}

	terminalUrl, err := webterminal.GetWebterminal(terminal)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError).SetInternal(err)
	}

	//return ctx.Redirect(http.StatusSeeOther, terminalUrl)
	return ctx.JSON(http.StatusOK, apigen.WebterminalInfo{Url: &terminalUrl})
}
