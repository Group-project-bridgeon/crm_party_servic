# CRM Party Service

Enterprise Resource Planning (ERP) & Customer Relationship Management (CRM) microservice for managing **Parties** (Organizations & Individuals), **Customers**, **Suppliers**, **Addresses**, **Contacts**, **Stakeholder Contact Persons**, and **Interaction Notes**.

---

## 🏛️ Architecture & Domain Design

This service implements Martin Fowler's **Party Archetype Pattern**:
* **Party**: The central entity representing any legal organization or individual person.
* **Profiles**:
  * `organizations`: Extended attributes for business entities (registration number, industry, employee count, founded date).
  * `individuals`: Extended attributes for natural persons (first/last name, DOB, gender, national ID, job title).
* **Roles Played by Parties**:
  * `customers`: B2B / B2C customer roles with credit limits, payment terms, and status (`ACTIVE`, `BLOCKED`, `HOLD`).
  * `suppliers`: Vendor roles with ratings ($0.00 - 5.00$), payment methods, and ISO currency codes.
  * `customer_user_accounts` & `supplier_user_accounts`: Mapping ERP authentication users to customer & vendor portals.
* **Subsidiary Entities**:
  * `party_addresses`: Multi-address registry (Billing, Shipping, Office, Warehouse) with primary flags and coordinates.
  * `party_contacts`: Multi-channel communication channels (Email, Phone, Mobile, WhatsApp, Fax).
  * `party_contact_persons`: Key account managers, executives, and points of contact.
  * `party_notes`: CRM interaction and chronological activity logs.

### 📁 Clean Architecture Directory Layout
```
crm_party_service/
├── cmd/service/main.go       # Composition root, DB pool, router execution, graceful shutdown
├── internal/
│   ├── api/                  # Categorized API routing & HTTP middleware pipeline
│   │   ├── router.go         # Central route registration
│   │   ├── middleware.go     # Fiber middlewares (recover, requestid, slog, metrics, cors)
│   │   ├── health_routes.go  # Health (/health) & Prometheus metrics (/metrics) endpoints
│   │   ├── party_routes.go   # Parties & 360 overview routes
│   │   ├── customer_routes.go# Customers, composite onboarding, user accounts
│   │   ├── supplier_routes.go# Suppliers, composite onboarding, vendor accounts
│   │   ├── organization_routes.go # B2B profile routes
│   │   ├── individual_routes.go   # B2C profile routes
│   │   ├── address_routes.go # Address management routes
│   │   ├── contact_routes.go # Contacts & stakeholder routes
│   │   └── note_routes.go    # CRM interaction note routes
│   ├── logger/               # Enterprise log/slog structured logger & Fiber middleware
│   ├── metrics/              # Prometheus RED metrics collectors (RPS, Latency, In-Flight)
│   ├── handler/http/         # Primary HTTP adapters (Fiber handlers & DTOs)
│   ├── service/              # Use cases & business logic
│   ├── repository/           # Secondary data access ports & PostgreSQL implementations
│   ├── domain/               # Core entities, value objects, and domain invariants
│   ├── db/                   # sqlc generated queries & models
│   └── config/               # Environment configuration loader
├── monitoring/               # Observability infrastructure
│   ├── prometheus/           # Prometheus scrape configuration
│   └── grafana/              # Pre-provisioned datasources & Party Service dashboard
├── postman/                  # Postman collection & environment for API testing
└── test/api/                 # Dedicated black-box E2E integration test suite
```

---

## 🚀 Getting Started

### 1. Requirements
* **Go 1.27+**
* **PostgreSQL 16+**
* **Docker & Docker Compose** (optional)

### 2. Environment Variables

| Variable | Description | Default |
| :--- | :--- | :--- |
| `PORT` | HTTP server listening port | `:8080` |
| `DATABASE_URL` | PostgreSQL connection string | `postgres://postgres:postgres@localhost:5432/party?sslmode=disable` |
| `ENV` | Environment name (`development`, `production`, `staging`) | `development` |
| `LOG_LEVEL` | Log level (`DEBUG`, `INFO`, `WARN`, `ERROR`) | `INFO` |
| `SHUTDOWN_TIMEOUT_SECONDS` | Maximum time to wait for in-flight requests to complete during shutdown | `10` |

