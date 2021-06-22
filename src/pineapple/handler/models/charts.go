package models

type ImageInfo struct {
	Repository string `json:"repository" structs:"repository"`
	Tag        string `json:"tag" structs:"tag"`
	PullPolicy string `json:"pullPolicy" structs:"pullPolicy"`
}

type MountInfo struct {
	MountPath string `json:"mountPath"`
	SubPath   string `json:"subPath"`
	Name      string `json:"name"`
}

type VolumeMounts []MountInfo
