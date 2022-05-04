package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"time"
)

func main() {
	T := flag.Int("timeout", 10, "timeout")
	flag.Parse()
	pars := flag.Args()
	if len(pars) != 2 {
		fmt.Println("нужены uri и порт")
		return
	}
	myC := tcpConnect(pars, *T)
	defer myC.closing()

	go myC.readOutConn()
	myC.scanInConn()

}

type myConn struct {
	net.Conn
}

//tcpConnect создание подключения, по заданным параметрам, с переданным таймаутом.
func tcpConnect(pars []string, t int) myConn {
	conn, err := net.DialTimeout("tcp", net.JoinHostPort(pars[0], pars[1]), time.Duration(t)*time.Second)
	if err != nil {
		time.Sleep(time.Duration(t) * time.Second)
		//При подключении к несуществующему серверу, программа должна завершаться через timeout
		os.Exit(0)
	}
	return myConn{conn}
}

//readOutConn скан из конекта и вывод в stdout.
func (m *myConn) readOutConn() {
	r := bufio.NewReader(m)
	for {
		message, err := r.ReadString('\n')
		if err == io.EOF {
			fmt.Println("Connection closed by foreign host")
			break
		}
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Print(message)
	}
}

//scanInConn скан из stdin и отправка в конект.
func (m *myConn) scanInConn() {
	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		_, err := fmt.Fprintf(m, s.Text()+" / HTTP/1.0\r\n\r\n")
		if err != nil {
			fmt.Println(err)
			break
		}
	}
}

func (m *myConn) closing() {
	err := m.Close()
	if err != nil {
		fmt.Println("connection closing err", err)
	}
}
