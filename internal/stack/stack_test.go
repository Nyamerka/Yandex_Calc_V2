package stack

import (
	"testing"
)

func TestPush(t *testing.T) {
	var s Stack

	s.Push(1)
	s.Push("hello")
	s.Push(3.14)

	expectedLen := 3
	actualLen := s.Len()

	if actualLen != expectedLen {
		t.Errorf("Expected length %d, but got %d", expectedLen, actualLen)
	}

	expectedCap := 4 // Initial capacity is 2, but it grows to 4 after third push
	actualCap := s.Cap()

	if actualCap < expectedCap {
		t.Errorf("Expected capacity at least %d, but got %d", expectedCap, actualCap)
	}
}

func TestPop(t *testing.T) {
	var s Stack
	s.Push(1)
	s.Push("hello")

	value, err := s.Pop()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	expectedValue := "hello"
	actualValue := value

	if actualValue != expectedValue {
		t.Errorf("Expected value %v, but got %v", expectedValue, actualValue)
	}

	expectedLen := 1
	actualLen := s.Len()

	if actualLen != expectedLen {
		t.Errorf("Expected length %d, but got %d", expectedLen, actualLen)
	}

	value, err = s.Pop()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	expectedValue2 := 1
	actualValue = value

	if actualValue != expectedValue2 {
		t.Errorf("Expected value %v, but got %v", expectedValue, actualValue)
	}

	expectedLen = 0
	actualLen = s.Len()

	if actualLen != expectedLen {
		t.Errorf("Expected length %d, but got %d", expectedLen, actualLen)
	}

	_, err = s.Pop()
	if err == nil {
		t.Errorf("Expected error when popping from an empty stack")
	}
}

func TestTop(t *testing.T) {
	var s Stack
	s.Push(1)
	s.Push("hello")

	value, err := s.Top()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	expectedValue := "hello"
	actualValue := value

	if actualValue != expectedValue {
		t.Errorf("Expected value %v, but got %v", expectedValue, actualValue)
	}

	expectedLen := 2
	actualLen := s.Len()

	if actualLen != expectedLen {
		t.Errorf("Expected length %d, but got %d", expectedLen, actualLen)
	}

	s.Pop()
	s.Pop()

	_, err = s.Top()
	if err == nil {
		t.Errorf("Expected error when calling Top on an empty stack")
	}
}

func TestIsEmpty(t *testing.T) {
	var s Stack

	if !s.IsEmpty() {
		t.Errorf("Expected stack to be empty")
	}

	s.Push(1)

	if s.IsEmpty() {
		t.Errorf("Expected stack not to be empty")
	}

	s.Pop()

	if !s.IsEmpty() {
		t.Errorf("Expected stack to be empty after popping")
	}
}

func TestString(t *testing.T) {
	var s Stack
	s.Push(1)
	s.Push("hello")
	s.Push(3.14)

	expectedString := "1  |  hello  |  3.140000  |  "
	actualString := s.String()

	if actualString != expectedString {
		t.Errorf("Expected string %q, but got %q", expectedString, actualString)
	}
}
