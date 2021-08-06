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

// Package environment provides controller for environment.
package environment

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/4paradigm/openaios-platform/src/pineapple/apigen"
	"github.com/4paradigm/openaios-platform/src/pineapple/conf"
	"github.com/4paradigm/openaios-platform/src/pineapple/utils"
	"github.com/4paradigm/openaios-platform/src/pineapple/utils/helm"
	"github.com/pkg/errors"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/release"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/duration"
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	envSSHURL = flag.String("env-sshurl", os.Getenv("PINEAPPLE_ENV_SSHURL"),
		"env-sshurl")
)

type EnvironmentImpl struct {
	*helm.HelmClientImpl
}

func NewEnvironmentImpl(kubeToken string, namespace string) (*EnvironmentImpl, error) {
	helmClientImpl, err := helm.NewImpl(kubeToken, namespace)
	if err != nil {
		return nil, errors.WithMessage(err, "new helmClientImpl error: ")
	}
	envImpl := &EnvironmentImpl{
		HelmClientImpl: helmClientImpl,
	}
	return envImpl, nil
}

func (e *EnvironmentImpl) Delete(name string) error {
	releaseName := EnvPrefix + name
	if err := e.HelmClientImpl.Delete(releaseName); err != nil {
		return errors.WithMessage(err, "env delete error: ")
	}
	return nil
}

func (e *EnvironmentImpl) DeleteWithKeepHistory(name string) error {
	releaseName := EnvPrefix + name
	if err := e.HelmClientImpl.DeleteWithKeepHistory(releaseName); err != nil {
		return errors.WithMessage(err, "env delete with KeepHistory error: ")
	}
	return nil
}

func (e *EnvironmentImpl) GetInfoList(limit int, offset int) (*EnvironmentReleaseInfos, error) {
	envRuntimeStaticInfos, total, err := e.getStaticInfoList(limit, offset)
	if err != nil {
		return nil, err
	}
	envPodInfos, err := e.getPodInfoList(envRuntimeStaticInfos)
	if err != nil {
		return nil, err
	}
	envSSHInfos, err := e.getSSHInfoList()
	if err != nil {
		return nil, err
	}
	envReleaseInfos := make([]EnvironmentReleaseInfo, len(envRuntimeStaticInfos))

	for i := 0; i < len(envRuntimeStaticInfos); i++ {
		//envReleaseInfos[i] = new(EnvironmentReleaseInfo)
		envReleaseInfos[i].StaticInfo = envRuntimeStaticInfos[i]
		releaseName := envReleaseInfos[i].StaticInfo.Name
		if envPodInfos[*releaseName] == nil {
			state := EnvironmentStateUnknown
			if *envRuntimeStaticInfos[i].Description == "ran out of credit" {
				state = EnvironmentStateKilled
			}
			envReleaseInfos[i].State = &state
			envReleaseInfos[i].PodName = ""
			envReleaseInfos[i].Events = nil
		} else {
			envReleaseInfos[i].State = envPodInfos[*releaseName].State
			envReleaseInfos[i].PodName = *envPodInfos[*releaseName].PodName
			envReleaseInfos[i].Events = envPodInfos[*releaseName].Events
		}
		envReleaseInfos[i].SSHInfo = envSSHInfos[*releaseName]
		envReleaseInfos[i].ReleaseName = *releaseName
		envReleaseInfos[i].Type = "environment"
		*envReleaseInfos[i].StaticInfo.Name = (*releaseName)[len(EnvPrefix):]
		//if envReleaseInfos[i].State == nil {
		//	envReleaseInfos[i].State = new(EnvironmentState)
		//	*envReleaseInfos[i].State = EnvironmentStateUnknown
		//}
	}

	envInfos := EnvironmentReleaseInfos{
		Total: &total,
		Item:  &envReleaseInfos,
	}
	return &envInfos, nil
}

