package main

import "fmt"

// Account represents a user account with an ID, username, and email
type Account struct {
	AccID    int
	Username string
	Email    string
}

// NewAccount creates a new account
func NewAccount(id int, username, email string) *Account {
	return &Account{
		AccID:    id,
		Username: username,
		Email:    email,
	}
}

// Display prints the account details
func (a *Account) Display() {
	fmt.Printf("Account ID: %d, Username: %s, Email: %s\n", a.AccID, a.Username, a.Email)
}
