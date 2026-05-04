# Horoscope Vault

A REST API service built for Attack/Defense CTF cybersecurity competitions. The service implements an encrypted record storage and contains intentional vulnerabilities in accordance with competition rules.



## API

| Method | Path | Access | Description |
|--------|------|--------|-------------|
| `POST` | `/signup` | Public | User registration |
| `POST` | `/login` | Public | Authentication, returns JWT |
| `POST` | `/vault` | JWT | Create an encrypted record |
| `GET` | `/vault/{sign}/{date}` | JWT | Retrieve a record |

## Build

```bash
# 1. copy and fill in the environment variables
cp .env.example .env

# 2. generate RSA keys (see secrets/README.md)
openssl genrsa -out secrets/jwt_private.pem 2048
openssl rsa -in secrets/jwt_private.pem -pubout -out secrets/jwt_public.pem

# 3. run
docker compose up --build
```

Service is available at `http://localhost:8888`.


## Vulnerabilities

| Vulnerability | Location | Type |
|---------------|----------|------|
| SQL Injection | `internal/models/models.go` — `FindByLogin`, `GetEntry` | `fmt.Sprintf` without parameterization |
| Static IV | `internal/crypto/aes.go` — `fixedIV` | IV reuse in AES-CBC |
| IDOR | `internal/handlers/valut.go` — `GetEntryHandler` | No `user_id` check on read |
