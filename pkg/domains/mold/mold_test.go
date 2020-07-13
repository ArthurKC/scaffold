package mold

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/go-cmp/cmp"
)

func TestMold_Pour(t *testing.T) {
	type fields struct {
		name         string
		components   []*Component
		requirements []*MaterialRequirement
	}
	type args struct {
		destDir string
	}
	type mockset struct {
		material Material
		mockset  func(material *MockMaterial)
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		mockset func(material *MockMaterial)
		want    *Product
		wantErr bool
	}{
		{
			name: "empty components",
			fields: fields{
				name:         "mold",
				components:   []*Component{},
				requirements: []*MaterialRequirement{},
			},
			args: args{
				destDir: "d/",
			},
			mockset: func(material *MockMaterial) {
				material.EXPECT().Parameters().Return(map[string]string{})
			},
			want: &Product{[]*ProductComponent{}},
		},
		{
			name: "no constituent mold",
			fields: fields{
				name: "mold",
				components: []*Component{
					{
						Path:     "a/b/c.yaml",
						Contents: "some content",
					},
				},
				requirements: []*MaterialRequirement{},
			},
			args: args{
				destDir: "d/",
			},
			mockset: func(material *MockMaterial) {
				material.EXPECT().Parameters().Return(map[string]string{})
			},
			want: &Product{
				[]*ProductComponent{
					{
						path:     "a/b/c.yaml",
						contents: "some content",
					},
				},
			},
		},
		{
			name: "single constituent content mold",
			fields: fields{
				name: "mold",
				components: []*Component{
					{
						Path:     "a/b/c.yaml",
						Contents: "user: {{.Name}}",
					},
				},
				requirements: []*MaterialRequirement{
					{
						name:        "Name",
						description: "D1",
					},
				},
			},
			args: args{
				destDir: "d/",
			},
			mockset: func(material *MockMaterial) {
				material.EXPECT().Parameters().Return(map[string]string{
					"Name": "testUser",
				})
			},
			want: &Product{
				[]*ProductComponent{
					{
						path:     "a/b/c.yaml",
						contents: "user: testUser",
					},
				},
			},
		},
		{
			name: "single constituent content and path mold",
			fields: fields{
				name: "mold",
				components: []*Component{
					{
						Path:     "a/b/{{.Name}}.yaml",
						Contents: "user: {{.Name}}",
					},
				},
				requirements: []*MaterialRequirement{
					{
						name:        "Name",
						description: "D1",
					},
				},
			},
			args: args{
				destDir: "d/",
			},
			mockset: func(material *MockMaterial) {
				material.EXPECT().Parameters().Return(map[string]string{
					"Name": "testUser",
				})
			},
			want: &Product{
				[]*ProductComponent{
					{
						path:     "a/b/testUser.yaml",
						contents: "user: testUser",
					},
				},
			},
		},
		{
			name: "multi parameters content and path mold",
			fields: fields{
				name: "mold",
				components: []*Component{
					{
						Path:     "a/{{.Name}}/{{.Score}}.yaml",
						Contents: "{{.Name}}: {{.Score}}",
					},
				},
				requirements: []*MaterialRequirement{
					{
						name:        "Name",
						description: "D1",
					},
					{
						name:        "Score",
						description: "D2",
					},
				},
			},
			args: args{
				destDir: "d/",
			},
			mockset: func(material *MockMaterial) {
				material.EXPECT().Parameters().Return(map[string]string{
					"Name":  "testUser",
					"Score": "100",
				})
			},
			want: &Product{
				[]*ProductComponent{
					{
						path:     "a/testUser/100.yaml",
						contents: "testUser: 100",
					},
				},
			},
		},
		{
			name: "multi parameters content and multi paths mold",
			fields: fields{
				name: "mold",
				components: []*Component{
					{
						Path:     "a/{{.Name}}/{{.Score}}.yaml",
						Contents: "{{.Name}}: {{.Score}}",
					},
					{
						Path:     "b/{{.Score}}/{{.Name}}.yaml",
						Contents: "{{.Score}}: {{.Name}}",
					},
				},
				requirements: []*MaterialRequirement{
					{
						name:        "Name",
						description: "D1",
					},
					{
						name:        "Score",
						description: "D2",
					},
				},
			},
			args: args{
				destDir: "d/",
			},
			mockset: func(material *MockMaterial) {
				material.EXPECT().Parameters().Return(map[string]string{
					"Name":  "testUser",
					"Score": "100",
				})
			},
			want: &Product{
				[]*ProductComponent{
					{
						path:     "a/testUser/100.yaml",
						contents: "testUser: 100",
					},
					{
						path:     "b/100/testUser.yaml",
						contents: "100: testUser",
					},
				},
			},
		},
		{
			name: "multi parameters content and multi paths mold",
			fields: fields{
				name: "mold",
				components: []*Component{
					{
						Path:     "a/{{.Name}}/{{.Score}}.yaml",
						Contents: "{{.Name}}: {{.Score}}",
					},
					{
						Path:     "b/{{.Score}}/{{.Name}}.yaml",
						Contents: "{{.Score}}: {{.Name}}",
					},
				},
				requirements: []*MaterialRequirement{
					{
						name:        "Name",
						description: "D1",
					},
					{
						name:        "Score",
						description: "D2",
					},
				},
			},
			args: args{
				destDir: "d/",
			},
			mockset: func(material *MockMaterial) {
				material.EXPECT().Parameters().Return(map[string]string{
					"Name":  "testUser",
					"Score": "100",
				})
			},
			want: &Product{
				[]*ProductComponent{
					{
						path:     "a/testUser/100.yaml",
						contents: "testUser: 100",
					},
					{
						path:     "b/100/testUser.yaml",
						contents: "100: testUser",
					},
				},
			},
		},
		{
			name: "mold functions",
			fields: fields{
				name: "mold",
				components: []*Component{
					{
						Path:     "users/{{snakecase .Name}}.yaml",
						Contents: "name: {{firstRuneToUpper .Name}}, dir: {{DestDir}}",
					},
				},
				requirements: []*MaterialRequirement{
					{
						name:        "Name",
						description: "D1",
					},
				},
			},
			args: args{
				destDir: "target/dir",
			},
			mockset: func(material *MockMaterial) {
				material.EXPECT().Parameters().Return(map[string]string{
					"Name": "testUser",
				})
			},
			want: &Product{
				[]*ProductComponent{
					{
						path:     "users/test_user.yaml",
						contents: "name: TestUser, dir: target/dir",
					},
				},
			},
		},
		{
			name: "invalid template path",
			fields: fields{
				name: "mold",
				components: []*Component{
					{
						Path:     "users/{{.Name .Name}}.yaml",
						Contents: "some content",
					},
				},
			},
			args: args{
				destDir: "target/dir",
			},
			mockset: func(material *MockMaterial) {
				material.EXPECT().Parameters().Return(map[string]string{
					"Name": "testUser",
				})
			},
			wantErr: true,
		},
		{
			name: "invalid template content",
			fields: fields{
				name: "mold",
				components: []*Component{
					{
						Path:     "users/{{.Name}}.yaml",
						Contents: "some {{.Name }",
					},
				},
			},
			args: args{
				destDir: "target/dir",
			},
			mockset: func(material *MockMaterial) {
				material.EXPECT().Parameters().Return(map[string]string{
					"Name": "testUser",
				})
			},
			wantErr: true,
		},
		{
			name: "not meet all requirements",
			fields: fields{
				name: "mold",
				components: []*Component{
					{
						Path:     "users/{{.Name}}.yaml",
						Contents: "score {{.Score}",
					},
				},
				requirements: []*MaterialRequirement{
					{
						name:        "Name",
						description: "D1",
					},
					{
						name:        "Score",
						description: "D2",
					},
				},
			},
			args: args{
				destDir: "target/dir",
			},
			mockset: func(material *MockMaterial) {
				material.EXPECT().Parameters().Return(map[string]string{
					"Name": "testUser",
				})
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			m := &Mold{
				name:         tt.fields.name,
				components:   tt.fields.components,
				requirements: tt.fields.requirements,
			}
			material := NewMockMaterial(ctrl)
			tt.mockset(material)
			got, err := m.Pour(tt.args.destDir, material)
			if err != nil != tt.wantErr {
				t.Errorf("Mold.Pour() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := cmp.Diff(tt.want, got, cmp.AllowUnexported(Product{}, ProductComponent{})); diff != "" {
				t.Errorf("Mold.Pour() is mismatch (-want, +got):\n%s", diff)
			}
		})
	}
}
