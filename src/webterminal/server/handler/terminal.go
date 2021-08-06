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
	"github.com/4paradigm/openaios-platform/src/webterminal/server/apigen"
	"github.com/4paradigm/openaios-platform/src/webterminal/server/controller/webterminal"
	"github.com/labstack/echo/v4"
	"net/http"
)

func (h *Handler) GetTerminal(ctx echo.Context, params apigen.GetTerminalParams) error {
	userID := ctx.Get("userID").(string)
	bearerToken := ctx.Get("bearerToken").(string)

	podName := params.Pod
	containerName := params.Container

	terminal := webterminal.Webterminal{
		Namespace:     userID,
		PodName:       podName,
		ContainerName: containerName,
		BearerToken:   bearerToken,
	}

	terminalURL, err := webterminal.GetWebterminal(terminal)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError).SetInternal(err)
	}

	//return ctx.Redirect(http.StatusSeeOther, terminalURL)
	return ctx.JSON(http.StatusOK, apigen.WebterminalInfo{Url: &terminalURL})
}
