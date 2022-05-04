package main

import (
	"fmt"
	"time"
)

type state interface {
	sleep(t time.Duration)
	watchYoutube(t time.Duration)
	work(t time.Duration)
}

type Yakov struct {
	noPower      state
	noMotivation state
	working      state

	currentState state

	power      int
	motivation int
}

func (y *Yakov) setState(s state) {
	y.currentState = s
}
func (y *Yakov) sleep(t time.Duration) {
	//fmt.Printf("before state =%T\n", y.currentState)
	y.currentState.sleep(t)
	//fmt.Printf("after current state =%T\n", y.currentState)
}
func (y *Yakov) watchYoutube(t time.Duration) {
	//fmt.Printf("before current state =%T\n", y.currentState)
	y.currentState.watchYoutube(t)
	//fmt.Printf("after current state =%T\n", y.currentState)
}
func (y *Yakov) work(t time.Duration) {
	//fmt.Printf("before current state =%T\n", y.currentState)
	y.currentState.work(t)
	//fmt.Printf("after current state =%T\n", y.currentState)
}

type noPower struct {
	Yakov *Yakov
}

func (n *noPower) sleep(t time.Duration) {
	n.Yakov.power += int(t / time.Hour)
	fmt.Println("Yakov slept well")
	n.Yakov.setState(n.Yakov.noMotivation)
}
func (n *noPower) watchYoutube(t time.Duration) {
	fmt.Println("Yakov need to sleep before watching smth")
	t /= t
}
func (n *noPower) work(t time.Duration) {
	fmt.Println("Yakov need to sleep before working smth")
	t /= t
}

type noMotivation struct {
	Yakov *Yakov
}

func (n *noMotivation) sleep(t time.Duration) {
	fmt.Println("Yakov need no sleep now")
	t /= t
}
func (n *noMotivation) watchYoutube(t time.Duration) {
	n.Yakov.motivation += int(t / time.Hour)
	fmt.Println("Yakov feel motivated")
	n.Yakov.setState(n.Yakov.working)
}
func (n *noMotivation) work(t time.Duration) {
	fmt.Println("Yakov has no motivation to work")
	t /= t
}

type working struct {
	Yakov *Yakov
}

func (w *working) sleep(t time.Duration) {
	fmt.Println("Yakov need no sleep now")
	t /= t
}
func (w *working) watchYoutube(t time.Duration) {
	fmt.Println("Yakov need no youtube now")
	t /= t
}
func (w *working) work(t time.Duration) {
	if w.Yakov.motivation-int(t/time.Hour) <= 1 {
		fmt.Println("Yakov cant work to much time.\nHe has no power for it.")
		return
	}
	if w.Yakov.motivation-int(t/time.Hour) <= 1 {
		fmt.Println("Yakov cant work to much time.\nHe has no motivation for it")
		return
	}
	for i := 0; i < int(t/time.Hour); i++ {
		fmt.Println("Yakov work hard")
		w.Yakov.motivation--
		w.Yakov.power--
	}
	if w.Yakov.power <= 1 {
		w.Yakov.setState(w.Yakov.noPower)
		return
	}
	if w.Yakov.motivation <= 1 {
		w.Yakov.setState(w.Yakov.noMotivation)
	}
}

func newYakov() *Yakov {
	y := &Yakov{}
	noPowerState := noPower{
		Yakov: y,
	}
	noMotivationState := noMotivation{
		Yakov: y,
	}
	workingState := working{
		Yakov: y,
	}
	y.noPower = &noPowerState
	y.noMotivation = &noMotivationState
	y.working = &workingState
	y.setState(&noPowerState)
	return y
}

func main() {
	Ya := newYakov()
	Ya.work(10 * time.Hour)
	Ya.watchYoutube(10 * time.Hour)
	Ya.sleep(10 * time.Hour)
	Ya.work(10 * time.Hour)
	Ya.work(8 * time.Hour)
	Ya.sleep(10 * time.Hour)
	Ya.watchYoutube(10 * time.Hour)
	Ya.work(10 * time.Hour)
	Ya.work(8 * time.Hour)

}

/*
Я не придумал примера лучше, чем модель себя, который, спит, смотрит ютуб и иногда работает. Есть несколько
состояний в каждом из которых поддерживается некий интеорфейс-список действий, которые можно отслеживать и менять
между собой.
*/
