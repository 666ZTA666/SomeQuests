package main

import "fmt"

type handler interface {
	execute(*request)
	setNext(handler)
}

type firstHandler struct {
	next handler
}

func (f *firstHandler) execute(r *request) {
	if r.first {
		fmt.Println(r.name, "Первая обработка уже произведена. Выходим.")
		return
	}
	fmt.Println(r.name, "Производим первую обработку запроса и передаем далее.")
	r.first = true
	f.next.execute(r)
}

func (f *firstHandler) setNext(next handler) {
	f.next = next
}

type secondHandler struct {
	next handler
}

func (s *secondHandler) execute(r *request) {
	if r.second {
		fmt.Println(r.name, "Вторая обработка уже произведена. Выходим.")
		return
	}
	fmt.Println(r.name, "Производим вторую обработку запроса и передаем далее.")
	r.second = true
	s.next.execute(r)
}

func (s *secondHandler) setNext(next handler) {
	s.next = next
}

type thirdHandler struct {
	next handler
}

func (t *thirdHandler) execute(r *request) {
	if r.third {
		fmt.Println(r.name, "Третья обработка произведена. Выходим.")
		return
	}
	fmt.Println(r.name, "Производим третью обработку запроса, работа завершена.")
	r.third = true
}

func (t *thirdHandler) setNext(next handler) {
	t.next = next
}

type request struct {
	name   string
	first  bool
	second bool
	third  bool
}

func main() {
	thirdHandler := &thirdHandler{}
	secondHandler := &secondHandler{}
	secondHandler.setNext(thirdHandler)
	firstHandler := &firstHandler{}
	firstHandler.setNext(secondHandler)

	request1 := &request{name: "req1"}
	request2 := &request{name: "req2", first: true}
	request3 := &request{name: "req3", second: true}
	request4 := &request{name: "req4", third: true}

	firstHandler.execute(request1)
	firstHandler.execute(request2)
	firstHandler.execute(request3)
	firstHandler.execute(request4)
}

/*
Паттерн "цепочка вызовов". Задается последовательность обработчиков, которые последовательно обрабатывают запрос.
В случае, если запрос не удовлетворяет правилам обработки, его обработка прерывается, т.к. запрос не валиден.
Может быть и наоборот, в случае если конкретный обработчик не может обработать запрос, он передает его далее.
Стоит понимать, что при использовании паттерна есть вероятность, что запрос не будет обработан.
*/
