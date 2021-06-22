package handler

import (
	"github.com/4paradigm/openaios-platform/src/billing/apigen"
	"github.com/4paradigm/openaios-platform/src/billing/utils"
	"github.com/4paradigm/openaios-platform/src/internal/mongodb"
	"github.com/4paradigm/openaios-platform/src/internal/response"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"net/http"
)

func (h Handler) PostInstance(ctx echo.Context, params apigen.PostInstanceParams) error {
	client, err := mongodb.GetMongodbClient(mongodbUrl)
	defer mongodb.KillMongodbClient(client)
	if err != nil {
		log.Error(err.Error())
		return ctx.String(http.StatusInternalServerError, "cannot connect to mongodb.")
	}

	instanceId, err := utils.CreateInstance(client, string(params.UserId), params.InstanceName)
	if err != nil {
		log.Error(err.Error())
		return response.BadRequestWithMessage(ctx, err.Error())
	}
	return ctx.String(http.StatusOK, instanceId)
}
