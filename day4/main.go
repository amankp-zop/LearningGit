package main

import "fmt"

type BankAccount struct {
	Owner   string
	Balance float64
}

// Method with a value receiver — won't modify the original struct
func (b BankAccount) DisplayBalance() {
	fmt.Printf("Owner: %s, Balance: $%.2f\n", b.Owner, b.Balance)
}

// Pointer receiver — modifies the original struct
func (b *BankAccount) Deposit(amount float64) {
	b.Balance += amount
	fmt.Printf("Deposited $%.2f to %s's account.\n", amount, b.Owner)
}

// Pointer receiver — modifies the original struct if sufficient balance
func (b *BankAccount) Withdraw(amount float64) {
	if b.Balance >= amount {
		b.Balance -= amount
		fmt.Printf("Withdrew $%.2f from %s's account.\n", amount, b.Owner)
	} else {
		fmt.Printf("Insufficient funds for %s. Withdrawal of $%.2f failed.\n", b.Owner, amount)
	}
}

func main() {
	account := BankAccount{Owner: "Alice", Balance: 100.0}

	// Initial balance
	account.DisplayBalance()

	// Try deposit and withdraw
	account.Deposit(50.0)
	account.DisplayBalance()

	account.Withdraw(30.0)
	account.DisplayBalance()

	// Withdraw more than balance
	account.Withdraw(150.0)
	account.DisplayBalance()
}
