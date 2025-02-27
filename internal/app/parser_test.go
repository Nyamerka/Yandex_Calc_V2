package app

import (
	"math"
	"testing"
)

func compareASTNodes(n1, n2 *ASTNode) bool {
	if n1 == nil && n2 == nil {
		return true
	}
	if n1 == nil || n2 == nil {
		return false
	}
	if n1.IsLeaf != n2.IsLeaf {
		return false
	}
	if n1.IsLeaf {
		return math.Abs(n1.Value-n2.Value) < 1e-9
	}
	if n1.Operator != n2.Operator {
		return false
	}
	return compareASTNodes(n1.Left, n2.Left) && compareASTNodes(n1.Right, n2.Right)
}

func createASTNode(isLeaf bool, value float64, operator string, left, right *ASTNode) *ASTNode {
	return &ASTNode{
		IsLeaf:   isLeaf,
		Value:    value,
		Operator: operator,
		Left:     left,
		Right:    right,
	}
}

func TestParseAST_ValidExpressions(t *testing.T) {
	tests := []struct {
		expression string
		expected   *ASTNode
	}{
		{"1+1", createASTNode(false, 0, "+", createASTNode(true, 1, "", nil, nil), createASTNode(true, 1, "", nil, nil))},
		{"2*2", createASTNode(false, 0, "*", createASTNode(true, 2, "", nil, nil), createASTNode(true, 2, "", nil, nil))},
		{"10/2", createASTNode(false, 0, "/", createASTNode(true, 10, "", nil, nil), createASTNode(true, 2, "", nil, nil))},
		{"3-1", createASTNode(false, 0, "-", createASTNode(true, 3, "", nil, nil), createASTNode(true, 1, "", nil, nil))},
		{"(1+2)*3", createASTNode(false, 0, "*", createASTNode(false, 0, "+", createASTNode(true, 1, "", nil, nil), createASTNode(true, 2, "", nil, nil)), createASTNode(true, 3, "", nil, nil))},
		{"(10-5)/2", createASTNode(false, 0, "/", createASTNode(false, 0, "-", createASTNode(true, 10, "", nil, nil), createASTNode(true, 5, "", nil, nil)), createASTNode(true, 2, "", nil, nil))},
		{"3*4+2", createASTNode(false, 0, "+", createASTNode(false, 0, "*", createASTNode(true, 3, "", nil, nil), createASTNode(true, 4, "", nil, nil)), createASTNode(true, 2, "", nil, nil))},
		{"10/5-2", createASTNode(false, 0, "-", createASTNode(false, 0, "/", createASTNode(true, 10, "", nil, nil), createASTNode(true, 5, "", nil, nil)), createASTNode(true, 2, "", nil, nil))},
		{"(3+4)*(5-2)", createASTNode(false, 0, "*", createASTNode(false, 0, "+", createASTNode(true, 3, "", nil, nil), createASTNode(true, 4, "", nil, nil)), createASTNode(false, 0, "-", createASTNode(true, 5, "", nil, nil), createASTNode(true, 2, "", nil, nil)))},
		{"100/10*5", createASTNode(false, 0, "*", createASTNode(false, 0, "/", createASTNode(true, 100, "", nil, nil), createASTNode(true, 10, "", nil, nil)), createASTNode(true, 5, "", nil, nil))},
		{"(1+2+3)*4", createASTNode(false, 0, "*", createASTNode(false, 0, "+", createASTNode(false, 0, "+", createASTNode(true, 1, "", nil, nil), createASTNode(true, 2, "", nil, nil)), createASTNode(true, 3, "", nil, nil)), createASTNode(true, 4, "", nil, nil))},
		{"10-(5-2)", createASTNode(false, 0, "-", createASTNode(true, 10, "", nil, nil), createASTNode(false, 0, "-", createASTNode(true, 5, "", nil, nil), createASTNode(true, 2, "", nil, nil)))},
		{"100/10/5", createASTNode(false, 0, "/", createASTNode(false, 0, "/", createASTNode(true, 100, "", nil, nil), createASTNode(true, 10, "", nil, nil)), createASTNode(true, 5, "", nil, nil))},
		{"(10+5)*2-3", createASTNode(false, 0, "-", createASTNode(false, 0, "*", createASTNode(false, 0, "+", createASTNode(true, 10, "", nil, nil), createASTNode(true, 5, "", nil, nil)), createASTNode(true, 2, "", nil, nil)), createASTNode(true, 3, "", nil, nil))},
		{"(10+5)*2/3", createASTNode(false, 0, "/", createASTNode(false, 0, "*", createASTNode(false, 0, "+", createASTNode(true, 10, "", nil, nil), createASTNode(true, 5, "", nil, nil)), createASTNode(true, 2, "", nil, nil)), createASTNode(true, 3, "", nil, nil))},
		{"(10+5)*2/3.0", createASTNode(false, 0, "/", createASTNode(false, 0, "*", createASTNode(false, 0, "+", createASTNode(true, 10, "", nil, nil), createASTNode(true, 5, "", nil, nil)), createASTNode(true, 2, "", nil, nil)), createASTNode(true, 3, "", nil, nil))},
		{"10+5*2/3", createASTNode(false, 0, "+", createASTNode(true, 10, "", nil, nil), createASTNode(false, 0, "/", createASTNode(false, 0, "*", createASTNode(true, 5, "", nil, nil), createASTNode(true, 2, "", nil, nil)), createASTNode(true, 3, "", nil, nil)))},
		{"10+5*2/3.0", createASTNode(false, 0, "+", createASTNode(true, 10, "", nil, nil), createASTNode(false, 0, "/", createASTNode(false, 0, "*", createASTNode(true, 5, "", nil, nil), createASTNode(true, 2, "", nil, nil)), createASTNode(true, 3, "", nil, nil)))},
		{"10+5*2/3.0+1", createASTNode(false, 0, "+", createASTNode(false, 0, "+", createASTNode(true, 10, "", nil, nil), createASTNode(false, 0, "/", createASTNode(false, 0, "*", createASTNode(true, 5, "", nil, nil), createASTNode(true, 2, "", nil, nil)), createASTNode(true, 3, "", nil, nil))), createASTNode(true, 1, "", nil, nil))},
		{"10+5*2/3.0-1", createASTNode(false, 0, "-", createASTNode(false, 0, "+", createASTNode(true, 10, "", nil, nil), createASTNode(false, 0, "/", createASTNode(false, 0, "*", createASTNode(true, 5, "", nil, nil), createASTNode(true, 2, "", nil, nil)), createASTNode(true, 3, "", nil, nil))), createASTNode(true, 1, "", nil, nil))},
		{"+1", createASTNode(true, 1, "", nil, nil)},
		{"-2", createASTNode(true, -2, "", nil, nil)},
		{"3.5", createASTNode(true, 3.5, "", nil, nil)},
		{"-3.5", createASTNode(true, -3.5, "", nil, nil)},
		{"+4.5", createASTNode(true, 4.5, "", nil, nil)},
		{"(-1+2)*3", createASTNode(false, 0, "*", createASTNode(false, 0, "+", createASTNode(true, -1, "", nil, nil), createASTNode(true, 2, "", nil, nil)), createASTNode(true, 3, "", nil, nil))},
		{"10+(-5+2)*3", createASTNode(false, 0, "+", createASTNode(true, 10, "", nil, nil), createASTNode(false, 0, "*", createASTNode(false, 0, "+", createASTNode(true, -5, "", nil, nil), createASTNode(true, 2, "", nil, nil)), createASTNode(true, 3, "", nil, nil)))},
		{"10+(5-2)*3", createASTNode(false, 0, "+", createASTNode(true, 10, "", nil, nil), createASTNode(false, 0, "*", createASTNode(false, 0, "-", createASTNode(true, 5, "", nil, nil), createASTNode(true, 2, "", nil, nil)), createASTNode(true, 3, "", nil, nil)))},
		{"10+5*(-2+3)", createASTNode(false, 0, "+", createASTNode(true, 10, "", nil, nil), createASTNode(false, 0, "*", createASTNode(true, 5, "", nil, nil), createASTNode(false, 0, "+", createASTNode(true, -2, "", nil, nil), createASTNode(true, 3, "", nil, nil))))},
		{"10+5*(2-3)", createASTNode(false, 0, "+", createASTNode(true, 10, "", nil, nil), createASTNode(false, 0, "*", createASTNode(true, 5, "", nil, nil), createASTNode(false, 0, "-", createASTNode(true, 2, "", nil, nil), createASTNode(true, 3, "", nil, nil))))},
	}

	for _, tt := range tests {
		node, err := ParseAST(tt.expression)
		if err != nil {
			t.Errorf("unexpected error for expression %q: %v", tt.expression, err)
		}
		if !compareASTNodes(node, tt.expected) {
			t.Errorf("for expression %q, AST mismatch", tt.expression)
		}
	}
}

func TestParseAST_InvalidExpressions(t *testing.T) {
	tests := []struct {
		expression string
	}{
		{"1//1"},
		{"abc"},
		{"1++"},
		{""},
		{"10/"},
		{"/10"},
		{"10*"},
		{"*10"},
		{"10+"},
		{"10-"},
		{"(10+5"},
		{"10+5)"},
		{"10+5*2/"},
		{"10+5*2/"},
		{"10+5*2/3*"},
		{"10+5*2/3*("},
		{"10+5*2/3*)"},
		{"10+5*2/3*(10+5"},
		{"10+5*2/3*(10+5))"},
		{"10+5*2/3*(10+5)+"},
		{"10+5*2/3*(10+5)+-"},
		{"10+5*2/3*(10+5)+-10*"},
		{"10+5*2/3*(10+5)+-10*/"},
		{"10+5*2/3*(10+5)+-10*/10"},
	}

	for _, tt := range tests {
		_, err := ParseAST(tt.expression)
		if err == nil {
			t.Errorf("expected an error for expression %q, got nil", tt.expression)
		}
	}
}
