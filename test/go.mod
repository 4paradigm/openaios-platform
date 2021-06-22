module github.com/4paradigm/openaios-platform/test

go 1.14

require (
	github.com/BurntSushi/toml v0.3.1
	github.com/deepmap/oapi-codegen v1.5.5
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/onsi/ginkgo v1.16.2
	github.com/onsi/gomega v1.12.0
	github.com/pkg/errors v0.8.1
	golang.org/x/sys v0.0.0-20210510120138-977fb7262007 // indirect
)

replace (
	github.com/docker/distribution => github.com/docker/distribution v0.0.0-20191216044856-a8371794149d
	github.com/docker/docker => github.com/moby/moby v17.12.0-ce-rc1.0.20200618181300-9dc6525e6118+incompatible
)
