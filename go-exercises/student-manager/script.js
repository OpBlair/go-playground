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

// Handle user Logout
function handleLogout() {
    localStorage.removeItem('westridge_token');
    localStorage.removeItem('westridge_role');

    window.location.replace('login.html');
}

// Opens a specific modal window
function openModal(modalId) {
    const modal = document.getElementById(modalId);
    if (modal) {
        modal.classList.remove('hidden');
        
        if (modalId === 'addUserModal') {
            loadProgramsIntoDropdown();
        } else if (modalId === 'addCourseModal') {
            loadCourseModalData();
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
            loadDashboardStats();
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
            loadDashboardStats();
        } else {
            alert("Failed to register user.");
        }

    } catch (error) {
        console.error("Connection error:", error);
        alert("Could not connect to the server. Is your Go backend running?");
    }
}

// Store loaded programs globally so we can reference their duration easily
let cachedPrograms = [];

// Load Programs and Lecturers when the Course Modal opens
async function loadCourseModalData() {
    const token = localStorage.getItem('westridge_token');
    if (!token) return;

    try {
        // Fetch programs list (make sure your Go backend returns duration!)
        const progResponse = await fetch("/api/programs-list", {
            headers: { "Authorization": `Bearer ${token}` }
        });
        if (progResponse.ok) {
            cachedPrograms = await progResponse.json();
            const progSelect = document.getElementById("courseProgramSelect");
            progSelect.innerHTML = '<option value="">Select a program...</option>';
            
            if (cachedPrograms && cachedPrograms.length > 0) {
                cachedPrograms.forEach(prog => {
                    const option = document.createElement("option");
                    option.value = prog.id;
                    option.textContent = `${prog.name} (${prog.code}) - [${prog.duration} Years]`;
                    progSelect.appendChild(option);
                });
            }
        }

        // Fetch lecturers list
        const lecResponse = await fetch("/api/lecturers-list", {
            headers: { "Authorization": `Bearer ${token}` }
        });
        if (lecResponse.ok) {
            const lecturers = await lecResponse.json();
            const lecSelect = document.getElementById("courseLecturerSelect");
            lecSelect.innerHTML = '<option value="">Assign a lecturer...</option>';
            
            if (lecturers && lecturers.length > 0) {
                lecturers.forEach(lec => {
                    const option = document.createElement("option");
                    option.value = lec.id; // UUID string
                    option.textContent = `${lec.first_name} ${lec.last_name}`;
                    lecSelect.appendChild(option);
                });
            }
        }
    } catch (error) {
        console.error("Error loading course modal data:", error);
    }
}

// Dynamically adjust Academic Year options based on selected program's duration
function handleProgramChangeForCourse() {
    const progSelect = document.getElementById("courseProgramSelect");
    const yearSelect = document.getElementById("courseAcademicYear");
    const selectedId = parseInt(progSelect.value, 10);

    yearSelect.innerHTML = '<option value="">Select academic year...</option>';
    if (!selectedId) return;

    // Find the selected program from our cache
    const program = cachedPrograms.find(p => p.id === selectedId);
    if (program && program.duration) {
        for (let i = 1; i <= program.duration; i++) {
            const option = document.createElement("option");
            option.value = i;
            option.textContent = `Year ${i}`;
            yearSelect.appendChild(option);
        }
    }
}

// Handle Course Creation Submission
async function handleCreateCourseSubmit(event) {
    event.preventDefault();

    const code = document.getElementById("courseCode").value;
    const title = document.getElementById("courseTitle").value;
    const programId = parseInt(document.getElementById("courseProgramSelect").value, 10);
    const academicYear = parseInt(document.getElementById("courseAcademicYear").value, 10);
    const semester = parseInt(document.getElementById("courseSemester").value, 10);
    const isElective = document.getElementById("courseIsElective").value === "true";
    
    // Lecturer ID can be empty/null if none selected
    const lecturerIdVal = document.getElementById("courseLecturerSelect").value;
    const lecturerId = lecturerIdVal ? lecturerIdVal : null;

    const token = localStorage.getItem('westridge_token');
    if (!token) {
        alert("Unauthorized: Please log in again.");
        window.location.href = 'index.html';
        return;
    }

    try {
        const response = await fetch("/api/courses", {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
                "Authorization": `Bearer ${token}`
            },
            body: JSON.stringify({
                code: code,
                title: title,
                program_id: programId,
                academic_year: academicYear,
                semester: semester,
                is_elective: isElective,
                lecturer_id: lecturerId
            })
        });

        if (!response.ok) {
            const errorText = await response.text();
            if (response.status === 409) {
                alert("Error: Course code already exists!");
            } else if (response.status === 403) {
                alert("Forbidden: Admin access required.");
            } else {
                alert("Failed to create course: " + errorText);
            }
            return;
        }

        const data = await response.json();
        if (data.success) {
            alert(data.message || "Course created successfully!");
            document.querySelectorAll('[id$="Modal"]').forEach(el => el.classList.add('hidden'));
            event.target.reset();
            // location.reload(); // Optional reload
            loadDashboardStats();
        } else {
            alert("Failed to create course.");
        }

    } catch (error) {
        console.error("Connection error:", error);
        alert("Could not connect to the server.");
    }
}

