package handler

import (
	"flag"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"github.com/4paradigm/openaios-platform/src/internal/response"
	"github.com/4paradigm/openaios-platform/src/pineapple/apigen"
	"github.com/4paradigm/openaios-platform/src/pineapple/controller/environment"
	"github.com/4paradigm/openaios-platform/src/pineapple/handler/models"
	"github.com/4paradigm/openaios-platform/src/pineapple/utils"
	"github.com/4paradigm/openaios-platform/src/pineapple/utils/helm"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/chart/loader"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

var (
	chartsDir = flag.String("env-chartsdir", os.Getenv("PINEAPPLE_ENV_CHARTSDIR"),
		"env-chartsdir")
)

type RequestBodyError struct {
	keyErrors map[string]string
}

func (r *RequestBodyError) Error() string {
	return "请求无效：输入不合法"
}

func (r *RequestBodyError) GetKeyErrors() map[string]string {
	return r.keyErrors
}

func (handler *Handler) CreateEnvironment(ctx echo.Context, name apigen.EnvironmentName) error {
	if len(name) >= 31 {
		return response.BadRequestWithMessagef(ctx, "环境名%s过长，请修改名称小于31字符后重试", name)
	}
	regex, err := regexp.Compile("[a-z]([-a-z0-9]*[a-z0-9])?")
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError).SetInternal(err)
	}

	if !regex.MatchString(string(name)) {
		return response.BadRequestWithMessagef(ctx, "环境名%s不合法，须符合正则表达式'[a-z]([-a-z0-9]*[a-z0-9])?'", name)
	}

	userId := ctx.Get("userID").(string)
	bearerToken := ctx.Get("bearerToken").(string)
	requestBody := new(apigen.CreateEnvironmentJSONRequestBody)
	if err := ctx.Bind(requestBody); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError).SetInternal(err)
	}

	if requestErr := checkCreateEnvironmentJSONRequestBody(requestBody); requestErr != nil {
		return response.BadRequestWithMessageWithJson(ctx, requestErr.Error(), requestErr.GetKeyErrors())
	}

	envInfo, err := createEnvironmentInfo(requestBody, string(name), userId)
	if err != nil {
		if strings.Contains(err.Error(), "User does not have enough balance") {
			return response.BadRequestWithMessage(ctx, "您的余额不足，无法创建环境")
		}
		return echo.NewHTTPError(http.StatusInternalServerError).SetInternal(err)
	}

	envImpl, err := environment.NewEnvironmentImpl(bearerToken, userId)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError).SetInternal(err)
	}
	chart, err := loadChart(*chartsDir)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError).SetInternal(err)
	}
	_, err = envImpl.Create(chart, envInfo)
	if err != nil {
		if strings.Contains(err.Error(), "cannot re-use a name that is still in use") {
			return response.BadRequestWithMessagef(ctx, "环境%s已经存在，请修改名称后重试", name)
		}
		return echo.NewHTTPError(http.StatusInternalServerError).SetInternal(err)
	}
	return response.StatusOKNoContent(ctx)
}

func loadChart(chartDir string) (*chart.Chart, error) {
	chart, err := loader.LoadDir(chartDir)
	if err != nil {
		return nil, errors.Wrap(err, "loadDir error: "+utils.GetRuntimeLocation())
	}
	return chart, nil
}

func createEnvironmentInfo(requestBody *apigen.CreateEnvironmentJSONRequestBody, envName string, userId string) (*environment.EnvironmentInfo, error) {

	volumeMounts := models.VolumeMounts{}
	if *requestBody.Mounts != nil {
		for _, vm := range *requestBody.Mounts {
			volumeMounts = append(volumeMounts, models.MountInfo{
				MountPath: *vm.Mountpath,
				SubPath:   *vm.Subpath,
				Name:      "user-storage",
			})
		}
	}

	imageRepository, err := RepoToURL(*requestBody.Image, userId)
	if err != nil {
		requestBodyError := new(RequestBodyError)
		requestBodyError.keyErrors = make(map[string]string)
		requestBodyError.keyErrors["image"] = "无效输入"
		return nil, requestBodyError
	}

	pineappleInfo, err := helm.NewPineappleInfo(envName, userId, environment.EnvPrefix)
	if err != nil {
		return nil, err
	}

	envInfo := environment.EnvironmentInfo{
		PineappleInfo: pineappleInfo,
		Image: models.ImageInfo{
			Repository: imageRepository,
			Tag:        *requestBody.Image.Tag,
			PullPolicy: "Always",
		},
		ServerType: environment.ServerType{
			Jupyter: strings.Title(strconv.FormatBool(*requestBody.Jupyter.Enable)),
			Ssh:     strings.Title(strconv.FormatBool(*requestBody.Ssh.Enable)),
		},
		SshKey:       *requestBody.Ssh.IdRsaPub,
		JupyterToken: *requestBody.Jupyter.Token,
		PvcClaimName: "remote-storage",
		VolumeMounts: volumeMounts,
		ResourceId:   string(*requestBody.ComputeUnit),
	}

	values, err := envInfo.CreateEnvValues()
	if err != nil {
		return nil, err
	}
	envInfo.SetValues(values)

	return &envInfo, nil
}

