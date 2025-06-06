import { addDays, subDays } from "date-fns";
import type {
  BookIssue,
  Fine,
  ReadingHall,
  HallVisit,
  Author,
  BookCopy,
  User,
  BookWithDetails,
  ReaderWithDetails,
  DashboardStats,
} from "../types";

// Mock Authors (российские и зарубежные писатели)
export const authors: Author[] = [
  {
    id: "auth-1",
    full_name: "Лев Николаевич Толстой",
    created_at: new Date("2020-01-15"),
  },
  {
    id: "auth-2",
    full_name: "Фёдор Михайлович Достоевский",
    created_at: new Date("2020-01-15"),
  },
  {
    id: "auth-3",
    full_name: "Антон Павлович Чехов",
    created_at: new Date("2020-01-15"),
  },
  {
    id: "auth-4",
    full_name: "Александр Сергеевич Пушкин",
    created_at: new Date("2020-01-15"),
  },
  {
    id: "auth-5",
    full_name: "Михаил Юрьевич Лермонтов",
    created_at: new Date("2020-01-15"),
  },
  {
    id: "auth-6",
    full_name: "Николай Васильевич Гоголь",
    created_at: new Date("2020-01-15"),
  },
  {
    id: "auth-7",
    full_name: "Иван Сергеевич Тургенев",
    created_at: new Date("2020-01-15"),
  },
  {
    id: "auth-8",
    full_name: "Михаил Афанасьевич Булгаков",
    created_at: new Date("2020-01-15"),
  },
  {
    id: "auth-9",
    full_name: "Александр Исаевич Солженицын",
    created_at: new Date("2020-01-15"),
  },
  {
    id: "auth-10",
    full_name: "Джордж Оруэлл",
    created_at: new Date("2020-01-15"),
  },
  {
    id: "auth-11",
    full_name: "Джейн Остин",
    created_at: new Date("2020-01-15"),
  },
  {
    id: "auth-12",
    full_name: "Агата Кристи",
    created_at: new Date("2020-01-15"),
  },
  {
    id: "auth-13",
    full_name: "Стивен Кинг",
    created_at: new Date("2020-01-15"),
  },
  {
    id: "auth-14",
    full_name: "Харуки Мураками",
    created_at: new Date("2020-01-15"),
  },
  {
    id: "auth-15",
    full_name: "Габриэль Гарсиа Маркес",
    created_at: new Date("2020-01-15"),
  },
];

