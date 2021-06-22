package handler

import (
	"encoding/json"
	"github.com/4paradigm/openaios-platform/src/billing/apigen"
	"github.com/4paradigm/openaios-platform/src/billing/utils"
	"github.com/4paradigm/openaios-platform/src/internal/mongodb"
	"github.com/4paradigm/openaios-platform/src/internal/response"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"net/http"
	"strings"
)

func (h Handler) GetAccount(ctx echo.Context) error {
	client, err := mongodb.GetMongodbClient(mongodbUrl)
	defer mongodb.KillMongodbClient(client)
	if err != nil {
		return echo.NewHTTPError(
			http.StatusInternalServerError, "cannot connect to mongodb.").SetInternal(
			errors.Wrap(err, response.GetRuntimeLocation()))
	}

	accountList, err := utils.GetAccountList(client)
	if err != nil {
		return echo.NewHTTPError(
			http.StatusInternalServerError, "get account list failed.").SetInternal(
			errors.Wrap(err, "get account list failed "+response.GetRuntimeLocation()))
	}

	newAccountList := []apigen.AccountInfo{}
	for _, item := range accountList {
		userID := item.UserId
		balance := item.Balance
		callbackURL := item.CallbackUrl
		computeunitGroup := item.ComputeunitGroup
		accountInfo := apigen.AccountInfo{
			UserID:           &userID,
			Balance:          balance,
			CallbackUrl:      &callbackURL,
			ComputeunitGroup: &computeunitGroup}
		newAccountList = append(newAccountList, accountInfo)
	}
	return ctx.JSON(http.StatusOK, newAccountList)
}

func (h Handler) PutAccountUserid(ctx echo.Context, userid string) error {
	client, err := mongodb.GetMongodbClient(mongodbUrl)
	defer mongodb.KillMongodbClient(client)
	if err != nil {
		return echo.NewHTTPError(
			http.StatusInternalServerError, "cannot connect to mongodb.").SetInternal(
			errors.Wrap(err, "cannot connect to mongodb "+response.GetRuntimeLocation()))
	}

	var body = apigen.AccountInfo{}
	err = json.NewDecoder(ctx.Request().Body).Decode(&body)
	if err != nil {
		return response.BadRequestWithMessage(ctx, "cannot get request body")
	}

	err = utils.UpdateUserAccount(client, userid, body.Balance, body.CallbackUrl, body.ComputeunitGroup)
	if err != nil {
		//if strings.Contains(err.Error(), "cannot find such user") {
		//	return utils.BadRequestWithMessage(ctx, "cannot find such user.")
		//} else {
		return echo.NewHTTPError(
			http.StatusInternalServerError, "update user account failed.").SetInternal(
			errors.Wrap(err, "update user account failed. "+response.GetRuntimeLocation()))
	}
	return response.StatusOKNoContent(ctx)
}

func (h Handler) DeleteAccountUserid(ctx echo.Context, userid string) error {
	client, err := mongodb.GetMongodbClient(mongodbUrl)
	defer mongodb.KillMongodbClient(client)
	if err != nil {
		return echo.NewHTTPError(
			http.StatusInternalServerError, "cannot connect to mongodb").SetInternal(
			errors.Wrap(err, response.GetRuntimeLocation()))
	}

	err = utils.DeleteUser(client, userid)
	if err != nil {
		return echo.NewHTTPError(
			http.StatusInternalServerError, "cannot delete user").SetInternal(
			errors.Wrap(err, response.GetRuntimeLocation()))
	}
	return response.StatusOKNoContent(ctx)
}

func (h Handler) PostAccountUserid(ctx echo.Context, userid string, params apigen.PostAccountUseridParams) error {
	client, err := mongodb.GetMongodbClient(mongodbUrl)
	defer mongodb.KillMongodbClient(client)
	if err != nil {
		return echo.NewHTTPError(
			http.StatusInternalServerError, "cannot connect to mongodb").SetInternal(
			errors.Wrap(err, response.GetRuntimeLocation()))
	}

	var balance = 0.0
	if params.Balance != nil {
		balance = *params.Balance
	}
	err = utils.CreateUserWithBalance(client, userid, balance, params.CallbackUrl)
	if err != nil && !strings.Contains(err.Error(), "duplicate key error") {
		return echo.NewHTTPError(
			http.StatusInternalServerError, "cannot create user").SetInternal(
			errors.Wrap(err, response.GetRuntimeLocation()))
	}
	return response.StatusOKNoContent(ctx)
}

func (h Handler) GetAccountUseridBalance(ctx echo.Context, userid string) error {
	client, err := mongodb.GetMongodbClient(mongodbUrl)
	defer mongodb.KillMongodbClient(client)
	if err != nil {
		return echo.NewHTTPError(
			http.StatusInternalServerError, "cannot connect to mongodb").SetInternal(
			errors.Wrap(err, response.GetRuntimeLocation()))
	}

	balance, err := utils.GetUserBalance(client, userid)
	if err != nil {
		return echo.NewHTTPError(
			http.StatusInternalServerError, "cannot get user balance").SetInternal(
			errors.Wrap(err, response.GetRuntimeLocation()))
	}
	return ctx.JSON(http.StatusOK, balance)
}

func (h Handler) PostAccountUseridBalance(ctx echo.Context, userid string, params apigen.PostAccountUseridBalanceParams) error {
	client, err := mongodb.GetMongodbClient(mongodbUrl)
	defer mongodb.KillMongodbClient(client)
	if err != nil {
		return echo.NewHTTPError(
			http.StatusInternalServerError, "cannot connect to mongodb").SetInternal(
			errors.Wrap(err, response.GetRuntimeLocation()))
	}

	balance := params.BuyBalance
	err = utils.ModifyUserBalance(client, userid, balance)
	if err != nil {
		return echo.NewHTTPError(
			http.StatusInternalServerError, "cannot buy balance for user").SetInternal(
			errors.Wrap(err, response.GetRuntimeLocation()))
	}
	return response.StatusOKNoContent(ctx)
}
