package runtime

type Scope struct {
	parent   *Scope
	declared map[string]Object
}

func NewScope() *Scope {
	return &Scope{
		parent:   nil,
		declared: make(map[string]Object),
	}
}

func (s *Scope) IsLocal(identifier string) bool {
	_, ok := s.declared[identifier]
	return ok
}

func (s *Scope) Set(identifier string, value Object) Object {
	if s.IsLocal(identifier) {
		s.declared[identifier] = value
		return value
	}

	if s.parent != nil {
		return s.parent.Set(identifier, value)
	}

	s.declared[identifier] = value
	return value

	// return &Exception{
	// 	Message: fmt.Sprintf(
	// 		"cannot set undeclared variable: %s",
	// 		identifier,
	// 	),
	// }
}

func (s *Scope) SetLocal(identifier string, value Object) {
	s.declared[identifier] = value
}

func (s *Scope) Get(identifier string) Object {
	val, ok := s.declared[identifier]
	if ok {
		return val
	}

	if s.parent != nil {
		return s.parent.Get(identifier)
	}

	return &Nil{}
}
