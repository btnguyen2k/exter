package main

import (
	"fmt"
	"sync"
)

func lock1(m sync.RWMutex) {
	fmt.Println("Start lock1")
	m.RLock()
	defer m.RUnlock()
	lock2(m)
	fmt.Println("End lock1")
}

func lock2(m sync.RWMutex) {
	fmt.Println("Start lock2")
	m.RLock()
	defer m.RUnlock()
	fmt.Println("End lock2")
}

func main() {
	var m sync.RWMutex
	lock1(m)
}
