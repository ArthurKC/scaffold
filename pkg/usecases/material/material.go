package material

//go:generate mockgen -source=$GOFILE -destination=./mock_${GOFILE} -package=$GOPACKAGE

import (
	"fmt"
	"log"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig"
	"github.com/huandu/xstrings"
)

type InputPort interface {
	Ask(p *Constituent) string
}

type OutputPort interface {
	Write(path string, content string)
	DestDir() string
}

type Constituent struct {
	Name        string
	Description string
}

type MoldSource interface {
	Parameters() []*Constituent
	Paths() []string
	Source(path string) string
}

type Material struct {
	tmpl MoldSource
	in   InputPort
	out  OutputPort
}

func New(tmpl MoldSource, in InputPort, out OutputPort) *Material {
	return &Material{tmpl, in, out}
}

func (g *Material) resolveParameters() map[string]string {
	paramNames := g.tmpl.Parameters()
	params := make(map[string]string, len(paramNames))
	for _, p := range paramNames {
		params[p.Name] = g.in.Ask(p)
	}
	return params
}

func (g *Material) executeMold(name string, tmpl string, params map[string]string) string {
	funcMap := make(template.FuncMap, 0)
	funcMap["firstRuneToLower"] = xstrings.FirstRuneToLower
	funcMap["firstRuneToUpper"] = xstrings.FirstRuneToUpper
	funcMap["DestDir"] = g.out.DestDir
	for k, v := range sprig.TxtFuncMap() {
		funcMap[k] = v
	}
	t, err := template.New(name).Funcs(funcMap).Parse(tmpl)
	if err != nil {
		log.Fatal("can not parse. %w", err)
	}
	var sb strings.Builder
	err = t.Execute(&sb, params)
	if err != nil {
		log.Fatal("can not execute mold. %w", err)
	}
	return sb.String()
}

func (g *Material) getOutPath(p string, params map[string]string) string {
	outPath := p
	if ext := filepath.Ext(p); strings.ToLower(ext) == ".gotmpl" {
		outPath = strings.TrimSuffix(outPath, ext)
	}
	return g.executeMold(p, outPath, params)
}

func (g *Material) Generate() {
	params := g.resolveParameters()
	for _, p := range g.tmpl.Paths() {
		outPath := g.getOutPath(p, params)
		content := g.executeMold(
			fmt.Sprintf("content of %s", p),
			g.tmpl.Source(p),
			params,
		)
		g.out.Write(outPath, content)
	}
}
