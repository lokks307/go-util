package djson

import "sync"

type element struct {
	data rune
	next *element
}

type RuneStack struct {
	lock *sync.RWMutex
	head *element
	Size int
}

func (m *RuneStack) Push(data rune) {
	m.lock.Lock()

	element := new(element)
	element.data = data
	temp := m.head
	element.next = temp
	m.head = element
	m.Size++

	m.lock.Unlock()
}

func (m *RuneStack) Pop() rune {
	if m.head == nil {
		return 0
	}
	m.lock.Lock()
	r := m.head.data
	m.head = m.head.next
	m.Size--

	m.lock.Unlock()

	return r
}

func NewRuneStack() *RuneStack {
	stk := new(RuneStack)
	stk.lock = &sync.RWMutex{}

	return stk
}

func (m *RuneStack) Peek() rune {
	if m.head == nil {
		return 0
	}

	m.lock.RLock()
	r := m.head.data
	m.lock.RUnlock()

	return r
}

func (m *RuneStack) IsEmpty() bool {
	m.lock.RLock()
	defer m.lock.RUnlock()
	return m.head == nil
}
