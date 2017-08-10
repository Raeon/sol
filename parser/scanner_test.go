package parser

import "testing"

func TestRemaining(t *testing.T) {

	tests := []struct {
		input    string
		expected int
	}{
		{"test", 4},
		{"example", 7},
		{"a much longer example", 21},
	}

	for i, tt := range tests {
		s := NewScanner(tt.input)
		if s.Remaining() != tt.expected {
			t.Fatalf("TestRemaining[%d]: got=%d expected=%d",
				i, s.Remaining(), tt.expected)
		}
	}
}

func TestRemainder(t *testing.T) {
	tests := []struct {
		input    string
		index    int
		expected string
	}{
		{"input string", 6, "string"},
	}

	for i, tt := range tests {
		s := NewScanner(tt.input)
		s.index = tt.index
		if s.Remainder() != tt.expected {
			t.Fatalf("TestRemainder[%d]: input=%d index=%d got=%s expected=%s",
				i, tt.index, tt.index, s.Remainder(), tt.expected)
		}
	}
}

func TestForward(t *testing.T) {
	tests := []struct {
		input  string
		before int
		amount int
		after  int
	}{
		{"this is a string", 0, 5, 5},
		{"short", 0, 6, 5},
	}

	for i, tt := range tests {
		s := NewScanner(tt.input)
		s.index = tt.before
		s.Forward(tt.amount)
		if s.index != tt.after {
			t.Fatalf("TestForward[%d]: expected=%d got=%d",
				i, tt.after, s.index)
		}
	}
}

func TestBackward(t *testing.T) {
	tests := []struct {
		input  string
		before int
		amount int
		after  int
	}{
		{"this is a string", 5, 5, 0},
		{"short", 5, 5, 0},
		{"short", 5, 10, 0},
	}

	for i, tt := range tests {
		s := NewScanner(tt.input)
		s.index = tt.before
		s.Backward(tt.amount)
		if s.index != tt.after {
			t.Fatalf("TestBackward[%d]: expected=%d got=%d",
				i, tt.after, s.index)
		}
	}
}

func TestPushPop(t *testing.T) {
	tests := []struct {
		push     bool
		length   int
		result   bool
		indexGet int
		indexSet int
	}{
		{true, 1, true, 0, 10},
		{true, 2, true, 10, 20},
		{true, 3, true, 20, 30},
		{false, 2, true, 30, 0},
		{false, 1, true, 20, 0},
		{false, 0, true, 10, 0},
		{false, 0, false, 0, 0},
	}

	s := NewScanner("")
	for i, tt := range tests {

		if s.index != tt.indexGet {
			t.Fatalf("TestPushPop[%d] get expected=%d got=%d",
				i, tt.indexGet, s.index)
		}

		var result bool
		if tt.push {
			result = s.Push()
			s.index = tt.indexSet
		} else {
			result = s.Pop()
		}
		len := s.stackSize()

		if tt.result != result {
			t.Fatalf("TestPushPop[%d]: result expected=%v got=%v",
				i, tt.result, result)
		}
		if len != tt.length {
			t.Fatalf("TestPushPop[%d]: len expected=%d got=%d",
				i, tt.length, len)
		}
	}
}

func TestMatchString(t *testing.T) {
	tests := []struct {
		input    string
		match    string
		expected int
	}{
		{"this is a string", "this", 4},
		{"this is a string", "test", 0},
		{"", "this", 0},
		{"this", "", 0},
	}

	for i, tt := range tests {
		s := NewScanner(tt.input)
		result := s.MatchString(tt.match)
		if result != tt.expected {
			t.Fatalf("TestMatchString[%d]: input=%s match=%s got=%d expected=%d",
				i, tt.input, tt.match, result, tt.expected)
		}
	}
}

func TestConsumeString(t *testing.T) {
	tests := []struct {
		input     string
		first     string
		second    string
		remainder string
	}{
		{"this is a string", "this ", "is ", "a string"},
		{"this is a string", "nomatch", "this ", "is a string"},
	}

	for i, tt := range tests {
		s := NewScanner(tt.input)
		s.ConsumeString(tt.first)
		s.ConsumeString(tt.second)
		remainder := s.Remainder()
		if remainder != tt.remainder {
			t.Fatalf("TestConsumeString[%d]: got=%s expected=%s",
				i, remainder, tt.remainder)
		}
	}
}

func TestMatchRegex(t *testing.T) {
	tests := []struct {
		input    string
		expr     string
		index    int
		expected int
	}{
		{"this is a string", "(this)", 0, 4},
		{"this is a string", "else", 0, 0},
		{"this is a string", ".{0,4} is", 0, 7},
	}

	for i, tt := range tests {
		s := NewScanner(tt.input)
		s.index = tt.index
		result := s.MatchRegex(tt.expr)
		if result != tt.expected {
			t.Fatalf("TestMatchRegex[%d]: input=%s expr=%s index=%d "+
				"got=%d expected=%d", i, tt.input, tt.expr, tt.index,
				result, tt.expected)
		}
	}
}

func TestLinePosition(t *testing.T) {
	tests := []struct {
		input      string
		index      int
		lineNumber int
		lineIndex  int
	}{
		{"this\nis\na\nstring", 0, 1, 1}, // t
		{"this\nis\na\nstring", 1, 1, 2}, // h
		{"this\nis\na\nstring", 2, 1, 3}, // i
		{"this\nis\na\nstring", 3, 1, 4}, // s
		{"this\nis\na\nstring", 4, 1, 5}, // \n
		{"this\nis\na\nstring", 5, 2, 1}, // i
		{"this\nis\na\nstring", 8, 3, 1}, // a
		{"this\nis\na\nstring", 15, 4, 6},
	}

	for i, tt := range tests {
		s := NewScanner(tt.input)
		s.index = tt.index
		s.updateLinePosition()

		if s.lineNumber != tt.lineNumber {
			t.Fatalf("TestLinePosition[%d]: lineNumber expected=%d got=%d",
				i, tt.lineNumber, s.lineNumber)
		}
		if s.lineIndex != tt.lineIndex {
			t.Fatalf("TestLinePosition[%d]: lineIndex expected=%d got=%d",
				i, tt.lineIndex, s.lineIndex)
		}
	}
}
