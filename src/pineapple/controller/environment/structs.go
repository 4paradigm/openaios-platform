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

package environment

import (
	"github.com/pkg/errors"
	"github.com/4paradigm/openaios-platform/src/pineapple/apigen"
	"github.com/4paradigm/openaios-platform/src/pineapple/utils"
	"helm.sh/helm/v3/pkg/time"
)

const (
	EnvPrefix = "env-"
)

type ServerType struct {
	Jupyter string `json:"jupyter" structs:"jupyter"`
	Ssh     string `json:"ssh" structs:"ssh"`
}

type EnvironmentState string

const (
	EnvironmentState_Failed    EnvironmentState = "Failed"
	EnvironmentState_Pending   EnvironmentState = "Pending"
	EnvironmentState_Running   EnvironmentState = "Running"
	EnvironmentState_Succeeded EnvironmentState = "Succeeded"
	EnvironmentState_Unknown   EnvironmentState = "Unknown"
	EnvironmentState_Killed    EnvironmentState = "Killed"
)

type EnvironmentRuntimeSshInfo struct {
	SshIp   *string `json:"ssh_ip,omitempty"`
	SshPort *string `json:"ssh_port,omitempty"`
}

type EnvironmentConfig struct {
	ComputeUnit *string `json:"compute_unit,omitempty"`
	Image       *struct {
		Repository *string `json:"repository,omitempty"`
		Tag        *string `json:"tag,omitempty"`
	} `json:"image,omitempty"`
	Jupyter *struct {
		Enable *bool   `json:"enable,omitempty"`
		Token  *string `json:"token,omitempty"`
	} `json:"jupyter,omitempty"`
	Mounts *[]apigen.StorageMapping `json:"mounts,omitempty"`
	Ssh    *struct {
		Enable   *bool   `json:"enable,omitempty"`
		IdRsaPub *string `json:"id_rsa.pub,omitempty"`
	} `json:"ssh,omitempty"`
}

type EnvironmentRuntimeStaticInfo struct {
	Name              *string            `json:"name,omitempty"`
	CreateTm          *time.Time         `json:"create_tm,omitempty"`
	EnvironmentConfig *EnvironmentConfig `json:"environmentConfig,omitempty"`
	NotebookUrl       *string            `json:"notebook_url,omitempty"`
	Description       *string            `json:"description,omitempty"`
}

type EnvironmentReleaseInfo struct {
	SshInfo     *EnvironmentRuntimeSshInfo         `json:"sshInfo,omitempty"`
	State       *EnvironmentState                  `json:"state,omitempty"`
	StaticInfo  *EnvironmentRuntimeStaticInfo      `json:"staticInfo,omitempty"`
	PodName     string                             `json:"pod_name,omitempty"`
	ReleaseName string                             `json:"-"`
	Type        string                             `json:"-"`
	Events      *[]apigen.ApplicationInstanceEvent `json:"events,omitempty"`
}

type EnvironmentReleaseInfos struct {
	Item  *[]EnvironmentReleaseInfo `json:"item,omitempty"`
	Total *int                      `json:"total,omitempty"`
}

type EnvironmentPodInfo struct {
	PodName *string
	State   *EnvironmentState
	Events  *[]apigen.ApplicationInstanceEvent
}

func (e *EnvironmentReleaseInfo) GetReleaseName() (string, error) {
	if e.ReleaseName == "" {
		return "", errors.New("Release Name cannot be found in EnvironmentReleaseInfo: " + utils.GetRuntimeLocation())
	}
	return e.ReleaseName, nil
}

func (e *EnvironmentReleaseInfo) GetType() (string, error) {
	if e.Type == "" {
		return "", errors.New("Release Type cannot be found in EnvironmentReleaseInfo: " + utils.GetRuntimeLocation())
	}
	return e.Type, nil
}
