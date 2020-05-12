# go channel practice

## 1 開啟20個 job並且在全部執行完才離開 main routine

利用 sync.WaitGroup 的Counter建立 20

透過 wg.Wait還有 wg.Done特性達到 
```golang===
wg := sync.WaitGroup{}
wg.Add(20)
for i:=0; i< 20;i++ {
    go func(val int, wg *sync.WaitGroup) {
        time.Sleep(time.Duration(rand.Int31n(1000)) * time.Millisecond)
        fmt.Printf("job %d finished", i)
        wg.Done() // wg count --
    }(i, &wg)
}
wg.Wait() // 只有wg count == 0才會從wg.Wait往下執行
```

2 利用 finishChan來 接收 channel close event
```golang==
outChan := make(chan string)
finishChan := make(chan struct{})
wg := sync.WaitGroup{}
	wg.Add(20)
	for i := 0; i < 20; i++ {
		go func(val int, wg *sync.WaitGroup, out chan<- string) {
			time.Sleep(time.Duration(rand.Int31n(1000)) * time.Millisecond)
			out <- fmt.Sprintf("finished job id: %d", val)
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
		case <-finishChan:
			break Loop // break when finish channel receive message
	}
```
3  利用 time.After來設定當超過某個執行時間就結束 整個job
```golang===
    outChan := make(chan string, 100)
	errChan := make(chan error, 100)
	finishChan := make(chan struct{})

	wg := sync.WaitGroup{}
	wg.Add(20)
	for i := 0; i < 20; i++ {
		go func(val int, wg *sync.WaitGroup, out chan<- string, err chan<- error) {
			time.Sleep(time.Duration(rand.Int31n(1000)) * time.Millisecond)
			out <- fmt.Sprintf("finished job id: %d", val)

			// if val == 15 {
			// 	err <- errors.New("fail  job in 15")
			// }
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
		case <-time.After(100 * time.Millisecond): //timeout machinism
			log.Println("timeout")
			break Loop
		}
	}
```