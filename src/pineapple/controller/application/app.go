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
	"fmt"
	"github.com/4paradigm/openaios-platform/src/pineapple/utils"
	"github.com/4paradigm/openaios-platform/src/pineapple/utils/helm"
	"github.com/pkg/errors"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/release"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/duration"
	"strconv"
	"strings"
	"time"
)

type ApplicationImpl struct {
	*helm.HelmClientImpl
}

func NewApplicationImpl(kubeToken string, namespace string) (*ApplicationImpl, error) {
	helmClientImpl, err := helm.NewImpl(kubeToken, namespace)
	if err != nil {
		return nil, errors.WithMessage(err, "new helmClientImpl error: ")
	}
	appImpl := &ApplicationImpl{
		HelmClientImpl: helmClientImpl,
	}
	return appImpl, nil
}

func (a *ApplicationImpl) Delete(name string) error {
	releaseName := AppPrefix + name
	if err := a.HelmClientImpl.Delete(releaseName); err != nil {
		return errors.WithMessage(err, "app delete error: ")
	}
	return nil
}

func (a *ApplicationImpl) DeleteWithKeepHistory(name string) error {
	releaseName := AppPrefix + name
	if err := a.HelmClientImpl.DeleteWithKeepHistory(releaseName); err != nil {
		return errors.WithMessage(err, "app delete with KeepHistory error: ")
	}
	return nil
}

func (a *ApplicationImpl) GetApplicationInstanceInfoList(limit int, offset int) (*ApplicationInstanceInfos, error) {
	releases, total, err := a.getReleaseList(limit, offset)
	if err != nil {
		return nil, errors.WithMessage(err, "getReleaseList error: ")
	}
	appInstanceInfos := new(ApplicationInstanceInfos)
	appInstanceInfos.Total = &total
	infos := make([]ApplicationInstanceInfo, len(releases))
	for i, _ := range releases {
		info, err := parseReleaseToInstanceInfo(releases[i])
		if err != nil {
			return nil, errors.WithMessage(err, "parseReleaseToInstanceInfo error, i="+strconv.Itoa(i)+": ")
		}
		infos[i] = *info
	}
	appInstanceInfos.Item = &infos
	return appInstanceInfos, nil
}

func (a *ApplicationImpl) GetApplicationInstanceInfo(appName string) (*ApplicationInstanceInfo, error) {
	releaseName := AppPrefix + appName
	client := action.NewGet(a.ActionConfig)
	result, err := client.Run(releaseName)
	if err != nil {
		if strings.Contains(err.Error(), "release: not found") {
			return nil, nil
		}
		return nil, errors.Wrap(err, "GetApplicationInstanceInfo run error: "+utils.GetRuntimeLocation())
	}
	appInstanceInfo, err := parseReleaseToInstanceInfo(result)
	if err != nil {
		return nil, errors.WithMessage(err, "parseReleaseToInstanceInfo error: ")
	}
	return appInstanceInfo, nil
}

func (a *ApplicationImpl) GetApplicationReleaseNotes(appName string) (*ApplicationNotes, error) {
	releaseName := AppPrefix + appName
	client := action.NewGet(a.ActionConfig)
	result, err := client.Run(releaseName)
	if err != nil {
		if strings.Contains(err.Error(), "release: not found") {
			return nil, nil
		}
		return nil, errors.Wrap(err, "GetApplicationReleaseNotes run error: "+utils.GetRuntimeLocation())
	}
	return &ApplicationNotes{Notes: result.Info.Notes}, nil
}

func (a *ApplicationImpl) GetApplicationPodsInfo(appName string) (*ApplicationPodInfo, error) {
	releaseName := AppPrefix + appName
	releaseType := "application"
	appPodInfo := ApplicationPodInfo{
		ApplicationReleaseInfo: &ApplicationReleaseInfo{
			ReleaseName: releaseName,
			Type:        releaseType,
		},
		Pods:  nil,
		Total: 0,
	}
	pods, err := a.GetPodList(appPodInfo)
	if err != nil {
		return nil, errors.WithMessage(err, "GetPodList error: ")
	}
	podInfos := make([]PodInfo, len(*pods))
	for i, _ := range *pods {
		podInfos[i].Name = (*pods)[i].Name
		podInfos[i].Status = (*pods)[i].Status.Phase
		podInfos[i].CreateTm = (*pods)[i].CreationTimestamp.Time

		// containers
		var containers = []ContainerInfo{}
		containersMap := make(map[string]*ContainerInfo)
		for _, c := range (*pods)[i].Spec.Containers {
			container := ContainerInfo{
				Name:  c.Name,
				Image: c.Image,
				State: "",
				Ports: nil,
			}
			containerPorts := make([]ContainerPort, len(c.Ports))
			for k, p := range c.Ports {
				containerPorts[k].ContainerPort = strconv.Itoa(int(p.ContainerPort))
				containerPorts[k].Protocol = string(p.Protocol)
			}
			container.Ports = containerPorts
			containersMap[c.Name] = &container
		}
		for _, c := range (*pods)[i].Status.ContainerStatuses {
			if c.State.Running != nil {
				containersMap[c.Name].State = "Running"
			} else if c.State.Terminated != nil {
				containersMap[c.Name].State = "Terminated"
			} else {
				containersMap[c.Name].State = "Waiting"
			}
		}
		for _, value := range containersMap {
			containers = append(containers, *value)
		}
		podInfos[i].Containers = containers

		// events
		events, err := a.GetSpecifyInvolvedObjectEventList((*pods)[i].Name)
		if err != nil {
			return nil, errors.WithMessage(err, "GetSpecifyInvolvedObjectEventList error: ")
		}
		var eventInfos = []EventInfo{}
		for _, e := range *events {
			eventInfos = append(eventInfos, parseEventToEventInfo(&e))
		}
		podInfos[i].Event = eventInfos
	}
	appPodInfo.Pods = podInfos
	appPodInfo.Total = len(podInfos)
	return &appPodInfo, nil
}

