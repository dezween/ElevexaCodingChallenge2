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

## Architecture

```
ElevexaCodingChallenge2/
├── .gitignore             # ignore rules for build artifacts and IDE files
├── LICENSE                # MIT license
├── main.go                # Bootstrap, only server startup
├── Makefile               # Common dev tasks: test, testrace, build, run
├── go.mod, go.sum         # Dependencies
├── README.md              # Documentation (this file)
└── internal/
    ├── config/
    │   ├── config.go
    │   └── config_test.go
    ├── handlers/
    │   ├── handlers.go
    │   └── handlers_test.go
    ├── kybertransit/
    │   ├── kyber.go
    │   └── kyber_test.go
    ├── routes/
    │   └── routes.go
    └── server/
        ├── server.go
        └── server_test.go
```

- **main.go**: Only starts the server, no business logic.
- **internal/config**: Loads configuration from environment variables.
- **internal/handlers**: HTTP layer, no business logic; errors are logged and safe for clients.
- **internal/kybertransit**: Kyber key management, encryption, decryption. All functions documented and errors are informative.
- **internal/routes**: Central place for route templates and route names used for named routes.
- **internal/server**: Router setup; routes are registered and named for safe URL building in tests.
- **.gitignore**: Excludes binaries, coverage, IDE files, etc.
- **LICENSE**: MIT license.

## API Endpoints

All endpoints are POST and accept/return JSON.

### 1. Create a new Kyber key pair
- **POST** `/transit/keys/{name}`
- **Request body:** `{}`
- **Response:**
```json
{
  "name": "my-key",
  "public_key": "...base64..."
}
```

**Example:**
```bash
curl -X POST http://localhost:8080/transit/keys/my-key -H "Content-Type: application/json" -d "{}"
```

### 2. Encrypt data with Kyber
- **POST** `/transit/encrypt/{name}`
- **Request body:**
```json
{
  "plaintext": "...base64..."
}
```
- **Response:**
```json
{
  "ciphertext": "...base64..."
}
```

**Example:**
```bash
curl -X POST http://localhost:8080/transit/encrypt/my-key -H "Content-Type: application/json" -d '{"plaintext":"SGVsbG8gd29ybGQ="}'
```

### 3. Decrypt data with Kyber
- **POST** `/transit/decrypt/{name}`
- **Request body:**
```json
{
  "ciphertext": "...base64..."
}
```
- **Response:**
```json
{
  "plaintext": "...base64..."
}
```

**Example:**
```bash
curl -X POST http://localhost:8080/transit/decrypt/my-key -H "Content-Type: application/json" -d '{"ciphertext":"..."}'
```

## Configuration

The server port can be configured via the `KYBER_SERVER_PORT` environment variable (default: `:8080`).

Example (Windows cmd):
```
set KYBER_SERVER_PORT=:9090
```

Example (Unix-like):
```
export KYBER_SERVER_PORT=:9090
```

## Build, Run, and Test

You can build and run the service directly with Go or use the provided `Makefile`.

### Using Makefile (recommended)

```cmd
make test        # Run all unit and integration tests
make build       # Build binary (kyber-server.exe on Windows, kyber-server on Unix)
make run         # Build and run the server
make clean       # Remove binaries
```

### With Go directly

```cmd
go build -o kyber-server.exe main.go   # Windows
go build -o kyber-server main.go       # Unix
./kyber-server.exe                     # Windows
./kyber-server                         # Unix
```

### Run with a custom port

Windows (cmd.exe):
```cmd
set KYBER_SERVER_PORT=:9090
make run
```

Unix-like (bash):
```bash
export KYBER_SERVER_PORT=:9090
make run
```

## License

This project is licensed under the MIT License - see the included `LICENSE` file for details.

## References
- [CRYSTALS-Kyber](https://pq-crystals.org/kyber/)
- [Cloudflare CIRCL](https://github.com/cloudflare/circl)
