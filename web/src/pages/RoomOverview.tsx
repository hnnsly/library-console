import React, { useState } from "react";
import { Link } from "react-router-dom";
import { motion } from "framer-motion";
import {
  BarChart3,
  Users,
  Clock,
  TrendingUp,
  Calendar,
  User,
  DoorOpen,
  DoorClosed,
} from "lucide-react";
import { members, rooms, todayVisitors, hourlyStats } from "../data/mockData";

const RoomOverview: React.FC = () => {
  const [selectedDate] = useState(new Date());

  const maxVisitors = Math.max(...hourlyStats.map((h) => h.visitors));
  const totalCapacity = rooms.reduce((sum, room) => sum + room.capacity, 0);
  const totalOccupancy = rooms.reduce(
    (sum, room) => sum + room.currentOccupancy,
    0,
  );
  const currentVisitors = todayVisitors.filter((v) => v.status === "in").length;

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
          Обзор читальных залов
        </h1>
        <p className="text-gray-600 mt-1">
          Мониторинг посещаемости и загруженности залов в системе libr
        </p>
      </div>

      <motion.div
        variants={containerVariants}
        initial="hidden"
        animate="visible"
        className="space-y-6"
      >
        {/* Общая статистика */}
        <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
          <motion.div variants={itemVariants} className="card p-6">
            <div className="flex items-center">
              <div className="p-3 rounded-full bg-primary-100 text-primary-600">
                <Users size={24} />
              </div>
              <div className="ml-4">
                <h3 className="text-sm font-medium text-gray-500">
                  Всего мест
                </h3>
                <p className="text-2xl font-semibold">{totalCapacity}</p>
              </div>
            </div>
          </motion.div>

          <motion.div variants={itemVariants} className="card p-6">
            <div className="flex items-center">
              <div className="p-3 rounded-full bg-green-100 text-green-600">
                <User size={24} />
              </div>
              <div className="ml-4">
                <h3 className="text-sm font-medium text-gray-500">
                  Сейчас в залах
                </h3>
                <p className="text-2xl font-semibold">{currentVisitors}</p>
              </div>
            </div>
          </motion.div>

          <motion.div variants={itemVariants} className="card p-6">
            <div className="flex items-center">
              <div className="p-3 rounded-full bg-blue-100 text-blue-600">
                <TrendingUp size={24} />
              </div>
              <div className="ml-4">
                <h3 className="text-sm font-medium text-gray-500">
                  Загруженность
                </h3>
                <p className="text-2xl font-semibold">
                  {Math.round((totalOccupancy / totalCapacity) * 100)}%
                </p>
              </div>
            </div>
          </motion.div>

          <motion.div variants={itemVariants} className="card p-6">
            <div className="flex items-center">
              <div className="p-3 rounded-full bg-accent-100 text-accent-600">
                <Calendar size={24} />
              </div>
              <div className="ml-4">
                <h3 className="text-sm font-medium text-gray-500">
                  Посещений сегодня
                </h3>
                <p className="text-2xl font-semibold">{todayVisitors.length}</p>
              </div>
            </div>
          </motion.div>
        </div>

        {/* Залы */}
        <motion.div variants={itemVariants} className="card">
          <div className="p-4 bg-primary-500 text-white">
            <h2 className="text-xl font-semibold">Состояние залов</h2>
          </div>
          <div className="p-6">
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              {rooms.map((room) => (
                <div key={room.id} className="border rounded-lg p-4">
                  <div className="flex justify-between items-start mb-3">
                    <h3 className="font-medium text-gray-900">{room.name}</h3>
                    <span
                      className={`px-2 py-1 rounded-full text-xs font-medium ${
                        room.currentOccupancy / room.capacity > 0.8
                          ? "bg-red-100 text-red-800"
                          : room.currentOccupancy / room.capacity > 0.6
                            ? "bg-yellow-100 text-yellow-800"
                            : "bg-green-100 text-green-800"
                      }`}
                    >
                      {Math.round(
                        (room.currentOccupancy / room.capacity) * 100,
                      )}
                      % загружен
                    </span>
                  </div>
                  <div className="mb-2">
                    <div className="flex justify-between text-sm text-gray-600 mb-1">
                      <span>
                        {room.currentOccupancy} из {room.capacity}
                      </span>
                      <span>
                        {room.capacity - room.currentOccupancy} свободных
                      </span>
                    </div>
                    <div className="w-full bg-gray-200 rounded-full h-2">
                      <div
                        className={`h-2 rounded-full transition-all duration-300 ${
                          room.currentOccupancy / room.capacity > 0.8
                            ? "bg-red-500"
                            : room.currentOccupancy / room.capacity > 0.6
                              ? "bg-yellow-500"
                              : "bg-green-500"
                        }`}
                        style={{
                          width: `${(room.currentOccupancy / room.capacity) * 100}%`,
                        }}
                      ></div>
                    </div>
                  </div>
                </div>
              ))}
            </div>
          </div>
        </motion.div>

        <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
          {/* Минималистичный график посещений */}
          <motion.div variants={itemVariants} className="card">
            <div className="p-4 bg-accent-500 text-white">
              <h2 className="text-xl font-semibold">
                График посещений сегодня
              </h2>
            </div>
            <div className="p-6">
              <div className="space-y-2">
                {hourlyStats.map((stat) => (
                  <div key={stat.hour} className="flex items-center group">
                    <div className="w-10 text-xs text-gray-500 font-medium">
                      {stat.hour}:00
                    </div>
                    <div className="flex-1 mx-3 bg-gray-100 rounded-full h-3 overflow-hidden">
                      <motion.div
                        initial={{ width: 0 }}
                        animate={{
                          width: `${(stat.visitors / maxVisitors) * 100}%`,
                        }}
                        transition={{
                          duration: 0.8,
                          delay: stat.hour * 0.05,
                          ease: "easeOut",
                        }}
                        className="bg-gradient-to-r from-accent-400 to-accent-600 h-full rounded-full relative"
                      >
                        <div className="absolute inset-0 bg-accent-300 opacity-0 group-hover:opacity-30 transition-opacity duration-200 rounded-full"></div>
                      </motion.div>
                    </div>
                    <div className="w-8 text-xs text-gray-700 font-medium text-right">
                      {stat.visitors}
                    </div>
                  </div>
                ))}
              </div>
              <div className="mt-4 pt-3 border-t border-gray-200">
                <div className="flex justify-between text-xs text-gray-500">
                  <span>Часы работы</span>
                  <span>
                    Пик: {Math.max(...hourlyStats.map((h) => h.visitors))} чел.
                  </span>
                </div>
              </div>
            </div>
          </motion.div>

          {/* Список посетителей */}
          <motion.div variants={itemVariants} className="card">
            <div className="p-4 bg-green-500 text-white">
              <h2 className="text-xl font-semibold">Посетители сегодня</h2>
            </div>
            <div className="p-4 max-h-96 overflow-y-auto">
              {todayVisitors.length > 0 ? (
                <div className="space-y-2">
                  {todayVisitors.map((visitor) => {
                    const member = members.find(
                      (m) => m.id === visitor.memberId,
                    );
                    return (
                      <div
                        key={visitor.id}
                        className="flex items-center justify-between p-3 bg-gray-50 rounded-lg hover:bg-gray-100 transition-colors duration-200"
                      >
                        <div className="flex items-center min-w-0 flex-1">
                          <div
                            className={`p-1.5 rounded-full flex-shrink-0 ${
                              visitor.status === "in"
                                ? "bg-green-100 text-green-600"
                                : "bg-gray-100 text-gray-600"
                            }`}
                          >
                            {visitor.status === "in" ? (
                              <DoorOpen size={14} />
                            ) : (
                              <DoorClosed size={14} />
                            )}
                          </div>
                          <div className="ml-3 min-w-0 flex-1">
                            <Link
                              to={`/members/${visitor.memberId}`}
                              className="font-medium text-gray-900 hover:text-primary-600 transition-colors duration-200 block truncate"
                            >
                              {visitor.memberName}
                            </Link>
                            <div className="text-xs text-gray-600">
                              <span className="inline-block">
                                Вход:{" "}
                                {visitor.entryTime.toLocaleTimeString("ru-RU", {
                                  hour: "2-digit",
                                  minute: "2-digit",
                                })}
                              </span>
                              {visitor.exitTime && (
                                <span className="inline-block ml-2">
                                  Выход:{" "}
                                  {visitor.exitTime.toLocaleTimeString(
                                    "ru-RU",
                                    { hour: "2-digit", minute: "2-digit" },
                                  )}
                                </span>
                              )}
                            </div>
                          </div>
                        </div>
                        <span
                          className={`px-2 py-1 rounded-full text-xs font-medium flex-shrink-0 ml-2 ${
                            visitor.status === "in"
                              ? "bg-green-100 text-green-800"
                              : "bg-gray-100 text-gray-800"
                          }`}
                        >
                          {visitor.status === "in" ? "В зале" : "Покинул(а)"}
                        </span>
                      </div>
                    );
                  })}
                </div>
              ) : (
                <p className="text-gray-500 text-center py-8">
                  Посетителей сегодня пока не было
                </p>
              )}
            </div>
          </motion.div>
        </div>
      </motion.div>
    </div>
  );
};

export default RoomOverview;
