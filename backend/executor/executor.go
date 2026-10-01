package executor

import (
	"fmt"
	"strconv"
	"strings"
	"unicode/utf8"
)

const (
	IF       = "如果右方判斷式結果為真，執行以下特定操作"
	WHILE    = "重複並持續驗證右方表達式之真實性，條件滿足時持續執行特定操作"
	BREAK    = "無視迴圈判斷式要求，直接跳脫迴圈"
	CONTINUE = "捨棄以下迴圈內容，直接繼續下一輪迴圈"
	EXIT     = "無視所有指令，直接退出程式"
	END      = "結束以上判斷式或迴圈"
)

type Executor struct {
	variables map[string]interface{}
	arrays    map[string][]interface{}
	output    func(string)
	input     func() string
	stopCheck func() bool
}

func NewExecutor(outputFunc func(string), inputFunc func() string) *Executor {
	return &Executor{
		variables: make(map[string]interface{}),
		arrays:    make(map[string][]interface{}),
		output:    outputFunc,
		input:     inputFunc,
	}
}

func (e *Executor) SetStopCheck(stopCheck func() bool) {
	e.stopCheck = stopCheck
}

func (e *Executor) shouldStop() bool {
	return e.stopCheck != nil && e.stopCheck()
}

func (e *Executor) toFloat(v interface{}) (float64, error) {
	switch val := v.(type) {
	case float64:
		return val, nil
	case int:
		return float64(val), nil
	case string:
		return strconv.ParseFloat(val, 64)
	default:
		return 0, fmt.Errorf("無法轉換為數字：%v", v)
	}
}

func (e *Executor) apply(a, b interface{}, op string) (interface{}, error) {
	switch op {
	case "加上":
		if aStr, ok := a.(string); ok {
			if bStr, ok := b.(string); ok {
				return aStr + bStr, nil
			}
			return aStr + fmt.Sprint(b), nil
		}
		aNum, err := e.toFloat(a)
		if err != nil {
			return nil, err
		}
		bNum, err := e.toFloat(b)
		if err != nil {
			return nil, err
		}
		return aNum + bNum, nil
	case "減去":
		aNum, err := e.toFloat(a)
		if err != nil {
			return nil, err
		}
		bNum, err := e.toFloat(b)
		if err != nil {
			return nil, err
		}
		return aNum - bNum, nil
	case "乘以":
		if aStr, ok := a.(string); ok {
			bNum, err := e.toFloat(b)
			if err != nil {
				return nil, err
			}
			return strings.Repeat(aStr, int(bNum)), nil
		}
		aNum, err := e.toFloat(a)
		if err != nil {
			return nil, err
		}
		bNum, err := e.toFloat(b)
		if err != nil {
			return nil, err
		}
		return aNum * bNum, nil
	case "除以":
		aNum, err := e.toFloat(a)
		if err != nil {
			return nil, err
		}
		bNum, err := e.toFloat(b)
		if err != nil {
			return nil, err
		}
		if bNum == 0 {
			return nil, fmt.Errorf("除以零錯誤")
		}
		return aNum / bNum, nil
	case "取模":
		aNum, err := e.toFloat(a)
		if err != nil {
			return nil, err
		}
		bNum, err := e.toFloat(b)
		if err != nil {
			return nil, err
		}
		return float64(int(aNum) % int(bNum)), nil
	case "等於":
		aNum, aErr := e.toFloat(a)
		bNum, bErr := e.toFloat(b)
		if aErr == nil && bErr == nil {
			if aNum == bNum {
				return float64(1), nil
			}
			return float64(0), nil
		}
		if fmt.Sprint(a) == fmt.Sprint(b) {
			return float64(1), nil
		}
		return float64(0), nil
	case "大於":
		aNum, err := e.toFloat(a)
		if err != nil {
			return nil, err
		}
		bNum, err := e.toFloat(b)
		if err != nil {
			return nil, err
		}
		if aNum > bNum {
			return float64(1), nil
		}
		return float64(0), nil
	case "小於":
		aNum, err := e.toFloat(a)
		if err != nil {
			return nil, err
		}
		bNum, err := e.toFloat(b)
		if err != nil {
			return nil, err
		}
		if aNum < bNum {
			return float64(1), nil
		}
		return float64(0), nil
	case "轉型":
		bStr := fmt.Sprint(b)
		switch bStr {
		case "數字":
			return e.toFloat(a)
		case "字串":
			return fmt.Sprint(a), nil
		default:
			return nil, fmt.Errorf("未知型態：%s", bStr)
		}
	default:
		return nil, fmt.Errorf("未知運算子：%s", op)
	}
}

func precedence(op string) int {
	switch op {
	case "等於", "大於", "小於":
		return 1
	case "加上", "減去":
		return 2
	case "乘以", "除以", "取模":
		return 3
	case "轉型":
		return 4
	default:
		return 0
	}
}

