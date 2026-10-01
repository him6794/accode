package main

import (
	"encoding/json"
	"fmt"
)

func main() {
	output := "你忘了加句點。\n"

	data := map[string]interface{}{
		"output":    output,
		"completed": true,
	}

	jsonBytes, err := json.Marshal(data)
	if err != nil {
		fmt.Printf("錯誤: %v\n", err)
		return
	}

	fmt.Printf("JSON: %s\n", string(jsonBytes))

	var decoded map[string]interface{}
	json.Unmarshal(jsonBytes, &decoded)
	fmt.Printf("解碼後: %v\n", decoded["output"])
}
