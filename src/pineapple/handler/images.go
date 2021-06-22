package handler

import (
	"code.cloudfoundry.org/bytefmt"
	"github.com/4paradigm/openaios-platform/src/internal/response"
	"github.com/4paradigm/openaios-platform/src/pineapple/apigen"
	"github.com/4paradigm/openaios-platform/src/pineapple/conf"
	"github.com/4paradigm/openaios-platform/src/pineapple/utils"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/pkg/errors"
	"math"
	"net/http"
	"path/filepath"
	"strings"
	"time"
)

var SourcePublic = "public"
var SourcePrivate = "private"

func (handler *Handler) ListImportingImages(ctx echo.Context) error {
	userID := ctx.Get("userID").(string)
	host, basepath, schemes := conf.GetHarborAddress()
	client := utils.GetHarborClient(host, basepath, schemes)

	policyList, err := utils.GetReplicationPolicies(client)
	if err != nil {
		return echo.NewHTTPError(
			http.StatusInternalServerError, "Get importing image list failed.").SetInternal(
			errors.Wrap(err, utils.GetRuntimeLocation()))
	}

	importingList := []apigen.ImageImportingInfo{}
	for _, policy := range policyList {
		if (policy == nil) || !strings.HasPrefix(policy.DestNamespace, userID) {
			continue
		}
		registry := apigen.ImageRegistryInfo{Id: &policy.SrcRegistry.ID,
			Url: &policy.SrcRegistry.URL}
		importingId := apigen.ImageImportingId(policy.ID)
		var imageRepo apigen.ImageRepo
		var imageTag apigen.ImageTag
		for _, filter := range policy.Filters {
			if filter.Type == "name" {
				imageRepo = apigen.ImageRepo(filter.Value)
				continue
			}
			if filter.Type == "tag" {
				imageTag = apigen.ImageTag(filter.Value)
			}
		}

		// Get policy execution
		execution, err := utils.GetPolicyExecution(client, policy.ID)
		if err != nil {
			return echo.NewHTTPError(
				http.StatusInternalServerError, "Get importing image list failed.").SetInternal(
				errors.Wrap(err, utils.GetRuntimeLocation()))
		}
		if execution == nil {
			ctx.Logger().Warn("Replication has no execution")
			continue
		}
		var status string
		var startTime, endTime time.Time
		status = execution.Status
		if status == "Succeed" && strings.Contains(execution.StatusText, "no resources") {
			status = "NotFound"
		}
		startTime, _ = time.Parse(time.RFC3339, execution.StartTime)
		endTime, _ = time.Parse(time.RFC3339, execution.EndTime)
		// TODO: currently one policy contains only one execution

		importingInfo := apigen.ImageImportingInfo{
			ImportingId: &importingId,
			Registry:    &registry,
			Status:      &status,
			Repo:        &imageRepo,
			Tag:         &imageTag,
			StartTime:   &startTime,
			EndTime:     &endTime,
		}
		importingList = append(importingList, importingInfo)
	}
	return ctx.JSON(http.StatusOK, &importingList)
}

func (handler *Handler) PostImagesImporting(ctx echo.Context, params apigen.PostImagesImportingParams) error {
	userID := ctx.Get("userID").(string)
	host, basepath, schemes := conf.GetHarborAddress()
	client := utils.GetHarborClient(host, basepath, schemes)
	name := userID + utils.RandStr(10)
	namespace, _ := filepath.Split(string(params.Repo))
	err := utils.CreateReplicationPolicy(client, filepath.Join(userID, namespace), string(params.Repo),
		string(params.Tag), name, params.RegistryId)
	if err != nil {
		return echo.NewHTTPError(
			http.StatusInternalServerError, "Create importing failed.").SetInternal(
			errors.Wrap(err, utils.GetRuntimeLocation()))
	}

	policyID, err := utils.GetReplicationPolicyIDByName(client, name)
	if err != nil {
		return echo.NewHTTPError(
			http.StatusInternalServerError, "Create importing failed.").SetInternal(
			errors.Wrap(err, utils.GetRuntimeLocation()))
	}

	err = utils.RunReplicationExecution(client, policyID)
	if err != nil {
		return echo.NewHTTPError(
			http.StatusInternalServerError, "Create importing failed.").SetInternal(
			errors.Wrap(err, utils.GetRuntimeLocation()))
	}
	return response.StatusOKNoContent(ctx)
}

