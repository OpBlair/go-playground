// Hardcoded dummy credentials for login testing
const dummy_credentials = {
    "student": { "email": "student@westridge.dev", "password": 12345 },
    "lecturer": { "email": "lecturer@westridge.dev", "password": 12345 },
    "admin": { "email": "admin@westridge.dev", "password": 12345 }
}

// Redirect paths for each user role
const roleRedirects = {
    "student": "student.html",
    "lecturer": "lecturer.html",
    "admin": "admin.html"
}

// Handles login form submission and authentication
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

// Opens a specific modal window
function openModal(modalId) {
    const modal = document.getElementById(modalId);
    if (modal) modal.classList.remove('hidden');
}

// Closes a specific modal window
function closeModal(modalId) {
    const modal = document.getElementById(modalId);
    if (modal) modal.classList.add('hidden');
}

// Toggles student fields on the signup/register form
function toggleStudentFields() {
    const roleSelect = document.getElementById('roleSelect');
    const studentFields = document.getElementById('studentFields');
    if (roleSelect && studentFields) {
        studentFields.style.display = (roleSelect.value === 'student') ? 'block' : 'none';
    }
}

// Opens the student course detail view
function openCourseView(code, title) {
    document.getElementById('dashboardView').classList.add('hidden');
    document.getElementById('courseDetailView').classList.remove('hidden');
    document.getElementById('courseDetailCode').innerText = code;
    document.getElementById('courseDetailTitle').innerText = title;
    window.scrollTo({ top: 0, behavior: 'smooth' });
}

// Closes student course detail view and returns to dashboard
function closeCourseView() {
    document.getElementById('courseDetailView').classList.add('hidden');
    document.getElementById('dashboardView').classList.remove('hidden');
    window.scrollTo({ top: 0, behavior: 'smooth' });
}

// Opens lecturer course management workspace and builds week dropdowns dynamically
function openCourseManagement(code, title, count, totalWeeks) {
    document.getElementById('dashboardView').classList.add('hidden');
    document.getElementById('courseWorkspaceView').classList.remove('hidden');
    document.getElementById('workspaceCode').innerText = code;
    document.getElementById('workspaceTitle').innerText = title;
    document.getElementById('workspaceStudentCount').innerText = count;
    document.getElementById('workspaceWeeksCount').innerText = totalWeeks;

    const outlineSelect = document.getElementById('outlineWeekSelect');
    const quizSelect = document.getElementById('quizWeekSelect');
    
    if (outlineSelect && quizSelect) {
        outlineSelect.innerHTML = '<option value="">Select Week...</option>';
        quizSelect.innerHTML = '<option value="">Select Week...</option>';

        for (let i = 1; i <= totalWeeks; i++) {
            const optionHtml = `<option value="week${i}">Week ${i}</option>`;
            outlineSelect.insertAdjacentHTML('beforeend', optionHtml);
            quizSelect.insertAdjacentHTML('beforeend', optionHtml);
        }
    }

    window.scrollTo({ top: 0, behavior: 'smooth' });
}

// Closes lecturer course workspace and returns to course list
function closeCourseManagement() {
    document.getElementById('courseWorkspaceView').classList.add('hidden');
    document.getElementById('dashboardView').classList.remove('hidden');
    window.scrollTo({ top: 0, behavior: 'smooth' });
}

// Toggles expandable course outline sections (accordions)
function toggleAccordion(id) {
    const el = document.getElementById(id);
    if (el) el.classList.toggle('hidden');
}

// Simulates frontend backend console logs for UI mockups
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

// Handles generic modal form submissions and closes modals
function handleModalSubmit(event, successMessage) {
    event.preventDefault();
    const consoleLog = document.getElementById('consoleLogText');
    if (consoleLog) {
        consoleLog.textContent = successMessage;
    }
            
    document.querySelectorAll('[id\$="Modal"]').forEach(el => el.classList.add('hidden'));
    alert(successMessage + "\n(Note: This is a non-functional frontend UI mockup connected to your Go backend logic)");
}
