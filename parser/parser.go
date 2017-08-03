package parser

type Parser interface {
	Parse(s *Scanner) (*ParseResult, error)
}

type ParseResult struct {
	Matched  string
	Name     string
	Children []*ParseResult
	Extra    interface{}
}

type ParseError struct {
	err string
}

func err(err string) *ParseError {
	return &ParseError{err: err}
}
func (pe *ParseError) Error() string {
	return pe.err
}
