package event

import "reflect"

type Combiner struct {
	InEvents          []AEvent
	IgnoreEvents      []AEvent
	InEventCheckCount int
	OutEvent          AEvent
}

func IsSameEvent(aevt, bevt AEvent) bool {
	return reflect.DeepEqual(aevt, bevt)
}

func NewCombiner() *Combiner {
	return &Combiner{}
}

func (m *Combiner) SetInEvents(evts ...AEvent) {
	m.InEvents = make([]AEvent, len(evts))
	copy(m.InEvents, evts)

	m.InEventCheckCount = 0
}

func (m *Combiner) SetIgnoreEvent(evts ...AEvent) {
	m.IgnoreEvents = make([]AEvent, len(evts))
	copy(m.IgnoreEvents, evts)
}

func (m *Combiner) SetOutEvent(evt AEvent) {
	m.OutEvent = evt
}

func (m *Combiner) Listen(evt *AEvent) {

	for idx := range m.IgnoreEvents {
		if IsSameEvent(m.IgnoreEvents[idx], *evt) {
			return
		}
	}

	if IsSameEvent(m.InEvents[m.InEventCheckCount], *evt) {
		m.InEventCheckCount++
	} else {
		m.InEventCheckCount = 0
	}

	if m.InEventCheckCount >= len(m.InEvents) {
		Bus <- m.OutEvent
		m.InEventCheckCount = 0
	}
}
