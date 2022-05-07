// Producer-consumer, variable sharing solution
package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// program variables, bufferSize >= 1
var producerCount int = 1
var consumerCount int = 2
var bufferSize int = 5
var buffer = make(chan int, bufferSize)
var runforever = make(chan bool)

// synchronisation structures
var mutex sync.Mutex
var semaphoreItems = make(chan int, bufferSize)
var semaphoreSpaces = make(chan int, bufferSize)

func sleep() {
	time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
}

func producer(id int) {
	for {
		// if there is empty space in the buffer continue, otherwise wait
		<-semaphoreSpaces
		mutex.Lock()

		// produce
		num := rand.Intn(100)
		buffer <- num
		fmt.Printf("Producer%d : %d\n", id, num)

		// new item in the buffer
		mutex.Unlock()
		semaphoreItems <- 1
	}
}

func consumer(id int) {
	for {
		// if there is item in the buffer continue, otherwise wait
		<-semaphoreItems
		mutex.Lock()

		// consume
		num := <-buffer
		fmt.Printf("\t\t\tConsumer%d: %d\n", id, num)
		sleep()

		// new empty space in the buffer
		mutex.Unlock()
		semaphoreSpaces <- 1
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())

	for i := 0; i < bufferSize; i++ {
		semaphoreSpaces <- 1
	}

	// launch goroutines
	for i := 0; i < producerCount; i++ {
		go producer(i + 1)
	}

	for i := 0; i < consumerCount; i++ {
		go consumer(i + 1)
	}

	// stop execution with keyboard interrupt
	<-runforever
}
