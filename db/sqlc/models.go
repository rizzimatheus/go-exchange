// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.16.0

package db

import (
	"time"

	"github.com/google/uuid"
)

type Account struct {
	ID        int64     `json:"id"`
	Owner     string    `json:"owner"`
	Balance   int64     `json:"balance"`
	Currency  string    `json:"currency"`
	CreatedAt time.Time `json:"created_at"`
}

type Ask struct {
	ID            int64  `json:"id"`
	Pair          string `json:"pair"`
	FromAccountID int64  `json:"from_account_id"`
	ToAccountID   int64  `json:"to_account_id"`
	Price         int64  `json:"price"`
	// it must be positive
	Amount    int64     `json:"amount"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

type Bid struct {
	ID            int64  `json:"id"`
	Pair          string `json:"pair"`
	FromAccountID int64  `json:"from_account_id"`
	ToAccountID   int64  `json:"to_account_id"`
	Price         int64  `json:"price"`
	// it must be positive
	Amount    int64     `json:"amount"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

type Entry struct {
	ID        int64 `json:"id"`
	AccountID int64 `json:"account_id"`
	// can be negative or positive
	Amount    int64     `json:"amount"`
	CreatedAt time.Time `json:"created_at"`
}

type Session struct {
	ID           uuid.UUID `json:"id"`
	Username     string    `json:"username"`
	RefreshToken string    `json:"refresh_token"`
	UserAgent    string    `json:"user_agent"`
	ClientIp     string    `json:"client_ip"`
	IsBlocked    bool      `json:"is_blocked"`
	ExpiresAt    time.Time `json:"expires_at"`
	CreatedAt    time.Time `json:"created_at"`
}

type Trade struct {
	ID               int64     `json:"id"`
	FirstTransferID  int64     `json:"first_transfer_id"`
	SecondTransferID int64     `json:"second_transfer_id"`
	CreatedAt        time.Time `json:"created_at"`
}

type Transfer struct {
	ID            int64 `json:"id"`
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	// it must be positive
	Amount    int64     `json:"amount"`
	CreatedAt time.Time `json:"created_at"`
}

type User struct {
	Username          string    `json:"username"`
	HashedPassword    string    `json:"hashed_password"`
	FullName          string    `json:"full_name"`
	Email             string    `json:"email"`
	PasswordChangedAt time.Time `json:"password_changed_at"`
	CreatedAt         time.Time `json:"created_at"`
}