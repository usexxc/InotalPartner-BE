# 📋 InotalPartner — API Contract v1

> **Dokumen kesepakatan Backend ↔ Frontend.**
> BE yang implement, FE yang konsumsi. Kalau mau ubah bentuk response — ubah dulu di sini, jangan langsung di kode.
>
> Base URL (dev): `http://localhost:8080/api` — FE sudah di-proxy lewat Vite (`/api` → `:8080`), jadi FE cukup panggil `/api/...`

---

## 🔑 Konvensi Global

### Response envelope (WAJIB, semua endpoint)

```json
// Sukses
{ "success": true, "message": "Login berhasil", "data": { ... } }

// Gagal
{ "success": false, "message": "Email sudah terdaftar" }
// Validasi gagal tambah field:
{ "success": false, "message": "Data tidak valid", "errors": "Key: 'xxx' Error:..." }
```

### Status code yang dipakai
| Kode | Arti |
|------|------|
| 200 | OK (login, get, update) |
| 201 | Created (register, create partner) |
| 400 | Payload tidak valid / validasi gagal |
| 401 | Belum login / token invalid-expired |
| 403 | Login tapi bukan admin (endpoint admin-only) |
| 404 | Data tidak ditemukan |
| 409 | Conflict (email/username sudah ada) |

### Autentikasi
- Semua endpoint kecuali `auth/*` wajib header: `Authorization: Bearer <token>` (sudah dihandle interceptor axios FE)
- Token expired → BE balas 401 → interceptor FE auto-logout. Jangan ubah perilaku ini.

---

## 1️⃣ AUTH — sudah jadi ✅ (referensi, jangan diubah)

### POST /auth/register
```json
// request
{ "full_name": "Nazril", "company": "INOTAL", "email": "n@mail.com", "password": "min8karakter" }
// 201 → data: { "token": "eyJ...", "user": { "id":1, "full_name":"...", "email":"...", "company":"...", "avatar_url":"", "role":"partner", "provider":"local" } }
```

### POST /auth/login
```json
// request
{ "email": "n@mail.com", "password": "..." }
// 200 → data: { token, user } — bentuk sama dengan register
// 401 → { "success": false, "message": "Email atau password salah" }
```

### GET /auth/me  (protected)
```json
// 200 → data: { "id":1, "full_name":"...", ... }  (bentuk PublicUser)
```

### GET /auth/google → redirect ke Google
### GET /auth/google/callback → handle callback, redirect balik ke FRONTEND_URL

---

## 2️⃣ PARTNERS — 🔨 yang dibangun (tugas BE: Nazril)

### Model data `partners` (GORM, AutoMigrate)

| Kolom | Tipe | Rules |
|-------|------|-------|
| id | uint PK | auto |
| name | varchar(150) | required, unique |
| logo_url | text | opsional |
| category | varchar(50) | enum: `client` / `vendor` / `reseller` / `affiliate` |
| industry | varchar(100) | opsional, mis. "Fintech" |
| pic_name | varchar(150) | required (nama orang di partner) |
| pic_email | varchar(255) | valid email kalau diisi |
| pic_phone | varchar(30) | opsional |
| address | text | opsional |
| status | varchar(20) | enum: `active` / `prospect` / `inactive`, default `prospect` |
| notes | text | opsional |
| created_at / updated_at / deleted_at | | standar GORM (soft delete) |

### GET /partners  (protected, semua role)
Query params (semua opsional): `search` (cocokkan name/industry/pic_name, ILIKE), `status`, `category`, `page` (default 1), `limit` (default 10, max 100)

```json
// 200 → data:
{
  "items": [
    { "id": 5, "name": "PT Maju Bersama", "logo_url": "", "category": "client",
      "industry": "Logistik", "pic_name": "Budi", "pic_email": "budi@maju.id",
      "pic_phone": "0812xxx", "address": "Bandung", "status": "active",
      "notes": "", "created_at": "2026-09-05T10:00:00Z", "updated_at": "..." }
  ],
  "pagination": { "page": 1, "limit": 10, "total": 27, "total_pages": 3 }
}
```

### GET /partners/:id  (protected)
```json
// 200 → data: { ...satu objek partner, bentuk sama seperti items[0] }
// 404 → { "success": false, "message": "Partner tidak ditemukan" }
```

### POST /partners  (admin only → 403 kalau role partner)
```json
// request — name, category, pic_name WAJIB
{ "name": "PT Baru Jaya", "category": "vendor", "industry": "Cloud",
  "pic_name": "Sinta", "pic_email": "sinta@baru.id", "pic_phone": "0813xxx",
  "address": "Jakarta", "status": "prospect", "notes": "kontak dari pameran" }
// 201 → data: { objek partner lengkap }
// 409 → name sudah dipakai
```

### PUT /partners/:id  (admin only)
```json
// request: field mana pun boleh dikirim sebagian (partial update)
{ "status": "active", "notes": "MoU ditandatangani 5 Sep" }
// 200 → data: { objek partner terbaru }
```

### DELETE /partners/:id  (admin only)
```json
// 200 → { "success": true, "message": "Partner dihapus" }  (soft delete)
```

---

## 3️⃣ PROJECTS & ACTIVITIES — v2, menyusul (jangan dibangun dulu)

> FE: untuk sekarang, dashboard boleh tetap pakai dummy. Fokus kontrak ini dulu.

---

## ✍️ Catatan kesepakatan
- [ ] BE: selesaikan GET list + GET by id dulu (FE bisa mulai halaman list), lalu POST/PUT/DELETE
- [ ] FE: halaman `Partners` pakai bentuk data persis seperti di atas — biar pas API jadi tinggal ganti mock → axios call
- [ ] Perubahan apapun di dokumen ini diskusi dulu di grup WA
