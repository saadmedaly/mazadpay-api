# MazadPay API

Backend API for [MazadPay](https://mazadpay.com) — أول منصة مزادات رقمية في موريتانيا.

| الحالة | التفاصيل |
|--------|---------|
| Milestone 1 | ✅ Backend skeleton + health endpoints |
| Milestone 2 | ✅ Database config + migrations schema |
| Milestone 3 | ⏳ Authentication (upcoming) |

---

## Stack

| المكوّن | التقنية |
|--------|---------|
| Language | Go 1.22+ |
| Router | Fiber v2 |
| Database driver | pgx/v5 (pgxpool) |
| Database | PostgreSQL 16 / Neon |
| Migrations | golang-migrate |
| Hosting | Render |

---

## التشغيل محلياً

```bash
# 1. انسخ متغيرات البيئة
cp .env.example .env
# عدّل DATABASE_URL إذا أردت اتصال DB، أو اتركه فارغاً للتشغيل بدون DB

# 2. شغّل السيرفر
go run ./cmd/server

# إذا port 8080 محجوز
PORT=3001 go run ./cmd/server
```

---

## Endpoints

| Method | Path | الوصف |
|--------|------|-------|
| GET | `/health` | Health check (Render monitor) |
| GET | `/api/v1/health` | Health check versioned |
| GET | `/api/v1/health/db` | Database connection status |

### مثال `/health`

```json
{
  "status": "ok",
  "service": "mazadpay-api",
  "env": "development",
  "version": "0.1.0"
}
```

### مثال `/api/v1/health/db` — بدون DATABASE_URL

```json
{ "status": "disabled", "database": "not_configured" }
```

### مثال `/api/v1/health/db` — مع DB متصل

```json
{ "status": "ok", "database": "connected" }
```

### مثال `/api/v1/health/db` — DB غير متاح

```json
{ "status": "unreachable", "database": "unreachable" }
```
HTTP status: `503 Service Unavailable`

---

## Database & Migrations

### متطلبات

- PostgreSQL 16+ أو **Neon** (يشترط `sslmode=require`)
- لا تستخدم SQLite في الإنتاج

### تشغيل migrations

```bash
# تثبيت golang-migrate (مرة واحدة)
go install github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# تطبيق كل migrations
migrate -path migrations -database "$DATABASE_URL" up

# التراجع خطوة واحدة
migrate -path migrations -database "$DATABASE_URL" down 1
```

> **تحذير:** لا تشغّل migrations على production مباشرة — اختبر على Neon branch أولاً.

### الجداول (Migration 000001)

| الجدول | الوصف |
|--------|-------|
| `categories` | تصنيفات المزادات |
| `users` | المستخدمون (password_hash جاهز لـ Milestone 3) |
| `auctions` | المزادات |
| `auction_images` | صور المزادات (R2 URLs) |
| `bids` | المزايدات |
| `favorites` | المفضّلات |
| `contact_messages` | رسائل التواصل |
| `audit_logs` | سجل العمليات الإدارية |

---

## متغيرات البيئة

| المتغير | الافتراضي | الوصف |
|---------|-----------|-------|
| `PORT` | `8080` | منفذ السيرفر |
| `APP_ENV` | `development` | البيئة |
| `APP_VERSION` | `0.1.0` | الإصدار |
| `CORS_ALLOWED_ORIGINS` | localhost + mazadpay.com | الدومينات المسموحة |
| `DATABASE_URL` | *(فارغ)* | Neon PostgreSQL URL |
| `DB_MAX_CONNS` | `10` | أقصى عدد connections |
| `DB_MIN_CONNS` | `1` | أدنى عدد connections |
| `DB_MAX_CONN_LIFETIME` | `1h` | عمر الـ connection |

للقائمة الكاملة (JWT، R2…) راجع [.env.example](.env.example).

---

## هيكل الملفات

```
mazadpay-api/
├── cmd/server/main.go
├── internal/
│   ├── config/config.go          ← env loader
│   ├── db/
│   │   ├── db.go                 ← pgxpool connect
│   │   └── health.go             ← DB ping
│   └── http/
│       ├── router.go             ← Fiber app + middleware + routes
│       └── handlers/
│           └── health.go         ← /health + /health/db handlers
├── migrations/
│   ├── 000001_create_core_tables.up.sql
│   └── 000001_create_core_tables.down.sql
├── .env.example
├── .gitignore
├── go.mod / go.sum
└── README.md
```

---

## ⚠️ ملاحظات مهمة

- **لا تضع `.env` في Git** — يحتوي على DATABASE_URL وSecrets.
- **CORS** مضبوط على دومينات محددة — لا wildcard `*` أبداً.
- **السيرفر يعمل بدون DB** — `DATABASE_URL` فارغ = وضع development بدون قاعدة بيانات.
- **Milestone 2** فقط: لا Auth، لا Login/Register، لا Auctions APIs.

---

© 2026 MazadPay — Nouakchott, Mauritania
