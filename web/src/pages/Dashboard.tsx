import React from "react";
import { Link } from "react-router-dom";
import { motion } from "framer-motion";
import { BookOpen, Users, Clock, AlertCircle } from "lucide-react";
import { books, members, checkoutRecords } from "../data/mockData";

const Dashboard: React.FC = () => {
  // Calculate statistics
  const totalBooks = books.length;
  const availableBooks = books.filter(
    (book) => book.status === "available",
  ).length;
  const totalMembers = members.length;
  const activeMembers = members.filter(
    (member) => member.membershipStatus === "active",
  ).length;

  const activeCheckouts = checkoutRecords.filter(
    (record) => record.status === "active",
  ).length;
  const overdueBooks = checkoutRecords.filter(
    (record) =>
      record.status === "active" && new Date(record.dueDate) < new Date(),
  ).length;

  // Books due this week
  const today = new Date();
  const endOfWeek = new Date(today);
  endOfWeek.setDate(today.getDate() + 7);

  const booksDueThisWeek = checkoutRecords.filter(
    (record) =>
      record.status === "active" &&
      new Date(record.dueDate) >= today &&
      new Date(record.dueDate) <= endOfWeek,
  );

  // Recent activities - using checkout records
  const recentActivities = checkoutRecords
    .sort((a, b) => {
      const dateA = a.returnedDate || a.checkoutDate;
      const dateB = b.returnedDate || b.checkoutDate;
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
              <p className="text-2xl font-semibold">{totalBooks}</p>
            </div>
          </div>
          <div className="mt-3 text-sm">
            <span className="text-green-600 font-medium">
              {availableBooks} доступно
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
              <p className="text-2xl font-semibold">{totalMembers}</p>
            </div>
          </div>
          <div className="mt-3 text-sm">
            <span className="text-green-600 font-medium">
              {activeMembers} активных
            </span>
          </div>
        </motion.div>

        <motion.div variants={itemVariants} className="card p-6">
          <div className="flex items-center">
            <div className="p-3 rounded-full bg-green-100 text-green-600">
              <BookOpen size={24} />
            </div>
            <div className="ml-4">
              <h3 className="text-lg font-medium text-gray-700">На руках</h3>
              <p className="text-2xl font-semibold">{activeCheckouts}</p>
            </div>
          </div>
          <div className="mt-3 text-sm">
            <Link
              to="/checkout"
              className="text-primary-500 hover:underline font-medium"
            >
              Подробнее
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
              <p className="text-2xl font-semibold">{overdueBooks}</p>
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
            {booksDueThisWeek.length > 0 ? (
              <div className="overflow-x-auto">
                <table className="min-w-full divide-y divide-gray-200">
                  <thead>
                    <tr>
                      <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                        Книга
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
                    {booksDueThisWeek.map((record) => {
                      const book = books.find((b) => b.id === record.bookId);
                      const member = members.find(
                        (m) => m.id === record.memberId,
                      );
                      return (
                        <tr key={record.id}>
                          <td className="px-4 py-3 whitespace-nowrap">
                            <div className="font-medium text-gray-900">
                              {book?.title}
                            </div>
                          </td>
                          <td className="px-4 py-3 whitespace-nowrap">
                            <div className="text-gray-700">
                              {member?.firstName} {member?.lastName}
                            </div>
                          </td>
                          <td className="px-4 py-3 whitespace-nowrap">
                            <div className="text-gray-700">
                              {new Date(record.dueDate).toLocaleDateString(
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
                const book = books.find((b) => b.id === activity.bookId);
                const member = members.find((m) => m.id === activity.memberId);
                const isReturn = activity.returnedDate !== null;

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
                      <div className="ml-3">
                        <p className="text-sm font-medium text-gray-800">
                          {member?.firstName} {member?.lastName}{" "}
                          {isReturn ? "вернул(а)" : "взял(а)"}
                        </p>
                        <p className="text-sm text-gray-600">
                          "{book?.title}"{" "}
                          {new Date(
                            isReturn
                              ? activity.returnedDate!
                              : activity.checkoutDate,
                          ).toLocaleDateString("ru-RU")}
                        </p>
                      </div>
                    </div>
                  </li>
                );
              })}
            </ul>
          </div>
        </motion.div>
      </div>
    </div>
  );
};

export default Dashboard;
