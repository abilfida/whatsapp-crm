# WhatsApp CRM (Go + Fiber + GORM + MariaDB)

Aplikasi REST API untuk CRM dengan integrasi WhatsApp (send chat + webhook). Framework: Go Fiber, ORM: GORM, DB: MariaDB, Cache/Queue: Redis.

## Fitur Utama
- Auth JWT (login/register/profile/change-password)
- Role: admin, agent, supervisor
- Manajemen User (admin-only)
- Manajemen Customer
- Percakapan (assign agent, status, priority, notes)
- Pesan (list per percakapan, kirim text; media, template, audio, video, upload multipart)
- Webhook endpoint (verify + handler)

## Struktur Endpoint (prefix versi: /api/v1)
- Auth: POST /auth/login, POST /auth/register, GET /auth/profile, POST /auth/change-password
- Users (admin): GET/POST/PUT/DELETE /users
- Customers: GET/POST/PUT/DELETE /customers, GET /customers/:id
- Conversations: GET/POST/GET/:id, PUT /:id/{assign|status|priority|notes}
- Messages: GET /messages/conversation/:id, POST /messages/conversation/:id/{text|media|template}
- Upload: POST /messages/conversation/:id/upload (multipart -> simpan ke storage -> kirim WA)
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

## Storage Configuration (Local / S3 / GCS)
Set driver via STORAGE_DRIVER=local|s3|gcs

### Local (default, plain URL untuk dev)
```
STORAGE_DRIVER=local
STORAGE_BASE_PATH=./uploads/media
PUBLIC_BASE_URL=
STORAGE_SIGNED_URL_EXP_SECONDS=86400
```
- Jika PUBLIC_BASE_URL diisi (mis. http://localhost:8080/static), URL file akan memakai base tersebut.

### AWS S3 (private object + presigned URL 24 jam)
```
STORAGE_DRIVER=s3
AWS_REGION=ap-southeast-1
AWS_S3_BUCKET=your-bucket
AWS_S3_PREFIX=whatsapp-crm
AWS_ACCESS_KEY_ID=...
AWS_SECRET_ACCESS_KEY=...
AWS_S3_SIGNED_URL_EXP_SECONDS=86400
```

### Google Cloud Storage (private object + signed URL 24 jam)
```
STORAGE_DRIVER=gcs
GCS_BUCKET=your-bucket
GCS_PREFIX=whatsapp-crm
GOOGLE_APPLICATION_CREDENTIALS=/secrets/service-account.json
GCS_SIGNED_URL_EXP_SECONDS=86400
```
- Docker Compose: mount kredensial
```
  api:
    environment:
      GOOGLE_APPLICATION_CREDENTIALS: /secrets/service-account.json
    volumes:
      - ./secrets/gcs-sa.json:/secrets/service-account.json:ro
```

## Upload Validation
- Batas ukuran default: UPLOAD_MAX_SIZE=16777216 (16MB)
- Whitelist ekstensi berdasarkan media type:
  - ALLOWED_IMAGE_TYPES=jpg,jpeg,png,gif,webp
  - ALLOWED_DOCUMENT_TYPES=pdf,doc,docx,xls,xlsx,ppt,pptx
  - ALLOWED_AUDIO_TYPES=mp3,ogg,m4a,wav,aac
  - ALLOWED_VIDEO_TYPES=mp4,3gp,mov,avi,mkv
- Jika file melebihi batas atau di luar whitelist → HTTP 400

## Contoh Upload (curl)
```
curl -X POST \
  -H "Authorization: Bearer <your_token>" \
  -F "file=@/path/image.jpg" \
  -F "caption=Optional caption" \
  http://localhost:8080/api/v1/messages/conversation/<conversation_uuid>/upload
```

## Catatan Integrasi WhatsApp
- Set WHATSAPP_API_URL dan WHATSAPP_API_TOKEN sesuai gateway kamu.
- Webhook verify token: WHATSAPP_WEBHOOK_VERIFY_TOKEN
- Payload webhook dapat disesuaikan (adapter) bila format gateway berbeda.

## Roadmap Lanjutan
- OpenAPI/Swagger
- Observability (structured logging, metrics)
- CI/CD dan image registry
