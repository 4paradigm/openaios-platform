/*
 * Copyright © 2021 peizhaoyou <peizhaoyou@4paradigm.com>
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

// Package webterminal provides controller for webterminal.
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

type GetTerminalAPIResponse struct {
	Success bool   `json:"success"`
	Token   string `json:"token"`
	Message string `json:"message"`
}

func GetWebterminal(terminalInfo Webterminal) (string, error) {
	reqBody, err := json.Marshal(terminalInfo)
	if err != nil {
		return "", errors.Wrap(err, "Marshal terminalInfo error: ")
	}
	gottyURL := "http://127.0.0.1:8080/api/get-terminal"
	req, err := http.NewRequest("POST", gottyURL, bytes.NewBuffer(reqBody))
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

	respBody := new(GetTerminalAPIResponse)
	if err := json.Unmarshal(body, respBody); err != nil {
		return "", errors.Wrap(err, "Unmarshal resp body error: ")
	}

	if !respBody.Success {
		return "", errors.New("get terminal url from gotty error: " + respBody.Message)
	}

	return "/terminal/?token=" + respBody.Token, nil
}
