package main

import (
	"accode-go/executor"
	"fmt"
	"strings"
)

func main() {
	var output strings.Builder
	outputFunc := func(s string) {
		output.WriteString(s)
	}
	inputFunc := func() string {
		return ""
	}

	exec := executor.NewExecutor(outputFunc, inputFunc)

	code := "設定一指定名稱變數的值，名稱與值分別為「x」，10。\n設定一指定名稱變數的值，名稱與值分別為「y」，20。\n要求在終端機（或控制台）輸出變數「x」加上變數「y」。"

	fmt.Printf("測試程式碼:\n%s\n\n", code)

	err := exec.Execute(code)
	if err != nil {
		fmt.Printf("執行錯誤: %v\n", err)
	} else {
		fmt.Printf("輸出結果: %s\n", output.String())
	}
}