func (handler *Handler) DeleteImagesImporting(ctx echo.Context, params apigen.DeleteImagesImportingParams) error {
	userID := ctx.Get("userID").(string)
	host, basepath, schemes := conf.GetHarborAddress()
	client := utils.GetHarborClient(host, basepath, schemes)

	policyID := int64(params.ImportingId)
	policy, err := utils.GetReplicationPolicyByID(client, policyID)
	if err != nil {
		return echo.NewHTTPError(
			http.StatusInternalServerError, "Delete Importing failed.").SetInternal(
			errors.Wrap(err, utils.GetRuntimeLocation()))
	}
	if !strings.HasPrefix(policy.DestNamespace, userID) || !strings.HasPrefix(policy.Name, userID) {
		ctx.Logger().Error("user " + userID + " cannot delete the policy of " + policy.DestNamespace)
		return response.BadRequestWithMessage(ctx, "Permission denied")
	}

	err = utils.DeleteReplicationPolicy(client, policyID)
	if err != nil {
		if strings.Contains(err.Error(), "has running executions") {
			return response.BadRequestWithMessage(ctx, "Importing is running.")
		} else {
			return echo.NewHTTPError(
				http.StatusInternalServerError, "Delete Importing failed.").SetInternal(
				errors.Wrap(err, utils.GetRuntimeLocation()))
		}
	}
	return response.StatusOKNoContent(ctx)
}

func (handler *Handler) PutImagesImporting(ctx echo.Context, params apigen.PutImagesImportingParams) error {
	userID := ctx.Get("userID").(string)
	host, basepath, schemes := conf.GetHarborAddress()
	client := utils.GetHarborClient(host, basepath, schemes)
	policyID := int64(params.ImportingId)
	policy, err := utils.GetReplicationPolicyByID(client, policyID)
	if err != nil {
		return echo.NewHTTPError(
			http.StatusInternalServerError, "Stop importing failed.").SetInternal(
			errors.Wrap(err, utils.GetRuntimeLocation()))
	}
	if !strings.HasPrefix(policy.DestNamespace, userID) || !strings.HasPrefix(policy.Name, userID) {
		ctx.Logger().Error("user " + userID + " cannot delete the policy of " + policy.DestNamespace)
		return response.BadRequestWithMessage(ctx, "Permission denied")
	}

	execution, err := utils.GetReplicationExecution(client, policyID)
	if err != nil {
		return echo.NewHTTPError(
			http.StatusInternalServerError, "Stop importing failed.").SetInternal(
			errors.Wrap(err, utils.GetRuntimeLocation()))
	}

	if execution.InProgress != 0 {
		err = utils.StopReplicationExecution(client, execution.ID)
		if err != nil {
			return echo.NewHTTPError(
				http.StatusInternalServerError, "Stop importing failed.").SetInternal(
				errors.Wrap(err, utils.GetRuntimeLocation()))
		}
	}
	return response.StatusOKNoContent(ctx)
}

func (handler *Handler) GetImagesRegistry(ctx echo.Context) error {
	host, basepath, schemes := conf.GetHarborAddress()
	client := utils.GetHarborClient(host, basepath, schemes)
	registries, err := utils.GetRegistries(client)
	if err != nil {
		return echo.NewHTTPError(
			http.StatusInternalServerError, "Get registry list failed.").SetInternal(
			errors.Wrap(err, utils.GetRuntimeLocation()))
	}

	registryList := []apigen.ImageRegistryInfo{}
	for _, registry := range registries {
		if registry == nil {
			continue
		}
		registryList = append(registryList, apigen.ImageRegistryInfo{Url: &registry.URL, Id: &registry.ID})
	}
	return ctx.JSON(http.StatusOK, registryList)
}

func (handler *Handler) GetPublicImagesInfo(ctx echo.Context, params apigen.GetPublicImagesInfoParams) error {
	host, basepath, schemes := conf.GetHarborAddress()
	client := utils.GetHarborClient(host, basepath, schemes)
	var imageCount int64 = 0
	// check filter
	var q *string = nil
	if params.Filter != nil {
		Q := filepath.Join("~public", *params.Filter)
		Q = "name=" + Q
		q = &Q
	}

	repositories, err := utils.GetProjectRepositories(client, "public", q)
	if err != nil {
		return echo.NewHTTPError(
			http.StatusInternalServerError, "Get images info failed.").SetInternal(
			errors.Wrap(err, utils.GetRuntimeLocation()))
	}
	for _, repo := range repositories {
		if repo == nil {
			continue
		}
		imageCount = imageCount + repo.ArtifactCount
	}
	imagesInfo := apigen.ProjectUserInfo{ImageCount: &imageCount}
	return ctx.JSON(http.StatusOK, imagesInfo)
}

