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
	"context"
	"flag"
	"github.com/4paradigm/openaios-platform/src/internal/billingclient"
	"github.com/4paradigm/openaios-platform/src/pineapple/conf"
	"github.com/4paradigm/openaios-platform/src/pineapple/controller/pvc"
	"github.com/4paradigm/openaios-platform/src/pineapple/handler"
	"github.com/4paradigm/openaios-platform/src/pineapple/utils"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"k8s.io/client-go/kubernetes"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	pvcChartsDir = flag.String("pvc-chartsdir", os.Getenv("PINEAPPLE_PVC_CHARTSDIR"),
		"pvc-chartsdir")
	pvcCephConfDir = flag.String("pvc-ceph-json-dir", os.Getenv("PINEAPPLE_PVC_CEPH_JSON_DIR"),
		"pvc-ceph-json-dir")
)

func kubernetesUserInitialize(ctx context.Context, client *kubernetes.Clientset, userID string) error {
	namespace := userID
	if err := utils.CreateNamespace(ctx, client, namespace); err != nil {
		return err
	}
	if err := utils.CreateUserRoleBindingWithEdit(ctx, client, namespace, userID); err != nil {
		return err
	}
	return nil
}

func userStoragePvcInit(ctx echo.Context) error {
	userID := ctx.Get("userID").(string)
	bearerToken, err := conf.GetKubeToken()
	if err != nil {
		return errors.WithMessage(err, "GetKubeToken error: ")
	}
	pvcImpl, err := pvc.NewPvcImpl(bearerToken, userID)
	if err != nil {
		return errors.WithMessage(err, "NewPvcImpl error: ")
	}
	pvcInfo, err := parseUserStorageConfigToPvcInfo(userID)
	if err != nil {
		return errors.WithMessage(err, "parseUserStorageConfigToPvcInfo error: ")
	}
	_, err = pvcImpl.Create(*pvcChartsDir, "user-storage", pvcInfo)
	if err != nil && !strings.Contains(err.Error(), "cannot re-use a name that is still in use") && !strings.Contains(err.Error(), "release: already exists") {
		return errors.WithMessage(err, "pvcImpl.Create error: ")
	}
	return nil
}

func parseUserStorageConfigToPvcInfo(userID string) (*pvc.PvcInfo, error) {

	storageQuota, err := strconv.ParseInt(conf.GetUserStorageQuotaBytes(), 10, 64)
	if err != nil {
		return nil, errors.Wrap(err, "ParseInt error: "+utils.GetRuntimeLocation())
	}
	storageQuota = storageQuota / 1000000000

	file, err := os.Open(*pvcCephConfDir)
	if err != nil {
		return nil, errors.Wrap(err, "os.Open error: "+utils.GetRuntimeLocation())
	}
	defer file.Close()
	byteFile, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, errors.Wrap(err, "ReadAll error: "+utils.GetRuntimeLocation())
	}
	//cephInfo := new(pvc.CephfsInfo)
	//if err := json.Unmarshal(byteFile, cephInfo); err != nil {
	//	return nil, errors.Wrap(err, "Unmarshal error: "+utils.GetRuntimeLocation())
	//}

	cephInfo := new(map[string]interface{})
	if err := yaml.Unmarshal(byteFile, cephInfo); err != nil {
		return nil, errors.Wrap(err, "Unmarshal error: "+utils.GetRuntimeLocation())
	}

	(*cephInfo)["path"] = (*cephInfo)["path"].(string) + "/" + userID

	pvcInfo := pvc.PvcInfo{
		UserID: userID,
		Cephfs: cephInfo,
		Capacity: &pvc.Capacity{
			Storage: strconv.FormatInt(storageQuota, 10) + "Gi",
		},
		CephSecret: &pvc.CephSecret{
			Key: os.Getenv("PINEAPPLE_PVC_CEPH_SECRET"),
		},
	}
	return &pvcInfo, nil
}

func InitUser(c echo.Context) error {
	userID := c.Get("userID").(string)
	// make the process idempotent
	// skip if user is existed
	k8sClient, err := utils.GetKubernetesClient()
	if err != nil {
		return err
	}

	if ok, err := utils.CheckUserReady(c.Request().Context(), k8sClient, userID); err != nil {
		return err
	} else if ok {
		return nil
	}
	time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
	c.Logger().Infof("initializing user %s", userID)
	// make user directory if not exists
	if err := handler.MkUserDir(userID); err != nil {
		return err
	}
	if err := kubernetesUserInitialize(c.Request().Context(), k8sClient, userID); err != nil {
		return err
	}
	if err := handler.CreateHarborUser(userID); err != nil {
		return err
	}
	if err := userStoragePvcInit(c); err != nil {
		return errors.WithMessage(err, "userStoragePvcInit error: ")
	}
	billingClient, err := billingclient.GetBillingClient(conf.GetBillingServerURL())
	if err != nil {
		return errors.Wrap(err, "cannot connect to billing server "+utils.GetRuntimeLocation())
	}
	if err := billingclient.InitUserBillingAccount(billingClient, userID, conf.GetInternalURL()); err != nil {
		return err
	}
	if err := utils.MarkUserAsReady(c.Request().Context(), k8sClient, userID); err != nil {
		return err
	}
	c.Logger().Infof("user %s is initialized", userID)
	return nil
}
