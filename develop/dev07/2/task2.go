package main

import (
	"fmt"
	"time"
)

func main() {
	or := Unite

	start := time.Now()
	<-or(
		sig(1*time.Second),
		sig(2*time.Second),
		sig(5*time.Second),
		sig(100*time.Millisecond),
	)
	fmt.Printf("done after %v\n", time.Since(start))
}

func Unite(channels ...<-chan interface{}) <-chan interface{} {
	superchan := make(chan interface{})
	for _, ch := range channels {
		go func(ch <-chan interface{}) {
			select {
			case <-ch:
				close(superchan)
			case <-superchan:
				return
			}
		}(ch)
	}
	return superchan
}

func sig(after time.Duration) <-chan interface{} {
	c := make(chan interface{})
	go func() {
		defer close(c)
		time.Sleep(after)
	}()
	return c
}
