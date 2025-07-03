# TripleA – Internal Money Transfer System

TripleA is a lightweight internal money transfer system written in Go. It handles account creation, balance updates, and atomic transactions between accounts. Designed with a layered architecture and clear separation of concerns.

## Features

- 🚀 RESTful APIs for account and transaction management
- 🔒 Atomic balance updates using SQL transactions
- 🧪 Table-driven unit tests and mocks
- 🧩 Layered architecture (handler → service → repository)
- 🐘 PostgreSQL as the underlying database

---

## Project Structure

```
/cmd/               # Main application entry
/internal/
  /handlers/        # HTTP layer
  /service/         # Business logic
  /repository/      # DB interactions
  /models/          # DB models
  /errors/          # Custom error types and HTTP mapping
/tests/             # Test files (table-driven and mocks)
```

---

## Setup

### Prerequisites

- Go 1.20+
- PostgreSQL
- [Go Mock](https://github.com/golang/mock) for tests

### Install

```bash
git clone https://github.com/your-org/TripleA.git
cd TripleA
go mod tidy
```

### Run the server

```bash
go run main.go
```

---

## API Endpoints

### Account

| Method | Endpoint           | Description        |
|--------|--------------------|--------------------|
| POST   | `/accounts`        | Create an account  |
| GET    | `/accounts/{id}`   | Get account details|

### Transaction

| Method | Endpoint            | Description               |
|--------|---------------------|---------------------------|
| POST   | `/transactions`     | Transfer between accounts |

---

## 🧪 API Usage Examples

### ✅ POST /accounts

**Request:**
```json
curl --location 'http://localhost:8080/accounts' \
--header 'Content-Type: application/json' \
--data '{
"account_id": 999,
"initial_balance": "1000"
}'
```

**Response:**
```json
{}
```

**Error Responses:**
```json
{ "error_message": "account already exists"}
```
```json
{ "error_message": "invalid amount"}
```
```json
{ "error_message": "invalid json format"}
```
```json
{ "error_message": "invalid account id"}
```

---

### ✅ GET /accounts/{account_id}

**Example:**
```
GET /accounts/1001
```

**Response:**
```json
curl --location 'http://localhost:8080/accounts/123'
```
**Response (Inactive/Deleted account):**
```json
{
  "account_id": 123,
  "balance": "100.00000"
}
```
---

### ✅ POST /transactions

**Request:**
```json
{
  "source_account_id": 1001,
  "destination_account_id": 1002,
  "amount": "250.00"
}
```

**Success Response:**
```json
{
  "source_account_id": 1001,
  "available_balance": "1250.00000"
}
```

**Error Response when source account is invalid:**
```json
{
  "error_message": "sender account not found" 
}
```

**Error Response when destination account is invalid:**
```json
{
  "error_message": "receiver account not found" 
}
```

**Error Response when there is insufficient fund in sender account:**
```json
{
  "error_message": "insufficient funds in sender account"
}
```

----


## Testing

```bash
go test ./... -v
```

To regenerate mocks:
```bash
mockgen -source=internal/repository/interfaces.go -destination=internal/repository/mocks/mock_repository.go -package=mocks
```

---

## License

MIT License

---

## Maintainers

- [@bhuvi1021](https://github.com/bhuvi1021)
