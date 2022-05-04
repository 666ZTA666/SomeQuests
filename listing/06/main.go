package main

import (
	"fmt"
)

func main() {
	var s = []string{"1", "2", "3"}
	modifySlice(s)
	fmt.Println(s)
}
func modifySlice(i []string) { //слайс не передан, а скопирован, так как срез - структура с указателем на массив.
	i[0] = "3"         //меняется первый элемент массива под капотом.
	i = append(i, "4") //после аппенда меняется подкапотный массив, из-за изменения емкости слайса.
	i[1] = "5"         //все изменения включая предыдущий аппенд не отразятся на срезе вне функции.
	i = append(i, "6")
	// чтобы поменять массив, который по ходу работы нужно было увеличить или еще каким-то образом проаппендить,
	//лучше использовать return []string
}
