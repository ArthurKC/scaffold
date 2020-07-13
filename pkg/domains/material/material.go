package material

// Material is aggregation root.
type Material struct {
	name       string
	parameters map[string]string
}

// New creates new Material instance.
func New(name string, parameters map[string]string) *Material {
	return &Material{
		name:       name,
		parameters: parameters,
	}
}

func (m *Material) Name() string {
	return m.name
}

func (m *Material) Parameters() map[string]string {
	return m.parameters
}
