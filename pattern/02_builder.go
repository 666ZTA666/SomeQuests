package main

import "fmt"

type Building struct {
	firstPart  string
	secondPart string
	thirdPart  string
}

type Builder interface {
	build1()
	build2()
	build3()
	buildAll() Building
}

type normalBuilder struct {
	building1 string
	building2 string
	building3 string
}

func NewNormalBuilder() *normalBuilder {
	return &normalBuilder{}
}

func (b *normalBuilder) build1() {
	b.building1 = "first building part is normal"
}

func (b *normalBuilder) build2() {
	b.building2 = "second building part is normal"
}

func (b *normalBuilder) build3() {
	b.building3 = "third building part is normal"
}

func (b *normalBuilder) buildAll() Building {
	fmt.Println("normal building is ready")
	return Building{
		firstPart:  b.building1,
		secondPart: b.building2,
		thirdPart:  b.building3,
	}
}

type anotherBuilder struct {
	building1 string
	building2 string
	building3 string
}

func NewAnotherBuilder() *anotherBuilder {
	return &anotherBuilder{}
}

func (b *anotherBuilder) build1() {
	b.building1 = "first building part is another"
}

func (b *anotherBuilder) build2() {
	b.building2 = "second building part is another"
}

func (b *anotherBuilder) build3() {

}

func (b *anotherBuilder) buildAll() Building {
	fmt.Println("another building is ready")
	return Building{
		firstPart:  b.building1,
		secondPart: b.building2,
		thirdPart:  b.building3,
	}
}

type director struct {
	builder Builder
}

func newDirector(b Builder) *director {
	return &director{
		builder: b,
	}
}

func (d *director) setBuilder(b Builder) {
	d.builder = b
}

func (d *director) buildHouse() Building {
	d.builder.build1()
	d.builder.build2()
	d.builder.build3()
	return d.builder.buildAll()
}

func main() {
	builder1 := NewNormalBuilder()
	builder2 := NewAnotherBuilder()

	director := newDirector(builder1)
	normalHouse := director.buildHouse()

	fmt.Println("Normal building 1:", normalHouse.firstPart)
	fmt.Println("Normal building 2:", normalHouse.secondPart)
	fmt.Println("Normal building 3:", normalHouse.thirdPart)
	fmt.Println()

	director.setBuilder(builder2)
	anotherHouse := director.buildHouse()

	fmt.Println("Another building 1:", anotherHouse.firstPart)
	fmt.Println("Another building 2:", anotherHouse.secondPart)
	fmt.Println("Another building 3:", anotherHouse.thirdPart)

}

/*
Насколько Я понимаю, паттерн билдер используется для создания экземпляров структур, в котором множество
не обязательных параметров. Так же есть директоры, которые позволяют создавать насколько-то законченный экземпляр.
У директора есть алгоритм построения экземпляра, который последовательно вызывает команды соответствующих билдеров.
Иначе говоря, если функция-конструктор принимает слишком много аргументов, из-за чего становится очень неудобной,
используется этот паттерн, который позволяет назначить значение каждого поля экземпляра по отдельности.
*/
