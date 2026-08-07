package main

import (
	"encoding/json"
	"fmt"
	"mime"
	"net/http"
)

type PaymentRequest struct {
	PaymentMethod string  `json:"payment_method"`
	Amount        float64 `json:"amount"`
}

type PaymentResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type PaymentGateway interface {
	Pay(amount float64) error
}

type Stripe struct{}
type PayPal struct{}
type MoMo struct{}

func (s Stripe) Pay(amount float64) error {
	if amount <= 0 {
		return fmt.Errorf("Amount must be greater than zero")
	}
	if amount >= 1000 {
		return fmt.Errorf("The amount $%.2f exceeds daily transaction limit", amount)
	}
	fmt.Printf("Paid %.2f using Stripe\n", amount)
	return nil
}

func (p PayPal) Pay(amount float64) error {
	if amount <= 0 {
		return fmt.Errorf("Amount must be greater than zero")
	}
	if amount >= 1200 {
		return fmt.Errorf("The amount $%.2f exceeds daily transaction limit", amount)
	}
	fmt.Printf("Paid %.2f using Paypal\n", amount)
	return nil
}

func (m MoMo) Pay(amount float64) error {
	if amount <= 0 {
		return fmt.Errorf("Amount must be greater than zero")
	}
	if amount >= 3000 {
		return fmt.Errorf("The amount $%.2f exceeds daily transaction limit", amount)
	}
	fmt.Printf("Paid %.2f using Mobile Money\n", amount)
	return nil
}

func Checkout(g PaymentGateway, amount float64) error {
	err := g.Pay(amount)
	if err != nil {
		fmt.Println()
		fmt.Println("Error:", err)
	}

	return err
}

func PaymentHandler(w http.ResponseWriter, r *http.Request) {
	// w http.ResponseWriter := sends html, text, status codes back to the web browser
	// r *http.Request := pointer that contains all the incoming data sent by the browser

	// check the Method
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	/*
		// Extract Form Data
		methodInput := r.FormValue("payment_method")
		amountInput := r.FormValue("amount")

		// convert amount from string to float
		amount, err := strconv.ParseFloat(amountInput, 64)
		if err != nil {
			http.Error(w, "Invalid amount entered", http.StatusBadRequest)
			return
		}
	*/

	var req PaymentRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid JSON data package", http.StatusBadRequest)
		return
	}

	var gateway PaymentGateway
	methodInput := req.PaymentMethod
	amount := req.Amount

	switch methodInput {
	case "stripe":
		gateway = Stripe{}
	case "paypal":
		gateway = PayPal{}
	case "momo":
		gateway = MoMo{}
	default:
		http.Error(w, "Invalid Payment method select", http.StatusBadRequest)
		return
	}

	err = Checkout(gateway, amount)
	if err != nil {
		fmt.Fprintf(w, "An error occured: %v", err)
		return
	}
	fmt.Fprintf(w, "Payment of $%.2f processed successfully!", amount)
}

func main() {
	// MIME registry bug for JavaScript files
	mime.AddExtensionType(".js", "text/javascript")

	// Serve static files from the current directory "."
	fs := http.FileServer(http.Dir("."))
	http.Handle("/", fs)

	// send any incoming requests at "/checkout" to your PaymentHandler
	http.HandleFunc("/checkout", PaymentHandler)

	fmt.Println("Server is running on http://localhost:8080...")

	err := http.ListenAndServe(":8080", nil)

	if err != nil {
		fmt.Println("Server failed to start", err)
	}
}
