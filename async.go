package edd

func Await(p *Promise) (c chan interface{}) {
	c = make(chan interface{}, 1)
	p.Then(func(d interface{}) interface{} {
		// println("get, put into chan")
		c <- d
		return d
	}, func(e error) interface{} {
		// println("close")
		close(c)
		return e
	})
	return
}

func AsyncRead(r *Reader) *Promise {
	return NewPromise(func(resolve func(interface{}) interface{}, reject func(error) interface{}) {
		// println("reading")
		r.Read(func(d string, e error) {
			if e != nil {
				reject(e)
			}
			resolve(d)
		})
	})
}
