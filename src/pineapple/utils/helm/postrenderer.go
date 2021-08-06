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
	"bytes"
	"fmt"
	"github.com/4paradigm/openaios-platform/src/pineapple/utils"
	"github.com/pkg/errors"
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
