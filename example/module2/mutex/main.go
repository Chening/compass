package main

import (
	"fmt"
	"sync"
)

func main() {
	rLock()

	lock()

	wLock()
}

func lock(){
	lock := sync.Mutex{}

	for i :=0; i<3 ; i ++ {
		lock.Lock()
		defer lock.Unlock()

		fmt.Println("lock", i)
	}
}

func rLock() {
	lock := sync.RWMutex{}

	for i := 0; i < 3; i++ {
		lock.RLock()
		defer lock.RUnlock()

		fmt.Println("rLock", i)
	}
}

func wLock () {
	lock := sync.RWMutex{}
	for i:=0; i<3; i++ {
		lock.Lock()
		defer lock.Unlock()

		fmt.Println("wLock", i)
	}
}
