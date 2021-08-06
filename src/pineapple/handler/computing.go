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
	"github.com/pkg/errors"
	"github.com/4paradigm/openaios-platform/src/internal/billingclient"
	"github.com/4paradigm/openaios-platform/src/pineapple/apigen"
	"github.com/4paradigm/openaios-platform/src/pineapple/conf"
	"github.com/4paradigm/openaios-platform/src/pineapple/utils"
	"net/http"
)

func (handler *Handler) GetComputingUnitSpecs(ctx echo.Context) error {
	userId := ctx.Get("userID").(string)
	billingClient, err := billingclient.GetBillingClient(conf.GetBillingServerURL())
	if err != nil {
		return echo.NewHTTPError(
			http.StatusInternalServerError, "cannot connect to billing server").SetInternal(
			errors.Wrap(err, utils.GetRuntimeLocation()))
	}
	computeunitInfos, err := billingclient.GetComputeUnitListByUserID(billingClient, userId)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError).SetInternal(err)
	}
	result := make([]apigen.ComputeUnitSpec, len(computeunitInfos))

	for i, info := range computeunitInfos {
		id := apigen.ComputeUnitId(*(info.Id))
		result[i].Name = info.Id
		result[i].Id = &id
		result[i].Description = info.Description
		result[i].Price = info.Price
	}

	return ctx.JSON(http.StatusOK, result)
}
