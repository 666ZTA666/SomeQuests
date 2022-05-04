package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"unicode"
)

func main() {
	s := os.Args
	if len(s) > 2 {
		_, _ = fmt.Fprintln(os.Stderr, "Too many args in command line")
		return
	}
	fmt.Println(Unpack("\\55qwe\\\\5"))
}

func Unpack(str string) string {
	var last rune           //последний символ
	var res strings.Builder //результат
	var es bool             //слэш или не слэш.
	var numR int            //символ в виде числа
	var err error           //ошибка

	if str == "" { // проверка на пустую строку
		return ""
	}
	if unicode.IsDigit([]rune(str)[0]) { //проверка, что первый символ не цифра
		_, _ = fmt.Fprintln(os.Stderr, "invalid string")
		return ""
	}

	for _, run := range []rune(str) { //бежим по строке доставая каждый символ как руну
		switch {
		case unicode.IsLetter(run): // кейс 1, символ - просто символ.
			res.WriteRune(run) //пишем его в результат
		case unicode.IsDigit(run): // кейс 2, символ - число.
			if unicode.IsDigit(last) && es {
				//если предыдущий символ тоже был числом,
				// без \ перед ним, то выходим.
				_, _ = fmt.Fprintln(os.Stderr, "invalid string")
				return ""
			}
			if es { //если прошлый символ является эскейпом
				res.WriteRune(run) //пишем результат
				es = false         //прошлый эскейп сбрасывается
			} else {
				numR, err = strconv.Atoi(string(run)) //преобразуем в число
				if err != nil {
					_, _ = fmt.Fprintln(os.Stderr, "conversion error", err.Error())
					return ""
				}
				for i := 0; i < numR-1; i++ {
					//пишем в результат, столько повторов прошлого символа
					// сколько значение цифры после него, -1, потому что,
					// когда мы встретили букву, её тоже записали.
					res.WriteRune(last)
				}
			}
		case string(run) == "\\": // кейс 3, символ - слэш.
			if es { //если у нас есть эскейп, расходуем его.
				res.WriteRune(run)
				es = false
			} else {
				//если нет эскейпа, то теперь есть
				es = true
			}
		default:
			// Здесь скорее всего будут все вылеты на знаках препинания, смайликах итп.
			// В задании нет условий обработки таких случаев.
			_, _ = fmt.Fprintln(os.Stderr, "something wrong, but I don't know what happened")
			return ""
		}
		last = run
	}
	return res.String() //объединяет руны наши в строку и возвращаем из функции
}