func (handler *Handler) GetPublicImages(ctx echo.Context, params apigen.GetPublicImagesParams) error {
	host, basepath, schemes := conf.GetHarborAddress()
	client := utils.GetHarborClient(host, basepath, schemes)
	var imageCount, startCount, endCount int64 = 0, 1, math.MaxInt64
	if params.Page != nil && params.PageSize != nil {
		startCount = *params.PageSize**params.Page - *params.PageSize + 1
		endCount = *params.PageSize * *params.Page
	}
	// check filter
	var q *string = nil
	if params.Filter != nil {
		Q := filepath.Join("~public", *params.Filter)
		Q = "name=" + Q
		q = &Q
	}

	repositories, err := utils.GetProjectRepositories(client, "public", q)
	if err != nil {
		return echo.NewHTTPError(
			http.StatusInternalServerError, "Get public images failed.").SetInternal(
			errors.Wrap(err, utils.GetRuntimeLocation()))
	}
	imagesList := []apigen.ImageInfo{}
	for _, repo := range repositories {
		if repo == nil {
			continue
		}
		if imageCount+repo.ArtifactCount < startCount {
			imageCount += repo.ArtifactCount
			continue
		}
		if imageCount >= endCount {
			break
		}

		// get artifact list
		repoName := strings.TrimPrefix(repo.Name, "public/")
		artifacts, err := utils.GetProjectRepositoryArtifacts(client, "public", repoName)
		if err != nil {
			return echo.NewHTTPError(
				http.StatusInternalServerError, "Get public images failed.").SetInternal(
				errors.Wrap(err, utils.GetRuntimeLocation()))
		}
		imageRepo := apigen.ImageRepo(repoName)
		for _, atf := range artifacts {
			if atf == nil {
				continue
			}
			imageCount = imageCount + 1
			if imageCount < startCount {
				continue
			}
			if imageCount > endCount {
				break
			}
			imageTags := []apigen.ImageTag{}
			for _, tag := range atf.Tags {
				imageTags = append(imageTags, apigen.ImageTag(tag.Name))
			}
			imageSize := bytefmt.ByteSize(uint64(atf.Size))
			importingTime, err := time.Parse(time.RFC3339, atf.PushTime.String())
			if err != nil {
				return echo.NewHTTPError(
					http.StatusInternalServerError, "Get public images failed.").SetInternal(
					errors.Wrap(err, utils.GetRuntimeLocation()))
			}
			artifactInfo := apigen.ImageInfo{
				Tags:          &imageTags,
				Repo:          &imageRepo,
				Size:          &imageSize,
				ImportingTime: &importingTime,
			}
			imagesList = append(imagesList, artifactInfo)
		}
	}
	return ctx.JSON(http.StatusOK, imagesList)
}

func (handler *Handler) GetImagesInfo(ctx echo.Context) error {
	userID := ctx.Get("userID").(string)
	host, basepath, schemes := conf.GetHarborAddress()
	client := utils.GetHarborClient(host, basepath, schemes)
	var imageCount int64 = 0

	repositories, err := utils.GetProjectRepositories(client, userID, nil)
	if err != nil {
		return echo.NewHTTPError(
			http.StatusInternalServerError, "Get images info failed.").SetInternal(
			errors.Wrap(err, utils.GetRuntimeLocation()))
	}
	for _, repo := range repositories {
		if repo == nil {
			continue
		}
		imageCount = imageCount + repo.ArtifactCount
	}
	imagesInfo := apigen.ProjectUserInfo{ImageCount: &imageCount}
	return ctx.JSON(http.StatusOK, &imagesInfo)
}

