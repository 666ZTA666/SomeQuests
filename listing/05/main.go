package main

import "fmt"

type customError struct {
	msg string
}

func (e *customError) Error() string {
	return e.msg
}
func test() *customError { // если поменять тип возвращаемого значения на error, то все будет ок.
	{
		// do something
	}
	return nil //nil интерфейса customError
}
func main() {
	var err error //у error как у builtin интерфейса может быть значение nil?
	err = test()
	if err != nil { //nil с которым мы сравниваем err совершенно пустой, а err это значение nil интерфейса CERR
		// интерфейс == nil когда не определен тип и данных нет, в данном случае определен тип - *customError
		fmt.Println("error")
		return
	}
	println("ok")
}

/*
если мы можем сравнивать ошибку и nil то всё ок, а сравнивать customError и nil.
*/
