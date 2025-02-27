package queue

import (
	"bytes"
	"fmt"
)

type Queue struct {
	rep    []interface{}
	front  int
	back   int
	length int
}

func New() *Queue {
	return new(Queue).Init()
}

func (q *Queue) Init() *Queue {
	q.rep = make([]interface{}, 1)
	q.front, q.back, q.length = 0, 0, 0
	return q
}

func (q *Queue) lazyInit() {
	if q.rep == nil {
		q.Init()
	}
}

func (q *Queue) Len() int {
	return q.length
}

func (q *Queue) empty() bool {
	return q.length == 0
}

func (q *Queue) full() bool {
	return q.length == len(q.rep)
}

func (q *Queue) sparse() bool {
	return 1 < q.length && q.length < len(q.rep)/4
}

func (q *Queue) resize(size int) {
	adjusted := make([]interface{}, size)
	if q.front < q.back {
		copy(adjusted, q.rep[q.front:q.back])
	} else {
		n := copy(adjusted, q.rep[q.front:])
		copy(adjusted[n:], q.rep[:q.back])
	}
	q.rep = adjusted
	q.front = 0
	q.back = q.length
}

func (q *Queue) lazyGrow() {
	if q.full() {
		q.resize(len(q.rep) * 2)
	}
}

func (q *Queue) lazyShrink() {
	if q.sparse() {
		q.resize(len(q.rep) / 2)
	}
}

func (q *Queue) String() string {
	var result bytes.Buffer
	result.WriteByte('[')
	j := q.front
	for i := 0; i < q.length; i++ {
		result.WriteString(fmt.Sprintf("%v", q.rep[j]))
		if i < q.length-1 {
			result.WriteByte(' ')
		}
		j = q.inc(j)
	}
	result.WriteByte(']')
	return result.String()
}

func (q *Queue) inc(i int) int {
	return (i + 1) & (len(q.rep) - 1) // requires l = 2^n
}

func (q *Queue) dec(i int) int {
	return (i - 1) & (len(q.rep) - 1) // requires l = 2^n
}

func (q *Queue) Front() interface{} {
	return q.rep[q.front]
}

func (q *Queue) Back() interface{} {
	return q.rep[q.dec(q.back)]
}

func (q *Queue) PushFront(v interface{}) {
	q.lazyInit()
	q.lazyGrow()
	q.front = q.dec(q.front)
	q.rep[q.front] = v
	q.length++
}

func (q *Queue) PushBack(v interface{}) {
	q.lazyInit()
	q.lazyGrow()
	q.rep[q.back] = v
	q.back = q.inc(q.back)
	q.length++
}

func (q *Queue) PopFront() interface{} {
	if q.empty() {
		return nil
	}
	v := q.rep[q.front]
	q.rep[q.front] = nil
	q.front = q.inc(q.front)
	q.length--
	q.lazyShrink()
	return v
}

func (q *Queue) PopBack() interface{} {
	if q.empty() {
		return nil
	}
	q.back = q.dec(q.back)
	v := q.rep[q.back]
	q.rep[q.back] = nil
	q.length--
	q.lazyShrink()
	return v
}
