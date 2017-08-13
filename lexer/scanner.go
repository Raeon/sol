package lexer

import (
	"regexp"
	"strings"
)

type Scanner struct {
	input string
	index int

	lineNumber int
	lineIndex  int

	stack *scannerState
}

type scannerState struct {
	index int

	lineNumber int
	lineIndex  int

	next *scannerState
}

func NewScanner(input string) *Scanner {
	return &Scanner{
		input: input,
		index: 0,

		lineNumber: 0,
		lineIndex:  0,
	}
}

func (s *Scanner) Remaining() int {
	return len(s.input) - s.index
}

func (s *Scanner) Remainder() string {
	return s.input[s.index:]
}

func (s *Scanner) Forward(amount int) string {
	before := s.index
	s.index += amount
	if s.index > len(s.input) {
		s.index = len(s.input)
	}
	s.updateLinePosition()
	return s.input[before:s.index]
}

func (s *Scanner) Backward(amount int) string {
	before := s.index
	s.index -= amount
	if s.index < 0 {
		s.index = 0
	}
	s.updateLinePosition()
	return s.input[s.index:before]
}

func (s *Scanner) Push() bool {
	cur := &scannerState{
		index:      s.index,
		lineNumber: s.lineNumber,
		lineIndex:  s.lineIndex,
		next:       s.stack,
	}
	s.stack = cur
	return true
}

func (s *Scanner) Pop() bool {
	if s.stack != nil {
		s.index = s.stack.index
		s.lineNumber = s.stack.lineNumber
		s.lineIndex = s.stack.lineIndex
		s.stack = s.stack.next
		return true
	}
	return false
}

func (s *Scanner) Discard() {
	if s.stack != nil {
		s.stack = s.stack.next
	}
}

func (s *Scanner) stackSize() int {
	cur := s.stack
	count := 0

	for cur != nil {
		count++
		cur = cur.next
	}

	return count
}

func (s *Scanner) MatchString(str string) int {
	if s.Remaining() < len(str) {
		return 0
	}

	if s.input[s.index:s.index+len(str)] == str {
		return len(str)
	}
	return 0
}

func (s *Scanner) ConsumeString(str string) string {
	length := s.MatchString(str)
	return s.Forward(length)
}

func (s *Scanner) MatchRegex(expr string) int {
	reg := regexp.MustCompile("^" + expr)
	rem := s.Remainder()
	match := reg.FindString(rem)
	matchLen := len(match)
	if rem[:matchLen] != match {
		return 0
	}
	return matchLen
}

func (s *Scanner) ConsumeRegex(expr string) string {
	length := s.MatchRegex(expr)
	return s.Forward(length)
}

func (s *Scanner) updateLinePosition() {

	// target shouldnt be higher than len(s.input)
	target := s.index
	if target > len(s.input) {
		target = len(s.input)
	}

	// get everything up to the cursor (exclusive)
	before := s.input[0:target]

	// get the index where the current line started
	lineStart := strings.LastIndex(before, "\n") + 1
	if lineStart == -1 {
		lineStart = 0
	}

	// calculate line number and index,
	// both offset by +1 for human readability
	s.lineNumber = strings.Count(before, "\n") + 1
	s.lineIndex = s.index - lineStart + 1
}

func (s *Scanner) LineNumber() int {
	return s.lineNumber
}

func (s *Scanner) LineIndex() int {
	return s.lineIndex
}
