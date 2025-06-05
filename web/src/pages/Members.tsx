import React, { useState } from "react";
import { Link } from "react-router-dom";
import { motion } from "framer-motion";
import {
  Search,
  PlusCircle,
  UserCheck,
  AlertTriangle,
  Clock,
} from "lucide-react";
import { members } from "../data/mockData";

const Members: React.FC = () => {
  const [searchTerm, setSearchTerm] = useState("");
  const [statusFilter, setStatusFilter] = useState("");

  // Filter members based on search and status
  const filteredMembers = members.filter((member) => {
    const matchesSearch =
      `${member.firstName} ${member.lastName}`
        .toLowerCase()
        .includes(searchTerm.toLowerCase()) ||
      member.email.toLowerCase().includes(searchTerm.toLowerCase()) ||
      member.phone.includes(searchTerm);

    const matchesStatus = statusFilter
      ? member.membershipStatus === statusFilter
      : true;

    return matchesSearch && matchesStatus;
  });

  // Count active borrowed books
  const getActiveBorrowedCount = (memberId: string) => {
    const member = members.find((m) => m.id === memberId);
    return member ? member.borrowedBooks.length : 0;
  };

  // Get status text in Russian
  // const getStatusText = (status: string) => {
  //   switch (status) {
  //     case "active":
  //       return "Активен";
  //     case "expired":
  //       return "Истёк";
  //     case "suspended":
  //       return "Заблокирован";
  //     default:
  //       return status;
  //   }
  // };

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
          to="/members/add"
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
              placeholder="Поиск по имени, email или телефону..."
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
            <option value="expired">Истёкшие</option>
            <option value="suspended">Заблокированные</option>
          </select>
        </div>
      </div>

      {filteredMembers.length > 0 ? (
        <motion.div
          variants={containerVariants}
          initial="hidden"
          animate="visible"
          className="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-3 gap-6"
        >
          {filteredMembers.map((member) => (
            <motion.div
              key={member.id}
              variants={itemVariants}
              className="card overflow-hidden hover:shadow-lg transition-shadow duration-300"
            >
              <div
                className={`h-2 ${
                  member.membershipStatus === "active"
                    ? "bg-green-500"
                    : member.membershipStatus === "expired"
                      ? "bg-orange-500"
                      : "bg-red-500"
                }`}
              ></div>

              <div className="p-5">
                {/* Header with avatar and status */}
                <div className="flex items-start justify-between mb-4">
                  <div className="flex items-center min-w-0 flex-1">
                    <div className="w-10 h-10 rounded-full bg-primary-500 text-white flex items-center justify-center text-sm font-medium flex-shrink-0">
                      {member.firstName.charAt(0)}
                      {member.lastName.charAt(0)}
                    </div>
                    <div className="ml-3 min-w-0 flex-1">
                      <Link to={`/members/${member.id}`} className="block">
                        <h3 className="font-medium text-gray-900 hover:text-primary-600 transition-colors truncate">
                          {member.firstName} {member.lastName}
                        </h3>
                      </Link>
                      <p className="text-sm text-gray-600 truncate">
                        {member.email}
                      </p>
                    </div>
                  </div>
                  <div className="flex-shrink-0 ml-2">
                    {member.membershipStatus === "active" ? (
                      <span className="inline-flex items-center px-2 py-1 rounded-full text-xs font-medium bg-green-100 text-green-800">
                        <UserCheck size={10} className="mr-1" />
                        <span className="hidden sm:inline">Активен</span>
                      </span>
                    ) : member.membershipStatus === "expired" ? (
                      <span className="inline-flex items-center px-2 py-1 rounded-full text-xs font-medium bg-orange-100 text-orange-800">
                        <Clock size={10} className="mr-1" />
                        <span className="hidden sm:inline">Истёк</span>
                      </span>
                    ) : (
                      <span className="inline-flex items-center px-2 py-1 rounded-full text-xs font-medium bg-red-100 text-red-800">
                        <AlertTriangle size={10} className="mr-1" />
                        <span className="hidden sm:inline">Заблок.</span>
                      </span>
                    )}
                  </div>
                </div>

                {/* Stats */}
                <div className="grid grid-cols-2 gap-3 mb-4">
                  <div className="bg-gray-50 p-3 rounded-lg text-center">
                    <p className="text-gray-500 text-xs mb-1">Взято книг</p>
                    <p className="text-lg font-semibold text-gray-900">
                      {getActiveBorrowedCount(member.id)}
                    </p>
                  </div>
                  <div className="bg-gray-50 p-3 rounded-lg text-center">
                    <p className="text-gray-500 text-xs mb-1">Штрафы</p>
                    <p
                      className={`text-lg font-semibold ${member.fines > 0 ? "text-red-600" : "text-gray-900"}`}
                    >
                      ₽{member.fines.toFixed(0)}
                    </p>
                  </div>
                </div>

                {/* Footer */}
                <div className="flex items-center justify-between text-xs text-gray-500">
                  <span className="truncate">
                    С {new Date(member.joinDate).toLocaleDateString("ru-RU")}
                  </span>
                  <Link
                    to={`/members/${member.id}`}
                    className="text-primary-600 hover:text-primary-800 font-medium whitespace-nowrap ml-2"
                  >
                    Подробнее
                  </Link>
                </div>
              </div>
            </motion.div>
          ))}
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
    </div>
  );
};

export default Members;
