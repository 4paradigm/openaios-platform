package handler

import (
	"encoding/json"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"github.com/4paradigm/openaios-platform/src/internal/billingclient"
	"github.com/4paradigm/openaios-platform/src/internal/response"
	"github.com/4paradigm/openaios-platform/src/pineapple/apigen"
	"github.com/4paradigm/openaios-platform/src/pineapple/conf"
	"github.com/4paradigm/openaios-platform/src/pineapple/controller/application"
	"github.com/4paradigm/openaios-platform/src/pineapple/controller/environment"
	"github.com/4paradigm/openaios-platform/src/pineapple/utils"
	"net/http"
)

func (handler *Handler) UserInfo(c echo.Context) error {
	idTokenMessage := c.Get("idToken").(*json.RawMessage)
	data, err := json.MarshalIndent(idTokenMessage, "", "    ")
	if err != nil {
		return echo.NewHTTPError(
			http.StatusInternalServerError).SetInternal(
			errors.Wrap(err, utils.GetRuntimeLocation()))
	}
	return c.String(http.StatusOK, string(data))
}

func (handler *Handler) GetUser(c echo.Context) error {
	userName := c.Get("userName").(string)
	userID := c.Get("userID").(string)
	billingClient, err := billingclient.GetBillingClient(conf.GetBillingServerURL())
	if err != nil {
		return echo.NewHTTPError(
			http.StatusInternalServerError, "cannot connect to billing server").SetInternal(
			errors.Wrap(err, utils.GetRuntimeLocation()))
	}
	balance, err := billingclient.GetUserBalance(billingClient, userID)
	if err != nil {
		return echo.NewHTTPError(
			http.StatusInternalServerError, "get user balance failed").SetInternal(
			errors.Wrap(err, utils.GetRuntimeLocation()))
	}
	return c.JSON(http.StatusOK, &apigen.UserInfo{
		Name:    &userName,
		Balance: balance,
	})
}

func (handler *Handler) GetUserTasks(c echo.Context) error {
	userID := c.Get("userID").(string)
	bearerToken := c.Get("bearerToken").(string)

	appClient, err := application.NewApplicationImpl(bearerToken, userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError).SetInternal(
			errors.Wrap(err, response.GetRuntimeLocation()))
	}
	appInfo, err := appClient.GetApplicationInstanceInfoList(0, 0)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError).SetInternal(
			errors.Wrap(err, response.GetRuntimeLocation()))
	}
	appTotal := int64(*appInfo.Total)

	envClient, err := environment.NewEnvironmentImpl(bearerToken, userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError).SetInternal(
			errors.Wrap(err, response.GetRuntimeLocation()))
	}
	envInfo, err := envClient.GetInfoList(0, 0)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError).SetInternal(
			errors.Wrap(err, response.GetRuntimeLocation()))
	}
	envTotal := int64(*envInfo.Total)

	client, err := utils.GetKubernetesClient()
	if err != nil {
		return echo.NewHTTPError(
			http.StatusInternalServerError, "get user tasks failed").SetInternal(
			errors.Wrap(err, "cannot get k8s client "+utils.GetRuntimeLocation()))
	}
	podList, err := utils.GetPodList(client, "", userID)
	if err != nil {
		return echo.NewHTTPError(
			http.StatusInternalServerError, "get user tasks failed").SetInternal(
			errors.Wrap(err, "cannot get user pod list "+utils.GetRuntimeLocation()))
	}

	computeUnitMap := map[string]int64{}
	for _, pod := range *podList {
		if pod.Status.Phase != "Running" || pod.DeletionTimestamp != nil {
			continue
		}
		var computeunitList []string
		computeunitString := pod.Annotations["openaios.4paradigm.com/computeunitList"]
		if computeunitString == "" {
			continue
		}
		err = json.Unmarshal([]byte(computeunitString), &computeunitList)
		if err != nil {
			c.Logger().Warn(err)
			continue
		}
		for _, computeunit := range computeunitList {
			computeUnitMap[computeunit] += 1
		}
	}
	computeUnitList := []apigen.UserTaskInfo{}
	billingClient, err := billingclient.GetBillingClient(conf.GetBillingServerURL())
	if err != nil {
		return echo.NewHTTPError(
			http.StatusInternalServerError, "cannot connect to billing server").SetInternal(
			errors.Wrap(err, utils.GetRuntimeLocation()))
	}
	for k, v := range computeUnitMap {
		var computeUnitId = apigen.ComputeUnitId(k)
		var computeUnitNum = v
		price, err := billingclient.GetComputeUnitPrice(billingClient, k)
		if err != nil {
			c.Logger().Warn(errors.Wrap(err, "get computeunit price failed "+utils.GetRuntimeLocation()))
		}
		computeUnitList = append(computeUnitList,
			apigen.UserTaskInfo{ComputeUnit: &computeUnitId, Number: &computeUnitNum, Price: &price})
	}
	UserTaskInfo := apigen.UserTasksInfo{EnvNum: &envTotal, AppNum: &appTotal, TaskList: &computeUnitList}
	return c.JSON(http.StatusOK, UserTaskInfo)
}

func (handler *Handler) PostUserInit(ctx echo.Context) error {
	return response.StatusOKNoContent(ctx)
}
