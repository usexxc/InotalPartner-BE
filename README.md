# InotalPartner-BE

Backend **InotalPartner** — platform kelola partner INOTAL.
Go + Gin + GORM + PostgreSQL.

## Menjalankan (development)

```bash
# 1. Siapkan database Postgres, lalu:
createdb inotal_partner        # atau lewat pgAdmin/DBeaver

# 2. Salin env & isi
cp .env.example .env           # edit JWT_SECRET minimal

# 3. Jalan
go mod tidy
go run ./cmd
# → [server] Starting on http://localhost:8080

# 4. Cek
curl http://localhost:8080/health
# {"status":"ok","service":"inotal-be"}
```

Tabel dibuat otomatis saat start (GORM AutoMigrate).

## Struktur

```
cmd/main.go                  entry point — load config, konek DB, dependency injection
config/config.go             baca .env
internal/model/              struct GORM (user.go, partner.go ...)
internal/repository/         query DB (satu repo per model)
internal/service/            logika bisnis + DTO (request/response)
internal/handler/            HTTP handler (baca request → panggil service)
internal/middleware/         AuthMiddleware (JWT), AdminOnly (role)
routes/routes.go             daftar semua endpoint
```

**Pola untuk fitur baru** (contoh: partners) — tiru alur auth:
1. `model/partner.go` → struct + AutoMigrate di `database.Connect`
2. `repository/partner_repository.go` → Create/Find/Update/Delete
3. `service/partner_service.go` → DTO + validasi + aturan bisnis
4. `handler/partner_handler.go` → baca JSON → panggil service → bungkus envelope
5. Daftarkan di `routes/routes.go` (group protected, pakai middleware)

## Kontrak API

Lihat **API_CONTRACT.md** — sumber kebenaran bentuk request/response.
Jangan ubah bentuk response tanpa update dokumen itu (FE ikut kena).
