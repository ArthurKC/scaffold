package metal

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
	Constituents() []*Constituent
	Paths() []string
	Source(path string) string
}

type Metal struct {
	tmpl MoldSource
	in   InputPort
	out  OutputPort
}

func New(tmpl MoldSource, in InputPort, out OutputPort) *Metal {
	return &Metal{tmpl, in, out}
}

func (g *Metal) resolveConstituents() map[string]string {
	paramNames := g.tmpl.Constituents()
	params := make(map[string]string, len(paramNames))
	for _, p := range paramNames {
		params[p.Name] = g.in.Ask(p)
	}
	return params
}

func (g *Metal) executeMold(name string, tmpl string, params map[string]string) string {
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

func (g *Metal) getOutPath(p string, params map[string]string) string {
	outPath := p
	if ext := filepath.Ext(p); strings.ToLower(ext) == ".gotmpl" {
		outPath = strings.TrimSuffix(outPath, ext)
	}
	return g.executeMold(p, outPath, params)
}

func (g *Metal) Generate() {
	params := g.resolveConstituents()
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
