package edd

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestPromise(t *testing.T) {
	p := NewPipe()
	r := NewReader(p)
	w := NewWriter(p)

	// i := 0
	onRead := func(resolve func(interface{}) interface{}, reject func(error) interface{}) {
		// println(i, "reading")
		// ii := i
		// i++
		r.Read(func(d string, e error) {
			// println(ii, "onRead")
			if e != nil {
				reject(e)
			}
			resolve(d)
		})
	}
	readOnce := func(d interface{}) interface{} {
		print(d.(string))
		return NewPromise(onRead)
	}
	NewPromise(onRead).Then(
		readOnce, PromisePass).Then(
		readOnce, PromisePass).Then(
		readOnce, PromisePass).Then(
		readOnce, PromisePass).Then(
		readOnce, PromisePass).Catch(func(e error) {
		if e == EOF {
			println("EOF")
		} else {
			t.Error(e.Error())
		}
	})

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		time.Sleep(1 * time.Second)
		for i := 0; i < 3; i++ {
			w.Write(fmt.Sprintf("No. %d\n", i))
			time.Sleep(1 * time.Second)
		}
		w.Close()
		time.Sleep(1 * time.Second)
		wg.Done()
	}()
	wg.Wait()
}
