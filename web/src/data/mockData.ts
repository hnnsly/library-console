import { Book, Member, CheckoutRecord, Reservation } from "../types";
import { addDays, subDays } from "date-fns";

// Generate mock books
export const books: Book[] = [
  {
    id: "1",
    title: "To Kill a Mockingbird",
    author: "Harper Lee",
    isbn: "9780060935467",
    coverImage:
      "https://images.pexels.com/photos/3646172/pexels-photo-3646172.jpeg?auto=compress&cs=tinysrgb&w=1260&h=750&dpr=2",
    publisher: "HarperCollins",
    publishedYear: 1960,
    genre: ["Fiction", "Classic", "Coming-of-age"],
    description:
      "To Kill a Mockingbird is a novel by Harper Lee published in 1960. It was immediately successful, winning the Pulitzer Prize, and has become a classic of modern American literature.",
    totalCopies: 5,
    availableCopies: 3,
    location: "Fiction - F LEE",
    addedDate: new Date("2020-01-15"),
    status: "available",
  },
  {
    id: "2",
    title: "1984",
    author: "George Orwell",
    isbn: "9780451524935",
    coverImage:
      "https://images.pexels.com/photos/6373305/pexels-photo-6373305.jpeg?auto=compress&cs=tinysrgb&w=1260&h=750&dpr=2",
    publisher: "Signet Classic",
    publishedYear: 1949,
    genre: ["Fiction", "Dystopian", "Political"],
    description:
      "1984 is a dystopian novel by English novelist George Orwell. It was published in June 1949 as Orwell's ninth and final book completed in his lifetime.",
    totalCopies: 4,
    availableCopies: 1,
    location: "Fiction - F ORW",
    addedDate: new Date("2020-02-10"),
    status: "checked-out",
  },
  {
    id: "3",
    title: "The Great Gatsby",
    author: "F. Scott Fitzgerald",
    isbn: "9780743273565",
    coverImage:
      "https://images.pexels.com/photos/1907785/pexels-photo-1907785.jpeg?auto=compress&cs=tinysrgb&w=1260&h=750&dpr=2",
    publisher: "Scribner",
    publishedYear: 1925,
    genre: ["Fiction", "Classic", "Literary Fiction"],
    description:
      "The Great Gatsby is a 1925 novel by American writer F. Scott Fitzgerald. Set in the Jazz Age on Long Island, the novel depicts narrator Nick Carraway's interactions with mysterious millionaire Jay Gatsby.",
    totalCopies: 3,
    availableCopies: 2,
    location: "Fiction - F FIT",
    addedDate: new Date("2020-03-05"),
    status: "available",
  },
  {
    id: "4",
    title: "Pride and Prejudice",
    author: "Jane Austen",
    isbn: "9780141439518",
    coverImage:
      "https://images.pexels.com/photos/1765033/pexels-photo-1765033.jpeg?auto=compress&cs=tinysrgb&w=1260&h=750&dpr=2",
    publisher: "Penguin Classics",
    publishedYear: 1813,
    genre: ["Fiction", "Classic", "Romance"],
    description:
      "Pride and Prejudice is an 1813 romantic novel of manners written by Jane Austen. The novel follows the character development of Elizabeth Bennet, the dynamic protagonist of the book.",
    totalCopies: 5,
    availableCopies: 4,
    location: "Fiction - F AUS",
    addedDate: new Date("2020-04-12"),
    status: "available",
  },
  {
    id: "5",
    title: "The Hobbit",
    author: "J.R.R. Tolkien",
    isbn: "9780547928227",
    coverImage:
      "https://images.pexels.com/photos/3377538/pexels-photo-3377538.jpeg?auto=compress&cs=tinysrgb&w=1260&h=750&dpr=2",
    publisher: "Houghton Mifflin Harcourt",
    publishedYear: 1937,
    genre: ["Fiction", "Fantasy", "Adventure"],
    description:
      "The Hobbit, or There and Back Again is a children's fantasy novel by J. R. R. Tolkien. It was published on 21 September 1937 to wide critical acclaim.",
    totalCopies: 6,
    availableCopies: 3,
    location: "Fantasy - F TOL",
    addedDate: new Date("2020-05-20"),
    status: "checked-out",
  },
  {
    id: "6",
    title: "The Catcher in the Rye",
    author: "J.D. Salinger",
    isbn: "9780316769488",
    coverImage:
      "https://images.pexels.com/photos/5834344/pexels-photo-5834344.jpeg?auto=compress&cs=tinysrgb&w=1260&h=750&dpr=2",
    publisher: "Little, Brown and Company",
    publishedYear: 1951,
    genre: ["Fiction", "Coming-of-age"],
    description:
      "The Catcher in the Rye is a novel by J. D. Salinger, partially published in serial form in 1945–1946 and as a novel in 1951. It was originally intended for adults but is often read by adolescents for its themes of teenage angst and alienation.",
    totalCopies: 4,
    availableCopies: 2,
    location: "Fiction - F SAL",
    addedDate: new Date("2020-06-15"),
    status: "available",
  },
];

