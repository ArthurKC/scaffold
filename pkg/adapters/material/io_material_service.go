package material

import (
	"bufio"
	"errors"
	"fmt"
	"io"

	"github.com/ArthurKC/foundry/pkg/domains/material"
	"github.com/ArthurKC/foundry/pkg/domains/mold"
)

type IOMaterialService struct {
	writer  io.Writer
	scanner *bufio.Scanner
}

func NewIOMaterialService(w io.Writer, r io.Reader) *IOMaterialService {
	return &IOMaterialService{
		writer:  w,
		scanner: bufio.NewScanner(r),
	}
}

func (i *IOMaterialService) GetMaterial(requirements []*mold.MaterialRequirement) (mold.Material, error) {
	params := make(map[string]string, len(requirements))
	for _, r := range requirements {
		fmt.Fprintf(i.writer, "%[1]s: %[2]s\n%[1]s?: ", r.Name(), r.Description())
		if i.scanner.Scan() {
			params[r.Name()] = i.scanner.Text()
			continue
		}
		if err := i.scanner.Err(); err != nil {
			return nil, err
		}
		return nil, errors.New("input is cancelled")
	}
	return material.New("console_material", params), nil
}