func checkCreateEnvironmentJSONRequestBody(requestBody *apigen.CreateEnvironmentJSONRequestBody) *RequestBodyError {
	// TODO(fuhao): check RequestBody
	requestBodyError := new(RequestBodyError)
	requestBodyError.keyErrors = make(map[string]string)
	if requestBody == nil {
		return requestBodyError
	}
	if requestBody.Image == nil {
		requestBodyError.keyErrors["image"] = "无效输入"
	} else {
		if requestBody.Image.Repo == nil {
			requestBodyError.keyErrors["image.repository"] = "无效输入"
		}
		if requestBody.Image.Tag == nil {
			requestBodyError.keyErrors["image.tag"] = "无效输入"
		}
		if requestBody.Image.Source == nil {
			requestBodyError.keyErrors["image.source"] = "无效输入"
		}
	}
	if requestBody.ComputeUnit == nil {
		requestBodyError.keyErrors["compute_unit"] = "无效输入"
	} else {
		// TODO(fuhao): compute unit
		//existUnit := false
		//if !existUnit {
		//	requestBodyError.keyErrors["compute_unit"] = "无效输入"
		//}
	}
	if requestBody.Mounts == nil {
		requestBodyError.keyErrors["mounts"] = "无效输入"
	} else {
	}
	if requestBody.Jupyter == nil {
		requestBodyError.keyErrors["jupyter"] = "无效输入"
	} else {
		if requestBody.Jupyter.Enable == nil {
			requestBodyError.keyErrors["jupyter.enable"] = "无效输入"
		}
		if requestBody.Jupyter.Token == nil {
			*requestBody.Jupyter.Token = ""
		}
	}
	if requestBody.Ssh == nil {
		requestBodyError.keyErrors["ssh"] = "无效输入"
	} else {
		if requestBody.Ssh.Enable == nil {
			requestBodyError.keyErrors["ssh.enable"] = "无效输入"
		}
		if requestBody.Ssh.IdRsaPub == nil {
			requestBodyError.keyErrors["ssh.id_rsa.pub"] = "无效输入"
		}
		if requestBody.Ssh.Enable != nil && requestBody.Ssh.IdRsaPub != nil {
			if *requestBody.Ssh.Enable {
				if *requestBody.Ssh.IdRsaPub == "" {
					requestBodyError.keyErrors["ssh.id_rsa.pub"] = "无效输入"
				}
			}
		}
	}

	if len(requestBodyError.keyErrors) != 0 {
		return requestBodyError
	}

	// check Mounts
	for i, vm := range *requestBody.Mounts {
		if vm.Subpath == nil || vm.Mountpath == nil {
			//return errors.New("vm.SubPath == nil || vm.MountPath == nil")
			if vm.Subpath == nil {
				requestBodyError.keyErrors["mounts.subpath"] = "无效输入"
				return requestBodyError
			}
			requestBodyError.keyErrors["mounts.mountpath"] = "无效输入"
			return requestBodyError
		}
		if filepath.IsAbs(*vm.Subpath) {
			tempPath, err := filepath.Rel("/", *vm.Subpath)
			if err != nil {
				requestBodyError.keyErrors["mounts.subpath"] = "无效输入"
				return requestBodyError
			}
			*((*requestBody.Mounts)[i].Subpath) = tempPath
		}
	}
	return nil

}

