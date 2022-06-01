package main

import (
	"sync/atomic"
	"testing"
)

func TestQueue(t *testing.T) {
	queue := MakeQueue[int32]()
	var sum int32
	for i := 0; i < 100_000; i++ {
		go func() {
			queue.Push(1)
		}()
	}
	for {
		go func() {
			for {
				val, ok := queue.Pop()
				if !ok {
					return
				}
				atomic.AddInt32(&sum, val)
			}
		}()
		if atomic.LoadInt32(&sum) == 100_000 {
			break
		}
	}
}
