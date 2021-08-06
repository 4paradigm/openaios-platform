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

// Package helm provides utils for helm.
package helm

import (
	"github.com/4paradigm/openaios-platform/src/pineapple/conf"
	"github.com/4paradigm/openaios-platform/src/pineapple/utils"
	"github.com/pkg/errors"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/chart/loader"
	"net/http"
)

func DownloadChart(url string) (*chart.Chart, error) {
	realURL, err := getChartrepoURL(url)
	if err != nil {
		return nil, errors.Wrap(err, "getChartrepoURL url error: "+utils.GetRuntimeLocation())
	}
	client := new(http.Client)
	request, err := http.NewRequest("GET", realURL, nil)
	if err != nil {
		return nil, errors.Wrap(err, "NewRequest error: "+utils.GetRuntimeLocation())
	}
	request.Header.Add("Accept-Encoding", "gzip")
	request.SetBasicAuth(conf.GetHarborAdmin())
	resp, err := client.Do(request)
	if err != nil {
		return nil, errors.Wrap(err, "Get charts from url error: "+utils.GetRuntimeLocation())
	}
	defer resp.Body.Close()
	chart, err := loader.LoadArchive(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "LoadArchive from reader error: "+utils.GetRuntimeLocation())
	}
	return chart, nil
}

func DownloadChartFiles(url string) ([]*chart.File, error) {
	realURL, err := getChartrepoURL(url)
	if err != nil {
		return nil, errors.Wrap(err, "getChartrepoURL url error: "+utils.GetRuntimeLocation())
	}
	client := new(http.Client)
	request, err := http.NewRequest("GET", realURL, nil)
	if err != nil {
		return nil, errors.Wrap(err, "NewRequest error: "+utils.GetRuntimeLocation())
	}
	request.Header.Add("Accept-Encoding", "gzip")
	request.SetBasicAuth(conf.GetHarborAdmin())
	resp, err := client.Do(request)
	if err != nil {
		return nil, errors.Wrap(err, "Get charts from url error: "+utils.GetRuntimeLocation())
	}
	defer resp.Body.Close()

	bufferedFiles, err := loader.LoadArchiveFiles(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "LoadArchiveFiles error: "+utils.GetRuntimeLocation())
	}
	files, err := expandCharts(bufferedFiles)
	if err != nil {
		return nil, errors.WithMessage(err, "expandCharts error: ")
	}
	return files, nil
}

func expandCharts(bufferedFiles []*loader.BufferedFile) ([]*chart.File, error) {
	var files []*chart.File
	for _, bf := range bufferedFiles {
		files = append(files, &chart.File{Name: bf.Name, Data: bf.Data})
	}
	return files, nil
}

func getChartrepoURL(url string) (string, error) {
	// harborURL, _, _ := conf.GetHarborAddress()
	harborURL := conf.GetHarborURL()
	if harborURL == "" {
		return "", errors.New("harbor url is empty: " + utils.GetRuntimeLocation())
	}
	realURL := harborURL + "/" + url
	return realURL, nil
}
