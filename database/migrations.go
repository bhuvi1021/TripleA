package database

import (
	"database/sql"
	"log"
)

// RunMigrations executes database migrations
func RunMigrations(db *sql.DB) error {
	// Create accounts table
	accountsTable := `
	CREATE TABLE IF NOT EXISTS accounts (
		id SERIAL PRIMARY KEY,
		account_id BIGINT UNIQUE NOT NULL,
		balance DECIMAL(20,5) NOT NULL DEFAULT 0,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		deleted_at TIMESTAMP DEFAULT NULL
	);
	
    CREATE INDEX IF NOT EXISTS idx_accounts_id ON accounts(account_id);
	`

	// Create transactions table
	transactionsTable := `
	CREATE TABLE IF NOT EXISTS transactions (
		id SERIAL PRIMARY KEY,
		account_id BIGINT NOT NULL,
		amount DECIMAL(20,5) NOT NULL,
	    currency_code VARCHAR(3) NOT NULL,
	    available_balance DECIMAL(20,5) NOT NULL,
		is_credit BOOL,
		reference VARCHAR(50) NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		deleted_at TIMESTAMP DEFAULT NULL,
		FOREIGN KEY (account_id) REFERENCES accounts(account_id)
	);
	
    CREATE INDEX IF NOT EXISTS idx_transactions_account ON transactions(account_id);
	`

	if _, err := db.Exec(accountsTable); err != nil {
		return err
	}

	if _, err := db.Exec(transactionsTable); err != nil {
		return err
	}

	log.Printf("Database migrations completed successfully")
	return nil
}
