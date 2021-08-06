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

package handler

import (
	"github.com/labstack/gommon/log"
	"github.com/4paradigm/openaios-platform/src/billing/conf"
	"github.com/4paradigm/openaios-platform/src/billing/utils"
	"github.com/4paradigm/openaios-platform/src/internal/mongodb"
)

type Handler struct{}

var mongodbUrl = conf.GetMongodbUrl()

func NewHandler() Handler {
	return Handler{}
}

func InitBillingServer() {
	// init mongodb collection
	client, err := mongodb.GetMongodbClient(mongodbUrl)
	defer mongodb.KillMongodbClient(client)
	if err != nil {
		log.Error(err.Error())
		return
	}
	if err = utils.InitColl(client); err != nil {
		log.Error(err.Error())
	}
}
