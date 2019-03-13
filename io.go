package edd

import (
	"errors"
)

var (
	EOF = errors.New("EOF")
)

type Pipe struct {
	e      *EventEmitter
	closed bool
}

func NewPipe() *Pipe {
	return &Pipe{
		e:      NewEventEmitter(),
		closed: false,
	}
}

type Reader struct {
	p *Pipe
}

func NewReader(p *Pipe) *Reader {
	return &Reader{
		p,
	}
}

func (r *Reader) Read(fn func(string, error)) {
	onData := func(e *Event) {
		defer func() {
			err := recover()
			if err != nil {
				fn("", err.(error))
			}
		}()
		fn(e.d.(string), nil)
	}
	onClose := func(e *Event) {
		fn("", EOF)
		r.Close()
	}
	r.p.e.Once("data", &onData)
	r.p.e.Clean("close")
	r.p.e.Once("close", &onClose)
}

func (r *Reader) Close() {
	r.p.closed = true
	r.p.e = nil
}

type Writer struct {
	p *Pipe
}

func NewWriter(p *Pipe) *Writer {
	return &Writer{
		p,
	}
}

func (r *Writer) Write(d string) {
	if r.p.closed {
		return
	}
	r.p.e.Emit("data", d)
}

func (r *Writer) Close() {
	r.p.closed = true
	r.p.e.Emit("close", nil)
}
