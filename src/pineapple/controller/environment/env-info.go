package environment

import (
	"github.com/fatih/structs"
	"github.com/4paradigm/openaios-platform/src/pineapple/conf"
	"github.com/4paradigm/openaios-platform/src/pineapple/handler/models"
	"github.com/4paradigm/openaios-platform/src/pineapple/utils/helm"
)

type EnvironmentInfo struct {
	*helm.PineappleInfo
	Image        models.ImageInfo
	ServerType   ServerType
	SshKey       string
	JupyterToken string
	PvcClaimName string
	VolumeMounts models.VolumeMounts
	ResourceId   string
}

func (e *EnvironmentInfo) CreateEnvValues() (map[string]interface{}, error) {
	envValues := make(map[string]interface{})
	envValues["image"] = structs.Map(e.Image)

	envValues["serverType"] = structs.Map(e.ServerType)
	envValues["ssh"] = map[string]string{
		"sshKey": e.SshKey,
	}
	envValues["jupyter"] = map[string]string{
		"token": e.JupyterToken,
	}
	envValues["pvc"] = map[string]string{
		"claimName": e.PvcClaimName,
	}
	envValues["volumeMounts"] = e.VolumeMounts
	envValues["ingress"] = map[string]interface{}{
		"host":      conf.GetExternalURLHost(),
		"enableTLS": conf.GetExternalTLS(),
	}

	envValues["pineapple"] = map[string]interface{}{
		"belongTo": "user",
		"default": map[string]interface{}{
			"resourceId": e.ResourceId,
		},
	}

	return envValues, nil
}
