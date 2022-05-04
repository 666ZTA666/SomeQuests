Что выведет программа? Объяснить вывод программы. Объяснить внутреннее устройство интерфейсов и их отличие от пустых интерфейсов.

```go
package main

import (
	"fmt"
	"os"
)

func Foo() error {
	var err *os.PathError = nil
	return err
}

func main() {
	err := Foo()
	fmt.Println(err)
	fmt.Println(err == nil)
}
```

Ответ:
```
nil
false

Значение интерфейса состоит из конкретного значения и динамического типа: [Value, Type]
Даже если значение nil, а тип не nil, то интерфейсы не равны.
nil в принте = [nil, nil] в интерфейсе.
nil в ошибке = [nil, os.PathError] в интерфейсе.

```