package generator

import (
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
)

type outWriteArgs struct {
	path    string
	content string
}

func TestNewGenerator(t *testing.T) {
	type args struct {
		tmpl TemplateSource
		in   InputPort
		out  OutputPort
	}
	tests := []struct {
		name string
		args args
		want *Generator
	}{
		{
			name: "correct case",
			args: args{
				tmpl: &MockTemplateSource{},
				in:   &MockInputPort{},
				out:  &MockOutputPort{},
			},
			want: &Generator{
				tmpl: &MockTemplateSource{},
				in:   &MockInputPort{},
				out:  &MockOutputPort{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(tt.args.tmpl, tt.args.in, tt.args.out); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewGenerator() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGenerator_Generate(t *testing.T) {
	tests := []struct {
		name    string
		mockset func(tmpl *MockTemplateSource, in *MockInputPort, out *MockOutputPort)
	}{
		{
			name: "empty template source",
			mockset: func(tmpl *MockTemplateSource, in *MockInputPort, out *MockOutputPort) {
				tmpl.EXPECT().Params().Return([]*Parameter{})
				tmpl.EXPECT().Paths().Return([]string{})
			},
		},
		{
			name: "no parameter template",
			mockset: func(tmpl *MockTemplateSource, in *MockInputPort, out *MockOutputPort) {
				tmpl.EXPECT().Params().Return([]*Parameter{})
				tmpl.EXPECT().Paths().Return([]string{
					"a/b/c.yaml.gotmpl",
				})
				tmpl.EXPECT().Source("a/b/c.yaml.gotmpl").Return("some content")
				out.EXPECT().Write("a/b/c.yaml", "some content")
			},
		},
		{
			name: "single parameter content template",
			mockset: func(tmpl *MockTemplateSource, in *MockInputPort, out *MockOutputPort) {
				tmpl.EXPECT().Params().Return([]*Parameter{
					{Name: "Name", Description: "Description"},
				})
				tmpl.EXPECT().Paths().Return([]string{
					"a/b/c.yaml.gotmpl",
				})
				tmpl.EXPECT().Source("a/b/c.yaml.gotmpl").Return("user: {{.Name}}")
				in.EXPECT().Ask(&Parameter{Name: "Name", Description: "Description"}).Return("testUser")
				out.EXPECT().Write("a/b/c.yaml", "user: testUser")
			},
		},
		{
			name: "single parameter content and path template",
			mockset: func(tmpl *MockTemplateSource, in *MockInputPort, out *MockOutputPort) {
				tmpl.EXPECT().Params().Return([]*Parameter{
					{Name: "Name", Description: "Description"},
				})
				tmpl.EXPECT().Paths().Return([]string{
					"a/b/{{.Name}}.yaml.gotmpl",
				})
				tmpl.EXPECT().Source("a/b/{{.Name}}.yaml.gotmpl").Return("user: {{.Name}}")
				in.EXPECT().Ask(&Parameter{Name: "Name", Description: "Description"}).Return("testUser")
				out.EXPECT().Write("a/b/testUser.yaml", "user: testUser")
			},
		},
		{
			name: "multi parameters content and path template",
			mockset: func(tmpl *MockTemplateSource, in *MockInputPort, out *MockOutputPort) {
				tmpl.EXPECT().Params().Return([]*Parameter{
					{Name: "Name", Description: "D1"},
					{Name: "Score", Description: "D2"},
				})
				tmpl.EXPECT().Paths().Return([]string{
					"a/{{.Name}}/{{.Score}}.yaml.gotmpl",
				})
				tmpl.EXPECT().Source("a/{{.Name}}/{{.Score}}.yaml.gotmpl").Return("{{.Name}}: {{.Score}}")
				in.EXPECT().Ask(&Parameter{Name: "Name", Description: "D1"}).Return("testUser")
				in.EXPECT().Ask(&Parameter{Name: "Score", Description: "D2"}).Return("100")
				out.EXPECT().Write("a/testUser/100.yaml", "testUser: 100")
			},
		},
		{
			name: "multi parameters content and multi paths template",
			mockset: func(tmpl *MockTemplateSource, in *MockInputPort, out *MockOutputPort) {
				tmpl.EXPECT().Params().Return([]*Parameter{
					{Name: "Name", Description: "D1"},
					{Name: "Score", Description: "D2"},
				})
				tmpl.EXPECT().Paths().Return([]string{
					"users/{{.Name}}.yaml.gotmpl",
					"scores/{{.Score}}.yaml.gotmpl",
				})
				tmpl.EXPECT().Source("users/{{.Name}}.yaml.gotmpl").Return("score: {{.Score}}")
				tmpl.EXPECT().Source("scores/{{.Score}}.yaml.gotmpl").Return("name: {{.Name}}")
				in.EXPECT().Ask(&Parameter{Name: "Name", Description: "D1"}).Return("testUser")
				in.EXPECT().Ask(&Parameter{Name: "Score", Description: "D2"}).Return("100")
				out.EXPECT().Write("users/testUser.yaml", "score: 100")
				out.EXPECT().Write("scores/100.yaml", "name: testUser")
			},
		},
		{
			name: "template functions",
			mockset: func(tmpl *MockTemplateSource, in *MockInputPort, out *MockOutputPort) {
				tmpl.EXPECT().Params().Return([]*Parameter{
					{Name: "Name", Description: "D1"},
				})
				tmpl.EXPECT().Paths().Return([]string{
					"users/{{snakecase .Name}}.yaml.gotmpl",
				})
				tmpl.EXPECT().Source("users/{{snakecase .Name}}.yaml.gotmpl").Return("name: {{firstRuneToUpper .Name}}, dir: {{DestDir}}")
				in.EXPECT().Ask(&Parameter{Name: "Name", Description: "D1"}).Return("testUser")
				out.EXPECT().Write("users/test_user.yaml", "name: TestUser, dir: target/dir")
				out.EXPECT().DestDir().Return("target/dir")
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			tmpl := NewMockTemplateSource(ctrl)
			in := NewMockInputPort(ctrl)
			out := NewMockOutputPort(ctrl)
			tt.mockset(tmpl, in, out)
			g := &Generator{
				tmpl: tmpl,
				in:   in,
				out:  out,
			}
			g.Generate()
		})
	}
}