func (handler *Handler) GetImages(ctx echo.Context, params apigen.GetImagesParams) error {
	userID := ctx.Get("userID").(string)
	host, basepath, schemes := conf.GetHarborAddress()
	client := utils.GetHarborClient(host, basepath, schemes)
	var imageCount, startCount, endCount int64 = 0, 1, math.MaxInt64
	if params.Page != nil && params.PageSize != nil {
		startCount = *params.PageSize**params.Page - *params.PageSize + 1
		endCount = *params.PageSize * *params.Page
	}

	repositories, err := utils.GetProjectRepositories(client, userID, nil)
	if err != nil {
		return echo.NewHTTPError(
			http.StatusInternalServerError, "Get images failed.").SetInternal(
			errors.Wrap(err, utils.GetRuntimeLocation()))
	}
	imagesList := []apigen.ImageInfo{}
	for _, repo := range repositories {
		if repo == nil {
			continue
		}
		if imageCount+repo.ArtifactCount < startCount {
			imageCount += repo.ArtifactCount
			continue
		}
		if imageCount >= endCount {
			break
		}

		// get artifact list
		repoName := strings.TrimPrefix(repo.Name, userID+"/")
		artifacts, err := utils.GetProjectRepositoryArtifacts(client, userID, repoName)
		if err != nil {
			return echo.NewHTTPError(
				http.StatusInternalServerError, "Get images failed.").SetInternal(
				errors.Wrap(err, utils.GetRuntimeLocation()))
		}
		imageRepo := apigen.ImageRepo(repoName)
		for _, atf := range artifacts {
			if atf == nil {
				continue
			}
			imageCount = imageCount + 1
			if imageCount < startCount {
				continue
			}
			if imageCount > endCount {
				break
			}
			imageTags := []apigen.ImageTag{}
			for _, tag := range atf.Tags {
				imageTags = append(imageTags, apigen.ImageTag(tag.Name))
			}
			imageSize := bytefmt.ByteSize(uint64(atf.Size))
			importingTime, err := time.Parse(time.RFC3339, atf.PushTime.String())
			if err != nil {
				return echo.NewHTTPError(
					http.StatusInternalServerError, "Get images failed.").SetInternal(
					errors.Wrap(err, utils.GetRuntimeLocation()))
			}
			artifactInfo := apigen.ImageInfo{
				Tags:          &imageTags,
				Repo:          &imageRepo,
				Size:          &imageSize,
				ImportingTime: &importingTime,
			}
			imagesList = append(imagesList, artifactInfo)
		}
	}
	return ctx.JSON(http.StatusOK, imagesList)
}

func (handler *Handler) PutImages(ctx echo.Context, params apigen.PutImagesParams) error {
	userID := ctx.Get("userID").(string)
	host, basepath, schemes := conf.GetHarborAddress()
	client := utils.GetHarborClient(host, basepath, schemes)
	src := userID + "/" + string(params.SrcRepo) + ":" + string(params.Tag)
	err := utils.CopyArtifact(client, userID, string(params.DestRepo), src)
	if err != nil {
		return echo.NewHTTPError(
			http.StatusInternalServerError, "Copy image failed.").SetInternal(
			errors.Wrap(err, utils.GetRuntimeLocation()))
	}
	return response.StatusOKNoContent(ctx)
}

func (handler *Handler) DeleteImages(ctx echo.Context, params apigen.DeleteImagesParams) error {
	userID := ctx.Get("userID").(string)
	host, basepath, schemes := conf.GetHarborAddress()
	client := utils.GetHarborClient(host, basepath, schemes)
	repoName := string(params.Repo)
	tag := string(params.Tag)
	err := utils.DeleteArtifact(client, userID, repoName, tag)
	if err != nil {
		ctx.Logger().Error(err.Error())
		if strings.Contains(err.Error(), "deleteArtifactNotFound") {
			return response.BadRequestWithMessage(ctx, "Image not found")
		}
		return echo.NewHTTPError(
			http.StatusInternalServerError, "Delete image failed.").SetInternal(
			errors.Wrap(err, utils.GetRuntimeLocation()))
	}

	repo, err := utils.GetRepository(client, userID, repoName)
	if err != nil {
		ctx.Logger().Warn(errors.Wrap(err, utils.GetRuntimeLocation()))
	}

	if repo != nil && repo.ArtifactCount == 0 {
		err = utils.DeleteRepository(client, userID, repoName)
		if err != nil {
			ctx.Logger().Warn(errors.Wrap(err, utils.GetRuntimeLocation()))
		}
	}
	return response.StatusOKNoContent(ctx)
}

