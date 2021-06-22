package utils

import (
	"context"
	"github.com/4paradigm/openaios-platform/src/pineapple/conf"
	"github.com/AimAlex/harbor-client/harborcli"
	"github.com/AimAlex/harbor-client/harborcli/artifact"
	"github.com/AimAlex/harbor-client/harborcli/chart_repository"
	"github.com/AimAlex/harbor-client/harborcli/products"
	"github.com/AimAlex/harbor-client/harborcli/project"
	"github.com/AimAlex/harbor-client/harborcli/repository"
	"github.com/AimAlex/harbor-client/models"
	"github.com/go-openapi/runtime"
	httptransport "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
	"github.com/pkg/errors"
	"math/rand"
	"strings"
	"time"
)

func GetHarborClient(host string, basepath string, schemes []string) *harborcli.HarborAPI {
	//url, basepath, schemes := conf.GetHarborAddress()
	return harborcli.New(httptransport.New(host, basepath, schemes), strfmt.Default)
}

func getHarborAdminAuth() runtime.ClientAuthInfoWriter {
	username, passwd := conf.GetHarborAdmin()
	return httptransport.BasicAuth(username, passwd)
}

func CreateHarborUser(client *harborcli.HarborAPI, userID string, passwd string, realName string, email string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	userInfo := models.User{Username: userID,
		Password: passwd,
		Realname: realName,
		Email:    email}
	postUserParams := products.PostUsersParams{User: &userInfo, Context: ctx}
	_, err := client.Products.PostUsers(&postUserParams, getHarborAdminAuth())
	if err != nil && !strings.Contains(err.Error(), "status 409") {
		return errors.Wrap(err, "initial user image repo failed "+GetRuntimeLocation())
	}
	return nil
}

func CheckHarborUserExist(client *harborcli.HarborAPI, userID string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var pageSize int32 = 0
	params := products.GetUsersSearchParams{Context: ctx, PageSize: &pageSize, Username: userID}
	resp, err := client.Products.GetUsersSearch(&params, getHarborAdminAuth())
	if err != nil {
		return false, errors.Wrap(err, "get harbor users failed "+GetRuntimeLocation())
	}
	if len(resp.Payload) == 1 {
		return true, nil
	} else {
		return false, nil
	}
}

func CreateHarborProject(client *harborcli.HarborAPI, name string, isPublic string, limit *int64) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	projectInfo := models.ProjectReq{ProjectName: name, RegistryID: nil, StorageLimit: limit,
		Metadata: &models.ProjectMetadata{Public: isPublic}}
	createProjectParams := project.CreateProjectParams{Context: ctx, Project: &projectInfo}
	_, err := client.Project.CreateProject(&createProjectParams, getHarborAdminAuth())
	if err != nil && !strings.Contains(err.Error(), "createProjectConflict") {
		return errors.Wrap(err, "create project failed. "+GetRuntimeLocation())
	}
	return nil
}

func AddUserToProject(client *harborcli.HarborAPI, userID string, projectID int32, role int64) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	member := models.ProjectMember{RoleID: role, MemberUser: &models.UserEntity{Username: userID}}
	params := products.PostProjectsProjectIDMembersParams{Context: ctx, ProjectID: int64(projectID), ProjectMember: &member}
	_, err := client.Products.PostProjectsProjectIDMembers(&params, getHarborAdminAuth())
	if err != nil && !strings.Contains(err.Error(), "postProjectsProjectIdMembersConflict ") {
		return errors.Wrap(err, "cannot add member to project "+GetRuntimeLocation())
	}
	return nil
}

func CheckProjectMemberExist(client *harborcli.HarborAPI, userID string, projectID int64) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	params := products.GetProjectsProjectIDMembersParams{Context: ctx, ProjectID: projectID, Entityname: &userID}
	resp, err := client.Products.GetProjectsProjectIDMembers(&params, getHarborAdminAuth())
	if err != nil {
		return false, errors.Wrap(err, "cannot get project members "+GetRuntimeLocation())
	}
	if len(resp.Payload) == 1 {
		return true, nil
	} else {
		return false, nil
	}
}

