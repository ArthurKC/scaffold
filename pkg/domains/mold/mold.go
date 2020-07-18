package mold

import "fmt"

// Mold is aggregation root.
type Mold struct {
	name         string
	components   []*Component
	requirements []*MaterialRequirement
}

// New creates new Mold instance.
func New(name string, components []*Component, requirements []*MaterialRequirement) *Mold {
	return &Mold{
		name:         name,
		components:   components,
		requirements: requirements,
	}
}

func (m *Mold) Name() string {
	return m.name
}

func (m *Mold) Components() []*Component {
	ret := make([]*Component, 0, len(m.components))
	return append(ret, m.components...)
}

func (m *Mold) Requirements() []*MaterialRequirement {
	ret := make([]*MaterialRequirement, 0, len(m.requirements))
	return append(ret, m.requirements...)
}

func (m *Mold) Pour(destDir string, material Material) (*Product, error) {
	cnts := material.Parameters()
	filtered := make(map[string]string, len(m.requirements))
	for _, r := range m.requirements {
		if !r.meets(cnts) {
			return nil, fmt.Errorf("a material doesn't meet the requirement. (Name: %s)", r.name)
		}
		filtered[r.name] = cnts[r.name]
	}

	ps := make([]*ProductComponent, 0, len(m.components))
	for _, c := range m.components {
		p, err := c.pour(destDir, filtered)
		if err != nil {
			return nil, err
		}
		ps = append(ps, p)
	}
	return &Product{ps}, nil
}
