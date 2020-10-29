package event

type Listener struct {
	skey     string
	callback func(ae AEvent)
	etype    []string
}

func NewListener() *Listener {
	return &Listener{}
}

func (m *Listener) subscribe() {
	if m.callback != nil && Manager != nil && len(m.etype) > 0 {
		if m.skey != "" {
			Manager.UpdateSubscription(m.skey, m.callback, m.etype...)
		} else {
			m.skey = Manager.Subscribe(m.callback, m.etype...)
		}
	}
}

func (m *Listener) Set(callback func(ae AEvent), etype ...string) {
	m.SetEventType(etype...)
	m.SetCallback(callback)
}

func (m *Listener) SetCallback(callback func(ae AEvent)) {
	m.callback = callback
	m.subscribe()
}

func (m *Listener) SetEventType(etype ...string) {
	m.etype = etype
	m.subscribe()
}
