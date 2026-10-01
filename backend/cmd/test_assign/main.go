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

	code := "設定一指定名稱變數的值，名稱與值分別為「x」，10。"

	fmt.Printf("測試程式碼: %s\n", code)

	err := exec.Execute(code)
	if err != nil {
		fmt.Printf("錯誤: %v\n", err)
	} else {
		fmt.Printf("成功\n")
	}
}