func (e *EnvironmentImpl) GetInfo(name string) (*EnvironmentReleaseInfo, error) {
	releaseName := EnvPrefix + name
	envReleaseInfo := new(EnvironmentReleaseInfo)
	envReleaseInfo.Type = "environment"
	envRuntimeStaticInfo, err := e.getStaticInfo(releaseName)
	if err != nil {
		return nil, err
	}
	envReleaseInfo.StaticInfo = envRuntimeStaticInfo
	envReleaseInfo.ReleaseName = *envReleaseInfo.StaticInfo.Name
	envReleaseInfo.StaticInfo.Name = &name

	pods, err := e.GetPodList(envReleaseInfo)
	if err != nil {
		return nil, errors.WithMessage(err, "GetPodList error: ")
	}
	if len(*pods) == 0 {
		envReleaseInfo.PodName = ""
		*envReleaseInfo.State = EnvironmentStateKilled
		envReleaseInfo.Events = nil
	} else {
		envReleaseInfo.PodName = (*pods)[0].Name
		state := EnvironmentState((*pods)[0].Status.Phase)
		envReleaseInfo.State = &state
		events, err := e.GetSpecifyInvolvedObjectEventList((*pods)[0].Name)
		if err != nil {
			return nil, errors.WithMessage(err, "GetSpecifyInvolvedObjectEventList error: ")
		}
		var eventInfos = []apigen.ApplicationInstanceEvent{}
		for _, e := range *events {
			eventInfos = append(eventInfos, parseEventToEventInfo(&e))
		}
		envReleaseInfo.Events = &eventInfos
	}

	envReleaseInfo.SSHInfo, err = e.getSSHInfo(releaseName)
	if err != nil {
		return envReleaseInfo, err
	}
	return envReleaseInfo, nil
}

