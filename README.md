# Kyber Transit API (Coding Challenge 2)

This project implements a Vault-like Transit Secrets Engine with support for the post-quantum cryptographic algorithm CRYSTALS-Kyber (Kyber-1024), using Go and the [Cloudflare CIRCL library](https://github.com/cloudflare/circl).

## Features
- **Key Generation**: Create Kyber key pairs.
- **Encryption**: Encrypt data using Kyber public key (with a demo symmetric layer).
- **Decryption**: Decrypt data using Kyber private key.
- **REST API**: Endpoints compatible with typical Vault Transit API style.
- **Unit & Integration Tests**: High coverage, edge cases, error handling.
- **Clean Architecture**: Separation of HTTP, business logic, and bootstrap layers.
- **Configurable Port**: Server port is configurable via environment variable.
- **Health Check**: GET `/health` returns 200 OK.

## Architecture

```
ElevexaCodingChallenge2/
├── .gitignore             # ignore rules for build artifacts and IDE files
├── LICENSE                # MIT license
├── main.go                # Bootstrap, only server startup (with graceful shutdown)
├── Makefile               # Common dev tasks: test, testrace, coverage, build, run
├── go.mod, go.sum         # Dependencies
├── README.md              # Documentation (this file)
└── internal/
    ├── config/
    │   ├── config.go
    │   └── config_test.go
    ├── handlers/
    │   ├── handlers.go      # HTTP handlers + KeyStoreManager (thread-safe in-memory store)
    │   └── handlers_test.go # Table-driven tests (use handlers.ResetKeyStore for isolation)
    ├── kybertransit/
    │   ├── kyber.go         # Kyber logic (CIRCL), SECURITY WARNING about XOR (demo-only)
    │   └── kyber_test.go    # Table-driven tests, edge cases
    ├── routes/
    │   └── routes.go
    └── server/
        ├── server.go        # Router setup (includes GET /health)
        └── server_test.go
```

- **main.go**: Starts the server with graceful shutdown.
- **internal/config**: Loads configuration from env; validates port format (":8080").
- **internal/handlers**: HTTP handlers; encapsulated key storage via KeyStoreManager; errors logged and safe for clients.
- **internal/kybertransit**: Kyber key management, encryption, decryption. Clear error wrapping; defensive checks; demo XOR note.
- **internal/routes**: Central place for route templates and names.
- **internal/server**: Router setup; named routes; includes a health check endpoint.

## API Endpoints

All endpoints are POST and accept/return JSON unless noted.

### 1. Create a new Kyber key pair
- **POST** `/transit/keys/{name}`
- Request: `{}`
- Response:
```json
{
  "name": "my-key",
  "public_key": "...base64..."
}
```

### 2. Encrypt data with Kyber
- **POST** `/transit/encrypt/{name}`
- Request:
```json
{ "plaintext": "...base64 or text..." }
```
- Response:
```json
{ "ciphertext": "...base64...", "encdata": "...base64..." }
```

### 3. Decrypt data with Kyber
- **POST** `/transit/decrypt/{name}`
- Request:
```json
{ "ciphertext": "...base64...", "encdata": "...base64..." }
```
- Response:
```json
{ "plaintext": "...base64 or text..." }
```

### 4. Health check
- **GET** `/health`
- Response: `200 OK`, body: `ok`

## Configuration

Server port via `KYBER_SERVER_PORT` (default: `:8080`).

Windows (cmd.exe):
```
set KYBER_SERVER_PORT=:9090
```

Unix-like (bash):
```
export KYBER_SERVER_PORT=:9090
```

## Build, Run, and Test

Using Makefile (recommended):

```
make test       # Run all unit and integration tests
make testrace   # Run tests with race detector
make coverage   # Generate coverage report
make build      # Build the binary (kyber-server.exe on Windows)
make run        # Build and run the server
make clean      # Remove binaries
```

With Go directly:
```
go build -o kyber-server.exe main.go   # Windows
go build -o kyber-server main.go       # Unix
```

## Testing Notes
- Tests are table-driven and cover success and failure scenarios.
- Handlers use an in-memory key store encapsulated by `KeyStoreManager`.
- Test isolation: call `handlers.ResetKeyStore()` before tests to clear the in-memory store.
- SECURITY: The symmetric layer uses XOR with the Kyber shared secret for demonstration only. Do not use in production; switch to a KEM→DEM construction with an AEAD cipher and secure key storage.

## License
MIT — see `LICENSE`.
