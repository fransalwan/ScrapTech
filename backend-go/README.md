# 🏛️ ScrapFlow Core Fintech Engine (Golang)

> **High-Throughput Commodity Escrow & Supply Chain Financing (SCF) Core**  
> Built for the Scrap Metal Industrial Market in Indonesia.

[![Go Version](https://img.shields.io/badge/Go-1.22+-00ADD8?style=flat&logo=go)](https://golang.org)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-ACID-336791?style=flat&logo=postgresql)](https://www.postgresql.org)
[![Architecture](https://img.shields.io/badge/Architecture-Clean%20%2F%20DDD-orange)](https://github.com/fransalwan)
[![Double-Entry Ledger](https://img.shields.io/badge/Accounting-Zero--Sum%20Ledger-brightgreen)](https://en.wikipedia.org/wiki/Double-entry_bookkeeping)

---

## 🎯 Background & Business Problem

In the heavy industrial scrap metal market (pabrik peleburan, BUMN, galangan kapal, dan juragan lapak besi tua), transactions frequently reach **billions of Rupiah per lot**, yet suffer from:
1. **Severe Trust Deficit:** Fear of fake tenders, default on advance payments, or material theft.
2. **Liquidity Choke:** Scrap aggregators (lapak) win large tenders but lack instant cash flow to fulfill immediate payments while waiting for smelter payout cycles (14–30 day invoice terms).
3. **Weight & Refraction Disputes:** Manual weighbridge tickets (*jembatan timbang*) are vulnerable to tare manipulation and tare inflation fraud.

**ScrapFlow** addresses these issues with a dedicated financial core:
- **Milestone-based Escrow Vault** (Bid-Bond and tender advance fund locking).
- **Supply Chain Financing (PO/SPK Talangan)** injecting liquidity directly into verified tender escrows.
- **Atomic Weighbridge Settlement** (Gross - Tare - Refraction formula calculated and disbursed per truckload).

---

## 🏗️ Architectural Highlights (Fintech Primitives)

### 1. Double-Entry General Ledger Engine
- Every money movement is recorded as balanced **Debit** and **Credit** postings.
- **Zero-Sum Invariant:** $\sum \text{Debit} - \sum \text{Credit} = 0$.
- No direct `UPDATE balance` statements; account balances are dynamically computed from immutable journal entries.
- Currency is strictly handled in integer units (`int64` / Sen Rupiah) to prevent IEEE 754 floating-point errors.

### 2. Idempotency Key Middleware
- Mutating financial operations (`POST /bids`, `POST /tickets/:id/settle`, `POST /ledger/deposit`) mandate the `Idempotency-Key` HTTP header.
- Protects against duplicate deductions and race conditions caused by intermittent field connectivity.

### 3. Concurrency & Pessimistic Locking
- Database transactions use **PostgreSQL `SELECT ... FOR UPDATE`** row locking (`clause.Locking{Strength: "UPDATE"}`) during tender awarding and weighbridge settlements to prevent double-settlement.

---

## 📐 Chart of Accounts (COA)

| Code | Account Name | Type | Normal Balance |
| :--- | :--- | :--- | :--- |
| **1010** | Cash in Bank - Platform Custody | `ASSET` | Debit |
| **1020** | SCF Financing Loan Receivables | `ASSET` | Debit |
| **2010** | Escrow Bid-Bond Payable (Lapak Deposits) | `LIABILITY` | Credit |
| **2020** | Escrow Milestone Payable (Pabrik Contracts) | `LIABILITY` | Credit |
| **2030-xxx** | Withdrawable Wallet Balance (Company) | `LIABILITY` | Credit |
| **4010** | Escrow & Settlement Fee Revenue | `REVENUE` | Credit |
| **4020** | SCF Financing Fee & Interest Revenue | `REVENUE` | Credit |

---

## 🚀 API Endpoints

### Authentication & Company KYC
- `POST /api/v1/auth/register` - Register Company (Pabrik, Lapak, Smelter) & User
- `POST /api/v1/auth/login` - Authenticate & obtain JWT

### Tender & Bidding Engine
- `GET /api/v1/tenders` - List published tenders (Filter by status)
- `GET /api/v1/tenders/:id` - Get tender details & active bids
- `POST /api/v1/tenders` - Publish scrap tender *(Auth required)*
- `POST /api/v1/tenders/:id/bid` - Submit bid + lock bid-bond *(Idempotency-Key required)*
- `POST /api/v1/tenders/:id/award` - Award tender & generate SPK *(Idempotency-Key required)*

### Weighbridge & Field Settlement
- `POST /api/v1/contracts/:contract_id/tickets` - Submit truck weighbridge slip
- `GET /api/v1/contracts/:contract_id/tickets` - View truck receipts for a contract
- `POST /api/v1/tickets/:ticket_id/settle` - Atomic Double-Entry Settlement *(Idempotency-Key required)*

### Ledger & Wallet
- `GET /api/v1/ledger/balance` - Get verified company wallet balance
- `GET /api/v1/ledger/transactions` - Audit trail of all double-entry journals
- `POST /api/v1/ledger/deposit` - Deposit into escrow vault *(Idempotency-Key required)*

### Supply Chain Financing (SCF)
- `POST /api/v1/financing/apply` - Apply for PO/SPK tender financing
- `GET /api/v1/financing/my-facilities` - View active financing facilities
- `POST /api/v1/financing/:id/disburse` - Disburse loan directly into tender escrow *(Idempotency-Key required)*

---

## 🧪 Testing

```bash
cd backend-go
go test -v ./...
```
