# MazadPay API

Backend API for [MazadPay](https://mazadpay.com) — أول منصة مزادات رقمية في موريتانيا.

**الحالة الحالية:** Milestone 1 — Skeleton + Health Check فقط.
لا يوجد Database أو Auth في هذه المرحلة.

---

## Stack

| المكوّن | التقنية |
|--------|---------|
| Language | Go 1.22+ |
| Router | Fiber v2 |
| Middleware | recover, requestid, logger, cors, helmet |
| Database | — (Milestone 2) |
| Auth | — (Milestone 3) |

---

## التشغيل محلياً

```bash
# 1. انسخ متغيرات البيئة
cp .env.example .env

# 2. شغّل السيرفر
go run ./cmd/server

# أو مع تحديد port
PORT=9090 go run ./cmd/server
```

السيرفر يعمل على: `http://localhost:8080`

---

## Endpoints الحالية

| Method | Path | الوصف |
|--------|------|-------|
| GET | `/health` | Health check (Render / uptime monitors) |
| GET | `/api/v1/health` | Health check versioned |

### مثال على الرد

```json
{
  "status": "ok",
  "service": "mazadpay-api",
  "env": "development",
  "version": "0.1.0"
}
```

### 404 Response

```json
{
  "error": "not_found",
  "message": "Route not found"
}
```

---

## متغيرات البيئة

| المتغير | القيمة الافتراضية | الوصف |
|---------|------------------|-------|
| `PORT` | `8080` | منفذ السيرفر |
| `APP_ENV` | `development` | البيئة: development / production |
| `APP_VERSION` | `0.1.0` | إصدار التطبيق |
| `CORS_ALLOWED_ORIGINS` | localhost + mazadpay.com | الدومينات المسموح بها |

للمتغيرات المستقبلية (Database، JWT، R2) راجع [.env.example](.env.example).

---

## هيكل الملفات

```
mazadpay-api/
├── cmd/
│   └── server/
│       └── main.go             ← Entry point
├── internal/
│   ├── config/
│   │   └── config.go           ← Environment config loader
│   └── http/
│       ├── router.go           ← Fiber app + middleware + routes
│       └── handlers/
│           └── health.go       ← Health check handler
├── .env.example
├── .gitignore
├── go.mod
├── go.sum
└── README.md
```

---

## ملاحظات مهمة

- **لا تضع `.env` في Git** — يحتوي على secrets.
- **CORS** مضبوط على دومينات محددة — لا wildcard `*`.
- **Milestone 1** فقط: لا Database، لا Auth، لا Auctions.
- للخطة الكاملة راجع [PHASE2_BACKEND_PLAN.md](../plateforme-mazadpay/PHASE2_BACKEND_PLAN.md).

---

© 2026 MazadPay — Nouakchott, Mauritania
