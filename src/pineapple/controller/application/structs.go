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

package application

import (
	"github.com/pkg/errors"
	"github.com/4paradigm/openaios-platform/src/pineapple/utils"
	v1 "k8s.io/api/core/v1"
	"time"
)

const (
	AppPrefix = "app-"
)

type ApplicationInstanceInfo struct {
	AppName      string        `json:"instance_name,omitempty"`
	ChartName    string        `json:"chart_name,omitempty"`
	ChartVersion string        `json:"chart_version,omitempty"`
	CreateTm     time.Time     `json:"create_tm,omitempty"`
	Duration     time.Duration `json:"-"`
	Status       string        `json:"status,omitempty"`
}

type ApplicationInstanceInfos struct {
	Item  *[]ApplicationInstanceInfo `json:"item,omitempty"`
	Total *int                       `json:"total,omitempty"`
}

type ApplicationReleaseInfo struct {
	ReleaseName string `json:"-"`
	Type        string `json:"-"`
}

type ApplicationPodInfo struct {
	*ApplicationReleaseInfo
	Total int       `json:"total,omitempty"`
	Pods  []PodInfo `json:"item,omitempty"`
}

type PodInfo struct {
	Name       string          `json:"name,omitempty"`
	Status     v1.PodPhase     `json:"state,omitempty"`
	Containers []ContainerInfo `json:"containers,omitempty"`
	Event      []EventInfo     `json:"events,omitempty"`
	CreateTm   time.Time       `json:"create_tm,omitempty"`
}

type ContainerInfo struct {
	Name  string          `json:"name,omitempty"`
	Image string          `json:"image,omitempty"`
	State string          `json:"state,omitempty"`
	ports []ContainerPort `json:"ports,omitempty"`
}

type ContainerPort struct {
	ContainerPort string `json:"container_port,omitempty"`
	Protocol      string `json:"protocol,omitempty"`
}

type EventInfo struct {
	Age     string `json:"age,omitempty"`
	From    string `json:"from,omitempty"`
	Message string `json:"message,omitempty"`
	Reason  string `json:"reason,omitempty"`
	Type    string `json:"type,omitempty"`
}

type ApplicationServiceInfo struct {
	*ApplicationReleaseInfo
	Total    int           `json:"total,omitempty"`
	Services []ServiceInfo `json:"item,omitempty"`
}

type ServiceInfo struct {
	Name        string         `json:"name,omitempty"`
	ClusterIP   string         `json:"cluster_ip,omitempty"`
	ExternalIPs []string       `json:"external_ips,omitempty"`
	Type        v1.ServiceType `json:"type,omitempty"`
	Ports       []ServicePort  `json:"ports,omitempty"`
}

type ServicePort struct {
	Name     string `json:"name,omitempty"`
	Port     string `json:"port,omitempty"`
	NodePort string `json:"node_port,omitempty"`
	Protocol string `json:"protocol,omitempty"`
}

type ApplicationNotes struct {
	Notes string `json:"notes,omitempty"`
}

func (a *ApplicationReleaseInfo) GetReleaseName() (string, error) {
	if a.ReleaseName == "" {
		return "", errors.New("Release Name cannot be found in ApplicationReleaseInfo: " + utils.GetRuntimeLocation())
	}
	return a.ReleaseName, nil
}

func (a *ApplicationReleaseInfo) GetType() (string, error) {
	if a.Type == "" {
		return "", errors.New("Release Type cannot be found in ApplicationReleaseInfo: " + utils.GetRuntimeLocation())
	}
	return a.Type, nil
}
