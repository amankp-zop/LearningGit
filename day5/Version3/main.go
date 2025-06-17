package main

import (
	"fmt"
)

type PaymentMethod interface {
	Pay(amount float64) string
}

// Optional interface for OTP support
type OTPGenerator interface {
	GenerateOTP() string
}

// CreditCard implements both PaymentMethod and OTPGenerator
type CreditCard struct {
	CardNumber string
}

func (c CreditCard) Pay(amount float64) string {
	last4 := c.CardNumber[len(c.CardNumber)-4:]
	return fmt.Sprintf("[CreditCard] Paid ₹%.2f using card ending with %s", amount, last4)
}

func (c CreditCard) GenerateOTP() string {
	return "[CreditCard] OTP sent to registered number"
}

// PayPal implements only PaymentMethod
type PayPal struct {
	Email string
}

func (p PayPal) Pay(amount float64) string {
	return fmt.Sprintf("[PayPal] Paid ₹%.2f using PayPal account: %s", amount, p.Email)
}

// UPI implements both PaymentMethod and OTPGenerator
type UPI struct {
	UPIID string
}

func (u UPI) Pay(amount float64) string {
	return fmt.Sprintf("[UPI] Paid ₹%.2f using UPI: %s", amount, u.UPIID)
}

func (u UPI) GenerateOTP() string {
	return "[UPI] OTP sent to registered device"
}

func main() {
	payments := []PaymentMethod{
		CreditCard{CardNumber: "1234567812341234"},
		PayPal{Email: "user@example.com"},
		UPI{UPIID: "user@upi"},
	}

	for _, method := range payments {
		// Check if the method supports OTP
		if otpCapable, ok := method.(OTPGenerator); ok {
			fmt.Println(otpCapable.GenerateOTP())
		}
		fmt.Println(method.Pay(500.0))
		fmt.Println()
	}
}
