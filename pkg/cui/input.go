package cui

import (
	"bufio"
	"fmt"
	"os"
)

type Input struct {
	scanner *bufio.Scanner
}

func NewInput() *Input {
	return &Input{bufio.NewScanner(os.Stdin)}
}

func (i *Input) Ask(message string) string {
	fmt.Printf("%s?: ", message)
	for i.scanner.Scan() {
		return i.scanner.Text()
	}
	if err := i.scanner.Err(); err != nil {
		panic(err)
	}
	return ""
}
