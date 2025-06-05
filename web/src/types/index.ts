// Book types
export interface Book {
  id: string;
  title: string;
  author: string;
  isbn: string;
  coverImage: string;
  publisher: string;
  publishedYear: number;
  genre: string[];
  description: string;
  totalCopies: number;
  availableCopies: number;
  location: string;
  addedDate: Date;
  status: 'available' | 'checked-out' | 'reserved' | 'lost' | 'damaged';
}

// Member types
export interface Member {
  id: string;
  firstName: string;
  lastName: string;
  email: string;
  phone: string;
  address: string;
  joinDate: Date;
  membershipStatus: 'active' | 'expired' | 'suspended';
  borrowedBooks: BorrowedBook[];
  borrowingHistory: BorrowingRecord[];
  fines: number;
}

export interface BorrowedBook {
  bookId: string;
  borrowDate: Date;
  dueDate: Date;
  renewCount: number;
}

export interface BorrowingRecord {
  bookId: string;
  borrowDate: Date;
  returnDate: Date;
  wasLate: boolean;
  condition: 'good' | 'damaged';
}

// Checkout/Return types
export interface CheckoutRecord {
  id: string;
  memberId: string;
  bookId: string;
  checkoutDate: Date;
  dueDate: Date;
  returnedDate: Date | null;
  status: 'active' | 'returned' | 'overdue';
  fineAmount: number;
}

// Reservation types
export interface Reservation {
  id: string;
  memberId: string;
  bookId: string;
  reservationDate: Date;
  expirationDate: Date;
  status: 'pending' | 'fulfilled' | 'expired' | 'cancelled';
}

// Notification types
export interface Notification {
  id: string;
  recipientId: string;
  message: string;
  type: 'due-date' | 'reservation' | 'fine' | 'system';
  createdAt: Date;
  isRead: boolean;
}