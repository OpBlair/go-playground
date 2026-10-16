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

// Listen for the login form submission if it exists on the page
const loginForm = document.getElementById("loginForm");
if (loginForm) {
    loginForm.addEventListener("submit", function(event) {
        event.preventDefault();

        const inputEmail = document.getElementById("email").value;
        const inputPassword = document.getElementById("password").value;

        let authenticatedRole = null;

        for (const [role, credentials] of Object.entries(dummy_credentials)) {
            if (credentials.email === inputEmail && credentials.password == inputPassword) {
                authenticatedRole = role;
                break;
            }
        }

        if (authenticatedRole) {
            localStorage.setItem("userRole", authenticatedRole);
            localStorage.setItem("userEmail", inputEmail);
            window.location.href = roleRedirects[authenticatedRole];
        } else {
            alert("Invalid credentials!");
        }
    });
}

function openModal(modalId) {
    const modal = document.getElementById(modalId);
    if (modal) modal.classList.remove('hidden');
}

function closeModal(modalId) {
    const modal = document.getElementById(modalId);
    if (modal) modal.classList.add('hidden');
}

function toggleStudentFields() {
    const roleSelect = document.getElementById('roleSelect');
    const studentFields = document.getElementById('studentFields');
    if (roleSelect && studentFields) {
        if (roleSelect.value === 'student') {
            studentFields.style.display = 'block';
        } else {
            studentFields.style.display = 'none';
        }
    }
}

function simulateAction(actionName) {
    const consoleLog = document.getElementById('consoleLogText');
    if (consoleLog) {
        consoleLog.textContent = `Executed: ${actionName} [UI simulation mockup mode]`;
        consoleLog.classList.remove('text-emerald-400');
        consoleLog.classList.add('text-amber-400');
        setTimeout(() => {
            consoleLog.classList.remove('text-amber-400');
            consoleLog.classList.add('text-emerald-400');
        }, 1000);
    }
}

function handleModalSubmit(event, successMessage) {
    event.preventDefault();
    const consoleLog = document.getElementById('consoleLogText');
    if (consoleLog) {
        consoleLog.textContent = successMessage;
    }
            
    // Close all modals
    document.querySelectorAll('[id$="Modal"]').forEach(el => el.classList.add('hidden'));
            
    alert(successMessage + "\n(Note: This is a non-functional frontend UI mockup connected to your Go backend logic)");
}
