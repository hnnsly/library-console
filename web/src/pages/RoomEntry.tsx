import React, { useState } from "react";
import { motion } from "framer-motion";
import {
  DoorOpen,
  DoorClosed,
  User,
  Clock,
  AlertCircle,
  CheckCircle,
  Building2,
} from "lucide-react";
import { readers, hallVisits, readingHalls } from "../data/mockData";
import type { Reader, HallVisit, ReadingHall } from "../types";

interface ExtendedHallVisit extends HallVisit {
  reader?: Reader;
  hall?: ReadingHall;
}

const RoomEntry: React.FC = () => {
  const [ticketNumber, setTicketNumber] = useState("");
  const [selectedHallId, setSelectedHallId] = useState(
    readingHalls[0]?.id || "",
  );
  const [visits, setVisits] = useState<ExtendedHallVisit[]>(
    hallVisits.map((visit) => ({
      ...visit,
      reader: readers.find((r) => r.id === visit.reader_id),
      hall: readingHalls.find((h) => h.id === visit.hall_id),
    })),
  );
  const [message, setMessage] = useState<{
    type: "success" | "error";
    text: string;
  } | null>(null);

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();

    if (!ticketNumber.trim()) {
      setMessage({ type: "error", text: "Введите номер читательского билета" });
      return;
    }

    if (!selectedHallId) {
      setMessage({ type: "error", text: "Выберите читальный зал" });
      return;
    }

    // Найти читателя по номеру билета
    const reader = readers.find(
      (r) =>
        r.ticket_number === ticketNumber ||
        r.full_name.toLowerCase().includes(ticketNumber.toLowerCase()) ||
        r.id === ticketNumber,
    );

    if (!reader) {
      setMessage({
        type: "error",
        text: "Читатель с таким номером билета не найден",
      });
      return;
    }

    if (!reader.is_active) {
      setMessage({
        type: "error",
        text: "Читательский билет неактивен",
      });
      return;
    }

    const selectedHall = readingHalls.find((h) => h.id === selectedHallId);

    // Проверить, находится ли читатель уже в этом зале
    const activeVisit = visits.find(
      (v) =>
        v.reader_id === reader.id &&
        v.hall_id === selectedHallId &&
        v.visit_type === "entry" &&
        !visits.some(
          (exitVisit) =>
            exitVisit.reader_id === reader.id &&
            exitVisit.hall_id === selectedHallId &&
            exitVisit.visit_type === "exit" &&
            exitVisit.visit_time > v.visit_time,
        ),
    );

    if (activeVisit) {
      // Выход из зала
      const exitVisit: ExtendedHallVisit = {
        id: `visit-${Date.now()}`,
        reader_id: reader.id,
        hall_id: selectedHallId,
        visit_type: "exit",
        visit_time: new Date(),
        reader,
        hall: selectedHall,
      };

      setVisits((prev) => [exitVisit, ...prev]);
      setMessage({
        type: "success",
        text: `${reader.full_name} покинул(а) ${selectedHall?.hall_name}`,
      });
    } else {
      // Проверить лимит мест в зале
      const currentOccupancy = getCurrentOccupancy(selectedHallId);
      if (selectedHall && currentOccupancy >= selectedHall.total_seats) {
        setMessage({
          type: "error",
          text: `${selectedHall.hall_name} заполнен (${selectedHall.total_seats}/${selectedHall.total_seats} мест)`,
        });
        return;
      }

      // Вход в зал
      const entryVisit: ExtendedHallVisit = {
        id: `visit-${Date.now()}`,
        reader_id: reader.id,
        hall_id: selectedHallId,
        visit_type: "entry",
        visit_time: new Date(),
        reader,
        hall: selectedHall,
      };

      setVisits((prev) => [entryVisit, ...prev]);
      setMessage({
        type: "success",
        text: `${reader.full_name} вошёл(ла) в ${selectedHall?.hall_name}`,
      });
    }

    setTicketNumber("");
    setTimeout(() => setMessage(null), 3000);
  };

  // Получить текущую занятость зала
  const getCurrentOccupancy = (hallId: string): number => {
    const entries = visits.filter(
      (v) => v.hall_id === hallId && v.visit_type === "entry",
    );
    const exits = visits.filter(
      (v) => v.hall_id === hallId && v.visit_type === "exit",
    );

    const currentVisitors = new Set();

    // Сначала добавляем всех, кто вошел
    entries.forEach((entry) => {
      currentVisitors.add(entry.reader_id);
    });

    // Затем убираем тех, кто вышел
    exits.forEach((exit) => {
      // Проверяем, есть ли более поздний вход после этого выхода
      const laterEntry = entries.find(
        (entry) =>
          entry.reader_id === exit.reader_id &&
          entry.visit_time > exit.visit_time,
      );

      if (!laterEntry) {
        currentVisitors.delete(exit.reader_id);
      }
    });

    return currentVisitors.size;
  };

  // Получить текущих посетителей зала
  const getCurrentVisitors = (hallId: string): Reader[] => {
    const entries = visits.filter(
      (v) => v.hall_id === hallId && v.visit_type === "entry",
    );
    const exits = visits.filter(
      (v) => v.hall_id === hallId && v.visit_type === "exit",
    );

    const currentVisitorIds = new Set<string>();

    entries.forEach((entry) => {
      currentVisitorIds.add(entry.reader_id);
    });

    exits.forEach((exit) => {
      const laterEntry = entries.find(
        (entry) =>
          entry.reader_id === exit.reader_id &&
          entry.visit_time > exit.visit_time,
      );

      if (!laterEntry) {
        currentVisitorIds.delete(exit.reader_id);
      }
    });

    return readers.filter((r) => currentVisitorIds.has(r.id));
  };

  const totalCurrentVisitors = readingHalls.reduce(
    (sum, hall) => sum + getCurrentOccupancy(hall.id),
    0,
  );

  return (
    <div className="max-w-6xl mx-auto">
      <div className="mb-6">
        <h1 className="text-3xl font-serif font-bold text-gray-900">
          Вход в читальные залы
        </h1>
        <p className="text-gray-600 mt-1">
          Регистрация входа и выхода читателей в читальные залы системы libr
        </p>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6 mb-8">
        {/* Форма ввода */}
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          className="card p-6 lg:col-span-2"
        >
          <div className="flex items-center mb-4">
            <DoorOpen className="text-primary-600 mr-3" size={24} />
            <h2 className="text-xl font-semibold">Регистрация входа/выхода</h2>
          </div>

          <form onSubmit={handleSubmit} className="space-y-4">
            <div>
              <label
                htmlFor="ticketNumber"
                className="block text-sm font-medium text-gray-700 mb-2"
              >
                Номер читательского билета
              </label>
              <input
                type="text"
                id="ticketNumber"
                value={ticketNumber}
                onChange={(e) => setTicketNumber(e.target.value)}
                className="input w-full"
                placeholder="Введите номер билета..."
                autoFocus
              />
            </div>

            <div>
              <label
                htmlFor="hallSelect"
                className="block text-sm font-medium text-gray-700 mb-2"
              >
                Читальный зал
              </label>
              <select
                id="hallSelect"
                value={selectedHallId}
                onChange={(e) => setSelectedHallId(e.target.value)}
                className="input w-full"
              >
                {readingHalls.map((hall) => {
                  const occupancy = getCurrentOccupancy(hall.id);
                  return (
                    <option key={hall.id} value={hall.id}>
                      {hall.hall_name} ({occupancy}/{hall.total_seats} мест)
                    </option>
                  );
                })}
              </select>
            </div>

            <button type="submit" className="btn btn-primary w-full">
              Зарегистрировать
            </button>
          </form>

          {message && (
            <motion.div
              initial={{ opacity: 0, scale: 0.95 }}
              animate={{ opacity: 1, scale: 1 }}
              className={`mt-4 p-3 rounded-lg flex items-center ${
                message.type === "success"
                  ? "bg-green-100 text-green-800"
                  : "bg-red-100 text-red-800"
              }`}
            >
              {message.type === "success" ? (
                <CheckCircle size={16} className="mr-2" />
              ) : (
                <AlertCircle size={16} className="mr-2" />
              )}
              {message.text}
            </motion.div>
          )}
        </motion.div>

        {/* Статистика */}
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ delay: 0.1 }}
          className="card p-6"
        >
          <div className="flex items-center mb-4">
            <User className="text-accent-600 mr-3" size={24} />
            <h2 className="text-xl font-semibold">Статистика</h2>
          </div>

          <div className="space-y-4">
            <div className="bg-green-50 p-4 rounded-lg">
              <div className="flex items-center justify-between">
                <span className="text-green-800 font-medium">
                  Всего посетителей
                </span>
                <span className="text-2xl font-bold text-green-600">
                  {totalCurrentVisitors}
                </span>
              </div>
            </div>

            <div className="bg-blue-50 p-4 rounded-lg">
              <div className="flex items-center justify-between">
                <span className="text-blue-800 font-medium">
                  Записей сегодня
                </span>
                <span className="text-2xl font-bold text-blue-600">
                  {visits.length}
                </span>
              </div>
            </div>
          </div>
        </motion.div>
      </div>

      {/* Статистика по залам */}
      <motion.div
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ delay: 0.2 }}
        className="card mb-6"
      >
        <div className="p-4 bg-purple-500 text-white">
          <h2 className="text-xl font-semibold flex items-center">
            <Building2 size={20} className="mr-2" />
            Состояние читальных залов
          </h2>
        </div>
        <div className="p-4">
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
            {readingHalls.map((hall) => {
              const occupancy = getCurrentOccupancy(hall.id);
              const occupancyPercent = (occupancy / hall.total_seats) * 100;

              return (
                <div key={hall.id} className="border rounded-lg p-4">
                  <h3 className="font-medium text-gray-900 mb-2">
                    {hall.hall_name}
                  </h3>
                  {hall.specialization && (
                    <p className="text-sm text-gray-600 mb-2">
                      {hall.specialization}
                    </p>
                  )}
                  <div className="flex items-center justify-between mb-2">
                    <span className="text-sm text-gray-600">
                      Посетители: {occupancy}/{hall.total_seats}
                    </span>
                    <span className="text-sm font-medium text-gray-900">
                      {occupancyPercent.toFixed(0)}%
                    </span>
                  </div>
                  <div className="w-full bg-gray-200 rounded-full h-2">
                    <div
                      className={`h-2 rounded-full ${
                        occupancyPercent >= 90
                          ? "bg-red-500"
                          : occupancyPercent >= 70
                            ? "bg-yellow-500"
                            : "bg-green-500"
                      }`}
                      style={{
                        width: `${Math.min(occupancyPercent, 100)}%`,
                      }}
                    ></div>
                  </div>
                </div>
              );
            })}
          </div>
        </div>
      </motion.div>

      {/* Последние записи */}
      <motion.div
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ delay: 0.3 }}
        className="card"
      >
        <div className="p-4 bg-primary-500 text-white">
          <h2 className="text-xl font-semibold">Последние записи</h2>
        </div>
        <div className="p-4">
          {visits.length > 0 ? (
            <div className="overflow-x-auto">
              <table className="min-w-full divide-y divide-gray-200">
                <thead>
                  <tr>
                    <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                      Читатель
                    </th>
                    <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                      Читальный зал
                    </th>
                    <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                      Время
                    </th>
                    <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                      Действие
                    </th>
                  </tr>
                </thead>
                <tbody className="divide-y divide-gray-200">
                  {visits.slice(0, 10).map((visit) => (
                    <tr key={visit.id} className="hover:bg-gray-50">
                      <td className="px-4 py-3 whitespace-nowrap">
                        <div>
                          <div className="font-medium text-gray-900">
                            {visit.reader?.full_name}
                          </div>
                          <div className="text-sm text-gray-500">
                            Билет: {visit.reader?.ticket_number}
                          </div>
                        </div>
                      </td>
                      <td className="px-4 py-3 whitespace-nowrap">
                        <div className="text-sm text-gray-900">
                          {visit.hall?.hall_name}
                        </div>
                        {visit.hall?.specialization && (
                          <div className="text-xs text-gray-500">
                            {visit.hall.specialization}
                          </div>
                        )}
                      </td>
                      <td className="px-4 py-3 whitespace-nowrap text-gray-700">
                        {visit.visit_time.toLocaleTimeString("ru-RU")}
                      </td>
                      <td className="px-4 py-3 whitespace-nowrap">
                        {visit.visit_type === "entry" ? (
                          <span className="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-green-100 text-green-800">
                            <DoorOpen size={12} className="mr-1" />
                            Вход
                          </span>
                        ) : (
                          <span className="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-gray-100 text-gray-800">
                            <DoorClosed size={12} className="mr-1" />
                            Выход
                          </span>
                        )}
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          ) : (
            <p className="text-gray-500 text-center py-4">
              Записей о посещениях пока нет
            </p>
          )}
        </div>
      </motion.div>
    </div>
  );
};

export default RoomEntry;
