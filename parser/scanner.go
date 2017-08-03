package parser

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
	reg := regexp.MustCompile(expr)
	matched := reg.FindString(s.Remainder())
	return len(matched)
}

func (s *Scanner) ConsumeRegex(expr string) string {
	length := s.MatchRegex(expr)
	return s.Forward(length)
}

func (s *Scanner) updateLinePosition() {

	before := s.input[0 : s.index+1] // +1 to include \n
	lineStart := strings.LastIndex(before, "\n")
	if lineStart == -1 {
		lineStart = 0
	}

	s.lineNumber = 1 + strings.Count(before, "\n")
	s.lineIndex = 1 + s.index - lineStart
}

func (s *Scanner) LineNumber() int {
	return s.lineNumber
}

func (s *Scanner) LineIndex() int {
	return s.lineIndex
}