// Mock Books (расширенный список с реальными книгами)
export const books: BookWithDetails[] = [
  {
    id: "book-1",
    title: "Война и мир",
    isbn: "9785170959525",
    publication_year: 1869,
    publisher: "АСТ",
    total_copies: 8,
    available_copies: 5,
    created_at: new Date("2020-01-15"),
    authors: [authors[0]],
  },
  {
    id: "book-2",
    title: "Анна Каренина",
    isbn: "9785170895472",
    publication_year: 1877,
    publisher: "АСТ",
    total_copies: 6,
    available_copies: 3,
    created_at: new Date("2020-01-20"),
    authors: [authors[0]],
  },
  {
    id: "book-3",
    title: "Преступление и наказание",
    isbn: "9785170901234",
    publication_year: 1866,
    publisher: "Эксмо",
    total_copies: 7,
    available_copies: 2,
    created_at: new Date("2020-02-10"),
    authors: [authors[1]],
  },
  {
    id: "book-4",
    title: "Братья Карамазовы",
    isbn: "9785170845673",
    publication_year: 1880,
    publisher: "Эксмо",
    total_copies: 5,
    available_copies: 4,
    created_at: new Date("2020-02-15"),
    authors: [authors[1]],
  },
  {
    id: "book-5",
    title: "Вишнёвый сад",
    isbn: "9785170923456",
    publication_year: 1904,
    publisher: "Азбука",
    total_copies: 4,
    available_copies: 4,
    created_at: new Date("2020-03-05"),
    authors: [authors[2]],
  },
  {
    id: "book-6",
    title: "Евгений Онегин",
    isbn: "9785170778901",
    publication_year: 1833,
    publisher: "АСТ",
    total_copies: 6,
    available_copies: 1,
    created_at: new Date("2020-03-12"),
    authors: [authors[3]],
  },
  {
    id: "book-7",
    title: "Герой нашего времени",
    isbn: "9785170856789",
    publication_year: 1840,
    publisher: "Эксмо",
    total_copies: 5,
    available_copies: 3,
    created_at: new Date("2020-04-18"),
    authors: [authors[4]],
  },
  {
    id: "book-8",
    title: "Мёртвые души",
    isbn: "9785170934567",
    publication_year: 1842,
    publisher: "АСТ",
    total_copies: 4,
    available_copies: 2,
    created_at: new Date("2020-04-25"),
    authors: [authors[5]],
  },
  {
    id: "book-9",
    title: "Отцы и дети",
    isbn: "9785170812345",
    publication_year: 1862,
    publisher: "Азбука",
    total_copies: 5,
    available_copies: 5,
    created_at: new Date("2020-05-10"),
    authors: [authors[6]],
  },
  {
    id: "book-10",
    title: "Мастер и Маргарита",
    isbn: "9785170789012",
    publication_year: 1967,
    publisher: "АСТ",
    total_copies: 9,
    available_copies: 4,
    created_at: new Date("2020-05-20"),
    authors: [authors[7]],
  },
  {
    id: "book-11",
    title: "Белая гвардия",
    isbn: "9785170945123",
    publication_year: 1925,
    publisher: "Эксмо",
    total_copies: 3,
    available_copies: 1,
    created_at: new Date("2020-06-01"),
    authors: [authors[7]],
  },
  {
    id: "book-12",
    title: "Архипелаг ГУЛАГ",
    isbn: "9785170867890",
    publication_year: 1973,
    publisher: "АСТ",
    total_copies: 4,
    available_copies: 2,
    created_at: new Date("2020-06-15"),
    authors: [authors[8]],
  },
  {
    id: "book-13",
    title: "1984",
    isbn: "9785170923789",
    publication_year: 1949,
    publisher: "АСТ",
    total_copies: 7,
    available_copies: 0,
    created_at: new Date("2020-07-01"),
    authors: [authors[9]],
  },
  {
    id: "book-14",
    title: "Скотный двор",
    isbn: "9785170834567",
    publication_year: 1945,
    publisher: "АСТ",
    total_copies: 5,
    available_copies: 3,
    created_at: new Date("2020-07-10"),
    authors: [authors[9]],
  },
  {
    id: "book-15",
    title: "Гордость и предубеждение",
    isbn: "9785170745632",
    publication_year: 1813,
    publisher: "Эксмо",
    total_copies: 4,
    available_copies: 2,
    created_at: new Date("2020-08-05"),
    authors: [authors[10]],
  },
  {
    id: "book-16",
    title: "Убийство в Восточном экспрессе",
    isbn: "9785170892134",
    publication_year: 1934,
    publisher: "АСТ",
    total_copies: 6,
    available_copies: 4,
    created_at: new Date("2020-08-20"),
    authors: [authors[11]],
  },
  {
    id: "book-17",
    title: "Сияние",
    isbn: "9785170756789",
    publication_year: 1977,
    publisher: "АСТ",
    total_copies: 5,
    available_copies: 1,
    created_at: new Date("2020-09-10"),
    authors: [authors[12]],
  },
  {
    id: "book-18",
    title: "Норвежский лес",
    isbn: "9785170823456",
    publication_year: 1987,
    publisher: "Эксмо",
    total_copies: 4,
    available_copies: 3,
    created_at: new Date("2020-09-25"),
    authors: [authors[13]],
  },
  {
    id: "book-19",
    title: "Сто лет одиночества",
    isbn: "9785170901567",
    publication_year: 1967,
    publisher: "АСТ",
    total_copies: 3,
    available_copies: 0,
    created_at: new Date("2020-10-12"),
    authors: [authors[14]],
  },
  {
    id: "book-20",
    title: "Кафка на пляже",
    isbn: "9785170834123",
    publication_year: 2002,
    publisher: "Эксмо",
    total_copies: 4,
    available_copies: 2,
    created_at: new Date("2020-10-28"),
    authors: [authors[13]],
  },
];

