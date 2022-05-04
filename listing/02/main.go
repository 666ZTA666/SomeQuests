package main

import (
	"fmt"
)

func test() (x int) { // объявление переменной
	defer func() {
		x++ //увеличение переменной
	}()
	x = 1  //инициализация переменной
	return //вернем 2
}
func anotherTest() int {
	var x int // объявление и инициализация переменной x = 0
	defer func() {
		x++ // берём x = 0 и увеличиваем на 1
	}()
	x = 1    // приравниваем единице и уже не важно что делает дефер.
	return x //1
}
func main() {
	fmt.Println(test())        //2
	fmt.Println(anotherTest()) //1
}
