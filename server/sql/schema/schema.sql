DROP SCHEMA IF EXISTS library CASCADE;

CREATE SCHEMA library;

SET
    search_path TO library;

-- Добавление модуля для UUID
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Создание типов данных
CREATE TYPE user_role AS ENUM ('administrator', 'librarian', 'reader');

CREATE TYPE book_status AS ENUM (
    'available',
    'issued',
    'reserved',
    'lost',
    'damaged'
);

CREATE TYPE action_type AS ENUM (
    'login',
    'logout',
    'book_issue',
    'book_return',
    'book_extend',
    'book_add',
    'book_remove',
    'reader_register',
    'reader_update',
    'fine_payment'
);

-- 1. Таблица пользователей системы (администраторы, библиотекари, читатели)
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4 (),
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    role user_role NOT NULL DEFAULT 'reader',
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 2. Таблица читальных залов
CREATE TABLE reading_halls (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4 (),
    library_name VARCHAR(100) NOT NULL,
    hall_name VARCHAR(100) NOT NULL,
    specialization VARCHAR(200),
    total_seats INTEGER NOT NULL CHECK (total_seats > 0),
    occupied_seats INTEGER DEFAULT 0 CHECK (occupied_seats >= 0),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT chk_occupied_seats CHECK (occupied_seats <= total_seats)
);

-- 3. Таблица читателей (профили пользователей с ролью reader)
CREATE TABLE readers (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4 (),
    user_id UUID UNIQUE REFERENCES users (id) ON DELETE CASCADE,
    ticket_number VARCHAR(20) UNIQUE NOT NULL,
    full_name VARCHAR(200) NOT NULL,
    birth_date DATE NOT NULL,
    phone VARCHAR(20),
    education VARCHAR(100),
    reading_hall_id UUID REFERENCES reading_halls (id),
    registration_date DATE DEFAULT CURRENT_DATE,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 4. Таблица авторов
CREATE TABLE authors (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4 (),
    full_name VARCHAR(200) NOT NULL,
    birth_year INTEGER,
    death_year INTEGER,
    biography TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT chk_birth_death_year CHECK (
        death_year IS NULL
        OR death_year >= birth_year
    )
);

-- 5. Таблица книг
CREATE TABLE books (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4 (),
    title VARCHAR(500) NOT NULL,
    isbn VARCHAR(17),
    publication_year INTEGER,
    publisher VARCHAR(200),
    pages INTEGER,
    language VARCHAR(50) DEFAULT 'Russian',
    description TEXT,
    total_copies INTEGER NOT NULL DEFAULT 1 CHECK (total_copies > 0),
    available_copies INTEGER NOT NULL DEFAULT 1,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT chk_available_copies CHECK (
        available_copies >= 0
        AND available_copies <= total_copies
    )
);

-- 6. Связующая таблица авторов и книг (многие ко многим)
CREATE TABLE book_authors (
    book_id UUID REFERENCES books (id) ON DELETE CASCADE,
    author_id UUID REFERENCES authors (id) ON DELETE CASCADE,
    PRIMARY KEY (book_id, author_id)
);

-- 7. Таблица экземпляров книг
CREATE TABLE book_copies (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4 (),
    book_id UUID NOT NULL REFERENCES books (id) ON DELETE CASCADE,
    copy_code VARCHAR(50) UNIQUE NOT NULL, -- Шифр книги
    status book_status DEFAULT 'available',
    reading_hall_id UUID REFERENCES reading_halls (id),
    condition_notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 8. Таблица выдач книг
CREATE TABLE book_issues (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4 (),
    reader_id UUID NOT NULL REFERENCES readers (id),
    book_copy_id UUID NOT NULL REFERENCES book_copies (id),
    issue_date DATE DEFAULT CURRENT_DATE,
    due_date DATE NOT NULL,
    return_date DATE,
    extended_count INTEGER DEFAULT 0,
    librarian_id UUID REFERENCES users (id),
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT chk_dates CHECK (
        due_date >= issue_date
        AND (
            return_date IS NULL
            OR return_date >= issue_date
        )
    )
);

-- 9. Таблица штрафов
CREATE TABLE fines (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4 (),
    reader_id UUID NOT NULL REFERENCES readers (id),
    book_issue_id UUID REFERENCES book_issues (id),
    amount DECIMAL(10, 2) NOT NULL CHECK (amount >= 0),
    reason VARCHAR(500) NOT NULL,
    fine_date DATE DEFAULT CURRENT_DATE,
    paid_date DATE,
    paid_amount DECIMAL(10, 2) DEFAULT 0 CHECK (paid_amount >= 0),
    is_paid BOOLEAN DEFAULT FALSE,
    librarian_id UUID REFERENCES users (id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 10. Таблица рейтингов книг
CREATE TABLE book_ratings (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4 (),
    book_id UUID NOT NULL REFERENCES books (id) ON DELETE CASCADE,
    reader_id UUID NOT NULL REFERENCES readers (id),
    rating INTEGER NOT NULL CHECK (
        rating >= 1
        AND rating <= 5
    ),
    review TEXT,
    rating_date DATE DEFAULT CURRENT_DATE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (book_id, reader_id)
);

-- 11. Таблица логов системы
CREATE TABLE system_logs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4 (),
    user_id UUID REFERENCES users (id),
    action_type action_type NOT NULL,
    entity_type VARCHAR(50), -- books, readers, fines, etc.
    entity_id UUID,
    old_values JSONB,
    new_values JSONB,
    ip_address INET,
    user_agent TEXT,
    action_timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    details TEXT
);

-- Создание индексов для оптимизации запросов
-- Основные индексы для поиска
CREATE INDEX idx_users_username ON users (username);

CREATE INDEX idx_users_email ON users (email);

CREATE INDEX idx_users_role ON users (role);

CREATE INDEX idx_readers_ticket_number ON readers (ticket_number);

CREATE INDEX idx_readers_user_id ON readers (user_id);

CREATE INDEX idx_readers_hall_id ON readers (reading_hall_id);

-- Индексы для книг и авторов
CREATE INDEX idx_books_title ON books USING gin (to_tsvector ('russian', title));

CREATE INDEX idx_books_isbn ON books (isbn);

CREATE INDEX idx_authors_name ON authors USING gin (to_tsvector ('russian', full_name));

CREATE INDEX idx_book_copies_code ON book_copies (copy_code);

CREATE INDEX idx_book_copies_status ON book_copies (status);

CREATE INDEX idx_book_copies_book_id ON book_copies (book_id);

-- Индексы для выдач и штрафов
CREATE INDEX idx_book_issues_reader_id ON book_issues (reader_id);

CREATE INDEX idx_book_issues_copy_id ON book_issues (book_copy_id);

CREATE INDEX idx_book_issues_dates ON book_issues (issue_date, due_date, return_date);

CREATE INDEX idx_book_issues_active ON book_issues (reader_id)
WHERE
    return_date IS NULL;

CREATE INDEX idx_fines_reader_id ON fines (reader_id);

CREATE INDEX idx_fines_unpaid ON fines (reader_id)
WHERE
    is_paid = FALSE;

-- Индексы для рейтингов и логов
CREATE INDEX idx_book_ratings_book_id ON book_ratings (book_id);

CREATE INDEX idx_book_ratings_rating ON book_ratings (rating);

CREATE INDEX idx_system_logs_user_id ON system_logs (user_id);

CREATE INDEX idx_system_logs_timestamp ON system_logs (action_timestamp);

CREATE INDEX idx_system_logs_action_type ON system_logs (action_type);
