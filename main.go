package main

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"
)

func main() {
	outChan := make(chan string, 100)
	errChan := make(chan error, 100)
	finishChan := make(chan struct{})

	wg := sync.WaitGroup{}
	wg.Add(20)
	for i := 0; i < 20; i++ {
		go func(val int, wg *sync.WaitGroup, out chan<- string, err chan<- error) {
			time.Sleep(time.Duration(rand.Int31n(1000)) * time.Millisecond)
			out <- fmt.Sprintf("finished job id: %d", val)

			if val == 15 {
				err <- errors.New("fail  job in 15")
			}
			wg.Done()
		}(i, &wg, outChan, errChan)
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
		case err := <-errChan:
			log.Println(err)
			break Loop
		case <-finishChan:
			break Loop // break when finish channel receive message
		}
	}
}
