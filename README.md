# ACcode - 中文程式語言直譯器

前後端分離架構的 ACcode 專案。

## 專案結構

```
accode-go/
├── backend/          # Go API 伺服器
│   ├── main.go
│   ├── handlers/
│   ├── executor/
│   ├── models/
│   └── go.mod
└── frontend/         # React 前端應用
    ├── src/
    ├── public/
    └── package.json
```

## 架構設計

### 後端（Go）
- **職責**：程式碼執行引擎、進程管理、安全隔離
- **技術**：Gin 框架、Goroutine、Channel
- **API 端點**：
  - `POST /api/run` - 執行程式碼
  - `GET /api/status/:sessionId` - 查詢執行狀態
  - `POST /api/input` - 提交使用者輸入
  - `POST /api/stop` - 停止執行

### 前端（React）
- **職責**：UI 邏輯、程式碼編輯、輸出顯示、使用者互動
- **技術**：React、Fetch API
- **功能**：
  - 程式碼編輯器
  - 即時輸出顯示
  - 輸入框管理
  - 執行控制

## 開發指南

### 後端開發

```bash
cd backend
go run .
```

後端將在 `http://localhost:5001` 啟動。

### 前端開發

```bash
cd frontend
npm install
npm start
```

前端將在 `http://localhost:3000` 啟動。

## API 文件

詳見 `backend/README.md`

## 授權

與原 Python 版本相同。