// Generate mock members
export const members: Member[] = [
  {
    id: "1",
    firstName: "Иван",
    lastName: "Петров",
    email: "ivan.petrov@example.com",
    phone: "555-123-4567",
    address: "123 Главная ул., Москва",
    joinDate: new Date("2020-01-10"),
    membershipStatus: "active",
    borrowedBooks: [
      {
        bookId: "2",
        borrowDate: subDays(new Date(), 15),
        dueDate: addDays(new Date(), 15),
        renewCount: 0,
      },
    ],
    borrowingHistory: [
      {
        bookId: "1",
        borrowDate: subDays(new Date(), 60),
        returnDate: subDays(new Date(), 45),
        wasLate: false,
        condition: "good",
      },
      {
        bookId: "3",
        borrowDate: subDays(new Date(), 30),
        returnDate: subDays(new Date(), 10),
        wasLate: true,
        condition: "good",
      },
    ],
    fines: 0,
  },
  {
    id: "2",
    firstName: "Мария",
    lastName: "Сидорова",
    email: "maria.sidorova@example.com",
    phone: "555-987-6543",
    address: "456 Дубовая ул., Санкт-Петербург",
    joinDate: new Date("2020-03-15"),
    membershipStatus: "active",
    borrowedBooks: [
      {
        bookId: "5",
        borrowDate: subDays(new Date(), 10),
        dueDate: addDays(new Date(), 4),
        renewCount: 1,
      },
    ],
    borrowingHistory: [
      {
        bookId: "4",
        borrowDate: subDays(new Date(), 90),
        returnDate: subDays(new Date(), 75),
        wasLate: false,
        condition: "good",
      },
    ],
    fines: 0,
  },
  {
    id: "3",
    firstName: "Александр",
    lastName: "Козлов",
    email: "alexander.kozlov@example.com",
    phone: "555-456-7890",
    address: "789 Сосновая ул., Новосибирск",
    joinDate: new Date("2020-05-22"),
    membershipStatus: "expired",
    borrowedBooks: [],
    borrowingHistory: [
      {
        bookId: "1",
        borrowDate: subDays(new Date(), 120),
        returnDate: subDays(new Date(), 100),
        wasLate: false,
        condition: "good",
      },
      {
        bookId: "6",
        borrowDate: subDays(new Date(), 80),
        returnDate: subDays(new Date(), 50),
        wasLate: true,
        condition: "damaged",
      },
    ],
    fines: 15.5,
  },
  {
    id: "4",
    firstName: "Елена",
    lastName: "Волкова",
    email: "elena.volkova@example.com",
    phone: "555-321-9876",
    address: "321 Березовая ул., Екатеринбург",
    joinDate: new Date("2020-07-10"),
    membershipStatus: "active",
    borrowedBooks: [],
    borrowingHistory: [
      {
        bookId: "2",
        borrowDate: subDays(new Date(), 40),
        returnDate: subDays(new Date(), 20),
        wasLate: false,
        condition: "good",
      },
    ],
    fines: 0,
  },
  {
    id: "5",
    firstName: "Дмитрий",
    lastName: "Орлов",
    email: "dmitry.orlov@example.com",
    phone: "555-654-3210",
    address: "654 Кленовая ул., Казань",
    joinDate: new Date("2020-09-05"),
    membershipStatus: "active",
    borrowedBooks: [],
    borrowingHistory: [
      {
        bookId: "3",
        borrowDate: subDays(new Date(), 25),
        returnDate: subDays(new Date(), 5),
        wasLate: false,
        condition: "good",
      },
    ],
    fines: 0,
  },
];

