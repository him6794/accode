package executor

import (
	"strings"
	"testing"
)

func TestExecute_StringLiteralWithCJK(t *testing.T) {
	var output strings.Builder
	exec := NewExecutor(func(s string) {
		output.WriteString(s)
	}, func() string {
		return ""
	})

	code := "要求在終端機（或控制台）輸出「搶 21」。"
	if err := exec.Execute(code); err != nil {
		t.Fatalf("execute failed: %v", err)
	}

	if got, want := output.String(), "搶 21"; got != want {
		t.Fatalf("unexpected output: got %q, want %q", got, want)
	}
}

func TestExecute_ConcatStringAndNumber(t *testing.T) {
	var output strings.Builder
	exec := NewExecutor(func(s string) {
		output.WriteString(s)
	}, func() string {
		return ""
	})

	code := strings.Join([]string{
		"設定一指定名稱變數的值，名稱與值分別為「x」，21。",
		"要求在終端機（或控制台）輸出（「輪到玩家，目前剩下 」加上（變數「x」轉型「字串」））。",
	}, "\n")

	if err := exec.Execute(code); err != nil {
		t.Fatalf("execute failed: %v", err)
	}

	if got, want := output.String(), "輪到玩家，目前剩下 21"; got != want {
		t.Fatalf("unexpected output: got %q, want %q", got, want)
	}
}
