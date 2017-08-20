package shift

import (
	"fmt"
	"regexp"
)

type TokenType struct {
	name       string
	pattern    string
	precedence int
}

func NewTokenType(name, pattern string, precedence int) *TokenType {
	return &TokenType{
		name:       name,
		pattern:    pattern,
		precedence: precedence,
	}
}

func (tt *TokenType) NewToken(literal string, index, line, char int) *Token {
	return &Token{
		Literal: literal,
		Type:    tt,

		Index:      index,
		LineNumber: line,
		LineIndex:  index,
	}
}

type Token struct {
	Literal string
	Type    *TokenType

	Index      int
	LineNumber int
	LineIndex  int
}

func (t *Token) ToString() string {
	return fmt.Sprintf("%s(%s)", t.Type.name, t.Literal)
}

type Tokenizer struct {
	input string
	types []*TokenType

	index      int
	lineNumber int
	lineIndex  int

	current *Token
	next    *Token
}

func NewTokenizer(input string) *Tokenizer {
	t := &Tokenizer{
		input:      input,
		index:      0,
		lineNumber: 1,
		lineIndex:  1,
	}
	return t
}

func (t *Tokenizer) Register(tta *TokenType) *TokenType {

	// Type array has descending precedence.
	// When we find an item in the array
	// with a lower precendence than the current
	// type, we insert it! Or, if there is none,
	// we just append it.

	// Find index to insert type at
	for i, ttb := range t.types {
		if ttb.precedence < tta.precedence {
			// Insert at 'i'
			t.types = append(t.types, nil)
			copy(t.types[i+1:], t.types[i:])
			t.types[i] = tta
			return tta
		}
	}

	// Otherwise, append
	t.types = append(t.types, tta)
	return tta
}

func (t *Tokenizer) Next() *Token {
	t.fetch()
	return t.current
}

func (t *Tokenizer) Current() *Token {
	return t.current
}

func (t *Tokenizer) Peek() *Token {
	return t.next
}

func (t *Tokenizer) fetch() bool {

	// Shift existing tokens forward
	t.current = t.next
	t.next = nil

	// Skip any whitespace
	t.skipWhitespace()

	// Try to match token in order of
	// descending token type precedence.
	// This also takes into account the EOF token.
	for _, tt := range t.types {
		token := t.matchType(tt)
		if token != nil {
			t.next = token
			return true
		}
	}

	return false
}

func (t *Tokenizer) matchType(tt *TokenType) *Token {
	loc := t.matchPattern(tt.pattern)
	if loc == nil {
		return nil
	}

	match := t.input[t.index : t.index+loc[1]]
	tok := tt.NewToken(match, t.index,
		t.lineNumber, t.lineIndex)
	t.index += loc[1]
	return tok
}

func (t *Tokenizer) matchPattern(pattern string) []int {
	reg := regexp.MustCompile("^" + pattern)
	rem := t.Remainder()
	return reg.FindStringIndex(rem)
}

func (t *Tokenizer) skipWhitespace() bool {
	loc := t.matchPattern("[\r\n\t ]+")
	if loc != nil {
		t.index += loc[1]
	}
	// TODO: Update line number and index
	return loc != nil
}

func (t *Tokenizer) Remaining() int {
	return len(t.input) - t.index
}

func (t *Tokenizer) Remainder() string {
	return t.input[t.index:]
}

func (t *Tokenizer) EOF() bool {
	return t.index == len(t.input)
}
