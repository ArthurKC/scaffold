package generator

import (
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

	"github.com/ArthurKC/scaffold/pkg/domains/generator"
	"gopkg.in/yaml.v2"
)

type Parameter struct {
	Name        string
	Description string
}

type TemplateSourceMeta struct {
	Params []*Parameter
}

type TemplateSource struct {
	meta        *TemplateSourceMeta
	sources     []string
	templateDir string
}

func NewTemplateSource(templateDir string) (*TemplateSource, error) {
	metaPath := path.Join(templateDir, "scaffold.yaml")
	f, err := ioutil.ReadFile(metaPath)
	if err != nil {
		return nil, err
	}
	m := TemplateSourceMeta{}
	if err := yaml.Unmarshal(f, &m); err != nil {
		return nil, err
	}

	sources := make([]string, 0)
	err = filepath.Walk(templateDir, func(p string, info os.FileInfo, err error) error {
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

		src, err := filepath.Rel(templateDir, p)
		if err != nil {
			return err
		}

		sources = append(sources, src)
		return nil
	})
	if err != nil {
		return nil, err
	}

	return &TemplateSource{&m, sources, templateDir}, nil
}

func (t *TemplateSource) Params() []*generator.Parameter {
	ret := make([]*generator.Parameter, 0, len(t.meta.Params))
	for _, p := range t.meta.Params {
		ret = append(ret, &generator.Parameter{Name: p.Name, Description: p.Description})
	}
	return ret
}

func (t *TemplateSource) Paths() []string {
	return t.sources
}

func (t *TemplateSource) Source(path string) string {
	f, err := ioutil.ReadFile(filepath.Join(t.templateDir, path))
	if err != nil {
		panic(err)
	}
	return string(f)
}
