import React from "react";
import { Link } from "react-router-dom";
import { motion } from "framer-motion";
import { BookOpen, Users, Clock, AlertCircle, Building2 } from "lucide-react";
import {
  dashboardStats,
  bookIssues,
  readers,
  booksWithDetails,
  bookCopies,
  readingHalls,
} from "../data/mockData";

const Dashboard: React.FC = () => {
  // Books due this week
  const today = new Date();
  const endOfWeek = new Date(today);
  endOfWeek.setDate(today.getDate() + 7);

  const issuesDueThisWeek = bookIssues.filter(
    (issue) =>
      !issue.return_date &&
      new Date(issue.due_date) >= today &&
      new Date(issue.due_date) <= endOfWeek,
  );

  // Recent activities - using book issues
  const recentActivities = bookIssues
    .sort((a, b) => {
      const dateA = a.return_date || a.issue_date;
      const dateB = b.return_date || b.issue_date;
      return new Date(dateB).getTime() - new Date(dateA).getTime();
    })
    .slice(0, 5);

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

  return (
    <div className="max-w-7xl mx-auto">
      <div className="mb-6">
        <h1 className="text-3xl font-serif font-bold text-gray-900">
          Панель управления
        </h1>
        <p className="text-gray-600 mt-1">
          Добро пожаловать в систему управления библиотекой libr.
        </p>
      </div>

      <motion.div
        variants={containerVariants}
        initial="hidden"
        animate="visible"
        className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8"
      >
        {/* Stat Cards */}
        <motion.div variants={itemVariants} className="card p-6">
          <div className="flex items-center">
            <div className="p-3 rounded-full bg-primary-100 text-primary-600">
              <BookOpen size={24} />
            </div>
            <div className="ml-4">
              <h3 className="text-lg font-medium text-gray-700">Всего книг</h3>
              <p className="text-2xl font-semibold">
                {dashboardStats.total_books}
              </p>
            </div>
          </div>
          <div className="mt-3 text-sm">
            <span className="text-green-600 font-medium">
              {booksWithDetails.reduce(
                (sum, book) => sum + book.available_copies,
                0,
              )}{" "}
              доступно
            </span>
          </div>
        </motion.div>

        <motion.div variants={itemVariants} className="card p-6">
          <div className="flex items-center">
            <div className="p-3 rounded-full bg-accent-100 text-accent-600">
              <Users size={24} />
            </div>
            <div className="ml-4">
              <h3 className="text-lg font-medium text-gray-700">Читатели</h3>
              <p className="text-2xl font-semibold">
                {dashboardStats.total_readers}
              </p>
            </div>
          </div>
          <div className="mt-3 text-sm">
            <span className="text-green-600 font-medium">
              {readers.filter((r) => r.is_active).length} активных
            </span>
          </div>
        </motion.div>

        <motion.div variants={itemVariants} className="card p-6">
          <div className="flex items-center">
            <div className="p-3 rounded-full bg-green-100 text-green-600">
              <BookOpen size={24} />
            </div>
            <div className="ml-4">
              <h3 className="text-lg font-medium text-gray-700">Выдано</h3>
              <p className="text-2xl font-semibold">
                {dashboardStats.active_issues}
              </p>
            </div>
          </div>
          <div className="mt-3 text-sm">
            <Link
              to="/checkout"
              className="text-primary-500 hover:underline font-medium"
            >
              Новая выдача
            </Link>
          </div>
        </motion.div>

        <motion.div variants={itemVariants} className="card p-6">
          <div className="flex items-center">
            <div className="p-3 rounded-full bg-red-100 text-red-600">
              <AlertCircle size={24} />
            </div>
            <div className="ml-4">
              <h3 className="text-lg font-medium text-gray-700">Просрочено</h3>
              <p className="text-2xl font-semibold">
                {dashboardStats.overdue_issues}
              </p>
            </div>
          </div>
          <div className="mt-3 text-sm">
            <Link
              to="/returns"
              className="text-primary-500 hover:underline font-medium"
            >
              Обработать возврат
            </Link>
          </div>
        </motion.div>
      </motion.div>

      {/* Second Row Stats - Reading Halls */}
      <motion.div
        variants={containerVariants}
        initial="hidden"
        animate="visible"
        className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6 mb-8"
      >
        <motion.div variants={itemVariants} className="card p-6">
          <div className="flex items-center">
            <div className="p-3 rounded-full bg-purple-100 text-purple-600">
              <Building2 size={24} />
            </div>
            <div className="ml-4">
              <h3 className="text-lg font-medium text-gray-700">
                Читальные залы
              </h3>
              <p className="text-2xl font-semibold">{readingHalls.length}</p>
            </div>
          </div>
          <div className="mt-3 text-sm">
            <span className="text-blue-600 font-medium">
              {readingHalls.reduce((sum, hall) => sum + hall.total_seats, 0)}{" "}
              мест
            </span>
          </div>
        </motion.div>

        <motion.div variants={itemVariants} className="card p-6">
          <div className="flex items-center">
            <div className="p-3 rounded-full bg-blue-100 text-blue-600">
              <Users size={24} />
            </div>
            <div className="ml-4">
              <h3 className="text-lg font-medium text-gray-700">Посетители</h3>
              <p className="text-2xl font-semibold">
                {dashboardStats.current_hall_visitors}
              </p>
            </div>
          </div>
          <div className="mt-3 text-sm">
            <span className="text-green-600 font-medium">сейчас в залах</span>
          </div>
        </motion.div>

        <motion.div variants={itemVariants} className="card p-6">
          <div className="flex items-center">
            <div className="p-3 rounded-full bg-yellow-100 text-yellow-600">
              <AlertCircle size={24} />
            </div>
            <div className="ml-4">
              <h3 className="text-lg font-medium text-gray-700">Штрафы</h3>
              <p className="text-2xl font-semibold">
                {dashboardStats.total_fines.toFixed(0)} ₽
              </p>
            </div>
          </div>
          <div className="mt-3 text-sm">
            <Link
              to="/fines"
              className="text-primary-500 hover:underline font-medium"
            >
              Управление штрафами
            </Link>
          </div>
        </motion.div>
      </motion.div>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        {/* Books Due This Week */}
        <motion.div
          variants={itemVariants}
          className="card col-span-1 lg:col-span-2 overflow-hidden"
        >
          <div className="p-4 bg-primary-500 text-white">
            <h2 className="text-xl font-semibold">
              Книги к возврату на этой неделе
            </h2>
          </div>
          <div className="p-4">
            {issuesDueThisWeek.length > 0 ? (
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
                        Читатель
                      </th>
                      <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                        Дата возврата
                      </th>
                    </tr>
                  </thead>
                  <tbody className="divide-y divide-gray-200">
                    {issuesDueThisWeek.map((issue) => {
                      const bookCopy = bookCopies.find(
                        (copy) => copy.id === issue.book_copy_id,
                      );
                      const book = booksWithDetails.find(
                        (b) => b.id === bookCopy?.book_id,
                      );
                      const reader = readers.find(
                        (r) => r.id === issue.reader_id,
                      );

                      return (
                        <tr key={issue.id}>
                          <td className="px-4 py-3">
                            <div className="font-medium text-gray-900">
                              {book?.title}
                            </div>
                            <div className="text-xs text-gray-500">
                              {book?.authors
                                .map((author) => author.full_name)
                                .join(", ")}
                            </div>
                          </td>
                          <td className="px-4 py-3 whitespace-nowrap">
                            <div className="text-sm text-gray-700">
                              {bookCopy?.copy_code}
                            </div>
                          </td>
                          <td className="px-4 py-3 whitespace-nowrap">
                            <div className="text-gray-700">
                              {reader?.full_name}
                            </div>
                            <div className="text-xs text-gray-500">
                              {reader?.ticket_number}
                            </div>
                          </td>
                          <td className="px-4 py-3 whitespace-nowrap">
                            <div className="text-gray-700">
                              {new Date(issue.due_date).toLocaleDateString(
                                "ru-RU",
                              )}
                            </div>
                          </td>
                        </tr>
                      );
                    })}
                  </tbody>
                </table>
              </div>
            ) : (
              <p className="text-gray-500 text-center py-4">
                Нет книг к возврату на этой неделе
              </p>
            )}
          </div>
        </motion.div>

        {/* Recent Activity */}
        <motion.div variants={itemVariants} className="card overflow-hidden">
          <div className="p-4 bg-accent-500 text-white">
            <h2 className="text-xl font-semibold">Последние операции</h2>
          </div>
          <div className="p-4">
            <ul className="divide-y divide-gray-200">
              {recentActivities.map((activity) => {
                const bookCopy = bookCopies.find(
                  (copy) => copy.id === activity.book_copy_id,
                );
                const book = booksWithDetails.find(
                  (b) => b.id === bookCopy?.book_id,
                );
                const reader = readers.find((r) => r.id === activity.reader_id);
                const isReturn = activity.return_date !== null;

                return (
                  <li key={activity.id} className="py-3">
                    <div className="flex items-start">
                      <div
                        className={`mt-1 p-2 rounded-full ${isReturn ? "bg-green-100 text-green-600" : "bg-blue-100 text-blue-600"}`}
                      >
                        {isReturn ? (
                          <BookOpen size={16} />
                        ) : (
                          <Clock size={16} />
                        )}
                      </div>
                      <div className="ml-3 flex-1">
                        <p className="text-sm font-medium text-gray-800">
                          {reader?.full_name}{" "}
                          {isReturn ? "вернул(а)" : "взял(а)"}
                        </p>
                        <p className="text-sm text-gray-600 line-clamp-1">
                          "{book?.title}"
                        </p>
                        <div className="flex items-center justify-between mt-1">
                          <span className="text-xs text-gray-500">
                            {bookCopy?.copy_code}
                          </span>
                          <span className="text-xs text-gray-500">
                            {new Date(
                              isReturn
                                ? activity.return_date!
                                : activity.issue_date,
                            ).toLocaleDateString("ru-RU")}
                          </span>
                        </div>
                      </div>
                    </div>
                  </li>
                );
              })}
            </ul>
          </div>
        </motion.div>
      </div>

      {/* Reading Halls Status */}
      <motion.div variants={itemVariants} className="card mt-6 overflow-hidden">
        <div className="p-4 bg-purple-500 text-white">
          <h2 className="text-xl font-semibold">Состояние читальных залов</h2>
        </div>
        <div className="p-4">
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
            {readingHalls.map((hall) => (
              <div key={hall.id} className="border rounded-lg p-4">
                <h3 className="font-medium text-gray-900 mb-2">
                  {hall.hall_name}
                </h3>
                {hall.specialization && (
                  <p className="text-sm text-gray-600 mb-2">
                    {hall.specialization}
                  </p>
                )}
                <div className="flex items-center justify-between">
                  <span className="text-sm text-gray-600">
                    Посетители: {hall.current_visitors}/{hall.total_seats}
                  </span>
                  <div className="w-16 bg-gray-200 rounded-full h-2">
                    <div
                      className="bg-blue-500 h-2 rounded-full"
                      style={{
                        width: `${Math.min((hall.current_visitors / hall.total_seats) * 100, 100)}%`,
                      }}
                    ></div>
                  </div>
                </div>
              </div>
            ))}
          </div>
        </div>
      </motion.div>
    </div>
  );
};

export default Dashboard;
