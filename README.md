Aquanote
溫度TDS監控系統
硬體ESP32-S3 + TDS BOARD V1.0 + DS18B20 + TDS探針

後端專案結構
aquanote-backend/
├── cmd/
│   └── server/
│       └── main.go          # 應用程式入口
├── internal/
│   ├── handler/             # HTTP 請求處理（即 Controller）
│   │   └── user_handler.go
│   ├── service/             # 業務邏輯層
│   │   └── user_service.go
│   ├── repository/          # 資料存取層（DB 操作）
│   │   └── user_repo.go
│   ├── model/               # 資料結構 / DB Schema
│   │   └── user.go
│   ├── middleware/          # 自訂 Middleware（JWT, Logging...）
│   │   └── auth.go
│   └── router/              # 路由設定
│       └── router.go
├── pkg/
│   ├── config/              # 設定檔讀取（env, yaml）
│   │   └── config.go
│   └── utils/               # 共用工具函式
│       └── response.go
├── configs/                 # 設定檔（.yaml, .env）
│   └── app.yaml
├── docs/                    # API 文件（Swagger 等）
├── go.mod
├── go.sum
└── README.md

前端專案結構 Vite React