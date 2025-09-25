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

- **Note**: Some files/folders (for example `.idea/`, `kyber-server`, `coverage/`) may be present only in development or after running build/tests; they are listed here to reflect the current workspace state.

- **main.go**: Only starts the server, no business logic.
- **internal/config**: Loads configuration from environment variables.
- **internal/handlers**: HTTP layer, no business logic; errors are logged and safe for clients.
- **internal/kybertransit**: Kyber key management, encryption, decryption. All functions documented and errors are informative.
- **internal/routes**: Central place for route templates and route names used for named routes.
- **internal/server**: Router setup; routes are registered and named for safe URL building in tests.

## Named routes and tests

This project uses gorilla/mux named routes. Routes are registered with `.Name(...)` in `internal/server.NewRouter()` and tests build concrete URLs using `router.Get("routeName").URL("name", value)`. This keeps route registration and URL building in sync (single source of truth).

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

## Deployment & Run

This section has quick instructions to build and run the standalone service locally (useful for development and testing)

### Local build and run

You can build and run the service directly with Go or use the provided `Makefile`.

1) Build with Go:

```bash
# Build binary
go build -o kyber-server main.go

# Run binary (default port :8080)
./kyber-server
```

2) Using the Makefile (recommended for development):

```bash
# Run unit and integration tests
make test

# Build binary
make build

# Run the built binary
make run
```

3) Running with a custom port (example):

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
