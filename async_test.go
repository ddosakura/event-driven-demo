package edd

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestAsync(t *testing.T) {
	p := NewPipe()
	r := NewReader(p)
	w := NewWriter(p)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		time.Sleep(1 * time.Second)
		for i := 0; i < 3; i++ {
			// println("writing")
			w.Write(fmt.Sprintf("No. %d\n", i))
			time.Sleep(1 * time.Second)
		}
		w.Close()
		time.Sleep(1 * time.Second)
		wg.Done()
	}()

	for {
		// println("waiting write")
		d, ok := <-Await(AsyncRead(r))
		// println("get from chan")
		if ok {
			print(d.(string))
		} else {
			println("EOF")
			break
		}
	}

	wg.Wait()
}
