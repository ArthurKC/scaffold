package mold

type Product struct {
	components []*ProductComponent
}

func (p *Product) Components() []*ProductComponent {
	ret := make([]*ProductComponent, 0, len(p.components))
	return append(ret, p.components...)
}
