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

// Package helm provides utils for helm.
package helm

import (
	"github.com/4paradigm/openaios-platform/src/pineapple/utils"
	"github.com/pkg/errors"
	v1 "k8s.io/api/core/v1"
)

func (h *HelmClientImpl) GetPodList(helmReleaseInfo IHelmReleaseInfo) (*[]v1.Pod, error) {
	releaseName, err := helmReleaseInfo.GetReleaseName()
	if err != nil {
		return nil, errors.WithMessage(err, "GetReleaseName error: ")
	}
	//appType, err := helmReleaseInfo.GetType()
	if err != nil {
		return nil, errors.WithMessage(err, "GetType error: ")
	}
	client, err := utils.GetKubernetesClient()
	if err != nil {
		return nil, errors.WithMessage(err, "GetKubernetesClient error: ")
	}

	labelSelector := "openaios.4paradigm.com/app=true" + "," + "app.kubernetes.io/instance=" + releaseName

	podList, err := utils.GetPodList(client, labelSelector, *h.Config.Namespace)
	if err != nil {
		return nil, errors.WithMessage(err, "GetPodList error: ")
	}
	return podList, nil
}

func (h *HelmClientImpl) GetServiceList(helmReleaseInfo IHelmReleaseInfo) (*[]v1.Service, error) {
	releaseName, err := helmReleaseInfo.GetReleaseName()
	if err != nil {
		return nil, errors.WithMessage(err, "GetReleaseName error: ")
	}
	//appType, err := helmReleaseInfo.GetType()
	if err != nil {
		return nil, errors.WithMessage(err, "GetType error: ")
	}
	client, err := utils.GetKubernetesClient()
	if err != nil {
		return nil, errors.WithMessage(err, "GetKubernetesClient error: ")
	}

	labelSelector := "openaios.4paradigm.com/app=true" + "," + "app.kubernetes.io/instance=" + releaseName

	svcList, err := utils.GetServiceList(client, labelSelector, *h.Config.Namespace)
	if err != nil {
		return nil, errors.WithMessage(err, "GetPodList error: ")
	}
	return svcList, nil
}

func (h *HelmClientImpl) GetSpecifyInvolvedObjectEventList(involvedObjectName string) (*[]v1.Event, error) {
	client, err := utils.GetKubernetesClient()
	if err != nil {
		return nil, errors.WithMessage(err, "GetKubernetesClient error: ")
	}
	eventList, err := utils.GetSpecifyInvolvedObjectEventList(client, involvedObjectName, *h.Config.Namespace)
	if err != nil {
		return nil, errors.WithMessage(err, "GetSpecifyInvolvedObjectEventList error: ")
	}
	return eventList, nil
}