// Generate mock checkout records
export const checkoutRecords: CheckoutRecord[] = [
  {
    id: "1",
    memberId: "1",
    bookId: "2",
    checkoutDate: subDays(new Date(), 15),
    dueDate: addDays(new Date(), 15),
    returnedDate: null,
    status: "active",
    fineAmount: 0,
  },
  {
    id: "2",
    memberId: "2",
    bookId: "5",
    checkoutDate: subDays(new Date(), 10),
    dueDate: addDays(new Date(), 4),
    returnedDate: null,
    status: "active",
    fineAmount: 0,
  },
  {
    id: "3",
    memberId: "1",
    bookId: "1",
    checkoutDate: subDays(new Date(), 60),
    dueDate: subDays(new Date(), 30),
    returnedDate: subDays(new Date(), 45),
    status: "returned",
    fineAmount: 0,
  },
  {
    id: "4",
    memberId: "1",
    bookId: "3",
    checkoutDate: subDays(new Date(), 30),
    dueDate: subDays(new Date(), 0),
    returnedDate: subDays(new Date(), 10),
    status: "returned",
    fineAmount: 5.0,
  },
  {
    id: "5",
    memberId: "3",
    bookId: "6",
    checkoutDate: subDays(new Date(), 80),
    dueDate: subDays(new Date(), 50),
    returnedDate: subDays(new Date(), 50),
    status: "returned",
    fineAmount: 15.5,
  },
];

// Generate mock reservations
export const reservations: Reservation[] = [
  {
    id: "1",
    memberId: "3",
    bookId: "2",
    reservationDate: subDays(new Date(), 5),
    expirationDate: addDays(new Date(), 25),
    status: "pending",
  },
  {
    id: "2",
    memberId: "1",
    bookId: "5",
    reservationDate: subDays(new Date(), 3),
    expirationDate: addDays(new Date(), 27),
    status: "pending",
  },
];

// Reading room data
export interface RoomData {
  id: string;
  name: string;
  capacity: number;
  currentOccupancy: number;
}

export const rooms: RoomData[] = [
  {
    id: "1",
    name: "Основной читальный зал",
    capacity: 150,
    currentOccupancy: 99,
  },
  {
    id: "2",
    name: "Тихий зал",
    capacity: 50,
    currentOccupancy: 2,
  },
  {
    id: "3",
    name: "Компьютерный зал",
    capacity: 30,
    currentOccupancy: 18,
  },
  {
    id: "4",
    name: "Групповые занятия",
    capacity: 40,
    currentOccupancy: 16,
  },
];

// Visitor records for room entry/exit
export interface VisitorRecord {
  id: string;
  memberId: string;
  memberName: string;
  entryTime: Date;
  exitTime?: Date;
  status: "in" | "out";
}