func (e *EnvironmentImpl) getStaticInfoList(limit int, offset int) ([]*EnvironmentRuntimeStaticInfo, int, error) {
	client := action.NewList(e.ActionConfig)
	client.Filter = EnvPrefix
	client.Uninstalled = true
	client.Deployed = true
	client.Failed = true
	client.SetStateMask()
	results, err := client.Run()
	if err != nil {
		return nil, 0, errors.Wrap(err, "getStaticInfoList run error: "+utils.GetRuntimeLocation())
	}
	total := len(results)

	if offset >= 0 && limit >= 0 {
		// Guard on offset
		if offset >= len(results) {
			return nil, total, nil
		}

		// Calculate the limit and offset, and then truncate results if necessary.
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

	envRuntimeStaticInfos := []*EnvironmentRuntimeStaticInfo{}
	for _, r := range results {
		envRuntimeStaticInfo, err := e.getEnvironmentRuntimeStaticInfoFromRelease(r)
		if err != nil {
			return nil, total, errors.WithMessage(err, "getEnvironmentRuntimeStaticInfoFromRelease error: ")
		}
		envRuntimeStaticInfos = append(envRuntimeStaticInfos, envRuntimeStaticInfo)
	}
	return envRuntimeStaticInfos, total, nil
}

func (e *EnvironmentImpl) getStaticInfo(releaseName string) (*EnvironmentRuntimeStaticInfo, error) {
	client := action.NewGet(e.ActionConfig)
	results, err := client.Run(releaseName)
	if err != nil {
		return nil, errors.Wrap(err, "getStaticInfo run error: "+utils.GetRuntimeLocation())
	}
	envRuntimeStaticInfo, err := e.getEnvironmentRuntimeStaticInfoFromRelease(results)
	if err != nil {
		return nil, err
	}
	return envRuntimeStaticInfo, nil
}

func (e *EnvironmentImpl) getEnvironmentRuntimeStaticInfoFromRelease(release *release.Release) (*EnvironmentRuntimeStaticInfo, error) {
	manifest := release.Manifest
	indexStart := strings.Index(manifest, "PINEAPPLE_ENV_INFO_START_HERE<<<")
	indexEnd := strings.Index(manifest, "<<<PINEAPPLE_ENV_INFO_END_HERE")
	envRuntimeStaticInfoStr := manifest[indexStart+32 : indexEnd]
	envRuntimeStaticInfo := new(EnvironmentRuntimeStaticInfo)
	if err := json.Unmarshal([]byte(envRuntimeStaticInfoStr), envRuntimeStaticInfo); err != nil {
		return nil, errors.Wrap(err, "Unmarshal env Static Info error: "+utils.GetRuntimeLocation())

	}
	*envRuntimeStaticInfo.CreateTm = release.Info.FirstDeployed
	*envRuntimeStaticInfo.NotebookURL = conf.GetExternalURL() + *envRuntimeStaticInfo.NotebookURL
	envRuntimeStaticInfo.Description = new(string)
	*envRuntimeStaticInfo.Description = release.Info.Description
	return envRuntimeStaticInfo, nil
}

func (e *EnvironmentImpl) getPodInfoList(envRuntimeStaticInfos []*EnvironmentRuntimeStaticInfo) (map[string]*EnvironmentPodInfo, error) {
	client, err := utils.GetKubernetesClient()
	if err != nil {
		return nil, errors.WithMessage(err, "GetKubernetesClient error: ")
	}
	labelSelector := "app.kubernetes.io/instance in ("
	for _, eI := range envRuntimeStaticInfos {
		labelSelector += *eI.Name + ","
	}
	labelSelector += "),"
	labelSelector += "openaios.4paradigm.com/app in (true)"

	podList, err := utils.GetPodList(client, labelSelector, *e.Config.Namespace)
	if err != nil {
		return nil, errors.WithMessage(err, "Get pod list error: ")
	}
	var envPodInfos = map[string]*EnvironmentPodInfo{}
	envPodInfos = make(map[string]*EnvironmentPodInfo)
	for _, pod := range *podList {
		envPodInfos[pod.Labels["name"]] = new(EnvironmentPodInfo)
		state := EnvironmentState(pod.Status.Phase)
		name := pod.Name
		events, err := e.GetSpecifyInvolvedObjectEventList(pod.Name)
		if err != nil {
			return nil, errors.WithMessage(err, "GetSpecifyInvolvedObjectEventList error: ")
		}
		var eventInfos = []apigen.ApplicationInstanceEvent{}
		for _, e := range *events {
			eventInfos = append(eventInfos, parseEventToEventInfo(&e))
		}
		envPodInfos[pod.Labels["name"]].PodName = &name
		envPodInfos[pod.Labels["name"]].State = &state
		envPodInfos[pod.Labels["name"]].Events = &eventInfos
	}
	return envPodInfos, nil
}

func parseEventToEventInfo(e *v1.Event) apigen.ApplicationInstanceEvent {
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
	message := strings.TrimSpace(e.Message)
	eventInfo := apigen.ApplicationInstanceEvent{
		Age:     &interval,
		From:    &source,
		Message: &message,
		Reason:  &e.Reason,
		Type:    &e.Type,
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

func (e *EnvironmentImpl) getSSHInfoList() (map[string]*EnvironmentRuntimeSSHInfo, error) {
	client, err := utils.GetKubernetesClient()
	if err != nil {
		return nil, errors.WithMessage(err, "GetKubernetesClient error: ")
	}
	labelSelector := "role=ssh-service"
	services, err := client.CoreV1().Services(*e.Config.Namespace).List(context.TODO(), metav1.ListOptions{
		LabelSelector: labelSelector,
	})
	if err != nil {
		return nil, errors.Wrap(err, "ssh service list error: "+utils.GetRuntimeLocation())
	}
	var envSSHList = map[string]*EnvironmentRuntimeSSHInfo{}
	envSSHList = make(map[string]*EnvironmentRuntimeSSHInfo)
	for _, svc := range services.Items {
		port := strconv.Itoa(int(svc.Spec.Ports[0].NodePort))
		envSSHList[svc.Labels["name"]] = new(EnvironmentRuntimeSSHInfo)
		*envSSHList[svc.Labels["name"]] = EnvironmentRuntimeSSHInfo{
			SSHIP:   envSSHURL,
			SSHPort: &port,
		}
	}
	return envSSHList, nil
}

func (e *EnvironmentImpl) getSSHInfo(releaseName string) (*EnvironmentRuntimeSSHInfo, error) {
	client, err := utils.GetKubernetesClient()
	if err != nil {
		return nil, errors.WithMessage(err, "GetKubernetesClient error: ")
	}
	labelSelector := "role=ssh-service, name=" + releaseName
	services, err := client.CoreV1().Services(*e.Config.Namespace).List(context.TODO(), metav1.ListOptions{
		LabelSelector: labelSelector,
	})
	if err != nil {
		return nil, errors.Wrap(err, "ssh service list error: "+utils.GetRuntimeLocation())
	}
	if len(services.Items) == 0 {
		return nil, nil
	}
	port := strconv.Itoa(int(services.Items[0].Spec.Ports[0].NodePort))
	envRuntimeSSHInfo := EnvironmentRuntimeSSHInfo{
		SSHIP:   envSSHURL,
		SSHPort: &port,
	}
	return &envRuntimeSSHInfo, nil
}

//func (e *EnvironmentImpl) GetExistEnvNames() ([]string, error) {
//	return nil, nil
//}
