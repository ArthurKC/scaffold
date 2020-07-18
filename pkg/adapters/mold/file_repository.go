package mold

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/ArthurKC/foundry/pkg/domains/mold"
	"gopkg.in/yaml.v2"
)

type FileRepository struct {
}

func NewFileRepository() *FileRepository {
	return &FileRepository{}
}

func (f *FileRepository) metaFilePath(name string) string {
	return path.Join(name, "mold.yaml")
}

func (f *FileRepository) withTempFile(fnc func(string) error) error {
	tmpPath, err := ioutil.TempDir("", "*")
	if err != nil {
		return fmt.Errorf("cannot create temp directory: %w", err)
	}
	defer os.RemoveAll(tmpPath)

	return fnc(tmpPath)
}

func (f *FileRepository) saveMeta(m *mold.Mold, tmpPath string) error {
	metaPath := f.metaFilePath(tmpPath)
	conts := make([]*Parameter, 0, len(m.Requirements()))
	for _, r := range m.Requirements() {
		conts = append(conts, &Parameter{Name: r.Name(), Description: r.Description()})
	}
	y := moldYaml{Parameters: conts}
	data, err := yaml.Marshal(y)
	if err != nil {
		return fmt.Errorf("cannot marshal yaml: %w", err)
	}
	err = ioutil.WriteFile(metaPath, data, 0644)
	if err != nil {
		return fmt.Errorf("cannot save metafile, filePath=%s, %w", metaPath, err)
	}
	return nil
}

func (f *FileRepository) saveComponent(tmpPath string, c *mold.Component) error {
	cPath := path.Join(tmpPath, c.Path)
	if _, err := os.Stat(filepath.Dir(cPath)); os.IsNotExist(err) {
		err = os.MkdirAll(filepath.Dir(cPath), 0755)
		if err != nil {
			return fmt.Errorf("cannot prepare directory, filePath=%s: %w ", cPath, err)
		}
	}
	err := ioutil.WriteFile(cPath, []byte(c.Contents), 0644)
	if err != nil {
		return fmt.Errorf("cannot save component, filePath=%s: %w ", cPath, err)
	}
	return nil
}

// retreat the dst files not to delete when some error is occured in this process.
func (f *FileRepository) safeSwap(src string, dst string) error {
	tmp := ""
	if _, err := os.Stat(dst); !os.IsNotExist(err) {
		tmp = dst + ".backup"
		err = os.Rename(dst, tmp)
		if err != nil {
			return fmt.Errorf("cannot move original directory to temporal directory, tempPath=%s: %w", tmp, err)
		}
	}

	err := os.Rename(src, dst)
	if err != nil {
		return fmt.Errorf("cannot move to destination, src=%s, dest=%s: %w", src, dst, err)
	}

	if tmp != "" {
		// delete by defer
		err := os.Rename(tmp, src)
		if err != nil {
			return fmt.Errorf("cannot delete original files, filePath=%s: %w", tmp, err)
		}
	}
	return nil
}

func (f *FileRepository) Save(m *mold.Mold) error {
	return f.withTempFile(func(tmpPath string) error {
		err := f.saveMeta(m, tmpPath)
		if err != nil {
			return err
		}

		for _, c := range m.Components() {
			err := f.saveComponent(tmpPath, c)
			if err != nil {
				return err
			}
		}

		err = f.safeSwap(tmpPath, m.Name())
		if err != nil {
			return err
		}

		return nil
	})
}

func (f *FileRepository) FindByName(name string) (*mold.Mold, error) {
	metaPath := f.metaFilePath(name)
	bin, err := ioutil.ReadFile(metaPath)
	if err != nil {
		return nil, fmt.Errorf("cannot readfile: filePath = %s, %w", metaPath, err)
	}
	m := moldYaml{}
	if err := yaml.Unmarshal(bin, &m); err != nil {
		return nil, fmt.Errorf("cannot readfile: filePath = %s, %w", metaPath, err)
	}
	requirements := make([]*mold.MaterialRequirement, 0, len(m.Parameters))
	for _, c := range m.Parameters {
		requirements = append(requirements, mold.NewMaterialRequirement(c.Name, c.Description))
	}

	sources := make([]*mold.Component, 0)
	err = filepath.Walk(name, func(p string, info os.FileInfo, err error) error {
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

		src, err := filepath.Rel(name, p)
		if err != nil {
			return err
		}

		// remove if not gotmpl file
		if strings.ToLower(filepath.Ext(src)) != ".gotmpl" {
			return nil
		}

		cont, err := ioutil.ReadFile(filepath.Join(name, src))
		if err != nil {
			return err
		}

		sources = append(
			sources,
			mold.NewComponent(strings.TrimSuffix(src, filepath.Ext(src)), string(cont)),
		)
		return nil
	})
	if err != nil {
		return nil, err
	}

	return mold.New(name, sources, requirements), nil
}

type Parameter struct {
	Name        string
	Description string
}

type moldYaml struct {
	Parameters []*Parameter
}
