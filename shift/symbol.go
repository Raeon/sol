package shift

type Symbol interface {
	IsTerminal() bool
	IsEqual(other Symbol) bool
	ToString() string
}

type TokenSymbol struct {
	tokenType *TokenType
}

func NewTokenSymbol(tokenType *TokenType) Symbol {
	return &TokenSymbol{
		tokenType: tokenType,
	}
}

func (s *TokenSymbol) IsTerminal() bool {
	return true
}

func (s *TokenSymbol) IsEqual(other Symbol) bool {
	ts, ok := other.(*TokenSymbol)
	if !ok {
		return false
	}
	return s.tokenType == ts.tokenType
}

func (s *TokenSymbol) ToString() string {
	return s.tokenType.name
}

type ReferenceSymbol struct {
	rule *Rule
}

func NewReferenceSymbol(rule *Rule) Symbol {
	return &ReferenceSymbol{
		rule: rule,
	}
}

func (s *ReferenceSymbol) IsTerminal() bool {
	return false
}

func (s *ReferenceSymbol) IsEqual(other Symbol) bool {
	rs, ok := other.(*ReferenceSymbol)
	if !ok {
		return false
	}
	return s.rule == rs.rule
}

func (s *ReferenceSymbol) ToString() string {
	return s.rule.name
}
