package main

import "fmt"

type Fasad struct {
	st1 *struct1
	st2 *struct2
	st3 *struct3
}

func NewFasad() *Fasad {
	return &Fasad{st1: newSt1(), st2: newSt2(), st3: newSt3()}
}

func (f *Fasad) WorkWithData(data string) {
	fmt.Println("fasad work with", data)
	f.st1.work1(data)
	f.st2.work2(data)
	f.st3.work3(data)
}

type struct1 struct {
}

func newSt1() *struct1 {
	fmt.Println("creating first struct")
	return &struct1{}
}

func (s struct1) work1(data string) {
	fmt.Println("first struct work with", data)
}

type struct2 struct {
}

func newSt2() *struct2 {
	fmt.Println("creating second struct")
	return &struct2{}
}

func (s struct2) work2(data string) {
	fmt.Println("second struct work with", data)
}

type struct3 struct {
}

func newSt3() *struct3 {
	fmt.Println("creating third struct")
	return &struct3{}
}

func (s struct3) work3(data string) {
	fmt.Println("third struct work with", data)
}

func main() {
	fasadForClient := NewFasad()
	fmt.Println()
	fasadForClient.WorkWithData("Data")
}

/*
Разбор всех паттернов Я стараюсь делать максимально безликим и абстрактным, для того, чтобы понять суть и идею
паттернов в общем, а не в частности для каких-то конкретных случаев. Но так как безликий фасад
выглядит не примечательно добавлю от себя комментарий. При решении задания lo мной был использован паттерн фасад.
Была создана структура для хранения данных о подключении к бд, к стрим-серверу и модель получаемых данных.
И через один из методов этой структуры решалась задачи чтения сообщений из стрим-сервера, в том же методе
это сообщение из json парсилось в структуру данных, и эти данные записывались в БД и кэш.
Соответственно с точки зрения "условного" клиента, вызван был один метод-фасад. За которым скрывалась гора различных
вызовов функций, методов итп.

Подводя итог, суть паттерна в том, чтобы скрыть от глаз клиента все сложности реализации процессов,
оставляя простой интерфейс. И упростить читабельность и поддержку кода.
*/