func GetProjectIDByName(client *harborcli.HarborAPI, projectName string) (int32, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	listProjectParams := project.ListProjectsParams{Context: ctx, Name: &projectName}
	resp, err := client.Project.ListProjects(&listProjectParams, getHarborAdminAuth())
	if err != nil {
		return 0, errors.Wrap(err, "cannot list project "+GetRuntimeLocation())
	}
	for _, item := range resp.Payload {
		if item.Name == projectName {
			return item.ProjectID, nil
		}
	}
	return 0, errors.New("Cannot get user project " + GetRuntimeLocation())
}

func GetReplicationPolicies(client *harborcli.HarborAPI) ([]*models.ReplicationPolicy, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	params := products.GetReplicationPoliciesParams{Context: ctx}
	resp, err := client.Products.GetReplicationPolicies(&params, getHarborAdminAuth())
	if err != nil {
		return nil, errors.Wrap(err, "Get replication policy failed. "+GetRuntimeLocation())
	}
	return resp.Payload, nil
}

func GetPolicyExecution(client *harborcli.HarborAPI, policyID int64) (*models.ReplicationExecution, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	params := products.GetReplicationExecutionsParams{PolicyID: &policyID, Context: ctx}
	resp, err := client.Products.GetReplicationExecutions(&params, getHarborAdminAuth())
	if err != nil {
		return nil, errors.Wrap(err, "cannot get execution "+GetRuntimeLocation())
	}
	if len(resp.Payload) == 0 {
		return nil, nil
	}
	return resp.Payload[0], nil
}

func CreateReplicationPolicy(client *harborcli.HarborAPI, namespace string, repo string, tag string,
	name string, registryID int64) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	repoFilter := models.ReplicationFilter{Type: "name", Value: repo}
	tagFilter := models.ReplicationFilter{Type: "tag", Value: tag}
	registry := models.Registry{ID: registryID}
	trigger := models.ReplicationTrigger{Type: "manual"}
	policy := models.ReplicationPolicy{
		Deletion:      false,
		DestNamespace: namespace,
		Enabled:       true,
		Filters:       []*models.ReplicationFilter{&repoFilter, &tagFilter},
		Override:      true,
		Name:          name,
		SrcRegistry:   &registry,
		Trigger:       &trigger,
		DestRegistry:  nil,
	}

	params := products.PostReplicationPoliciesParams{Context: ctx, Policy: &policy}
	_, err := client.Products.PostReplicationPolicies(&params, getHarborAdminAuth())
	if err != nil {
		return errors.Wrap(err, "cannot create replication policy. "+GetRuntimeLocation())
	}
	return nil
}

func GetReplicationPolicyIDByName(client *harborcli.HarborAPI, policyName string) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	for {
		params := products.GetReplicationPoliciesParams{Context: ctx, Name: &policyName}
		resp, err := client.Products.GetReplicationPolicies(&params, getHarborAdminAuth())
		if err != nil {
			return 0, errors.Wrap(err, "cannot get replication policy. "+GetRuntimeLocation())
		}
		if len(resp.Payload) != 0 {
			return resp.Payload[0].ID, nil
		}
	}
}

func RunReplicationExecution(client *harborcli.HarborAPI, policyID int64) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	params := products.PostReplicationExecutionsParams{Context: ctx,
		Execution: &models.ReplicationExecution{PolicyID: policyID}}
	_, err := client.Products.PostReplicationExecutions(&params, getHarborAdminAuth())
	if err != nil {
		return errors.Wrap(err, "cannot execute replication policy "+GetRuntimeLocation())
	}
	return nil
}

