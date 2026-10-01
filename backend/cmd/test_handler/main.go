package main

import (
	"encoding/json"
	"fmt"
)

func main() {

	jsonStr := `{"code":"設定一指定名稱變數的值，名稱與值分別為「x」，10。\n設定一指定名稱變數的值，名稱與值分別為「y」，20。\n要求在終端機（或控制台）輸出變數「x」加上變數「y」。"}`

	var req struct {
		Code string `json:"code"`
	}

	err := json.Unmarshal([]byte(jsonStr), &req)
	if err != nil {
		fmt.Printf("JSON 解析錯誤: %v\n", err)
		return
	}

	fmt.Printf("解析後的程式碼:\n%s\n", req.Code)
	fmt.Printf("程式碼長度: %d bytes\n", len(req.Code))

	for i, ch := range req.Code {
		if ch == '\n' {
			fmt.Printf("換行符位置: %d\n", i)
		}
	}
}