// Load Dashboard Stats
async function loadDashboardStats() {
    try {
        // get token
        const token = localStorage.getItem('westridge_token')

        // make request with a valid token since it's a protected route
        const response = await fetch("/api/stats", {
            method: "GET",
            headers: {
                "Authorization": `Bearer ${token}`
            }
        });

        if (!response.ok) {
            const errorText = await response.text();
            throw new Error(`Failed to fetch stats: ${response.status} - ${errorText}`);
        }
        
        const stats = await response.json();

        document.getElementById("totalUsersStat").innerHTML = stats.total_users;
        document.getElementById("totalAdminsStat").innerHTML = stats.total_admins;
        document.getElementById("totalLecturersStat").innerHTML = stats.total_lecturers;
        document.getElementById("activeStudentsStat").innerHTML = stats.total_students;
        document.getElementById("activeProgramsStat").innerHTML = stats.total_programs;
        document.getElementById("activeCoursesStat").innerHTML = stats.active_courses;

    } catch (error) {
        console.error('Error loading dashboard stats', error);
    }
}

// Program Frontend Preview 
async function loadProgramsTable() {
    const token = localStorage.getItem('westridge_token');
    if (!token) return;

    try {
        const response = await fetch("/api/programs-list", {
            method: "GET",
            headers: { "Authorization": `Bearer ${token}` }
        });

        if (!response.ok) throw new Error("Failed to fetch programs");

        const programs = await response.json();
        const tbody = document.getElementById("programsTableBody");
        tbody.innerHTML = "";

        if (!programs || programs.length === 0) {
            tbody.innerHTML = `<tr><td colspan="4" class="py-4 px-6 text-center text-slate-400">No programs found.</td></tr>`;
            return;
        }

        programs.forEach(prog => {
            const tr = document.createElement("tr");
            tr.className = "hover:bg-slate-50/50 transition";
            tr.innerHTML = `
                <td class="py-4 px-6 font-mono font-semibold text-slate-800">${prog.code}</td>
                <td class="py-4 px-6 font-medium text-slate-800">${prog.name}</td>
                <td class="py-4 px-6">${prog.duration} Years</td>
                <td class="py-4 px-6 text-right space-x-1">
                    <button onclick="deleteProgram(${prog.id})" class="p-2 text-slate-400 hover:text-red-600 hover:bg-red-50 rounded-lg transition"><i class="fa-solid fa-trash-can"></i></button>
                </td>
            `;
            tbody.appendChild(tr);
        });

    } catch (error) {
        console.error("Error loading programs table:", error);
    }
}

// Load Users Table
async function loadUsersTable() {
    const token = localStorage.getItem('westridge_token');
    if (!token) return;

    try {
        const response = await fetch("/api/users-list", {
            method: "GET",
            headers: { "Authorization": `Bearer ${token}` }
        });

        if (!response.ok) throw new Error("Failed to fetch users");

        const users = await response.json();
        const tbody = document.getElementById("usersTableBody");
        if (!tbody) return;
        tbody.innerHTML = "";

        if (!users || users.length === 0) {
            tbody.innerHTML = `<tr><td colspan="5" class="py-4 px-6 text-center text-slate-400">No users found.</td></tr>`;
            return;
        }

        users.forEach(user => {
            const tr = document.createElement("tr");
            tr.className = "hover:bg-slate-50/50 transition";
            
            // Badge color styling based on role
            let roleBadge = `<span class="px-2.5 py-1 rounded-lg text-xs font-semibold bg-slate-100 text-slate-700">${user.role}</span>`;
            if (user.role === 'admin') roleBadge = `<span class="px-2.5 py-1 rounded-lg text-xs font-semibold bg-rose-50 text-rose-700 border border-rose-200">Admin</span>`;
            if (user.role === 'lecturer') roleBadge = `<span class="px-2.5 py-1 rounded-lg text-xs font-semibold bg-indigo-50 text-indigo-700 border border-indigo-200">Lecturer</span>`;
            if (user.role === 'student') roleBadge = `<span class="px-2.5 py-1 rounded-lg text-xs font-semibold bg-emerald-50 text-emerald-700 border border-emerald-200">Student</span>`;

            tr.innerHTML = `
                <td class="py-4 px-6 font-medium text-slate-800">${user.first_name} ${user.last_name}</td>
                <td class="py-4 px-6 text-slate-600">${user.email}</td>
                <td class="py-4 px-6">${roleBadge}</td>
                <td class="py-4 px-6 text-slate-600">${user.assigned_program}</td>
                <td class="py-4 px-6 text-right space-x-1">
                    <button class="p-2 text-slate-400 hover:text-red-600 hover:bg-red-50 rounded-lg transition"><i class="fa-solid fa-trash-can"></i></button>
                </td>
            `;
            tbody.appendChild(tr);
        });

    } catch (error) {
        console.error("Error loading users table:", error);
    }
}

