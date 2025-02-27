package stack

import (
	"errors"
	"fmt"
	"strings"
)

type Stack []interface{}

func (stack *Stack) String() string {
	result := strings.Builder{}
	for _, value := range *stack {
		switch v := value.(type) {
		case int:
			result.WriteString(fmt.Sprintf("%d  |  ", v))
		case string:

			result.WriteString(fmt.Sprintf("%s  |  ", v))
		case float64:

			result.WriteString(fmt.Sprintf("%f  |  ", v))
		default:

			result.WriteString(fmt.Sprintf("%v  |  ", v))
		}
	}

	return result.String()
}

func (stack *Stack) Pop() (interface{}, error) {
	theStack := *stack
	if len(theStack) == 0 {
		return nil, errors.New("can't Pop() an empty stack")
	}
	x := theStack[len(theStack)-1]
	*stack = theStack[:len(theStack)-1]
	return x, nil
}

func (stack *Stack) Push(x interface{}) {
	*stack = append(*stack, x)
}

func (stack *Stack) Top() (interface{}, error) {
	if len(*stack) == 0 {
		return nil, errors.New("can't Top() an empty stack")
	}
	return (*stack)[len(*stack)-1], nil
}

func (stack *Stack) Cap() int {
	return cap(*stack)
}

func (stack *Stack) Len() int {
	return len(*stack)
}

func (stack *Stack) IsEmpty() bool {
	return len(*stack) == 0
}
