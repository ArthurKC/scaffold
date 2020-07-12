package mold

type ProductComponent struct {
	path     string
	contents string
}

func (p *ProductComponent) Path() string {
	return p.path
}

func (p *ProductComponent) Contents() string {
	return p.contents
}
