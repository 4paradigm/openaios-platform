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
	"archive/tar"
	"archive/zip"
	"github.com/klauspost/compress/gzip"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"github.com/4paradigm/openaios-platform/src/internal/response"
	"github.com/4paradigm/openaios-platform/src/pineapple/apigen"
	"github.com/4paradigm/openaios-platform/src/pineapple/conf"
	"github.com/4paradigm/openaios-platform/src/pineapple/utils"
	"github.com/4paradigm/openaios-platform/src/pineapple/utils/helm"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type NamedFile struct {
	io.ReadCloser
	name string
}

func (nF *NamedFile) Name() string {
	return nF.name
}

func (handler *Handler) GetAppstoreChartList(ctx echo.Context) error {
	userId := ctx.Get("userID").(string)

	privateChartList, err := getChartList(apigen.ChartCategory_private, userId)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError).SetInternal(err)
	}
	publicGeneralChartList, err := getChartList(apigen.ChartCategory_public_community, string(apigen.ChartCategory_public_community))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError).SetInternal(err)
	}
	publicOfficialChartList, err := getChartList(apigen.ChartCategory_public_official, string(apigen.ChartCategory_public_official))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError).SetInternal(err)
	}
	publicPracticalChartList, err := getChartList(apigen.ChartCategory_public_practical, string(apigen.ChartCategory_public_practical))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError).SetInternal(err)
	}

	chartList := append(append(append(privateChartList, publicGeneralChartList...), publicOfficialChartList...), publicPracticalChartList...)
	resp := apigen.ChartMetaDataList{
		Items: &chartList,
	}
	return ctx.JSON(http.StatusOK, resp)
}

func getChartList(category apigen.ChartCategory, repo string) ([]apigen.ChartMetadata, error) {
	host, basepath, schemes := conf.GetHarborV1Address()
	client := utils.GetHarborClient(host, basepath, schemes)
	chartInfoEntryList, err := utils.GetChartrepoRepoCharts(client, repo)
	if err != nil {
		errors.WithMessage(err, "GetChartrepoRepoCharts error: ")
	}
	var chartMetadatas = []apigen.ChartMetadata{}
	for _, c := range chartInfoEntryList {
		name := c.Name
		version := c.LatestVersion
		icon := c.Icon
		chartVersionDetails, err := utils.GetChartrepoRepoChartsNameVersion(client, repo, *name, version)
		if err != nil {
			return nil, errors.WithMessagef(err, "GetChartrepoRepoChartsNameVersion err: Name: %s, Version: %s", name, version)
		}
		description := chartVersionDetails.Metadata.Description
		url := filepath.Join("chartrepo", repo, chartVersionDetails.Metadata.Urls[0])
		chartMetadatas = append(chartMetadatas, apigen.ChartMetadata{
			Category:    &category,
			Description: &description,
			IconLink:    &icon,
			Name:        name,
			Url:         &url,
			Version:     &version,
		})
	}

	//log.Printf("==========getChartList: repo: %s ==========\n%+v", repo, chartMetadatas)

	return chartMetadatas, nil
}

