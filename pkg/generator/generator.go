package generator

import (
	"fmt"
	"log"
	"path/filepath"
	"strings"
	"text/template"
)

type InputPort interface {
	Ask(p *Parameter) string
}

type OutputPort interface {
	Write(path string, content string)
}

type Parameter struct {
	Name        string
	Description string
}

type TemplateSource interface {
	Params() []*Parameter
	Paths() []string
	Source(path string) string
}

type Generator struct {
	tmpl TemplateSource
	in   InputPort
	out  OutputPort
}

func New(tmpl TemplateSource, in InputPort, out OutputPort) *Generator {
	return &Generator{tmpl, in, out}
}

func (g *Generator) resolveParams() map[string]string {
	paramNames := g.tmpl.Params()
	params := make(map[string]string, len(paramNames))
	for _, p := range paramNames {
		params[p.Name] = g.in.Ask(p)
	}
	return params
}

func (g *Generator) executeTemplate(name string, tmpl string, params map[string]string) string {
	t, err := template.New(name).Parse(tmpl)
	if err != nil {
		log.Fatal("can not parse. %w", err)
	}
	var sb strings.Builder
	err = t.Execute(&sb, params)
	if err != nil {
		log.Fatal("can not execute template. %w", err)
	}
	return sb.String()
}

func (g *Generator) getOutPath(p string, params map[string]string) string {
	outPath := p
	if ext := filepath.Ext(p); strings.ToLower(ext) == ".gotmpl" {
		outPath = strings.TrimSuffix(outPath, ext)
	}
	return g.executeTemplate(p, outPath, params)
}

func (g *Generator) Generate() {
	params := g.resolveParams()
	for _, p := range g.tmpl.Paths() {
		outPath := g.getOutPath(p, params)
		content := g.executeTemplate(
			fmt.Sprintf("content of %s", p),
			g.tmpl.Source(p),
			params,
		)
		g.out.Write(outPath, content)
	}
}
