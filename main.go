package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"sync"
	"time"
)

var waitingClients = 0;
var waitingClientsMutex sync.Mutex; // Lepiej użyć semafora ale w go nie jest w bibliottece standardowej
var nextCustomerMutex sync.Mutex;
var barberMutex sync.Mutex;
var clientsToday = 0;

func main () {
	barberMutex.Lock()
	go barber();
	for range 10 {
		go nextCustomer();
		time.Sleep(time.Duration(rand.Float64() * 2 * float64(time.Second)))
	}
}

func nextCustomer() {
	nextCustomerMutex.Lock()
	waitingClientsMutex.Lock()
	clientsToday++;
	defer nextCustomerMutex.Unlock()
	if waitingClients > 3 { // Jeśli równo max klientów wychodzimy
		waitingClientsMutex.Unlock()
		fmt.Println("Klient nr. ", clientsToday, " właśnie wyszedł nieobsłużony.")
	} else if waitingClients > 0 {
		waitingClients++;
		waitingClientsMutex.Unlock()
	} else {
		waitingClients++;
		barberMutex.Unlock();
	}
}

func barber() {
	fmt.Println("Właśnie otwarto zakład fryzjerski")
	defer fmt.Println("Barber zakończył pracę, idzie do domu")
	for {
		waitingClientsMutex.Lock(); // Zappewnia wyłącznie bezpieczny odczyt zmiennej
		if waitingClients > 0 {
			waitingClientsMutex.Unlock();
			barberMutex.Lock()
			waitingClients--;
			time.Sleep(1 * time.Second)
			fmt.Println("Obsługa " + strconv.Itoa(clientsToday) + " clienta dzisiaj.")
			barberMutex.Unlock();
		} else {
			waitingClientsMutex.Unlock();
			barberMutex.Lock() // * Zasypia i czeka na obudzenie
		}
		if clientsToday > 5 {
			break;
		}
	}
}


