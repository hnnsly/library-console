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
} from "lucide-react";
import {
  books,
  checkoutRecords,
  reservations,
  members,
} from "../data/mockData";

const BookDetails: React.FC = () => {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();

  const book = books.find((book) => book.id === id);

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

  // Get checkout and reservation info
  const activeCheckouts = checkoutRecords.filter(
    (record) => record.bookId === id && record.status === "active",
  );

  const activeReservations = reservations.filter(
    (reservation) =>
      reservation.bookId === id && reservation.status === "pending",
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
        <div className="md:flex">
          <motion.div variants={itemVariants} className="md:w-1/3 p-6">
            <div className="aspect-[2/3] overflow-hidden rounded-lg shadow-md">
              <img
                src={book.coverImage}
                alt={book.title}
                className="w-full h-full object-cover"
              />
            </div>

            <div className="mt-6 space-y-4">
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
                <div className="space-y-2">
                  <button
                    disabled={book.availableCopies === 0}
                    className={`btn w-full flex items-center justify-center ${
                      book.availableCopies > 0
                        ? "bg-primary-500 text-white hover:bg-primary-600"
                        : "bg-gray-200 text-gray-500 cursor-not-allowed"
                    }`}
                  >
                    <BookOpen size={16} className="mr-2" /> Выдать
                  </button>

                  <button className="btn btn-accent w-full flex items-center justify-center">
                    <Calendar size={16} className="mr-2" /> Забронировать
                  </button>
                </div>
              </div>
            </div>
          </motion.div>

          <motion.div
            variants={itemVariants}
            className="md:w-2/3 p-6 md:border-l border-gray-200"
          >
            <div className="flex flex-col md:flex-row md:items-center justify-between mb-4">
              <h1 className="text-3xl font-serif font-bold">{book.title}</h1>
              <div className="mt-2 md:mt-0">
                {book.status === "available" ? (
                  <span className="inline-flex items-center px-3 py-1 rounded-full text-sm font-medium bg-green-100 text-green-800">
                    Доступно
                  </span>
                ) : (
                  <span className="inline-flex items-center px-3 py-1 rounded-full text-sm font-medium bg-red-100 text-red-800">
                    {getStatusText(book.status)}
                  </span>
                )}
              </div>
            </div>

            <div className="text-xl text-gray-600 mb-6">
              Автор: {book.author}
            </div>

            <div className="grid grid-cols-1 md:grid-cols-2 gap-4 mb-6">
              <div className="flex items-start">
                <div className="p-2 rounded-full bg-gray-100 text-gray-600">
                  <BookOpen size={18} />
                </div>
                <div className="ml-3">
                  <h3 className="text-sm text-gray-500">ISBN</h3>
                  <p className="font-medium">{book.isbn}</p>
                </div>
              </div>

              <div className="flex items-start">
                <div className="p-2 rounded-full bg-gray-100 text-gray-600">
                  <Calendar size={18} />
                </div>
                <div className="ml-3">
                  <h3 className="text-sm text-gray-500">Издано</h3>
                  <p className="font-medium">
                    {book.publishedYear}, {book.publisher}
                  </p>
                </div>
              </div>

              <div className="flex items-start">
                <div className="p-2 rounded-full bg-gray-100 text-gray-600">
                  <Map size={18} />
                </div>
                <div className="ml-3">
                  <h3 className="text-sm text-gray-500">Местоположение</h3>
                  <p className="font-medium">{book.location}</p>
                </div>
              </div>

              <div className="flex items-start">
                <div className="p-2 rounded-full bg-gray-100 text-gray-600">
                  <AlertTriangle size={18} />
                </div>
                <div className="ml-3">
                  <h3 className="text-sm text-gray-500">Наличие</h3>
                  <p className="font-medium">
                    {book.availableCopies} из {book.totalCopies} экземпляров
                    доступно
                  </p>
                </div>
              </div>
            </div>

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

            <div className="mb-6">
              <h3 className="text-lg font-medium mb-2">Описание</h3>
              <p className="text-gray-700 leading-relaxed">
                {book.description}
              </p>
            </div>

            {/* Current checkouts */}
            {activeCheckouts.length > 0 && (
              <div className="mb-6">
                <h3 className="text-lg font-medium mb-2">Текущие выдачи</h3>
                <div className="bg-blue-50 rounded-lg p-4">
                  <table className="min-w-full">
                    <thead>
                      <tr>
                        <th className="text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                          Читатель
                        </th>
                        <th className="text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                          Дата возврата
                        </th>
                      </tr>
                    </thead>
                    <tbody>
                      {activeCheckouts.map((checkout) => {
                        const member = members.find(
                          (m) => m.id === checkout.memberId,
                        );
                        return (
                          <tr key={checkout.id}>
                            <td className="py-2">
                              <Link
                                to={`/members/${member?.id}`}
                                className="text-primary-600 hover:underline"
                              >
                                {member?.firstName} {member?.lastName}
                              </Link>
                            </td>
                            <td className="py-2">
                              {new Date(checkout.dueDate).toLocaleDateString(
                                "ru-RU",
                              )}
                            </td>
                          </tr>
                        );
                      })}
                    </tbody>
                  </table>
                </div>
              </div>
            )}

            {/* Reservations */}
            {activeReservations.length > 0 && (
              <div>
                <h3 className="text-lg font-medium mb-2">Бронирования</h3>
                <div className="bg-purple-50 rounded-lg p-4">
                  <table className="min-w-full">
                    <thead>
                      <tr>
                        <th className="text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                          Читатель
                        </th>
                        <th className="text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                          Истекает
                        </th>
                      </tr>
                    </thead>
                    <tbody>
                      {activeReservations.map((reservation) => {
                        const member = members.find(
                          (m) => m.id === reservation.memberId,
                        );
                        return (
                          <tr key={reservation.id}>
                            <td className="py-2">
                              <Link
                                to={`/members/${member?.id}`}
                                className="text-primary-600 hover:underline"
                              >
                                {member?.firstName} {member?.lastName}
                              </Link>
                            </td>
                            <td className="py-2">
                              {new Date(
                                reservation.expirationDate,
                              ).toLocaleDateString("ru-RU")}
                            </td>
                          </tr>
                        );
                      })}
                    </tbody>
                  </table>
                </div>
              </div>
            )}
          </motion.div>
        </div>
      </motion.div>
    </div>
  );
};

export default BookDetails;
