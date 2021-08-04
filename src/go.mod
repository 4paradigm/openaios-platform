module github.com/4paradigm/openaios-platform/src

go 1.14

require (
	code.cloudfoundry.org/bytefmt v0.0.0-20200131002437-cf55d5288a48
	github.com/AimAlex/harbor-client v0.2.0
	github.com/NYTimes/gziphandler v1.1.1
	github.com/asaskevich/govalidator v0.0.0-20210307081110-f21760c49a8d // indirect
	github.com/coreos/go-oidc/v3 v3.0.0
	github.com/creack/pty v1.1.11
	github.com/deepmap/oapi-codegen v1.5.5
	github.com/elazarl/go-bindata-assetfs v1.0.1
	github.com/fatih/structs v1.1.0
	github.com/go-openapi/runtime v0.19.5 // do not use go-openapi >= v0.20.0
	github.com/go-openapi/strfmt v0.19.5
	github.com/gorilla/websocket v1.4.2
	github.com/klauspost/compress v1.9.5
	github.com/labstack/echo-contrib v0.9.0
	github.com/labstack/echo/v4 v4.2.1
	github.com/labstack/gommon v0.3.0
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/patrickmn/go-cache v2.1.0+incompatible
	github.com/pkg/errors v0.9.1
	github.com/robfig/cron/v3 v3.0.1
	github.com/urfave/cli/v2 v2.3.0
	go.mongodb.org/mongo-driver v1.5.0
	golang.org/x/net v0.0.0-20210316092652-d523dce5a7f4 // indirect
	gopkg.in/yaml.v3 v3.0.0-20200615113413-eeeca48fe776
	helm.sh/helm/v3 v3.6.1
	k8s.io/api v0.21.2
	k8s.io/apimachinery v0.21.2
	k8s.io/cli-runtime v0.21.2
	k8s.io/client-go v0.21.2
	sigs.k8s.io/kustomize/api v0.8.10
)

replace (
	github.com/docker/distribution => github.com/docker/distribution v0.0.0-20191216044856-a8371794149d
	github.com/docker/docker => github.com/moby/moby v17.12.0-ce-rc1.0.20200618181300-9dc6525e6118+incompatible
)
