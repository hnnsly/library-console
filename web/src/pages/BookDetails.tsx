import React from "react";
import { useParams, Link, useNavigate } from "react-router-dom";
import { motion } from "framer-motion";
import {
  ArrowLeft,
  Edit,
  Trash2,
  BookOpen,
  Calendar,
  AlertTriangle,
  Map,
  User,
  Building2,
} from "lucide-react";
import {
  booksWithDetails,
  bookIssues,
  readers,
  bookCopies,
  users,
} from "../data/mockData";
import type { BookStatus } from "../types";

const BookDetails: React.FC = () => {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();

  const book = booksWithDetails.find((book) => book.id === id);

  if (!book) {
    return (
      <div className="max-w-4xl mx-auto text-center py-12">
        <h1 className="text-2xl font-bold text-gray-900 mb-4">
          Книга не найдена
        </h1>
        <p className="text-gray-600 mb-6">
          Книга, которую вы ищете, не существует или была удалена из системы
          libr.
        </p>
        <button onClick={() => navigate("/books")} className="btn btn-primary">
          Вернуться к книгам
        </button>
      </div>
    );
  }

  // Get book copies for this book
  const bookCopiesForBook = bookCopies.filter(
    (copy) => copy.book_id === book.id,
  );

  // Get current issues for this book
  const activeIssues = bookIssues.filter(
    (issue) =>
      bookCopiesForBook.some((copy) => copy.id === issue.book_copy_id) &&
      !issue.return_date,
  );

  // Animation variants
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

  const getStatusText = (status: BookStatus): string => {
    const statusMap: Record<BookStatus, string> = {
      available: "Доступно",
      issued: "Выдано",
      reserved: "Забронировано",
      lost: "Утеряно",
      damaged: "Повреждено",
    };
    return statusMap[status];
  };

  const getStatusColor = (availableCopies: number, totalCopies: number) => {
    if (availableCopies === 0) {
      return "bg-red-100 text-red-800";
    } else if (availableCopies < totalCopies / 2) {
      return "bg-yellow-100 text-yellow-800";
    }
    return "bg-green-100 text-green-800";
  };

  const getStatusLabel = (availableCopies: number) => {
    if (availableCopies === 0) return "Недоступно";
    return "Доступно";
  };

  return (
    <div className="max-w-7xl mx-auto">
      <div className="mb-6">
        <button
          onClick={() => navigate("/books")}
          className="flex items-center text-primary-600 hover:underline"
        >
          <ArrowLeft size={16} className="mr-1" /> Вернуться к книгам
        </button>
      </div>

      <motion.div
        variants={containerVariants}
        initial="hidden"
        animate="visible"
        className="bg-white rounded-lg shadow-sm overflow-hidden"
      >
        <motion.div variants={itemVariants} className="p-6">
          <div className="flex flex-col md:flex-row md:items-center justify-between mb-4">
            <h1 className="text-3xl font-serif font-bold">{book.title}</h1>
            <div className="mt-2 md:mt-0">
              <span
                className={`inline-flex items-center px-3 py-1 rounded-full text-sm font-medium ${getStatusColor(
                  book.available_copies,
                  book.total_copies,
                )}`}
              >
                {getStatusLabel(book.available_copies)}
              </span>
            </div>
          </div>

          <div className="text-xl text-gray-600 mb-6">
            Автор{book.authors.length > 1 ? "ы" : ""}:{" "}
            {book.authors.map((author) => author.full_name).join(", ")}
          </div>

          <div className="mb-6 space-y-4">
            <div className="flex justify-between">
              <Link
                to={`/books/${book.id}/edit`}
                className="btn btn-secondary flex items-center"
              >
                <Edit size={16} className="mr-2" /> Изменить
              </Link>
              <button className="btn bg-red-500 text-white hover:bg-red-600 flex items-center">
                <Trash2 size={16} className="mr-2" /> Удалить
              </button>
            </div>

            <div className="border-t border-gray-200 pt-4">
              <h3 className="text-lg font-medium mb-3">Быстрые действия</h3>
              <div className="flex space-x-2">
                <button
                  disabled={book.available_copies === 0}
                  className={`btn flex items-center justify-center ${
                    book.available_copies > 0
                      ? "bg-primary-500 text-white hover:bg-primary-600"
                      : "bg-gray-200 text-gray-500 cursor-not-allowed"
                  }`}
                >
                  <BookOpen size={16} className="mr-2" /> Выдать
                </button>

                <button className="btn btn-accent flex items-center justify-center">
                  <Calendar size={16} className="mr-2" /> Забронировать
                </button>
              </div>
            </div>
          </div>

          <div className="grid grid-cols-1 md:grid-cols-2 gap-4 mb-6">
            {book.isbn && (
              <div className="flex items-start">
                <div className="p-2 rounded-full bg-gray-100 text-gray-600">
                  <BookOpen size={18} />
                </div>
                <div className="ml-3">
                  <h3 className="text-sm text-gray-500">ISBN</h3>
                  <p className="font-medium">{book.isbn}</p>
                </div>
              </div>
            )}

            {(book.publication_year || book.publisher) && (
              <div className="flex items-start">
                <div className="p-2 rounded-full bg-gray-100 text-gray-600">
                  <Calendar size={18} />
                </div>
                <div className="ml-3">
                  <h3 className="text-sm text-gray-500">Издано</h3>
                  <p className="font-medium">
                    {book.publication_year && book.publisher
                      ? `${book.publication_year}, ${book.publisher}`
                      : book.publication_year || book.publisher}
                  </p>
                </div>
              </div>
            )}

            <div className="flex items-start">
              <div className="p-2 rounded-full bg-gray-100 text-gray-600">
                <Map size={18} />
              </div>
              <div className="ml-3">
                <h3 className="text-sm text-gray-500">Местоположение</h3>
                <p className="font-medium">
                  {bookCopiesForBook[0]?.location_info || "Не указано"}
                </p>
              </div>
            </div>

            <div className="flex items-start">
              <div className="p-2 rounded-full bg-gray-100 text-gray-600">
                <AlertTriangle size={18} />
              </div>
              <div className="ml-3">
                <h3 className="text-sm text-gray-500">Наличие</h3>
                <p className="font-medium">
                  {book.available_copies} из {book.total_copies} экземпляров
                  доступно
                </p>
              </div>
            </div>
          </div>

          {book.genre && book.genre.length > 0 && (
            <div className="mb-6">
              <h3 className="text-lg font-medium mb-2">Жанры</h3>
              <div className="flex flex-wrap gap-2">
                {book.genre.map((genre, index) => (
                  <span
                    key={index}
                    className="bg-primary-100 text-primary-800 px-3 py-1 rounded-full text-sm"
                  >
                    {genre}
                  </span>
                ))}
              </div>
            </div>
          )}

          {book.description && (
            <div className="mb-6">
              <h3 className="text-lg font-medium mb-2">Описание</h3>
              <p className="text-gray-700 leading-relaxed">
                {book.description}
              </p>
            </div>
          )}

          {/* Book copies section */}
          {bookCopiesForBook.length > 0 && (
            <div className="mb-6">
              <h3 className="text-lg font-medium mb-2">Экземпляры книги</h3>
              <div className="bg-gray-50 rounded-lg p-4">
                <div className="grid gap-3">
                  {bookCopiesForBook.map((copy) => (
                    <div
                      key={copy.id}
                      className="flex items-center justify-between bg-white p-3 rounded border"
                    >
                      <div className="flex items-center space-x-3">
                        <Building2 size={16} className="text-gray-500" />
                        <div>
                          <p className="font-medium">{copy.copy_code}</p>
                          <p className="text-sm text-gray-500">
                            {copy.location_info}
                          </p>
                        </div>
                      </div>
                      <span
                        className={`px-2 py-1 rounded-full text-xs font-medium ${
                          copy.status === "available"
                            ? "bg-green-100 text-green-800"
                            : copy.status === "issued"
                              ? "bg-blue-100 text-blue-800"
                              : copy.status === "reserved"
                                ? "bg-purple-100 text-purple-800"
                                : "bg-red-100 text-red-800"
                        }`}
                      >
                        {getStatusText(copy.status)}
                      </span>
                    </div>
                  ))}
                </div>
              </div>
            </div>
          )}

          {/* Current issues */}
          {activeIssues.length > 0 && (
            <div className="mb-6">
              <h3 className="text-lg font-medium mb-2">Текущие выдачи</h3>
              <div className="bg-blue-50 rounded-lg p-4">
                <div className="space-y-3">
                  {activeIssues.map((issue) => {
                    const reader = readers.find(
                      (r) => r.id === issue.reader_id,
                    );
                    const copy = bookCopiesForBook.find(
                      (c) => c.id === issue.book_copy_id,
                    );
                    const librarian = users.find(
                      (u) => u.id === issue.librarian_id,
                    );
                    const isOverdue = new Date(issue.due_date) < new Date();

                    return (
                      <div
                        key={issue.id}
                        className="flex items-center justify-between bg-white p-3 rounded border"
                      >
                        <div className="flex items-center space-x-3">
                          <User size={16} className="text-gray-500" />
                          <div>
                            <p className="font-medium">
                              <Link
                                to={`/readers/${reader?.id}`}
                                className="text-primary-600 hover:underline"
                              >
                                {reader?.full_name}
                              </Link>
                            </p>
                            <p className="text-sm text-gray-500">
                              Экземпляр: {copy?.copy_code}
                            </p>
                            {librarian && (
                              <p className="text-xs text-gray-400">
                                Выдал: {librarian.username}
                              </p>
                            )}
                          </div>
                        </div>
                        <div className="text-right">
                          <p
                            className={`text-sm font-medium ${isOverdue ? "text-red-600" : "text-gray-900"}`}
                          >
                            {isOverdue ? "Просрочено" : "Возврат до:"}
                          </p>
                          <p
                            className={`text-sm ${isOverdue ? "text-red-600" : "text-gray-600"}`}
                          >
                            {new Date(issue.due_date).toLocaleDateString(
                              "ru-RU",
                            )}
                          </p>
                        </div>
                      </div>
                    );
                  })}
                </div>
              </div>
            </div>
          )}

          {/* Additional info */}
          <div className="border-t border-gray-200 pt-4">
            <div className="grid grid-cols-2 md:grid-cols-4 gap-4 text-center">
              <div>
                <p className="text-2xl font-bold text-primary-600">
                  {book.total_copies}
                </p>
                <p className="text-sm text-gray-500">Всего экземпляров</p>
              </div>
              <div>
                <p className="text-2xl font-bold text-green-600">
                  {book.available_copies}
                </p>
                <p className="text-sm text-gray-500">Доступно</p>
              </div>
              <div>
                <p className="text-2xl font-bold text-blue-600">
                  {activeIssues.length}
                </p>
                <p className="text-sm text-gray-500">Выдано</p>
              </div>
              <div>
                <p className="text-2xl font-bold text-gray-600">
                  {new Date(book.created_at).toLocaleDateString("ru-RU")}
                </p>
                <p className="text-sm text-gray-500">Добавлено</p>
              </div>
            </div>
          </div>
        </motion.div>
      </motion.div>
    </div>
  );
};

export default BookDetails;
