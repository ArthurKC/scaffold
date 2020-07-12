package metal

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

type Output struct {
	destDir string
}

func NewOutput(destDir string) *Output {
	return &Output{destDir}
}

func (o *Output) Write(path string, content string) {
	destPath := filepath.Join(o.destDir, path)
	_, err := os.Stat(destPath)
	if err == nil {
		fmt.Printf("Skipped (already exists) %s\n", destPath)
		return
	}

	err = os.MkdirAll(filepath.Dir(destPath), os.ModePerm)
	if err != nil {
		panic(err)
	}
	if err = ioutil.WriteFile(destPath, []byte(content), 0644); err != nil {
		panic(err)
	}
	fmt.Printf("created %s\n", destPath)
}

func (o *Output) DestDir() string {
	return o.destDir
}
