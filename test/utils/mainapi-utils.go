package utils

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/4paradigm/openaios-platform/test/openapi/main/apigen/restclient"
)

func getRestClient() (*restclient.ClientWithResponses, error) {
	endpointURL := GetConfig("../test-cicd.toml").Env.EndpointURL
	return restclient.NewClientWithResponses(endpointURL)
}

func reqAddAuth(ctx context.Context, req *http.Request) error {
	var key string = "Authorization"
	req.Header.Add("Authorization", ctx.Value(key).(string))
	return nil
}

func GetUserTasks(tokenString string) (*restclient.UserTasksInfo, error) {
	client, err := getRestClient()
	if err != nil {
		return nil, errors.New("cannot connect to api gateway.")
	}
	authValue := "Bearer " + tokenString
	var key string = "Authorization"
	ctx, cancel := context.WithTimeout(context.WithValue(context.Background(), key, authValue), 10*time.Second)
	defer cancel()
	resp, err := client.GetUserTasksWithResponse(ctx, reqAddAuth)
	if err != nil {
		return nil, errors.New("cannot get user tasks")
	}
	if resp.StatusCode() != http.StatusOK {
		return nil, errors.New("cannot get user tasks")
	}
	return resp.JSON200, err
}

func GetUserCostPerMinute(tokenString string) (float64, error) {
	tasksInfo, err := GetUserTasks(tokenString)
	if err != nil {
		return 0, errors.New("GetUserTasks failed")
	}
	var cost float64 = 0.0
	for _, taskInfo := range *tasksInfo.TaskList {
		cost += *taskInfo.Price * float64(*taskInfo.Number)
	}
	return cost, nil
}

func CreateBasicEnvironment(tokenString string, computeUnitId string, repo string, tag string, source string) (string, error) {
	client, err := getRestClient()
	if err != nil {
		return "", errors.New("cannot connect to api gateway")
	}
	authValue := "Bearer " + tokenString
	var key interface{} = "Authorization"
	ctx, cancel := context.WithTimeout(context.WithValue(context.Background(), key, authValue), 10*time.Second)
	defer cancel()
	randEnvName := GenRandEnvName(20, 9)
	var envConfig restclient.CreateEnvironmentJSONRequestBody
	jsonStr := fmt.Sprintf("{\"compute_unit\":\"%s\",\"image\":{\"repo\":\"%s\",\"tag\":\"%s\",\"source\":\"%s\"},\"jupyter\":{\"enable\":false, \"token\":\"%s\"},\"mounts\":[],\"ssh\":{\"enable\":false, \"id_rsa.pub\":\"%s\"}}", computeUnitId, repo, tag, source, "", "")
	json.Unmarshal([]byte(jsonStr), &envConfig)
	resp, err := client.CreateEnvironment(ctx, randEnvName, envConfig, reqAddAuth)
	if resp.StatusCode == 200 {
		return string(randEnvName), nil
	}
	return "", errors.New("cannot create environment")
}

func DeleteEnvironment(tokenString string, envNameString string) error {
	client, err := getRestClient()
	if err != nil {
		return errors.New("cannot connect to api gateway")
	}
	authValue := "Bearer " + tokenString
	var key interface{} = "Authorization"
	ctx, cancel := context.WithTimeout(context.WithValue(context.Background(), key, authValue), 10*time.Second)
	defer cancel()
	var envName restclient.EnvironmentName = restclient.EnvironmentName(envNameString)
	_, err = client.DeleteEnvironment(ctx, envName, reqAddAuth)
	if err == nil {
		return nil
	}
	return errors.New("cannot delete environment")
}
