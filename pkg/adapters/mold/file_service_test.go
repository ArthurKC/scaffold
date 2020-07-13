package mold

import (
	"testing"

	"github.com/ArthurKC/foundry/pkg/domains/mold"

	"github.com/google/go-cmp/cmp"
)

func TestFileService_ImportFrom(t *testing.T) {
	type args struct {
		path    string
		dstName string
	}
	tests := []struct {
		name    string
		f       *FileService
		args    args
		want    *mold.Mold
		wantErr bool
	}{
		{
			name: "not found file",
			f:    &FileService{},
			args: args{
				path:    "testdata/not_found_path",
				dstName: "dest/dir",
			},
			wantErr: true,
		},
		{
			name: "import from not mold",
			f:    &FileService{},
			args: args{
				path:    "testdata/import",
				dstName: "dest/dir",
			},
			want: mold.New(
				"dest/dir",
				[]*mold.Component{
					mold.NewComponent("dir/f1.go.gotmpl", "package dir\n"),
					mold.NewComponent("dir/f2.gotmpl.gotmpl", "package testf2\n"),
					mold.NewComponent("f1.go.gotmpl", "package testf1\n"),
				},
				[]*mold.MaterialRequirement{},
			),
		},
		{
			name: "import from already mold",
			f:    &FileService{},
			args: args{
				path:    "testdata/sample",
				dstName: "dest/dir",
			},
			want: mold.New(
				"dest/dir",
				[]*mold.Component{
					mold.NewComponent("dir1/f1_1.go.GOtmpl", "f1_1:content"),
					mold.NewComponent("dir1/f1_2.go.gotmpl", "f1_2:content"),
					mold.NewComponent("dir2/.gitkeep", ""),
					mold.NewComponent("f1.gotmpl.gotmpl", "f1:content"),
					mold.NewComponent("f2.noext", "f2:content"),
					mold.NewComponent("{{.Name}}.yaml.gotmpl", "{{.Name}}:content"),
				},
				[]*mold.MaterialRequirement{
					mold.NewMaterialRequirement("Name", "D1"),
				},
			),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &FileService{}
			got, err := f.ImportFrom(tt.args.path, tt.args.dstName)
			if (err != nil) != tt.wantErr {
				t.Errorf("FileService.ImportFrom() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := cmp.Diff(tt.want, got, cmp.AllowUnexported(mold.Mold{}, mold.MaterialRequirement{})); diff != "" {
				t.Errorf("FileService.ImportFrom() mismatch (-want +got) = \n%s", diff)
			}
		})
	}
}
