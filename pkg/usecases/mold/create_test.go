package mold

import (
	"errors"
	"testing"

	"github.com/ArthurKC/foundry/pkg/domains/mold"
	gomock "github.com/golang/mock/gomock"
)

func TestCreateInteractor_ExecuteCreate(t *testing.T) {
	type args struct {
		req *CreateRequest
	}
	tests := []struct {
		name    string
		mockset func(
			repo *mold.MockRepository,
			service *mold.MockService,
			out *MockCreateOutputPort,
		)
		args args
	}{
		{
			name: "empty source",
			mockset: func(repo *mold.MockRepository, service *mold.MockService, out *MockCreateOutputPort) {
				service.EXPECT().ImportFrom("import/path", "mold/path").Return(nil, errors.New("import error"))
				out.EXPECT().RenderCreateError(&CreateError{
					Req: &CreateRequest{
						ImportPath: "import/path",
						MoldName:   "mold/path",
					},
					Err: errors.New("import error"),
				})
			},
			args: args{
				req: &CreateRequest{
					ImportPath: "import/path",
					MoldName:   "mold/path",
				},
			},
		},
		{
			name: "can not save",
			mockset: func(repo *mold.MockRepository, service *mold.MockService, out *MockCreateOutputPort) {
				m := mold.New(
					"mold/path",
					[]*mold.Component{
						mold.NewComponent("a/file.go.gotmpl", "contentA"),
						mold.NewComponent("b/file.go.gotmpl", "contentB"),
					},
					[]*mold.MaterialRequirement{},
				)
				service.EXPECT().ImportFrom("import/path", "mold/path").Return(m, nil)
				out.EXPECT().RenderCreateError(&CreateError{
					Req: &CreateRequest{
						ImportPath: "import/path",
						MoldName:   "mold/path",
					},
					Err: errors.New("save error"),
				})
				repo.EXPECT().Save(m).Return(errors.New("save error"))
			},
			args: args{
				req: &CreateRequest{
					ImportPath: "import/path",
					MoldName:   "mold/path",
				},
			},
		},
		{
			name: "can import mold",
			mockset: func(repo *mold.MockRepository, service *mold.MockService, out *MockCreateOutputPort) {
				m := mold.New(
					"mold/path",
					[]*mold.Component{
						mold.NewComponent("a/file.go.gotmpl", "contentA"),
						mold.NewComponent("b/file.go.gotmpl", "contentB"),
					},
					[]*mold.MaterialRequirement{},
				)
				service.EXPECT().ImportFrom("import/path", "mold/path").Return(m, nil)
				out.EXPECT().RenderCreateResponse(&CreateResponse{MoldName: "mold/path"})
				repo.EXPECT().Save(m).Return(nil)
			},
			args: args{
				req: &CreateRequest{
					ImportPath: "import/path",
					MoldName:   "mold/path",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			repo := mold.NewMockRepository(ctrl)
			service := mold.NewMockService(ctrl)
			out := NewMockCreateOutputPort(ctrl)
			tt.mockset(repo, service, out)
			r := &CreateInteractor{
				repo:    repo,
				service: service,
				output:  out,
			}
			r.ExecuteCreate(tt.args.req)
		})
	}
}
