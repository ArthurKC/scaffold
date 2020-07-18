package mold

import (
	"errors"
	"testing"

	"github.com/ArthurKC/foundry/pkg/domains/mold"
	gomock "github.com/golang/mock/gomock"
)

func TestPourInteractor_ExecutePour(t *testing.T) {
	type args struct {
		req *PourRequest
	}
	tests := []struct {
		name    string
		mockset func(
			repo *mold.MockRepository,
			materialService *mold.MockMaterialService,
			productService *mold.MockProductService,
			hout *MockPourOutputPort,
			genMaterial func(setup func(m *mold.MockMaterial)) *mold.MockMaterial,
		)
		args args
	}{
		{
			name: "not found",
			mockset: func(
				repo *mold.MockRepository,
				materialService *mold.MockMaterialService,
				productService *mold.MockProductService,
				out *MockPourOutputPort,
				genMaterial func(setup func(m *mold.MockMaterial)) *mold.MockMaterial,
			) {
				repo.EXPECT().FindByName("mold/dir").Return(nil, errors.New("not found"))
				out.EXPECT().RenderPourError(
					&PourError{
						&PourRequest{
							DestDir:  "dest/dir",
							MoldName: "mold/dir",
						},
						errors.New("not found"),
					},
				)
			},
			args: args{
				&PourRequest{
					DestDir:  "dest/dir",
					MoldName: "mold/dir",
				},
			},
		},
		{
			name: "no material",
			mockset: func(
				repo *mold.MockRepository,
				materialService *mold.MockMaterialService,
				productService *mold.MockProductService,
				out *MockPourOutputPort,
				genMaterial func(setup func(m *mold.MockMaterial)) *mold.MockMaterial,
			) {
				repo.EXPECT().FindByName("mold/dir").Return(
					mold.New(
						"mold/dir",
						[]*mold.Component{
							mold.NewComponent("a/b/c.yaml", "{{.Name}}"),
						},
						[]*mold.MaterialRequirement{
							mold.NewMaterialRequirement("Name", "D1"),
						},
					),
					nil,
				)
				materialService.EXPECT().GetMaterial([]*mold.MaterialRequirement{
					mold.NewMaterialRequirement("Name", "D1"),
				}).Return(nil, errors.New("no material"))
				out.EXPECT().RenderPourError(
					&PourError{
						&PourRequest{
							DestDir:  "dest/dir",
							MoldName: "mold/dir",
						},
						errors.New("no material"),
					},
				)
			},
			args: args{
				&PourRequest{
					DestDir:  "dest/dir",
					MoldName: "mold/dir",
				},
			},
		},
		{
			name: "cannot pour",
			mockset: func(
				repo *mold.MockRepository,
				materialService *mold.MockMaterialService,
				productService *mold.MockProductService,
				out *MockPourOutputPort,
				genMaterial func(setup func(m *mold.MockMaterial)) *mold.MockMaterial,
			) {
				repo.EXPECT().FindByName("mold/dir").Return(
					mold.New(
						"mold/dir",
						[]*mold.Component{
							// template error
							mold.NewComponent("a/b/c.yaml", "{{.Name}"),
						},
						[]*mold.MaterialRequirement{
							mold.NewMaterialRequirement("Name", "D1"),
						},
					),
					nil,
				)
				materialService.EXPECT().GetMaterial([]*mold.MaterialRequirement{
					mold.NewMaterialRequirement("Name", "D1"),
				}).Return(genMaterial(func(m *mold.MockMaterial) {
					m.EXPECT().Parameters().Return(map[string]string{
						"Name": "testUser",
					})
				}), nil)
				out.EXPECT().RenderPourError(
					gomock.Not(gomock.Nil()),
				)
			},
			args: args{
				&PourRequest{
					DestDir:  "dest/dir",
					MoldName: "mold/dir",
				},
			},
		},
		{
			name: "cannot save",
			mockset: func(
				repo *mold.MockRepository,
				materialService *mold.MockMaterialService,
				productService *mold.MockProductService,
				out *MockPourOutputPort,
				genMaterial func(setup func(m *mold.MockMaterial)) *mold.MockMaterial,
			) {
				repo.EXPECT().FindByName("mold/dir").Return(
					mold.New(
						"mold/dir",
						[]*mold.Component{
							// template error
							mold.NewComponent("a/b/c.yaml", "{{.Name}}"),
						},
						[]*mold.MaterialRequirement{
							mold.NewMaterialRequirement("Name", "D1"),
						},
					),
					nil,
				)
				materialService.EXPECT().GetMaterial([]*mold.MaterialRequirement{
					mold.NewMaterialRequirement("Name", "D1"),
				}).Return(genMaterial(func(m *mold.MockMaterial) {
					m.EXPECT().Parameters().Return(map[string]string{
						"Name": "testUser",
					})
				}), nil)
				productService.EXPECT().SaveProduct("dest/dir", gomock.Any()).Return(
					errors.New("cannot save"),
				)
				out.EXPECT().RenderPourError(
					&PourError{
						&PourRequest{
							DestDir:  "dest/dir",
							MoldName: "mold/dir",
						},
						errors.New("cannot save"),
					},
				)
			},
			args: args{
				&PourRequest{
					DestDir:  "dest/dir",
					MoldName: "mold/dir",
				},
			},
		},
		{
			name: "success",
			mockset: func(
				repo *mold.MockRepository,
				materialService *mold.MockMaterialService,
				productService *mold.MockProductService,
				out *MockPourOutputPort,
				genMaterial func(setup func(m *mold.MockMaterial)) *mold.MockMaterial,
			) {
				repo.EXPECT().FindByName("mold/dir").Return(
					mold.New(
						"mold/dir",
						[]*mold.Component{
							// template error
							mold.NewComponent("a/b/c.yaml", "{{.Name}}"),
						},
						[]*mold.MaterialRequirement{
							mold.NewMaterialRequirement("Name", "D1"),
						},
					),
					nil,
				)
				materialService.EXPECT().GetMaterial([]*mold.MaterialRequirement{
					mold.NewMaterialRequirement("Name", "D1"),
				}).Return(genMaterial(func(m *mold.MockMaterial) {
					m.EXPECT().Parameters().Return(map[string]string{
						"Name": "testUser",
					})
				}), nil)
				productService.EXPECT().SaveProduct("dest/dir", gomock.Any()).Return(nil)
				out.EXPECT().RenderPourResponse(
					&PourResponse{
						"mold/dir",
					},
				)
			},
			args: args{
				&PourRequest{
					DestDir:  "dest/dir",
					MoldName: "mold/dir",
				},
			},
		},
		// service.EXPECT().GetMaterial([]*mold.MaterialRequirement{
		// 	mold.NewMaterialRequirement("Name", "D1"),
		// }).Return(
		// 	genMaterial(func(m *mold.MockMaterial) {
		// 		m.EXPECT().Parameters().Return(map[string]string{
		// 			"Name": "testUser",
		// 		})
		// 	}),
		// )
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			repo := mold.NewMockRepository(ctrl)
			materialService := mold.NewMockMaterialService(ctrl)
			productService := mold.NewMockProductService(ctrl)
			out := NewMockPourOutputPort(ctrl)
			tt.mockset(repo, materialService, productService, out, func(setup func(m *mold.MockMaterial)) *mold.MockMaterial {
				ret := mold.NewMockMaterial(ctrl)
				setup(ret)
				return ret
			})
			r := &PourInteractor{
				repo:            repo,
				materialService: materialService,
				productService:  productService,
				output:          out,
			}
			r.ExecutePour(tt.args.req)
		})
	}
}
