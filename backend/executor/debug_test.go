package executor

import (
	"fmt"
	"strings"
	"testing"
)

func TestDebugUTF8(t *testing.T) {
	var output strings.Builder
	outputFunc := func(s string) {
		output.WriteString(s)
	}
	inputFunc := func() string {
		return ""
	}

	exec := NewExecutor(outputFunc, inputFunc)

	code := `要求在終端機（或控制台）輸出「Hello」加上「World」。`

	fmt.Printf("原始程式碼: %s\n", code)
	fmt.Printf("程式碼長度: %d bytes\n", len(code))
	fmt.Printf("程式碼 runes: %d\n", len([]rune(code)))

	idx := strings.Index(code, "加上")
	fmt.Printf("'加上' 位置: %d\n", idx)
	if idx >= 0 {
		fmt.Printf("'加上' 前後: ...%s...\n", code[idx-3:idx+9])
	}

	err := exec.Execute(code)
	if err != nil {
		fmt.Printf("執行錯誤: %v\n", err)
	}

	result := output.String()
	fmt.Printf("輸出結果: %s\n", result)
}
