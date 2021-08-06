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
	"github.com/labstack/gommon/log"
	"github.com/pkg/errors"
	"github.com/robfig/cron/v3"
	"github.com/4paradigm/openaios-platform/src/billing/conf"
	"github.com/4paradigm/openaios-platform/src/billing/utils"
	"github.com/4paradigm/openaios-platform/src/internal/mongodb"
	"github.com/4paradigm/openaios-platform/src/internal/response"
	"time"
)

func run() {
	mongodbUrl := conf.GetMongodbUrl()

	// init mongodb client
	client, err := mongodb.GetMongodbClient(mongodbUrl)
	defer mongodb.KillMongodbClient(client)
	if err != nil {
		log.Error(err.Error())
		return
	}

	// init k8s client
	kubeClient, err := utils.GetKubernetesClient()
	if err != nil {
		log.Error(err.Error())
	}

	// get user pod list in k8s
	podList, err := utils.GetPodList(kubeClient, "openaios.4paradigm.com/app=true", "")
	if err != nil {
		log.Error(err.Error())
		return
	}

	// init user bill map and price list
	billMap := map[string]float64{}
	priceMap, err := utils.GetPriceMap(client)
	if err != nil {
		log.Error(err.Error())
		return
	}

	// deal with user account for each pod
	for _, pod := range *podList {
		// read pod information
		userID := pod.Namespace
		podName := pod.Name
		podUID := pod.UID
		instanceID := pod.Labels["instanceID"]
		updateTime := time.Now()
		startTime := updateTime
		if pod.Status.StartTime != nil {
			startTime = pod.Status.StartTime.Time
		}
		status := string(pod.Status.Phase)

		// only care running pod
		if status != "Running" || pod.DeletionTimestamp != nil {
			continue
		}

		var computeunitList []string
		computeunitString := pod.Annotations["openaios.4paradigm.com/computeunitList"]
		if computeunitString == "" {
			log.Warnf("user %s's pod %s does not have computeunit list", userID, podName)
			continue
		}
		err = json.Unmarshal([]byte(computeunitString), &computeunitList)
		if err != nil {
			log.Error(err)
			continue
		}
		for _, computeunit := range computeunitList {
			computeunitPrice, ok := priceMap[computeunit]
			if !ok {
				log.Warn("price map does not have such compute price " + computeunit)
			}
			billMap[userID] -= computeunitPrice
		}

		// update mongodb
		err = utils.UpdateOrInsertPod(client, utils.PodInfo{UserId: userID, PodName: podName, PodUID: podUID,
			InstanceId: instanceID, ComputeunitList: computeunitList, StartTime: startTime, UpdateTime: updateTime})
		if err != nil {
			log.Error(err)
		}
	}

	// check user account map
	for userId, cost := range billMap {
		err = utils.ModifyUserBalance(client, userId, cost)
		if err != nil {
			log.Error(errors.Wrap(err, response.GetRuntimeLocation()))
		}
	}

	// check user no balance
	err = utils.CheckUserNoBalance(client, billMap)
	if err != nil {
		log.Error(errors.Wrap(err, response.GetRuntimeLocation()))
	}
}

func heartbeat() {
	c := cron.New()
	_, err := c.AddFunc("* * * * *", run)
	//_, err := c.AddFunc("@every 1s", run)
	if err != nil {
		log.Error(errors.Wrap(err, response.GetRuntimeLocation()))
		return
	}
	c.Start()
	select {}
}
