# 📚 Book Catalog API

Selamat datang di Book Catalog API! 🚀  
Project ini adalah backend service untuk mengelola **katalog buku digital** lengkap dengan fitur-fitur modern seperti:

- 🔐 Autentikasi JWT
- 📦 Upload file ke S3 (via MinIO)
- 🗃️ CRUD Buku, Penulis, Kategori, Penerbit, Pengguna
- 🧠 Validasi, Filter, Pagination
- 🐳 Docker + Seeder Otomatis
- 🧱 Dibangun dengan **Hexagonal Architecture**

> Dibangun dari template starter milikku sendiri:  
> 👉 [**starter-go**](https://github.com/ziruiproject/starter-go)

---

## 🧠 Tech Stack

| Layer        | Tools/Libs                 |
|--------------|----------------------------|
| Language     | Go 🧬                      |
| HTTP Server  | Fiber ⚡                    |
| Auth         | JWT 🔐                     |
| Database     | PostgreSQL 🐘              |
| Storage      | MinIO (S3 compatible) ☁️   |
| ORM          | Gorm ⚙️                    |
| Arch Design  | Hexagonal Architecture 🛠️ |
| Container    | Docker + Docker Compose 🐳 |

---

## 🧾 Struktur Project (Hexagonal Style)

```

📦 project-root
├── cmd/                  # Entry point
├── config/               # Konfigurasi & .env
├── internal/
│   ├── adapters/         # Infrastruktur (handler, repo, middleware, etc.)
│   └── core/             # Domain (entity, usecase, dto, interface)
├── pkg/                  # Utilitas (hasher, helper)
├── Makefile              # Shortcut command
├── docker-compose.yaml   # Container orchestration
└── README.md             # You're here 😎

````

---

## ⚙️ Konfigurasi `.env`

Salin file contoh:

```bash
cp config/example.env .env
````

Isi dengan konfigurasi seperti berikut:

```env
APP_URL=http://localhost
APP_PORT=8000
TIMEZONE=Asia/Jakarta

# PostgreSQL
POSTGRES_DB=postgres
POSTGRES_PORT=5432
POSTGRES_USER=postgres
POSTGRES_HOST=localhost
POSTGRES_PASSWORD=password123

# MinIO (S3-compatible)
MINIO_ROOT_USER=minioadmin
MINIO_ROOT_PASSWORD=minioadmin
MINIO_PORT_API=9000
MINIO_PORT_UI=9001
MINIO_BUCKET=bucket
MINIO_HOST=localhost
MINIO_REGION=us-west-2

# JWT Secret
JWT_SECRET=my_super_secret_key
```

---

## 🚀 Cara Menjalankan (Dev Mode)

### 1️⃣ Clone dan masuk ke direktori project

```bash
git clone https://github.com/kamu/projectmu.git
cd projectmu
```

### 2️⃣ Install dependency Go

```bash
go mod tidy
```

### 3️⃣ Jalankan semua service + seeder otomatis

```bash
make up-and-seed
```

➡️ Ini akan:

* 🐳 Menjalankan PostgreSQL & MinIO
* 🏗️ Build dan jalankan service Go
* 🌱 Jalankan migrasi & seeder

### 4️⃣ Akses layanan

| Service      | URL                                            |
| ------------ | ---------------------------------------------- |
| API Server   | [http://localhost:8000](http://localhost:8000) |
| MinIO UI     | [http://localhost:9001](http://localhost:9001) |
| MinIO Bucket | `bucket` (akses via API)                       |

MinIO Login:

```txt
Username: minioadmin
Password: minioadmin
```

---

## ☁️ Upload File ke S3 via MinIO

* Gunakan endpoint upload API
* File akan dikirim ke bucket S3 (`bucket`)
* Response mengembalikan URL akses file
* Menggunakan SDK: `https://github.com/aws/aws-sdk-go-v2`

---

## 🧪 Seed & Sample Data

Seeder disimpan di:

```
internal/adapters/database/seeder/
```

Contoh file:

* `202508020859_user_seeder.sql`
* `202508051233_book_seeder.sql`

Seeder akan dieksekusi otomatis saat kamu menjalankan:

```bash
make up-and-seed
```

---

## 💡 Build Manual (Tanpa Docker)

Kalau kamu prefer jalanin manual:

```bash
go mod tidy
go build -o bin/app cmd/app/main.go
./bin/app
```

Pastikan PostgreSQL & MinIO aktif yaa! 🔥

---

## 🔌 Impor API ke Postman / Insomnia

Project ini menyertakan (`api-specs.yaml`) yang diekspor langsung dari Insomnia. Kamu bisa menggunakannya untuk testing dan eksplorasi endpoint secara cepat.
### 🧪 Postman

1. Buka Postman
2. Klik `Import`
3. Pilih `File` → pilih `api-specs.yaml`
4. Done! Kamu bisa langsung eksplorasi endpointnya

### 🛌 Insomnia

1. Buka Insomnia
2. Klik `Create → Import`
3. Pilih `From File` → `api-specs.yaml`
4. Mulai testing 😴

>  `api-specs.yaml` berada di root direktori.

---

## 📌 TODO / Fitur Selanjutnya

* [ ] 🔍 Integrasi Swagger UI endpoint
* [ ] 🧪 Unit & Integration Tests
* [ ] 🚦 CI/CD Pipeline

---

## 🧑‍💻 Kontribusi

Kalau kamu tertarik pakai atau kontribusi:

* Fork repo ini 🍴
* Bikin fitur baru? PR welcome!
* Jangan lupa kasih ⭐ di sini dan di  [starter-go](https://github.com/ziruiproject/starter-go)

---

## 📬 Kontak

Buat diskusi, kolaborasi, atau sekadar ngopi ☕
📧 DM via GitHub atau buka issue ya!

---

> Dibuat dengan ❤️ oleh [Yudha Sugiharto](https://github.com/ziruiproject)