// Mock Book Copies (реалистичные коды и локации)
export const bookCopies: BookCopy[] = [
  // Копии для "Война и мир"
  {
    id: "copy-1",
    book_id: "book-1",
    copy_code: "РЛ-821.161-Т53-001",
    status: "available",
    hall_id: "hall-1",
    location_info: "Зал 1, Стеллаж А-15, Полка 3",
    created_at: new Date("2020-01-15"),
  },
  {
    id: "copy-2",
    book_id: "book-1",
    copy_code: "РЛ-821.161-Т53-002",
    status: "issued",
    hall_id: "hall-1",
    location_info: "Зал 1, Стеллаж А-15, Полка 3",
    created_at: new Date("2020-01-15"),
  },
  {
    id: "copy-3",
    book_id: "book-1",
    copy_code: "РЛ-821.161-Т53-003",
    status: "available",
    hall_id: "hall-1",
    location_info: "Зал 1, Стеллаж А-15, Полка 3",
    created_at: new Date("2020-01-15"),
  },

  // Копии для "Преступление и наказание"
  {
    id: "copy-4",
    book_id: "book-3",
    copy_code: "РЛ-821.161-Д67-001",
    status: "issued",
    hall_id: "hall-1",
    location_info: "Зал 1, Стеллаж А-12, Полка 2",
    created_at: new Date("2020-02-10"),
  },
  {
    id: "copy-5",
    book_id: "book-3",
    copy_code: "РЛ-821.161-Д67-002",
    status: "issued",
    hall_id: "hall-1",
    location_info: "Зал 1, Стеллаж А-12, Полка 2",
    created_at: new Date("2020-02-10"),
  },

  // Копии для "1984"
  {
    id: "copy-6",
    book_id: "book-13",
    copy_code: "ЗЛ-820-О79-001",
    status: "issued",
    hall_id: "hall-2",
    location_info: "Зал 2, Стеллаж Б-08, Полка 1",
    created_at: new Date("2020-07-01"),
  },
  {
    id: "copy-7",
    book_id: "book-13",
    copy_code: "ЗЛ-820-О79-002",
    status: "issued",
    hall_id: "hall-2",
    location_info: "Зал 2, Стеллаж Б-08, Полка 1",
    created_at: new Date("2020-07-01"),
  },

  // Копии для "Мастер и Маргарита"
  {
    id: "copy-8",
    book_id: "book-10",
    copy_code: "РЛ-821.161-Б90-001",
    status: "issued",
    hall_id: "hall-1",
    location_info: "Зал 1, Стеллаж А-18, Полка 4",
    created_at: new Date("2020-05-20"),
  },
  {
    id: "copy-9",
    book_id: "book-10",
    copy_code: "РЛ-821.161-Б90-002",
    status: "available",
    hall_id: "hall-1",
    location_info: "Зал 1, Стеллаж А-18, Полка 4",
    created_at: new Date("2020-05-20"),
  },

  // Копии для "Евгений Онегин"
  {
    id: "copy-10",
    book_id: "book-6",
    copy_code: "РЛ-821.161-П87-001",
    status: "issued",
    hall_id: "hall-1",
    location_info: "Зал 1, Стеллаж А-14, Полка 1",
    created_at: new Date("2020-03-12"),
  },
];

// Mock Readers (реалистичные российские имена и данные)
export const readers: ReaderWithDetails[] = [
  {
    id: "reader-1",
    ticket_number: "2024-001",
    full_name: "Петров Иван Александрович",
    email: "i.petrov@yandex.ru",
    phone: "+7 (903) 123-45-67",
    registration_date: new Date("2024-01-15"),
    is_active: true,
    created_at: new Date("2024-01-15"),
    total_fines_amount: 0,
  },
  {
    id: "reader-2",
    ticket_number: "2024-002",
    full_name: "Сидорова Мария Петровна",
    email: "m.sidorova@gmail.com",
    phone: "+7 (916) 987-65-43",
    registration_date: new Date("2024-01-22"),
    is_active: true,
    created_at: new Date("2024-01-22"),
    total_fines_amount: 50.0,
  },
  {
    id: "reader-3",
    ticket_number: "2024-003",
    full_name: "Козлов Александр Сергеевич",
    email: "a.kozlov@mail.ru",
    phone: "+7 (925) 456-78-90",
    registration_date: new Date("2024-02-03"),
    is_active: true,
    created_at: new Date("2024-02-03"),
    total_fines_amount: 125.5,
  },
  {
    id: "reader-4",
    ticket_number: "2023-087",
    full_name: "Смирнова Елена Викторовна",
    email: "e.smirnova@yandex.ru",
    phone: "+7 (909) 234-56-78",
    registration_date: new Date("2023-09-12"),
    is_active: true,
    created_at: new Date("2023-09-12"),
    total_fines_amount: 0,
  },
  {
    id: "reader-5",
    ticket_number: "2023-156",
    full_name: "Морозов Дмитрий Андреевич",
    email: "d.morozov@gmail.com",
    phone: "+7 (917) 345-67-89",
    registration_date: new Date("2023-11-28"),
    is_active: false,
    created_at: new Date("2023-11-28"),
    total_fines_amount: 75.0,
  },
  {
    id: "reader-6",
    ticket_number: "2024-004",
    full_name: "Волкова Анна Михайловна",
    email: "a.volkova@yandex.ru",
    phone: "+7 (926) 567-89-01",
    registration_date: new Date("2024-02-14"),
    is_active: true,
    created_at: new Date("2024-02-14"),
    total_fines_amount: 0,
  },
  {
    id: "reader-7",
    ticket_number: "2023-203",
    full_name: "Лебедев Сергей Николаевич",
    email: "s.lebedev@mail.ru",
    phone: "+7 (903) 678-90-12",
    registration_date: new Date("2023-12-05"),
    is_active: true,
    created_at: new Date("2023-12-05"),
    total_fines_amount: 25.0,
  },
  {
    id: "reader-8",
    ticket_number: "2024-005",
    full_name: "Новикова Ольга Дмитриевна",
    email: "o.novikova@gmail.com",
    phone: "+7 (915) 789-01-23",
    registration_date: new Date("2024-02-20"),
    is_active: true,
    created_at: new Date("2024-02-20"),
    total_fines_amount: 0,
  },
];

