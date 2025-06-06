import React from "react";
import { useParams, Link, useNavigate } from "react-router-dom";
import { motion } from "framer-motion";
import {
  ArrowLeft,
  Edit,
  Trash2,
  Mail,
  Phone,
  Calendar,
  AlertTriangle,
  Clock,
  Check,
  CreditCard,
  BookOpen,
} from "lucide-react";
import {
  readers,
  bookIssues,
  booksWithDetails,
  bookCopies,
  fines,
} from "../data/mockData";
import type { ReaderWithDetails, BookIssueWithDetails } from "../types";

const ReaderDetails: React.FC = () => {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();

  const reader = readers.find((reader) => reader.id === id);

  if (!reader) {
    return (
      <div className="max-w-4xl mx-auto text-center py-12">
        <h1 className="text-2xl font-bold text-gray-900 mb-4">
          Читатель не найден
        </h1>
        <p className="text-gray-600 mb-6">
          Читатель, которого вы ищете, не существует или был удален из системы
          libr.
        </p>
        <button
          onClick={() => navigate("/readers")}
          className="btn btn-primary"
        >
          Вернуться к читателям
        </button>
      </div>
    );
  }

  // Получаем текущие выдачи читателя
  const currentIssues = bookIssues.filter(
    (issue) => issue.reader_id === reader.id && !issue.return_date,
  );

  // Получаем историю выдач читателя
  const issueHistory = bookIssues.filter(
    (issue) => issue.reader_id === reader.id && issue.return_date,
  );

  // Получаем штрафы читателя
  const readerFines = fines.filter((fine) => fine.reader_id === reader.id);
  const unpaidFines = readerFines.filter((fine) => !fine.is_paid);
  const totalUnpaidFines = unpaidFines.reduce(
    (sum, fine) => sum + fine.amount,
    0,
  );

  // Функция для получения деталей книги по экземпляру
  const getBookDetailsByIssue = (issue: any) => {
    const bookCopy = bookCopies.find((copy) => copy.id === issue.book_copy_id);
    const book = booksWithDetails.find((book) => book.id === bookCopy?.book_id);
    return { book, bookCopy };
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
    <div className="max-w-7xl mx-auto">
      <div className="mb-6">
        <button
          onClick={() => navigate("/readers")}
          className="flex items-center text-primary-600 hover:underline"
        >
          <ArrowLeft size={16} className="mr-1" /> Назад к читателям
        </button>
      </div>

      <motion.div
        variants={containerVariants}
        initial="hidden"
        animate="visible"
      >
        <motion.div
          variants={itemVariants}
          className="bg-white rounded-lg shadow-sm overflow-hidden mb-6"
        >
          <div
            className={`h-2 ${
              reader.is_active ? "bg-green-500" : "bg-red-500"
            }`}
          ></div>

          <div className="p-6">
            <div className="flex flex-col md:flex-row md:items-center justify-between mb-6">
              <div className="flex items-center">
                <div className="w-16 h-16 rounded-full bg-primary-500 text-white flex items-center justify-center text-xl font-medium">
                  {reader.full_name
                    .split(" ")
                    .map((name) => name.charAt(0))
                    .join("")
                    .slice(0, 2)}
                </div>
                <div className="ml-4">
                  <div className="flex items-center">
                    <h1 className="text-3xl font-serif font-bold">
                      {reader.full_name}
                    </h1>
                    {reader.is_active ? (
                      <span className="ml-3 inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-green-100 text-green-800">
                        Активен
                      </span>
                    ) : (
                      <span className="ml-3 inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-red-100 text-red-800">
                        Неактивен
                      </span>
                    )}
                  </div>
                  <div className="mt-1 space-y-1">
                    <p className="text-gray-600">
                      Билет: {reader.ticket_number}
                    </p>
                    <p className="text-gray-600">ID: {reader.id}</p>
                  </div>
                </div>
              </div>
              <div className="mt-4 md:mt-0 flex space-x-3">
                <Link
                  to={`/readers/${reader.id}/edit`}
                  className="btn btn-secondary flex items-center"
                >
                  <Edit size={16} className="mr-2" /> Редактировать
                </Link>
                <button className="btn bg-red-500 text-white hover:bg-red-600 flex items-center">
                  <Trash2 size={16} className="mr-2" /> Удалить
                </button>
              </div>
            </div>

            <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
              {reader.email && (
                <div className="flex items-start">
                  <div className="p-2 rounded-full bg-gray-100 text-gray-600">
                    <Mail size={18} />
                  </div>
                  <div className="ml-3">
                    <h3 className="text-sm text-gray-500">Эл. почта</h3>
                    <p className="font-medium">{reader.email}</p>
                  </div>
                </div>
              )}

              {reader.phone && (
                <div className="flex items-start">
                  <div className="p-2 rounded-full bg-gray-100 text-gray-600">
                    <Phone size={18} />
                  </div>
                  <div className="ml-3">
                    <h3 className="text-sm text-gray-500">Телефон</h3>
                    <p className="font-medium">{reader.phone}</p>
                  </div>
                </div>
              )}

              <div className="flex items-start">
                <div className="p-2 rounded-full bg-gray-100 text-gray-600">
                  <Calendar size={18} />
                </div>
                <div className="ml-3">
                  <h3 className="text-sm text-gray-500">Дата регистрации</h3>
                  <p className="font-medium">
                    {new Date(reader.registration_date).toLocaleDateString(
                      "ru-RU",
                    )}
                  </p>
                </div>
              </div>
            </div>

            {/* Статистика читателя */}
            <div className="mt-6 grid grid-cols-1 md:grid-cols-4 gap-4">
              <div className="bg-blue-50 rounded-lg p-4 text-center">
                <div className="text-2xl font-bold text-blue-600">
                  {currentIssues.length}
                </div>
                <div className="text-sm text-blue-600">Книг на руках</div>
              </div>
              <div className="bg-green-50 rounded-lg p-4 text-center">
                <div className="text-2xl font-bold text-green-600">
                  {issueHistory.length}
                </div>
                <div className="text-sm text-green-600">Всего выдач</div>
              </div>
              <div className="bg-purple-50 rounded-lg p-4 text-center">
                <div className="text-2xl font-bold text-purple-600">
                  {readerFines.length}
                </div>
                <div className="text-sm text-purple-600">Штрафов</div>
              </div>
              <div className="bg-orange-50 rounded-lg p-4 text-center">
                <div className="text-2xl font-bold text-orange-600">
                  {totalUnpaidFines.toFixed(0)} ₽
                </div>
                <div className="text-sm text-orange-600">К доплате</div>
              </div>
            </div>

            {totalUnpaidFines > 0 && (
              <div className="mt-6 bg-red-50 border border-red-100 rounded-lg p-4 flex items-center">
                <AlertTriangle size={20} className="text-red-500 mr-3" />
                <div className="flex-1">
                  <h3 className="font-medium text-red-800">
                    Неоплаченные штрафы
                  </h3>
                  <p className="text-red-700">
                    У данного читателя есть неоплаченные штрафы на сумму{" "}
                    {totalUnpaidFines.toFixed(2)} ₽.
                  </p>
                </div>
                <button className="btn bg-white text-red-600 border border-red-200 hover:bg-red-50 flex items-center">
                  <CreditCard size={16} className="mr-2" />
                  Принять оплату
                </button>
              </div>
            )}
          </div>
        </motion.div>

        {/* Текущие выдачи */}
        <motion.div
          variants={itemVariants}
          className="bg-white rounded-lg shadow-sm overflow-hidden mb-6"
        >
          <div className="p-4 bg-primary-500 text-white font-medium flex items-center">
            <BookOpen size={20} className="mr-2" />
            Книги на руках ({currentIssues.length})
          </div>
          <div className="p-6">
            {currentIssues.length > 0 ? (
              <div className="overflow-x-auto">
                <table className="min-w-full divide-y divide-gray-200">
                  <thead>
                    <tr>
                      <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                        Книга
                      </th>
                      <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                        Экземпляр
                      </th>
                      <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                        Дата выдачи
                      </th>
                      <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                        Срок возврата
                      </th>
                      <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                        Статус
                      </th>
                      <th className="px-4 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">
                        Действия
                      </th>
                    </tr>
                  </thead>
                  <tbody className="divide-y divide-gray-200">
                    {currentIssues.map((issue) => {
                      const { book, bookCopy } = getBookDetailsByIssue(issue);
                      const isOverdue = new Date(issue.due_date) < new Date();

                      return (
                        <tr key={issue.id}>
                          <td className="px-4 py-4">
                            <div>
                              <Link
                                to={`/books/${book?.id}`}
                                className="text-primary-600 hover:underline font-medium"
                              >
                                {book?.title}
                              </Link>
                              <div className="text-sm text-gray-500">
                                {book?.authors
                                  .map((author) => author.full_name)
                                  .join(", ")}
                              </div>
                            </div>
                          </td>
                          <td className="px-4 py-4 whitespace-nowrap">
                            <div className="text-sm font-medium text-gray-900">
                              {bookCopy?.copy_code}
                            </div>
                            {bookCopy?.location_info && (
                              <div className="text-sm text-gray-500">
                                {bookCopy.location_info}
                              </div>
                            )}
                          </td>
                          <td className="px-4 py-4 whitespace-nowrap text-gray-700">
                            {new Date(issue.issue_date).toLocaleDateString(
                              "ru-RU",
                            )}
                          </td>
                          <td className="px-4 py-4 whitespace-nowrap text-gray-700">
                            {new Date(issue.due_date).toLocaleDateString(
                              "ru-RU",
                            )}
                          </td>
                          <td className="px-4 py-4 whitespace-nowrap">
                            {isOverdue ? (
                              <span className="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-red-100 text-red-800">
                                <AlertTriangle size={12} className="mr-1" />
                                Просрочено
                              </span>
                            ) : (
                              <span className="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-green-100 text-green-800">
                                <Check size={12} className="mr-1" />В срок
                              </span>
                            )}
                          </td>
                          <td className="px-4 py-4 whitespace-nowrap text-right space-x-2">
                            <button className="text-primary-600 hover:text-primary-900 text-sm">
                              Продлить
                            </button>
                            <button className="text-green-600 hover:text-green-900 text-sm">
                              Вернуть
                            </button>
                          </td>
                        </tr>
                      );
                    })}
                  </tbody>
                </table>
              </div>
            ) : (
              <p className="text-gray-500 text-center py-6">
                У данного читателя нет книг на руках.
              </p>
            )}
          </div>
        </motion.div>

        {/* История выдач */}
        <motion.div
          variants={itemVariants}
          className="bg-white rounded-lg shadow-sm overflow-hidden mb-6"
        >
          <div className="p-4 bg-accent-500 text-white font-medium flex items-center">
            <Clock size={20} className="mr-2" />
            История выдач ({issueHistory.length})
          </div>
          <div className="p-6">
            {issueHistory.length > 0 ? (
              <div className="overflow-x-auto">
                <table className="min-w-full divide-y divide-gray-200">
                  <thead>
                    <tr>
                      <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                        Книга
                      </th>
                      <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                        Экземпляр
                      </th>
                      <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                        Выдана
                      </th>
                      <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                        Возвращена
                      </th>
                      <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                        Статус
                      </th>
                    </tr>
                  </thead>
                  <tbody className="divide-y divide-gray-200">
                    {issueHistory
                      .sort(
                        (a, b) =>
                          new Date(b.return_date!).getTime() -
                          new Date(a.return_date!).getTime(),
                      )
                      .map((issue) => {
                        const { book, bookCopy } = getBookDetailsByIssue(issue);
                        const wasOverdue =
                          new Date(issue.return_date!) >
                          new Date(issue.due_date);

                        return (
                          <tr key={issue.id}>
                            <td className="px-4 py-4">
                              <div>
                                <Link
                                  to={`/books/${book?.id}`}
                                  className="text-primary-600 hover:underline font-medium"
                                >
                                  {book?.title}
                                </Link>
                                <div className="text-sm text-gray-500">
                                  {book?.authors
                                    .map((author) => author.full_name)
                                    .join(", ")}
                                </div>
                              </div>
                            </td>
                            <td className="px-4 py-4 whitespace-nowrap text-sm text-gray-700">
                              {bookCopy?.copy_code}
                            </td>
                            <td className="px-4 py-4 whitespace-nowrap text-gray-700">
                              {new Date(issue.issue_date).toLocaleDateString(
                                "ru-RU",
                              )}
                            </td>
                            <td className="px-4 py-4 whitespace-nowrap text-gray-700">
                              {new Date(issue.return_date!).toLocaleDateString(
                                "ru-RU",
                              )}
                            </td>
                            <td className="px-4 py-4 whitespace-nowrap">
                              {wasOverdue ? (
                                <span className="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-red-100 text-red-800">
                                  <Clock size={12} className="mr-1" />{" "}
                                  Просрочено
                                </span>
                              ) : (
                                <span className="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-green-100 text-green-800">
                                  <Check size={12} className="mr-1" /> В срок
                                </span>
                              )}
                            </td>
                          </tr>
                        );
                      })}
                  </tbody>
                </table>
              </div>
            ) : (
              <p className="text-gray-500 text-center py-6">
                У данного читателя нет истории выдач.
              </p>
            )}
          </div>
        </motion.div>

        {/* Штрафы */}
        {readerFines.length > 0 && (
          <motion.div
            variants={itemVariants}
            className="bg-white rounded-lg shadow-sm overflow-hidden"
          >
            <div className="p-4 bg-orange-500 text-white font-medium flex items-center">
              <AlertTriangle size={20} className="mr-2" />
              Штрафы ({readerFines.length})
            </div>
            <div className="p-6">
              <div className="overflow-x-auto">
                <table className="min-w-full divide-y divide-gray-200">
                  <thead>
                    <tr>
                      <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                        Сумма
                      </th>
                      <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                        Причина
                      </th>
                      <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                        Дата штрафа
                      </th>
                      <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                        Статус
                      </th>
                      <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                        Дата оплаты
                      </th>
                    </tr>
                  </thead>
                  <tbody className="divide-y divide-gray-200">
                    {readerFines.map((fine) => (
                      <tr key={fine.id}>
                        <td className="px-4 py-4 whitespace-nowrap font-medium text-gray-900">
                          {fine.amount.toFixed(2)} ₽
                        </td>
                        <td className="px-4 py-4 text-gray-700">
                          {fine.reason}
                        </td>
                        <td className="px-4 py-4 whitespace-nowrap text-gray-700">
                          {new Date(fine.fine_date).toLocaleDateString("ru-RU")}
                        </td>
                        <td className="px-4 py-4 whitespace-nowrap">
                          {fine.is_paid ? (
                            <span className="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-green-100 text-green-800">
                              <Check size={12} className="mr-1" />
                              Оплачено
                            </span>
                          ) : (
                            <span className="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-red-100 text-red-800">
                              <AlertTriangle size={12} className="mr-1" />
                              Не оплачено
                            </span>
                          )}
                        </td>
                        <td className="px-4 py-4 whitespace-nowrap text-gray-700">
                          {fine.paid_date
                            ? new Date(fine.paid_date).toLocaleDateString(
                                "ru-RU",
                              )
                            : "—"}
                        </td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>
            </div>
          </motion.div>
        )}
      </motion.div>
    </div>
  );
};

export default ReaderDetails;
