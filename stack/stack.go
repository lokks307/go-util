package stack

import "sync"

type element struct {
    data interface{}
    next *element
}

type stack struct {
    lock *sync.Mutex
    head *element
    Size int
}

func (thiz *stack) Push(data interface{}) {
    thiz.lock.Lock()

    element := new(element)
    element.data = data
    temp := thiz.head
    element.next = temp
    thiz.head = element
    thiz.Size++

    thiz.lock.Unlock()
}

func (thiz *stack) Pop() interface{} {
    if thiz.head == nil {
        return nil
    }
    thiz.lock.Lock()
    r := thiz.head.data
    thiz.head = thiz.head.next
    thiz.Size--

    thiz.lock.Unlock()

    return r
}

func New() *stack {
    stk := new(stack)
    stk.lock = &sync.Mutex{}

    return stk
}

func (thiz *stack) Peek() interface{} {
	if thiz.head == nil {
        return nil
	}

	return thiz.head.data
}