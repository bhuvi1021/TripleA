# TripleA – Internal Money Transfer System

TripleA is a lightweight internal money transfer system written in Go. It handles account creation, balance updates, and money transfer transactions between accounts. Designed with a layered architecture.

## Features

-  RESTful APIs for account and transaction management
-  Atomic balance updates using SQL transactions
-  Table-driven unit tests and mocks
-  Layered architecture (handler → service → repository)
-  PostgreSQL as the underlying database

---

## Project Assumptions
- The amount in the payload is considered as string as mentioned in the requirement, though using float would be more appropriate
- The account id in the payload is considered as int as mentioned in the requirement, though using string would be more appropriate
- No transaction activities are recorded


## Project Explanation

# Accounts
- For Get Accounts api, we will also return IsDeleted flag if the account is soft-deleted. This is to support backward support for the user's past activities or transactions
- For Create Transaction api, we will create two records for each money transfer. for eg, if money is transfered from Account 123 to 124, then two record will be recorded 1) debit entry for account 123 and 2) credit entry for account 124. This is to support the ledger/transaction history for the user.
- To link these two credit and debit transaction, reference can be used. 
- We also maintain available_balance in transaction table to support ledger functions.

## Project Structure

```
/main.go                # Main application entry
/config/                # Configs
/database/              # Database migraation
/routes/                # Routes for registration
/internal/
  /server/handlers/     # HTTP layer
  /service/             # Business logic
  /repository/          # DB interactions
  /models/              # DB models
  /errors/              # Custom error types and HTTP response code mapping
```

---

## Database Schema
accounts Table


## Setup

### Prerequisites

- Go 1.20+
- PostgreSQL
- [Go Mock](https://github.com/golang/mock) for tests
- Update the postgres password and port number in config/config.go line number #15
``` bash 
DatabaseURL: getEnv("DATABASE_URL", "postgres://postgres:yourpassword@localhost:5432/postgres?sslmode=disable"),
```

### Install

```bash
git clone https://github.com/bhuvi1021/TripleA.git
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
