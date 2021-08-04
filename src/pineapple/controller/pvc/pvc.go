package pvc

import (
	"github.com/fatih/structs"
	"github.com/pkg/errors"
	"github.com/4paradigm/openaios-platform/src/pineapple/utils"
	"github.com/4paradigm/openaios-platform/src/pineapple/utils/helm"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/release"
)

type SecretRef struct {
	Name string `json:"name" structs:"name"`
}

//type CephfsInfo struct {
//	Monitors  []string   `json:"monitors" structs:"monitors"`
//	Path      string     `json:"path" structs:"path"`
//	User      string     `json:"user" structs:"user"`
//	SecretRef *SecretRef `json:"secretRef" structs:"secretRef"`
//}

type Capacity struct {
	Storage string `structs:"storage"`
}

type CephSecret struct {
	Key string `structs:"key"`
}

type PvcInfo struct {
	UserId     string                  `structs:"userId"`
	Cephfs     *map[string]interface{} `structs:"cephfs"`
	Capacity   *Capacity               `structs:"capacity"`
	CephSecret *CephSecret             `structs:"cephSecret"`
}

type PvcImpl struct {
	*helm.HelmClientImpl
}

func NewPvcImpl(kubeToken string, namespace string) (*PvcImpl, error) {
	helmClientImpl, err := helm.NewImpl(kubeToken, namespace)
	if err != nil {
		return nil, errors.WithMessage(err, "new helmClientImpl error: ")
	}
	pvcImpl := &PvcImpl{
		HelmClientImpl: helmClientImpl,
	}
	return pvcImpl, nil
}

func (p *PvcImpl) Create(chartDir string, releaseName string, pvcInfo *PvcInfo) (*release.Release, error) {
	chartValues, err := p.parsePvcInfoToChartValues(pvcInfo)
	if err != nil {
		return nil, err
	}
	envChart, err := loader.LoadDir(chartDir)
	if err != nil {
		return nil, errors.Wrap(err, "loadDir error: "+utils.GetRuntimeLocation())
	}
	client := action.NewInstall(p.ActionConfig)
	client.Namespace = *p.Config.Namespace
	client.ReleaseName = releaseName
	results, err := client.Run(envChart, chartValues)
	if err != nil {
		return nil, errors.Wrap(err, "Create run error: "+utils.GetRuntimeLocation())
	}
	return results, nil
}

func (p *PvcImpl) Delete(releaseName string) error {
	client := action.NewUninstall(p.ActionConfig)
	_, err := client.Run(releaseName)
	if err != nil {
		return errors.Wrap(err, "Delete run error: "+utils.GetRuntimeLocation())

	}
	return nil
}

func (p *PvcImpl) parsePvcInfoToChartValues(pvcInfo *PvcInfo) (map[string]interface{}, error) {
	chartValues := structs.Map(pvcInfo)
	return chartValues, nil
}
