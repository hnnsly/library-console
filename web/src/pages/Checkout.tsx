import React, { useState } from "react";
import { motion } from "framer-motion";
import {
  Search,
  BookOpen,
  User,
  ArrowRight,
  Calendar,
  Check,
  Book,
} from "lucide-react";
import { booksWithDetails, readers, bookCopies } from "../data/mockData";
import { addDays, format } from "date-fns";
import { ru } from "date-fns/locale";
import type { BookWithDetails, Reader, BookCopy } from "../types";

interface CheckoutFormData {
  book: BookWithDetails | null;
  reader: Reader | null;
  bookCopy: BookCopy | null;
  dueDate: string;
}

const Checkout: React.FC = () => {
  const [bookSearchTerm, setBookSearchTerm] = useState("");
  const [readerSearchTerm, setReaderSearchTerm] = useState("");
  const [formData, setFormData] = useState<CheckoutFormData>({
    book: null,
    reader: null,
    bookCopy: null,
    dueDate: format(addDays(new Date(), 14), "yyyy-MM-dd"),
  });
  const [checkoutComplete, setCheckoutComplete] = useState(false);

  // Фильтрация книг по поисковому запросу (только те, у которых есть доступные экземпляры)
  const filteredBooks = booksWithDetails
    .filter((book) => {
      const hasAvailableCopies = book.available_copies > 0;
      const matchesSearch =
        book.title.toLowerCase().includes(bookSearchTerm.toLowerCase()) ||
        book.authors.some((author) =>
          author.full_name.toLowerCase().includes(bookSearchTerm.toLowerCase()),
        ) ||
        (book.isbn && book.isbn.includes(bookSearchTerm));

      return hasAvailableCopies && matchesSearch;
    })
    .slice(0, 5);

  // Фильтрация читателей по поисковому запросу (только активные)
  const filteredReaders = readers
    .filter((reader) => {
      const isActive = reader.is_active;
      const matchesSearch =
        reader.full_name
          .toLowerCase()
          .includes(readerSearchTerm.toLowerCase()) ||
        (reader.email &&
          reader.email
            .toLowerCase()
            .includes(readerSearchTerm.toLowerCase())) ||
        (reader.phone && reader.phone.includes(readerSearchTerm)) ||
        reader.ticket_number.includes(readerSearchTerm);

      return isActive && matchesSearch;
    })
    .slice(0, 5);

  // Получить доступные экземпляры для выбранной книги
  const getAvailableCopiesForBook = (bookId: string): BookCopy[] => {
    return bookCopies.filter(
      (copy) => copy.book_id === bookId && copy.status === "available",
    );
  };

  const handleBookSelect = (book: BookWithDetails) => {
    const availableCopies = getAvailableCopiesForBook(book.id);
    const firstAvailableCopy = availableCopies[0] || null;

    setFormData((prev) => ({
      ...prev,
      book,
      bookCopy: firstAvailableCopy,
    }));
  };

  const handleReaderSelect = (reader: Reader) => {
    setFormData((prev) => ({
      ...prev,
      reader,
    }));
  };

  const handleCheckout = () => {
    if (formData.book && formData.reader && formData.bookCopy) {
      // В реальном приложении это обновило бы базу данных
      setCheckoutComplete(true);

      // Сбросить форму через 3 секунды
      setTimeout(() => {
        setFormData({
          book: null,
          reader: null,
          bookCopy: null,
          dueDate: format(addDays(new Date(), 14), "yyyy-MM-dd"),
        });
        setBookSearchTerm("");
        setReaderSearchTerm("");
        setCheckoutComplete(false);
      }, 3000);
    }
  };

  const updateDueDate = (date: string) => {
    setFormData((prev) => ({ ...prev, dueDate: date }));
  };

  const clearBookSelection = () => {
    setFormData((prev) => ({ ...prev, book: null, bookCopy: null }));
  };

  const clearReaderSelection = () => {
    setFormData((prev) => ({ ...prev, reader: null }));
  };

  // Варианты анимации
  const containerVariants = {
    hidden: { opacity: 0 },
    visible: {
      opacity: 1,
      transition: {
        staggerChildren: 0.1,
      },
    },
  };

  const itemVariants = {
    hidden: { y: 20, opacity: 0 },
    visible: {
      y: 0,
      opacity: 1,
      transition: { type: "spring", stiffness: 300, damping: 24 },
    },
  };

  return (
    <div className="max-w-5xl mx-auto">
      <div className="mb-6">
        <h1 className="text-3xl font-serif font-bold text-gray-900">
          Выдача книг
        </h1>
        <p className="text-gray-600 mt-1">
          Выдача книг читателям библиотеки libr.
        </p>
      </div>

      <motion.div
        variants={containerVariants}
        initial="hidden"
        animate="visible"
        className="bg-white rounded-lg shadow-sm overflow-hidden"
      >
        {checkoutComplete ? (
          <motion.div
            initial={{ opacity: 0, scale: 0.9 }}
            animate={{ opacity: 1, scale: 1 }}
            className="p-8 text-center"
          >
            <div className="w-16 h-16 bg-green-100 rounded-full flex items-center justify-center mx-auto mb-4">
              <Check size={32} className="text-green-600" />
            </div>
            <h2 className="text-2xl font-bold text-gray-900 mb-2">
              Выдача завершена!
            </h2>
            <p className="text-gray-600 mb-6">
              <span className="font-medium">«{formData.book?.title}»</span>{" "}
              выдана читателю{" "}
              <span className="font-medium">{formData.reader?.full_name}</span>
            </p>
            <div className="text-gray-700 space-y-1">
              <p>
                Экземпляр:{" "}
                <span className="font-medium">
                  {formData.bookCopy?.copy_code}
                </span>
              </p>
              <p>
                Билет читателя:{" "}
                <span className="font-medium">
                  {formData.reader?.ticket_number}
                </span>
              </p>
              <p>
                Дата возврата:{" "}
                <span className="font-medium">
                  {format(new Date(formData.dueDate), "d MMMM yyyy", {
                    locale: ru,
                  })}
                </span>
              </p>
            </div>
          </motion.div>
        ) : (
          <div className="grid grid-cols-1 md:grid-cols-2 divide-y md:divide-y-0 md:divide-x divide-gray-200">
            {/* Левая колонка: Выбор книги */}
            <motion.div variants={itemVariants} className="p-6">
              <h2 className="text-xl font-semibold mb-4 flex items-center">
                <BookOpen size={20} className="mr-2 text-primary-500" />
                Выбор книги
              </h2>

              <div className="mb-4">
                <label htmlFor="bookSearch" className="label">
                  Поиск книги по названию, автору или ISBN
                </label>
                <div className="relative">
                  <Search
                    className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400"
                    size={18}
                  />
                  <input
                    id="bookSearch"
                    type="text"
                    placeholder="Начните вводить..."
                    value={bookSearchTerm}
                    onChange={(e) => setBookSearchTerm(e.target.value)}
                    className="input pl-10 py-2"
                  />
                </div>
              </div>

              {bookSearchTerm && (
                <div className="mb-4 border rounded-lg overflow-hidden">
                  {filteredBooks.length > 0 ? (
                    <ul className="divide-y divide-gray-200">
                      {filteredBooks.map((book) => (
                        <li
                          key={book.id}
                          className={`p-3 cursor-pointer hover:bg-gray-50 transition-colors ${
                            formData.book?.id === book.id ? "bg-primary-50" : ""
                          }`}
                          onClick={() => handleBookSelect(book)}
                        >
                          <div className="flex items-start">
                            <div className="h-12 w-8 bg-gradient-to-br from-primary-50 to-primary-100 rounded overflow-hidden mr-3 flex-shrink-0 flex items-center justify-center">
                              <Book size={16} className="text-primary-400" />
                            </div>
                            <div>
                              <p className="font-medium text-gray-900">
                                {book.title}
                              </p>
                              <p className="text-sm text-gray-600">
                                {book.authors
                                  .map((author) => author.full_name)
                                  .join(", ")}
                              </p>
                              <div className="flex items-center mt-1">
                                <span className="text-xs bg-green-100 text-green-800 px-2 py-0.5 rounded-full">
                                  {book.available_copies} доступно
                                </span>
                              </div>
                            </div>
                          </div>
                        </li>
                      ))}
                    </ul>
                  ) : (
                    <div className="p-4 text-center text-gray-500">
                      Доступных книг не найдено. Попробуйте другой поисковый
                      запрос.
                    </div>
                  )}
                </div>
              )}

              {formData.book && (
                <div className="mt-6 p-4 border border-primary-100 rounded-lg bg-primary-50">
                  <h3 className="font-medium text-primary-900 mb-2">
                    Выбранная книга
                  </h3>
                  <div className="flex items-start">
                    <div className="h-24 w-16 bg-gradient-to-br from-primary-100 to-primary-200 rounded overflow-hidden mr-4 flex-shrink-0 flex items-center justify-center">
                      <Book size={24} className="text-primary-500" />
                    </div>
                    <div>
                      <p className="font-medium text-gray-900">
                        {formData.book.title}
                      </p>
                      <p className="text-sm text-gray-600 mb-1">
                        {formData.book.authors
                          .map((author) => author.full_name)
                          .join(", ")}
                      </p>
                      {formData.book.isbn && (
                        <p className="text-sm text-gray-600">
                          ISBN: {formData.book.isbn}
                        </p>
                      )}
                      {formData.bookCopy && (
                        <div className="mt-2">
                          <p className="text-sm text-gray-600">
                            Экземпляр:{" "}
                            <span className="font-medium">
                              {formData.bookCopy.copy_code}
                            </span>
                          </p>
                          {formData.bookCopy.location_info && (
                            <p className="text-sm text-gray-600">
                              Расположение: {formData.bookCopy.location_info}
                            </p>
                          )}
                        </div>
                      )}
                      <button
                        onClick={clearBookSelection}
                        className="text-sm text-red-600 hover:underline mt-2"
                      >
                        Убрать
                      </button>
                    </div>
                  </div>
                </div>
              )}
            </motion.div>

            {/* Правая колонка: Выбор читателя */}
            <motion.div variants={itemVariants} className="p-6">
              <h2 className="text-xl font-semibold mb-4 flex items-center">
                <User size={20} className="mr-2 text-primary-500" />
                Выбор читателя
              </h2>

              <div className="mb-4">
                <label htmlFor="readerSearch" className="label">
                  Поиск читателя по имени, email, телефону или номеру билета
                </label>
                <div className="relative">
                  <Search
                    className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400"
                    size={18}
                  />
                  <input
                    id="readerSearch"
                    type="text"
                    placeholder="Начните вводить..."
                    value={readerSearchTerm}
                    onChange={(e) => setReaderSearchTerm(e.target.value)}
                    className="input pl-10 py-2"
                  />
                </div>
              </div>

              {readerSearchTerm && (
                <div className="mb-4 border rounded-lg overflow-hidden">
                  {filteredReaders.length > 0 ? (
                    <ul className="divide-y divide-gray-200">
                      {filteredReaders.map((reader) => (
                        <li
                          key={reader.id}
                          className={`p-3 cursor-pointer hover:bg-gray-50 transition-colors ${
                            formData.reader?.id === reader.id
                              ? "bg-primary-50"
                              : ""
                          }`}
                          onClick={() => handleReaderSelect(reader)}
                        >
                          <div className="flex items-center">
                            <div className="w-10 h-10 rounded-full bg-primary-500 text-white flex items-center justify-center mr-3 text-sm font-medium">
                              {reader.full_name
                                .split(" ")
                                .map((name) => name.charAt(0))
                                .join("")
                                .slice(0, 2)}
                            </div>
                            <div>
                              <p className="font-medium text-gray-900">
                                {reader.full_name}
                              </p>
                              <p className="text-sm text-gray-600">
                                Билет: {reader.ticket_number}
                              </p>
                              {reader.email && (
                                <p className="text-sm text-gray-600">
                                  {reader.email}
                                </p>
                              )}
                            </div>
                          </div>
                        </li>
                      ))}
                    </ul>
                  ) : (
                    <div className="p-4 text-center text-gray-500">
                      Активных читателей не найдено. Попробуйте другой поисковый
                      запрос.
                    </div>
                  )}
                </div>
              )}

              {formData.reader && (
                <div className="mt-6 p-4 border border-primary-100 rounded-lg bg-primary-50">
                  <h3 className="font-medium text-primary-900 mb-2">
                    Выбранный читатель
                  </h3>
                  <div className="flex items-center">
                    <div className="w-12 h-12 rounded-full bg-primary-500 text-white flex items-center justify-center mr-4 font-medium">
                      {formData.reader.full_name
                        .split(" ")
                        .map((name) => name.charAt(0))
                        .join("")
                        .slice(0, 2)}
                    </div>
                    <div>
                      <p className="font-medium text-gray-900">
                        {formData.reader.full_name}
                      </p>
                      <p className="text-sm text-gray-600">
                        Билет: {formData.reader.ticket_number}
                      </p>
                      {formData.reader.email && (
                        <p className="text-sm text-gray-600">
                          {formData.reader.email}
                        </p>
                      )}
                      {formData.reader.phone && (
                        <p className="text-xs text-gray-500">
                          {formData.reader.phone}
                        </p>
                      )}
                      <button
                        onClick={clearReaderSelection}
                        className="text-sm text-red-600 hover:underline mt-2"
                      >
                        Убрать
                      </button>
                    </div>
                  </div>
                </div>
              )}

              <div className="mt-6">
                <label htmlFor="dueDate" className="label flex items-center">
                  <Calendar size={18} className="mr-2" /> Дата возврата
                </label>
                <input
                  id="dueDate"
                  type="date"
                  value={formData.dueDate}
                  onChange={(e) => updateDueDate(e.target.value)}
                  min={format(addDays(new Date(), 1), "yyyy-MM-dd")}
                  className="input py-2"
                />
              </div>
            </motion.div>
          </div>
        )}

        {!checkoutComplete && (
          <div className="bg-gray-50 p-6 flex justify-end">
            <button
              className={`btn flex items-center ${
                formData.book && formData.reader && formData.bookCopy
                  ? "bg-primary-500 text-white hover:bg-primary-600"
                  : "bg-gray-200 text-gray-500 cursor-not-allowed"
              }`}
              disabled={
                !formData.book || !formData.reader || !formData.bookCopy
              }
              onClick={handleCheckout}
            >
              Завершить выдачу <ArrowRight size={18} className="ml-2" />
            </button>
          </div>
        )}
      </motion.div>
    </div>
  );
};

export default Checkout;