func GetReplicationPolicyByID(client *harborcli.HarborAPI, policyID int64) (*models.ReplicationPolicy, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	params := products.GetReplicationPoliciesIDParams{Context: ctx, ID: policyID}
	resp, err := client.Products.GetReplicationPoliciesID(&params, getHarborAdminAuth())
	if err != nil {
		return nil, errors.Wrap(err, "cannot get replication policy "+GetRuntimeLocation())
	}
	return resp.Payload, nil
}

func DeleteReplicationPolicy(client *harborcli.HarborAPI, policyID int64) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	params := products.DeleteReplicationPoliciesIDParams{Context: ctx, ID: policyID}
	_, err := client.Products.DeleteReplicationPoliciesID(&params, getHarborAdminAuth())
	if err != nil {
		return errors.Wrap(err, "cannot delete replication policy. "+GetRuntimeLocation())
	}
	return nil
}

func GetReplicationExecution(client *harborcli.HarborAPI, policyID int64) (*models.ReplicationExecution, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	params := products.GetReplicationExecutionsParams{Context: ctx, PolicyID: &policyID}
	resp, err := client.Products.GetReplicationExecutions(&params, getHarborAdminAuth())
	if err != nil {
		return nil, errors.Wrap(err, "cannot get replication execution "+GetRuntimeLocation())
	}
	return resp.Payload[0], nil
}

func StopReplicationExecution(client *harborcli.HarborAPI, executionID int64) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	params := products.PutReplicationExecutionsIDParams{Context: ctx, ID: executionID}
	_, err := client.Products.PutReplicationExecutionsID(&params, getHarborAdminAuth())
	if err != nil {
		return errors.Wrap(err, "cannot stop replication execution. "+GetRuntimeLocation())
	}
	return nil
}

func GetRegistries(client *harborcli.HarborAPI) ([]*models.Registry, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	params := products.GetRegistriesParams{Context: ctx}
	resp, err := client.Products.GetRegistries(&params, getHarborAdminAuth())
	if err != nil {
		return nil, errors.Wrap(err, "cannot get registries. "+GetRuntimeLocation())
	}
	return resp.Payload, nil
}

func GetProjectRepositories(client *harborcli.HarborAPI, projectName string, filter *string) ([]*models.Repository, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var pageSize int64 = 0
	params := repository.ListRepositoriesParams{Context: ctx, ProjectName: projectName, Q: filter, PageSize: &pageSize}
	resp, err := client.Repository.ListRepositories(&params, getHarborAdminAuth())
	if err != nil {
		return nil, errors.Wrap(err, "cannot list repositories. "+GetRuntimeLocation())
	}
	return resp.Payload, nil
}

func GetProjectRepositoryArtifacts(client *harborcli.HarborAPI, projectName string, repoName string) ([]*models.Artifact, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	withTag := true
	var pageSize int64 = 0
	params := artifact.ListArtifactsParams{ProjectName: projectName, RepositoryName: repoName,
		WithTag: &withTag, Context: ctx, PageSize: &pageSize}
	resp, err := client.Artifact.ListArtifacts(&params, getHarborAdminAuth())
	if err != nil {
		return nil, errors.Wrap(err, "cannot list artifacts. "+GetRuntimeLocation())
	}
	return resp.Payload, nil
}

func CopyArtifact(client *harborcli.HarborAPI, projectName string, destRepo string, src string) error {
	params := artifact.CopyArtifactParams{ProjectName: projectName, RepositoryName: destRepo,
		From: src, Context: context.Background()}
	_, err := client.Artifact.CopyArtifact(&params, getHarborAdminAuth())
	if err != nil {
		return errors.Wrap(err, "cannot copy image "+GetRuntimeLocation())
	}
	return nil
}

func DeleteArtifact(client *harborcli.HarborAPI, projectName string, repo string, tag string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	params := artifact.DeleteArtifactParams{Context: ctx, ProjectName: projectName,
		RepositoryName: repo, Reference: tag}
	_, err := client.Artifact.DeleteArtifact(&params, getHarborAdminAuth())
	if err != nil {
		return errors.Wrap(err, "cannot delete image. "+GetRuntimeLocation())
	}
	return nil
}

