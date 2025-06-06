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
  Building2,
} from "lucide-react";
import {
  readers,
  readingHalls,
  hallVisits,
  hourlyStats,
} from "../data/mockData";
import type { Reader, ReadingHall, HallVisit } from "../types";

interface VisitorData {
  id: string;
  readerId: string;
  readerName: string;
  ticketNumber: string;
  entryTime: Date;
  exitTime?: Date;
  status: "in" | "out";
  hallId: string;
  hallName: string;
}

const RoomOverview: React.FC = () => {
  const [selectedDate] = useState(new Date());

  // Обработка посещений для получения текущего состояния
  const processVisits = (): VisitorData[] => {
    const visitorMap = new Map<string, VisitorData>();

    // Сортируем посещения по времени
    const sortedVisits = [...hallVisits].sort(
      (a, b) =>
        new Date(a.visit_time).getTime() - new Date(b.visit_time).getTime(),
    );

    sortedVisits.forEach((visit) => {
      const reader = readers.find((r) => r.id === visit.reader_id);
      const hall = readingHalls.find((h) => h.id === visit.hall_id);

      if (!reader || !hall) return;

      const key = `${visit.reader_id}-${visit.hall_id}`;

      if (visit.visit_type === "entry") {
        visitorMap.set(key, {
          id: visit.id,
          readerId: reader.id,
          readerName: reader.full_name,
          ticketNumber: reader.ticket_number,
          entryTime: new Date(visit.visit_time),
          status: "in",
          hallId: hall.id,
          hallName: hall.hall_name,
        });
      } else if (visit.visit_type === "exit") {
        const existingVisit = visitorMap.get(key);
        if (existingVisit) {
          existingVisit.exitTime = new Date(visit.visit_time);
          existingVisit.status = "out";
        }
      }
    });

    return Array.from(visitorMap.values()).sort(
      (a, b) => b.entryTime.getTime() - a.entryTime.getTime(),
    );
  };

  const todayVisitors = processVisits();

  // Подсчет текущих посетителей по залам
  const getCurrentOccupancy = (hallId: string): number => {
    return todayVisitors.filter((v) => v.hallId === hallId && v.status === "in")
      .length;
  };

  // Обновляем данные о залах с текущей загруженностью
  const roomsWithOccupancy = readingHalls.map((hall) => ({
    ...hall,
    currentOccupancy: getCurrentOccupancy(hall.id),
  }));

  const maxVisitors = Math.max(...hourlyStats.map((h) => h.visitors));
  const totalCapacity = readingHalls.reduce(
    (sum, hall) => sum + hall.total_seats,
    0,
  );
  const totalOccupancy = roomsWithOccupancy.reduce(
    (sum, hall) => sum + hall.currentOccupancy,
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
                <Building2 size={24} />
              </div>
              <div className="ml-4">
                <h3 className="text-sm font-medium text-gray-500">
                  Всего мест
                </h3>
                <p className="text-2xl font-semibold">{totalCapacity}</p>
                <p className="text-xs text-gray-500">
                  в {readingHalls.length} залах
                </p>
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
                <p className="text-xs text-gray-500">активных посетителей</p>
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
                  {totalCapacity > 0
                    ? Math.round((totalOccupancy / totalCapacity) * 100)
                    : 0}
                  %
                </p>
                <p className="text-xs text-gray-500">общая по всем залам</p>
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
                <p className="text-xs text-gray-500">уникальных визитов</p>
              </div>
            </div>
          </motion.div>
        </div>

        {/* Залы */}
        <motion.div variants={itemVariants} className="card">
          <div className="p-4 bg-primary-500 text-white">
            <h2 className="text-xl font-semibold flex items-center">
              <Building2 size={20} className="mr-2" />
              Состояние залов
            </h2>
          </div>
          <div className="p-6">
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
              {roomsWithOccupancy.map((hall) => (
                <div
                  key={hall.id}
                  className="border rounded-lg p-4 hover:shadow-md transition-shadow"
                >
                  <div className="flex justify-between items-start mb-3">
                    <div>
                      <h3 className="font-medium text-gray-900">
                        {hall.hall_name}
                      </h3>
                      {hall.specialization && (
                        <p className="text-sm text-gray-600">
                          {hall.specialization}
                        </p>
                      )}
                    </div>
                    <span
                      className={`px-2 py-1 rounded-full text-xs font-medium flex-shrink-0 ${
                        hall.currentOccupancy / hall.total_seats > 0.8
                          ? "bg-red-100 text-red-800"
                          : hall.currentOccupancy / hall.total_seats > 0.6
                            ? "bg-yellow-100 text-yellow-800"
                            : "bg-green-100 text-green-800"
                      }`}
                    >
                      {Math.round(
                        (hall.currentOccupancy / hall.total_seats) * 100,
                      )}
                      %
                    </span>
                  </div>
                  <div className="mb-2">
                    <div className="flex justify-between text-sm text-gray-600 mb-1">
                      <span>
                        {hall.currentOccupancy} из {hall.total_seats}
                      </span>
                      <span>
                        {hall.total_seats - hall.currentOccupancy} свободных
                      </span>
                    </div>
                    <div className="w-full bg-gray-200 rounded-full h-2">
                      <div
                        className={`h-2 rounded-full transition-all duration-300 ${
                          hall.currentOccupancy / hall.total_seats > 0.8
                            ? "bg-red-500"
                            : hall.currentOccupancy / hall.total_seats > 0.6
                              ? "bg-yellow-500"
                              : "bg-green-500"
                        }`}
                        style={{
                          width: `${Math.min((hall.currentOccupancy / hall.total_seats) * 100, 100)}%`,
                        }}
                      ></div>
                    </div>
                  </div>

                  {/* Список текущих посетителей зала */}
                  {hall.currentOccupancy > 0 && (
                    <div className="mt-3 pt-3 border-t border-gray-100">
                      <p className="text-xs text-gray-500 mb-2">
                        Сейчас в зале:
                      </p>
                      <div className="space-y-1 max-h-20 overflow-y-auto">
                        {todayVisitors
                          .filter(
                            (v) => v.hallId === hall.id && v.status === "in",
                          )
                          .slice(0, 3)
                          .map((visitor) => (
                            <div
                              key={`${visitor.readerId}-${hall.id}`}
                              className="text-xs text-gray-700"
                            >
                              {visitor.readerName}
                            </div>
                          ))}
                        {todayVisitors.filter(
                          (v) => v.hallId === hall.id && v.status === "in",
                        ).length > 3 && (
                          <div className="text-xs text-gray-500">
                            и еще{" "}
                            {todayVisitors.filter(
                              (v) => v.hallId === hall.id && v.status === "in",
                            ).length - 3}
                            ...
                          </div>
                        )}
                      </div>
                    </div>
                  )}
                </div>
              ))}
            </div>
          </div>
        </motion.div>

        <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
          {/* Минималистичный график посещений */}
          <motion.div variants={itemVariants} className="card">
            <div className="p-4 bg-accent-500 text-white">
              <h2 className="text-xl font-semibold flex items-center">
                <BarChart3 size={20} className="mr-2" />
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
                          width: `${maxVisitors > 0 ? (stat.visitors / maxVisitors) * 100 : 0}%`,
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
              <h2 className="text-xl font-semibold flex items-center">
                <Users size={20} className="mr-2" />
                Посетители сегодня
              </h2>
            </div>
            <div className="p-4 max-h-96 overflow-y-auto">
              {todayVisitors.length > 0 ? (
                <div className="space-y-2">
                  {todayVisitors.map((visitor) => (
                    <div
                      key={`${visitor.readerId}-${visitor.hallId}-${visitor.entryTime.getTime()}`}
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
                            to={`/readers/${visitor.readerId}`}
                            className="font-medium text-gray-900 hover:text-primary-600 transition-colors duration-200 block truncate"
                          >
                            {visitor.readerName}
                          </Link>
                          <div className="text-xs text-gray-600">
                            <div className="flex flex-wrap gap-2">
                              <span>Билет: {visitor.ticketNumber}</span>
                              <span>Зал: {visitor.hallName}</span>
                            </div>
                            <div className="mt-1">
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
                  ))}
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
