import React, { useState } from "react";
import { Link } from "react-router-dom";
import { motion } from "framer-motion";
import { Search, PlusCircle, Check, X, Book } from "lucide-react";
import { booksWithDetails } from "../data/mockData";
import type { BookWithDetails } from "../types";

const Books: React.FC = () => {
  const [searchTerm, setSearchTerm] = useState("");
  const [selectedStatus, setSelectedStatus] = useState<string>("");
  const [viewMode, setViewMode] = useState<"grid" | "list">("grid");

  // Filter books based on search term and filters
  const filteredBooks = booksWithDetails.filter((book) => {
    const authorNames = book.authors
      .map((author) => author.full_name)
      .join(" ");

    const matchesSearch =
      book.title.toLowerCase().includes(searchTerm.toLowerCase()) ||
      authorNames.toLowerCase().includes(searchTerm.toLowerCase()) ||
      (book.isbn && book.isbn.includes(searchTerm));

    const matchesStatus = selectedStatus
      ? getBookStatus(book) === selectedStatus
      : true;

    return matchesSearch && matchesStatus;
  });

  const handleSearchChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setSearchTerm(e.target.value);
  };

  const clearFilters = () => {
    setSelectedStatus("");
  };

  const containerVariants = {
    hidden: { opacity: 0 },
    visible: {
      opacity: 1,
      transition: {
        staggerChildren: 0.05,
      },
    },
  };

  // Helper function to get book status based on available copies
  const getBookStatus = (book: BookWithDetails): string => {
    if (book.available_copies === 0) return "unavailable";
    return "available";
  };

  return (
    <div className="max-w-7xl mx-auto">
      <div className="flex flex-col md:flex-row md:items-center justify-between mb-6">
        <div>
          <h1 className="text-3xl font-serif font-bold text-gray-900">Книги</h1>
          <p className="text-gray-600 mt-1">
            Просмотр и управление коллекцией библиотеки libr.
          </p>
        </div>
        <Link
          to="/books/add"
          className="btn btn-primary mt-3 md:mt-0 flex items-center justify-center md:justify-start"
        >
          <PlusCircle size={18} className="mr-2" /> Добавить книгу
        </Link>
      </div>

      <div className="bg-white rounded-lg shadow-sm p-4 mb-6">
        <div className="flex flex-col md:flex-row md:items-center gap-4">
          {/* Search */}
          <div className="relative flex-grow">
            <Search
              className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400"
              size={18}
            />
            <input
              type="text"
              placeholder="Поиск по названию, автору или ISBN..."
              value={searchTerm}
              onChange={handleSearchChange}
              className="input pl-10 py-2"
            />
          </div>

          {/* Filters */}
          <div className="flex flex-col sm:flex-row gap-3">
            <select
              value={selectedStatus}
              onChange={(e) => setSelectedStatus(e.target.value)}
              className="input py-2"
            >
              <option value="">Все статусы</option>
              <option value="available">Доступно</option>
              <option value="unavailable">Недоступно</option>
            </select>
          </div>

          {/* View toggle */}
          <div className="flex border rounded-md overflow-hidden">
            <button
              className={`px-3 py-2 ${viewMode === "grid" ? "bg-primary-500 text-white" : "bg-white text-gray-700"}`}
              onClick={() => setViewMode("grid")}
            >
              Сетка
            </button>
            <button
              className={`px-3 py-2 ${viewMode === "list" ? "bg-primary-500 text-white" : "bg-white text-gray-700"}`}
              onClick={() => setViewMode("list")}
            >
              Список
            </button>
          </div>
        </div>

        {/* Active filters */}
        {selectedStatus && (
          <div className="mt-3 flex items-center">
            <span className="text-sm text-gray-600 mr-2">
              Активные фильтры:
            </span>
            <div className="flex flex-wrap gap-2">
              <div className="flex items-center bg-primary-100 text-primary-800 text-sm rounded-full px-3 py-1">
                <span>Статус: {getStatusText(selectedStatus)}</span>
                <button onClick={() => setSelectedStatus("")} className="ml-2">
                  <X size={14} />
                </button>
              </div>
              <button
                onClick={clearFilters}
                className="text-sm text-primary-600 hover:underline"
              >
                Очистить все
              </button>
            </div>
          </div>
        )}
      </div>

      {/* Books Grid/List View */}
      {filteredBooks.length > 0 ? (
        viewMode === "grid" ? (
          <motion.div
            variants={containerVariants}
            initial="hidden"
            animate="visible"
            className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-6"
          >
            {filteredBooks.map((book) => (
              <BookCard key={book.id} book={book} />
            ))}
          </motion.div>
        ) : (
          <div className="bg-white rounded-lg shadow-sm overflow-hidden">
            <div className="overflow-x-auto">
              <table className="min-w-full divide-y divide-gray-200">
                <thead className="bg-gray-50">
                  <tr>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                      Название
                    </th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                      Автор(ы)
                    </th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                      Год издания
                    </th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                      Издательство
                    </th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                      Статус
                    </th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                      Экземпляры
                    </th>
                    <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">
                      Действия
                    </th>
                  </tr>
                </thead>
                <tbody className="bg-white divide-y divide-gray-200">
                  {filteredBooks.map((book) => (
                    <tr key={book.id} className="hover:bg-gray-50">
                      <td className="px-6 py-4">
                        <Link
                          to={`/books/${book.id}`}
                          className="text-primary-600 hover:underline font-medium"
                        >
                          {book.title}
                        </Link>
                        {book.isbn && (
                          <p className="text-xs text-gray-500 mt-1">
                            ISBN: {book.isbn}
                          </p>
                        )}
                      </td>
                      <td className="px-6 py-4 text-gray-700">
                        {book.authors
                          .map((author) => author.full_name)
                          .join(", ")}
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap text-gray-700">
                        {book.publication_year || "—"}
                      </td>
                      <td className="px-6 py-4 text-gray-700">
                        {book.publisher || "—"}
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap">
                        <StatusBadge book={book} />
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap text-gray-700">
                        {book.available_copies}/{book.total_copies}
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
                        <Link
                          to={`/books/${book.id}`}
                          className="text-primary-600 hover:text-primary-900 mr-3"
                        >
                          Просмотр
                        </Link>
                        <Link
                          to={`/books/${book.id}/edit`}
                          className="text-accent-600 hover:text-accent-900"
                        >
                          Изменить
                        </Link>
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          </div>
        )
      ) : (
        <div className="text-center py-12 bg-white rounded-lg shadow-sm">
          <p className="text-gray-500 text-lg">
            Книги не найдены. Попробуйте изменить поисковый запрос или фильтры.
          </p>
          <button
            onClick={clearFilters}
            className="mt-3 text-primary-600 hover:underline"
          >
            Очистить все фильтры
          </button>
        </div>
      )}
    </div>
  );
};

// Helper function to get status text in Russian
const getStatusText = (status: string): string => {
  const statusMap: { [key: string]: string } = {
    available: "Доступно",
    unavailable: "Недоступно",
  };
  return statusMap[status] || status;
};

// Book Card Component
const BookCard: React.FC<{ book: BookWithDetails }> = ({ book }) => {
  const cardVariants = {
    hidden: { opacity: 0, y: 20 },
    visible: {
      opacity: 1,
      y: 0,
      transition: { type: "spring", stiffness: 300, damping: 24 },
    },
    hover: {
      y: -5,
      boxShadow:
        "0 10px 15px -3px rgba(0, 0, 0, 0.1), 0 4px 6px -2px rgba(0, 0, 0, 0.05)",
      transition: { type: "spring", stiffness: 300, damping: 24 },
    },
  };

  return (
    <motion.div
      variants={cardVariants}
      whileHover="hover"
      className="card overflow-hidden flex flex-col h-full"
    >
      <div className="relative h-32 bg-gradient-to-br from-primary-50 to-primary-100 flex items-center justify-center">
        <Book size={48} className="text-primary-400" />
        <div className="absolute top-2 right-2">
          <StatusBadge book={book} />
        </div>
      </div>
      <div className="p-4 flex flex-col flex-grow">
        <Link to={`/books/${book.id}`} className="hover:underline">
          <h3 className="font-medium text-lg line-clamp-2 mb-2">
            {book.title}
          </h3>
        </Link>
        <p className="text-gray-600 mb-2 text-sm">
          {book.authors.map((author) => author.full_name).join(", ")}
        </p>

        <div className="text-xs text-gray-500 space-y-1 mb-3">
          {book.publication_year && <p>Год: {book.publication_year}</p>}
          {book.publisher && <p>Издательство: {book.publisher}</p>}
          {book.isbn && <p>ISBN: {book.isbn}</p>}
        </div>

        <div className="mt-auto flex items-center justify-between">
          <span className="text-sm text-gray-600">
            {book.available_copies}/{book.total_copies} экз.
          </span>
          <Link
            to={`/books/${book.id}`}
            className="text-sm text-primary-600 hover:underline"
          >
            Подробнее
          </Link>
        </div>
      </div>
    </motion.div>
  );
};

// Status Badge Component
const StatusBadge: React.FC<{ book: BookWithDetails }> = ({ book }) => {
  const isAvailable = book.available_copies > 0;

  let bgColor, textColor, text;

  if (isAvailable) {
    bgColor = "bg-green-100";
    textColor = "text-green-800";
    text = "Доступно";
  } else {
    bgColor = "bg-red-100";
    textColor = "text-red-800";
    text = "Недоступно";
  }

  return (
    <span
      className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${bgColor} ${textColor}`}
    >
      {isAvailable && <Check size={12} className="mr-1" />}
      {text}
    </span>
  );
};

export default Books;
