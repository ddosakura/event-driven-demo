package edd

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func reading(r *Reader) {
	r.Read(func(d string, e error) {
		if e == EOF {
			println("EOF")
			return
		}
		if e != nil {
			panic(e)
		}
		print(d)
		reading(r)
	})
}

func TestMain(t *testing.T) {
	p := NewPipe()
	r := NewReader(p)
	w := NewWriter(p)

	reading(r)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
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
