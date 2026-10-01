package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func main() {

	payload := map[string]string{
		"code": "要求在終端機（或控制台）輸出 1 加上 2。",
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		fmt.Printf("JSON marshal error: %v\n", err)
		return
	}

	fmt.Printf("Sending JSON: %s\n", string(jsonData))

	resp, err := http.Post("http://localhost:9000/api/run", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Printf("Request error: %v\n", err)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	fmt.Printf("Response: %s\n", string(body))
}
