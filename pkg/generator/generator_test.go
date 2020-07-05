package generator

import (
	"reflect"
	"testing"
)

type MockTemplateSource struct {
	ParamsReturn []string
	PathsReturn  []string
	SourceCalled []string
	SourceReturn map[string]string
}

func (m *MockTemplateSource) Params() []string {
	return m.ParamsReturn
}

func (m *MockTemplateSource) Paths() []string {
	return m.PathsReturn
}

func (m *MockTemplateSource) Source(path string) string {
	m.SourceCalled = append(m.SourceCalled, path)
	return m.SourceReturn[path]
}

type MockInputPort struct {
	AskCalled []string
	AskReturn map[string]string
}

func (m *MockInputPort) Ask(message string) string {
	m.AskCalled = append(m.AskCalled, message)
	return m.AskReturn[message]
}

type outWriteArgs struct {
	path    string
	content string
}

type MockOutputPort struct {
	WriteCalled []outWriteArgs
}

func (m *MockOutputPort) Write(path string, content string) {
	m.WriteCalled = append(m.WriteCalled, outWriteArgs{path, content})
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
	type fields struct {
		tmpl TemplateSource
		in   InputPort
		out  OutputPort
	}
	type want struct {
		tmplSourceCalled []string
		inAskCalled      []string
		outWriteCalled   []outWriteArgs
	}
	tests := []struct {
		name   string
		fields fields
		want   want
	}{
		{
			name: "empty template source",
			fields: fields{
				tmpl: &MockTemplateSource{
					ParamsReturn: []string{},
					PathsReturn:  []string{},
					SourceCalled: []string{},
					SourceReturn: map[string]string{},
				},
				in: &MockInputPort{
					AskCalled: []string{},
					AskReturn: map[string]string{},
				},
				out: &MockOutputPort{
					WriteCalled: []outWriteArgs{},
				},
			},
			want: want{
				tmplSourceCalled: []string{},
				inAskCalled:      []string{},
				outWriteCalled:   []outWriteArgs{},
			},
		},
		{
			name: "no parameter template",
			fields: fields{
				tmpl: &MockTemplateSource{
					ParamsReturn: []string{},
					PathsReturn:  []string{"a/b/c.yaml.gotmpl"},
					SourceCalled: []string{},
					SourceReturn: map[string]string{"a/b/c.yaml.gotmpl": "some content"},
				},
				in: &MockInputPort{
					AskCalled: []string{},
					AskReturn: map[string]string{},
				},
				out: &MockOutputPort{
					WriteCalled: []outWriteArgs{},
				},
			},
			want: want{
				tmplSourceCalled: []string{"a/b/c.yaml.gotmpl"},
				inAskCalled:      []string{},
				outWriteCalled: []outWriteArgs{
					{path: "a/b/c.yaml", content: "some content"},
				},
			},
		},
		{
			name: "single parameter content template",
			fields: fields{
				tmpl: &MockTemplateSource{
					ParamsReturn: []string{"Name"},
					PathsReturn:  []string{"a/b/c.yaml.gotmpl"},
					SourceCalled: []string{},
					SourceReturn: map[string]string{"a/b/c.yaml.gotmpl": "user: {{.Name}}"},
				},
				in: &MockInputPort{
					AskCalled: []string{},
					AskReturn: map[string]string{
						"Name": "testUser",
					},
				},
				out: &MockOutputPort{
					WriteCalled: []outWriteArgs{},
				},
			},
			want: want{
				tmplSourceCalled: []string{"a/b/c.yaml.gotmpl"},
				inAskCalled:      []string{"Name"},
				outWriteCalled: []outWriteArgs{
					{path: "a/b/c.yaml", content: "user: testUser"},
				},
			},
		},
		{
			name: "single parameter content and path template",
			fields: fields{
				tmpl: &MockTemplateSource{
					ParamsReturn: []string{"Name"},
					PathsReturn:  []string{"a/b/{{.Name}}.yaml.gotmpl"},
					SourceCalled: []string{},
					SourceReturn: map[string]string{"a/b/{{.Name}}.yaml.gotmpl": "name: {{.Name}}"},
				},
				in: &MockInputPort{
					AskCalled: []string{},
					AskReturn: map[string]string{
						"Name": "testUser",
					},
				},
				out: &MockOutputPort{
					WriteCalled: []outWriteArgs{},
				},
			},
			want: want{
				tmplSourceCalled: []string{"a/b/{{.Name}}.yaml.gotmpl"},
				inAskCalled:      []string{"Name"},
				outWriteCalled: []outWriteArgs{
					{path: "a/b/testUser.yaml", content: "name: testUser"},
				},
			},
		},
		{
			name: "multi parameters content and path template",
			fields: fields{
				tmpl: &MockTemplateSource{
					ParamsReturn: []string{"Name", "Score"},
					PathsReturn:  []string{"a/{{.Name}}/{{.Score}}.yaml.gotmpl"},
					SourceCalled: []string{},
					SourceReturn: map[string]string{"a/{{.Name}}/{{.Score}}.yaml.gotmpl": "{{.Name}}: {{.Score}}"},
				},
				in: &MockInputPort{
					AskCalled: []string{},
					AskReturn: map[string]string{
						"Name":  "testUser",
						"Score": "100",
					},
				},
				out: &MockOutputPort{
					WriteCalled: []outWriteArgs{},
				},
			},
			want: want{
				tmplSourceCalled: []string{"a/{{.Name}}/{{.Score}}.yaml.gotmpl"},
				inAskCalled:      []string{"Name", "Score"},
				outWriteCalled: []outWriteArgs{
					{path: "a/testUser/100.yaml", content: "testUser: 100"},
				},
			},
		},
		{
			name: "multi parameters content and multi paths template",
			fields: fields{
				tmpl: &MockTemplateSource{
					ParamsReturn: []string{"Name", "Score"},
					PathsReturn: []string{
						"users/{{.Name}}.yaml.gotmpl",
						"scores/{{.Score}}.yaml.gotmpl",
					},
					SourceCalled: []string{},
					SourceReturn: map[string]string{
						"users/{{.Name}}.yaml.gotmpl":   "score: {{.Score}}",
						"scores/{{.Score}}.yaml.gotmpl": "name: {{.Name}}",
					},
				},
				in: &MockInputPort{
					AskCalled: []string{},
					AskReturn: map[string]string{
						"Name":  "testUser",
						"Score": "100",
					},
				},
				out: &MockOutputPort{
					WriteCalled: []outWriteArgs{},
				},
			},
			want: want{
				tmplSourceCalled: []string{
					"users/{{.Name}}.yaml.gotmpl",
					"scores/{{.Score}}.yaml.gotmpl",
				},
				inAskCalled: []string{"Name", "Score"},
				outWriteCalled: []outWriteArgs{
					{path: "users/testUser.yaml", content: "score: 100"},
					{path: "scores/100.yaml", content: "name: testUser"},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &Generator{
				tmpl: tt.fields.tmpl,
				in:   tt.fields.in,
				out:  tt.fields.out,
			}
			g.Generate()
			if got := g.tmpl.(*MockTemplateSource).SourceCalled; !reflect.DeepEqual(got, tt.want.tmplSourceCalled) {
				t.Errorf("tmpl.Source() is called with args %v, want %v", got, tt.want.tmplSourceCalled)
			}
			if got := g.in.(*MockInputPort).AskCalled; !reflect.DeepEqual(got, tt.want.inAskCalled) {
				t.Errorf("in.Ask() is called with args %v, want %v", got, tt.want.inAskCalled)
			}
			if got := g.out.(*MockOutputPort).WriteCalled; !reflect.DeepEqual(got, tt.want.outWriteCalled) {
				t.Errorf("out.Write() is called with args %v, want %v", got, tt.want.outWriteCalled)
			}
		})
	}
}