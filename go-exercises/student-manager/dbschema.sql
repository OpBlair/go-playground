-- =========================================================================
-- Westridge Student Management System - Master Database Schema
-- =========================================================================

-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Strict ENUM type for roles
CREATE TYPE user_role AS ENUM ('admin', 'lecturer', 'student');

-- Programs Table
CREATE TABLE programs (
    id SERIAL PRIMARY KEY,
    code VARCHAR(20) UNIQUE NOT NULL,
    name VARCHAR(100) NOT NULL,
    duration INT NOT NULL
);

-- Users Table (Includes composite unique constraint for role-based foreign keys)
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    password_hash VARCHAR(255) NOT NULL,
    first_name VARCHAR(30) NOT NULL,
    last_name VARCHAR(30) NOT NULL,
    email VARCHAR(50) NOT NULL,
    role user_role NOT NULL,
    program_id INT REFERENCES programs(id) ON DELETE SET NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    CONSTRAINT unique_user_id_role UNIQUE (id, role)
);

-- Courses Table (Strictly enforces that assigned lecturer has the 'lecturer' role)
CREATE TABLE courses (
    id SERIAL PRIMARY KEY,
    code VARCHAR(10) UNIQUE NOT NULL, 
    title VARCHAR(150) NOT NULL,
    program_id INT REFERENCES programs(id) ON DELETE CASCADE,
    is_elective BOOLEAN DEFAULT FALSE, -- FALSE = Core, TRUE = Elective
    academic_year INT NOT NULL CHECK (academic_year > 0), 
    semester INT NOT NULL CHECK (semester IN (1, 2)),      
    total_weeks INT NOT NULL DEFAULT 12, 
    lecturer_id UUID,
    lecturer_role user_role DEFAULT 'lecturer' CHECK (lecturer_role = 'lecturer'),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    FOREIGN KEY (lecturer_id, lecturer_role) REFERENCES users(id, role) ON DELETE SET NULL
);

-- Course Modules / Weekly Outlines Table (Published by Lecturers)
CREATE TABLE course_modules (
    id SERIAL PRIMARY KEY,
    course_id INT REFERENCES courses(id) ON DELETE CASCADE,
    week_number INT NOT NULL CHECK (week_number > 0),
    topic_title VARCHAR(200) NOT NULL,
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(course_id, week_number)
);

-- Enrollments Table (Many-to-Many: Students to Courses)
CREATE TABLE enrollments (
    student_id UUID REFERENCES users(id) ON DELETE CASCADE,
    course_id INT REFERENCES courses(id) ON DELETE CASCADE,
    enrolled_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (student_id, course_id)
);

-- Quizzes Table (Tied to specific target weeks and restricted to lecturers)
CREATE TABLE quizzes (
    id SERIAL PRIMARY KEY,
    title VARCHAR(200) NOT NULL,
    course_id INT REFERENCES courses(id) ON DELETE CASCADE,
    week_number INT NOT NULL CHECK (week_number > 0), 
    lecturer_id UUID,
    lecturer_role user_role DEFAULT 'lecturer' CHECK (lecturer_role = 'lecturer'),
    questions JSONB NOT NULL, 
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    FOREIGN KEY (lecturer_id, lecturer_role) REFERENCES users(id, role) ON DELETE SET NULL
);

-- Grades Table
CREATE TABLE grades (
    id SERIAL PRIMARY KEY,
    student_id UUID REFERENCES users(id) ON DELETE CASCADE,
    quiz_id INT REFERENCES quizzes(id) ON DELETE CASCADE,
    score NUMERIC(5, 2) NOT NULL CHECK (score >= 0 AND score <= 100),
    submitted_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Notifications Tables (Multi-user broadcast)
-- Core Notifications Table (Stores the message content once)
CREATE TABLE notifications (
    id SERIAL PRIMARY KEY,
    title VARCHAR(150) NOT NULL,
    message TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Junction Table: Maps notifications to multiple users & tracks read status per user
CREATE TABLE notification_recipients (
    notification_id INT REFERENCES notifications(id) ON DELETE CASCADE,
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    is_read BOOLEAN DEFAULT FALSE,
    read_at TIMESTAMP WITH TIME ZONE,
    PRIMARY KEY (notification_id, user_id)
);

-- Audit Trail Log
CREATE TABLE IF NOT EXISTS audit_logs (
    id SERIAL PRIMARY KEY,
    actor_id UUID REFERENCES users(id) ON DELETE SET NULL,
    action VARCHAR(100) NOT NULL,
    details TEXT,
    ip_address VARCHAR(45),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
