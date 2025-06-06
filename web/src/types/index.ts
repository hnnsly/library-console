// Database enums
export type BookStatus =
  | "available"
  | "issued"
  | "reserved"
  | "lost"
  | "damaged";
export type UserRole = "administrator" | "librarian";
export type VisitType = "entry" | "exit";

// User types
export interface User {
  id: string;
  username: string;
  email: string;
  password_hash: string;
  role: UserRole;
  is_active: boolean;
  created_at: Date;
}

// Reader types
export interface Reader {
  id: string;
  ticket_number: string;
  full_name: string;
  email?: string;
  phone?: string;
  registration_date: Date;
  is_active: boolean;
  created_at: Date;
}

// Author types
export interface Author {
  id: string;
  full_name: string;
  created_at: Date;
}

// Book types
export interface Book {
  id: string;
  title: string;
  isbn?: string;
  publication_year?: number;
  publisher?: string;
  total_copies: number;
  available_copies: number;
  created_at: Date;
  // Additional fields from joins
  authors?: Author[];
  copies?: BookCopy[];
}

export interface BookAuthor {
  book_id: string;
  author_id: string;
}

export interface BookCopy {
  id: string;
  book_id: string;
  copy_code: string;
  status: BookStatus;
  hall_id?: string;
  location_info?: string;
  created_at: Date;
}

// Book issue types (аналог checkout)
export interface BookIssue {
  id: string;
  reader_id: string;
  book_copy_id: string;
  issue_date: Date;
  due_date: Date;
  return_date?: Date;
  librarian_id?: string;
  created_at: Date;
  // Additional fields from joins
  reader?: Reader;
  book_copy?: BookCopy;
  librarian?: User;
}

// Fine types
export interface Fine {
  id: string;
  reader_id: string;
  book_issue_id?: string;
  amount: number; // decimal в TypeScript будем представлять как number
  reason: string;
  fine_date: Date;
  paid_date?: Date;
  is_paid: boolean;
  created_at: Date;
  // Additional fields from joins
  reader?: Reader;
  book_issue?: BookIssue;
}

// Reading hall types
export interface ReadingHall {
  id: string;
  hall_name: string;
  specialization?: string;
  total_seats: number;
  current_visitors: number;
  created_at: Date;
}

// Hall visit types
export interface HallVisit {
  id: string;
  reader_id: string;
  hall_id: string;
  visit_type: VisitType;
  visit_time: Date;
  librarian_id?: string;
  // Additional fields from joins
  reader?: Reader;
  hall?: ReadingHall;
  librarian?: User;
}

// Extended types for UI
export interface BookWithDetails extends Book {
  authors: Author[];
  active_issues?: BookIssue[];
  total_issues_count?: number;
}

export interface ReaderWithDetails extends Reader {
  active_issues?: BookIssue[];
  unpaid_fines?: Fine[];
  total_fines_amount?: number;
  current_hall_visits?: HallVisit[];
}

export interface BookIssueWithDetails extends BookIssue {
  reader: Reader;
  book_copy: BookCopy & {
    book: Book & {
      authors: Author[];
    };
  };
  librarian?: User;
  is_overdue?: boolean;
  days_overdue?: number;
}

// Dashboard statistics
export interface DashboardStats {
  total_books: number;
  total_readers: number;
  active_issues: number;
  overdue_issues: number;
  total_fines: number;
  current_hall_visitors: number;
}

// Search and filter types
export interface BookFilter {
  title?: string;
  author?: string;
  isbn?: string;
  status?: BookStatus;
  hall_id?: string;
}

export interface ReaderFilter {
  full_name?: string;
  email?: string;
  ticket_number?: string;
  is_active?: boolean;
}

export interface IssueFilter {
  reader_id?: string;
  book_id?: string;
  is_returned?: boolean;
  is_overdue?: boolean;
  librarian_id?: string;
}

// API response types
export interface PaginatedResponse<T> {
  data: T[];
  total: number;
  page: number;
  per_page: number;
  total_pages: number;
}

export interface ApiResponse<T> {
  success: boolean;
  data?: T;
  error?: string;
  message?: string;
}

// Form types for creating/updating
export interface CreateBookRequest {
  title: string;
  isbn?: string;
  publication_year?: number;
  publisher?: string;
  total_copies: number;
  author_ids: string[];
}

export interface UpdateBookRequest extends Partial<CreateBookRequest> {
  id: string;
}

export interface CreateReaderRequest {
  ticket_number: string;
  full_name: string;
  email?: string;
  phone?: string;
}

export interface UpdateReaderRequest extends Partial<CreateReaderRequest> {
  id: string;
  is_active?: boolean;
}

export interface CreateBookIssueRequest {
  reader_id: string;
  book_copy_id: string;
  due_date: Date;
}

export interface ReturnBookRequest {
  book_issue_id: string;
  condition?: "good" | "damaged";
  notes?: string;
}

export interface CreateUserRequest {
  username: string;
  email: string;
  password: string;
  role: UserRole;
}

export interface UpdateUserRequest {
  id: string;
  username?: string;
  email?: string;
  role?: UserRole;
  is_active?: boolean;
}

export interface RegisterHallVisitRequest {
  reader_id: string;
  hall_id: string;
  visit_type: VisitType;
}

// Notification types (если нужны)
export interface Notification {
  id: string;
  user_id: string;
  message: string;
  type: "info" | "warning" | "error" | "success";
  is_read: boolean;
  created_at: Date;
}
