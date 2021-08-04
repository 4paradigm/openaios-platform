package handler

import (
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"github.com/4paradigm/openaios-platform/src/billing/apigen"
	"github.com/4paradigm/openaios-platform/src/billing/utils"
	"github.com/4paradigm/openaios-platform/src/internal/mongodb"
	"github.com/4paradigm/openaios-platform/src/internal/response"
	"net/http"
)

func (h Handler) GetComputeunitUserid(ctx echo.Context, userid string) error {
	client, err := mongodb.GetMongodbClient(mongodbUrl)
	defer mongodb.KillMongodbClient(client)
	if err != nil {
		return echo.NewHTTPError(
			http.StatusInternalServerError, "cannot connect to mongodb.").SetInternal(
			errors.Wrap(err, "cannot connect to mongodb "+response.GetRuntimeLocation()))
	}

	groupList, err := utils.GetUserComputeunitGroup(client, userid)
	if err != nil {
		return echo.NewHTTPError(
			http.StatusInternalServerError, "get computeunit group failed.").SetInternal(
			errors.Wrap(err, "get computeunit group failed "+response.GetRuntimeLocation()))
	}
	computeunitMap := map[string]bool{}
	computeunitList := []map[string]interface{}{}
	for _, groupName := range groupList {
		groupList, avl, err := utils.GetComputeunitInfoByGroup(client, groupName)
		if err != nil {
			ctx.Logger().Warn(err)
			continue
		}
		if !avl {
			continue
		}
		for _, item := range groupList {
			if _, ok := computeunitMap[item["id"].(string)]; ok {
				continue
			}
			computeunitMap[item["id"].(string)] = true
			computeunitList = append(computeunitList, item)
		}
	}
	return ctx.JSON(http.StatusOK, computeunitList)
}

func (h Handler) DeleteComputeunitUserid(ctx echo.Context, userid string, params apigen.DeleteComputeunitUseridParams) error {
	if params.GroupName == "default" {
		return response.BadRequestWithMessage(ctx, "cannot delete group default.")
	}

	client, err := mongodb.GetMongodbClient(mongodbUrl)
	defer mongodb.KillMongodbClient(client)
	if err != nil {
		return echo.NewHTTPError(
			http.StatusInternalServerError, "cannot connect to mongodb.").SetInternal(
			errors.Wrap(err, response.GetRuntimeLocation()))
	}

	groupList, err := utils.GetUserComputeunitGroup(client, userid)
	if err != nil {
		return echo.NewHTTPError(
			http.StatusInternalServerError, "get computeunit group failed.").SetInternal(
			errors.Wrap(err, "get computeunit group failed "+response.GetRuntimeLocation()))
	}

	index := 0
	for _, val := range groupList {
		if val != params.GroupName {
			groupList[index] = val
			index = index + 1
		}
	}

	err = utils.ModifyUserComputeunitGroup(client, userid, groupList[:index])
	if err != nil {
		return echo.NewHTTPError(
			http.StatusInternalServerError, "update user computeunit group failed").SetInternal(
			errors.Wrap(err, response.GetRuntimeLocation()))
	}
	return response.StatusOKNoContent(ctx)
}

func (h Handler) PostComputeunitUserid(ctx echo.Context, userid string, params apigen.PostComputeunitUseridParams) error {
	if params.GroupName == "default" {
		return response.StatusOKNoContent(ctx)
	}

	client, err := mongodb.GetMongodbClient(mongodbUrl)
	defer mongodb.KillMongodbClient(client)
	if err != nil {
		return echo.NewHTTPError(
			http.StatusInternalServerError, "cannot connect to mongodb.").SetInternal(
			errors.Wrap(err, response.GetRuntimeLocation()))
	}

	groupList, err := utils.GetUserComputeunitGroup(client, userid)
	if err != nil {
		return echo.NewHTTPError(
			http.StatusInternalServerError, "get computeunit group failed.").SetInternal(
			errors.Wrap(err, "get computeunit group failed "+response.GetRuntimeLocation()))
	}
	if in, err := utils.In(groupList, params.GroupName); err != nil {
		return echo.NewHTTPError(
			http.StatusInternalServerError, "get computeunit group failed.").SetInternal(
			errors.Wrap(err, response.GetRuntimeLocation()))
	} else if in {
		return response.StatusOKNoContent(ctx)
	}

	groupList = append(groupList, params.GroupName)
	err = utils.ModifyUserComputeunitGroup(client, userid, groupList)
	if err != nil {
		return echo.NewHTTPError(
			http.StatusInternalServerError, "update user computeunit group failed").SetInternal(
			errors.Wrap(err, response.GetRuntimeLocation()))
	}
	return response.StatusOKNoContent(ctx)
}

func (h Handler) GetComputeunitUseridComputeunitIdComputeunitId(ctx echo.Context, userid string, computeunitId string) error {
	client, err := mongodb.GetMongodbClient(mongodbUrl)
	defer mongodb.KillMongodbClient(client)
	if err != nil {
		return echo.NewHTTPError(
			http.StatusInternalServerError, "cannot connect to mongodb.").SetInternal(
			errors.Wrap(err, response.GetRuntimeLocation()))
	}

	groupList, err := utils.GetUserComputeunitGroup(client, userid)
	if err != nil {
		return echo.NewHTTPError(
			http.StatusInternalServerError, "get computeunit group failed.").SetInternal(
			errors.Wrap(err, "get computeunit group failed "+response.GetRuntimeLocation()))
	}
	for _, groupName := range groupList {
		result, avl, err := utils.GetComputeunitInGroupByID(client, groupName, computeunitId)
		if err != nil {
			continue
		}
		if !avl {
			continue
		}
		return ctx.JSON(http.StatusOK, result)
	}
	return response.BadRequestWithMessage(ctx, "computeunit not exists")
}

func (h Handler) GetComputeunitGroupGroupName(ctx echo.Context, groupName string) error {
	client, err := mongodb.GetMongodbClient(mongodbUrl)
	defer mongodb.KillMongodbClient(client)
	if err != nil {
		return echo.NewHTTPError(
			http.StatusInternalServerError, "cannot connect to mongodb.").SetInternal(
			errors.Wrap(err, response.GetRuntimeLocation()))
	}

	computeunitList, avl, err := utils.GetComputeunitInfoByGroup(client, groupName)
	if err != nil {
		return echo.NewHTTPError(
			http.StatusInternalServerError, "get computeunit failed.").SetInternal(
			errors.Wrap(err, response.GetRuntimeLocation()))
	}
	if !avl {
		return response.BadRequestWithMessage(ctx, "computeunit group is not available.")
	}
	return ctx.JSON(http.StatusOK, computeunitList)
}

func (h Handler) GetComputeunitPrice(ctx echo.Context, params apigen.GetComputeunitPriceParams) error {
	client, err := mongodb.GetMongodbClient(mongodbUrl)
	defer mongodb.KillMongodbClient(client)
	if err != nil {
		return echo.NewHTTPError(
			http.StatusInternalServerError, "cannot connect to mongodb.").SetInternal(
			errors.Wrap(err, response.GetRuntimeLocation()))
	}

	priceMap, err := utils.GetPriceMap(client)
	if err != nil {
		return echo.NewHTTPError(
			http.StatusInternalServerError, "cannot get price map.").SetInternal(
			errors.Wrap(err, response.GetRuntimeLocation()))
	}

	computeUnitPrice, ok := priceMap[params.ComputeunitId]
	if !ok {
		return response.BadRequestWithMessage(ctx, "computeunit not exists.")
	}
	return ctx.JSON(http.StatusOK, computeUnitPrice)
}
