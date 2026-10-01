package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

func testAPI(code string, description string) {
	fmt.Printf("\n=== %s ===\n", description)
	fmt.Printf("Code: %s\n", code)

	payload := map[string]string{"code": code}
	jsonData, _ := json.Marshal(payload)

	resp, err := http.Post("http://localhost:9000/api/run", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	defer resp.Body.Close()

	var result map[string]string
	json.NewDecoder(resp.Body).Decode(&result)
	sessionID := result["sessionId"]
	fmt.Printf("Session ID: %s\n", sessionID)

	time.Sleep(500 * time.Millisecond)

	statusResp, _ := http.Get(fmt.Sprintf("http://localhost:9000/api/status/%s", sessionID))
	defer statusResp.Body.Close()

	body, _ := io.ReadAll(statusResp.Body)
	fmt.Printf("Result: %s\n", string(body))
}

func main() {

	testAPI("要求在終端機（或控制台）輸出 1 加上 2。", "基本運算")

	code2 := `設定一指定名稱變數的值，名稱與值分別為「x」，10。
設定一指定名稱變數的值，名稱與值分別為「y」，20。
要求在終端機（或控制台）輸出變數「x」加上變數「y」。`
	testAPI(code2, "變數運算")

	code3 := `設定一指定名稱變數的值，名稱與值分別為「name」，「World」。
要求在終端機（或控制台）輸出「Hello」加上變數「name」。`
	testAPI(code3, "字串連接")
}
