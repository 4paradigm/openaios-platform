package environment

import (
	"context"
	"encoding/json"
	"flag"
	"github.com/4paradigm/openaios-platform/src/pineapple/conf"
	"github.com/4paradigm/openaios-platform/src/pineapple/utils"
	"github.com/4paradigm/openaios-platform/src/pineapple/utils/helm"
	"github.com/pkg/errors"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/release"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"os"
	"strconv"
	"strings"
)

var (
	envSshUrl = flag.String("env-sshurl", os.Getenv("PINEAPPLE_ENV_SSHURL"),
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
	envSshInfos, err := e.getSshInfoList()
	if err != nil {
		return nil, err
	}
	envReleaseInfos := make([]EnvironmentReleaseInfo, len(envRuntimeStaticInfos))

	for i := 0; i < len(envRuntimeStaticInfos); i++ {
		//envReleaseInfos[i] = new(EnvironmentReleaseInfo)
		envReleaseInfos[i].StaticInfo = envRuntimeStaticInfos[i]
		releaseName := envReleaseInfos[i].StaticInfo.Name
		if envPodInfos[*releaseName] == nil {
			state := EnvironmentState_Unknown
			if *envRuntimeStaticInfos[i].Description == "ran out of credit" {
				state = EnvironmentState_Killed
			}
			envReleaseInfos[i].State = &state
			envReleaseInfos[i].PodName = ""
		} else {
			envReleaseInfos[i].State = envPodInfos[*releaseName].State
			envReleaseInfos[i].PodName = *envPodInfos[*releaseName].PodName
		}
		envReleaseInfos[i].SshInfo = envSshInfos[*releaseName]
		envReleaseInfos[i].ReleaseName = *releaseName
		envReleaseInfos[i].Type = "environment"
		*envReleaseInfos[i].StaticInfo.Name = (*releaseName)[len(EnvPrefix):]
		//if envReleaseInfos[i].State == nil {
		//	envReleaseInfos[i].State = new(EnvironmentState)
		//	*envReleaseInfos[i].State = EnvironmentState_Unknown
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
		*envReleaseInfo.State = EnvironmentState_Killed
	} else {
		envReleaseInfo.PodName = (*pods)[0].Name
		state := EnvironmentState((*pods)[0].Status.Phase)
		envReleaseInfo.State = &state
	}

	envReleaseInfo.SshInfo, err = e.getSshInfo(releaseName)
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
	*envRuntimeStaticInfo.NotebookUrl = conf.GetExternalURL() + *envRuntimeStaticInfo.NotebookUrl
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
	for _, e := range envRuntimeStaticInfos {
		labelSelector += *e.Name + ","
	}
	labelSelector += ")"
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
		envPodInfos[pod.Labels["name"]].PodName = &name
		envPodInfos[pod.Labels["name"]].State = &state
	}
	return envPodInfos, nil
}

func (e *EnvironmentImpl) getSshInfoList() (map[string]*EnvironmentRuntimeSshInfo, error) {
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
	var envSshList = map[string]*EnvironmentRuntimeSshInfo{}
	envSshList = make(map[string]*EnvironmentRuntimeSshInfo)
	for _, svc := range services.Items {
		port := strconv.Itoa(int(svc.Spec.Ports[0].NodePort))
		envSshList[svc.Labels["name"]] = new(EnvironmentRuntimeSshInfo)
		*envSshList[svc.Labels["name"]] = EnvironmentRuntimeSshInfo{
			SshIp:   envSshUrl,
			SshPort: &port,
		}
	}
	return envSshList, nil
}

func (e *EnvironmentImpl) getSshInfo(releaseName string) (*EnvironmentRuntimeSshInfo, error) {
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
	envRuntimeSshInfo := EnvironmentRuntimeSshInfo{
		SshIp:   envSshUrl,
		SshPort: &port,
	}
	return &envRuntimeSshInfo, nil
}

//func (e *EnvironmentImpl) GetExistEnvNames() ([]string, error) {
//	return nil, nil
//}
