package main

import (
	"fmt"
	"sync"
)

func lock1(m sync.Mutex) {
	fmt.Println("Start lock1")
	m.Lock()
	defer m.Unlock()
	lock2(m)
	fmt.Println("End lock1")
}

func lock2(m sync.Mutex) {
	fmt.Println("Start lock2")
	m.Lock()
	defer m.Unlock()
	fmt.Println("End lock2")
}

func main() {
	var m sync.Mutex
	lock1(m)
}
