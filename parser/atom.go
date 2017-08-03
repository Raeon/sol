package parser

type AtomParser struct {
	atom string
}

func Atom(atom string) Parser {
	return &AtomParser{atom: atom}
}
func (ap *AtomParser) Parse(s *Scanner) (*ParseResult, error) {
	if s.ConsumeString(ap.atom) != ap.atom {
		return nil, err("")
	}
	return nil, nil
}