func (e *Executor) evaluate(expression string) (interface{}, error) {
	ops := []string{}
	values := []interface{}{}
	i := 0

	for i < len(expression) {
		r, size := utf8.DecodeRuneInString(expression[i:])
		if r == utf8.RuneError && size == 1 {
			return nil, fmt.Errorf("無法解析字元：%q", expression[i:i+1])
		}
		ch := string(r)

		if ch == " " || ch == "　" {
			i += size
			continue
		}

		if ch == "（" {
			ops = append(ops, ch)
			i += size
			continue
		}

		if ch == "）" {
			for len(ops) > 0 && ops[len(ops)-1] != "（" {
				if len(values) < 2 {
					return nil, fmt.Errorf("運算元不足")
				}
				val2 := values[len(values)-1]
				val1 := values[len(values)-2]
				values = values[:len(values)-2]
				op := ops[len(ops)-1]
				ops = ops[:len(ops)-1]
				result, err := e.apply(val1, val2, op)
				if err != nil {
					return nil, err
				}
				values = append(values, result)
			}
			if len(ops) > 0 {
				ops = ops[:len(ops)-1]
			}
			i += size
			continue
		}

		if len(ch) == 1 && ((ch[0] >= '0' && ch[0] <= '9') || ch[0] == '.' || ch[0] == '-') {
			val := ""
			for i < len(expression) && ((expression[i] >= '0' && expression[i] <= '9') || expression[i] == '.' || expression[i] == '-') {
				val += string(expression[i])
				i++
			}
			num, err := strconv.ParseFloat(val, 64)
			if err != nil {
				return nil, err
			}
			values = append(values, num)
			continue
		}

		if ch == "「" {
			token := ""
			i += size
			closed := false
			for i < len(expression) {
				tokenRune, tokenRuneSize := utf8.DecodeRuneInString(expression[i:])
				if tokenRune == utf8.RuneError && tokenRuneSize == 1 {
					return nil, fmt.Errorf("無法解析字串內容：%q", expression[i:i+1])
				}
				tokenCh := string(tokenRune)
				if tokenCh == "」" {
					i += tokenRuneSize
					closed = true
					break
				}
				token += tokenCh
				i += tokenRuneSize
			}
			if !closed {
				return nil, fmt.Errorf("字串缺少結尾引號：%s", expression)
			}
			values = append(values, token)
			continue
		}

		if strings.HasPrefix(expression[i:], "變數「") {
			i += len("變數「")
			token := ""
			closed := false
			for i < len(expression) {
				tokenRune, tokenRuneSize := utf8.DecodeRuneInString(expression[i:])
				if tokenRune == utf8.RuneError && tokenRuneSize == 1 {
					return nil, fmt.Errorf("無法解析變數名稱：%q", expression[i:i+1])
				}
				tokenCh := string(tokenRune)
				if tokenCh == "」" {
					i += tokenRuneSize
					closed = true
					break
				}
				token += tokenCh
				i += tokenRuneSize
			}
			if !closed {
				return nil, fmt.Errorf("變數名稱缺少結尾引號：%s", expression)
			}
			if val, ok := e.variables[token]; ok {
				values = append(values, val)
			} else {
				return nil, fmt.Errorf("未定義的變數：%s", token)
			}
			continue
		}

		if strings.HasPrefix(expression[i:], "陣列「") {
			i += len("陣列「")
			token := ""
			closed := false
			for i < len(expression) {
				tokenRune, tokenRuneSize := utf8.DecodeRuneInString(expression[i:])
				if tokenRune == utf8.RuneError && tokenRuneSize == 1 {
					return nil, fmt.Errorf("無法解析陣列名稱：%q", expression[i:i+1])
				}
				tokenCh := string(tokenRune)
				if tokenCh == "」" {
					i += tokenRuneSize
					closed = true
					break
				}
				token += tokenCh
				i += tokenRuneSize
			}
			if !closed {
				return nil, fmt.Errorf("陣列名稱缺少結尾引號：%s", expression)
			}

			if strings.HasPrefix(expression[i:], "的索引（") {
				i += len("的索引（")
				tmpExp := ""
				stk := 1
				for i < len(expression) {
					tmpRune, tmpRuneSize := utf8.DecodeRuneInString(expression[i:])
					if tmpRune == utf8.RuneError && tmpRuneSize == 1 {
						return nil, fmt.Errorf("無法解析陣列索引：%q", expression[i:i+1])
					}
					tmpCh := string(tmpRune)
					if tmpCh == "（" {
						stk++
					} else if tmpCh == "）" {
						stk--
					}
					tmpExp += tmpCh
					i += tmpRuneSize
					if stk == 0 {
						break
					}
				}
				if stk != 0 {
					return nil, fmt.Errorf("陣列索引缺少右括號：%s", expression)
				}
				indexVal, err := e.evaluate(strings.TrimSuffix(tmpExp, "）"))
				if err != nil {
					return nil, err
				}
				index, err := e.toFloat(indexVal)
				if err != nil {
					return nil, err
				}
				if arr, ok := e.arrays[token]; ok {
					idx := int(index)
					if idx < 0 || idx >= len(arr) {
						return nil, fmt.Errorf("陣列索引超出範圍：%d", idx)
					}
					values = append(values, arr[idx])
				} else {
					return nil, fmt.Errorf("未定義的陣列：%s", token)
				}
				continue
			} else if strings.HasPrefix(expression[i:], "的長度") {
				i += len("的長度")
				if arr, ok := e.arrays[token]; ok {
					values = append(values, float64(len(arr)))
				} else {
					return nil, fmt.Errorf("未定義的陣列：%s", token)
				}
				continue
			} else {
				return nil, fmt.Errorf("未知的算式：%s", expression)
			}
		}

		keywords := []string{"加上", "減去", "乘以", "除以", "取模", "等於", "大於", "小於", "轉型"}
		found := false
		for _, kw := range keywords {
			if i+len(kw) <= len(expression) && expression[i:i+len(kw)] == kw {
				for len(ops) > 0 && ops[len(ops)-1] != "（" && precedence(ops[len(ops)-1]) >= precedence(kw) {
					if len(values) < 2 {
						return nil, fmt.Errorf("運算元不足")
					}
					val2 := values[len(values)-1]
					val1 := values[len(values)-2]
					values = values[:len(values)-2]
					op := ops[len(ops)-1]
					ops = ops[:len(ops)-1]
					result, err := e.apply(val1, val2, op)
					if err != nil {
						return nil, err
					}
					values = append(values, result)
				}
				ops = append(ops, kw)
				i += len(kw)
				found = true
				break
			}
		}

		if !found {
			remainingRunes := []rune(expression[i:])
			if len(remainingRunes) > 10 {
				remainingRunes = remainingRunes[:10]
			}
			remaining := string(remainingRunes)
			return nil, fmt.Errorf("未知的運算子：%s", remaining)
		}
	}

	for len(ops) > 0 {
		if len(values) < 2 {
			return nil, fmt.Errorf("運算元不足")
		}
		val2 := values[len(values)-1]
		val1 := values[len(values)-2]
		values = values[:len(values)-2]
		op := ops[len(ops)-1]
		ops = ops[:len(ops)-1]
		result, err := e.apply(val1, val2, op)
		if err != nil {
			return nil, err
		}
		values = append(values, result)
	}

	if len(values) != 1 {
		return nil, fmt.Errorf("表達式錯誤")
	}

	final := values[0]
	if num, ok := final.(float64); ok && num == float64(int(num)) {
		return int(num), nil
	}
	return final, nil
}

