const dummy_credentials = {
    "student": {
        "email": "student@westridge.dev",
        "password": 12345
    },
    "staff": {
        "email": "staff@westridge.dev",
        "password": 12345
    },
    "admin": {
        "email": "admin@westridge.dev",
        "password": 12345
    }
}

// Map roles to specific destinations
const roleRedirects = {
    "student": "student.html",
    "staff": "staff.html",
    "admin": "admin.html"
}

// Listen for the form submission
document.getElementById("loginForm").addEventListener("submit", function(event) {
    // Stop the page from refreshing immediately
    event.preventDefault();

    // 2. Grab the inputs at the exact moment of click
    const inputEmail = document.getElementById("email").value;
    const inputPassword = document.getElementById("password").value;

    let authenticatedRole = null;

    // Check credentials
    for (const [role, credentials] of Object.entries(dummy_credentials)) {
        if (credentials.email === inputEmail && credentials.password == inputPassword) {
            authenticatedRole = role;
            break;
        }
    }

    // Handle success or failure
    if (authenticatedRole) {
        localStorage.setItem("userRole", authenticatedRole);
        localStorage.setItem("userEmail", inputEmail);

        // Redirect to the role's page
        window.location.href = roleRedirects[authenticatedRole];
    } else {
        alert("Invalid credentials!");
    }
});

function openModal(modalId) {
    document.getElementById(modalId).classList.remove('hidden');
}

function closeModal(modalId) {
    document.getElementById(modalId).classList.add('hidden');
}

function simulateAction(actionName) {
    const consoleLog = document.getElementById('consoleLogText');
    consoleLog.textContent = `Executed: ${actionName} [UI simulation mockup mode]`;
    consoleLog.classList.remove('text-emerald-400');
    consoleLog.classList.add('text-amber-400');
    setTimeout(() => {
        consoleLog.classList.remove('text-amber-400');
        consoleLog.classList.add('text-emerald-400');
    }, 1000);
}

function handleModalSubmit(event, successMessage) {
    event.preventDefault();
    const consoleLog = document.getElementById('consoleLogText');
    consoleLog.textContent = successMessage;
            
    // Close all modals
    document.querySelectorAll('[id$="Modal"]').forEach(el => el.classList.add('hidden'));
            
    // Show alert notice
    alert(successMessage + "\n(Note: This is a non-functional frontend UI mockup connected to your Go backend logic)");
}
