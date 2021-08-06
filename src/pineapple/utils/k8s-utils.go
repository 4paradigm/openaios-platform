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

package utils

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"io"
	v1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"path/filepath"
	"time"
)

type k8sDockerRegistrySecretData struct {
	Auths map[string]k8sDockerRegistrySecret `json:"auths"`
}

type k8sDockerRegistrySecret struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email,omitempty"`
	Auth     string `json:"auth"`
}

func GetKubernetesClient() (*kubernetes.Clientset, error) {
	var config *rest.Config
	var err error

	config, err = rest.InClusterConfig()
	if err == rest.ErrNotInCluster {
		var kubeconfig string
		if home := homedir.HomeDir(); home != "" {
			kubeconfig = filepath.Join(home, ".kube", "config")
		} else {
			return nil, errors.New("HOME env not found")
		}

		// use the current context in kubeconfig
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, errors.WithMessagef(err, "trying to load in cluster kubernetes config, but failed")
	}

	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func CreateNamespace(ctx context.Context, client *kubernetes.Clientset, ns string) error {
	namespaceSpec := &v1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: ns,
		},
		Status: v1.NamespaceStatus{
			Phase: v1.NamespaceActive,
		},
	}
	_, err := client.CoreV1().Namespaces().Create(ctx, namespaceSpec, metav1.CreateOptions{})
	if err != nil && !k8serrors.IsAlreadyExists(err) {
		return err
	}
	return nil
}

// TODO: check all networking policy
func CreateNetworkPolicy(ctx context.Context, client *kubernetes.Clientset, ns string) error {
	networkPolicySpec := &networkingv1.NetworkPolicy{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "isolated-network-policy",
			Namespace: ns,
		},
		Spec: networkingv1.NetworkPolicySpec{
			// match all pods in 'ns'
			PodSelector: metav1.LabelSelector{},
			PolicyTypes: []networkingv1.PolicyType{
				networkingv1.PolicyTypeIngress,
				networkingv1.PolicyTypeEgress,
			},
			Ingress: []networkingv1.NetworkPolicyIngressRule{
				{
					// default to match all ports
					Ports: []networkingv1.NetworkPolicyPort{},
					From: []networkingv1.NetworkPolicyPeer{
						{
							// match all pod in the 'ns'
							PodSelector: &metav1.LabelSelector{},
						},
						{
							// match ingress controller
						},
					},
				},
			},
			Egress: []networkingv1.NetworkPolicyEgressRule{
				{
					Ports: []networkingv1.NetworkPolicyPort{},
					To: []networkingv1.NetworkPolicyPeer{
						{
							PodSelector: &metav1.LabelSelector{},
						},
					},
				},
			},
		},
	}
	_, err := client.NetworkingV1().NetworkPolicies(ns).Create(ctx, networkPolicySpec, metav1.CreateOptions{})
	if err != nil {
		return errors.Wrap(err, "create network policy failed")
	}
	return nil
}

func CreateUserRoleBindingWithEdit(ctx context.Context, client *kubernetes.Clientset, ns, userID string) error {
	rb := rbacv1.RoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name: userID + "-edit",
		},
		Subjects: []rbacv1.Subject{
			{
				Kind:     "User",
				Name:     userID,
				APIGroup: "rbac.authorization.k8s.io",
			},
			{
				Kind:      "ServiceAccount",
				Name:      userID + "-svc-account",
				APIGroup:  "",
				Namespace: userID,
			},
		},
		RoleRef: rbacv1.RoleRef{
			Kind:     "ClusterRole",
			Name:     "edit",
			APIGroup: "rbac.authorization.k8s.io",
		},
	}
	_, err := client.RbacV1().RoleBindings(ns).Create(ctx, &rb, metav1.CreateOptions{})
	if err != nil && !k8serrors.IsAlreadyExists(err) {
		return err
	}

	svcAcc := v1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name: userID + "-svc-account",
		},
	}
	_, err = client.CoreV1().ServiceAccounts(ns).Create(ctx, &svcAcc, metav1.CreateOptions{})
	if err != nil && !k8serrors.IsAlreadyExists(err) {
		return err
	}
	return nil
}

func GetPodStatus(client *kubernetes.Clientset, labelSelector string, ns string) (v1.PodPhase, error) {
	pods, err := client.CoreV1().Pods(ns).List(context.TODO(), metav1.ListOptions{
		LabelSelector: labelSelector,
	})
	if err != nil {
		return v1.PodUnknown, errors.Wrap(err, "pod list error: "+GetRuntimeLocation())
	}
	if len(pods.Items) == 0 {
		return v1.PodUnknown, nil
	}
	return pods.Items[0].Status.Phase, nil
}

func GetPodList(client *kubernetes.Clientset, labelSelector string, ns string) (*[]v1.Pod, error) {
	pods, err := client.CoreV1().Pods(ns).List(context.TODO(), metav1.ListOptions{
		LabelSelector: labelSelector,
	})
	if err != nil {
		return nil, errors.Wrap(err, "pod list error: "+GetRuntimeLocation())
	}
	return &pods.Items, nil
}

