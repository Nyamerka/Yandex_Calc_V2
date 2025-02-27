package eval

import (
	"math/big"
	"testing"
)

func TestEval_ValidExpressions(t *testing.T) {
	tests := []struct {
		expression string
		expected   float64
	}{
		{"1+1", 2},
		{"2*2", 4},
		{"10/2", 5},
		{"3-1", 2},
		{"(1+2)*3", 9},
		{"(10-5)/2", 2.5},
		{"3*4+2", 14},
		{"10/5-2", 0},
		{"(3+4)*(5-2)", 21},
		{"100/10*5", 50},
		{"(1+2+3)*4", 24},
		{"10-(5-2)", 7},
		{"100/10/5", 2},
		{"(10+5)*2-3", 27},
		{"(10+5)*2/3", 10},
		{"(10+5)*2/3.0", 10},
		{"10+5*2/3", 13.3333333333},
		{"10+5*2/3.0", 13.3333333333},
		{"10+5*2/3.0+1", 14.3333333333},
		{"10+5*2/3.0-1", 12.3333333333},
		{"10+5*2/3.0*2", 16.6666666667},
		{"10+5*2/3.0*2.0", 16.6666666667},
		{"10+5*2/3.0*2.0+1", 17.6666666667},
		{"10+5*2/3.0*2.0-1", 15.6666666667},
	}
	for _, tt := range tests {
		result, err := Eval(tt.expression)
		if err != nil {
			t.Errorf("unexpected error for expression %q: %v", tt.expression, err)
		}
		actual := BigratToFloat(result)
		if actual != tt.expected {
			t.Errorf("for expression %q, expected %v, got %v", tt.expression, tt.expected, actual)
		}
	}
}

func TestEval_InvalidExpressions(t *testing.T) {
	tests := []string{
		"1//1",
		"abc",
		"1++",
		"))",
		"10/",
		"/10",
		"10*",
		"*10",
		"10+",
		"+10",
		"10-",
		"((10+5)",
		"10+5*2/",
		"10+5*2/",
		"10+5*2/3*",
		"10+5*2/3*(",
		"10+5*2/3*)",
		"10+5*2/3*(10+5",
		"10+5*2/3*(10+5)+",
		"10+5*2/3*(10+5)+-",
		"10+5*2/3*(10+5)++-10",
		"10+5*2/3*(10+5)+-10*",
		"10+5*2/3*(10+5)+-10*/",
		"10+5*2/3*(10+5)+-10*/10",
	}
	for _, expr := range tests {
		_, err := Eval(expr)
		if err == nil {
			t.Errorf("expected an error for expression %q, got nil", expr)
		}
	}
}

func TestBigratToFloat(t *testing.T) {
	tests := []struct {
		bigrat   *big.Rat
		expected float64
	}{
		{big.NewRat(10, 2), 5},
		{big.NewRat(1, 3), 0.3333333333},
	}
	for _, tt := range tests {
		actual := BigratToFloat(tt.bigrat)
		if actual != tt.expected {
			t.Errorf("expected %v, got %v", tt.expected, actual)
		}
	}
}

func TestBigratToInt(t *testing.T) {
	tests := []struct {
		bigrat   *big.Rat
		expected int64
	}{
		{big.NewRat(10, 2), 5},
		{big.NewRat(1, 3), 0},
		{big.NewRat(7, 2), 4},
		{big.NewRat(-5, 2), -3},
	}
	for _, tt := range tests {
		actual, err := BigratToInt(tt.bigrat)
		if err != nil {
			t.Errorf("unexpected error for big.Rat %v: %v", tt.bigrat, err)
		}
		if actual != tt.expected {
			t.Errorf("expected %d, got %d for big.Rat %v", tt.expected, actual, tt.bigrat)
		}
	}
}

func TestBigratToBigint(t *testing.T) {
	tests := []struct {
		bigrat   *big.Rat
		expected *big.Int
	}{
		{big.NewRat(10, 2), big.NewInt(5)},
		{big.NewRat(1, 3), big.NewInt(0)},
		{big.NewRat(7, 2), big.NewInt(4)},
		{big.NewRat(-5, 2), big.NewInt(-3)},
	}
	for _, tt := range tests {
		actual := BigratToBigint(tt.bigrat)
		if actual.Cmp(tt.expected) != 0 {
			t.Errorf("expected %v, got %v for big.Rat %v", tt.expected, actual, tt.bigrat)
		}
	}
}

func TestFloatToBigrat(t *testing.T) {
	tests := []struct {
		input    float64
		expected *big.Rat
	}{
		{5.0, big.NewRat(5, 1)},
		{0.3333333333, big.NewRat(3333333333, 10000000000)},
		{3.14, big.NewRat(314, 100)},
		{-2.5, big.NewRat(-25, 10)},
	}
	for _, tt := range tests {
		actual := FloatToBigrat(tt.input)
		if actual.Cmp(tt.expected) != 0 {
			t.Errorf("expected %v, got %v for float %v", tt.expected, actual, tt.input)
		}
	}
}
