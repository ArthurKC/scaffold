package metal

import (
	"bufio"
	"fmt"
	"os"

	"github.com/ArthurKC/foundry/pkg/domains/metal"
)

type Input struct {
	scanner *bufio.Scanner
}

func NewInput() *Input {
	return &Input{bufio.NewScanner(os.Stdin)}
}

func (i *Input) Ask(p *metal.Constituent) string {
	fmt.Printf("%[1]s: %[2]s\n%[1]s?: ", p.Name, p.Description)
	for i.scanner.Scan() {
		return i.scanner.Text()
	}
	if err := i.scanner.Err(); err != nil {
		panic(err)
	}
	return ""
}
