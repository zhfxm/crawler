package main

import (
	"fmt"
	"time"
)

func main() {
	demo1()
	demo2()
}

func demo1() {
	var Ball int
	table := make(chan int)
	go player(table)
	go player(table)
	table <- Ball
	time.Sleep(1 * time.Second)
	<-table
}

func player(table chan int) {
	ball := <-table
	fmt.Println(ball)
	ball++
	time.Sleep(100 * time.Millisecond)
	table <- ball
}

func demo2() {
	ch1 := search("jonson")
	ch2 := search("alan")
	select{
	case msg := <-ch1: 
	fmt.Println(msg)
	case msg := <-ch2:
		fmt.Println(msg)
	}
}

func search(msg string) chan string {
	var ch = make(chan string)
	go func() {
		var i int
	for {
		ch <- fmt.Sprintf("get %s %d", msg, i)
		i++
		time.Sleep(100 * time.Millisecond)
	}
	}()
	return ch
}
