package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	// demo1()
	// fanIn()
	fanOut1()
	// time.Sleep(10 * time.Second)
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

func fanIn() {
	ch1 := search("jonson")
	ch2 := search("alan")
	for {
		select{
			case msg := <-ch1: 
			fmt.Println(msg)
			case msg := <-ch2:
				fmt.Println(msg)
			}
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

func fanOut1() {
	var wg sync.WaitGroup
	wg.Add(36)
	go pool1(&wg, 36, 50)
	wg.Wait()
}

func pool1(wg *sync.WaitGroup, workers, tasks int) {
	taskCh := make(chan int)
	for i := 0; i < workers; i++ {
		go worker1(taskCh, wg)
	}
	for i := 0; i < tasks; i++ {
		taskCh <- i
	}
	close(taskCh)
}

func worker1(taskCh <-chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		task, ok := <- taskCh
		if !ok {
			return
		}
		// d := time.Duration(task) * time.Millisecond
		// time.Sleep(d)
		fmt.Println("processing task", task)
	}
}