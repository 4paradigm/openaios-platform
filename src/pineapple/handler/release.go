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
	"github.com/4paradigm/openaios-platform/src/internal/response"
	"github.com/4paradigm/openaios-platform/src/pineapple/apigen/internalapigen"
	"github.com/4paradigm/openaios-platform/src/pineapple/conf"
	"github.com/4paradigm/openaios-platform/src/pineapple/controller/application"
	"github.com/4paradigm/openaios-platform/src/pineapple/controller/environment"
	"github.com/4paradigm/openaios-platform/src/pineapple/utils/helm"
	"net/http"
)

func (handler *Handler) DeleteReleases(ctx echo.Context, params internalapigen.DeleteReleasesParams) error {
	if params.User == nil {
		return response.BadRequestWithMessagef(ctx, "request query user")
	}
	userId := *params.User
	bearerToken, err := conf.GetKubeToken()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError).SetInternal(err)
	}
	helmImpl, err := helm.NewImpl(bearerToken, userId)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError).SetInternal(err)
	}
	envs, err := helmImpl.List(environment.EnvPrefix)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError).SetInternal(err)
	}
	apps, err := helmImpl.List(application.AppPrefix)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError).SetInternal(err)
	}
	releases := append(envs, apps...)
	description := "ran out of credit"
	if err := helmImpl.DeleteListWithKeepHistory(releases, description); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError).SetInternal(err)
	}
	return response.StatusOKNoContent(ctx)
}
