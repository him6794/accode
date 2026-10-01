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

	code1 := "要求在終端機（或控制台）輸出「Hello」。"
	fmt.Printf("測試 1: %s\n", code1)
	err := exec.Execute(code1)
	if err != nil {
		fmt.Printf("錯誤: %v\n", err)
	} else {
		fmt.Printf("輸出: %s\n\n", output.String())
	}

	output.Reset()
	code2 := "設定一指定名稱變數的值，名稱與值分別為「x」，10。\n設定一指定名稱變數的值，名稱與值分別為「y」，20。\n要求在終端機（或控制台）輸出變數「x」加上變數「y」。"
	fmt.Printf("測試 2:\n%s\n", code2)
	err = exec.Execute(code2)
	if err != nil {
		fmt.Printf("錯誤: %v\n", err)
	} else {
		fmt.Printf("輸出: %s\n", output.String())
	}
}
