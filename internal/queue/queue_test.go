package queue

import (
	"testing"
)

func TestNewQueue(t *testing.T) {
	q := New()
	if q == nil {
		t.Errorf("Expected non-nil queue, got nil")
	}
	if q.Len() != 0 {
		t.Errorf("Expected length 0, got %d", q.Len())
	}
	if q.front != 0 {
		t.Errorf("Expected front index 0, got %d", q.front)
	}
	if q.back != 0 {
		t.Errorf("Expected back index 0, got %d", q.back)
	}
}

func TestPushFront(t *testing.T) {
	q := New()
	q.PushFront(1)
	q.PushFront("hello")
	q.PushFront(3.14)

	expectedLen := 3
	actualLen := q.Len()
	if actualLen != expectedLen {
		t.Errorf("Expected length %d, got %d", expectedLen, actualLen)
	}

	expectedFront := 3.14
	actualFront := q.Front()
	if actualFront != expectedFront {
		t.Errorf("Expected front %v, got %v", expectedFront, actualFront)
	}

	expectedBack := 1
	actualBack := q.Back()
	if actualBack != expectedBack {
		t.Errorf("Expected back %v, got %v", expectedBack, actualBack)
	}

	expectedString := "[3.14 hello 1]"
	actualString := q.String()
	if actualString != expectedString {
		t.Errorf("Expected string %q, got %q", expectedString, actualString)
	}
}

func TestPushBack(t *testing.T) {
	q := New()
	q.PushBack(1)
	q.PushBack("hello")
	q.PushBack(3.14)

	expectedLen := 3
	actualLen := q.Len()
	if actualLen != expectedLen {
		t.Errorf("Expected length %d, got %d", expectedLen, actualLen)
	}

	expectedFront := 1
	actualFront := q.Front()
	if actualFront != expectedFront {
		t.Errorf("Expected front %v, got %v", expectedFront, actualFront)
	}

	expectedBack := 3.14
	actualBack := q.Back()
	if actualBack != expectedBack {
		t.Errorf("Expected back %v, got %v", expectedBack, actualBack)
	}

	expectedString := "[1 hello 3.14]"
	actualString := q.String()
	if actualString != expectedString {
		t.Errorf("Expected string %q, got %q", expectedString, actualString)
	}
}

func TestPopFront(t *testing.T) {
	q := New()
	q.PushBack(1)
	q.PushBack("hello")
	q.PushBack(3.14)

	value := q.PopFront()
	expectedValue := 1
	actualValue := value
	if actualValue != expectedValue {
		t.Errorf("Expected value %v, got %v", expectedValue, actualValue)
	}

	expectedLen := 2
	actualLen := q.Len()
	if actualLen != expectedLen {
		t.Errorf("Expected length %d, got %d", expectedLen, actualLen)
	}

	expectedFront := "hello"
	actualFront := q.Front()
	if actualFront != expectedFront {
		t.Errorf("Expected front %v, got %v", expectedFront, actualFront)
	}

	expectedBack := 3.14
	actualBack := q.Back()
	if actualBack != expectedBack {
		t.Errorf("Expected back %v, got %v", expectedBack, actualBack)
	}

	expectedString := "[hello 3.14]"
	actualString := q.String()
	if actualString != expectedString {
		t.Errorf("Expected string %q, got %q", expectedString, actualString)
	}

	// Pop remaining elements
	q.PopFront()
	q.PopFront()

	expectedLen = 0
	actualLen = q.Len()
	if actualLen != expectedLen {
		t.Errorf("Expected length %d, got %d", expectedLen, actualLen)
	}

	expectedString = "[]"
	actualString = q.String()
	if actualString != expectedString {
		t.Errorf("Expected string %q, got %q", expectedString, actualString)
	}

	// Pop from empty queue
	value = q.PopFront()
	if value != nil {
		t.Errorf("Expected nil when popping from empty queue, got %v", value)
	}
}

