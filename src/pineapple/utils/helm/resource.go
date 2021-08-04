package helm

import (
	"github.com/pkg/errors"
	"github.com/4paradigm/openaios-platform/src/pineapple/utils"
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
