package main

import (
	"fmt"
	"math"
)

/*shape интерфейс фигур, в который мы внесли изменение в виде функции accept, которой не было до применения паттерна.
type shape interface {
	getType() string
	accept(visitor)
}
*/

//sqare - квадрат, у которого из параметров только длина стороны.
type square struct {
	side int
}

//accept принимающий в качестве аргумента интерфейс-посетитель, и вызываем метод посетителя,
//передавая в него нашу фигуру, в данном случае квадрат.
func (s *square) accept(v visitor) {
	v.visitForSquare(s)
}

//getType - изначальный метод интерфейса, чтобы показать, что в интерфейсе есть и другие методы,
//кроме метода для посетителя.
func (s *square) getType() string {
	return "square"
}

//circle - круг, у которого есть радиус.
type circle struct {
	radius int
}

//accept - метод для посетителя.
func (c *circle) accept(v visitor) {
	v.visitForCircle(c)
}

//getType - метод интерфейса.
func (c *circle) getType() string {
	return "circle"
}

//rectangle прямоугольник. Параметры - две стороны, h - вертикальная, l - горизонтальная.
type rectangle struct {
	l int
	h int
}

//accept - метод для посетителя.
func (t *rectangle) accept(v visitor) {
	v.visitForRectangle(t)
}

//getType - метод интерфейса.
func (t *rectangle) getType() string {
	return "rectangle"
}

//Интерфейс посетителя.
type visitor interface {
	visitForSquare(*square)
	visitForCircle(*circle)
	visitForRectangle(*rectangle)
}

//Структура для расчёта площади — первый посетитель.
type areaCalculator struct {
	area int
}

//Метод для квадрата.
func (a *areaCalculator) visitForSquare(s *square) {
	a.area = s.side * s.side
	fmt.Println("area of square", a.area)
}

//Метод для круга
func (a *areaCalculator) visitForCircle(s *circle) {
	a.area = int(math.Pi * float64(s.radius) * float64(s.radius))
	fmt.Println("area of circle:", a.area)
}

//Метод для прямоугольника
func (a *areaCalculator) visitForRectangle(s *rectangle) {
	a.area = s.h * s.l
	fmt.Println("area of rectangle:", a.area)
}

//Структура для расчёта координат середины фигуры - второй посетитель
type middleCoordinates struct {
	x int
	y int
}

//Метод для квадрата.
func (a *middleCoordinates) visitForSquare(s *square) {
	a.x = s.side / 2
	a.y = s.side / 2
	fmt.Println("middle point coordinates of square:", a.x, a.y)
}

//Метод для круга. Можно усложнить, но для примера подойдёт.
func (a *middleCoordinates) visitForCircle(s *circle) {
	a.x = s.radius
	a.y = s.radius
	fmt.Println("middle point coordinates of circle:", a.x, a.y)
}

//Метод для прямоугольника.
func (a *middleCoordinates) visitForRectangle(s *rectangle) {
	a.y = s.h / 2
	a.x = s.l / 2
	fmt.Println("middle point coordinates of rectangle:", a.x, a.y)
}

func main() {
	//создаём фигуры
	square := &square{side: 6}
	circle := &circle{radius: 3}
	rectangle := &rectangle{l: 4, h: 8}
	//создаём посетителя
	areaCalculator := &areaCalculator{}
	//Через каждую из фигур вызываем метод для посетителя, в который передаем нужного посетителя.
	//В данном случае это расчёт площади.
	square.accept(areaCalculator)
	circle.accept(areaCalculator)
	rectangle.accept(areaCalculator)
	//Создаём второго посетителя.
	middleCoordinates := &middleCoordinates{}
	//И рассчитываем координаты, через те же методы, но другой параметр.
	square.accept(middleCoordinates)
	circle.accept(middleCoordinates)
	rectangle.accept(middleCoordinates)
}

/*
Идея вызова метода из метода в целом схожа с паттерном "Адаптер" из предыдущего задания, где главная задача стоит в том,
чтобы не переписывать уйму кода, который потом надо поддерживать. В данном случае, приходится написать кучу кода для
посетителя, но если эти дополнительные функции не основные для посещаемого интерфейса, и нужны не очень часто,
то можно и написать.
*/
