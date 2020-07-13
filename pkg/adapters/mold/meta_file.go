package mold

import (
	"io/ioutil"
	"path"

	"gopkg.in/yaml.v2"
)

type MetaFile struct {
	dirPath string
	meta    *metaYaml
}

type metaYaml struct {
	Parameters []*constituent
}

type constituent struct {
	Name        string
	Description string
}

func NewMetaFile(dirPath string) *MetaFile {
	return &MetaFile{
		dirPath: dirPath,
	}
}

func (m *MetaFile) filePath() string {
	return path.Join(m.dirPath, "mold.yaml")
}

func (m *MetaFile) Load() error {
	f, err := ioutil.ReadFile(m.filePath())
	if err != nil {
		return err
	}
	meta := metaYaml{}
	if err := yaml.Unmarshal(f, &meta); err != nil {
		return err
	}
	m.meta = &meta
	return nil
}

func (m *MetaFile) Save() error {
	f, err := yaml.Marshal(m.meta)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(m.filePath(), f, 0644)
}

func (m *MetaFile) Initialize() error {
	m.meta = &metaYaml{[]*constituent{}}
	return m.Save()
}
