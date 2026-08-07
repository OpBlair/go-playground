'use strict';

const form = document.querySelector("form");

form.addEventListener("submit", async (e) => {
    e.preventDefault();
    await processPayment();
});

async function processPayment(){
    const methodSelect = document.getElementById("payment-method");
    const amountInput = document.getElementById("amount");

    const selectedMethod = methodSelect.value;
    const enteredAmount = parseFloat(amountInput.value); 

    const paymentPackage = {
        payment_method: selectedMethod,
        amount: enteredAmount
    }

    const response = await fetch('/checkout', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(paymentPackage)
    })

    const message = await response.text();
    alert(message);
}