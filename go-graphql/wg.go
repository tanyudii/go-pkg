package go_graphql

import "sync"

type Wg struct {
	sync.WaitGroup
}

func (w *Wg) Wrap(cb func()) {
	w.Add(1)
	go func() {
		cb()
		w.Done()
	}()
}
