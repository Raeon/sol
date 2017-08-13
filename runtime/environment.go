package runtime

type Environment struct {
	scope *Scope
}

type Scope struct {
	parent   *Scope
	declared map[string]*Object
}
