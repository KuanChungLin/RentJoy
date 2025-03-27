# RentJoy

RentJoy 是一個場地租借平台，讓使用者可以輕鬆找到並預訂適合的場地。

## 功能特點

- 場地搜尋與預訂
- 社群媒體登入（Facebook、Google）
- 場地管理系統
- 即時訂單處理
- 金流整合（ECPay）
- 圖片上傳（Cloudinary）
- 地圖整合（Google Maps）

## 技術棧

- 後端：Go
- 資料庫：SQL Server
- 快取：Redis
- 認證：JWT
- 圖片存儲：Cloudinary
- 金流：ECPay
- 地圖：Google Maps API

## 環境要求

- Go 1.21 或更高版本
- SQL Server
- Redis
- SSL 證書（用於 HTTPS）

## 安裝步驟

1. 克隆專案
```bash
git clone https://github.com/yourusername/RentJoy.git
cd RentJoy
```

2. 安裝依賴
```bash
go mod download
```

3. 設定環境變數
```bash
cp .env.example .env
```
編輯 `.env` 文件，填入必要的環境變數值。

4. 運行專案
```bash
go run cmd/app/main.go
```

## 環境變數說明

專案需要設定以下環境變數：

### 資料庫設定
- `DB_USER`: 資料庫使用者名稱
- `DB_PASSWORD`: 資料庫密碼
- `DB_HOST`: 資料庫主機地址
- `DB_PORT`: 資料庫端口
- `DB_NAME`: 資料庫名稱

### Redis 設定
- `REDIS_ADDR`: Redis 服務器地址
- `REDIS_PASSWORD`: Redis 密碼（如果有）
- `REDIS_DB`: Redis 資料庫編號

### OAuth 設定
- `FACEBOOK_APP_ID`: Facebook 應用程式 ID
- `FACEBOOK_APP_SECRET`: Facebook 應用程式密鑰
- `GOOGLE_CLIENT_ID`: Google 客戶端 ID
- `GOOGLE_CLIENT_SECRET`: Google 客戶端密鑰

### 其他服務設定
- `CLOUDINARY_URL`: Cloudinary 服務 URL
- `GOOGLE_MAPS_API_KEY`: Google Maps API 金鑰
- `ECPAY_MERCHANT_ID`: ECPay 商店編號
- `ECPAY_HASH_KEY`: ECPay HashKey
- `ECPAY_HASH_IV`: ECPay HashIV

## 開發團隊

- [您的名字] - 主要開發者

## 授權

本專案採用 MIT 授權條款 - 詳見 [LICENSE](LICENSE) 文件