func parseEventToEventInfo(e *v1.Event) EventInfo {
	// Age
	var interval string
	if e.Count > 1 {
		interval = fmt.Sprintf("%s (x%d over %s)", translateTimestampSince(e.LastTimestamp), e.Count, translateTimestampSince(e.FirstTimestamp))
	} else {
		interval = translateTimestampSince(e.FirstTimestamp)
		if e.FirstTimestamp.IsZero() {
			interval = translateMicroTimestampSince(e.EventTime)
		}
	}
	// From
	source := e.Source.Component
	if source == "" {
		source = e.ReportingController
	}
	eventInfo := EventInfo{
		Age:     interval,
		From:    source,
		Message: strings.TrimSpace(e.Message),
		Reason:  e.Reason,
		Type:    e.Type,
	}
	return eventInfo
}

func translateMicroTimestampSince(timestamp metav1.MicroTime) string {
	if timestamp.IsZero() {
		return "<unknown>"
	}

	return duration.HumanDuration(time.Since(timestamp.Time))
}

func translateTimestampSince(timestamp metav1.Time) string {
	if timestamp.IsZero() {
		return "<unknown>"
	}

	return duration.HumanDuration(time.Since(timestamp.Time))
}

func (a *ApplicationImpl) GetApplicationServicesInfo(appName string) (*ApplicationServiceInfo, error) {
	releaseName := AppPrefix + appName
	releaseType := "application"
	appSvcInfo := ApplicationServiceInfo{
		ApplicationReleaseInfo: &ApplicationReleaseInfo{
			ReleaseName: releaseName,
			Type:        releaseType,
		},
		Services: nil,
		Total:    0,
	}
	svcs, err := a.GetServiceList(appSvcInfo)
	if err != nil {
		return nil, errors.WithMessage(err, "GetServiceList error: ")
	}
	svcInfos := make([]ServiceInfo, len(*svcs))
	for i, _ := range *svcs {
		svcInfos[i].Name = (*svcs)[i].Name
		svcInfos[i].Type = (*svcs)[i].Spec.Type
		svcInfos[i].ClusterIP = (*svcs)[i].Spec.ClusterIP
		svcInfos[i].ExternalIPs = (*svcs)[i].Spec.ExternalIPs
		ports := make([]ServicePort, len((*svcs)[i].Spec.Ports))
		for j, _ := range (*svcs)[i].Spec.Ports {
			ports[j].Name = (*svcs)[i].Spec.Ports[j].Name
			ports[j].Port = strconv.Itoa(int(((*svcs)[i].Spec.Ports[j].Port)))
			ports[j].NodePort = strconv.Itoa(int(((*svcs)[i].Spec.Ports[j].NodePort)))
			ports[j].Protocol = string((*svcs)[i].Spec.Ports[j].Protocol)
		}
		svcInfos[i].Ports = ports
	}
	appSvcInfo.Services = svcInfos
	appSvcInfo.Total = len(svcInfos)
	return &appSvcInfo, nil
}

func parseReleaseToInstanceInfo(r *release.Release) (*ApplicationInstanceInfo, error) {
	info := ApplicationInstanceInfo{
		AppName:      r.Name[len(AppPrefix):],
		ChartName:    r.Chart.Name(),
		ChartVersion: r.Chart.Metadata.Version,
		CreateTm:     r.Info.FirstDeployed.Time,
		Duration:     time.Now().Sub(r.Info.FirstDeployed.Time),
		Status:       string(r.Info.Status),
	}
	return &info, nil
}

func (a *ApplicationImpl) getReleaseList(limit int, offset int) ([]*release.Release, int, error) {
	client := action.NewList(a.ActionConfig)
	client.Filter = AppPrefix
	client.Uninstalled = true
	client.Deployed = true
	client.Failed = true
	client.SetStateMask()
	results, err := client.Run()
	if err != nil {
		return nil, 0, errors.Wrap(err, "GetReleaseList run error: "+utils.GetRuntimeLocation())
	}
	total := len(results)

	if offset >= 0 && limit >= 0 {
		if offset >= len(results) {
			return nil, total, nil
		}
		realLimit := len(results)
		if limit > 0 && limit < realLimit {
			realLimit = limit
		}
		last := offset + realLimit
		if l := len(results); l < last {
			last = l
		}
		results = results[offset:last]
	}

	return results, total, nil
}
