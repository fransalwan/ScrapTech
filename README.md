# 🏛️ ScrapFlow — B2B Scrap Metal Tender & Commodity Fintech Platform

> **Enterprise-grade B2B Commodity Auction, Milestone Escrow Vault, & Supply Chain Financing (SCF) Core**  
> Tailored for the heavy scrap metal (*besi tua*) industrial ecosystem in Indonesia (Pabrik Manufaktur/BUMN, Juragan Lapak, dan Pabrik Peleburan / Smelter).

[![Go Version](https://img.shields.io/badge/Go-1.24-00ADD8?style=flat&logo=go)](https://golang.org)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-17_ACID-336791?style=flat&logo=postgresql)](https://www.postgresql.org)
[![Vue 3](https://img.shields.io/badge/Vue.js-3.x_Composition-4FC08D?style=flat&logo=vuedotjs)](https://vuejs.org)
[![Tailwind CSS](https://img.shields.io/badge/Tailwind-3.x-38B2AC?style=flat&logo=tailwind-css)](https://tailwindcss.com)
[![Accounting](https://img.shields.io/badge/Ledger-Double--Entry%20Zero--Sum-brightgreen)](https://en.wikipedia.org/wiki/Double-entry_bookkeeping)
[![Theme](https://img.shields.io/badge/UI-Dark_%26_Light_Mode-blueviolet)](#-antarmuka-dan-tema-dark--light-mode)

---

## 📌 Executive Summary & Business Rationale

In the heavy industrial scrap market (pabrik manufaktur, BUMN, galangan kapal, dan juragan lapak besi tua), single tender transactions often reach **Rp 500 Million to Rp 5 Billion**. However, the industry suffers from three structural bottlenecks:

1. **High Trust Deficit & Tender Fraud:** Fear of fictitious scrap auctions, default on advance payments, or material abandonment by unreliable contractors.
2. **Severe Liquidity Choke (Working Capital Strain):** Scrap aggregators (lapak) win multi-ton tenders but lack instant cash flow to settle daily weighbridge pickups while waiting for smelter (*pabrik peleburan*) invoice payment cycles (14–30 day payment terms).
3. **Weight & Refraction Sengketa (Disputes):** Manual weighbridge tickets (*jembatan timbang*) and arbitrary dirt/rust deduction (*refraksi*) frequently trigger financial disputes between factories and buyers.

**ScrapFlow** transforms this informal commodity flow into a regulated fintech ecosystem with **cryptographic auditability, immutable accounting, automated liquidity, and AI-assisted fraud detection**.

---

## 🏗️ Core Fintech Primitives & Topology

```
+----------------------------------------------------------------------------------------------------+
|                                    SCRAPFLOW FINTECH TOPOLOGY                                      |
+----------------------------------------------------------------------------------------------------+
|  [Pabrik / BUMN Seller]         [Juragan Lapak Buyer]             [Peleburan / Smelter Offtaker]   |
+----------------------------------------------------------------------------------------------------+
                                                ▲
                                                │
+------------------------------------- FINTECH CORE ENGINE (GOLANG) --------------------------------+
|  1. Double-Entry General Ledger (Debit == Credit Zero-Sum Invariant)                               |
|  2. Idempotency Key Engine (Anti-double deduction & resilient field offline sync)                  |
|  3. Milestone Escrow Vault (Bid-Bond Lock, Auto-Refund, and Multi-Stage Release)                   |
|  4. Digital Weighbridge Settlement (Formula: (Gross - Tare - Refraksi) * PricePerKg)               |
|  5. Supply Chain Financing (PO/SPK Loan injection directly into tender escrow)                    |
|  6. AI Agents Layer (Tare Fraud OCR Anomaly Detection & Spec Parser)                               |
|  7. Background Bank Reconciliation Worker & Payment Gateway VA Webhook                             |
+----------------------------------------------------------------------------------------------------+
                                                ▲
                                                │
+-------------------------------------- PERSISTENCE & ACID LAYER ------------------------------------+
|  - PostgreSQL 17 with Repeatable Read Isolation & Pessimistic Row Locking (`SELECT ... FOR UPDATE`)|
|  - BigInt currency precision (IDR Integer) to prevent floating-point rounding bugs                 |
+----------------------------------------------------------------------------------------------------+
```

### 1. Double-Entry General Ledger Engine
* **No Direct Mutating Queries:** Balances are never modified with naive `UPDATE balance = balance + x`. All monetary movements are calculated from immutable **Journal Postings**.
* **Zero-Sum Balance Invariant:** Every journal entry enforces $\sum \text{Debit} - \sum \text{Credit} = 0$.
* **Standard Chart of Accounts (COA):**
  - `1010` (Asset): Platform Bank Custody (Penampungan)
  - `1020` (Asset): Supply Chain Financing Loan Receivables
  - `2010` (Liability): Escrow Bid-Bond Payable (Titipan Jaminan Lapak)
  - `2020` (Liability): Escrow Milestone Payable (Kontrak Tender)
  - `2030-xxx` (Liability): User Withdrawable Balance
  - `4010` (Revenue): Escrow & Settlement Fee Revenue (0.5%)
  - `4020` (Revenue): SCF Loan Origination Fee (1.0%)

### 2. Idempotency Key Protection
* All mutating financial endpoints (`POST /tenders/:id/bid`, `POST /tickets/:id/settle`, `POST /ledger/deposit`) enforce the `Idempotency-Key` HTTP header.
* Prevents duplicate deductions when field operators or truck drivers experience intermittent mobile connectivity.

### 3. Concurrency Control with Pessimistic Row Locking
* Critical operations utilize **PostgreSQL `SELECT ... FOR UPDATE`** row locking (`clause.Locking{Strength: "UPDATE"}`) to prevent race conditions during tender awarding and weighbridge settlements.

### 4. Digital Weighbridge Integration & Instant Settlement
* Per-truckload automated formula:
  $$\text{Net Weight} = \text{Gross} - \text{Tare}$$
  $$\text{Billable Weight} = \text{Net} - \left(\text{Net} \times \frac{\text{Refraction } \%}{100}\right)$$
  $$\text{Net Settlement} = (\text{Billable Weight} \times \text{PricePerKg}) - \text{Platform Fee}$$
* Generates balanced multi-split journal entries automatically upon ticket verification.

### 5. AI Agent Layer (Autonomous Risk & Spec Audit)
* **Slip Tare Anomaly Detector:** Membandingkan berat tara truk aktual terhadap basis data historis armada. Jika deviasi $>3\%$, transaksi ditandai `SUSPECT_FRAUD` (mencegah kecurangan tangki air/pemberat).
* **Tender Spec Document Parser:** Mengurai dokumen PDF tender afkir BUMN secara otomatis ke dalam parameter material (HMS-1, Plate, Cu, tonase, dan reserve price).

---

## 🎨 Antarmuka dan Tema (Dark & Light Mode)

ScrapFlow dirancang khusus untuk kenyamanan operasional industri:
- **Mode Gelap (Dark Mode):** Estetika terminal fintech modern, nyaman untuk monitoring lelang intensif.
- **Mode Terang (Light Mode / Field Mode):** Kontras tinggi berstandar keterbacaan luar ruangan (*high sunlight readability*) bagi petugas jembatan timbang di lapangan pabrik.
- **Peralihan 1-Klik:** Dapat diakses langsung melalui tombol toggle di Navbar atas maupun halaman Login.

---

## 🛠️ Repository Structure

```
.
├── backend-go/                      # Core Financial Engine (Golang Clean Architecture)
│   ├── cmd/api/main.go              # Gin Engine & Dependency Injection Wiring
│   ├── internal/
│   │   ├── domain/                  # Entity Models, Chart of Accounts, & Zero-Sum Math
│   │   │   ├── ledger.go            # Double-Entry Ledger Core
│   │   │   ├── ledger_test.go       # 🧪 Unit Tests for Financial Invariants (5/5 Passing)
│   │   │   ├── tender.go            # Scrap Lots, Bids, Reserve Prices
│   │   │   ├── contract.go          # SPK Agreements & Weighbridge Tickets
│   │   │   └── financing.go         # Supply Chain Financing (SCF) Contracts
│   │   ├── repository/              # PostgreSQL Repositories with SELECT FOR UPDATE
│   │   ├── usecase/                 # Atomic Business Logic & Ledger Postings
│   │   │   ├── agent_usecase.go     # AI Slip Tare Fraud & Tender Spec Analyzer
│   │   │   └── agent_usecase_test.go# 🧪 Unit Tests for AI Agents (2/2 Passing)
│   │   ├── worker/                  # Background Bank Reconciliation Worker
│   │   ├── handler/                 # HTTP Handlers & JSON Envelopes
│   │   └── middleware/              # JWT Auth, Idempotency-Key, CORS
│   └── pkg/                         # Database connection, Utilities & Responses
│
├── frontend/                        # Modern B2B Web Portal & Mobile First UI (Vue 3)
│   ├── src/
│   │   ├── views/
│   │   │   ├── HomePage.vue         # Live Ticker, Tender Board, Weighbridge, SCF & Audit
│   │   │   ├── LoginPage.vue        # 1-Click Demo Profiles & Auth Flow
│   │   │   └── RegisterPage.vue     # Profil Usaha & Verifikasi KYC
│   │   ├── components/
│   │   │   └── Navbar.vue           # Role Switcher, Light/Dark Toggle, Top Navigation
│   │   ├── stores/
│   │   │   └── mainStore.js         # Pinia Store (Theme, RBAC, Auth, Sesi)
│   │   └── assets/
│   │       └── main.css             # Tailwind & Semantic Theme Variables
│   └── vite.config.js               # Vite Bundler Setup
│
└── README.md                        # Platform Documentation
```

---

## 🧪 Testing & Verification

Jalankan seluruh test suite unit di backend Golang:

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
PASS: 7/7 tests passed
```

---

## 🚀 Menjalankan Aplikasi Secara Lokal

### Prasyarat
- Go 1.22+
- Node.js 18+ & npm
- PostgreSQL 15+ (Database `scrapflow_db`)

### 1. Jalankan Backend (Golang)
```bash
cd backend-go
$env:DATABASE_URL="postgres://postgres:postgres@localhost:5432/scrapflow_db?sslmode=disable"
go run cmd/api/main.go
# Server berjalan di http://localhost:8080
```

### 2. Jalankan Frontend (Vue 3 / Vite)
```bash
cd frontend
npm install
npm run dev -- --host
# Frontend berjalan di http://localhost:5173
```

---

## 👥 Akun Demo Pengujian Cepat (1-Click Login)

Anda dapat langsung mencoba semua alur peran pengguna tanpa registrasi manual:

| Peran (Role) | Email | Password | Kapabilitas Utama |
| :--- | :--- | :--- | :--- |
| **Pabrik / BUMN** (`PABRIK_MANAGER`) | `budi@krakatausteel.co.id` | `password123` | Terbitkan tender baru, tentukan pemenang SPK, terima pencairan |
| **Juragan Lapak** (`LAPAK_OWNER`) | `haji.slamet@besiabadi.com` | `password123` | Tawar lelang, kunci bid-bond, ajukan talangan SCF modal kerja |
| **Operator Timbangan** (`WEIGH_OPERATOR`)| `operator@timbangan.com` | `password123` | Input slip bruto/tara, deteksi anomali AI, pencairan per rit |

---

## 📡 Ringkasan API Endpoint

| Method | Endpoint | Deskripsi | Idempotency |
| :--- | :--- | :--- | :---: |
| `POST` | `/api/v1/auth/register` | Pendaftaran perusahaan & legalitas KYC | No |
| `POST` | `/api/v1/auth/login` | Autentikasi pengguna & terbitkan JWT token | No |
| `GET` | `/api/v1/tenders` | Ambil daftar lelang scrap aktif | No |
| `POST` | `/api/v1/tenders` | Publikasikan tender scrap baru *(Pabrik)* | No |
| `POST` | `/api/v1/tenders/:id/bid` | Tawar lelang + kunci jaminan ke Escrow | **Mandatory** |
| `POST` | `/api/v1/tenders/:id/award` | Tunjuk pemenang & terbitkan SPK resmi | **Mandatory** |
| `POST` | `/api/v1/contracts/:id/tickets` | Input tiket timbang jembatan timbang | No |
| `POST` | `/api/v1/tickets/:id/settle` | Eksekusi pencairan multi-split buku besar | **Mandatory** |
| `GET` | `/api/v1/ledger/balance` | Cek saldo kas & jaminan dari buku besar | No |
| `GET` | `/api/v1/ledger/transactions` | Jejak audit lengkap entri jurnal debit/kredit | No |
| `POST` | `/api/v1/financing/apply` | Pengajuan talangan modal SCF berbasis SPK | No |
| `POST` | `/api/v1/financing/:id/disburse`| Pencairan talangan langsung ke escrow tender | **Mandatory** |
| `POST` | `/api/v1/agent/parse-spec` | AI Agent ekstraksi dokumen tender | No |
| `POST` | `/api/v1/agent/detect-fraud` | AI Agent deteksi deviasi bobot tara armada | No |
| `POST` | `/api/v1/payment/webhook` | Webhook pembayaran VA dengan verifikasi HMAC | **Mandatory** |

---

## 👨‍💻 Pengembang
- **Developer:** Frans Alwan
- **Fokus Keahlian:** Fintech Backend & Full-Stack Systems (Go, PostgreSQL, Vue 3, Clean Architecture)
