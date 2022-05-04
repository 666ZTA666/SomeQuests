package main

import "fmt"

type IExm interface {
	setName(name string)
	setPar(par int)
	getName() string
	getPar() int
}

type Exm struct {
	name string
	par  int
}

func (e *Exm) setName(name string) {
	e.name = name
}

func (e *Exm) getName() string {
	return e.name
}

func (e *Exm) setPar(par int) {
	e.par = par
}

func (e *Exm) getPar() int {
	return e.par
}

type FirstExm struct {
	Exm
}

func newFExm() *FirstExm {
	return &FirstExm{
		Exm: Exm{
			name: "first Exz",
			par:  1,
		},
	}
}

type SecondExm struct {
	Exm
}

func newSecondExm() *SecondExm {
	return &SecondExm{
		Exm: Exm{
			name: "second Exz",
			par:  2,
		},
	}
}

func getExm(exzType string) (IExm, error) {
	if exzType == "first" {
		return newFExm(), nil
	}
	if exzType == "second" {
		return newSecondExm(), nil
	}
	return nil, fmt.Errorf("wrong Exz type passed")
}

func main() {
	first, _ := getExm("first")
	second, _ := getExm("second")

	printDetails(first)
	printDetails(second)
}

func printDetails(g IExm) {
	fmt.Println("exm name:", g.getName())
	fmt.Println("par:", g.getPar())
}

/*
Паттерн "фабричный метод". Очень схож со строителем, но является чуть более простой версией порождащего паттерна.
Через фабрику можно делать близко похожие структуры с небольшими отличиями и выполняющие одинаковые методы,
а с помощью строителя можно создавать более сильно различные единицы.
*/
