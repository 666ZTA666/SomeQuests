package main

import (
	"fmt"
	"time"
)

func main() {
	start := time.Now()
	s, l := or(
		sig(1*time.Second),
		sig(2*time.Second),
		sig(5*time.Second),
		sig(100*time.Millisecond),
	)
	for i := 0; i < l; i++ {
		<-s
		fmt.Printf("done after %v\n", time.Since(start))
	}

	//time.Sleep(10 * time.Second)
}
func sig(after time.Duration) <-chan interface{} {
	c := make(chan interface{})
	go func() {
		defer close(c)
		time.Sleep(after)
		fmt.Println("closed")
	}()
	return c
}

func or(channels ...<-chan interface{}) (<-chan interface{}, int) {
	superchan := make(chan interface{})
	for _, ch := range channels {
		go func(ch <-chan interface{}) {
			superchan <- ch
		}(ch)
	}

	return superchan, len(channels)
}
