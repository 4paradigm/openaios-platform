package webterminal

import (
	"bytes"
	"encoding/json"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
)

type Webterminal struct {
	Namespace     string `json:"namespace"`
	PodName       string `json:"pod_name"`
	ContainerName string `json:"container_name"`
	BearerToken   string `json:"bearer_token"`
}

type GetTerminalApiResponse struct {
	Success bool   `json:"success"`
	Token   string `json:"token"`
	Message string `json:"message"`
}

func GetWebterminal(terminalInfo Webterminal) (string, error) {
	reqBody, err := json.Marshal(terminalInfo)
	if err != nil {
		return "", errors.Wrap(err, "Marshal terminalInfo error: ")
	}
	gottyUrl := "http://127.0.0.1:8080/api/get-terminal"
	req, err := http.NewRequest("POST", gottyUrl, bytes.NewBuffer(reqBody))
	if err != nil {
		return "", errors.Wrap(err, "http.NewRequest error: ")
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", errors.Wrap(err, "send req to gotty error: ")
	}
	defer resp.Body.Close()

	if statusCode := resp.StatusCode; statusCode != 200 {
		return "", errors.New("get terminal url from gotty error")
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", errors.Wrap(err, "read resp body error: ")
	}

	respBody := new(GetTerminalApiResponse)
	if err := json.Unmarshal(body, respBody); err != nil {
		return "", errors.Wrap(err, "Unmarshal resp body error: ")
	}

	if !respBody.Success {
		return "", errors.New("get terminal url from gotty error: " + respBody.Message)
	}

	return "/terminal/?token=" + respBody.Token, nil
}
