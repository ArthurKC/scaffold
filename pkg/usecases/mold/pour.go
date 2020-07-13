package mold

//go:generate mockgen -source=$GOFILE -destination=./mock_${GOFILE} -package=$GOPACKAGE

import (
	"github.com/ArthurKC/foundry/pkg/domains/mold"
)

type PourInteractor struct {
	repo            mold.Repository
	materialService mold.MaterialService
	productService  mold.ProductService
	output          PourOutputPort
}

func NewPourInteractor(repo mold.Repository, materialService mold.MaterialService, productService mold.ProductService, output PourOutputPort) *PourInteractor {
	return &PourInteractor{
		repo,
		materialService,
		productService,
		output,
	}
}

func (r *PourInteractor) ExecutePour(req *PourRequest) {
	mold, err := r.repo.FindByName(req.MoldName)
	if err != nil {
		r.output.RenderPourError(&PourError{req, err})
		return
	}

	m, err := r.materialService.GetMaterial(mold.Requirements())
	if err != nil {
		r.output.RenderPourError(&PourError{req, err})
		return
	}

	p, err := mold.Pour(req.DestDir, m)
	if err != nil {
		r.output.RenderPourError(&PourError{req, err})
		return
	}

	err = r.productService.SaveProduct(req.DestDir, p)
	if err != nil {
		r.output.RenderPourError(&PourError{req, err})
		return
	}

	r.output.RenderPourResponse(&PourResponse{MoldName: req.MoldName})
}

type PourOutputPort interface {
	RenderPourResponse(r *PourResponse)
	RenderPourError(e *PourError)
}

type PourRequest struct {
	MoldName string
	DestDir  string
}

type PourResponse struct {
	MoldName string
}

type PourError struct {
	Req *PourRequest
	Err error
}