### 3. Running with Docker & Docker Compose

#### Option A: Full Stack with Docker Compose (Recommended)
Automatically initializes PostgreSQL with migrations and boots the Party service:
```bash
# Build and start both PostgreSQL and CRM Party Service
docker compose up --build -d

# View live logs
docker compose logs -f party-service

# Check health status
curl http://localhost:8080/health

# Stop stack
docker compose down
```

#### Option B: Build and Run Docker Container Individually
```bash
# Build image
docker build -t crm-party-service:latest .

# Run container (connecting to host or external Postgres)
docker run -p 8080:8080 \
  -e DATABASE_URL="postgres://postgres:postgres@host.docker.internal:5432/party?sslmode=disable" \
  -e ENV="production" \
  crm-party-service:latest
```

#### Option C: Local Development with Standalone Postgres
```bash
# Start PostgreSQL database only
make docker-db

# Start party service locally
make run
```

---

## 📊 Observability, Logging & Monitoring

The service incorporates production-grade telemetry conforming to cloud-native microservice standards:

### 1. Structured Logging (`log/slog`)
* **Production/Staging**: Emits structured JSON logs containing `timestamp`, `level`, `service`, `env`, `request_id`, `route`, `latency_ms`, and `ip`. Seamlessly ingestible by Grafana Loki, Datadog, AWS CloudWatch, or ELK.
* **Development**: Colorized, human-readable text logs with source code file/line references.
* **Correlation**: Injects `request_id` from Fiber's `requestid` middleware into every log entry for tracing across distributed services.
* **Smart Noise Reduction**: High-frequency scrape requests (`/health`, `/metrics`) are logged at `DEBUG` level to prevent log pollution.

### 2. Prometheus Metrics (`/metrics`)
Standard RED (Rate, Errors, Duration) metrics collected automatically:
* `crm_party_service_http_requests_total`: Counter partitioned by `method`, `path` (route template), and HTTP `status`.
* `crm_party_service_http_request_duration_seconds`: High-resolution histogram tracking request latency percentiles (p50, p95, p99).
* `crm_party_service_http_requests_in_flight`: Gauge tracking active concurrent HTTP connections.
* **Go Runtime Telemetry**: Memory allocation (Heap, Stack, Sys), GC pause duration, and active goroutines.

