package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
)

func main() {
	var (
		A  = flag.Int("A", 0, "after")
		B  = flag.Int("B", 0, "before")
		C  = flag.Int("C", 0, "context")
		c  = flag.Bool("c", false, "count")
		iC = flag.Bool("i", false, "ignoreCase")
		R  = flag.Bool("v", false, "invert")
		F  = flag.Bool("F", false, "fixed")
		N  = flag.Bool("n", false, "lineNum")
	)
	flag.Parse()
	file := flag.Args()

	if len(file) < 2 {
		log.Fatal("error args")
	}

	app := newMyStr(file[1], file[0], newFlags(*A, *B, *C, *c, *iC, *R, *F, *N))
	app.readFile()
	app.run()
	app.out()

}

// флаги
type flags struct {
	after      int
	before     int
	context    int
	count      bool
	ignoreCase bool
	invert     bool
	fixed      bool
	lineNum    bool
}

//newFlags конструктор флагов
func newFlags(after, before, context int, count, ignoreCase, invert, fixed, lineNum bool) *flags {
	return &flags{after: after, before: before, context: context, count: count, ignoreCase: ignoreCase, invert: invert, fixed: fixed, lineNum: lineNum}
}

// MyStr наша структура, с которой будем работать
type MyStr struct {
	reg  string
	file string
	flags
	str  []string
	ints []int
}

//newMyStr конструктор
func newMyStr(file, reg string, f *flags) *MyStr {
	return &MyStr{
		reg:   reg,
		file:  file,
		flags: *f,
	}
}

//readFile читает строки из файла
func (m *MyStr) readFile() {
	data, err := os.ReadFile(m.file)
	if err != nil {
		_, _ = fmt.Fprint(os.Stderr, "grep: cannot read: ", m.file, ": No such file\n")
		os.Exit(1)
	}
	m.str = strings.Split(string(data), "\r\n")

}

//run выполняем всю логику
func (m *MyStr) run() {
	var pattern *regexp.Regexp //тут будет храниться наша регулярка.

	if m.flags.ignoreCase { // -i - "ignore-case" (игнорировать регистр)
		pattern, _ = regexp.Compile("(?i)" + m.reg)
	} else {
		pattern, _ = regexp.Compile(m.reg)
	}

	for i := 0; i < len(m.str); i++ {
		if m.flags.fixed { //-F - "fixed", точное совпадение со строкой, не паттерн
			if strings.Contains(m.str[i], m.reg) { //Contains проверяет есть ли в строке то что передаем вторым параметром.
				m.ints = append(m.ints, i)
			}
		} else {
			if pattern.MatchString(m.str[i]) { //ищем в строке по нашему патерну
				m.ints = append(m.ints, i)
			}
		}
	}
	//fmt.Println(m.ints)
}

//out выводит все что нужно в консоль
func (m *MyStr) out() {
	var ifPr bool

	low := m.context   //-C - "context" (A+B) печатать ±N строк вокруг совпадения
	top := m.context   //-C - "context" (A+B) печатать ±N строк вокруг совпадения
	if m.after > low { //-A - "after" печатать +N строк после совпадения
		low = m.after
	}
	if m.before > top { //-B - "before" печатать +N строк до совпадения
		top = m.before
	}

	switch {
	case m.count && !m.invert: //-c - "count" (количество строк)
		fmt.Println(len(m.ints))
	case m.invert && !m.count: //-v - "invert" (вместо совпадения, исключать)
		for i := 0; i < len(m.str); i++ {
			if !Check(m.ints, i) {
				if m.lineNum { //-n - "line num", напечатать номер строки
					fmt.Print(i+1, m.str[i], "\n")
				} else {
					fmt.Print(m.str[i], "\n")
				}
			}
		}
	case m.invert && m.count:
		fmt.Println(len(m.str) - len(m.ints))
		if m.lineNum { //-n - "line num", напечатать номер строки
			for i := 0; i < len(m.ints); i++ {
				fmt.Println(m.ints[i] + 1)
			}
		}
	default: //если тех ключей нет, то печатаем все найденные строки
		num := -1
		for idxStr := 0; idxStr < len(m.str); idxStr++ {
			for idxInts := 0; idxInts < len(m.ints); idxInts++ {
				//условие вывода строки
				ifPr = idxStr-low <= m.ints[idxInts] && m.ints[idxInts] <= idxStr+top && idxStr != num
				if m.lineNum && ifPr { //-n - "line num", напечатать номер строки
					if Check(m.ints, idxStr) {
						fmt.Print(">")
					}
					fmt.Print(idxStr+1, m.str[idxStr], "\n")
					num = idxStr
				} else if ifPr {
					if Check(m.ints, idxStr) {
						fmt.Print(">")
					}
					fmt.Print(m.str[idxStr], "\n")
					num = idxStr
				}
			}
		}
	}
}

//Check проверяет, находится ли число в массиве,
func Check(arr []int, ind int) bool {
	for i := 0; i < len(arr); i++ {
		if ind == arr[i] {
			return true
		}
	}
	return false
}

/*

 */
