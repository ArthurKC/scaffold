package mold

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

	"github.com/ArthurKC/foundry/pkg/domains/mold"
	"gopkg.in/yaml.v2"
)

type FileService struct {
}

func NewFileService() *FileService {
	return &FileService{}
}

func (f *FileService) metaFilePath(name string) string {
	return path.Join(name, "mold.yaml")
}

func (f *FileService) newMoldComponent(path, content string, alreadMold bool) *mold.Component {
	if alreadMold {
		return mold.NewComponent(path, content)
	} else {
		return mold.NewComponent(path+".gotmpl", content)
	}
}

func (f *FileService) ImportFrom(path string, dstName string) (*mold.Mold, error) {
	alreadyMold := false
	metaPath := f.metaFilePath(path)
	m := moldYaml{make([]*Parameter, 0)}
	bin, err := ioutil.ReadFile(metaPath)
	if err == nil {
		alreadyMold = true
		if err := yaml.Unmarshal(bin, &m); err != nil {
			return nil, fmt.Errorf("cannot readfile: filePath = %s, %w", metaPath, err)
		}
	}
	requirements := make([]*mold.MaterialRequirement, 0, len(m.Parameters))
	for _, c := range m.Parameters {
		requirements = append(requirements, mold.NewMaterialRequirement(c.Name, c.Description))
	}

	sources := make([]*mold.Component, 0)
	err = filepath.Walk(path, func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		// skip matafile
		if p == metaPath {
			return nil
		}

		src, err := filepath.Rel(path, p)
		if err != nil {
			return err
		}

		cont, err := ioutil.ReadFile(filepath.Join(path, src))
		if err != nil {
			return err
		}

		sources = append(sources, f.newMoldComponent(src, string(cont), alreadyMold))
		return nil
	})
	if err != nil {
		return nil, err
	}

	return mold.New(dstName, sources, requirements), nil
}
