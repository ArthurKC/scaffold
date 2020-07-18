package mold

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/ArthurKC/foundry/pkg/domains/mold"
)

type FileProductService struct {
	notifier io.Writer
}

func NewFileProductService(notifier io.Writer) *FileProductService {
	return &FileProductService{notifier}
}

func (f *FileProductService) SaveProduct(destDir string, product *mold.Product) error {
	for _, c := range product.Components() {
		destPath := filepath.Join(destDir, c.Path())

		if _, err := os.Stat(destPath); !os.IsNotExist(err) {
			fmt.Fprintf(f.notifier, "Skipped (already exists) %s\n", destPath)
			continue
		}
		err := os.MkdirAll(filepath.Dir(destPath), os.ModePerm)
		if err != nil {
			return err
		}
		if err = ioutil.WriteFile(destPath, []byte(c.Contents()), 0644); err != nil {
			return err
		}
		fmt.Fprintf(f.notifier, "created %s\n", destPath)
	}
	return nil
}
