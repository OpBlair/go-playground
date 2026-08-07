# 🛒 JSON-Driven Interface-Based Payment System

A lightweight, modern web payment processing simulation built in Go and JavaScript. This project demonstrates object grouping, decoupling via interfaces, and modern asynchronous JSON data exchanges between an HTML frontend and a Go backend server.

---

## 🛠 Architecture Overview

The system follows a strict decoupling pattern. The web handler does not know the internal mechanics of individual payment methods; it strictly accepts data bundles and processes transactions through a generic contract.

### 1. The Core Contract (Interfaces)
Instead of forcing the application to adapt to specific payment tools, we create an engine interface contract:
```go
type PaymentGateway interface {
    Pay(amount float64) error
}
```
Any structural entity that implements a `.Pay(amount float64) error` method automatically satisfies this contract.

### 2. Concrete Adapters (The Structs)
We created three unique transaction endpoints that implement the interface contract with customized business validation rules:
*   **Stripe:** Rejects transactions equal to or exceeding `$1,000.00`.
*   **PayPal:** Rejects transactions equal to or exceeding `$1,200.00`.
*   **Mobile Money (MoMo):** Rejects transactions equal to or exceeding `$3,000.00`.

---

## 🔄 Data Lifecycle & Flow Map

```text
[ Frontend: index.html ] 
       │  (User selects gateway & enters amount)
       ▼
[ Frontend: script.js ] 
       │  (Extracts DOM elements via IDs)
       │  (Packages variables into a single JavaScript Object)
       │  (Fires async fetch Request using JSON.stringify)
       ▼
[ Backend: Go HTTP Server /checkout ]
       │  (Guards against non-POST requests)
       │  (Decodes incoming JSON stream directly into PaymentRequest struct)
       │  (Evaluates choice using a clean switch router)
       │  (Assigns structural reference to PaymentGateway interface slot)
       ▼
[ Backend: Checkout Engine ]
       │  (Executes .Pay() method on the active interface asset)
       ▼
[ Response Feedback ]
       └─► Success: Prints success confirmation to browser string stream.
       └─► Failure: Gracefully passes validation error message back up.
```

---

## 💻 Tech Stack Implementation Details

### The Frontend Grouping & Async Fetch
Instead of transmitting archaic, raw form text streams, the user interface packages its parameters cleanly using an asynchronous block.

*   **Extraction:** Individual DOM assets are captured explicitly via unique node identities.
*   **Object Modeling:** Data values are converted to proper logical types (e.g., parsing an amount text into an executable browser float) prior to serialization.
*   **Async Network Exchange:** Modern `async/await` syntax ensures non-blocking browser threads. The data bundle payload is flagged with a strict `application/json` Content-Type header.

### The Backend Streaming Decoder
The Go endpoint completely bypasses loose string parsing routines by implementing a target data architecture mapping layer.

*   **Data Transport Object:** The `PaymentRequest` structure defines explicit structural tags (`json:"field"`) allowing automated parsing execution.
*   **Stream Unmarshalling:** We use Go's dynamic `json.NewDecoder(r.Body).Decode(&req)` pipeline. This reads the incoming payload directly from the network buffer and assigns fields to our data structure in a single operation.
*   **Linux Environment Patch:** To ensure stability on systems like Fedora, we explicitly bind the application's local MIME registry extension parameters (`mime.AddExtensionType(".js", "text/javascript")`) so the client browser never intercepts or blocks internal execution components.

---

## 🚀 How to Run Locally

### 1. File Structure Setup
Ensure your local project directory looks exactly like this:
```text
payment-system/
├── index.html
├── script.js
└── main.go
```

### 2. Booting the Application
Open a terminal in your project directory and execute the unified file routine:
```bash
# Clear any background zombie processes locking port 8080 (Linux specific)
fuser -k 8080/tcp

# Run your Go server
go run main.go
```

### 3. Accessing the System
Open your modern web browser and navigate directly to:
👉 **`http://localhost:8080`**

*(Note: Do not double-click the `index.html` file to open it locally as a file URL; accessing it directly via the local server origin is necessary to bypass browser CORS security checks).*
