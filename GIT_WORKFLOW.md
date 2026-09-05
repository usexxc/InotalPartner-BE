# 🔧 Setup Git — InotalPartner (Kerja Berdua)

> Panduan biar Nazril (BE/Go) dan temen (FE/React) nyambung lewat GitHub.
> Ikuti urutannya: **Bagian 1 sekali aja**, Bagian 3 tiap hari.

---

## 1. Setup Sekali Aja

### 1a. Identitas git (udah ✅)

```bash
git config --global user.name  "usexxc"
git config --global user.email "nazrilmalik33@gmail.com"
```
> Email ini harus sama dengan yang terdaftar di akun GitHub lo — biar commit lo kehitung sebagai kontribusi lo di profile GitHub.

### 1b. Kredensial GitHub (biar boleh push)

Cara termudah: **gh CLI** (login sekali, beres selamanya).

```bash
# Pasang (sekali)
sudo pacman -S github-cli

# Login — akan muncul KODE 8 huruf
gh auth login
# jawab: GitHub.com → HTTPS → Yes (authenticate Git) → Login with a web browser
# buka github.com/login/device → ketik kodenya → Authorize
```

**Alternatif tanpa install:** pakai *Personal Access Token* (PAT):
1. GitHub → foto profil → **Settings** → **Developer settings** (paling bawah) → **Personal access tokens → Tokens (classic)** → *Generate new token (classic)*
2. Centang scope **`repo`** → Generate → **copy tokennya** (cuma muncul 1x!)
3. Set credential helper sekali:
   ```bash
   git config --global credential.helper store
   ```
4. Pas pertama `git push`: username = `usexxc`, password = **token** (bukan password GitHub). Tersimpan, gak ditanya lagi.

### 1c. Akses repo temen (collaborator)

- **Yang punya repo (BiZa7):** GitHub → repo → *Settings → Collaborators → Add people* → ketik `usexxc` → invite
- **Lo:** cek email / notifikasi GitHub → **Accept invitation**
- Status: bisa clone, push, buat PR di repo asli. (Kalau belum di-invite, sementara pakai **Fork** — lihat catatan di bawah.)

### 1d. Clone (lo udah ✅ di `~/Projects/inotalpartner/`)

```bash
git clone https://github.com/BiZa7/InotalPartner-BE.git
git clone https://github.com/BiZa7/InotalPartner-FE.git
```

---

## 2. Pembagian & Aturan Main

| | Nazril | Temen |
|---|---|---|
| Repo | `InotalPartner-BE` | `InotalPartner-FE` |
| Isi | Go: model → repo → service → handler → routes | React: pages, components, api |
| Kontrak bersama | `API_CONTRACT.md` | `API_CONTRACT.md` |

1. **Repo BE & FE terpisah** → hampir mustahil konflik. Titik temu cuma `API_CONTRACT.md`
2. **Ubah kontrak = commit di repo masing-masing + kabari** (jangan edit diam-diam)
3. **Commit kecil-kecil** — 1 commit = 1 hal yang beres
4. Pesan commit: `feat: endpoint list partner`, `fix: validasi email duplikat`, `docs: API contract v1`

---

## 3. Alur Harian (hapalin — Level 1)

```bash
git pull origin main      # 1. AMBIL versi terbaru (selalu sebelum mulai)
                          # 2. NGODING...
git add .                 # 3. Masukin semua perubahan ke "keranjang"
git commit -m "feat: ..." # 4. Simpan titik restore
git push                  # 5. Kirim ke GitHub
```

> **VS Code bisa semua ini klik-klik:** panel Source Control (`Ctrl+Shift+G`) — `+` = add, ketik pesan → ✓ = commit, tombol ⚠️/↓ = pull, ↑ = push.

---

## 4. Level 2 — Cabang & PR (pas udah nyaman)

```bash
git checkout -b feat/partner-model   # bikin cabang
# ...ngoding, add, commit...
git push -u origin feat/partner-model
# GitHub → tombol "Compare & pull request" → temen review → Merge
```

Cocok dipake pas lo mulai sering eksperimen (cabang = dunia sandbox, main = dunia resmi).

---

## 5. Darurat: Conflict

Kalau push ditolak → `git pull` dulu. Kalau muncul tanda konflik di file:
```
<<<<<<< HEAD
kode lo
=======
kode dia
>>>>>>> ...
```
Buka file → pilih/gabungin yang bener → hapus penanda → save → `git add` → commit. Di VS Code ada tombol **Accept Current / Incoming / Both**.

Jangan pernah: `git push --force` (menghapus kerjaan orang lain).

---

## 6. Catatan Fork (hanya kalau belum di-invite)

```bash
# di GitHub: tombol "Fork" di repo BiZa7 → jadi repo lo
git clone https://github.com/usexxc/InotalPartner-BE.git
cd InotalPartner-BE
git remote add upstream https://github.com/BiZa7/InotalPartner-BE.git
# narik update dari repo asli:
git fetch upstream && git merge upstream/main
# ngajuin kerjaan: push ke fork → PR ke repo asli
```

---

## 7. Checklist Mulai

- [x] Identitas git
- [ ] `gh auth login` atau PAT
- [ ] Accept invite collaborator (atau fork dulu)
- [ ] Commit pertama lo: masukin `GIT_WORKFLOW.md` + `API_CONTRACT.md` ke repo BE
  ```bash
  git pull origin main
  git checkout -b docs/setup
  git add GIT_WORKFLOW.md API_CONTRACT.md
  git commit -m "docs: git workflow + API contract v1"
  git push -u origin docs/setup
  # GitHub → Compare & pull request → minta temen review → merge!
  ```
- [ ] Kirim `API_CONTRACT.md` ke temen (WA) + dia commit juga ke repo FE-nya
