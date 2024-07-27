package main

import (
	"io"
	"log"
	"os"
)

func main() {
	file, err := os.Create("log.txt") // ファイル生成
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()
	flags := log.Lshortfile // 発生したエラーの行数追加
	// flags := log.Lshortfile | log.LstdFlags
	warnLogger := log.New(io.MultiWriter(file, os.Stderr), "WARN: ", flags) // ログ作成
	errorLogger := log.New(io.MultiWriter(file, os.Stderr), "ERROR: ", flags)

	warnLogger.Println("warning A")

	errorLogger.Fatalln("critical error") // エラーログ書き込んだ上でプログラム強制終了
}