export const todayVisitors: VisitorRecord[] = [
  {
    id: "1",
    memberId: "1",
    memberName: "Иван Петров",
    entryTime: new Date(new Date().setHours(9, 30, 0, 0)),
    exitTime: new Date(new Date().setHours(14, 45, 0, 0)),
    status: "out",
  },
  {
    id: "2",
    memberId: "2",
    memberName: "Мария Сидорова",
    entryTime: new Date(new Date().setHours(10, 15, 0, 0)),
    status: "in",
  },
  {
    id: "3",
    memberId: "3",
    memberName: "Александр Козлов",
    entryTime: new Date(new Date().setHours(11, 0, 0, 0)),
    exitTime: new Date(new Date().setHours(15, 30, 0, 0)),
    status: "out",
  },
  {
    id: "4",
    memberId: "4",
    memberName: "Елена Волкова",
    entryTime: new Date(new Date().setHours(12, 20, 0, 0)),
    status: "in",
  },
  {
    id: "5",
    memberId: "5",
    memberName: "Дмитрий Орлов",
    entryTime: new Date(new Date().setHours(13, 10, 0, 0)),
    exitTime: new Date(new Date().setHours(16, 0, 0, 0)),
    status: "out",
  },
  {
    id: "6",
    memberId: "1",
    memberName: "Иван Петров",
    entryTime: new Date(new Date().setHours(8, 45, 0, 0)),
    exitTime: new Date(new Date().setHours(12, 15, 0, 0)),
    status: "out",
  },
  {
    id: "7",
    memberId: "2",
    memberName: "Мария Сидорова",
    entryTime: new Date(new Date().setHours(14, 30, 0, 0)),
    status: "in",
  },
];

// Hourly statistics for visitor chart
export interface HourlyStats {
  hour: number;
  visitors: number;
}

export const hourlyStats: HourlyStats[] = [
  { hour: 8, visitors: 8 },
  { hour: 9, visitors: 12 },
  { hour: 10, visitors: 25 },
  { hour: 11, visitors: 38 },
  { hour: 12, visitors: 45 },
  { hour: 13, visitors: 52 },
  { hour: 14, visitors: 48 },
  { hour: 15, visitors: 41 },
  { hour: 16, visitors: 35 },
  { hour: 17, visitors: 28 },
  { hour: 18, visitors: 18 },
  { hour: 19, visitors: 12 },
  { hour: 20, visitors: 6 },
];

// Recent entry records for RoomEntry page
export interface EntryRecord {
  id: string;
  memberId: string;
  memberName: string;
  entryTime: Date;
  exitTime?: Date;
  status: "in" | "out";
}

export const recentEntryRecords: EntryRecord[] = [
  {
    id: "1",
    memberId: "1",
    memberName: "Иван Петров",
    entryTime: new Date(Date.now() - 2 * 60 * 60 * 1000), // 2 часа назад
    status: "in",
  },
  {
    id: "2",
    memberId: "2",
    memberName: "Мария Сидорова",
    entryTime: new Date(Date.now() - 3 * 60 * 60 * 1000), // 3 часа назад
    exitTime: new Date(Date.now() - 1 * 60 * 60 * 1000), // 1 час назад
    status: "out",
  },
  {
    id: "3",
    memberId: "4",
    memberName: "Елена Волкова",
    entryTime: new Date(Date.now() - 4 * 60 * 60 * 1000), // 4 часа назад
    status: "in",
  },
  {
    id: "4",
    memberId: "5",
    memberName: "Дмитрий Орлов",
    entryTime: new Date(Date.now() - 5 * 60 * 60 * 1000), // 5 часов назад
    exitTime: new Date(Date.now() - 30 * 60 * 1000), // 30 минут назад
    status: "out",
  },
  {
    id: "5",
    memberId: "3",
    memberName: "Александр Козлов",
    entryTime: new Date(Date.now() - 6 * 60 * 60 * 1000), // 6 часов назад
    exitTime: new Date(Date.now() - 2 * 60 * 60 * 1000), // 2 часа назад
    status: "out",
  },
];
