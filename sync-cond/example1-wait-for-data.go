/*
 * We use sync.Cond to make multiple goroutines to
 * pause their execution and wait before the Cond
 * is free again.
 *
 * c.Wait() has to happen before c.Broadcast()
 */
package main

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"time"
)

func listen(name string, age map[string]int, c *sync.Cond) {
	c.L.Lock()
    defer c.L.Unlock()
	c.Wait()
	fmt.Println(name, " age: ", age["T"])
}

func broadcast(name string, age map[string]int, c *sync.Cond) {
	time.Sleep(time.Second)
	c.L.Lock()
    defer c.L.Unlock()
	age["T"] = 25
	c.Broadcast()
}

func main() {
	var age = make(map[string]int)

	m := sync.Mutex{}
	c := sync.NewCond(&m)

	// listener 1
	go listen("listener1", age, c)

	// listener 2
	go listen("listener1", age, c)

	// broadcast
	go broadcast("broadcast1", age, c)

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	<-ch
}
