package metal

import (
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

	"github.com/ArthurKC/foundry/pkg/domains/metal"
	"gopkg.in/yaml.v2"
)

type Constituent struct {
	Name        string
	Description string
}

type MoldSourceMeta struct {
	Constituents []*Constituent
}

type MoldSource struct {
	meta    *MoldSourceMeta
	sources []string
	moldDir string
}

func NewMoldSource(moldDir string) (*MoldSource, error) {
	metaPath := path.Join(moldDir, "foundry.yaml")
	f, err := ioutil.ReadFile(metaPath)
	if err != nil {
		return nil, err
	}
	m := MoldSourceMeta{}
	if err := yaml.Unmarshal(f, &m); err != nil {
		return nil, err
	}

	sources := make([]string, 0)
	err = filepath.Walk(moldDir, func(p string, info os.FileInfo, err error) error {
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

		src, err := filepath.Rel(moldDir, p)
		if err != nil {
			return err
		}

		sources = append(sources, src)
		return nil
	})
	if err != nil {
		return nil, err
	}

	return &MoldSource{&m, sources, moldDir}, nil
}

func (t *MoldSource) Constituents() []*metal.Constituent {
	ret := make([]*metal.Constituent, 0, len(t.meta.Constituents))
	for _, p := range t.meta.Constituents {
		ret = append(ret, &metal.Constituent{Name: p.Name, Description: p.Description})
	}
	return ret
}

func (t *MoldSource) Paths() []string {
	return t.sources
}

func (t *MoldSource) Source(path string) string {
	f, err := ioutil.ReadFile(filepath.Join(t.moldDir, path))
	if err != nil {
		panic(err)
	}
	return string(f)
}
