package handler

import (
	"github.com/4paradigm/openaios-platform/src/internal/billingclient"
	"github.com/4paradigm/openaios-platform/src/pineapple/apigen"
	"github.com/4paradigm/openaios-platform/src/pineapple/conf"
	"github.com/4paradigm/openaios-platform/src/pineapple/utils"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"net/http"
)

func (handler *Handler) GetComputingUnitSpecs(ctx echo.Context) error {
	userId := ctx.Get("userID").(string)
	billingClient, err := billingclient.GetBillingClient(conf.GetBillingServerURL())
	if err != nil {
		return echo.NewHTTPError(
			http.StatusInternalServerError, "cannot connect to billing server").SetInternal(
			errors.Wrap(err, utils.GetRuntimeLocation()))
	}
	computeunitInfos, err := billingclient.GetComputeUnitListByUserID(billingClient, userId)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError).SetInternal(err)
	}
	result := make([]apigen.ComputeUnitSpec, len(computeunitInfos))

	for i, info := range computeunitInfos {
		id := apigen.ComputeUnitId(*(info.Id))
		result[i].Name = info.Id
		result[i].Id = &id
		result[i].Description = info.Description
		result[i].Price = info.Price
	}

	return ctx.JSON(http.StatusOK, result)
}
