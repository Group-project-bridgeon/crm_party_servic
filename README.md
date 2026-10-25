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
│   │   ├── middleware.go     # Fiber middlewares (recover, requestid, logger, cors)
│   │   ├── health_routes.go  # Health & status endpoints
│   │   ├── party_routes.go   # Parties & 360 overview routes
│   │   ├── customer_routes.go# Customers, composite onboarding, user accounts
│   │   ├── supplier_routes.go# Suppliers, composite onboarding, vendor accounts
│   │   ├── organization_routes.go # B2B profile routes
│   │   ├── individual_routes.go   # B2C profile routes
│   │   ├── address_routes.go # Address management routes
│   │   ├── contact_routes.go # Contacts & stakeholder routes
│   │   └── note_routes.go    # CRM interaction note routes
│   ├── handler/http/         # Primary HTTP adapters (Fiber handlers & DTOs)
│   ├── service/              # Use cases & business logic
│   ├── repository/           # Secondary data access ports & PostgreSQL implementations
│   ├── domain/               # Core entities, value objects, and domain invariants
│   ├── db/                   # sqlc generated queries & models
│   └── config/               # Environment configuration loader
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
| `ENV` | Environment name (`development`, `production`) | `development` |
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
docker compose up -d postgres

# Start party service locally
go run cmd/service/main.go
```

---

## 📡 API Reference

### Health Check
| Method | Endpoint | Description |
| :--- | :--- | :--- |
| `GET` | `/health` | Service health status |

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

## 🧪 Testing

The test suite follows standard Go architectural conventions:
- **Unit & Handler Tests**: Co-located directly alongside domain logic in `internal/service/` and HTTP handlers in `internal/handler/http/`.
- **End-to-End API Integration Tests**: Located in root `test/api/` with a high-fidelity in-memory store simulating the complete Fiber HTTP router, request serialization, and response codes.

Run all tests across the entire repository:
```bash
go test -v ./...
```

Run only the root End-to-End API integration test suite:
```bash
go test -v ./test/api/...
```

Run only unit tests:
```bash
go test -v ./internal/...
```

Regenerate SQL queries:
```bash
/home/yadu/go/bin/sqlc generate
```

Compile binary:
```bash
go build -o bin/party_service ./cmd/service
```
