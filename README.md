# 🏛️ ScrapTech — Platform B2B Tender & Pembiayaan Besi Tua Indonesia

> **Sistem Lelang Komoditas Industri, Rekening Bersama Multi-Tahap (Escrow Vault), dan Fasilitas Talangan Modal Kerja (Supply Chain Financing)**  
> Dirancang khusus untuk ekosistem besi tua dan skrap logam industri Indonesia: menghubungkan Pabrik Manufaktur/BUMN, Juragan Lapak, serta Pabrik Peleburan (*Smelter*).

[![Go Version](https://img.shields.io/badge/Golang-1.24_Clean_Arch-00ADD8?style=flat&logo=go)](https://golang.org)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-17_ACID-336791?style=flat&logo=postgresql)](https://www.postgresql.org)
[![Vue 3](https://img.shields.io/badge/Vue.js-3.x_Composition-4FC08D?style=flat&logo=vuedotjs)](https://vuejs.org)
[![Tailwind CSS](https://img.shields.io/badge/Tailwind-3.x_Responsive-38B2AC?style=flat&logo=tailwind-css)](https://tailwindcss.com)
[![Buku Besar](https://img.shields.io/badge/Buku_Besar-Double--Entry_Zero--Sum-brightgreen)](https://en.wikipedia.org/wiki/Double-entry_bookkeeping)
[![Tema](https://img.shields.io/badge/Tampilan-Mode_Terang_%26_Gelap-blueviolet)](#-antarmuka-dan-tema-mode-terang--mode-gelap)

---

## 📌 Latar Belakang & Masalah Industri

Di pasar besi tua industri skala besar (pabrik manufaktur, BUMN minyak & gas, galangan kapal kapal afkir, dan pembongkaran pabrik), satu transaksi tender scrap bernilai antara **Rp 500 Juta hingga di atas Rp 5 Miliar**. Namun, transaksi konvensional menghadapi tiga kendala struktural utama:

1. **Defisit Kepercayaan & Risiko Wanprestasi (Fraud):** Pabrik khawatir pemenang tender tidak melunasi material atau meninggalkan scrap tidak terangkut, sementara pembeli (lapak) khawatir uang muka dibawa kabur oleh calo tender bodong.
2. **Keterbatasan Likuiditas Lapak (Modal Kerja Terkunci):** Juragan lapak sering memenangkan lot puluhan ton, namun kas harian terkuras untuk membayar tarikan timbangan harian per rit truk, sementara pembayaran dari pabrik peleburan (*smelter*) baru cair dalam termin 14–30 hari (*invoice financing gap*).
3. **Sengketa Timbangan & Refraksi Kotoran:** Tiket timbang fisik rentan dimanipulasi (misal: pengisian tangki air sebelum timbang kosong), serta potongan refraksi kotoran/karat yang diputuskan sepihak sering memicu sengketa pembayaran.

**ScrapTech** mendigitalisasi dan mengamankan rantai pasok besi tua ini menjadi ekosistem *Fintech B2B* terpercaya dengan **pencatatan buku besar berpasangan yang tidak dapat diubah (*immutable*), rekening bersama aman, dan deteksi kecurangan timbangan berbasis AI**.

---

## 🏗️ Topologi & Arsitektur Sistem

```
+----------------------------------------------------------------------------------------------------+
|                                    TOPOLOGI SISTEM SCRAPTECH                                       |
+----------------------------------------------------------------------------------------------------+
|  [Pabrik / BUMN Seller]         [Juragan Lapak Pembeli]           [Pabrik Peleburan / Smelter]     |
+----------------------------------------------------------------------------------------------------+
                                                ▲
                                                │
+---------------------------------- FINTECH CORE ENGINE (GOLANG) ------------------------------------+
|  1. Buku Besar Berpasangan (Double-Entry Ledger dengan Invarian Keseimbangan Zero-Sum)             |
|  2. Middleware Idempotensi (Mencegah pemotongan ganda saat sinyal internet lapangan terganggu)     |
|  3. Rekening Bersama Multi-Tahap (Kunci Jaminan Bid-Bond & Pencairan Bertahap Berbasis Progres)    |
|  4. Settlement Timbangan Digital (Formula: (Bruto - Tara - Refraksi) * Harga/Kg - Fee Platform)    |
|  5. Talangan Modal Kerja (Supply Chain Financing berbasis SPK Kontrak Aktif)                       |
|  6. Lapisan AI Agent (Deteksi Anomali Tara Truk & Ekstraksi Dokumen Spesifikasi Tender)            |
|  7. Background Worker Rekonsiliasi Bank & Webhook Payment Gateway VA (HMAC Secure)                 |
+----------------------------------------------------------------------------------------------------+
                                                ▲
                                                │
+------------------------------------ BASIS DATA & KEAMANAN ACID -----------------------------------+
|  - PostgreSQL 17 dengan Isolasi Repeatable Read & Pessimistic Row Locking (`SELECT ... FOR UPDATE`)|
|  - Presisi mata uang bilangan bulat IDR untuk mencegah kesalahan pembulatan desimal                |
+----------------------------------------------------------------------------------------------------+
```

---

## ⚙️ Prinsip Teknis Utama (Fintech Primitives)

### 1. Mesin Buku Besar Berpasangan (*Double-Entry General Ledger*)
* **Tidak Ada Query Mutasi Langsung:** Saldo pengguna tidak pernah diubah menggunakan query naif `UPDATE balance = balance + x`. Semua pergerakan dana dihitung dari riwayat entri jurnal yang kekal (*append-only ledger*).
* **Invarian Keseimbangan Zero-Sum:** Setiap transaksi wajib memenuhi prinsip matematika debit dan kredit seimbang:
  $$\sum \text{Debit} - \sum \text{Credit} = 0$$
* **Bagan Akun Standar (*Chart of Accounts / COA*):**
  - `1010` (Asset): Rekening Bank Kustodi Penampungan
  - `1020` (Asset): Piutang Talangan Modal Kerja (SCF Loan Receivables)
  - `2010` (Liability): Titipan Jaminan Lelang Lapak (Escrow Bid-Bond Payable)
  - `2020` (Liability): Dana Kontrak Tender Terkunci (Escrow Milestone Payable)
  - `2030-xxx` (Liability): Saldo Kas Pengguna yang Dapat Ditarik
  - `4010` (Revenue): Pendapatan Fee Rekening Bersama & Jembatan Timbang (0.5%)
  - `4020` (Revenue): Pendapatan Provisi Talangan Modal SCF (1.0%)

### 2. Perlindungan Idempotensi Transaksi Lapangan
* Seluruh endpoint yang mengubah kondisi keuangan (`POST /tenders/:id/bid`, `POST /tickets/:id/settle`, `POST /financing/apply`) mewajibkan header HTTP `Idempotency-Key`.
* Menjamin tidak ada pemotongan saldo berulang meskipun operator menekan tombol berkali-kali saat koneksi jaringan seluler tidak stabil.

### 3. Penguncian Baris Konkuren (*Pessimistic Row Locking*)
* Menggunakan fitur PostgreSQL `SELECT ... FOR UPDATE` (`clause.Locking{Strength: "UPDATE"}`) untuk mencegah kondisi balapan (*race condition*) saat penetapan pemenang lelang dan pencairan tiket timbangan.

### 4. Integrasi Jembatan Timbang Digital & Pencairan Instan
* Perhitungan otomatis per ritase truk scrap:
  $$\text{Berat Bersih} = \text{Bruto} - \text{Tara}$$
  $$\text{Berat Ditagih} = \text{Berat Bersih} - \left(\text{Berat Bersih} \times \frac{\text{Refraksi } \%}{100}\right)$$
  $$\text{Pencairan Bersih ke Penjual} = (\text{Berat Ditagih} \times \text{Harga/Kg}) - \text{Fee Platform}$$
* Sistem langsung menerbitkan entri jurnal berpasangan multi-split ke rekening penampungan dan penjual dalam hitungan milidetik.

### 5. Lapisan AI Agent (Deteksi Kecurangan & Dokumen)
* **AI Deteksi Anomali Tara Truk:** Membandingkan berat tara truk aktual terhadap riwayat bobot armada kendaraan. Jika terdeteksi deviasi $>3\%$, sistem memberi status `PERINGATAN ANOMALI` untuk mencegah manipulasi berat (seperti tangki air pemberat yang dikosongkan sebelum timbang isi).
* **AI Parser Spesifikasi Tender:** Mengekstrak dokumen tender PDF resmi secara otomatis menjadi parameter siap tawar (jenis logam HMS-1/Plate/Cast Iron, tonase lot, dan batas harga dasar).

---

## 🎨 Antarmuka dan Tema (Mode Terang & Mode Gelap)

ScrapTech dirancang dengan prinsip *UX-friendly* dan *mobile-first*:
- **Mode Gelap (*Dark Mode*):** Palet estetika terminal modern bernuansa *navy/slate*, nyaman untuk analisis tender dan pemantauan jangka panjang.
- **Mode Terang (*Light Mode / Field Mode*):** Kontras tinggi berstandar keterbacaan sinar matahari (*high daylight readability*) menggunakan latar belakang putih bersih (`#ffffff` & `#f8fafc`) dengan tipografi kontras tegas (`#0f172a`), dirancang agar operator jembatan timbang di lapangan pabrik atau dermaga dapat membaca tiket dengan nyaman di bawah terik matahari.
- **Peralihan 1-Klik:** Tombol ganti tema tersedia di pojok kanan atas halaman login maupun di bar navigasi utama. Preferensi tema tersimpan permanen di peramban pengguna.

---

## 🛠️ Struktur Direktori Proyek

```
.
├── backend-go/                      # Inti Mesin Keuangan (Golang Clean Architecture)
│   ├── cmd/api/main.go              # Entry Point Gin HTTP & Dependency Injection
│   ├── internal/
│   │   ├── domain/                  # Entitas Bisnis, Bagan Akun (COA), & Invarian Zero-Sum
│   │   │   ├── ledger.go            # Logika Inti Buku Besar Berpasangan
│   │   │   ├── ledger_test.go       # 🧪 Unit Test Validasi Invarian Akuntansi (Lolos 5/5)
│   │   │   ├── tender.go            # Model Lot Scrap, Tawaran, & Harga Dasar
│   │   │   ├── contract.go          # Surat Perjanjian Kerja (SPK) & Tiket Timbang
│   │   │   └── financing.go         # Kontrak Talangan Modal SCF
│   │   ├── repository/              # Repositori PostgreSQL dengan SELECT FOR UPDATE
│   │   ├── usecase/                 # Logika Bisnis Atomik & Posting Jurnal
│   │   │   ├── agent_usecase.go     # AI Agent Deteksi Anomali Tara & Parser Dokumen
│   │   │   └── agent_usecase_test.go# 🧪 Unit Test AI Agent (Lolos 2/2)
│   │   ├── worker/                  # Background Worker Rekonsiliasi Bank Otomatis
│   │   ├── handler/                 # HTTP Handler & Format Respon JSON Standar
│   │   └── middleware/              # JWT Auth, Header Idempotensi, CORS
│   └── pkg/                         # Koneksi Database, Enkripsi HMAC, & Utilitas
│
├── frontend/                        # Portal Web B2B & PWA Lapangan (Vue 3 + Vite)
│   ├── src/
│   │   ├── views/
│   │   │   ├── HomePage.vue         # Ticker Indeks Harga, Papan Lelang, Simulator Timbangan & SCF
│   │   │   ├── LoginPage.vue        # 1-Klik Profil Demo & Autentikasi Pengguna
│   │   │   └── RegisterPage.vue     # Pendaftaran Perusahaan & Legalitas Usaha
│   │   ├── components/
│   │   │   ├── Navbar.vue           # Navigasi Lega, Switcher Peran, & Toggle Mode Terang/Gelap
│   │   │   └── Footer.vue           # Hak Cipta & Informasi Regulasi ISO/Escrow
│   │   ├── stores/
│   │   │   └── mainStore.js         # State Management Pinia (Tema, Sesi, RBAC)
│   │   └── assets/
│   │       └── main.css             # Desain Tailwind & Token Tema UX-Friendly
│   └── vite.config.js               # Konfigurasi Bundler Vite
│
└── README.md                        # Dokumentasi Lengkap Proyek
```

---

## 🧪 Verifikasi & Pengujian Kode (*Testing*)

Jalankan seluruh test suite unit di backend Golang untuk memverifikasi keabsahan invarian akuntansi:

```bash
cd backend-go
go test -v ./...
```

**Hasil Pengujian:**
```
=== RUN   TestValidateZeroSum_Success
--- PASS: TestValidateZeroSum_Success (0.00s)
=== RUN   TestValidateZeroSum_MultiSplitSuccess
--- PASS: TestValidateZeroSum_MultiSplitSuccess (0.00s)
=== RUN   TestValidateZeroSum_Unbalanced
--- PASS: TestValidateZeroSum_Unbalanced (0.00s)
=== RUN   TestValidateZeroSum_NonPositiveAmount
--- PASS: TestValidateZeroSum_NonPositiveAmount (0.00s)
=== RUN   TestValidateZeroSum_InsufficientEntries
--- PASS: TestValidateZeroSum_InsufficientEntries (0.00s)
=== RUN   TestParseTenderDocument_HMS1
--- PASS: TestParseTenderDocument_HMS1 (0.00s)
=== RUN   TestAnalyzeWeighbridgeSlip_TareFraudDetection
--- PASS: TestAnalyzeWeighbridgeSlip_TareFraudDetection (0.00s)
PASS: 7/7 tes unit lolos verifikasi
```

---

## 🚀 Panduan Menjalankan Program Secara Lokal

### Prasyarat
- Go 1.22 atau lebih baru
- Node.js 18+ dan npm
- PostgreSQL 15+ (Database: `scrapflow_db`)

### 1. Menjalankan Backend (Golang)
```bash
cd backend-go
$env:DATABASE_URL="postgres://postgres:postgres@localhost:5432/scrapflow_db?sslmode=disable"
go run cmd/api/main.go
# Server aktif di http://localhost:8080
```

### 2. Menjalankan Frontend (Vue 3 / Vite)
```bash
cd frontend
npm install
npm run dev -- --host
# Portal frontend aktif di http://localhost:5173
```

---

## 👥 Akun Demo Pengujian Cepat (1-Klik Login)

Aplikasi telah dilengkapi profil demo bawaan yang dapat langsung digunakan tanpa registrasi:

| Peran (*Role*) | Email | Kata Sandi | Wewenang & Kapabilitas Fitur |
| :--- | :--- | :--- | :--- |
| **Pabrik / BUMN** (`PABRIK_MANAGER`) | `budi@krakatausteel.co.id` | `password123` | Terbitkan tender baru, tunjuk pemenang SPK, pantau pencairan dana |
| **Juragan Lapak** (`LAPAK_OWNER`) | `haji.slamet@besiabadi.com` | `password123` | Tawar lelang scrap, kunci jaminan bid-bond, ajukan talangan SCF modal kerja |
| **Operator Timbangan** (`WEIGH_OPERATOR`)| `operator@timbangan.com` | `password123` | Input bruto/tara jembatan timbang, deteksi curang AI, eksekusi pencairan rit |

---

## 📡 Daftar Endpoint API Utama

| Metode | Endpoint | Deskripsi Transaksi | Idempotensi |
| :--- | :--- | :--- | :---: |
| `POST` | `/api/v1/auth/register` | Pendaftaran akun perusahaan & profil KYC | Tidak |
| `POST` | `/api/v1/auth/login` | Autentikasi dan penerbitan token JWT sesi | Tidak |
| `GET` | `/api/v1/tenders` | Mengambil daftar tender scrap aktif | Tidak |
| `POST` | `/api/v1/tenders` | Menerbitkan tender scrap baru *(Pabrik)* | Tidak |
| `POST` | `/api/v1/tenders/:id/bid` | Memasukkan penawaran & mengunci jaminan bid-bond | **Wajib** |
| `POST` | `/api/v1/tenders/:id/award` | Menetapkan pemenang lelang & menerbitkan SPK resmi | **Wajib** |
| `POST` | `/api/v1/contracts/:id/tickets` | Input tiket timbangan jembatan timbang | Tidak |
| `POST` | `/api/v1/tickets/:id/settle` | Eksekusi pencairan multi-split buku besar per rit truk | **Wajib** |
| `GET` | `/api/v1/ledger/balance` | Pengecekan saldo kas dan jaminan dari buku besar | Tidak |
| `GET` | `/api/v1/ledger/transactions` | Mengambil jejak audit entri jurnal debit/kredit | Tidak |
| `POST` | `/api/v1/financing/apply` | Pengajuan fasilitas talangan modal SCF berbasis SPK | Tidak |
| `POST` | `/api/v1/financing/:id/disburse`| Pencairan dana talangan langsung ke escrow kontrak | **Wajib** |
| `POST` | `/api/v1/agent/parse-spec` | AI Agent ekstraksi parameter dokumen tender PDF | Tidak |
| `POST` | `/api/v1/agent/detect-fraud` | AI Agent deteksi deviasi bobot tara armada truk | Tidak |
| `POST` | `/api/v1/payment/webhook` | Webhook pembayaran Virtual Account dengan verifikasi HMAC | **Wajib** |

---

## 👨‍💻 Pengembang & Kontak
- **Pengembang:** Frans Alwan
- **Fokus Keahlian:** Rekayasa Sistem Finansial, Backend Berkinerja Tinggi & Full-Stack (Golang, PostgreSQL, Vue 3, Clean Architecture)
