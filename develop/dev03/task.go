package main

import (
	"flag"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

const razd = " "

func main() {
	R := flag.Bool("r", false, "нормальная сортировка - 1, ревёрс - 0")
	K := flag.Int("k", 1, "колонка по которой сортируем")
	N := flag.Bool("n", false, "сортировка по числам")
	U := flag.Bool("u", false, "не выводим одинаковые строки")

	flag.Parse()

	data := Read(flag.Arg(0))

	ss, h, l := Split(data)

	if *K > l {
		*K = 1
	}
	ss = Addition(ss, h, l)

	if *U {
		ss = Hash(ss, h)
	}

	ss = Sort(ss, *K, *N, *R)

	Output(ss)
}

func Read(name string) string {
	data, err := os.ReadFile(name)
	if err != nil {
		_, _ = fmt.Fprint(os.Stderr, "sort: cannot read: ", name, ": No such file\n")
		os.Exit(1)
	}
	return string(data)
}

func Split(data string) ([][]string, int, int) {
	sep := "\n"
	if strings.Contains(data, "\r") && strings.Contains(data, "\n") {
		sep = "\r\n"
	} else if strings.Contains(data, "\r") && !strings.Contains(data, "\n") {
		sep = "\r"
	}
	str := strings.Split(data, sep)
	highStr, lenStr := len(str), 0
	ss := make([][]string, highStr)
	for i := 0; i < highStr; i++ {
		ss[i] = strings.Split(str[i], razd)
		if len(ss[i]) > lenStr {
			lenStr = len(ss[i])
		}
	}
	return ss, highStr, lenStr
}

func Addition(ss [][]string, highStr, lenStr int) [][]string {
	for i := 0; i < highStr; i++ {
		for len(ss[i]) < lenStr {
			ss[i] = append(ss[i], "")
		}
	}
	return ss
}

func Hash(ss [][]string, highStr int) [][]string {
	hash := make(map[string]struct{})
	for i := 0; i < highStr; i++ {
		hash[strings.ToLower(strings.Join(ss[i], razd))] = struct{}{}
	}
	str := make([][]string, len(hash))
	i := 0
	for key := range hash {
		str[i] = strings.Split(key, razd)
		i++
	}
	return str
}

func Sort(ss [][]string, k int, n, r bool) [][]string {
	sort.Slice(ss, func(i, j int) bool {
		if !n {
			return (strings.Join(ss[i][k-1:], razd) < strings.Join(ss[j][k-1:], razd)) != r
		} else { //if n
			a1, err1 := strconv.Atoi(ss[i][k-1])
			a2, err2 := strconv.Atoi(ss[j][k-1])
			if err1 != nil && err2 != nil {
				return (strings.Join(ss[i][k-1:], razd) < strings.Join(ss[j][k-1:], razd)) != r
			} else if err1 != nil {
				return true
			} else if err2 != nil {
				return false
			}
			return (a1 < a2) != r
		}
	})
	return ss
}

func Output(ss [][]string) {
	for _, s1 := range ss {
		for _, s2 := range s1 {
			fmt.Print(s2, " ")
		}
		fmt.Println()
	}
}