// Mock Users (сотрудники библиотеки)
export const users: User[] = [
  {
    id: "user-1",
    username: "admin",
    email: "admin@library.ru",
    password_hash:
      "$2b$12$LQv3c1yqBWVHxkd0LHAkCOYz6TtxMQJqhN8/LeWNcIrDtpKIWTB5m",
    role: "administrator",
    is_active: true,
    created_at: new Date("2020-01-01"),
  },
  {
    id: "user-2",
    username: "bibliothekar1",
    email: "l.ivanova@library.ru",
    password_hash:
      "$2b$12$LQv3c1yqBWVHxkd0LHAkCOYz6TtxMQJqhN8/LeWNcIrDtpKIWTB5m",
    role: "librarian",
    is_active: true,
    created_at: new Date("2020-06-15"),
  },
  {
    id: "user-3",
    username: "bibliothekar2",
    email: "n.fedorova@library.ru",
    password_hash:
      "$2b$12$LQv3c1yqBWVHxkd0LHAkCOYz6TtxMQJqhN8/LeWNcIrDtpKIWTB5m",
    role: "librarian",
    is_active: true,
    created_at: new Date("2021-03-10"),
  },
];

// Mock Book Issues (реалистичные даты выдачи)
export const bookIssues: BookIssue[] = [
  {
    id: "issue-1",
    reader_id: "reader-1",
    book_copy_id: "copy-2",
    issue_date: subDays(new Date(), 12),
    due_date: addDays(new Date(), 18),
    librarian_id: "user-2",
    created_at: subDays(new Date(), 12),
  },
  {
    id: "issue-2",
    reader_id: "reader-2",
    book_copy_id: "copy-4",
    issue_date: subDays(new Date(), 8),
    due_date: addDays(new Date(), 22),
    librarian_id: "user-2",
    created_at: subDays(new Date(), 8),
  },
  {
    id: "issue-3",
    reader_id: "reader-3",
    book_copy_id: "copy-5",
    issue_date: subDays(new Date(), 35),
    due_date: subDays(new Date(), 5), // Просрочена
    librarian_id: "user-3",
    created_at: subDays(new Date(), 35),
  },
  {
    id: "issue-4",
    reader_id: "reader-4",
    book_copy_id: "copy-6",
    issue_date: subDays(new Date(), 18),
    due_date: addDays(new Date(), 12),
    librarian_id: "user-2",
    created_at: subDays(new Date(), 18),
  },
  {
    id: "issue-5",
    reader_id: "reader-6",
    book_copy_id: "copy-7",
    issue_date: subDays(new Date(), 25),
    due_date: addDays(new Date(), 5),
    librarian_id: "user-3",
    created_at: subDays(new Date(), 25),
  },
  {
    id: "issue-6",
    reader_id: "reader-7",
    book_copy_id: "copy-8",
    issue_date: subDays(new Date(), 42),
    due_date: subDays(new Date(), 12), // Просрочена
    librarian_id: "user-2",
    created_at: subDays(new Date(), 42),
  },
  {
    id: "issue-7",
    reader_id: "reader-8",
    book_copy_id: "copy-10",
    issue_date: subDays(new Date(), 5),
    due_date: addDays(new Date(), 25),
    librarian_id: "user-3",
    created_at: subDays(new Date(), 5),
  },
];