func (e *Executor) assignVariable(args string) error {
	parts := strings.Split(args, "，")
	if len(parts) != 2 {
		return fmt.Errorf("變數賦值格式錯誤")
	}

	varName := strings.Trim(parts[0], "「」")
	value, err := e.evaluate(parts[1])
	if err != nil {
		return err
	}

	e.variables[varName] = value
	return nil
}

func (e *Executor) initArray(args string) error {
	parts := strings.Split(args, "，")
	if len(parts) != 2 {
		return fmt.Errorf("陣列初始化格式錯誤")
	}

	arrayName := strings.Trim(parts[0], "「」")
	sizeVal, err := e.evaluate(parts[1])
	if err != nil {
		return err
	}

	size, err := e.toFloat(sizeVal)
	if err != nil {
		return err
	}

	arr := make([]interface{}, int(size))
	for i := range arr {
		arr[i] = nil
	}

	e.arrays[arrayName] = arr
	return nil
}

func (e *Executor) assignArray(args string) error {
	parts := strings.Split(args, "，")
	if len(parts) != 3 {
		return fmt.Errorf("陣列賦值格式錯誤")
	}

	arrayName := strings.Trim(parts[0], "「」")
	indexVal, err := e.evaluate(parts[1])
	if err != nil {
		return err
	}

	index, err := e.toFloat(indexVal)
	if err != nil {
		return err
	}

	value, err := e.evaluate(parts[2])
	if err != nil {
		return err
	}

	if arr, ok := e.arrays[arrayName]; ok {
		idx := int(index)
		if idx < 0 || idx >= len(arr) {
			return fmt.Errorf("陣列索引超出範圍：%d", idx)
		}
		arr[idx] = value
		return nil
	}

	return fmt.Errorf("未定義的陣列：%s", arrayName)
}

func (e *Executor) getUserInput(args string) error {
	varName := strings.Trim(args, "「」")
	input := e.input()
	e.variables[varName] = strings.TrimSpace(input)
	return nil
}

