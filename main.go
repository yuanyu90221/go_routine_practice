package main

import (
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"
)

func main() {
	outChan := make(chan string, 100)
	finishChan := make(chan struct{})

	wg := sync.WaitGroup{}
	wg.Add(20)
	for i := 0; i < 20; i++ {
		go func(val int, wg *sync.WaitGroup) {
			time.Sleep(time.Duration(rand.Int31n(1000)) * time.Millisecond)
			fmt.Println("finished job id:", val)
			wg.Done()
		}(i, &wg)
	}
	go func() {
		wg.Wait()
		close(finishChan)
	}()
Loop:
	for {
		select {
		case out := <-outChan:
			log.Println(out)
		case <-finishChan:
			break Loop // break when finish channel receive message
		}
	}
}
