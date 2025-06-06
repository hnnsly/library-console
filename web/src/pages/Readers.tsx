import React, { useState } from "react";
import { Link } from "react-router-dom";
import { motion } from "framer-motion";
import {
  Search,
  PlusCircle,
  UserCheck,
  AlertTriangle,
  BookOpen,
  CreditCard,
} from "lucide-react";
import { readers, bookIssues, fines } from "../data/mockData";
import type { Reader } from "../types";

const Readers: React.FC = () => {
  const [searchTerm, setSearchTerm] = useState("");
  const [statusFilter, setStatusFilter] = useState("");

  // Filter readers based on search and status
  const filteredReaders = readers.filter((reader) => {
    const matchesSearch =
      reader.full_name.toLowerCase().includes(searchTerm.toLowerCase()) ||
      (reader.email &&
        reader.email.toLowerCase().includes(searchTerm.toLowerCase())) ||
      (reader.phone && reader.phone.includes(searchTerm)) ||
      reader.ticket_number.includes(searchTerm);

    const matchesStatus = statusFilter
      ? statusFilter === "active"
        ? reader.is_active
        : !reader.is_active
      : true;

    return matchesSearch && matchesStatus;
  });

  // Count active borrowed books for a reader
  const getActiveBorrowedCount = (readerId: string) => {
    return bookIssues.filter(
      (issue) => issue.reader_id === readerId && !issue.return_date,
    ).length;
  };

  // Get unpaid fines total for a reader
  const getUnpaidFinesTotal = (readerId: string) => {
    return fines
      .filter((fine) => fine.reader_id === readerId && !fine.is_paid)
      .reduce((sum, fine) => sum + fine.amount, 0);
  };

  // Animation variants
  const containerVariants = {
    hidden: { opacity: 0 },
    visible: {
      opacity: 1,
      transition: {
        staggerChildren: 0.05,
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
      <div className="flex flex-col md:flex-row md:items-center justify-between mb-6">
        <div>
          <h1 className="text-3xl font-serif font-bold text-gray-900">
            Читатели
          </h1>
          <p className="text-gray-600 mt-1">
            Управление читателями и их учётными записями в системе libr.
          </p>
        </div>
        <Link
          to="/readers/add"
          className="btn btn-primary mt-3 md:mt-0 flex items-center justify-center md:justify-start"
        >
          <PlusCircle size={18} className="mr-2" /> Добавить читателя
        </Link>
      </div>

      <div className="bg-white rounded-lg shadow-sm p-4 mb-6">
        <div className="flex flex-col md:flex-row gap-4">
          {/* Search */}
          <div className="relative flex-grow">
            <Search
              className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400"
              size={18}
            />
            <input
              type="text"
              placeholder="Поиск по имени, email, телефону или номеру билета..."
              value={searchTerm}
              onChange={(e) => setSearchTerm(e.target.value)}
              className="input pl-10 py-2"
            />
          </div>

          {/* Status filter */}
          <select
            value={statusFilter}
            onChange={(e) => setStatusFilter(e.target.value)}
            className="input py-2 md:w-48"
          >
            <option value="">Все статусы</option>
            <option value="active">Активные</option>
            <option value="inactive">Неактивные</option>
          </select>
        </div>
      </div>

      {filteredReaders.length > 0 ? (
        <motion.div
          variants={containerVariants}
          initial="hidden"
          animate="visible"
          className="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-3 gap-6"
        >
          {filteredReaders.map((reader) => {
            const activeBooksCount = getActiveBorrowedCount(reader.id);
            const unpaidFines = getUnpaidFinesTotal(reader.id);

            return (
              <motion.div
                key={reader.id}
                variants={itemVariants}
                className="card overflow-hidden hover:shadow-lg transition-shadow duration-300"
              >
                <div
                  className={`h-2 ${
                    reader.is_active ? "bg-green-500" : "bg-red-500"
                  }`}
                ></div>

                <div className="p-5">
                  {/* Header with avatar and status */}
                  <div className="flex items-start justify-between mb-4">
                    <div className="flex items-center min-w-0 flex-1">
                      <div className="w-10 h-10 rounded-full bg-primary-500 text-white flex items-center justify-center text-sm font-medium flex-shrink-0">
                        {reader.full_name
                          .split(" ")
                          .map((name) => name.charAt(0))
                          .join("")
                          .slice(0, 2)}
                      </div>
                      <div className="ml-3 min-w-0 flex-1">
                        <Link to={`/readers/${reader.id}`} className="block">
                          <h3 className="font-medium text-gray-900 hover:text-primary-600 transition-colors truncate">
                            {reader.full_name}
                          </h3>
                        </Link>
                        <p className="text-sm text-gray-600 truncate">
                          Билет: {reader.ticket_number}
                        </p>
                        {reader.email && (
                          <p className="text-xs text-gray-500 truncate">
                            {reader.email}
                          </p>
                        )}
                      </div>
                    </div>
                    <div className="flex-shrink-0 ml-2">
                      {reader.is_active ? (
                        <span className="inline-flex items-center px-2 py-1 rounded-full text-xs font-medium bg-green-100 text-green-800">
                          <UserCheck size={10} className="mr-1" />
                          <span className="hidden sm:inline">Активен</span>
                        </span>
                      ) : (
                        <span className="inline-flex items-center px-2 py-1 rounded-full text-xs font-medium bg-red-100 text-red-800">
                          <AlertTriangle size={10} className="mr-1" />
                          <span className="hidden sm:inline">Неактивен</span>
                        </span>
                      )}
                    </div>
                  </div>

                  {/* Contact info */}
                  {reader.phone && (
                    <div className="mb-3">
                      <p className="text-xs text-gray-500">Телефон:</p>
                      <p className="text-sm text-gray-700">{reader.phone}</p>
                    </div>
                  )}

                  {/* Stats */}
                  <div className="grid grid-cols-2 gap-3 mb-4">
                    <div className="bg-blue-50 p-3 rounded-lg text-center">
                      <div className="flex items-center justify-center mb-1">
                        <BookOpen size={14} className="text-blue-500 mr-1" />
                        <p className="text-blue-600 text-xs font-medium">
                          Книг на руках
                        </p>
                      </div>
                      <p className="text-lg font-semibold text-blue-700">
                        {activeBooksCount}
                      </p>
                    </div>
                    <div className="bg-orange-50 p-3 rounded-lg text-center">
                      <div className="flex items-center justify-center mb-1">
                        <CreditCard
                          size={14}
                          className="text-orange-500 mr-1"
                        />
                        <p className="text-orange-600 text-xs font-medium">
                          Штрафы
                        </p>
                      </div>
                      <p
                        className={`text-lg font-semibold ${
                          unpaidFines > 0 ? "text-red-600" : "text-orange-700"
                        }`}
                      >
                        {unpaidFines.toFixed(0)} ₽
                      </p>
                    </div>
                  </div>

                  {/* Warnings */}
                  {unpaidFines > 0 && (
                    <div className="mb-3 p-2 bg-red-50 border border-red-100 rounded-lg">
                      <div className="flex items-center">
                        <AlertTriangle
                          size={14}
                          className="text-red-500 mr-2"
                        />
                        <span className="text-xs text-red-700 font-medium">
                          Есть неоплаченные штрафы
                        </span>
                      </div>
                    </div>
                  )}

                  {/* Footer */}
                  <div className="flex items-center justify-between text-xs text-gray-500">
                    <span className="truncate">
                      Регистрация:{" "}
                      {new Date(reader.registration_date).toLocaleDateString(
                        "ru-RU",
                      )}
                    </span>
                    <Link
                      to={`/readers/${reader.id}`}
                      className="text-primary-600 hover:text-primary-800 font-medium whitespace-nowrap ml-2"
                    >
                      Подробнее
                    </Link>
                  </div>
                </div>
              </motion.div>
            );
          })}
        </motion.div>
      ) : (
        <div className="text-center py-12 bg-white rounded-lg shadow-sm">
          <p className="text-gray-500 text-lg">
            Читатели не найдены. Попробуйте изменить критерии поиска.
          </p>
          {statusFilter && (
            <button
              onClick={() => setStatusFilter("")}
              className="mt-3 text-primary-600 hover:underline"
            >
              Очистить фильтр статуса
            </button>
          )}
        </div>
      )}

      {/* Statistics summary */}
      <div className="mt-8 grid grid-cols-1 md:grid-cols-4 gap-4">
        <div className="bg-white rounded-lg shadow-sm p-4 text-center">
          <div className="text-2xl font-bold text-gray-900">
            {readers.length}
          </div>
          <div className="text-sm text-gray-600">Всего читателей</div>
        </div>
        <div className="bg-white rounded-lg shadow-sm p-4 text-center">
          <div className="text-2xl font-bold text-green-600">
            {readers.filter((r) => r.is_active).length}
          </div>
          <div className="text-sm text-gray-600">Активных</div>
        </div>
        <div className="bg-white rounded-lg shadow-sm p-4 text-center">
          <div className="text-2xl font-bold text-blue-600">
            {bookIssues.filter((issue) => !issue.return_date).length}
          </div>
          <div className="text-sm text-gray-600">Книг выдано</div>
        </div>
        <div className="bg-white rounded-lg shadow-sm p-4 text-center">
          <div className="text-2xl font-bold text-orange-600">
            {fines.filter((fine) => !fine.is_paid).length}
          </div>
          <div className="text-sm text-gray-600">Неоплаченных штрафов</div>
        </div>
      </div>
    </div>
  );
};

export default Readers;
