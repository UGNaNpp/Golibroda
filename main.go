package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

const (
	MAX_CUSTOMERS = 8
	NUM_CHAIRS    = 3
)

var (
	waitingCustomers []int
	mutex            sync.Mutex
	cond             = sync.NewCond(&mutex)
	barberSleeping   = true
)

func barber() {
	for {
		mutex.Lock()
		for len(waitingCustomers) == 0 {
			fmt.Println("The barber is sleeping...")
			barberSleeping = true
			cond.Wait()
		}

		customer := waitingCustomers[0]
		waitingCustomers = waitingCustomers[1:]
		fmt.Printf("The barber is cutting hair for customer %d\n", customer)
		mutex.Unlock()

		time.Sleep(2 * time.Second)
		fmt.Printf("The barber has finished cutting hair for customer %d\n", customer)

		mutex.Lock()
		barberSleeping = false
		cond.Signal()
		mutex.Unlock()
	}
}

func customer(index int) {
	time.Sleep(time.Duration(rand.Intn(5)+1) * time.Second)

	mutex.Lock()
	if len(waitingCustomers) < NUM_CHAIRS {
		waitingCustomers = append(waitingCustomers, index)
		fmt.Printf("Customer %d is waiting in the waiting room\n", index)

		if barberSleeping {
			cond.Signal()
		}
		mutex.Unlock()

		mutex.Lock()
		for barberSleeping {
			cond.Wait()
		}
		mutex.Unlock()

		fmt.Printf("Customer %d has finished getting a haircut\n", index)
	} else {
		fmt.Printf("Customer %d is leaving because the waiting room is full\n", index)
		mutex.Unlock()
	}
}

func main() {
	go barber()
	var wg sync.WaitGroup
	for i := 0; i < MAX_CUSTOMERS; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			customer(index)
		}(i)
	}

	wg.Wait()
}
