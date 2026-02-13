# CallFlow

Full-stack call management platform — Go API backend + Flutter mobile app.

## Prerequisites

- Go 1.24+
- Flutter 3.16+ / Dart 3.2+
- PostgreSQL 16+
- [golang-migrate](https://github.com/golang-migrate/migrate) (`brew install golang-migrate`)
- [sqlc](https://sqlc.dev/) (`brew install sqlc`)
- Docker & Docker Compose (optional, for local Postgres)

## Project Structure

```
callflow/
├── api/                     # Go backend
│   ├── cmd/api/main.go      # Entry point
│   ├── config/              # DB config
│   ├── internal/
│   │   ├── api/             # Router, handlers, middleware
│   │   ├── domain/          # Business entities
│   │   ├── repository/      # Data access layer
│   │   ├── service/         # Business logic
│   │   └── sql/
│   │       ├── migrations/  # SQL migration files
│   │       ├── queries/     # sqlc query definitions
│   │       └── db/          # sqlc generated code
│   ├── sqlc.yaml
│   ├── Dockerfile
│   └── .env.example
├── callflow_app/            # Flutter mobile app
├── docker-compose.yml
└── docs/
```

---

## API Backend

### Setup

```bash
cd api
cp .env.example .env
# Edit .env with your values (JWT_SECRET is required)
```

### Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `PORT` | `8080` | Server port |
| `DB_HOST` | `localhost` | PostgreSQL host |
| `DB_PORT` | `5432` | PostgreSQL port |
| `DB_USER` | `postgres` | DB username |
| `DB_PASSWORD` | `postgres` | DB password |
| `DB_NAME` | `callflow_db` | DB name |
| `DB_SSL_MODE` | `disable` | SSL mode |
| `JWT_SECRET` | — | **Required.** JWT signing secret |
| `MSG91_AUTH_KEY` | — | MSG91 auth key (optional) |
| `MSG91_OTP_TEMPLATE_ID` | — | MSG91 OTP template (optional) |
| `CORS_ALLOW_ORIGINS` | `http://localhost:3000` | Allowed CORS origins |

### Database

```bash
# Run all up migrations
migrate -path api/internal/sql/migrations -database 'YOUR_DATABASE_URL' up

# Rollback all migrations
migrate -path api/internal/sql/migrations -database 'YOUR_DATABASE_URL' down

# Drop everything and re-migrate
migrate -path api/internal/sql/migrations -database 'YOUR_DATABASE_URL' drop -f
migrate -path api/internal/sql/migrations -database 'YOUR_DATABASE_URL' up

# Check current migration version
migrate -path api/internal/sql/migrations -database 'YOUR_DATABASE_URL' version

# Force a specific version (to fix dirty state)
migrate -path api/internal/sql/migrations -database 'YOUR_DATABASE_URL' force VERSION
```

### Code Generation (sqlc)

```bash
cd api
sqlc generate
```

### Run

```bash
cd api
go run cmd/api/main.go
```

### Build & Run with Docker

```bash
cd api
docker build -t callflow-api .
docker run -p 8080:8080 --env-file .env callflow-api
```

---

## Flutter App

### Setup

```bash
cd callflow_app
flutter pub get
```

### Code Generation

The app uses build_runner for Riverpod, Drift, Freezed, and JSON serialization:

```bash
cd callflow_app
dart run build_runner build --delete-conflicting-outputs
```

### Run (debug)

```bash
cd callflow_app

# List available devices
flutter devices

# Run on connected device / emulator
flutter run

# Run on specific device
flutter run -d <device_id>
```

### Build APK (Android)

```bash
cd callflow_app

# Debug APK
flutter build apk --debug

# Release APK
flutter build apk --release

# Split APKs per ABI (smaller size)
flutter build apk --split-per-abi --release
```

Output: `callflow_app/build/app/outputs/flutter-apk/`

### Build App Bundle (Play Store)

```bash
cd callflow_app
flutter build appbundle --release
```

Output: `callflow_app/build/app/outputs/bundle/release/`

### Install APK on Device

```bash
# Install debug APK
flutter install

# Or directly via adb
adb install callflow_app/build/app/outputs/flutter-apk/app-release.apk
```

### Clean

```bash
cd callflow_app
flutter clean
flutter pub get
dart run build_runner build --delete-conflicting-outputs
```
