package main

import "testing"

var testStr1 = []struct {
	name   string
	input1 string
	input2 string
	output bool
}{
	{name: "1", input1: "пятка", input2: "тяпка", output: true},
	{name: "2", input1: "пятак", input2: "пятка", output: true},
	{name: "3", input1: "пятак", input2: "тяпка", output: true},
	{"4", "листок", "слиток", true},
	{"5", "листок", "столик", true},
	{"6", "слиток", "столик", true},
	{"7", "пятка", "листок", false},
	{"8", "пятка", "пятка", false},
	// Проверка поиска анаграмм без учета повторов.
}

func TestAnagr(t *testing.T) {
	for _, v := range testStr1 {
		t.Run(v.name, func(t *testing.T) {
			if v.output != anagr(v.input1, v.input2) {
				t.Error("Error", v.name)
			}
		})
	}
}

var testStr2 = []struct {
	name   string
	input  []string
	output []string
}{
	{"1", []string{"тяпка", "пятак", "пятка"}, []string{"пятак", "пятка", "тяпка"}},
	// проверка на сортировку длинных слов
	{"2", []string{"в", "б", "а"}, []string{"а", "б", "в"}},
	// проверка на сортировку по буквам
	{"3", []string{"аа", "аа", "аб"}, []string{"аа", "аб"}},
	//проверка на удаление повторов
}

func TestHashAndSort(t *testing.T) {
	for _, val := range testStr2 {
		t.Run(val.name, func(t *testing.T) {
			res := hashAndSort(val.input)
			for i, v := range res {
				if v != val.output[i] {
					t.Error("Error", val.name, i, "order")
				}
				if len(res) != len(val.output) {
					t.Error("Error", val.name, i, "length")
				}
			}
		})
	}
}

var testStr3 = []struct {
	name   string
	input  []string
	output []string
	res    bool
}{
	{"1", []string{"тяпка", "пятак", "пятка"}, []string{"тяпка", "пятак", "пятка", "тяпка"}, true},
	// Проверка на факт того, что ключом берется первое слово, у которого есть анаграммы. 3.
	{"2", []string{"аймак", "майка", "кайма"}, []string{"аймак", "аймак", "кайма", "майка"}, true},
	// Проверка на алфавитную сортировку результатов. 4.
	{"3", []string{"абв", "абв", "авб", "авб", "вба", "абв", "вба", "авб"}, []string{"абв", "абв", "авб", "вба"}, true},
	// Проверка факта удаления повторов в результатах. 7.
	{"4", []string{"АБв", "Абв", "аВб", "авБ", "ВбА", "абв", "вба", "авб"}, []string{"абв", "абв", "авб", "вба"}, true},
	// Проверка на приведение результатов к нижнему регистру. 6.
	{"5", []string{"ааа", "ааа", "ааа"}, []string{"ааа", ""}, false},
	//Проверка на множества из одного элемента. 5.
}

func TestSearch(t *testing.T) {
	for _, v := range testStr3 {
		t.Run(v.name, func(t *testing.T) {
			res := search(v.input)
			i := 1
			for _, k := range res[v.output[0]] {
				if k != v.output[i] {
					t.Error("Error", v.name, "strings")
				}
				if _, ok := res[v.output[0]]; ok != v.res {
					t.Error("Error", v.name, "bool")
				}
				i++
			}
			i = 1
		})
	}

}