// Mock Reading Halls (реалистичные залы российской библиотеки)
export const readingHalls: ReadingHall[] = [
  {
    id: "hall-1",
    hall_name: "Общий читальный зал",
    specialization: "Художественная и научная литература",
    total_seats: 180,
    current_visitors: 142,
    created_at: new Date("2019-09-01"),
  },
  {
    id: "hall-2",
    hall_name: "Зал тихого чтения",
    specialization: "Индивидуальная работа с документами",
    total_seats: 60,
    current_visitors: 38,
    created_at: new Date("2019-09-01"),
  },
  {
    id: "hall-3",
    hall_name: "Электронный зал",
    specialization: "Работа с цифровыми ресурсами",
    total_seats: 40,
    current_visitors: 24,
    created_at: new Date("2020-02-15"),
  },
  {
    id: "hall-4",
    hall_name: "Зал периодических изданий",
    specialization: "Газеты, журналы, справочники",
    total_seats: 80,
    current_visitors: 31,
    created_at: new Date("2019-09-01"),
  },
  {
    id: "hall-5",
    hall_name: "Детский читальный зал",
    specialization: "Литература для детей и подростков",
    total_seats: 50,
    current_visitors: 18,
    created_at: new Date("2020-01-10"),
  },
];

// Mock Hall Visits (реалистичные посещения в течение дня)
export const hallVisits: HallVisit[] = [
  // Утренние посещения
  {
    id: "visit-1",
    reader_id: "reader-1",
    hall_id: "hall-1",
    visit_type: "entry",
    visit_time: new Date(new Date().setHours(9, 15, 0, 0)),
    librarian_id: "user-2",
  },
  {
    id: "visit-2",
    reader_id: "reader-2",
    hall_id: "hall-2",
    visit_type: "entry",
    visit_time: new Date(new Date().setHours(9, 32, 0, 0)),
    librarian_id: "user-2",
  },
  {
    id: "visit-3",
    reader_id: "reader-3",
    hall_id: "hall-1",
    visit_type: "entry",
    visit_time: new Date(new Date().setHours(10, 8, 0, 0)),
    librarian_id: "user-3",
  },
  {
    id: "visit-4",
    reader_id: "reader-4",
    hall_id: "hall-3",
    visit_type: "entry",
    visit_time: new Date(new Date().setHours(10, 45, 0, 0)),
    librarian_id: "user-2",
  },

  // Дневные посещения
  {
    id: "visit-5",
    reader_id: "reader-5",
    hall_id: "hall-4",
    visit_type: "entry",
    visit_time: new Date(new Date().setHours(11, 20, 0, 0)),
    librarian_id: "user-3",
  },
  {
    id: "visit-6",
    reader_id: "reader-6",
    hall_id: "hall-1",
    visit_type: "entry",
    visit_time: new Date(new Date().setHours(12, 15, 0, 0)),
    librarian_id: "user-2",
  },
  {
    id: "visit-7",
    reader_id: "reader-7",
    hall_id: "hall-2",
    visit_type: "entry",
    visit_time: new Date(new Date().setHours(13, 30, 0, 0)),
    librarian_id: "user-3",
  },

  // Выходы
  {
    id: "visit-8",
    reader_id: "reader-2",
    hall_id: "hall-2",
    visit_type: "exit",
    visit_time: new Date(new Date().setHours(14, 20, 0, 0)),
    librarian_id: "user-2",
  },
  {
    id: "visit-9",
    reader_id: "reader-5",
    hall_id: "hall-4",
    visit_type: "exit",
    visit_time: new Date(new Date().setHours(15, 45, 0, 0)),
    librarian_id: "user-3",
  },

  // Вечерние посещения
  {
    id: "visit-10",
    reader_id: "reader-8",
    hall_id: "hall-1",
    visit_type: "entry",
    visit_time: new Date(new Date().setHours(16, 10, 0, 0)),
    librarian_id: "user-2",
  },
];

