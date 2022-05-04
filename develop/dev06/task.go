package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	// колонки начинаются с 1, варианты f флага (int), (-int), (int-), (- int, int, int-int, int-)
	var (
		F = flag.String("f", "0", "fields")
		D = flag.String("d", "\t", "delimiter")
		S = flag.Bool("s", false, "separated")
		// если в строке нет разделителя из -d, то мы не выведем эту строку, в случае с флагом -s,
		//но, если флага нет, то выведем строку без разделителей, как единую первую колонку.

	)
	flag.Parse()
	fmt.Print("\"", *D, "\"\n")
	str := Read(flag.Arg(0))
	fmt.Println("прочитали")
	ck := WorkWithF(*F)
	fmt.Println("разобрались с флагами", ck)
	res := Split(str, ck, *D)
	fmt.Println("разделили")
	Out(res, *D, *S)
}

//Read читаем файл, возвращаем построчно
func Read(name string) []string {
	data, err := os.ReadFile(name)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "wrong name of file", err)
	}
	str := strings.Split(string(data), "\r\n")
	return str
}

//WorkWithF обрабатываем строку флага f в массив колонок
func WorkWithF(f string) []int {
	var (
		a   int
		b   int
		err error
	)
	f = strings.TrimRight(strings.TrimLeft(f, "("), ")")

	if strings.Contains(f, ",") {
		//если порядок введения столбцов в норме, то все будет работать
		// если порядок типа (2, -1, 5-2) итп, то за результат не отвечаю.
		// в случае, если номер столбца переданный в флаг f, будет повторятся,
		//например, (-2, 2-5), норм результата не будет,
		//но с повторами еще можно разобраться.
		ss := strings.Split(f, ",")
		mas := make([]int, 0, 2*len(ss))
		for i := 0; i < len(ss); i++ {
			if strings.Contains(ss[i], "-") {
				if string(ss[i][0]) == "-" {
					ss[i] = strings.TrimLeft(ss[i], "-")
					a, err = strconv.Atoi(ss[i])
					if err != nil {
						fmt.Println("первый символ минус, а второй не число...")
					} else if a != 0 {
						for i := 1; i <= a; i++ {
							mas = append(mas, i)
						}
						//return mas
					} else if a == 0 {
						fmt.Println("колонки начинаются с 1")
						os.Exit(1)
					}
				} else if string(ss[i][len(ss[i])-1]) == "-" {
					ss[i] = strings.TrimRight(ss[i], "-")
					a, err = strconv.Atoi(ss[i])
					if err != nil {
						fmt.Println("последний символ минус, а осталось не число...")
					} else if a != 0 {
						mas = append(mas, a)
						mas = append(mas, 0)
					} else if a == 0 {
						fmt.Println("колонки начинаются с 1")
						os.Exit(1)
					}
				} else {
					a, err = strconv.Atoi(ss[i][:strings.Index(ss[i], "-")])
					if err != nil {
						fmt.Println("перед минусом не число...")
					}
					b, err = strconv.Atoi(ss[i][strings.Index(ss[i], "-")+1:])
					if err != nil {
						fmt.Println("после минуса не число...")
					}
					if a == 0 {
						fmt.Println("колонки начинаются с 1")
						os.Exit(1)
					}
					if b == 0 {
						fmt.Println("колонки начинаются с 1")
						os.Exit(1)
					}
					for i := 0; i <= b-a; i++ {
						mas = append(mas, a+i)
					}
				}
			} else {
				a, err = strconv.Atoi(ss[i])
				if err != nil {
					fmt.Println("ошибка с", i, "элементом, когда они через запятую", err)
				}
				if a == 0 {
					fmt.Println("колонки начинаются с 1")
					os.Exit(1)
				}
				mas = append(mas, a)
			}
		}
		return mas
	} else if strings.Contains(f, "-") {
		if string(f[0]) == "-" {
			f = strings.TrimLeft(f, "-")
			a, err = strconv.Atoi(f)
			if err != nil {
				fmt.Println("первый символ минус, а второй не число...")
			} else {
				if a != 0 {
					mas := make([]int, a)
					for i := 1; i <= a; i++ {
						mas[i-1] = i
					}
					return mas
				} else {
					fmt.Println("колонки начинаются с 1")
					os.Exit(1)
				}
			}
		} else if string(f[len(f)-1]) == "-" {
			f = strings.TrimRight(f, "-")
			a, err = strconv.Atoi(f)
			if err != nil {
				fmt.Println("последний символ минус, а осталось не число...")
			} else if a != 0 {
				mas := make([]int, 2)
				mas[0] = a
				mas[1] = 0
				return mas
			} else {
				fmt.Println("колонки начинаются с 1")
				os.Exit(1)
			}
		} else {
			a, err = strconv.Atoi(f[:strings.Index(f, "-")])
			if err != nil {
				fmt.Println("перед минусом не число...")
			}
			b, err = strconv.Atoi(f[strings.Index(f, "-"):])
			if err != nil {
				fmt.Println("после минуса не число...")
			}
			if a == 0 {
				fmt.Println("колонки начинаются с 1")
				os.Exit(1)
			}
			if b == 0 {
				fmt.Println("колонки начинаются с 1")
				os.Exit(1)
			}
			mas := make([]int, a-b+1)
			for i := 0; i <= a-b; i++ {
				mas[i] = a + i
			}
			return mas
		}
	} else {
		a, err = strconv.Atoi(f)
		if err != nil {
			fmt.Println("не просто число")
		} else if a != 0 {
			return []int{a}
		} else if a <= 0 {
			fmt.Println("ошибка колонки, они с 1 начинаются")
			os.Exit(1)
		}
	}
	fmt.Println("массив с 0")
	os.Exit(1)
	return nil
}

//Split выделяет нужные колонки согласно флагам d и f
func Split(str []string, ck []int, d string) []string {
	res := make([]string, len(str))
	for i := 0; i < len(str); i++ {
		ss := strings.Split(str[i], d)
		for j := 0; j < len(ck); j++ {
			if ck[j] >= len(ss) {
				res[i] = strings.Join(ss, d)
				break
			}
			if ck[j] != 0 {
				res[i] += ss[ck[j]-1] + d
			} else {
				res[i] = strings.Join(ss[ck[j-1]-1:], d)
			}
		}
	}
	return res
}

//Out выводит колонки согласно флагу s
func Out(str []string, d string, s bool) {
	for i := 0; i < len(str); i++ {
		if s {
			if strings.Contains(str[i], d) {
				fmt.Println(str[i])
			}
		} else {
			fmt.Println(str[i])
		}
	}
}
