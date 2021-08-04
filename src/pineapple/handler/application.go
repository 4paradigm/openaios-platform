package handler

import (
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"github.com/4paradigm/openaios-platform/src/internal/response"
	"github.com/4paradigm/openaios-platform/src/pineapple/apigen"
	"github.com/4paradigm/openaios-platform/src/pineapple/controller/application"
	"github.com/4paradigm/openaios-platform/src/pineapple/utils/helm"
	"net/http"
	"regexp"
	"strings"
)

func (handler *Handler) GetApplicationList(ctx echo.Context, params apigen.GetApplicationListParams) error {
	userId := ctx.Get("userID").(string)
	bearerToken := ctx.Get("bearerToken").(string)
	appImpl, err := application.NewApplicationImpl(bearerToken, userId)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError).SetInternal(err)
	}
	var offset int
	var limit int
	if params.Offset == nil || params.Limit == nil {
		offset = -1
		limit = -1
	} else {
		offset = *params.Offset
		limit = *params.Limit
	}
	appPodInfo, err := appImpl.GetApplicationInstanceInfoList(limit, offset)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError).SetInternal(err)
	}
	return ctx.JSON(http.StatusOK, appPodInfo)
}

func (handler *Handler) DeleteApplication(ctx echo.Context, name string) error {
	userId := ctx.Get("userID").(string)
	bearerToken := ctx.Get("bearerToken").(string)
	appImpl, err := application.NewApplicationImpl(bearerToken, userId)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError).SetInternal(err)
	}
	if err := appImpl.Delete(name); err != nil {
		return response.BadRequestWithMessagef(ctx, "应用: %s 删除失败，请重试", name)
	}
	return response.StatusOKNoContent(ctx)
}

func (handler *Handler) CreateApplication(ctx echo.Context, name string) error {
	if len(name) >= 31 {
		return response.BadRequestWithMessagef(ctx, "应用名%s过长，请修改名称小于31字符后重试", name)
	}
	regex, err := regexp.Compile("[a-z]([-a-z0-9]*[a-z0-9])?")
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError).SetInternal(err)
	}
	if !regex.MatchString(string(name)) {
		return response.BadRequestWithMessagef(ctx, "应用名%s不合法，须符合正则表达式'[a-z]([-a-z0-9]*[a-z0-9])?'", name)
	}
	userId := ctx.Get("userID").(string)
	bearerToken := ctx.Get("bearerToken").(string)
	requestBody := new(apigen.CreateApplicationJSONRequestBody)
	if err := ctx.Bind(requestBody); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError).SetInternal(err)
	}

	if err := checkCreateApplicationJSONRequestBody(requestBody); err != nil {
		return response.BadRequestWithMessage(ctx, err.Error())
	}

	appInfo, err := createApplicationInfo(requestBody, name, userId)
	if err != nil {
		if strings.Contains(err.Error(), "User does not have enough balance") {
			return response.BadRequestWithMessage(ctx, "您的余额不足，无法创建应用")
		}
		return echo.NewHTTPError(http.StatusInternalServerError).SetInternal(err)
	}

	appImpl, err := application.NewApplicationImpl(bearerToken, userId)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError).SetInternal(err)
	}
	chart, err := helm.DownloadChart(*requestBody.Url)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError).SetInternal(err)
	}
	_, err = appImpl.Create(chart, appInfo)
	if err != nil {
		if strings.Contains(err.Error(), "cannot re-use a name that is still in use") {
			return response.BadRequestWithMessagef(ctx, "应用%s已经存在，请修改名称后重试", name)
		}
		return echo.NewHTTPError(http.StatusInternalServerError).SetInternal(err)
	}
	return response.StatusOKNoContent(ctx)
}

func checkCreateApplicationJSONRequestBody(requestBody *apigen.CreateApplicationJSONRequestBody) error {
	if requestBody.Answers == nil {
		return errors.New("Answers 无效")
	}
	if requestBody.Url == nil {
		return errors.New("Url 无效")
	}
	return nil
}

func createApplicationInfo(requestBody *apigen.CreateApplicationJSONRequestBody, appName string, userId string) (*application.ApplicationInfo, error) {
	pineappleInfo, err := helm.NewPineappleInfo(appName, userId, application.AppPrefix)
	if err != nil {
		return nil, err
	}
	appInfo := application.ApplicationInfo{
		pineappleInfo,
	}
	appInfo.SetValues(*requestBody.Answers)

	return &appInfo, nil
}

func (handler *Handler) GetApplicationMetadata(ctx echo.Context, name string) error {
	userId := ctx.Get("userID").(string)
	bearerToken := ctx.Get("bearerToken").(string)
	appImpl, err := application.NewApplicationImpl(bearerToken, userId)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError).SetInternal(err)
	}
	appInstanceInfo, err := appImpl.GetApplicationInstanceInfo(name)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError).SetInternal(err)
	}
	return ctx.JSON(http.StatusOK, appInstanceInfo)
}

func (handler *Handler) GetApplicationPods(ctx echo.Context, name string) error {
	userId := ctx.Get("userID").(string)
	bearerToken := ctx.Get("bearerToken").(string)
	appImpl, err := application.NewApplicationImpl(bearerToken, userId)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError).SetInternal(err)
	}
	appPodInfo, err := appImpl.GetApplicationPodsInfo(name)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError).SetInternal(err)
	}
	return ctx.JSON(http.StatusOK, appPodInfo)
}

func (handler *Handler) GetApplicationServices(ctx echo.Context, name string) error {
	userId := ctx.Get("userID").(string)
	bearerToken := ctx.Get("bearerToken").(string)
	appImpl, err := application.NewApplicationImpl(bearerToken, userId)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError).SetInternal(err)
	}
	appSvcInfo, err := appImpl.GetApplicationServicesInfo(name)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError).SetInternal(err)
	}
	return ctx.JSON(http.StatusOK, appSvcInfo)
}

func (handler *Handler) GetApplicationNotes(ctx echo.Context, name string) error {
	userId := ctx.Get("userID").(string)
	bearerToken := ctx.Get("bearerToken").(string)
	appImpl, err := application.NewApplicationImpl(bearerToken, userId)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError).SetInternal(err)
	}
	notes, err := appImpl.GetApplicationReleaseNotes(name)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError).SetInternal(err)
	}
	return ctx.JSON(http.StatusOK, notes)
}
