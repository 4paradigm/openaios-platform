package conf

import (
	"flag"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"net"
	"os"
	"strconv"
	"strings"
)

var (
	externalURL = flag.String("external-url", os.Getenv("PINEAPPLE_EXTERNAL_URL"),
		"external-url")
	externalHost = flag.String("external-host", os.Getenv("PINEAPPLE_EXTERNAL_HOST"),
		"external-host")
	externalTLS = flag.String("external-tls", os.Getenv("PINEAPPLE_EXTERNAL_TLS"),
		"external-tls")
	internalURL = flag.String("internal-url", os.Getenv("PINEAPPLE_INTERNAL_URL"),
		"internal-url")
	kubeApiServer = flag.String("kube-apiserver", getKubeApiServer(),
		"kube-apiserver")
	kubeCaFile = flag.String("kube-cafile", getKubeCaFile(),
		"kube-cafile")
	kubeTokenFile = flag.String("kube-token", getKubeTokenFile(),
		"kube-token")
	userStorageQuotaBytes = flag.String("storage-user-quota", os.Getenv("PINEAPPLE_STORAGE_USER_QUOTA"),
		"set user directory quota")
	billingServerUrl = flag.String("billing-server-url", os.Getenv("PINEAPPLE_BILLING_SERVER_URL"),
		"billing server url")
	harborURL = flag.String("images-harbor-url", os.Getenv("PINEAPPLE_HARBOR_URL"),
		"harbor url")
	harborBasepath = flag.String("images-harbor-basepath", os.Getenv("PINEAPPLE_HARBOR_BASEPATH"),
		"harbor basepath")
	harborV1Basepath = flag.String("images-harborV1-basepath", os.Getenv("PINEAPPLE_HARBORV1_BASEPATH"),
		"harborV1 basepath")
	harborAdminUsername = flag.String("images-harbor-admin-username", os.Getenv("PINEAPPLE_HARBOR_ADMIN_USERNAME"),
		"harbor admin username")
	harborAdminPassword = flag.String("images-harbor-admin-password", os.Getenv("PINEAPPLE_HARBOR_ADMIN_PASSWORD"),
		"harbor admin password")
	harborStorageLimit = flag.String("images-storage-limit", os.Getenv("PINEAPPLE_HARBOR_STORAGE_LIMIT"),
		"user image storage limit")
	mongodbUrl = flag.String("mongodb-url", os.Getenv("PINEAPPLE_MONGODB_URL"),
		"mongodb url")
	mongodbDatabase = flag.String("mongodb-database", os.Getenv("PINEAPPLE_MONGODB_DATABASE"),
		"mongodb database")
)

func GetExternalURL() string {
	return *externalURL
}

func GetExternalURLHost() string {
	return *externalHost
}

func GetExternalTLS() bool {
	b, err := strconv.ParseBool(*externalTLS)
	if err != nil {
		return false
	}
	return b
}

func GetInternalURL() string {
	return *internalURL
}

func GetKubeApiServer() string {
	return *kubeApiServer
}

func GetKubeCaFile() string {
	return *kubeCaFile
}

func GetKubeToken() (string, error) {
	token, err := ioutil.ReadFile(*kubeTokenFile)
	if err != nil {
		return "", err
	}
	return string(token), nil
}
func GetUserStorageQuotaBytes() string {
	return *userStorageQuotaBytes
}

func GetBillingServerURL() string {
	return *billingServerUrl
}

func GetHarborURL() (string) {
	return *harborURL
}

func GetHarborAddress() (string, string, []string) {
	splitHost := strings.Split(*harborURL, "://")
	return splitHost[1], *harborBasepath, []string{splitHost[0]}
}

func GetHarborV1Address() (string, string, []string) {
	splitHost := strings.Split(*harborURL, "://")
	return splitHost[1], *harborV1Basepath, []string{splitHost[0]}
}

func GetHarborAdmin() (string, string) {
	return *harborAdminUsername, *harborAdminPassword
}

func GetHarborStorageLimit() (*int64, error) {
	var limit *int64
	if *harborStorageLimit == "" {
		limit = nil
	} else {
		parse, err := strconv.ParseInt(*harborStorageLimit, 10, 64)
		if err != nil {
			return nil, err
		}
		limit = &parse
	}
	return limit, nil
}

func GetAppConf() (map[string]interface{}, error) {
	file, err := os.Open("/root/config/appConf.yaml")
	if err != nil {
		return nil, errors.Wrap(err, "os.Open error: ")
	}
	defer file.Close()
	byteFile, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, errors.Wrap(err, "ReadAll error: ")
	}
	appConf := make(map[string]interface{})
	if err := yaml.Unmarshal(byteFile, appConf); err != nil {
		return nil, errors.Wrap(err, "unmarshal appConf error: ")
	}
	return appConf, nil
}

func getKubeApiServer() string {
	if os.Getenv("PINEAPPLE_ENV_KUBEAPISERVER") != "" {
		return os.Getenv("PINEAPPLE_ENV_KUBEAPISERVER")
	}
	return "https://" + net.JoinHostPort(os.Getenv("KUBERNETES_SERVICE_HOST"), os.Getenv("KUBERNETES_SERVICE_PORT"))
}

func getKubeCaFile() string {
	if os.Getenv("PINEAPPLE_ENV_KUBECAFILE") != "" {
		return os.Getenv("PINEAPPLE_ENV_KUBECAFILE")
	}
	return "/var/run/secrets/kubernetes.io/serviceaccount/ca.crt"
}

func getKubeTokenFile() string {
	if os.Getenv("PINEAPPLE_KUBE_TOKENFILE") != "" {
		return os.Getenv("PINEAPPLE_KUBE_TOKENFILE")
	}
	return "/var/run/secrets/kubernetes.io/serviceaccount/token"
}

func GetMongodbUrl() string {
	return *mongodbUrl
}

func GetMongodbDatabase() string {
	return *mongodbDatabase
}