func CreateHarborUser(userID string) error {
	host, basepath, schemes := conf.GetHarborAddress()
	client := utils.GetHarborClient(host, basepath, schemes)
	password := "4Paradigm" + userID
	realName := userID
	email := userID + "@pineapple.com"
	err := utils.CreateHarborUser(client, userID, password, realName, email)
	if err != nil {
		log.Warn(err)
		maxTime := 10
		for i := 0; i <= maxTime; i++ {
			time.Sleep(1 * time.Second)
			exist, e := utils.CheckHarborUserExist(client, userID)
			if e == nil {
				if !exist {
					return errors.Wrap(e, "cannot create harbor user "+utils.GetRuntimeLocation())
				} else {
					break
				}
			}
			log.Warn(e)
			if i == maxTime {
				return errors.Wrap(e, "cannot create harbor user "+utils.GetRuntimeLocation())
			}
		}
	}

	// create user secret in k8s namespace
	// TODO: Maybe k8s can use robot account
	harborUrl, _, _ := conf.GetHarborAddress()
	kubeClient, err := utils.GetKubernetesClient()
	if err != nil {
		return errors.Wrap(err, "cannot get kubenetes client. "+utils.GetRuntimeLocation())
	}
	err = utils.CreateK8sDockerRegistrySecret(kubeClient, userID, harborUrl, userID, password)
	if err != nil && !strings.Contains(err.Error(), "already exists") {
		return errors.Wrap(err, "create user secret failed "+utils.GetRuntimeLocation())
	}

	limit, err := conf.GetHarborStorageLimit()
	if err != nil {
		return errors.Wrap(err, "get user harbor quota failed "+utils.GetRuntimeLocation())
	}

	err = utils.CreateHarborProject(client, userID, "false", limit)
	if err != nil {
		log.Warn(err)
	}

	userProjectID, err := utils.GetProjectIDByName(client, userID)
	if err != nil {
		return errors.Wrap(err, "get project id failed "+utils.GetRuntimeLocation())
	}
	err = utils.AddUserToProject(client, userID, userProjectID, 2)
	if err != nil {
		exist, e := utils.CheckProjectMemberExist(client, userID, int64(userProjectID))
		if e != nil {
			return errors.Wrap(e, "add user to project failed "+utils.GetRuntimeLocation())
		}
		if !exist {
			return errors.Wrap(err, "user not a project member "+utils.GetRuntimeLocation())
		}
	}

	publicProjectID, err := utils.GetProjectIDByName(client, "public")
	if err != nil {
		return errors.Wrap(err, "get public project id failed "+utils.GetRuntimeLocation())
	}
	err = utils.AddUserToProject(client, userID, publicProjectID, 5)
	if err != nil {
		exist, e := utils.CheckProjectMemberExist(client, userID, int64(publicProjectID))
		if e != nil {
			return errors.Wrap(e, "add user to public project failed "+utils.GetRuntimeLocation())
		}
		if !exist {
			return errors.Wrap(err, "user not public project member "+utils.GetRuntimeLocation())
		}
	}
	return nil
}

func RepoToURL(config apigen.ImageConfig, userId string) (string, error) {
	harborUrl, _, _ := conf.GetHarborAddress()
	if harborUrl == "" {
		return "", errors.New("harbor url is empty " + utils.GetRuntimeLocation())
	}
	if config.Source == nil || config.Repo == nil {
		return "", errors.New("ImageConfig has empty attribute " + utils.GetRuntimeLocation())
	}
	if *config.Source == SourcePublic {
		return filepath.Join(harborUrl, "public", *config.Repo), nil
	} else if *config.Source == SourcePrivate {
		return filepath.Join(harborUrl, userId, *config.Repo), nil
	} else {
		return "", errors.New("ImageConfig source is invalid " + utils.GetRuntimeLocation())
	}
}

func URLToRepo(url string, userId string, tag string) (*apigen.ImageConfig, error) {
	harborUrl, _, _ := conf.GetHarborAddress()
	if harborUrl == "" {
		return nil, errors.New("harbor url is empty " + utils.GetRuntimeLocation())
	}
	if strings.HasPrefix(url, harborUrl) {
		url = strings.TrimPrefix(url, harborUrl+"/")
	} else {
		return nil, errors.New("harbor registry is wrong " + utils.GetRuntimeLocation())
	}
	if strings.HasPrefix(url, SourcePublic) {
		url = strings.TrimPrefix(url, "public/")
		return &apigen.ImageConfig{Repo: &url, Source: &SourcePublic, Tag: &tag}, nil
	} else if strings.HasPrefix(url, userId) {
		url = strings.TrimPrefix(url, userId+"/")
		return &apigen.ImageConfig{Repo: &url, Source: &SourcePrivate, Tag: &tag}, nil
	} else {
		return nil, errors.New("project name is wrong " + utils.GetRuntimeLocation())
	}
}
