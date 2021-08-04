package helm

import (
	"bytes"
	"fmt"
	"github.com/pkg/errors"
	"github.com/4paradigm/openaios-platform/src/pineapple/utils"
	"path/filepath"
	"sigs.k8s.io/kustomize/api/filesys"
	"sigs.k8s.io/kustomize/api/konfig"
)
import "sigs.k8s.io/kustomize/api/krusty"

type PostRendererImpl struct {
	fSys filesys.FileSystem
}

func NewPostRendererImpl() *PostRendererImpl {
	return &PostRendererImpl{
		fSys: filesys.MakeFsInMemory(),
	}
}

func (p *PostRendererImpl) WriteKustomzation(path string, content string) error {
	err := p.fSys.WriteFile(
		filepath.Join(
			path,
			konfig.DefaultKustomizationFileName()), []byte(`
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization
`+content))
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("unexpected error while writing Kustomization to %s: ", path)+utils.GetRuntimeLocation())
	}
	return nil
}

func (p *PostRendererImpl) WriteFile(path string, content string) error {
	err := p.fSys.WriteFile(path, []byte(content))
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("unexpected error while writing file to %s: ", path)+utils.GetRuntimeLocation())
	}
	return nil
}

func (p *PostRendererImpl) Run(renderedManifests *bytes.Buffer) (*bytes.Buffer, error) {
	p.WriteFile("all.yaml", renderedManifests.String())
	options := krusty.MakeDefaultOptions()
	kustomizer := krusty.MakeKustomizer(options)
	result, err := kustomizer.Run(p.fSys, ".")
	if err != nil {
		return nil, err
	}
	b, err := result.AsYaml()
	if err != nil {
		return nil, err
	}
	buf := bytes.NewBuffer(b)
	return buf, nil
}
