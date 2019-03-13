package edd

type Event struct {
	e string
	d interface{}
}

type EventCallback *func(*Event)

type EventEmitter struct {
	eventLoop map[string][]EventCallback
}

func NewEventEmitter() *EventEmitter {
	return &EventEmitter{
		eventLoop: make(map[string][]EventCallback),
	}
}

func (ee *EventEmitter) On(e string, fn EventCallback) {
	// println("on", e, fn)
	if ee.eventLoop[e] == nil {
		ee.eventLoop[e] = []EventCallback{fn}
	} else {
		ee.eventLoop[e] = append(ee.eventLoop[e], fn)
	}
}

func (ee *EventEmitter) Clean(e string) {
	ee.eventLoop[e] = nil
}

func (ee *EventEmitter) CleanAll() {
	ee.eventLoop = make(map[string][]EventCallback)
}

func (ee *EventEmitter) Remove(e string, fn EventCallback) {
	// println("remove", e, fn)
	index := -1
	for i, f := range ee.eventLoop[e] {
		if f == fn {
			index = i
		}
	}
	if index > -1 {
		ee.eventLoop[e] = append(ee.eventLoop[e][:index], ee.eventLoop[e][index+1:]...)
	}
}

func (ee *EventEmitter) Once(e string, fn EventCallback) {
	var t EventCallback
	f := func(ev *Event) {
		ee.Remove(e, t)
		(*fn)(ev)
	}
	t = &f
	ee.On(e, t)
}

func (ee *EventEmitter) Emit(e string, d interface{}) {
	if ee.eventLoop[e] != nil {
		for _, fn := range ee.eventLoop[e] {
			(*fn)(&Event{
				e,
				d,
			})
		}
	}
}
