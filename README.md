# Golang TodoList App

Aplikasi **TodoList** sederhana berbasis **Golang + Vue** dengan arsitektur **fullstack** (backend REST API dan frontend SPA). Proyek ini ditujukan sebagai latihan membangun aplikasi end-to-end dengan stack modern, sekaligus bisa dikembangkan lebih lanjut untuk kebutuhan produksi.

---

## 🚀 Stack yang Digunakan

### Backend

* **Golang** 1.24+
* **Gin** – Web framework untuk REST API
* **GORM** – ORM untuk akses database
* **MySQL** – Database relasional
* **Zap** – Structured Logging (logging JSON, lebih efisien & siap produksi)
* **GoDotEnv** – Load konfigurasi `.env`
* **Rate Limit Middleware** (custom, pakai `golang.org/x/time/rate`)
* **Docker Compose** – Menjalankan MySQL + API dalam container (opsional)

### Frontend

* **Vue 3** (SPA, Composition API, Vite bundler)
* **Vuetify 3** – UI Component Framework (sementara menggantikan TDesign)
* **Pinia** – State management
* **Vue Router** – Routing SPA
* **Axios** – HTTP client untuk komunikasi dengan API
* **Vite** – Build tool cepat, mendukung mode dev/prod
* **pnpm** – Package manager frontend

---

## ✨ Fitur yang Tersedia

### Backend

* CRUD Todo (Create, Read, Update, Delete)
* REST API dengan respons JSON
* Middleware:

  * **CORS** – support konfigurasi domain tertentu via `.env`
  * **Rate Limiting** – membatasi request per IP
  * **Structured Logging** dengan `zap` (request log JSON dengan method, path, status, latency, IP, user-agent)
* Health check endpoint: `/healthz`

### Frontend

* Tampilan TodoList responsif dengan **Vuetify**
* Tambah task baru
* Tandai task selesai (checkbox → strike-through)
* Edit task
* Hapus task dengan konfirmasi
* Skeleton loading state (next)
* Empty state (jika belum ada task) (next)
* Error handling sederhana (alert) (next)

---

## 📂 Struktur Proyek

```
my-todolist/
  backend/        # API Golang
    main.go
    db/           # koneksi & migrasi DB
    models/       # definisi entity (Todo)
    handlers/     # logika CRUD
    routes/       # definisi route
    middleware/   # rate-limit, dll
  frontend/       # SPA Vue 3 + Vuetify
    src/
      views/      # halaman (Home, dst)
      stores/     # Pinia store
      components/ # komponen UI
      api/        # axios client
  docker-compose.yml   # opsional, jalankan MySQL + API
```

---

## ⚙️ Cara Menjalankan

### Jalur Local

1. Jalankan MySQL lokal:

   ```sql
   CREATE DATABASE todolist;
   ```
2. Konfigurasi `backend/.env`:

   ```env
   PORT=8080
   DB_DSN=user:password@tcp(127.0.0.1:3306)/todolist?parseTime=true&loc=Local
   CORS_ALLOW_ORIGINS=http://localhost:5173
   ```
3. Jalankan backend:

   ```bash
   cd backend
   go run .
   ```
4. Jalankan frontend:

   ```bash
   cd frontend
   pnpm i
   pnpm dev
   ```

### Jalur Docker Compose

```bash
docker compose up -d --build
```

* API: `http://localhost:8080`
* Web: `http://localhost:5173`

---

## 🧩 Catatan Pengembangan Selanjutnya

* **ENV terpisah**: `.env.development`, `.env.production` ✅ done
* **Testing**: Unit test untuk handlers dengan `httptest` ✅ done, e2e sederhana untuk API
* **CI/CD**: GitHub Actions (build & test backend & frontend)
* **Security**:

  * Validasi input lebih ketat dengan `validator/v10` ✅ done
  * Sanitasi data sebelum masuk DB ✅ done
* **Auth & Multi-user**:

  * Tambah tabel `users`
  * JWT Authentication
  * Kaitkan `todos.user_id`
* **Frontend**:

  * Form validation lebih lengkap
  * Global error boundary
  * Mode dark/light 
* **Monitoring**:

  * Integrasi log dengan sistem observabilitas (mis. ELK, Loki, Grafana)

---