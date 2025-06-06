DROP SCHEMA IF EXISTS library CASCADE;

CREATE SCHEMA library;

SET search_path TO library;

-- Добавление модуля для UUID
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Создание типов данных
CREATE TYPE user_role AS ENUM ('administrator', 'librarian');

CREATE TYPE book_status AS ENUM (
    'available',
    'issued',
    'reserved',
    'lost',
    'damaged'
);

CREATE TYPE visit_type AS ENUM ('entry', 'exit');

-- 1. Таблица пользователей системы (только администраторы и библиотекари)
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    role user_role NOT NULL DEFAULT 'librarian',
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 2. Таблица читальных залов
CREATE TABLE reading_halls (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    hall_name VARCHAR(100) NOT NULL,
    specialization VARCHAR(200),
    total_seats INTEGER NOT NULL CHECK (total_seats > 0),
    current_visitors INTEGER DEFAULT 0 CHECK (current_visitors >= 0),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT chk_current_visitors CHECK (current_visitors <= total_seats)
);

-- 3. Таблица читателей (отдельная сущность, не связанная с пользователями)
CREATE TABLE readers (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    ticket_number VARCHAR(20) UNIQUE NOT NULL,
    full_name VARCHAR(200) NOT NULL,
    email VARCHAR(256), -- для уведомлений
    phone VARCHAR(20),
    registration_date DATE DEFAULT CURRENT_DATE,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 4. Таблица истории посещений залов
CREATE TABLE hall_visits (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    reader_id UUID NOT NULL REFERENCES readers(id),
    hall_id UUID NOT NULL REFERENCES reading_halls(id),
    visit_type visit_type NOT NULL,
    visit_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    librarian_id UUID REFERENCES users(id)
);

-- 5. Таблица авторов
CREATE TABLE authors (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    full_name VARCHAR(200) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 6. Таблица книг
CREATE TABLE books (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    title VARCHAR(500) NOT NULL,
    isbn VARCHAR(17),
    publication_year INTEGER,
    publisher VARCHAR(200),
    total_copies INTEGER NOT NULL DEFAULT 1 CHECK (total_copies > 0),
    available_copies INTEGER NOT NULL DEFAULT 1,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT chk_available_copies CHECK (
        available_copies >= 0 AND available_copies <= total_copies
    )
);

-- 7. Связующая таблица авторов и книг (многие ко многим)
CREATE TABLE book_authors (
    book_id UUID REFERENCES books(id) ON DELETE CASCADE,
    author_id UUID REFERENCES authors(id) ON DELETE CASCADE,
    PRIMARY KEY (book_id, author_id)
);

-- 8. Таблица экземпляров книг
CREATE TABLE book_copies (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    book_id UUID NOT NULL REFERENCES books(id) ON DELETE CASCADE,
    copy_code VARCHAR(50) UNIQUE NOT NULL,
    status book_status DEFAULT 'available',
    hall_id UUID REFERENCES reading_halls(id), -- в каком зале находится книга
    location_info TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 9. Таблица выдач книг
CREATE TABLE book_issues (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    reader_id UUID NOT NULL REFERENCES readers(id),
    book_copy_id UUID NOT NULL REFERENCES book_copies(id),
    issue_date DATE DEFAULT CURRENT_DATE,
    due_date DATE NOT NULL,
    return_date DATE,
    librarian_id UUID REFERENCES users(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT chk_dates CHECK (
        due_date >= issue_date AND
        (return_date IS NULL OR return_date >= issue_date)
    )
);

-- 10. Таблица штрафов
CREATE TABLE fines (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    reader_id UUID NOT NULL REFERENCES readers(id),
    book_issue_id UUID REFERENCES book_issues(id),
    amount DECIMAL(10, 2) NOT NULL CHECK (amount >= 0),
    reason VARCHAR(500) NOT NULL,
    fine_date DATE DEFAULT CURRENT_DATE,
    paid_date DATE,
    is_paid BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Создание основных индексов
CREATE INDEX idx_users_username ON users(username);
CREATE INDEX idx_readers_ticket_number ON readers(ticket_number);
CREATE INDEX idx_readers_email ON readers(email);

-- Индексы для залов и посещений
CREATE INDEX idx_reading_halls_specialization ON reading_halls(specialization);
CREATE INDEX idx_hall_visits_reader_id ON hall_visits(reader_id);
CREATE INDEX idx_hall_visits_hall_id ON hall_visits(hall_id);
CREATE INDEX idx_hall_visits_time ON hall_visits(visit_time);
CREATE INDEX idx_hall_visits_date ON hall_visits(DATE(visit_time));

-- Индексы для книг
CREATE INDEX idx_books_title ON books USING gin(to_tsvector('russian', title));
CREATE INDEX idx_books_isbn ON books(isbn);
CREATE INDEX idx_book_copies_code ON book_copies(copy_code);
CREATE INDEX idx_book_copies_status ON book_copies(status);
CREATE INDEX idx_book_copies_hall_id ON book_copies(hall_id);

-- Индексы для выдач и штрафов
CREATE INDEX idx_book_issues_reader_id ON book_issues(reader_id);
CREATE INDEX idx_book_issues_active ON book_issues(reader_id) WHERE return_date IS NULL;
CREATE INDEX idx_fines_reader_id ON fines(reader_id);
CREATE INDEX idx_fines_unpaid ON fines(reader_id) WHERE is_paid = FALSE;

-- Триггер для автоматического обновления счетчика посетителей в залах
CREATE OR REPLACE FUNCTION update_hall_visitors()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.visit_type = 'entry' THEN
        UPDATE reading_halls
        SET current_visitors = current_visitors + 1
        WHERE id = NEW.hall_id;
    ELSIF NEW.visit_type = 'exit' THEN
        UPDATE reading_halls
        SET current_visitors = GREATEST(current_visitors - 1, 0)
        WHERE id = NEW.hall_id;
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_update_hall_visitors
    AFTER INSERT ON hall_visits
    FOR EACH ROW
    EXECUTE FUNCTION update_hall_visitors();
