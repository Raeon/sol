package runtime

import (
	"sol/parser"
	"testing"
)

func TestEvaluate(t *testing.T) {

	tests := []struct {
		input  string
		output string
	}{
		// {"5 + 5", "10"},
		// {"5 * 5", "25"},
		// {"5 + 1 * 3", "8"},
		// {"1 + 3 * 5", "16"},
		// {"let x = 5", "5"},
		// {"let x = 5 x + 10", "15"},
		{"x = y = 5", "(x = (y = 5))"},
	}

	for i, tt := range tests {
		env := NewEnv()
		parser := parser.NewParser()

		node, err := parser.Parse(tt.input)
		if err != nil {
			t.Fatalf("TestEvaluate[%d]: error: %s",
				i, err.Error())
		}

		result := env.Evaluate(node)
		resultStr := result.ToString()
		if resultStr != tt.output {
			t.Fatalf("TestEvaluate[%d]: expected=\"%s\" got=\"%s\"",
				i, tt.output, resultStr)
		}
	}

}
