package mold

import (
	"fmt"
	"io"

	"github.com/ArthurKC/foundry/pkg/usecases/mold"
)

type IOOutputPort struct {
	writer io.Writer
}

func NewIOOutputPort(w io.Writer) *IOOutputPort {
	return &IOOutputPort{
		writer: w,
	}
}
func (i *IOOutputPort) RenderPourResponse(r *mold.PourResponse) {
	fmt.Fprintf(i.writer, "completed casting! (mold = %s)", r.MoldName)
}

func (i *IOOutputPort) RenderPourError(e *mold.PourError) {
	fmt.Fprintf(i.writer, "failed casting. (mold = %s, err = %v)", e.Req.MoldName, e.Err)
}

func (i *IOOutputPort) RenderCreateResponse(r *mold.CreateResponse) {
	fmt.Fprintf(i.writer, "completed create! (mold = %s)", r.MoldName)
}

func (i *IOOutputPort) RenderCreateError(e *mold.CreateError) {
	fmt.Fprintf(i.writer, "failed create. (mold = %s, err = %v)", e.Req.MoldName, e.Err)
}
