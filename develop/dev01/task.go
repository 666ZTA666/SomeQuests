package main

import (
	"fmt"
	"github.com/beevik/ntp"
	"os"
	"time"
)

func main() {
	ntpTime, err := ntp.Time("0.beevik-ntp.pool.ntp.org")
	if err != nil {
		_, err := fmt.Fprintln(os.Stderr, "Error:", err.Error())
		if err != nil {
			return
		}
		os.Exit(1)
	}
	ntpTimeFormatted := ntpTime.Format(time.UnixDate)

	fmt.Println("Network time:", ntpTime)
	fmt.Println("Unix Date Network time:", ntpTimeFormatted)
	fmt.Println("+++++++++++++++++++++++++++++++")
	timeFormatted := time.Now().Local().Format(time.UnixDate)
	fmt.Println("System time:", time.Now())
	fmt.Println("Unix Date System time:", timeFormatted)
}
