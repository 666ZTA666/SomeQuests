package main

import (
	"fmt"
	"sort"
)

func search(input []string) map[string][]string {
	result := make(map[string][]string)

	for i := 0; i < len(input); i++ {
		for j := i + 1; j < len(input); j++ {
			if anagr(input[i], input[j]) {
				result[input[i]] = append(result[input[i]], input[i], input[j])
				input = deleteSl(input, j)
				j--
			}
		}
	}

	for k := range result {
		result[k] = hashAndSort(result[k])
	}

	return result
}

func deleteSl(input []string, j int) []string {
	return append(input[:j], input[j+1:]...)
}

func anagr(one, two string) bool {
	if one == two { //одинаковые слова НЕ являются анаграммами
		return false
	}
	a := []rune(one)
	sort.Slice(a, func(i, j int) bool {
		return a[i] > a[j]
	})
	b := []rune(two)
	sort.Slice(b, func(i, j int) bool {
		return b[i] > b[j]
	})

	if string(a) == string(b) {
		return true
	}
	return false
}

func hashAndSort(input []string) []string {
	uMap := make(map[string]struct{})
	for _, v := range input {
		uMap[v] = struct{}{}
	}
	res := make([]string, len(uMap))
	i := 0
	for key := range uMap {
		res[i] = key
		i++
	}
	sort.Strings(res)
	return res

}

func main() {
	s := []string{"тяпка", "пятак", "пятка"}
	fmt.Print(search(s))
}
