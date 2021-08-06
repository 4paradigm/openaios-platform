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
	"github.com/4paradigm/openaios-platform/src/pineapple/conf"
	"github.com/4paradigm/openaios-platform/src/pineapple/utils"
	"github.com/pkg/errors"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/release"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"log"
	"os"
	"strings"
)

type IHelmReleaseInfo interface {
	GetReleaseName() (string, error)
	GetType() (string, error)
}

type HelmClientImpl struct {
	Config       *genericclioptions.ConfigFlags
	ActionConfig *action.Configuration
}

func NewImpl(kubeToken string, namespace string) (*HelmClientImpl, error) {
	kubeAPIServer := conf.GetKubeAPIServer()
	kubeCaFile := conf.GetKubeCaFile()
	config := &genericclioptions.ConfigFlags{
		APIServer:   &kubeAPIServer,
		CAFile:      &kubeCaFile,
		Namespace:   &namespace,
		BearerToken: &kubeToken,
	}
	actionConfig := new(action.Configuration)
	if err := actionConfig.Init(config, namespace, os.Getenv("HELM_DRIVER"), log.Printf); err != nil {
		return nil, errors.Wrap(err, "actionConfig init error:"+utils.GetRuntimeLocation())
	}
	helmClientImpl := &HelmClientImpl{
		Config:       config,
		ActionConfig: actionConfig,
	}
	return helmClientImpl, nil
}

func (h *HelmClientImpl) Create(chart *chart.Chart, info IPineappleInfo) (*release.Release, error) {
	chartValues, err := info.CreateChartValues()
	if err != nil {
		return nil, errors.WithMessage(err, "convert Values to chartValues error: ")
	}
	client := action.NewInstall(h.ActionConfig)
	client.Namespace = info.GetUserID()
	client.ReleaseName = info.GetPrefix() + info.GetName()
	postRenderer := NewPostRendererImpl()
	postRenderer.WriteKustomzation(".", `
commonLabels:
  app.kubernetes.io/managed-by: Helm
  app.kubernetes.io/instance: `+client.ReleaseName+`
  openaios.4paradigm.com/app: "true"
resources:
- all.yaml
`)
	client.PostRenderer = postRenderer
	results, err := client.Run(chart, chartValues)
	// TODO(fuhao): Graceful handling this error
	if err != nil {
		if strings.Contains(err.Error(), "cannot re-use a name that is still in use") {
			return nil, err
		}
		log.Printf("========== Warning: helm install error: ==========\n%+v", err)
		return nil, nil
	}
	return results, nil
}

func (h *HelmClientImpl) Delete(releaseName string) error {
	client := action.NewUninstall(h.ActionConfig)
	_, err := client.Run(releaseName)
	if err != nil {
		return errors.Wrap(err, "Delete run error: "+utils.GetRuntimeLocation())
	}
	return nil
}

func (h *HelmClientImpl) DeleteWithKeepHistory(releaseName string) error {
	client := action.NewUninstall(h.ActionConfig)
	client.KeepHistory = true
	_, err := client.Run(releaseName)
	if err != nil {
		return errors.Wrap(err, "Delete run with KeepHistory error: "+utils.GetRuntimeLocation())
	}
	return nil
}

func (h *HelmClientImpl) DeleteListWithKeepHistory(releases []*release.Release, description string) error {
	client := action.NewUninstall(h.ActionConfig)
	client.KeepHistory = true
	client.Description = description
	for i := range releases {
		_, err := client.Run(releases[i].Name)
		if err != nil {
			errMsg := "Error happened when delete: " + releases[i].Name + ": "
			return errors.Wrap(err, errMsg+utils.GetRuntimeLocation())
		}
	}
	return nil
}

func (h *HelmClientImpl) List(filter string) ([]*release.Release, error) {
	client := action.NewList(h.ActionConfig)
	client.Filter = filter
	results, err := client.Run()
	if err != nil {
		return nil, errors.Wrap(err, "List run error: "+utils.GetRuntimeLocation())
	}
	return results, nil
}
