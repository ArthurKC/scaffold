package mold

//go:generate mockgen -source=$GOFILE -destination=./mock_${GOFILE} -package=$GOPACKAGE

import (
	"github.com/ArthurKC/foundry/pkg/domains/mold"
)

type CreateInteractor struct {
	repo    mold.Repository
	service mold.Service
	output  CreateOutputPort
}

func NewCreateInteractor(repo mold.Repository, service mold.Service, presenter CreateOutputPort) *CreateInteractor {
	return &CreateInteractor{
		repo,
		service,
		presenter,
	}
}

func (r *CreateInteractor) ExecuteCreate(req *CreateRequest) {
	m, err := r.service.ImportFrom(req.ImportPath, req.MoldName)
	if err != nil {
		r.output.RenderCreateError(&CreateError{req, err})
		return
	}

	if err := r.repo.Save(m); err != nil {
		r.output.RenderCreateError(&CreateError{req, err})
		return
	}

	r.output.RenderCreateResponse(&CreateResponse{req.MoldName})
}

type CreateOutputPort interface {
	RenderCreateResponse(r *CreateResponse)
	RenderCreateError(e *CreateError)
}

type CreateRequest struct {
	ImportPath string
	MoldName   string
}

type CreateResponse struct {
	MoldName string
}

type CreateError struct {
	Req *CreateRequest
	Err error
}
