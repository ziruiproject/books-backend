
# Starter Project GO — Hexagonal Architecture (Ports & Adapters)

Starter adalah template atau boilerplate project Golang yang menggunakan pendekatan **Hexagonal Architecture** (Ports & Adapters).

## 🚀 Fitur Utama

- ✨ Arsitektur Hexagonal (clean, scalable)
- 🌐 HTTP dan ~~gRPC~~ endpoint (Fiber + gRPC ready)
- 🔐 Auth dengan JWT
- 🧩 Middleware (auth, logging)
- 🧪 Validasi input
- 🗃️ PostgreSQL ready
- 🧰 Modular dan ~~testable~~

---

## 📁 Struktur Project

```
.
├── cmd/                # Entry point aplikasi (bootstrap & main)
│   └── app/
├── config/             # Loader konfigurasi dari env
├── internal/           # Core + Adapter
│   ├── adapters/       # Adapter input/output (http, grpc, db, jwt, validator, dll)
│   └── core/           # Domain + usecase (business logic)
├── pkg/                # Helper dan library kecil
├── docker-compose.yaml
├── go.mod / go.sum
└── README.md
````

---

## 📦 Dependencies

- [Fiber](https://github.com/gofiber/fiber) — HTTP framework
- [gorm](https://gorm.io/gorm) — Database driver
- [zerolog](https://github.com/rs/zerolog) — Structured logger
- [jwt-go](https://github.com/golang-jwt/jwt) — JWT auth
- [go-playground/validator](https://github.com/go-playground/validator) — Input validation

---

## ⚙️ Menjalankan Project

### 1. Clone dan install dependencies

```bash
git clone https://github.com/ziruiproject/starter-go.git
cd starter-go
go mod tidy
````

### 2. Copy & edit konfigurasi

```bash
cp config/example.env .env
```

Isi variabel `.env`:

```env
APP_URL=http://localhost
APP_PORT=8080
TIMEZONE=Asia/Jakarta

POSTGRES_DB=postgres
POSTGRES_PORT=5432
POSTGRES_USER=postgres
POSTGRES_HOST=localhost
POSTGRES_PASSWORD=password123

JWT_SECRET=my_super_secret_key

```

### 3. Jalankan PostgreSQL

Jika ingin memakai Docker:

```bash
docker-compose up -d
```

### 4. Jalankan aplikasi

```bash
go run cmd/app/main.go cmd/app/bootstrap.go
```

[//]: # (---)

[//]: # ()
[//]: # (## 🧪 Menjalankan Testing)

[//]: # ()
[//]: # (```bash)

[//]: # (go test ./...)

[//]: # (```)

---

## 📌 TODO

* [ ] Memberikan komentar dengan godoc
* [ ] Unit test untuk usecase
* [ ] Integrasi object storage
* [ ] Tambah gRPC
* [ ] Deployment script
* [ ] Containerization dengan Dockerfile

---

## 👤 Author

Built by [Yudha Sugiharto](https://github.com/ziruiproject) — feel free to fork or contribute ✨

---

## 📄 License

MIT License — bebas digunakan untuk project pribadi maupun komersial.
