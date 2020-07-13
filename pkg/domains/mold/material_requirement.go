package mold

type MaterialRequirement struct {
	name        string
	description string
}

func NewMaterialRequirement(name string, description string) *MaterialRequirement {
	return &MaterialRequirement{
		name:        name,
		description: description,
	}
}

func (m *MaterialRequirement) Name() string {
	return m.name
}

func (m *MaterialRequirement) Description() string {
	return m.description
}

func (m *MaterialRequirement) meets(parameters map[string]string) bool {
	if _, ok := parameters[m.name]; ok {
		return true
	}
	return false
}