func (e *Executor) splitArray(args string) error {
	parts := strings.Split(args, "，")
	if len(parts) != 2 {
		return fmt.Errorf("字串拆分格式錯誤")
	}

	strVal, err := e.evaluate(parts[0])
	if err != nil {
		return err
	}

	arrayName := strings.Trim(parts[1], "「」")
	str := fmt.Sprint(strVal)

	arr := make([]interface{}, len(str))
	for i, ch := range str {
		arr[i] = string(ch)
	}

	e.arrays[arrayName] = arr
	return nil
}

func (e *Executor) runCommand(command string) error {
	commands := map[string]func(string) error{
		"要求在終端機（或控制台）輸出": func(args string) error {
			val, err := e.evaluate(args)
			if err != nil {
				return err
			}
			e.output(fmt.Sprint(val))
			return nil
		},
		"設定一指定名稱變數的值，名稱與值分別為":                 e.assignVariable,
		"初始化一個指定名稱和長度的陣列，名稱和長度分別為":            e.initArray,
		"設定指定名稱陣列其中一項的值，名稱、項數與值分別為":           e.assignArray,
		"讀取使用者輸入字串，並儲存至變數":                    e.getUserInput,
		"將一字串拆分成字元陣列，並將結果儲存至指定名稱的陣列，字串和陣列分別為": e.splitArray,
	}

	for prefix, handler := range commands {
		if strings.HasPrefix(command, prefix) {
			return handler(command[len(prefix):])
		}
	}

	return fmt.Errorf("未知指令：%s", command)
}

func getClause(lines []string, startLine int) int {
	j := startLine
	stk := 1
	for stk > 0 {
		j++
		if j >= len(lines) {
			break
		}
		if strings.HasPrefix(lines[j], IF) || strings.HasPrefix(lines[j], WHILE) {
			stk++
		} else if lines[j] == END {
			stk--
		}
	}
	return j
}

func (e *Executor) execSeparateLines(lines []string) (interface{}, error) {
	i := 0
	for i < len(lines) {
		if e.shouldStop() {
			return 0, nil
		}

		if lines[i] == BREAK {
			return -1, nil
		}
		if lines[i] == CONTINUE {
			return -2, nil
		}
		if lines[i] == EXIT {
			return 0, nil
		}

		if strings.HasPrefix(lines[i], IF) {
			j := getClause(lines, i)
			condition := lines[i][len(IF):]
			if strings.HasPrefix(condition, "（") && strings.HasSuffix(condition, "）") {
				condition = condition[len("（") : len(condition)-len("）")]
			}
			condVal, err := e.evaluate(condition)
			if err != nil {
				return nil, err
			}

			condNum, err := e.toFloat(condVal)
			if err != nil {
				return nil, err
			}

			if condNum != 0 {
				ret, err := e.execSeparateLines(lines[i+1 : j])
				if err != nil {
					return nil, err
				}
				if ret != nil {
					return ret, nil
				}
			}
			i = j
		} else if strings.HasPrefix(lines[i], WHILE) {
			j := getClause(lines, i)
			condition := lines[i][len(WHILE):]
			if strings.HasPrefix(condition, "（") && strings.HasSuffix(condition, "）") {
				condition = condition[len("（") : len(condition)-len("）")]
			}

			for {
				if e.shouldStop() {
					return 0, nil
				}

				condVal, err := e.evaluate(condition)
				if err != nil {
					return nil, err
				}

				condNum, err := e.toFloat(condVal)
				if err != nil {
					return nil, err
				}

				if condNum == 0 {
					break
				}

				ret, err := e.execSeparateLines(lines[i+1 : j])
				if err != nil {
					return nil, err
				}

				if ret != nil {
					if retNum, ok := ret.(int); ok {
						if retNum == -1 {
							break
						}
						if retNum == 0 {
							return 0, nil
						}
					}
				}
			}
			i = j
		} else {
			if err := e.runCommand(lines[i]); err != nil {
				return nil, err
			}
		}
		i++
	}

	return nil, nil
}

func (e *Executor) Execute(code string) error {
	e.variables = make(map[string]interface{})
	e.arrays = make(map[string][]interface{})

	lines := strings.Split(code, "\n")
	var splitCode []string

	for _, line := range lines {
		stripped := strings.TrimSpace(line)
		if stripped == "" {
			continue
		}
		if !strings.HasSuffix(stripped, "。") {
			return fmt.Errorf("你忘了加句點。")
		}
		splitCode = append(splitCode, stripped[:len(stripped)-len("。")])
	}

	res, err := e.execSeparateLines(splitCode)
	if err != nil {
		return err
	}

	if res != nil {
		if retNum, ok := res.(int); ok && retNum != 0 {
			return fmt.Errorf("程式碼錯誤執行。")
		}
	}

	return nil
}
