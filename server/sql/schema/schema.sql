CREATE TABLE halls (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL, -- Название зала
    library_name VARCHAR(150) NOT NULL, -- Название библиотеки
    specialization VARCHAR(100) NOT NULL, -- Специализация зала
    total_seats INTEGER NOT NULL DEFAULT 0, -- Общее количество мест
    occupied_seats INTEGER NOT NULL DEFAULT 0, -- Занятые места
    working_hours VARCHAR(50) DEFAULT '09:00-18:00', -- Время работы
    equipment TEXT, -- Оборудование
    status VARCHAR(20) DEFAULT 'open' CHECK (status IN ('open', 'closed', 'maintenance')), -- Статус зала
    visit_statistics INTEGER DEFAULT 0, -- Статистика посещений
    average_occupancy DECIMAL(5, 2) DEFAULT 0.00, -- Средняя загруженность %
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Создание триггера для обновления updated_at
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_halls_updated_at BEFORE UPDATE ON halls
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- 2. Таблица читателей
CREATE TABLE readers (
    id SERIAL PRIMARY KEY,
    full_name VARCHAR(150) NOT NULL, -- ФИО читателя
    ticket_number VARCHAR(20) UNIQUE NOT NULL, -- Номер читательского билета
    birth_date DATE NOT NULL, -- Дата рождения
    phone VARCHAR(20), -- Телефон
    email VARCHAR(100), -- Email
    education VARCHAR(100), -- Образование
    hall_id INTEGER NOT NULL, -- Закрепленный зал
    max_books_allowed INTEGER DEFAULT 5, -- Максимальное количество книг
    max_renewals_allowed INTEGER DEFAULT 2, -- Максимальное количество продлений
    total_debt DECIMAL(10, 2) DEFAULT 0.00, -- Общая задолженность
    status VARCHAR(20) DEFAULT 'active' CHECK (status IN ('active', 'blocked', 'debtor', 'inactive')), -- Статус читателя
    reader_rating INTEGER DEFAULT 0, -- Рейтинг читателя
    registration_date DATE DEFAULT CURRENT_DATE, -- Дата регистрации
    last_activity_date DATE, -- Дата последней активности
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_readers_hall FOREIGN KEY (hall_id) REFERENCES halls (id) ON DELETE RESTRICT ON UPDATE CASCADE
);

CREATE INDEX idx_readers_ticket_number ON readers (ticket_number);
CREATE INDEX idx_readers_status ON readers (status);
CREATE INDEX idx_readers_hall_id ON readers (hall_id);

CREATE TRIGGER update_readers_updated_at BEFORE UPDATE ON readers
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- 3. Таблица категорий книг
CREATE TABLE book_categories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE, -- Название категории
    description TEXT, -- Описание категории
    default_loan_days INTEGER DEFAULT 30, -- Стандартный срок выдачи для категории
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 4. Таблица книг
CREATE TABLE books (
    id SERIAL PRIMARY KEY,
    title VARCHAR(200) NOT NULL, -- Название книги
    author VARCHAR(150) NOT NULL, -- Автор
    publication_year INTEGER NOT NULL, -- Год издания
    isbn VARCHAR(20), -- ISBN
    book_code VARCHAR(50) UNIQUE NOT NULL, -- Шифр книги
    category_id INTEGER, -- Категория книги
    hall_id INTEGER NOT NULL, -- Зал, где находится книга
    total_copies INTEGER NOT NULL DEFAULT 1, -- Общее количество экземпляров
    available_copies INTEGER NOT NULL DEFAULT 1, -- Доступные экземпляры
    condition_status VARCHAR(20) DEFAULT 'good' CHECK (condition_status IN ('excellent', 'good', 'fair', 'poor')), -- Состояние книги
    location_info VARCHAR(100), -- Местоположение (стеллаж, полка)
    max_loan_days INTEGER DEFAULT 30, -- Максимальный срок выдачи
    max_renewals INTEGER DEFAULT 2, -- Максимальное количество продлений
    popularity_score INTEGER DEFAULT 0, -- Счетчик популярности
    rating DECIMAL(3, 2) DEFAULT 0.00, -- Рейтинг книги (0-10)
    acquisition_date DATE DEFAULT CURRENT_DATE, -- Дата поступления
    status VARCHAR(20) DEFAULT 'available' CHECK (status IN ('available', 'loaned', 'reserved', 'maintenance', 'lost')), -- Статус книги
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_books_category FOREIGN KEY (category_id) REFERENCES book_categories (id) ON DELETE SET NULL ON UPDATE CASCADE,
    CONSTRAINT fk_books_hall FOREIGN KEY (hall_id) REFERENCES halls (id) ON DELETE RESTRICT ON UPDATE CASCADE
);

CREATE INDEX idx_books_book_code ON books (book_code);
CREATE INDEX idx_books_author ON books (author);
CREATE INDEX idx_books_title ON books (title);
CREATE INDEX idx_books_category_id ON books (category_id);
CREATE INDEX idx_books_hall_id ON books (hall_id);
CREATE INDEX idx_books_status ON books (status);

CREATE TRIGGER update_books_updated_at BEFORE UPDATE ON books
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- 5. Таблица сотрудников (библиотекарей)
CREATE TABLE librarians (
    id SERIAL PRIMARY KEY,
    full_name VARCHAR(150) NOT NULL, -- ФИО сотрудника
    employee_id VARCHAR(20) UNIQUE NOT NULL, -- Табельный номер
    position VARCHAR(100) DEFAULT 'Библиотекарь', -- Должность
    phone VARCHAR(20), -- Телефон
    email VARCHAR(100), -- Email
    hire_date DATE DEFAULT CURRENT_DATE, -- Дата приема на работу
    status VARCHAR(20) DEFAULT 'active' CHECK (status IN ('active', 'inactive')), -- Статус сотрудника
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_librarians_employee_id ON librarians (employee_id);
CREATE INDEX idx_librarians_status ON librarians (status);

CREATE TRIGGER update_librarians_updated_at BEFORE UPDATE ON librarians
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- 6. Основная таблица истории выдачи книг
CREATE TABLE loan_history (
    id SERIAL PRIMARY KEY,
    book_id INTEGER NOT NULL, -- ID книги
    reader_id INTEGER NOT NULL, -- ID читателя
    librarian_id INTEGER NOT NULL, -- ID библиотекаря, выдавшего книгу
    loan_date DATE NOT NULL DEFAULT CURRENT_DATE, -- Дата выдачи
    due_date DATE NOT NULL, -- Планируемая дата возврата
    return_date DATE, -- Фактическая дата возврата
    renewals_count INTEGER DEFAULT 0, -- Количество продлений
    status VARCHAR(20) DEFAULT 'active' CHECK (status IN ('active', 'returned', 'overdue', 'lost', 'damaged')), -- Статус выдачи
    fine_amount DECIMAL(10, 2) DEFAULT 0.00, -- Сумма штрафа
    fine_paid BOOLEAN DEFAULT FALSE, -- Штраф оплачен
    comments TEXT, -- Комментарии
    return_librarian_id INTEGER, -- ID библиотекаря, принявшего книгу
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_loan_history_book FOREIGN KEY (book_id) REFERENCES books (id) ON DELETE RESTRICT ON UPDATE CASCADE,
    CONSTRAINT fk_loan_history_reader FOREIGN KEY (reader_id) REFERENCES readers (id) ON DELETE RESTRICT ON UPDATE CASCADE,
    CONSTRAINT fk_loan_history_librarian FOREIGN KEY (librarian_id) REFERENCES librarians (id) ON DELETE RESTRICT ON UPDATE CASCADE,
    CONSTRAINT fk_loan_history_return_librarian FOREIGN KEY (return_librarian_id) REFERENCES librarians (id) ON DELETE RESTRICT ON UPDATE CASCADE
);

CREATE INDEX idx_loan_history_book_id ON loan_history (book_id);
CREATE INDEX idx_loan_history_reader_id ON loan_history (reader_id);
CREATE INDEX idx_loan_history_loan_date ON loan_history (loan_date);
CREATE INDEX idx_loan_history_due_date ON loan_history (due_date);
CREATE INDEX idx_loan_history_status ON loan_history (status);
CREATE INDEX idx_loan_history_active_loans ON loan_history (status, due_date); -- Составной индекс для активных выдач

CREATE TRIGGER update_loan_history_updated_at BEFORE UPDATE ON loan_history
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- 7. Таблица продлений
CREATE TABLE renewals (
    id SERIAL PRIMARY KEY,
    loan_history_id INTEGER NOT NULL, -- ID записи в истории выдач
    renewal_date DATE NOT NULL DEFAULT CURRENT_DATE, -- Дата продления
    old_due_date DATE NOT NULL, -- Старая дата возврата
    new_due_date DATE NOT NULL, -- Новая дата возврата
    librarian_id INTEGER NOT NULL, -- ID библиотекаря, сделавшего продление
    reason VARCHAR(200), -- Причина продления
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_renewals_loan_history FOREIGN KEY (loan_history_id) REFERENCES loan_history (id) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT fk_renewals_librarian FOREIGN KEY (librarian_id) REFERENCES librarians (id) ON DELETE RESTRICT ON UPDATE CASCADE
);

CREATE INDEX idx_renewals_loan_history_id ON renewals (loan_history_id);
CREATE INDEX idx_renewals_renewal_date ON renewals (renewal_date);

-- 8. Таблица бронирования книг
CREATE TABLE reservations (
    id SERIAL PRIMARY KEY,
    book_id INTEGER NOT NULL, -- ID книги
    reader_id INTEGER NOT NULL, -- ID читателя
    reservation_date DATE NOT NULL DEFAULT CURRENT_DATE, -- Дата бронирования
    expiration_date DATE NOT NULL, -- Дата истечения брони
    status VARCHAR(20) DEFAULT 'active' CHECK (status IN ('active', 'fulfilled', 'cancelled', 'expired')), -- Статус брони
    priority_order INTEGER DEFAULT 1, -- Порядок в очереди
    notification_sent BOOLEAN DEFAULT FALSE, -- Уведомление отправлено
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_reservations_book FOREIGN KEY (book_id) REFERENCES books (id) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT fk_reservations_reader FOREIGN KEY (reader_id) REFERENCES readers (id) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE INDEX idx_reservations_book_id ON reservations (book_id);
CREATE INDEX idx_reservations_reader_id ON reservations (reader_id);
CREATE INDEX idx_reservations_status ON reservations (status);
CREATE INDEX idx_reservations_priority ON reservations (book_id, priority_order);

CREATE TRIGGER update_reservations_updated_at BEFORE UPDATE ON reservations
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- 9. Таблица штрафов
CREATE TABLE fines (
    id SERIAL PRIMARY KEY,
    loan_history_id INTEGER NOT NULL, -- ID записи в истории выдач
    reader_id INTEGER NOT NULL, -- ID читателя
    fine_type VARCHAR(20) NOT NULL CHECK (fine_type IN ('overdue', 'damage', 'loss', 'other')), -- Тип штрафа
    amount DECIMAL(10, 2) NOT NULL, -- Сумма штрафа
    fine_date DATE NOT NULL DEFAULT CURRENT_DATE, -- Дата начисления штрафа
    payment_date DATE, -- Дата оплаты
    status VARCHAR(20) DEFAULT 'unpaid' CHECK (status IN ('unpaid', 'paid', 'waived')), -- Статус оплаты
    description TEXT, -- Описание штрафа
    librarian_id INTEGER NOT NULL, -- ID библиотекаря, начислившего штраф
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_fines_loan_history FOREIGN KEY (loan_history_id) REFERENCES loan_history (id) ON DELETE RESTRICT ON UPDATE CASCADE,
    CONSTRAINT fk_fines_reader FOREIGN KEY (reader_id) REFERENCES readers (id) ON DELETE RESTRICT ON UPDATE CASCADE,
    CONSTRAINT fk_fines_librarian FOREIGN KEY (librarian_id) REFERENCES librarians (id) ON DELETE RESTRICT ON UPDATE CASCADE
);

CREATE INDEX idx_fines_reader_id ON fines (reader_id);
CREATE INDEX idx_fines_status ON fines (status);
CREATE INDEX idx_fines_fine_date ON fines (fine_date);

CREATE TRIGGER update_fines_updated_at BEFORE UPDATE ON fines
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- 10. Таблица операций/логов
CREATE TABLE operation_logs (
    id SERIAL PRIMARY KEY,
    operation_type VARCHAR(20) NOT NULL CHECK (operation_type IN ('loan', 'return', 'renewal', 'reservation', 'fine', 'registration', 'book_add', 'book_remove')),
    entity_type VARCHAR(20) NOT NULL CHECK (entity_type IN ('book', 'reader', 'hall', 'loan', 'reservation')),
    entity_id INTEGER NOT NULL, -- ID сущности, с которой производилась операция
    librarian_id INTEGER, -- ID библиотекаря, выполнившего операцию
    operation_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    details JSONB, -- Детали операции в формате JSON
    description TEXT -- Описание операции
);

CREATE INDEX idx_operation_logs_operation_type ON operation_logs (operation_type);
CREATE INDEX idx_operation_logs_entity ON operation_logs (entity_type, entity_id);
CREATE INDEX idx_operation_logs_operation_date ON operation_logs (operation_date);
CREATE INDEX idx_operation_logs_librarian_id ON operation_logs (librarian_id);

-- 11. Таблица статистики по дням (для аналитики)
CREATE TABLE daily_statistics (
    id SERIAL PRIMARY KEY,
    stat_date DATE NOT NULL UNIQUE,
    total_loans INTEGER DEFAULT 0, -- Выдач за день
    total_returns INTEGER DEFAULT 0, -- Возвратов за день
    total_renewals INTEGER DEFAULT 0, -- Продлений за день
    total_reservations INTEGER DEFAULT 0, -- Бронирований за день
    total_new_readers INTEGER DEFAULT 0, -- Новых читателей за день
    total_fines_amount DECIMAL(10, 2) DEFAULT 0.00, -- Сумма штрафов за день
    overdue_books INTEGER DEFAULT 0, -- Просроченных книг на конец дня
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_daily_statistics_stat_date ON daily_statistics (stat_date);

CREATE TRIGGER update_daily_statistics_updated_at BEFORE UPDATE ON daily_statistics
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Комментарии к таблицам
COMMENT ON TABLE halls IS 'Таблица читальных залов';
COMMENT ON TABLE readers IS 'Таблица читателей';
COMMENT ON TABLE book_categories IS 'Таблица категорий книг';
COMMENT ON TABLE books IS 'Таблица книг';
COMMENT ON TABLE librarians IS 'Таблица сотрудников (библиотекарей)';
COMMENT ON TABLE loan_history IS 'Основная таблица истории выдачи книг';
COMMENT ON TABLE renewals IS 'Таблица продлений';
COMMENT ON TABLE reservations IS 'Таблица бронирования книг';
COMMENT ON TABLE fines IS 'Таблица штрафов';
COMMENT ON TABLE operation_logs IS 'Таблица операций/логов';
COMMENT ON TABLE daily_statistics IS 'Таблица статистики по дням (для аналитики)';