func TestPopBack(t *testing.T) {
	q := New()
	q.PushBack(1)
	q.PushBack("hello")
	q.PushBack(3.14)

	value := q.PopBack()
	expectedValue := 3.14
	actualValue := value
	if actualValue != expectedValue {
		t.Errorf("Expected value %v, got %v", expectedValue, actualValue)
	}

	expectedLen := 2
	actualLen := q.Len()
	if actualLen != expectedLen {
		t.Errorf("Expected length %d, got %d", expectedLen, actualLen)
	}

	expectedFront := 1
	actualFront := q.Front()
	if actualFront != expectedFront {
		t.Errorf("Expected front %v, got %v", expectedFront, actualFront)
	}

	expectedBack := "hello"
	actualBack := q.Back()
	if actualBack != expectedBack {
		t.Errorf("Expected back %v, got %v", expectedBack, actualBack)
	}

	expectedString := "[1 hello]"
	actualString := q.String()
	if actualString != expectedString {
		t.Errorf("Expected string %q, got %q", expectedString, actualString)
	}

	// Pop remaining elements
	q.PopBack()
	q.PopBack()

	expectedLen = 0
	actualLen = q.Len()
	if actualLen != expectedLen {
		t.Errorf("Expected length %d, got %d", expectedLen, actualLen)
	}

	expectedString = "[]"
	actualString = q.String()
	if actualString != expectedString {
		t.Errorf("Expected string %q, got %q", expectedString, actualString)
	}

	// Pop from empty queue
	value = q.PopBack()
	if value != nil {
		t.Errorf("Expected nil when popping from empty queue, got %v", value)
	}
}

func TestFront(t *testing.T) {
	q := New()
	q.PushBack(1)
	q.PushBack("hello")

	expectedFront := 1
	actualFront := q.Front()
	if actualFront != expectedFront {
		t.Errorf("Expected front %v, got %v", expectedFront, actualFront)
	}

	q.PopFront()

	expectedFront2 := "hello"
	actualFront = q.Front()
	if actualFront != expectedFront2 {
		t.Errorf("Expected front %v, got %v", expectedFront, actualFront)
	}

	q.PopFront()

	actualFront = q.Front()
	if actualFront != nil {
		t.Errorf("Expected nil when getting front of empty queue, got %v", actualFront)
	}
}

func TestBack(t *testing.T) {
	q := New()
	q.PushBack(1)
	q.PushBack("hello")

	expectedBack := "hello"
	actualBack := q.Back()
	if actualBack != expectedBack {
		t.Errorf("Expected back %v, got %v", expectedBack, actualBack)
	}

	q.PopBack()

	expectedBack2 := 1
	actualBack = q.Back()
	if actualBack != expectedBack2 {
		t.Errorf("Expected back %v, got %v", expectedBack, actualBack)
	}

	q.PopBack()

	actualBack = q.Back()
	if actualBack != nil {
		t.Errorf("Expected nil when getting back of empty queue, got %v", actualBack)
	}
}

func TestLen(t *testing.T) {
	q := New()
	if q.Len() != 0 {
		t.Errorf("Expected length 0, got %d", q.Len())
	}

	q.PushBack(1)
	if q.Len() != 1 {
		t.Errorf("Expected length 1, got %d", q.Len())
	}

	q.PushBack("hello")
	if q.Len() != 2 {
		t.Errorf("Expected length 2, got %d", q.Len())
	}

	q.PopFront()
	if q.Len() != 1 {
		t.Errorf("Expected length 1, got %d", q.Len())
	}

	q.PopFront()
	if q.Len() != 0 {
		t.Errorf("Expected length 0, got %d", q.Len())
	}
}

func TestString(t *testing.T) {
	q := New()
	expectedString := "[]"
	actualString := q.String()
	if actualString != expectedString {
		t.Errorf("Expected string %q, got %q", expectedString, actualString)
	}

	q.PushBack(1)
	expectedString = "[1]"
	actualString = q.String()
	if actualString != expectedString {
		t.Errorf("Expected string %q, got %q", expectedString, actualString)
	}

	q.PushBack("hello")
	expectedString = "[1 hello]"
	actualString = q.String()
	if actualString != expectedString {
		t.Errorf("Expected string %q, got %q", expectedString, actualString)
	}

	q.PushBack(3.14)
	expectedString = "[1 hello 3.14]"
	actualString = q.String()
	if actualString != expectedString {
		t.Errorf("Expected string %q, got %q", expectedString, actualString)
	}

	q.PopFront()
	expectedString = "[hello 3.14]"
	actualString = q.String()
	if actualString != expectedString {
		t.Errorf("Expected string %q, got %q", expectedString, actualString)
	}

	q.PopBack()
	expectedString = "[hello]"
	actualString = q.String()
	if actualString != expectedString {
		t.Errorf("Expected string %q, got %q", expectedString, actualString)
	}

	q.PopBack()
	expectedString = "[]"
	actualString = q.String()
	if actualString != expectedString {
		t.Errorf("Expected string %q, got %q", expectedString, actualString)
	}
}
