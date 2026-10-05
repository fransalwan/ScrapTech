# 🏛️ ScrapFlow — B2B Scrap Metal Tender & Commodity Fintech Platform

> **Enterprise-grade B2B Commodity Auction, Milestone Escrow Vault, & Supply Chain Financing (SCF) Core**  
> Tailored for the heavy scrap metal (*besi tua*) industrial ecosystem in Indonesia.

[![Go Version](https://img.shields.io/badge/Go-1.22+-00ADD8?style=flat&logo=go)](https://golang.org)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-ACID-336791?style=flat&logo=postgresql)](https://www.postgresql.org)
[![Vue 3](https://img.shields.io/badge/Vue.js-3.x-4FC08D?style=flat&logo=vuedotjs)](https://vuejs.org)
[![Tailwind CSS](https://img.shields.io/badge/Tailwind-3.x-38B2AC?style=flat&logo=tailwind-css)](https://tailwindcss.com)
[![Accounting](https://img.shields.io/badge/Ledger-Double--Entry%20Zero--Sum-brightgreen)](https://en.wikipedia.org/wiki/Double-entry_bookkeeping)

---

## 📌 Executive Summary & Business Rationale

In the heavy industrial scrap market (pabrik manufaktur, BUMN, galangan kapal, dan juragan lapak besi tua), single tender transactions often reach **Rp 500 Million to Rp 5 Billion**. However, the industry suffers from three structural bottlenecks:

1. **High Trust Deficit & Tender Fraud:** Fear of fake scrap auctions, default on advance payments, or material abandonment by unreliable contractors.
2. **Severe Liquidity Choke (Working Capital Strain):** Scrap aggregators (lapak) win multi-ton tenders but lack instant cash flow to settle daily weighbridge pickups while waiting for smelter (*pabrik peleburan*) invoice payment cycles (14–30 day payment terms).
3. **Weight & Refraction Sengketa (Disputes):** Manual weighbridge tickets (*jembatan timbang*) and arbitrary dirt/rust deduction (*refraksi*) frequently trigger financial disputes between factories and buyers.

**ScrapFlow** transforms this informal commodity flow into a regulated fintech ecosystem with **cryptographic auditability, immutable accounting, and automated liquidity**.

---

## 🏗️ Core Fintech Primitives

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
|  2. Idempotency Key Engine (Anti-double settlement & resilient webhook processing)                 |
|  3. Milestone Escrow Vault (Bid-Bond Lock, Auto-Refund, and Multi-Stage Release)                   |
|  4. Digital Weighbridge Settlement (Formula: (Gross - Tare - Refraksi) * PricePerKg)               |
|  5. Supply Chain Financing (PO/SPK Loan injection directly into tender escrow)                    |
+----------------------------------------------------------------------------------------------------+
                                                ▲
                                                │
+-------------------------------------- PERSISTENCE & ACID LAYER ------------------------------------+
|  - PostgreSQL with Repeatable Read Isolation & Pessimistic Row Locking (`SELECT ... FOR UPDATE`)   |
|  - BigInt currency precision (Sen / Integer IDR) to prevent floating-point rounding bugs           |
+----------------------------------------------------------------------------------------------------+
```

### 1. Double-Entry General Ledger Engine
* **No Direct Mutating Queries:** Balances are never modified with naive `UPDATE balance = balance + x`. All monetary movements are calculated from immutable **Journal Postings**.
* **Zero-Sum Balance Invariant:** Every journal entry enforces $\sum \text{Debit} - \sum \text{Credit} = 0$.
* **Standard Chart of Accounts (COA):**
  - `1010` (Asset): Platform Bank Custody
  - `1020` (Asset): Supply Chain Financing Loan Receivables
  - `2010` (Liability): Escrow Bid-Bond Payable (Lapak Deposits)
  - `2020` (Liability): Escrow Milestone Payable (Tender Contracts)
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

---

## 🛠️ Repository Structure

```
.
├── backend-go/                      # Core Financial Engine (Golang Clean Architecture)
│   ├── cmd/api/main.go              # Gin HTTP Engine & DI Wiring
│   ├── internal/
│   │   ├── domain/                  # Enterprise Models, COA, & Zero-Sum Invariants
│   │   │   ├── ledger.go            # Double-Entry Core
│   │   │   ├── ledger_test.go       # 🧪 Unit Tests for Financial Invariants
│   │   │   ├── tender.go            # Scrap Lots, Bids, Reserve Prices
│   │   │   ├── contract.go          # SPK Agreements & Weighbridge Tickets
│   │   │   └── financing.go         # PO/SPK Supply Chain Financing
│   │   ├── repository/              # PostgreSQL Repositories with SELECT FOR UPDATE
│   │   ├── usecase/                 # Atomic Business Logic & Double-Entry Postings
│   │   ├── handler/                 # HTTP Handlers & JSON Envelopes
│   │   └── middleware/              # Auth JWT, Idempotency-Key, & CORS
│   ├── pkg/
│   │   ├── database/postgres.go     # PostgreSQL AutoMigrate & COA Seeding
│   │   └── response/response.go     # Standardized JSON Envelopes
│   └── README.md                    # Dedicated Backend Documentation
│
├── frontend/                        # Frontend B2B Portal & Field PWA (Vue 3 + Vite)
│   ├── src/views/HomePage.vue       # Live Tender Board, Weighbridge Simulator, & SCF Calc
│   ├── src/components/Navbar.vue    # High-Density B2B Navigation
│   └── src/stores/mainStore.js      # Pinia Store State Management
│
└── README.md                        # Project Overview & Portfolio Specification
```

---

## 🧪 Testing & Verification

Run the Golang unit test suite covering double-entry accounting invariants:

```bash
cd backend-go
go test -v ./...
```

**Test Output:**
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
PASS: 5/5 tests passed
```

---

## 🚀 Running the Services Locally

### 1. Backend (Golang)
```bash
cd backend-go
go run cmd/api/main.go
# Server runs on http://localhost:8080
```

### 2. Frontend (Vue 3 / Vite)
```bash
cd frontend
npm install
npm run dev
# Dashboard opens on http://localhost:5173
```

---

## 📡 API Reference Summary

| Method | Endpoint | Description | Idempotency |
| :--- | :--- | :--- | :---: |
| `POST` | `/api/v1/auth/register` | Register Company & KYC profile | No |
| `POST` | `/api/v1/auth/login` | Authenticate & issue JWT | No |
| `GET` | `/api/v1/tenders` | List active scrap tenders | No |
| `POST` | `/api/v1/tenders` | Publish scrap tender *(Pabrik)* | No |
| `POST` | `/api/v1/tenders/:id/bid` | Place bid + lock bid-bond into Escrow | **Mandatory** |
| `POST` | `/api/v1/tenders/:id/award` | Award tender & generate official SPK contract | **Mandatory** |
| `POST` | `/api/v1/contracts/:id/tickets` | Submit digital weighbridge slip | No |
| `POST` | `/api/v1/tickets/:id/settle` | Execute atomic double-entry settlement | **Mandatory** |
| `GET` | `/api/v1/ledger/balance` | Query verified wallet balance from ledger | No |
| `GET` | `/api/v1/ledger/transactions` | Full double-entry audit trail | No |
| `POST` | `/api/v1/financing/apply` | Apply for PO/SPK Supply Chain Financing | No |
| `POST` | `/api/v1/financing/:id/disburse`| Disburse loan directly into tender escrow | **Mandatory** |

---

## 👨‍💻 Author & Contact
- **Developer:** Frans Alwan
- **Role Target:** Fintech Backend / Full-Stack Engineer (Go / PostgreSQL / Vue)