func (handler *Handler) DeleteEnvironment(ctx echo.Context, name apigen.EnvironmentName) error {
	userId := ctx.Get("userID").(string)
	bearerToken := ctx.Get("bearerToken").(string)
	envImpl, err := environment.NewEnvironmentImpl(bearerToken, userId)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError).SetInternal(err)
	}
	if err := envImpl.Delete(string(name)); err != nil {
		return response.BadRequestWithMessagef(ctx, "环境: %s 删除失败，请重试", name)
	}
	return response.StatusOKNoContent(ctx)
}

func (handler *Handler) GetEnvironment(ctx echo.Context, name apigen.EnvironmentName) error {
	userId := ctx.Get("userID").(string)
	bearerToken := ctx.Get("bearerToken").(string)
	envImpl, err := environment.NewEnvironmentImpl(bearerToken, userId)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError).SetInternal(err)
	}
	envInfo, err := envImpl.GetInfo(string(name))
	if err != nil {
		return response.BadRequestWithMessagef(ctx, "环境: %s 获取失败，请重试", name)
	}

	envInfoResponse, err := parseEnvironmentRuntimeInfoToResponse(envInfo, userId)
	if err != nil {
		return response.BadRequestWithMessagef(ctx, "环境: %s 获取失败，请重试", name)
	}
	return ctx.JSON(http.StatusOK, envInfoResponse)
}

func (handler *Handler) GetEnvironmentList(ctx echo.Context, params apigen.GetEnvironmentListParams) error {
	userId := ctx.Get("userID").(string)
	bearerToken := ctx.Get("bearerToken").(string)
	envImpl, err := environment.NewEnvironmentImpl(bearerToken, userId)
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
	envInfos, err := envImpl.GetInfoList(limit, offset)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError).SetInternal(err)
	}

	var items []apigen.EnvironmentRuntimeInfo
	for _, e := range *envInfos.Item {
		item, err := parseEnvironmentRuntimeInfoToResponse(&e, userId)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError).SetInternal(err)
		}
		items = append(items, *item)
	}

	envInfosResponse := apigen.EnvironmentRuntimeInfos{
		Item:  &items,
		Total: envInfos.Total,
	}

	return ctx.JSON(http.StatusOK, envInfosResponse)
}

func parseEnvironmentRuntimeInfoToResponse(envInfo *environment.EnvironmentReleaseInfo, userId string) (*apigen.EnvironmentRuntimeInfo, error) {
	envName := apigen.EnvironmentName(*envInfo.StaticInfo.Name)
	envState := apigen.EnvironmentState(*envInfo.State)
	computeUnit := apigen.ComputeUnitId(*envInfo.StaticInfo.EnvironmentConfig.ComputeUnit)
	podName := envInfo.PodName
	image, err := URLToRepo(*envInfo.StaticInfo.EnvironmentConfig.Image.Repository, userId, *envInfo.StaticInfo.EnvironmentConfig.Image.Tag)
	if err != nil {
		return nil, errors.Wrap(err, "Url to Repo error: "+utils.GetRuntimeLocation())
	}
	sshInfo := new(apigen.EnvironmentRuntimeSshInfo)
	if envInfo.SshInfo != nil {
		sshInfo.SshIp = envInfo.SshInfo.SshIp
		sshInfo.SshPort = envInfo.SshInfo.SshPort
	}
	envInfoResponse := apigen.EnvironmentRuntimeInfo{
		SshInfo: sshInfo,
		State:   &envState,
		PodName: &podName,
		StaticInfo: &apigen.EnvironmentRuntimeStaticInfo{
			CreateTm:    &envInfo.StaticInfo.CreateTm.Time,
			Description: envInfo.StaticInfo.Description,
			EnvironmentConfig: &apigen.EnvironmentConfig{
				ComputeUnit: &computeUnit,
				Image:       image,
				Jupyter:     envInfo.StaticInfo.EnvironmentConfig.Jupyter,
				Mounts:      envInfo.StaticInfo.EnvironmentConfig.Mounts,
				Ssh:         envInfo.StaticInfo.EnvironmentConfig.Ssh,
			},
			Name:        &envName,
			NotebookUrl: envInfo.StaticInfo.NotebookUrl,
		},
		Events: envInfo.Events,
	}
	return &envInfoResponse, nil
}