// Load Courses Table
async function loadCoursesTable() {
    const token = localStorage.getItem('westridge_token');
    if (!token) return;

    try {
        const response = await fetch("/api/courses-list", {
            method: "GET",
            headers: { "Authorization": `Bearer ${token}` }
        });

        if (!response.ok) throw new Error("Failed to fetch courses");

        const courses = await response.json();
        const tbody = document.getElementById("coursesTableBody");
        if (!tbody) return;
        tbody.innerHTML = "";

        if (!courses || courses.length === 0) {
            tbody.innerHTML = `<tr><td colspan="5" class="py-4 px-6 text-center text-slate-400">No courses found.</td></tr>`;
            return;
        }

        courses.forEach(course => {
            const tr = document.createElement("tr");
            tr.className = "hover:bg-slate-50/50 transition";
            
            const typeBadge = course.is_elective 
                ? `<span class="px-2.5 py-1 rounded-lg text-xs font-semibold bg-amber-50 text-amber-700 border border-amber-200">Elective</span>`
                : `<span class="px-2.5 py-1 rounded-lg text-xs font-semibold bg-emerald-50 text-emerald-700 border border-emerald-200">Core</span>`;

            tr.innerHTML = `
                <td class="py-4 px-6 font-mono font-semibold text-slate-800">${course.code}</td>
                <td class="py-4 px-6 font-medium text-slate-800">${course.title}</td>
                <td class="py-4 px-6">${typeBadge}</td>
                <td class="py-4 px-6 font-mono text-xs text-slate-600">Year ${course.academic_year}, Sem ${course.semester}</td>
                <td class="py-4 px-6 text-right space-x-1">
                    <button class="p-2 text-slate-400 hover:text-red-600 hover:bg-red-50 rounded-lg transition"><i class="fa-solid fa-trash-can"></i></button>
                </td>
            `;
            tbody.appendChild(tr);
        });

    } catch (error) {
        console.error("Error loading courses table:", error);
    }
}

// Switches dashboard table views dynamically based on dropdown selection
function filterDashboardViews(viewType) {
    const programsDiv = document.getElementById('programsSection');
    const usersDiv = document.getElementById('usersSection');
    const coursesDiv = document.getElementById('coursesSection');

    if (!programsDiv || !usersDiv || !coursesDiv) return;

    // Reset visibility by removing hidden class from all
    programsDiv.classList.remove('hidden');
    usersDiv.classList.remove('hidden');
    coursesDiv.classList.remove('hidden');

    // Filter based on selection
    if (viewType === 'programs') {
        usersDiv.classList.add('hidden');
        coursesDiv.classList.add('hidden');
    } else if (viewType === 'users') {
        programsDiv.classList.add('hidden');
        coursesDiv.classList.add('hidden');
    } else if (viewType === 'courses') {
        programsDiv.classList.add('hidden');
        usersDiv.classList.add('hidden');
    }
    // If 'all', all divs remain visible naturally in their separate cards
}

document.addEventListener('DOMContentLoaded', () => {
    const token = localStorage.getItem('westridge_token');
    const role = localStorage.getItem('westridge_role');

    // If there's no token or they aren't an admin on an admin page, bounce them
    if (!token || role !== 'admin') {
        window.location.replace('login.html');
        return;
    }
    
    loadDashboardStats();
    loadProgramsTable();
    loadUsersTable();
    loadCoursesTable();
});