package main

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

func main() {
	// channel作成
	ch := make(chan int)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		// channelはゴルーチン間でデータを送受信するときに用いる
		// <-ch 読み込み ch<- 書き込み
		// chへの読み込みも書き込みがあるまでブロックされる
		ch <- 10
		time.Sleep(500 * time.Millisecond)
	}()
	// chの読み込み
	fmt.Println(<-ch)
	wg.Wait()
	// goroutine leak
	ch1 := make(chan int)
	go func() {
		// chに書き込みが行われるのを待つ
		fmt.Println(<-ch1)
	}()
	ch1 <- 10
	fmt.Printf("num of working goroutines: %d\n", runtime.NumGoroutine())
	// deadlock
	ch2 := make(chan int, 1) // 第三引数でバッファ指定
	ch2 <- 2 // この時点でバッファサイズがいっぱいになる
	ch2 <- 3
	fmt.Println(<-ch2)
}