// Mock Fines (реалистичные штрафы)
export const fines: Fine[] = [
  {
    id: "fine-1",
    reader_id: "reader-2",
    book_issue_id: "issue-2",
    amount: 500.0,
    reason: "Просрочка возврата на 10 дней (5 руб./день)",
    fine_date: subDays(new Date(), 3),
    is_paid: false,
    created_at: subDays(new Date(), 3),
  },
  {
    id: "fine-2",
    reader_id: "reader-3",
    book_issue_id: "issue-3",
    amount: 1370,
    reason: "Просрочка возврата на 17 дней + повреждение обложки",
    fine_date: subDays(new Date(), 8),
    is_paid: false,
    created_at: subDays(new Date(), 8),
  },
  {
    id: "fine-3",
    reader_id: "reader-5",
    book_issue_id: "issue-4",
    amount: 750.0,
    reason: "Просрочка возврата на 15 дней",
    fine_date: subDays(new Date(), 15),
    is_paid: true,
    created_at: subDays(new Date(), 15),
  },
  {
    id: "fine-4",
    reader_id: "reader-7",
    book_issue_id: "issue-6",
    amount: 308.0,
    reason: "Просрочка возврата на 5 дней",
    fine_date: subDays(new Date(), 1),
    is_paid: false,
    created_at: subDays(new Date(), 1),
  },
];

// Mock Dashboard Stats
export const dashboardStats: DashboardStats = {
  total_books: books.length,
  total_readers: readers.filter((r) => r.is_active).length,
  active_issues: bookIssues.filter((issue) => !issue.return_date).length,
  overdue_issues: bookIssues.filter(
    (issue) => !issue.return_date && new Date(issue.due_date) < new Date(),
  ).length,
  total_fines: fines
    .filter((f) => !f.is_paid)
    .reduce((sum, fine) => sum + fine.amount, 0),
  current_hall_visitors: readingHalls.reduce(
    (sum, hall) => sum + hall.current_visitors,
    0,
  ),
};

// Helper functions
export const getBookStatus = (
  book: BookWithDetails,
): "available" | "checked-out" | "reserved" => {
  if (book.available_copies === 0) return "checked-out";
  return "available";
};

export const members = readers;
export const checkoutRecords = bookIssues;
export const rooms = readingHalls;

// Extended data with relationships
export const booksWithDetails = books.map((book) => ({
  ...book,
  status: getBookStatus(book),
  coverImage: `https://picsum.photos/200/300?random=${book.id}`,
  genre: book.authors[0]?.full_name.includes("Толстой")
    ? ["Классическая литература", "Роман"]
    : book.authors[0]?.full_name.includes("Оруэлл")
      ? ["Антиутопия", "Научная фантастика"]
      : book.authors[0]?.full_name.includes("Кристи")
        ? ["Детектив", "Мистика"]
        : book.authors[0]?.full_name.includes("Кинг")
          ? ["Ужасы", "Триллер"]
          : ["Художественная литература"],
  description: `Классическое произведение "${book.title}" - одна из жемчужин мировой литературы.`,
  location:
    bookCopies.find((copy) => copy.book_id === book.id)?.location_info ||
    "Местоположение уточняется",
}));

export const todayVisitors = hallVisits.map((visit) => ({
  id: visit.id,
  memberId: visit.reader_id,
  memberName:
    readers.find((r) => r.id === visit.reader_id)?.full_name ||
    "Неизвестный читатель",
  entryTime: visit.visit_type === "entry" ? visit.visit_time : new Date(),
  exitTime: visit.visit_type === "exit" ? visit.visit_time : undefined,
  status: visit.visit_type === "entry" ? ("in" as const) : ("out" as const),
}));

// Hourly statistics (реалистичная статистика посещений)
export interface HourlyStats {
  hour: number;
  visitors: number;
}

export const hourlyStats: HourlyStats[] = [
  { hour: 8, visitors: 12 },
  { hour: 9, visitors: 28 },
  { hour: 10, visitors: 45 },
  { hour: 11, visitors: 67 },
  { hour: 12, visitors: 89 },
  { hour: 13, visitors: 98 },
  { hour: 14, visitors: 112 },
  { hour: 15, visitors: 128 },
  { hour: 16, visitors: 118 },
  { hour: 17, visitors: 95 },
  { hour: 18, visitors: 76 },
  { hour: 19, visitors: 52 },
  { hour: 20, visitors: 24 },
];

export const recentEntryRecords = todayVisitors.slice(0, 8);
