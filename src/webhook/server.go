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

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/labstack/gommon/log"
	"github.com/pkg/errors"
	"github.com/4paradigm/openaios-platform/src/internal/billingclient"
	"github.com/4paradigm/openaios-platform/src/internal/response"
	"k8s.io/api/admission/v1beta1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	tlsDir      = `/run/secrets/tls`
	tlsCertFile = `tls.crt`
	tlsKeyFile  = `tls.key`
)

var (
	podResource      = metav1.GroupVersionResource{Version: "v1", Resource: "pods"}
	billingServerURL = flag.String("billing-server-url", os.Getenv("PINEAPPLE_BILLING_SERVER_URL"),
		"billing server url")
)

// applySecurityDefaults implements the logic of our example admission controller webhook. For every pod that is created
// (outside of Kubernetes namespaces), it first checks if `runAsNonRoot` is set. If it is not, it is set to a default
// value of `false`. Furthermore, if `runAsUser` is not set (and `runAsNonRoot` was not initially set), it defaults
// `runAsUser` to a value of 1234.
//
// To demonstrate how requests can be rejected, this webhook further validates that the `runAsNonRoot` setting does
// not conflict with the `runAsUser` setting - i.e., if the former is set to `true`, the latter must not be `0`.
// Note that we combine both the setting of defaults and the check for potential conflicts in one webhook; ideally,
// the latter would be performed in a validating webhook admission controller.
func applySecurityDefaults(req *v1beta1.AdmissionRequest) ([]patchOperation, error) {
	// This handler should only get called on Pod objects as per the MutatingWebhookConfiguration in the YAML file.
	// However, if (for whatever reason) this gets invoked on an object of a different kind, issue a log message but
	// let the object request pass through otherwise.
	if req.Resource != podResource {
		log.Printf("expect resource to be %s", podResource)
		return nil, nil
	}

	// Parse the Pod object.
	raw := req.Object.Raw
	pod := corev1.Pod{}
	if _, _, err := universalDeserializer.Decode(raw, nil, &pod); err != nil {
		return nil, fmt.Errorf("could not deserialize pod object: %v", err)
	}

	// check if pod is application
	if pod.Labels["openaios.4paradigm.com/app"] != "true" {
		return nil, nil
	}
	// get user ID
	userID := req.Namespace
	// get computeunitMap in annotation
	var containerComputeunitMap = map[string]string{}
	var defaultComputeunit string = ""
	for k, v := range pod.Annotations {
		if strings.HasPrefix(k, "openaios.4paradigm.com/computeunit.") {
			containerName := strings.TrimPrefix(k, "openaios.4paradigm.com/computeunit.")
			containerComputeunitMap[containerName] = v
		} else if k == "openaios.4paradigm.com/computeunit" {
			defaultComputeunit = v
		}
	}
	// Create patch operations to apply sensible defaults, if those options are not set explicitly.
	var patches []patchOperation
	var computeunitList = []string{}
	var volumesPod = []corev1.Volume{}
	var volumesExists = map[string]bool{}
	billingClient, err := billingclient.GetBillingClient(*billingServerURL)
	if err != nil {
		log.Error(err)
		return nil, errors.Wrap(err, response.GetRuntimeLocation())
	}
	for item, container := range pod.Spec.Containers {
		computeunit, ok := containerComputeunitMap[container.Name]
		if !ok {
			if defaultComputeunit == "" {
				log.Error("container does not have computeunit")
				return nil, errors.New("container does not have computeunit " + response.GetRuntimeLocation())
			} else {
				computeunit = defaultComputeunit
			}
		} else {
			delete(containerComputeunitMap, container.Name)
		}

		itemString := strconv.FormatInt(int64(item), 10)
		computeunitList = append(computeunitList, computeunit)

		// get computeunit from billing server
		computeunitInfo, err := billingclient.GetOneComputeUnit(billingClient, userID, computeunit)
		if err != nil {
			log.Error(err)
			return nil, errors.Wrap(err, response.GetRuntimeLocation())
		}
		spec := computeunitInfo.Spec
		if spec == nil {
			log.Error("computeunit has no spec")
			return nil, errors.New("computeunit has no spec " + response.GetRuntimeLocation())
		}
		specMap := *spec

		// deal with resources
		patches = append(patches, patchOperation{
			Op:    "add",
			Path:  "/spec/containers/" + itemString + "/resources",
			Value: specMap["resources"],
		})

		// deal with volumeMounts
		if volumeMounts, ok := specMap["volumeMounts"]; ok {
			specVMList, suc := volumeMounts.([]interface{})
			if !suc {
				log.Error("invalid volumeMounts")
				return nil, errors.New("invalid volumeMounts " + response.GetRuntimeLocation())
			}
			var vmList []corev1.VolumeMount
			for _, vm := range specVMList {
				resByre, resByteErr := json.Marshal(vm)
				if resByteErr != nil {
					log.Error(resByre)
					return nil, errors.New("invalid volumeMounts " + response.GetRuntimeLocation())
				}
				var newData corev1.VolumeMount
				jsonRes := json.Unmarshal(resByre, &newData)
				if jsonRes != nil {
					log.Error(jsonRes)
					return nil, errors.New("invalid volumeMounts " + response.GetRuntimeLocation())
				}
				vmList = append(vmList, newData)
			}
			vmList = append(vmList, container.VolumeMounts...)
			patches = append(patches, patchOperation{
				Op:    "add",
				Path:  "/spec/containers/" + itemString + "/volumeMounts",
				Value: vmList,
			})
		}

		// merge volumes
		if volumes, ok := specMap["volumes"]; ok {
			volumesList, suc := volumes.([]interface{})
			if !suc {
				log.Error("invalid volumes")
				return nil, errors.New("invalid volumes " + response.GetRuntimeLocation())
			}
			for _, volume := range volumesList {
				resByre, resByteErr := json.Marshal(volume)
				if resByteErr != nil {
					log.Error(resByre)
					return nil, errors.New("invalid volumes " + response.GetRuntimeLocation())
				}
				var newData corev1.Volume
				jsonRes := json.Unmarshal(resByre, &newData)
				if jsonRes != nil {
					log.Error(jsonRes)
					return nil, errors.New("invalid volumes " + response.GetRuntimeLocation())
				}
				if exists, eok := volumesExists[newData.Name]; !eok || !exists {
					volumesPod = append(volumesPod, newData)
					volumesExists[newData.Name] = true
				}
			}
		}
	}

	// add pod volumes
	if len(volumesPod) != 0 {
		volumesPod = append(volumesPod, pod.Spec.Volumes...)
		patches = append(patches, patchOperation{
			Op:    "add",
			Path:  "/spec/volumes",
			Value: volumesPod,
		})
	}

	if len(containerComputeunitMap) != 0 {
		log.Error("unexpected computeunit")
		return nil, errors.New("unexpected computeunit" + response.GetRuntimeLocation())
	}

	// add computeunit to annotation
	b, err := json.Marshal(computeunitList)
	if err != nil {
		log.Error(err)
		return nil, errors.New("unexpected computeunit" + response.GetRuntimeLocation())
	}
	annotations := map[string]string{}
	for k, v := range pod.Annotations {
		annotations[k] = v
	}
	annotations["openaios.4paradigm.com/computeunitList"] = string(b)
	log.Info(annotations)
	patches = append(patches, patchOperation{
		Op:    "add",
		Path:  "/metadata/annotations",
		Value: annotations,
	})

	return patches, nil
}

func main() {
	certPath := filepath.Join(tlsDir, tlsCertFile)
	keyPath := filepath.Join(tlsDir, tlsKeyFile)

	mux := http.NewServeMux()
	mux.Handle("/mutate", admitFuncHandler(applySecurityDefaults))
	server := &http.Server{
		// We listen on port 8443 such that we do not need root privileges or extra capabilities for this server.
		// The Service object will take care of mapping this port to the HTTPS port 443.
		Addr:    ":8443",
		Handler: mux,
	}
	//time.Sleep(100000 * time.Second)
	log.Fatal(server.ListenAndServeTLS(certPath, keyPath))
}
