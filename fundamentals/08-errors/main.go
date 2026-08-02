package main

import (
	"errors"
	"fmt"
)

func CheckAge(age int) error {
	if age < 18 {
		return errors.New("must be atleast 18 years old")
	}
	return nil
}

var ErrDatabaseNotFound = errors.New("database connection lost")

func QueryUser() error {
	return ErrDatabaseNotFound
}

func GetUserProfile() error {
	err := QueryUser()
	if err != nil {
		// %w wraps the error, allowing callers to inspect it using errors.Is
		return fmt.Errorf("failed to load user profile: %w", err)
	}
	return nil
}

var (
	ErrTooYoung     = errors.New("must be at least 18 years old")
	ErrBannedUser   = errors.New("this account has been revoked")
	ErrMissingEmail = errors.New("email address cannot be empty")
)

func RegisterUser(age int, isBanned bool, email string) error {
	if age < 18 {
		return ErrTooYoung
	}

	if isBanned {
		return ErrBannedUser
	}

	if email == "" {
		return ErrMissingEmail
	}
	return nil
}

// Custom Error Types & errors.As
// custom error struct with metadata fields
type ValidationError struct {
	Field   string
	Message string
}

// implement the error interface by adding the Error() method
func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation failed on field '%s': %s", e.Field, e.Message)
}

func ValidateUser(username string) error {
	if len(username) < 3 {
		// return a pointer to our custom error type
		return &ValidationError{
			Field:   "username",
			Message: "must be at least 3 characters long",
		}
	}
	return nil
}

func main() {
	fmt.Println("===== Go Errors Tutorial =====")
	// nil means no error
	err := CheckAge(17)

	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Access Granted")
	}

	err = GetUserProfile()
	if err != nil {
		fmt.Println("Error:", err)
		// Check if the wrapped error chain contains our sentinel error
		if errors.Is(err, ErrDatabaseNotFound) {
			fmt.Println("-> Action: Attempting to reconnect to the database...")
		}
	} else {
		fmt.Println("Profile loaded successfully")
	}

	err = RegisterUser(15, false, "test@example.com")
	if err != nil {
		fmt.Println("Registration failed:", err)
	}

	err = RegisterUser(20, true, "test2@example.com")
	if err != nil {
		fmt.Println("Registration failed", err)
	}

	err = RegisterUser(25, false, "")
	if err != nil {
		fmt.Println("Registration failed:", err)
	}

	fmt.Println("===== Advanced Go Errors: Custom Types & errors.As =====")

	err = ValidateUser("bo")
	if err != nil {
		fmt.Println("Standard Error Output:", err)

		// using errors.As to extract the custom error type and accessing its fields
		var valErr *ValidationError
		if errors.As(err, &valErr) {
			fmt.Printf("-> Extracted Metadata: Field [%s] caused the issue.\n", valErr.Field)
		}
	}

	fmt.Printf("\n ===== Errors Tutorial Exercise ===== \n")
	newAcc := Account{}
	newAcc.Balance = 1100

	err = newAcc.Withdraw(1200)
	if err != nil {
		fmt.Println(err)
	}

	// detect the specific error.
	if err != nil {
		if errors.Is(err, ErrInsufficientFunds) {
			fmt.Println("you don't have enough money.")
		} else {
			fmt.Println(err)
		}
	}

	err = newAcc.Withdraw(500)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("Withdrawal successful. New Balance: %.2f\n", newAcc.Balance)
	}
}

// Errors Exercise
type Account struct {
	Balance float64
}

var ErrInsufficientFunds = errors.New("insufficient funds")

func (a *Account) Withdraw(amount float64) error {
	if amount > a.Balance {
		return fmt.Errorf("withdrawal failed: %w", ErrInsufficientFunds)
	}
	if amount <= 0 {
		return errors.New("amount must be greater than zero")
	}
	a.Balance -= amount
	return nil
}
