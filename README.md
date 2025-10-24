# WhatsApp CRM (Go + Fiber + GORM + MariaDB)

Aplikasi REST API untuk CRM dengan integrasi WhatsApp (send chat + webhook). Framework: Go Fiber, ORM: GORM, DB: MariaDB, Cache/Queue: Redis.

## Fitur Utama
- Auth JWT (login/register/profile/change-password)
- Role: admin, agent, supervisor
- Manajemen User (admin-only)
- Manajemen Customer
- Percakapan (assign agent, status, priority, notes)
- Pesan (list per percakapan, kirim text; media & template menyusul)
- Webhook endpoint (verify + handler) — tersedia, eksekusi dapat diaktifkan setelah update ini

## Struktur Endpoint (prefix versi: /api/v1)
- Auth: POST /auth/login, POST /auth/register, GET /auth/profile, POST /auth/change-password
- Users (admin): GET/POST/PUT/DELETE /users
- Customers: GET/POST/PUT/DELETE /customers, GET /customers/:id
- Conversations: GET/POST/GET/:id, PUT /:id/{assign|status|priority|notes}
- Messages: GET /messages/conversation/:id, POST /messages/conversation/:id/{text|media|template}
- Webhook: GET/POST /webhook/whatsapp

## Menjalankan Secara Lokal
1. Salin .env.example menjadi .env dan sesuaikan nilai
2. Jalankan database & redis via docker-compose:
   docker-compose up -d db redis
3. Jalankan API:
   go mod tidy
   go run cmd/api/main.go
4. Health check:
   curl http://localhost:8080/health

## Docker All-in-one
- Build dan jalankan semua (db, redis, api):
  docker-compose up --build -d

## Seeding Admin Default
- Jalankan: go run cmd/seed/main.go
- Kredensial default: admin@example.com / Admin123!
- Ganti segera di production.

## Catatan Integrasi WhatsApp
- Set WHATSAPP_API_URL dan WHATSAPP_API_TOKEN sesuai gateway kamu.
- Webhook verify token: WHATSAPP_WEBHOOK_VERIFY_TOKEN
- Payload webhook dapat disesuaikan (adapter) bila format gateway berbeda; saat ini generic JSON.

## Roadmap Lanjutan
- Implementasi kirim media dan template (controller)
- Validasi request (validator)
- Pagination & filter yang lebih lengkap
- Observability: structured logging, metrics
- CI/CD, kontainerisasi produksi, OpenAPI/Swagger
