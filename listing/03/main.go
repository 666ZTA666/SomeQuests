package main

import (
	"fmt"
	"os"
)

func Foo() error {
	var err *os.PathError = nil
	return err
}
func main() {
	err := Foo()
	fmt.Println(err)        //nil интерфейса os.PathError
	fmt.Println(err == nil) //nil без интерфейса.
}
