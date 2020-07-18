package mold

import (
	"strings"
	"text/template"

	"github.com/Masterminds/sprig"
	"github.com/huandu/xstrings"
)

type Component struct {
	Path     string
	Contents string
}

func NewComponent(path string, contents string) *Component {
	return &Component{path, contents}
}

func (m *Component) solidify(name string, tmpl string, destDir string, params map[string]string) (string, error) {
	funcMap := make(template.FuncMap, 0)
	funcMap["firstRuneToLower"] = xstrings.FirstRuneToLower
	funcMap["firstRuneToUpper"] = xstrings.FirstRuneToUpper
	funcMap["DestDir"] = func() string { return destDir }
	for k, v := range sprig.TxtFuncMap() {
		funcMap[k] = v
	}
	t, err := template.New(name).Funcs(funcMap).Parse(tmpl)
	if err != nil {
		return "", err
	}
	var sb strings.Builder
	err = t.Execute(&sb, params)
	if err != nil {
		return "", err
	}
	return sb.String(), nil
}

func (c *Component) pour(destDir string, cnstituents map[string]string) (*ProductComponent, error) {
	path, err := c.solidify("path:"+c.Path, c.Path, destDir, cnstituents)
	if err != nil {
		return nil, err
	}
	contents, err := c.solidify("content:"+c.Path, c.Contents, destDir, cnstituents)
	if err != nil {
		return nil, err
	}
	return &ProductComponent{path: path, contents: contents}, nil
}
