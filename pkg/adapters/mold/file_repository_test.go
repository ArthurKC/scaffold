package mold

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/ArthurKC/foundry/pkg/domains/mold"
	"github.com/google/go-cmp/cmp"
)

func TestFileRepository_Save(t *testing.T) {
	type args struct {
		m *mold.Mold
	}
	tests := []struct {
		name    string
		f       *FileRepository
		args    args
		want    map[string]string
		wantErr bool
	}{
		{
			name: "empty mold",
			f:    &FileRepository{},
			args: args{
				m: mold.New(
					"testdata/saved",
					[]*mold.Component{},
					[]*mold.MaterialRequirement{},
				),
			},
			want: map[string]string{
				"testdata/saved/mold.yaml": "parameters: []\n",
			},
		},
		{
			name: "one component",
			f:    &FileRepository{},
			args: args{
				m: mold.New(
					"testdata/saved",
					[]*mold.Component{
						mold.NewComponent("a/b/c.go", "content"),
					},
					[]*mold.MaterialRequirement{},
				),
			},
			want: map[string]string{
				"testdata/saved/mold.yaml": "parameters: []\n",
				"testdata/saved/a/b/c.go":  "content",
			},
		},
		{
			name: "many components",
			f:    &FileRepository{},
			args: args{
				m: mold.New(
					"testdata/saved",
					[]*mold.Component{
						mold.NewComponent("a/b/c1.go", "content-c1"),
						mold.NewComponent("a/b/c2.go", "content-c2"),
						mold.NewComponent("a/d1.go", "content-d1"),
						mold.NewComponent("x/d2.go", "content-d2"),
					},
					[]*mold.MaterialRequirement{},
				),
			},
			want: map[string]string{
				"testdata/saved/mold.yaml": "parameters: []\n",
				"testdata/saved/a/b/c1.go": "content-c1",
				"testdata/saved/a/b/c2.go": "content-c2",
				"testdata/saved/a/d1.go":   "content-d1",
				"testdata/saved/x/d2.go":   "content-d2",
			},
		},
		{
			name: "mold.yaml",
			f:    &FileRepository{},
			args: args{
				m: mold.New(
					"testdata/saved",
					[]*mold.Component{
						mold.NewComponent("a/{{.Name}}/{{.Score}}.go", "{{.Name}}-{{.Score}}"),
					},
					[]*mold.MaterialRequirement{
						mold.NewMaterialRequirement("Name", "D1"),
						mold.NewMaterialRequirement("Score", "D2"),
					},
				),
			},
			want: map[string]string{
				"testdata/saved/mold.yaml":                 "parameters:\n- name: Name\n  description: D1\n- name: Score\n  description: D2\n",
				"testdata/saved/a/{{.Name}}/{{.Score}}.go": "{{.Name}}-{{.Score}}",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !strings.HasPrefix(tt.args.m.Name(), "testdata/") {
				t.Error("location which is put a mold must be started with testdasa/")
				return
			}

			defer os.RemoveAll(tt.args.m.Name())
			f := &FileRepository{}
			if err := f.Save(tt.args.m); (err != nil) != tt.wantErr {
				t.Errorf("FileRepository.Save() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			for p, contents := range tt.want {
				got, err := ioutil.ReadFile(p)
				if err != nil {
					t.Errorf("failed to verify FileRepository.Save() result, error = %v", err)
					return
				}
				if diff := cmp.Diff(contents, string(got)); diff != "" {
					t.Errorf("created file (=%s) by FileRepository.Save() mismatch (-want +got) = \n%s", p, diff)
					return
				}
			}
		})
	}
}

func TestFileRepository_FindByName(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		f       *FileRepository
		args    args
		want    *mold.Mold
		wantErr bool
	}{
		{
			name: "file not found",
			f:    &FileRepository{},
			args: args{
				name: "testdata/not_found_path",
			},
			wantErr: true,
		},
		{
			name: "empty mold",
			f:    &FileRepository{},
			args: args{
				name: "testdata/empty",
			},
			want: mold.New(
				"testdata/empty",
				[]*mold.Component{},
				[]*mold.MaterialRequirement{},
			),
		},
		{
			name: "sample mold",
			f:    &FileRepository{},
			args: args{
				name: "testdata/sample",
			},
			want: mold.New(
				"testdata/sample",
				[]*mold.Component{
					mold.NewComponent("dir1/f1_1.go", "f1_1:content"),
					mold.NewComponent("dir1/f1_2.go", "f1_2:content"),
					mold.NewComponent("f1.gotmpl", "f1:content"),
					mold.NewComponent("{{.Name}}.yaml", "{{.Name}}:content"),
				},
				[]*mold.MaterialRequirement{
					mold.NewMaterialRequirement("Name", "D1"),
				},
			),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &FileRepository{}
			got, err := f.FindByName(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("FileRepository.FindByName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := cmp.Diff(tt.want, got, cmp.AllowUnexported(mold.Mold{}, mold.Component{}, mold.MaterialRequirement{})); diff != "" {
				t.Errorf("FileRepository.FindByName() mismatch (-want +got) = \n%s", diff)
			}
		})
	}
}