func GetRepository(client *harborcli.HarborAPI, projectName string, repo string) (*models.Repository, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	params := repository.GetRepositoryParams{Context: ctx, ProjectName: projectName, RepositoryName: repo}
	resp, err := client.Repository.GetRepository(&params, getHarborAdminAuth())
	if err != nil {
		return nil, errors.Wrap(err, "cannot get repository "+GetRuntimeLocation())
	}
	return resp.Payload, nil
}

func DeleteRepository(client *harborcli.HarborAPI, projectName string, repo string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	params := repository.DeleteRepositoryParams{Context: ctx, ProjectName: projectName, RepositoryName: repo}
	_, err := client.Repository.DeleteRepository(&params, getHarborAdminAuth())
	if err != nil {
		return errors.Wrap(err, "cannot delete repository. "+GetRuntimeLocation())
	}
	return nil
}

// get random string for policy name
func RandStr(strlen int) string {
	rand.Seed(time.Now().UnixNano())
	data := make([]byte, strlen)
	var num int
	for i := 0; i < strlen; i++ {
		num = rand.Intn(57) + 65
		for {
			if num > 90 && num < 97 {
				num = rand.Intn(57) + 65
			} else {
				break
			}
		}
		data[i] = byte(num)
	}
	return "-" + string(data)
}

func GetChartrepoRepoCharts(client *harborcli.HarborAPI, repo string) ([]*models.ChartInfoEntry, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	params := chart_repository.GetChartrepoRepoChartsParams{
		Repo:       repo,
		Context:    ctx,
		HTTPClient: nil,
	}
	resp, err := client.ChartRepository.GetChartrepoRepoCharts(&params, getHarborAdminAuth())
	if err != nil {
		return nil, errors.Wrap(err, "GetChartrepoRepoCharts error: "+GetRuntimeLocation())
	}
	return resp.Payload, nil
}

func GetChartrepoRepoChartsName(client *harborcli.HarborAPI, repo string, name string) (models.ChartVersions, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	params := chart_repository.GetChartrepoRepoChartsNameParams{
		Name:       name,
		Repo:       repo,
		Context:    ctx,
		HTTPClient: nil,
	}
	resp, err := client.ChartRepository.GetChartrepoRepoChartsName(&params, getHarborAdminAuth())
	if err != nil {
		return nil, errors.Wrap(err, "GetChartrepoRepoChartsName error: "+GetRuntimeLocation())
	}
	return resp.Payload, nil
}

func GetChartrepoRepoChartsNameVersion(client *harborcli.HarborAPI, repo string, name string, version string) (*models.ChartVersionDetails, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	params := chart_repository.GetChartrepoRepoChartsNameVersionParams{
		Name:       name,
		Repo:       repo,
		Version:    version,
		Context:    ctx,
		HTTPClient: nil,
	}
	resp, err := client.ChartRepository.GetChartrepoRepoChartsNameVersion(&params, getHarborAdminAuth())
	if err != nil {
		return nil, errors.Wrap(err, "GetChartrepoRepoChartsNameVersion error: "+GetRuntimeLocation())
	}
	return resp.Payload, nil
}

func PostChartrepoRepoCharts(client *harborcli.HarborAPI, repo string, chart runtime.NamedReadCloser) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	params := chart_repository.PostChartrepoRepoChartsParams{
		Chart:      chart,
		Prov:       nil,
		Repo:       repo,
		Context:    ctx,
		HTTPClient: nil,
	}
	_, err := client.ChartRepository.PostChartrepoRepoCharts(&params, getHarborAdminAuth())
	if err != nil {
		return errors.Wrap(err, "PostChartrepoRepoCharts error: "+GetRuntimeLocation())
	}
	return nil
}
