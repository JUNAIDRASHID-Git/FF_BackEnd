# FF_BackEnd

FunFillers Go Backend Service.

## Tech Stack & Features
- **Go / Gin Framework**: High-performance HTTP REST API
- **GORM / PostgreSQL**: Database ORM with dynamic fallback to JSON store
- **Open Graph (OG) Link Preview**: `/share/product/:id` for social media share previews (WhatsApp, Twitter, Facebook)
- **JWT Auth & Role Control**: Customer and Admin endpoints

## Setup & Running
```bash
go run main.go
```
The backend listens on port `5050` by default. Health check endpoint: `http://localhost:5050/api/health`.