func (handler *Handler) GetAppstoreChart(ctx echo.Context, category apigen.ChartCategory, name string, version string) error {
	userId := ctx.Get("userID").(string)
	host, basepath, schemes := conf.GetHarborV1Address()
	client := utils.GetHarborClient(host, basepath, schemes)
	repo := userId
	if category != apigen.ChartCategory_private {
		repo = string(category)
	}

	// Metadata & Files
	chartVersionDetails, err := utils.GetChartrepoRepoChartsNameVersion(client, repo, name, version)
	if err != nil {
		if strings.Contains(err.Error(), "[404] getChartrepoRepoChartsNameVersionNotFound") {
			return response.BadRequestWithMessagef(ctx, "无法找到应用%s，版本%s", name, version)
		}
		return echo.NewHTTPError(http.StatusInternalServerError).SetInternal(err)
	}
	url := filepath.Join("chartrepo", repo, chartVersionDetails.Metadata.Urls[0])
	metadata := apigen.ChartMetadata{
		Category:    &category,
		Description: &chartVersionDetails.Metadata.Description,
		IconLink:    chartVersionDetails.Metadata.Icon,
		Name:        chartVersionDetails.Metadata.Name,
		Url:         &url,
		Version:     &version,
	}

	files := new(apigen.ChartFiles)
	chartFiles, err := helm.DownloadChartFiles(url)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError).SetInternal(err)
	}
	for _, f := range chartFiles {
		files.Set(f.Name, f.Data)
	}

	// VersionList
	chartVersions, err := utils.GetChartrepoRepoChartsName(client, repo, name)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError).SetInternal(err)
	}
	var versionList = []string{}
	for _, cv := range chartVersions {
		versionList = append(versionList, *cv.Version)
	}

	// response
	resp := apigen.Chart{
		Files:       files,
		Metadata:    &metadata,
		VersionList: &versionList,
	}

	return ctx.JSON(http.StatusOK, resp)
}

func (handler *Handler) UploadChart(ctx echo.Context, category apigen.ChartCategory) error {
	userID := ctx.Get("userID").(string)
	host, basepath, schemes := conf.GetHarborV1Address()
	client := utils.GetHarborClient(host, basepath, schemes)
	if category != apigen.ChartCategory_private {
		return response.BadRequestWithMessagef(ctx, "没有权限向%s上传chart", string(category))
	}

	file, err := ctx.FormFile("file")
	if err != nil {
		return response.BadRequestWithMessage(ctx, err.Error())
	}
	fileSuffix := path.Ext(file.Filename)
	if fileSuffix != ".zip" {
		return response.BadRequestWithMessagef(ctx, "没有权限上传%s文件", fileSuffix)
	}
	filePrefix := file.Filename[0 : len(file.Filename)-len(fileSuffix)]

	tarDir, err := parseZipToTgz(file)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError).SetInternal(err)
	}
	tarFile, err := os.Open(tarDir)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError).SetInternal(err)
	}
	defer func() {
		tarFile.Close()
		os.Remove(tarDir)
	}()

	namedFile := NamedFile{
		ReadCloser: tarFile,
		name:       filePrefix + ".tgz",
	}

	if err := utils.PostChartrepoRepoCharts(client, userID, &namedFile); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError).SetInternal(err)
	}

	return response.StatusOKNoContent(ctx)
}

func parseZipToTgz(file *multipart.FileHeader) (string, error) {
	zipFile, err := file.Open()
	if err != nil {
		return "", errors.Wrap(err, "open zip file error: "+utils.GetRuntimeLocation())
	}
	defer zipFile.Close()

	tempFile, err := ioutil.TempFile(os.TempDir(), "tmp-chart-tgz")
	if err != nil {
		return "", errors.Wrap(err, "TempFile error: "+utils.GetRuntimeLocation())
	}
	defer tempFile.Close()
	gw := gzip.NewWriter(tempFile)
	defer gw.Close()
	tw := tar.NewWriter(gw)
	defer tw.Close()

	zr, err := zip.NewReader(zipFile, file.Size)
	for _, f := range zr.File {
		rc, err := f.Open()
		if err != nil {
			return "", errors.Wrap(err, "Open zip file error: "+utils.GetRuntimeLocation())
		}
		header, err := tar.FileInfoHeader(f.FileInfo(), "")
		if err != nil {
			return "", errors.Wrap(err, "FileInfoHeader error: "+utils.GetRuntimeLocation())
		}
		header.Name = f.Name
		if err := tw.WriteHeader(header); err != nil {
			return "", errors.Wrap(err, "WriteHeader error: "+utils.GetRuntimeLocation())

		}
		_, err = io.Copy(tw, rc)
		if err != nil {
			return "", errors.Wrap(err, "Copy error: "+utils.GetRuntimeLocation())
		}
		rc.Close()
	}

	return tempFile.Name(), nil
}
