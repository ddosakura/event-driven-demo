package edd

type PromiseState int

const (
	_ PromiseState = iota
	promisePending
	promiseResolved
	promiseRejected
)

var (
	PromisePass = func(e error) interface{} {
		// log.Print(e)
		return e
	}
)

type Promise struct {
	state       PromiseState
	value       interface{}
	onResolveds []func()
	onRejects   []func()
}

func NewPromise(fn func(resolve func(interface{}) interface{}, reject func(error) interface{})) *Promise {
	p := &Promise{
		state:       promisePending,
		onResolveds: make([]func(), 0),
		onRejects:   make([]func(), 0),
	}
	// println(p, "created")
	fn(func(d interface{}) interface{} {
		// println(p, "resolve")
		if p.state == promisePending {
			p.value = d
			p.state = promiseResolved
			for _, f := range p.onResolveds {
				f()
			}
		}
		return d
	}, func(e error) interface{} {
		// println(p, "reject")
		if p.state == promisePending {
			p.value = e
			p.state = promiseRejected
			for _, f := range p.onRejects {
				f()
			}
		}
		return e
	})
	return p
}

func (p *Promise) Then(onResolved func(interface{}) interface{}, onRejected func(error) interface{}) *Promise {
	onResolvedHandler := func(resolve func(interface{}) interface{}, reject func(error) interface{}) {
		// println(p, "resolveHandler")
		defer func() {
			e := recover()
			if e != nil {
				reject(e.(error))
			}
		}()
		x := onResolved(p.value)
		switch x.(type) {
		case *Promise:
			x.(*Promise).Then(resolve, reject)
		default:
			resolve(x)
		}

	}
	if p.state == promiseResolved {
		// println(p, "resolved")
		return NewPromise(onResolvedHandler)
	}

	onRejectedHandler := func(resolve func(interface{}) interface{}, reject func(error) interface{}) {
		// println(p, "rejectedHandler")
		defer func() {
			e := recover()
			if e != nil {
				reject(e.(error))
			}
		}()
		x := onRejected(p.value.(error))
		switch x.(type) {
		case *Promise:
			x.(*Promise).Then(resolve, reject)
		case error:
			reject(x.(error))
		default:
			panic("reject() should return Promise or error")
		}
	}
	if p.state == promiseRejected {
		// println(p, "rejected")
		return NewPromise(onRejectedHandler)
	}

	if p.state == promisePending {
		// println(p, "pending")
		return NewPromise(func(resolve func(interface{}) interface{}, reject func(error) interface{}) {
			p.onResolveds = append(p.onResolveds, func() {
				onResolvedHandler(resolve, reject)
			})
			p.onRejects = append(p.onRejects, func() {
				onRejectedHandler(resolve, reject)
			})
		})
	}

	panic("promise state error")
}

func (p *Promise) Catch(fn func(error)) {
	p.onRejects = append(p.onRejects, func() {
		switch p.value.(type) {
		case error:
			fn(p.value.(error))
		}
	})
}
