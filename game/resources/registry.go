package resources

type Registry[T any, I int] struct {
	Namespace map[string]I
	Elements  []T
	Names     []string
}

func (reg Registry[T, I]) Init(namespace map[string]I) *Registry[T, I] {
	reg.Namespace = namespace
	reg.Elements = make([]T, len(namespace))
	return &reg
}

func (reg *Registry[T, I]) Register(name string, element T) I {
	id, exists := reg.Namespace[name]
	if !exists {
		id = I(len(reg.Elements))
		reg.Namespace[name] = id
	}
	reg.Elements[id] = element
	reg.Names[id] = name
	return id
}
func (reg *Registry[T, I]) Set(id I, element T) {
	reg.Elements[id] = element
}

func (reg *Registry[T, I]) Get(id I) T {
	return reg.Elements[id]
}
func (reg *Registry[T, I]) GetName(id I) string {
	return reg.Names[id]
}
func (reg *Registry[T, I]) GetId(name string) I {
	return reg.Namespace[name]
}
