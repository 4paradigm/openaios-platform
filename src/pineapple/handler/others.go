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
	"context"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"github.com/4paradigm/openaios-platform/src/pineapple/apigen"
	"github.com/4paradigm/openaios-platform/src/pineapple/utils"
	"io"
	"net/http"
	"time"
)

func (handler *Handler) GetContainerLog(ctx echo.Context, podName apigen.PodName, params apigen.GetContainerLogParams) error {
	userID := ctx.Get("userID").(string)
	kubeClient, err := utils.GetKubernetesClient()
	if err != nil {
		return echo.NewHTTPError(
			http.StatusInternalServerError, "cannot connect to kubernetes cluster").SetInternal(
			errors.Wrap(err, utils.GetRuntimeLocation()))
	}

	stream, err := utils.GetContainerLog(context.TODO(), kubeClient, userID,
		string(podName), (*string)(params.ContainerName), params.TailLines)
	if err != nil {
		return echo.NewHTTPError(
			http.StatusInternalServerError, "cannot get log").SetInternal(
			errors.Wrap(err, utils.GetRuntimeLocation()))
	}
	defer func() {
		err = stream.Close()
		if err != nil {
			ctx.Logger().Error(err)
		}
	}()

	ctx.Response().Header().Set(echo.HeaderContentType, "text/plain")
	ctx.Response().Header().Set("Connection", "keep-alive")
	ctx.Response().Header().Set("Access-Control-Allow-Origin", "*")
	ctx.Response().Header().Set("Access-Control-Allow-Methods", "*")
	ctx.Response().Header().Set("Transfer-Encoding", "chunked")
	ctx.Response().Header().Set("X-Content-Type-Options", "nosniff")
	ctx.Response().WriteHeader(http.StatusOK)

	for {
		buf := make([]byte, 2000)
		byteNum, err := stream.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			return echo.NewHTTPError(
				http.StatusInternalServerError, "cannot read log").SetInternal(
				errors.Wrap(err, utils.GetRuntimeLocation()))
		}
		if byteNum == 0 {
			time.Sleep(1 * time.Second)
			continue
		}
		if _, err = ctx.Response().Writer.Write(buf); err != nil {
			return echo.NewHTTPError(
				http.StatusInternalServerError, "cannot read log").SetInternal(
				errors.Wrap(err, utils.GetRuntimeLocation()))
		}
		ctx.Response().Flush()
	}
	return nil
}
