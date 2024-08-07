package main

// 実行コマンド
// go tool trace trace.out

import (
	"context"
	"fmt"
	"log"
	"os"
	"runtime"
	"runtime/trace"
	"sync"
	"time"
)

func main() {
	// var wg sync.WaitGroup
	// wg.Add(1)
	// go func() {
	// 	defer wg.Done() //ゴルーチン終了合図
	// 	fmt.Println("goroutine invoked")
	// }()
	// wg.Wait() // ゴルーチン終了するのを待つ
	// fmt.Printf("num of working goroutines: %d\n", runtime.NumGoroutine()) 起動しているゴルーチンの数を取得
	// fmt.Println("main func finish")
	f, err := os.Create("trace.out")
	if err != nil {
		log.Fatalln("Error:", err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			log.Fatalln("Error:", err)
		}
	}()
	if err := trace.Start(f); err != nil {
		log.Fatalln("Error:", err)
	}
	defer trace.Stop()
	ctx, t := trace.NewTask(context.Background(), "main")
	defer t.End()
	fmt.Println("The number of logical CPU Cores:", runtime.NumCPU())

	// task(ctx, "Task1")
	// task(ctx, "Task2")
	// task(ctx, "Task3")
	var wg sync.WaitGroup
	wg.Add(3)
	go cTask(ctx, &wg, "Task1") // 先頭にgoを記載することで、ゴルーチンとして認識させることができる
	go cTask(ctx, &wg, "Task2")
	go cTask(ctx, &wg, "Task3")
	wg.Wait()

	s := []int{1, 2, 3}
	for _, i := range s {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			fmt.Println(i)
		}(i)
	}
	wg.Wait()
	fmt.Println("main func finish")
}
func task(ctx context.Context, name string) {
	defer trace.StartRegion(ctx, name).End()
	time.Sleep(time.Second)
	fmt.Println(name)
}
func cTask(ctx context.Context, wg *sync.WaitGroup, name string) {
	defer trace.StartRegion(ctx, name).End() // チェーンしている時は最後のメソッドのみ遅延する
	defer wg.Done()
	time.Sleep(time.Second)
	fmt.Println(name)
}
