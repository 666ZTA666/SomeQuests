package main

import "fmt"

func main() {
	ch := make(chan int)
	go func() {
		for i := 0; i < 10; i++ {
			ch <- i
		}
		//close(ch)
	}()
	for n := range ch {
		fmt.Println(n)
	} //читаем до закрытия, а закрытия нет.
	//deadlock
}
