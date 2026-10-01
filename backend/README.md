# ACcode Backend API

Go 語言實作的 ACcode 後端 API 伺服器。

## 功能

- 中文程式語言直譯器
- 並發執行隔離（Goroutine）
- 即時輸出串流（Channel）
- 使用者輸入支援
- Session 管理和自動清理

## API 端點

### POST /api/run
執行程式碼

**請求**：
```json
{
  "code": "要求在終端機（或控制台）輸出「Hello World」。"
}
```

**回應**：
```json
{
  "sessionId": "0"
}
```

### GET /api/status/:sessionId
查詢執行狀態

**回應**：
```json
{
  "output": "Hello World",
  "needsInput": false,
  "completed": true
}
```

### POST /api/input
提交使用者輸入

**請求**：
```json
{
  "sessionId": "0",
  "input": "使用者輸入的內容"
}
```

**回應**：
```json
{
  "message": "輸入已提交"
}
```

### POST /api/stop
停止程式執行

**請求**：
```json
{
  "sessionId": "0"
}
```

**回應**：
```json
{
  "message": "程式已成功終止"
}
```

## 執行

```bash
go run .
```

伺服器將在 `http://localhost:5001` 啟動。

## 開發

### 編譯
```bash
go build -o accode-backend.exe
```

### 測試
```bash
go test ./...
```

## CORS 設定

預設允許來自 `http://localhost:3000` 的請求（React 開發伺服器）。

生產環境請修改 `main.go` 中的 CORS 設定。
