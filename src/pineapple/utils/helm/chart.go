package helm

import (
	"github.com/pkg/errors"
	"github.com/4paradigm/openaios-platform/src/pineapple/conf"
	"github.com/4paradigm/openaios-platform/src/pineapple/utils"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/chart/loader"
	"net/http"
)

func DownloadChart(url string) (*chart.Chart, error) {
	realUrl, err := getChartrepoUrl(url)
	if err != nil {
		return nil, errors.Wrap(err, "getChartrepoUrl url error: "+utils.GetRuntimeLocation())
	}
	client := new(http.Client)
	request, err := http.NewRequest("GET", realUrl, nil)
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
	realUrl, err := getChartrepoUrl(url)
	if err != nil {
		return nil, errors.Wrap(err, "getChartrepoUrl url error: "+utils.GetRuntimeLocation())
	}
	client := new(http.Client)
	request, err := http.NewRequest("GET", realUrl, nil)
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

func getChartrepoUrl(url string) (string, error) {
	// harborUrl, _, _ := conf.GetHarborAddress()
	harborUrl := conf.GetHarborURL()
	if harborUrl == "" {
		return "", errors.New("harbor url is empty: " + utils.GetRuntimeLocation())
	}
	realUrl := harborUrl + "/" + url
	return realUrl, nil
}
