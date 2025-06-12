package main

import (
	"fmt"
	"os"
	"strconv"
)

var currencyStore = [...]string{
	"USD",
	"INR",
	"JPY",
	"EUR",
}

var currencyMap = map[string]float64{
	"USDINR": 85.00,
	"USDUSD": 1,
	"USDJPY": 145.00,
	"USDEUR": 0.90,
	"JPYUSD": 0.0069,
	"JPYJPY": 1,
	"JPYINR": 0.58,
	"JPYEUR": 0.0062,
	"INRUSD": 0.012,
	"INRJPY": 1.72,
	"INRINR": 1,
	"INREUR": 0.011,
	"EURUSD": 1.11,
	"EURJPY": 161.00,
	"EURINR": 91.00,
	"EUREUR": 1,
}

var amount float64

func isValidCurrency(curr string) bool {
	return curr == "USD" || curr == "INR" || curr == "JPY" || curr == "EUR"
}

func validateInput(arr []string) bool {

	if len(arr) == 1 && arr[0] == "--list" {
		return true
	}

	if len(arr) != 3 {
		return false
	}
	value, err := strconv.ParseFloat(arr[0], 64)
	if err != nil {
		fmt.Println("Invalid amount:", arr[0])
		return false
	}
	amount = value
	if !isValidCurrency(arr[1]) || !isValidCurrency(arr[2]) {
		fmt.Println("Invalid currency code.")
		return false
	}

	return true
}

func currencyConverter(amount float64, initialCurrency, finalCurrency string) float64 {
	return amount * currencyMap[initialCurrency+finalCurrency]
}

func main() {

	argsProg := os.Args
	argsWithoutPath := os.Args[1:]

	fmt.Println(argsProg)
	fmt.Println(argsWithoutPath)

	if validateInput(argsWithoutPath) {

		if argsWithoutPath[0] == "--list" {
			for i := 0; i < len(currencyStore); i++ {
				for j := 0; j < len(currencyStore); j++ {
					index := currencyStore[i] + currencyStore[j]
					value := currencyMap[index]

					fmt.Printf("The exchange rate for converting %v to %v is: %v \n", currencyStore[i], currencyStore[j], value)
				}
			}
		}

		convertedAmount := currencyConverter(amount, argsWithoutPath[1], argsWithoutPath[2])
		fmt.Println(amount, argsWithoutPath[1], "is equal to", convertedAmount)

	} else {
		fmt.Println("The input needs to be in the following format: go run fileName Amount IntialCurrency ConvertedCurrency")
	}

}
