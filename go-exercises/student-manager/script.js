// Redirect paths for each user role
const roleRedirects = {
    "student": "student.html",
    "lecturer": "lecturer.html",
    "admin": "admin.html"
}

// Handles login form submission and authentication
const loginForm = document.getElementById("loginForm");
if (loginForm) {
    loginForm.addEventListener("submit", async function(event) {
        event.preventDefault();

        const inputEmail = document.getElementById("email").value;
        const inputPassword = document.getElementById("password").value;

        try {
            // Send credentials to  Go backend
            const response = await fetch("/api/login", {
                method: "POST",
                headers: {
                    "Content-Type": "application/json"
                },
                body: JSON.stringify({
                    email: inputEmail,
                    password: inputPassword
                })
            });

            if (!response.ok) {
                // If Go returned 401 Unauthorized or 400 Bad Request
                alert("Invalid email or password!");
                return;
            }

            // Parse the JSON response sent from Go
            const data = await response.json();

            if (data.success) {
                localStorage.setItem('westridge_token', data.token);
                localStorage.setItem('westridge_role', data.role);
                
                // Redirect based on role
                if (data.role === 'admin') window.location.href = 'admin.html';
                else if (data.role === 'lecturer') window.location.href = 'lecturer.html';
                else window.location.href = 'student.html';
            } else {
                alert("Login failed. Please try again.");
            }

        } catch (error) {
            console.error("Connection error:", error);
            alert("Could not connect to the server. Is your Go backend running?");
        }
    });
}

// Opens a specific modal window
function openModal(modalId) {
    const modal = document.getElementById(modalId);
    if (modal) {
        modal.classList.remove('hidden');
        
        // If opening the user modal, fetch the latest programs from DB right away!
        if (modalId === 'addUserModal') {
            loadProgramsIntoDropdown();
        }
    }
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

// Handles creating a new program from the admin modal
async function handleCreateProgramSubmit(event) {
    event.preventDefault();

    // Grab input values from your modal form
    const programCode = document.getElementById("programCode").value;
    const programName = document.getElementById("programName").value;
    const duration = parseInt(document.getElementById("programDuration").value, 10);

    // Retrieve the JWT token stored during login
    const token = localStorage.getItem('westridge_token');
    if (!token) {
        alert("Unauthorized: Please log in again.");
        window.location.href = 'login.html';
        return;
    }

    try {
        // Send POST request to Go backend with Bearer token authentication
        const response = await fetch("/api/programs", {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
                "Authorization": `Bearer ${token}`
            },
            body: JSON.stringify({
                program_code: programCode,
                program_name: programName,
                duration: duration
            })
        });

        if (!response.ok) {
            const errorText = await response.text();
            if (response.status === 409) {
                alert("Error: Program code already exists!");
            } else if (response.status === 403) {
                alert("Forbidden: Admin access required.");
            } else {
                alert("Failed to create program: " + errorText);
            }
            return;
        }

        const data = await response.json();

        if (data.success) {
            alert(data.message || "Program created successfully!");
            
            // Close the modal
            document.querySelectorAll('[id$="Modal"]').forEach(el => el.classList.add('hidden'));
            
            // Reset form fields or reload the page/table to show the new program
            event.target.reset();
            // location.reload(); 
        } else {
            alert("Failed to create program.");
        }

    } catch (error) {
        console.error("Connection error:", error);
        alert("Could not connect to the server. Is your Go backend running?");
    }
}

// Fetches programs from the database and populates the user registration dropdown
async function loadProgramsIntoDropdown() {
    const programSelect = document.getElementById("userProgramSelect");
    if (!programSelect) return;

    const token = localStorage.getItem('westridge_token');
    if (!token) return;

    try {
        const response = await fetch("/api/programs-list", {
            method: "GET",
            headers: {
                "Authorization": `Bearer ${token}`
            }
        });

        if (!response.ok) {
            console.error("Failed to load programs list");
            return;
        }

        const programs = await response.json();

        // Clear existing options except the default prompt
        programSelect.innerHTML = '<option value="">Select a program...</option>';

        // Loop through database results and create options dynamically
        if (programs && programs.length > 0) {
            programs.forEach(prog => {
                const option = document.createElement("option");
                option.value = prog.id; // Database primary key (INT)
                option.textContent = `${prog.name} (${prog.code})`;
                programSelect.appendChild(option);
            });
        } else {
            const option = document.createElement("option");
            option.value = "";
            option.textContent = "No programs available (Create one first)";
            programSelect.appendChild(option);
        }

    } catch (error) {
        console.error("Error fetching programs:", error);
    }
}

// Handles registering a new user from the admin modal
async function handleCreateUserSubmit(event) {
    event.preventDefault();

    // Grab input values from your Register User modal form
    const firstName = document.getElementById("userFirstName").value;
    const lastName = document.getElementById("userLastName").value;
    const email = document.getElementById("userEmail").value;
    const password = document.getElementById("userPassword").value;
    const role = document.getElementById("roleSelect").value;
    
    // Program ID is only relevant for students or if selected, let's parse it safely
    const programSelect = document.getElementById("userProgramSelect");
    let programId = null;
    if (programSelect && programSelect.value) {
        programId = parseInt(programSelect.value, 10);
    }

    // Retrieve the JWT token stored during login
    const token = localStorage.getItem('westridge_token');
    if (!token) {
        alert("Unauthorized: Please log in again.");
        window.location.href = 'login.html';
        return;
    }

    try {
        // Send POST request to Go backend with Bearer token authentication
        const response = await fetch("/api/users", {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
                "Authorization": `Bearer ${token}`
            },
            body: JSON.stringify({
                first_name: firstName,
                last_name: lastName,
                user_email: email,
                user_password: password,
                user_role: role,
                program_id: programId
            })
        });

        if (!response.ok) {
            const errorText = await response.text();
            if (response.status === 409) {
                alert("Error: Email address already exists!");
            } else if (response.status === 403) {
                alert("Forbidden: Admin access required.");
            } else {
                alert("Failed to create user: " + errorText);
            }
            return;
        }

        const data = await response.json();

        if (data.success) {
            alert(data.message || "User registered successfully!");
            
            // Close the modal
            document.querySelectorAll('[id$="Modal"]').forEach(el => el.classList.add('hidden'));
            
            // Reset form fields
            event.target.reset();
            // location.reload(); // Uncomment to refresh the user list
        } else {
            alert("Failed to register user.");
        }

    } catch (error) {
        console.error("Connection error:", error);
        alert("Could not connect to the server. Is your Go backend running?");
    }
}