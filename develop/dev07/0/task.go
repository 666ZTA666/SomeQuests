package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	start := time.Now()
	<-or(
		sig(1*time.Second),
		sig(2*time.Second),
		sig(5*time.Second),
		sig(100*time.Millisecond),
	)
	fmt.Printf("done after %v", time.Since(start))
}
func or(channels ...<-chan interface{}) <-chan interface{} {
	wg := new(sync.WaitGroup)
	superchan := make(chan interface{})

	closeChan := func(channel <-chan interface{}) {
		for val := range channel {
			superchan <- val
		}
		wg.Done()
	}

	wg.Add(len(channels))
	for _, channel := range channels {

		go closeChan(channel)
	}
	go func() {
		// Тут ожидаем
		wg.Wait()
		close(superchan)
	}()
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
