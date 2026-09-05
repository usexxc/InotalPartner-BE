# 🔧 Alur Git — InotalPartner (Kerja Berdua)

> Proyek PKL INOTAL. Nazril pegang **Backend (Go)**, Fatan pegang **Frontend (React)**.
> Repo terpisah: `InotalPartner-BE` (Nazril) & `InotalPartner-FE` (Fatan) — titik temu cuma `API_CONTRACT.md`.

---

# 🟦 INI UNTUK NAZRIL (Backend — Go)

## Setup sekali aja

```bash
# Identitas git — udah ✅ (usexxc / nazrilmalik33@gmail.com)

# Login GitHub (biar boleh push) — paling gampang lewat gh CLI
sudo pacman -S github-cli
gh auth login
# → GitHub.com → HTTPS → Yes (authenticate Git) → Login with a web browser
# → muncul kode 8 huruf → buka github.com/login/device → ketik kodenya → Authorize

# Clone (udah ✅ di ~/Projects/inotalpartner/)
git clone https://github.com/BiZa7/InotalPartner-BE.git
```

Minta **Fatan** add lo jadi collaborator:
GitHub → repo BE → *Settings → Collaborators → Add people* → `usexxc`
→ lo terima invite via email → **Accept**.
*(Belum di-invite? Pakai Fork dulu — lihat bagian bawah.)*

## Tugas pertama lo

Commit kontrak + panduan ke repo BE:

```bash
git pull origin main
git checkout -b docs/setup
git add GIT_WORKFLOW.md API_CONTRACT.md
git commit -m "docs: git workflow + API contract v1"
git push -u origin docs/setup
# GitHub → tombol "Compare & pull request" → minta Fatan review → Merge!
```

## Kerjaan lo: CRUD /api/partners (urutan per file)

Tiru pola auth yang udah ada — alur wajibnya:

```
model/partner.go            → struct GORM (kolom = skema di API_CONTRACT)
repository/partner.go       → Create / Find / Update / Delete
service/partner_service.go  → DTO request/response + validasi (name, category, pic_name wajib)
handler/partner_handler.go  → baca JSON → panggil service → bungkus {success, message, data}
routes/routes.go            → daftarkan di group protected (POST/PUT/DELETE + AdminOnly)
```

Pengerjaan berurutan: **GET list (dgn search+pagination) → GET /:id → POST → PUT → DELETE**
Setiap 1 endpoint beres = commit + push. Jangan nunggu komplit semua.

## Alur harian (Level 1 — hapalin)

```bash
git pull origin main      # 1. ambil terbaru — SELALU sebelum mulai
                          # 2. ngoding...
git add .                 # 3. masukin perubahan
git commit -m "feat: endpoint list partner"
git push                  # 4. kirim ke GitHub
```

VS Code: panel **Source Control** (`Ctrl+Shift+G`) — `+`=add, pesan → ✓=commit, ↑=push, ↓=pull.

---

# 🟩 INI UNTUK FATAN (Frontend — React)

## Setup sekali aja

```bash
git config --global user.name  "username-github-lo"
git config --global user.email "email-github-lo"

sudo pacman -S github-cli
gh auth login
# → GitHub.com → HTTPS → Yes → Login with a web browser → ikutin kodenya

git clone https://github.com/BiZa7/InotalPartner-FE.git
cd InotalPartner-FE
npm install        # pasang dependencies
npm run dev        # jalan di http://localhost:5173
```

Repo FE punya Fatan (owner) — bebas push.
Tugas tambahan: add Nazril jadi collaborator (Settings → Collaborators → `usexxc`).

## Tugas pertama Fatan

Masukin kontrak API ke repo FE:

```bash
git pull origin main
git add API_CONTRACT.md
git commit -m "docs: tambah API contract v1 dari BE"
git push
```

## Kerjaan Fatan: halaman Partners (urutan)

1. **Halaman list** — tabel + search + pagination, mulai pakai **mock data** yang bentuknya
   PERSIS contoh response di `API_CONTRACT.md` (`items: [...]` + `pagination`)
2. **Form tambah/edit** — field sesuai skema tabel partner di kontrak
3. Pas backend Nazril jadi → ganti mock dengan axios:
   ```js
   import api from '../api/axios'
   const res = await api.get('/partners', { params: { search, page, limit } })
   // res.data.data = { items, pagination } — dijamin sesuai kontrak
   ```
4. Endpoint create/update/delete itu **admin-only** (403 kalau role partner) →
   sembunyikan tombol kalau `user.role !== 'admin'`

Dashboard yang masih dummy boleh dibiarkan dulu — fokus halaman Partners.

**PENTING:** kalau bentuk data di kontrak kurang pas buat UI → JANGAN ubah sendiri di kode.
Chat Nazril → update `API_CONTRACT.md` bareng → commit → pull.

## Alur harian (5 langkah, sama kayak Nazril)

```bash
git pull origin main      # 1. ambil terbaru — SELALU sebelum mulai
                          # 2. ngoding...
git add .                 # 3. masukin perubahan
git commit -m "feat: halaman list partner + mock data"
git push                  # 4. kirim ke GitHub
```

Pesan commit konsisten: `feat:` fitur, `fix:` bug, `docs:` dokumen, `style:` tampilan.

---

# 🤝 YANG SAMA UNTUK DUA-DUANYA

1. **Jangan pernah commit langsung ke `main`** setelah era cabang dimulai (awal-awal boleh, level 1 dulu)
2. **`git pull` sebelum mulai & sebelum push** — 90% masalah hilang dengan kebiasaan ini
3. **Commit kecil-kecil** = titik save; 1 commit = 1 hal yang beres
4. **Review PR silang** — GitHub kasi notifikasi "*...wants to merge*", baca diff (hijau=tambah, merah=hapus), komentarin, klik Merge kalau oke
5. **Conflict?** → `git pull` dulu → buka file yang ada `<<<<<<< HEAD` → pilih/gabungin yang bener → hapus penanda → add → commit. VS Code ada tombol **Accept Current / Incoming / Both**. Jangan pernah `git push --force`
6. **Ubah `API_CONTRACT.md` = diskusi dulu**, baru commit — itu kesepakatan dua pihak

## Level 2 — cabang & PR (pas udah nyaman)

```bash
git checkout -b feat/partner-model   # bikin cabang
# ...ngoding, add, commit...
git push -u origin feat/partner-model
# GitHub → "Compare & pull request" → partner review → Merge
```

## Kalau Nazril belum di-invite (Fork, sementara)

```bash
# GitHub: tombol "Fork" di repo BiZa7/InotalPartner-BE → jadi repo usexxc
git clone https://github.com/usexxc/InotalPartner-BE.git
cd InotalPartner-BE
git remote add upstream https://github.com/BiZa7/InotalPartner-BE.git
# narik update dari repo asli:
git fetch upstream && git merge upstream/main
# ngajuin kerjaan: push ke fork → PR ke repo asli
```

---

# ✅ Checklist Mulai

**Nazril:**
- [ ] `gh auth login`
- [ ] Accept invite collaborator dari Fatan (atau fork dulu)
- [ ] Commit `API_CONTRACT.md` ke repo BE → PR → minta Fatan merge
- [ ] Mulai `model/partner.go`

**Fatan:**
- [ ] `gh auth login`
- [ ] Clone FE + `npm install` + `npm run dev`
- [ ] Commit `API_CONTRACT.md` ke repo FE
- [ ] Mock data halaman list Partners (bentuk = sesuai kontrak)
