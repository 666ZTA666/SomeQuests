package main

import (
	"fmt"
	"math/rand"
	"time"
)

func asChan(vs ...int) <-chan int {
	c := make(chan int)
	go func() {
		for _, v := range vs {
			c <- v
			time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
		}
		close(c)
	}()
	return c
}
func merge(a, b <-chan int) <-chan int {
	c := make(chan int)
	go func() {
		for { // цикл бесконечный, выхода нет. Когда закончатся цифры из канала,
			// так как ни в один ничего не приходит в с будут лететь нули до бесконечности
			//  можно поставить условие на цикл, если каналы закрыты, то выходим из цикла
			select {
			case v := <-a:
				c <- v
			case v := <-b:
				c <- v
				/*default: так мы выходим в случае если ничего не пришло в канал,
				close(c) но с учетом небольшого сна этот сценарий может отработать рано
				*/
			}
		}
	}()
	return c
}
func main() {

	a := asChan(1, 3, 5, 7)
	b := asChan(2, 4, 6, 8)
	c := merge(a, b)
	for v := range c {
		fmt.Println(v) //тонна нулей

	}
}