func GetServiceList(client *kubernetes.Clientset, labelSelector string, ns string) (*[]v1.Service, error) {
	svcs, err := client.CoreV1().Services(ns).List(context.TODO(), metav1.ListOptions{
		LabelSelector: labelSelector,
	})
	if err != nil {
		return nil, errors.Wrap(err, "Service list error: "+GetRuntimeLocation())
	}
	return &svcs.Items, nil
}

func GetSpecifyInvolvedObjectEventList(client *kubernetes.Clientset, involvedObjectName string, ns string) (*[]v1.Event, error) {
	eventsInterface := client.CoreV1().Events(ns)
	selector := eventsInterface.GetFieldSelector(&involvedObjectName, &ns, nil, nil)
	options := metav1.ListOptions{FieldSelector: selector.String()}
	events, err := eventsInterface.List(context.TODO(), options)
	if err != nil {
		return nil, errors.Wrap(err, "get Event error: "+GetRuntimeLocation())
	}
	return &events.Items, nil
}

func GetConfigMapList(client *kubernetes.Clientset, ns string, labelSelector string) (*v1.ConfigMapList, error) {
	cms, err := client.CoreV1().ConfigMaps(ns).List(context.TODO(), metav1.ListOptions{
		LabelSelector: labelSelector,
	})
	if err != nil {
		return nil, errors.Wrap(err, "list configmap failed at "+GetRuntimeLocation())
	}
	return cms, nil
}

func GenerateConfigMap(
	name string,
	generateName string,
	labels map[string]string,
	annotations map[string]string,
	data map[string]string) *v1.ConfigMap {

	return &v1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:         name,
			GenerateName: generateName,
			Labels:       labels,
			Annotations:  annotations,
		},
		Data: data,
	}
}

func CreateConfigMap(
	ctx context.Context,
	client *kubernetes.Clientset,
	ns string,
	cm *v1.ConfigMap) (*v1.ConfigMap, error) {

	cm, err := client.CoreV1().ConfigMaps(ns).Create(ctx, cm, metav1.CreateOptions{})
	if err != nil {
		return nil, err
	}
	return cm, nil
}

func UpdateConfigMap(
	ctx context.Context,
	client *kubernetes.Clientset,
	ns string,
	cm *v1.ConfigMap) (*v1.ConfigMap, error) {

	cm, err := client.CoreV1().ConfigMaps(ns).Update(ctx, cm, metav1.UpdateOptions{})
	if err != nil {
		return nil, err
	}
	return cm, nil
}

func CreateK8sDockerRegistrySecret(
	client *kubernetes.Clientset, ns string, registry string, username string, password string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	authBase64 := genK8sDockerRegistrySecretData(registry,
		k8sDockerRegistrySecret{Username: username, Password: password})
	_, err := client.CoreV1().Secrets(ns).Create(
		ctx,
		&v1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "harbor-registry-secret",
				Namespace: ns},
			Data: map[string][]byte{".dockerconfigjson": authBase64},
			Type: v1.SecretTypeDockerConfigJson,
		},
		metav1.CreateOptions{})
	return err
}

func genK8sDockerRegistrySecretData(registry string, s k8sDockerRegistrySecret) []byte {
	s.Auth = base64.URLEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", s.Username, s.Password)))
	data := k8sDockerRegistrySecretData{Auths: map[string]k8sDockerRegistrySecret{registry: s}}
	authBytes, _ := json.Marshal(data)
	return authBytes
}

type mergePatch struct {
	MetaData struct {
		Labels struct {
			UserReady string `json:"user_ready"`
		} `json:"labels"`
	} `json:"metadata"`
}

func MarkUserAsReady(ctx context.Context, client *kubernetes.Clientset, userID string) error {
	payload := new(mergePatch)
	payload.MetaData.Labels.UserReady = "true"
	payloadBytes, _ := json.Marshal(payload)

	_, err := client.CoreV1().Namespaces().Patch(ctx, userID, types.MergePatchType, payloadBytes, metav1.PatchOptions{})
	if err != nil {
		return errors.Wrap(err, "patch namespace failed at "+GetRuntimeLocation())
	}
	return nil
}

func CheckUserReady(ctx context.Context, client *kubernetes.Clientset, userID string) (bool, error) {
	ns, err := client.CoreV1().Namespaces().Get(ctx, userID, metav1.GetOptions{})
	if err != nil {
		if k8serrors.IsNotFound(err) {
			return false, nil
		}
		return false, errors.Wrapf(err, "get namespace %v failed at %v", userID, GetRuntimeLocation())
	}
	_, exist := ns.ObjectMeta.Labels["user_ready"]
	return exist, nil
}

func GetContainerLog(ctx context.Context, client *kubernetes.Clientset, userID string,
	podName string, containerName *string, tailLines *int64) (io.ReadCloser, error) {
	podLogOptions := v1.PodLogOptions{Follow: true, TailLines: tailLines}
	if containerName != nil {
		podLogOptions.Container = *containerName
	}
	resp := client.CoreV1().Pods(userID).GetLogs(podName, &podLogOptions)
	stream, err := resp.Stream(ctx)
	if err != nil {
		return nil, errors.Wrap(err, GetRuntimeLocation())
	}
	return stream, nil
}