### 3. Grafana Dashboard
A pre-provisioned Grafana dashboard is included at `monitoring/grafana/dashboards/crm_party_service.json`:
* **URL**: [http://localhost:3000](http://localhost:3000)
* **Credentials**: `admin` / `admin`
* **Pre-configured Panels**:
  * Realtime Service Health & Status (UP/DOWN)
  * In-Flight Requests Gauge
  * 5xx Error Rate Percentage (%)
  * Request Throughput by Status Code (2xx, 4xx, 5xx)
  * p50 / p95 / p99 Latency Curves
  * Go Goroutines & OS Thread Counts
  * Go Heap Memory Consumption

To spin up the monitoring stack:
```bash
make docker-monitoring
# or: docker compose up -d prometheus grafana
```

---

## 📡 API Reference

### Health & Telemetry
| Method | Endpoint | Description |
| :--- | :--- | :--- |
| `GET` | `/health` | Service health status |
| `GET` | `/metrics` | Prometheus scrape endpoint (RED metrics & Go runtime) |

### 1. Parties
| Method | Endpoint | Description |
| :--- | :--- | :--- |
| `POST` | `/api/v1/parties` | Create a new party (`ORGANIZATION` or `INDIVIDUAL`) |
| `GET` | `/api/v1/parties` | List parties (query params: `limit`, `offset`, `party_type`, `is_active`, `search`) |
| `GET` | `/api/v1/parties/:id` | Get party by ID |
| `PUT` | `/api/v1/parties/:id` | Update party details |
| `PATCH` | `/api/v1/parties/:id/status` | Activate or deactivate party |
| `DELETE` | `/api/v1/parties/:id` | Deactivate party (soft delete) |
| `GET` | `/api/v1/parties/:id/360` | **Party 360° View**: Unified response with profile, roles, addresses, and contacts |

### 2. Organization Profiles
| Method | Endpoint | Description |
| :--- | :--- | :--- |
| `POST` | `/api/v1/organizations` | Create organization profile for a party |
| `GET` | `/api/v1/organizations/:party_id` | Get organization profile |
| `PUT` | `/api/v1/organizations/:party_id` | Update organization profile |

### 3. Individual Profiles
| Method | Endpoint | Description |
| :--- | :--- | :--- |
| `POST` | `/api/v1/individuals` | Create individual profile for a party |
| `GET` | `/api/v1/individuals/:party_id` | Get individual profile |
| `PUT` | `/api/v1/individuals/:party_id` | Update individual profile |

### 4. Customers
| Method | Endpoint | Description |
| :--- | :--- | :--- |
| `POST` | `/api/v1/customers` | Register an existing party as a customer |
| `POST` | `/api/v1/customers/composite` | **One-Shot Onboarding**: Create Party + Profile + Customer + Address atomically |
| `GET` | `/api/v1/customers` | List customers (filters: `customer_group`, `credit_status`, `search`, pagination) |
| `GET` | `/api/v1/customers/:id` | Get customer by customer ID |
| `GET` | `/api/v1/customers/party/:party_id` | Get customer by party ID |
| `PUT` | `/api/v1/customers/:id` | Update customer credit terms, status, and reps |
| `POST` | `/api/v1/customers/:id/user-accounts` | Link an auth user to customer portal |
| `GET` | `/api/v1/customers/:id/user-accounts` | List customer user portal accounts |
| `PATCH` | `/api/v1/customers/user-accounts/:account_id/status` | Activate/deactivate user account |

### 5. Suppliers
| Method | Endpoint | Description |
| :--- | :--- | :--- |
| `POST` | `/api/v1/suppliers` | Register an existing party as a supplier |
| `POST` | `/api/v1/suppliers/composite` | **One-Shot Onboarding**: Create Party + Profile + Supplier + Address atomically |
| `GET` | `/api/v1/suppliers` | List suppliers (filters: `supplier_group`, `currency_code`, `search`, pagination) |
| `GET` | `/api/v1/suppliers/:id` | Get supplier by supplier ID |
| `GET` | `/api/v1/suppliers/party/:party_id` | Get supplier by party ID |
| `PUT` | `/api/v1/suppliers/:id` | Update supplier rating, currency, payment terms |
| `POST` | `/api/v1/suppliers/:id/user-accounts` | Link an auth user to vendor portal |
| `GET` | `/api/v1/suppliers/:id/user-accounts` | List supplier vendor portal accounts |
| `PATCH` | `/api/v1/suppliers/user-accounts/:account_id/status` | Activate/deactivate user account |

### 6. Party Addresses
| Method | Endpoint | Description |
| :--- | :--- | :--- |
| `POST` | `/api/v1/parties/:party_id/addresses` | Add an address (Billing, Shipping, Office, Warehouse) |
| `GET` | `/api/v1/parties/:party_id/addresses` | List addresses for party (filter: `address_type`) |
| `GET` | `/api/v1/addresses/:address_id` | Get address by ID |
| `PUT` | `/api/v1/addresses/:address_id` | Update address |
| `DELETE` | `/api/v1/addresses/:address_id` | Delete address |

### 7. Party Contacts
| Method | Endpoint | Description |
| :--- | :--- | :--- |
| `POST` | `/api/v1/parties/:party_id/contacts` | Add communication channel (Email, Phone, WhatsApp) |
| `GET` | `/api/v1/parties/:party_id/contacts` | List communication channels (filter: `contact_type`) |
| `GET` | `/api/v1/contacts/:contact_id` | Get contact by ID |
| `PUT` | `/api/v1/contacts/:contact_id` | Update contact |
| `DELETE` | `/api/v1/contacts/:contact_id` | Delete contact |

### 8. Party Contact Persons (Account Managers / Stakeholders)
| Method | Endpoint | Description |
| :--- | :--- | :--- |
| `POST` | `/api/v1/parties/:party_id/contact-persons` | Add contact person / executive POC |
| `GET` | `/api/v1/parties/:party_id/contact-persons` | List contact persons |
| `GET` | `/api/v1/contact-persons/:person_id` | Get contact person by ID |
| `PUT` | `/api/v1/contact-persons/:person_id` | Update contact person |
| `DELETE` | `/api/v1/contact-persons/:person_id` | Delete contact person |

### 9. Party Notes
| Method | Endpoint | Description |
| :--- | :--- | :--- |
| `POST` | `/api/v1/parties/:party_id/notes` | Add CRM interaction note |
| `GET` | `/api/v1/parties/:party_id/notes` | List chronological CRM notes |
| `GET` | `/api/v1/notes/:note_id` | Get note by ID |
| `DELETE` | `/api/v1/notes/:note_id` | Delete note |

---

## 📮 API Testing with Postman

A complete Postman collection with 50 automated tests and an environment file is included in the [`postman/`](file:///home/yadu/project/Grp/crm_party_servic/postman) directory:

* [`postman/CRM_Party_Service.postman_collection.json`](file:///home/yadu/project/Grp/crm_party_servic/postman/CRM_Party_Service.postman_collection.json): Full collection with test assertions, pre-configured payloads, and automatic variable chaining.
* [`postman/CRM_Party_Service.postman_environment.json`](file:///home/yadu/project/Grp/crm_party_servic/postman/CRM_Party_Service.postman_environment.json): Local environment configuration (`baseUrl`: `http://localhost:8080`).

### Quick Start with Postman:
1. Open **Postman** and click **Import**.
2. Select both files inside the `postman/` directory.
3. In the top-right corner, choose the **CRM Party Service - Local** environment.
4. Execute individual requests or use Postman's **Run collection** to run all 50 tests in sequence. Dynamic IDs (`party_id`, `customer_id`, etc.) are captured and propagated automatically between requests.

---

## 🛠️ Makefile & Developer Commands

A comprehensive [`Makefile`](file:///home/yadu/project/Grp/crm_party_servic/Makefile) is provided to streamline development, testing, database code generation, and container lifecycle:

```bash
# View available Makefile targets
make help

# Build statically-linked binary to bin/party-service
make build

# Run service locally
make run

# Run all unit and E2E integration tests
make test

# Run unit tests only
make test-unit

# Run E2E API integration tests only
make test-e2e

# Generate test coverage report (generates coverage.html)
make test-coverage

# Code formatting & static analysis
make fmt
make vet

# Regenerate SQL queries via sqlc
make sqlc

# Docker Compose shortcuts
make docker-up    # Build and boot full stack (Postgres + Party Service)
make docker-db    # Start only PostgreSQL container
make docker-down  # Stop and remove containers
make docker-logs  # Stream container logs
make docker-ps    # View service status
```

---

## 🧪 Testing

The test suite follows standard Go architectural conventions:
- **Unit & Handler Tests**: Co-located directly alongside domain logic in `internal/service/` and HTTP handlers in `internal/handler/http/`.
- **End-to-End API Integration Tests**: Located in root `test/api/` with a high-fidelity in-memory store simulating the complete Fiber HTTP router, request serialization, and response codes.

Run all tests across the entire repository:
```bash
make test
# or: go test -v ./...
```

Run only the root End-to-End API integration test suite:
```bash
make test-e2e
# or: go test -v ./test/api/...
```

Run only unit tests:
```bash
make test-unit
# or: go test -v ./internal/...
```

Regenerate SQL queries:
```bash
make sqlc
# or: /home/yadu/go/bin/sqlc generate
```

Compile binary:
```bash
make build
# or: go build -o bin/party-service ./cmd/service
```
