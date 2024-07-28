package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	ch1 := make(chan int)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		fmt.Println(<-ch1)
	}()
	ch1 <- 10
	close(ch1)
	// 第二引数okでchannelが開いているか閉じているかを確認できる
	v, ok := <-ch1
	fmt.Printf("%v %v\n", v, ok)
	wg.Wait()

	ch2 := make(chan int, 2)
	ch2 <- 1
	ch2 <- 2
	close(ch2)
	v, ok = <-ch2
	fmt.Printf("%v %v\n", v, ok)
	v, ok = <-ch2
	fmt.Printf("%v %v\n", v, ok)
	// バッっファの値が全て読み込まれている
	v, ok = <-ch2
	fmt.Printf("%v %v\n", v, ok)

	ch3 := generateCountStream()
	for v := range ch3 {
		fmt.Println(v)
	}

	nCh := make(chan struct{}) // 空の構造体は0バイトしか消費しない
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			fmt.Printf("goroutine %v started\n", i)
			<-nCh
			fmt.Println(i)
		}(i)
	}
	time.Sleep(2 * time.Second)
	// channelがクローズされると、受信待ちのブロッキングが一斉に解除され、残りの処理を進める
	close(nCh)
	fmt.Println("unblocked by manual close")

	wg.Wait()
	fmt.Println("finish")
}
func generateCountStream() <-chan int { // <-chan intとすることで戻り値が読み込み専用の値となる
	ch := make(chan int)
	go func() {
		defer close(ch)
		for i := 0; i <= 5; i++ {
			ch <- i
		}
	}()
	return ch
}
