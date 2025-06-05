import React, { useState } from "react";
import { Link } from "react-router-dom";
import { motion } from "framer-motion";
import { Search, PlusCircle, Check, X } from "lucide-react";
import { books } from "../data/mockData";
import { Book } from "../types";

const Books: React.FC = () => {
  const [searchTerm, setSearchTerm] = useState("");
  const [selectedGenre, setSelectedGenre] = useState<string>("");
  const [selectedStatus, setSelectedStatus] = useState<string>("");
  const [viewMode, setViewMode] = useState<"grid" | "list">("grid");

  // Get unique genres from all books
  const allGenres = [...new Set(books.flatMap((book) => book.genre))].sort();

  // Filter books based on search term and filters
  const filteredBooks = books.filter((book) => {
    const matchesSearch =
      book.title.toLowerCase().includes(searchTerm.toLowerCase()) ||
      book.author.toLowerCase().includes(searchTerm.toLowerCase()) ||
      book.isbn.includes(searchTerm);

    const matchesGenre = selectedGenre
      ? book.genre.includes(selectedGenre)
      : true;
    const matchesStatus = selectedStatus
      ? book.status === selectedStatus
      : true;

    return matchesSearch && matchesGenre && matchesStatus;
  });

  const handleSearchChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setSearchTerm(e.target.value);
  };

  const clearFilters = () => {
    setSelectedGenre("");
    setSelectedStatus("");
  };

  // const cardVariants = {
  //   hidden: { opacity: 0, y: 20 },
  //   visible: {
  //     opacity: 1,
  //     y: 0,
  //     transition: { type: "spring", stiffness: 300, damping: 24 },
  //   },
  //   hover: {
  //     y: -5,
  //     boxShadow:
  //       "0 10px 15px -3px rgba(0, 0, 0, 0.1), 0 4px 6px -2px rgba(0, 0, 0, 0.05)",
  //     transition: { type: "spring", stiffness: 300, damping: 24 },
  //   },
  // };

  const containerVariants = {
    hidden: { opacity: 0 },
    visible: {
      opacity: 1,
      transition: {
        staggerChildren: 0.05,
      },
    },
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
              value={selectedGenre}
              onChange={(e) => setSelectedGenre(e.target.value)}
              className="input py-2"
            >
              <option value="">Все жанры</option>
              {allGenres.map((genre) => (
                <option key={genre} value={genre}>
                  {genre}
                </option>
              ))}
            </select>

            <select
              value={selectedStatus}
              onChange={(e) => setSelectedStatus(e.target.value)}
              className="input py-2"
            >
              <option value="">Все статусы</option>
              <option value="available">Доступно</option>
              <option value="checked-out">Выдано</option>
              <option value="reserved">Забронировано</option>
              <option value="lost">Утеряно</option>
              <option value="damaged">Повреждено</option>
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
        {(selectedGenre || selectedStatus) && (
          <div className="mt-3 flex items-center">
            <span className="text-sm text-gray-600 mr-2">
              Активные фильтры:
            </span>
            <div className="flex flex-wrap gap-2">
              {selectedGenre && (
                <div className="flex items-center bg-primary-100 text-primary-800 text-sm rounded-full px-3 py-1">
                  <span>Жанр: {selectedGenre}</span>
                  <button onClick={() => setSelectedGenre("")} className="ml-2">
                    <X size={14} />
                  </button>
                </div>
              )}
              {selectedStatus && (
                <div className="flex items-center bg-primary-100 text-primary-800 text-sm rounded-full px-3 py-1">
                  <span>Статус: {getStatusText(selectedStatus)}</span>
                  <button
                    onClick={() => setSelectedStatus("")}
                    className="ml-2"
                  >
                    <X size={14} />
                  </button>
                </div>
              )}
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
                      Автор
                    </th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                      Жанр
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
                      <td className="px-6 py-4 whitespace-nowrap">
                        <Link
                          to={`/books/${book.id}`}
                          className="text-primary-600 hover:underline font-medium"
                        >
                          {book.title}
                        </Link>
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap text-gray-700">
                        {book.author}
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap text-gray-700">
                        {book.genre[0]}
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap">
                        <StatusBadge status={book.status} />
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap text-gray-700">
                        {book.availableCopies}/{book.totalCopies}
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
    "checked-out": "Выдано",
    reserved: "Забронировано",
    lost: "Утеряно",
    damaged: "Повреждено",
  };
  return statusMap[status] || status;
};

// Book Card Component
const BookCard: React.FC<{ book: Book }> = ({ book }) => {
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
      <div className="relative h-48 overflow-hidden">
        <img
          src={book.coverImage}
          alt={book.title}
          className="w-full h-full object-cover transition-transform duration-300 hover:scale-105"
        />
        <div className="absolute top-2 right-2">
          <StatusBadge status={book.status} />
        </div>
      </div>
      <div className="p-4 flex flex-col flex-grow">
        <Link to={`/books/${book.id}`} className="hover:underline">
          <h3 className="font-medium text-lg line-clamp-1">{book.title}</h3>
        </Link>
        <p className="text-gray-600 mb-2">{book.author}</p>
        <div className="flex flex-wrap gap-1 mb-3">
          {book.genre.slice(0, 2).map((genre, index) => (
            <span
              key={index}
              className="text-xs bg-gray-100 text-gray-800 px-2 py-1 rounded-full"
            >
              {genre}
            </span>
          ))}
          {book.genre.length > 2 && (
            <span className="text-xs bg-gray-100 text-gray-800 px-2 py-1 rounded-full">
              +{book.genre.length - 2}
            </span>
          )}
        </div>
        <div className="mt-auto flex items-center justify-between">
          <span className="text-sm text-gray-600">
            {book.availableCopies}/{book.totalCopies} экз.
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
const StatusBadge: React.FC<{ status: string }> = ({ status }) => {
  let bgColor, textColor, text;

  switch (status) {
    case "available":
      bgColor = "bg-green-100";
      textColor = "text-green-800";
      text = "Доступно";
      break;
    case "checked-out":
      bgColor = "bg-blue-100";
      textColor = "text-blue-800";
      text = "Выдано";
      break;
    case "reserved":
      bgColor = "bg-purple-100";
      textColor = "text-purple-800";
      text = "Забронировано";
      break;
    case "lost":
      bgColor = "bg-red-100";
      textColor = "text-red-800";
      text = "Утеряно";
      break;
    case "damaged":
      bgColor = "bg-orange-100";
      textColor = "text-orange-800";
      text = "Повреждено";
      break;
    default:
      bgColor = "bg-gray-100";
      textColor = "text-gray-800";
      text = status;
  }

  return (
    <span
      className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${bgColor} ${textColor}`}
    >
      {status === "available" && <Check size={12} className="mr-1" />}
      {text}
    </span>
  );
};

export default Books;
