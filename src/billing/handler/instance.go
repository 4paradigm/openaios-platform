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
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/4paradigm/openaios-platform/src/billing/apigen"
	"github.com/4paradigm/openaios-platform/src/billing/utils"
	"github.com/4paradigm/openaios-platform/src/internal/mongodb"
	"github.com/4paradigm/openaios-platform/src/internal/response"
	"net/http"
)

func (h Handler) PostInstance(ctx echo.Context, params apigen.PostInstanceParams) error {
	client, err := mongodb.GetMongodbClient(mongodbUrl)
	defer mongodb.KillMongodbClient(client)
	if err != nil {
		log.Error(err.Error())
		return ctx.String(http.StatusInternalServerError, "cannot connect to mongodb.")
	}

	instanceId, err := utils.CreateInstance(client, string(params.UserId), params.InstanceName)
	if err != nil {
		log.Error(err.Error())
		return response.BadRequestWithMessage(ctx, err.Error())
	}
	return ctx.String(http.StatusOK, instanceId)
}
