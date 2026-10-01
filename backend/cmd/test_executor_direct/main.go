package main

import (
	"accode-go/executor"
	"fmt"
)

func main() {
	outputFunc := func(s string) {
		fmt.Print(s)
	}
	inputFunc := func() string {
		return ""
	}

	exec := executor.NewExecutor(outputFunc, inputFunc)

	fmt.Println("Test 1: 顯示 3。")
	err := exec.Execute("顯示 3。")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
	fmt.Println()

	fmt.Println("Test 2: 顯示 1 加上 2。")
	err = exec.Execute("顯示 1 加上 2。")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
	fmt.Println()

	fmt.Println("Test 3: 設定變數")
	code := `設定一指定名稱變數的值，名稱與值分別為「x」，10。
要求在終端機（或控制台）輸出變數「x」。`
	err = exec.Execute(code)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}
