-- Enable UUID extension (great for secure user IDs in URLs)
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

-- Users Table
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    password_hash VARCHAR(255) NOT NULL,
    first_name VARCHAR(30) NOT NULL,
    last_name VARCHAR(30) NOT NULL,
    email VARCHAR(50) NOT NULL,
    role user_role NOT NULL,
    program_id INT REFERENCES programs(id) ON DELETE SET NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Courses Table (Added total_weeks for dynamic week generation)
CREATE TABLE courses (
    id SERIAL PRIMARY KEY,
    code VARCHAR(10) UNIQUE NOT NULL, 
    title VARCHAR(150) NOT NULL,
    program_id INT REFERENCES programs(id) ON DELETE CASCADE,
    is_elective BOOLEAN DEFAULT FALSE, -- FALSE = Core, TRUE = Elective
    academic_year INT NOT NULL CHECK (academic_year > 0), 
    semester INT NOT NULL CHECK (semester IN (1, 2)),      
    total_weeks INT NOT NULL DEFAULT 12, -- Added: Total weeks for the course outline
    lecturer_id UUID REFERENCES users(id) ON DELETE SET NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
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

-- Quizzes Table (Added week_number to tie quizzes to specific weeks)
CREATE TABLE quizzes (
    id SERIAL PRIMARY KEY,
    title VARCHAR(200) NOT NULL,
    course_id INT REFERENCES courses(id) ON DELETE CASCADE,
    week_number INT NOT NULL CHECK (week_number > 0), -- Added: Target week for the quiz
    lecturer_id UUID REFERENCES users(id) ON DELETE SET NULL,
    questions JSONB NOT NULL, 
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Grades Table
CREATE TABLE grades (
    id SERIAL PRIMARY KEY,
    student_id UUID REFERENCES users(id) ON DELETE CASCADE,
    quiz_id INT REFERENCES quizzes(id) ON DELETE CASCADE,
    score NUMERIC(5, 2) NOT NULL CHECK (score >= 0 AND score <= 100),
    submitted_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
