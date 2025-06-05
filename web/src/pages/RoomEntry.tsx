import React, { useState } from "react";
import { motion } from "framer-motion";
import {
  DoorOpen,
  DoorClosed,
  User,
  Clock,
  AlertCircle,
  CheckCircle,
} from "lucide-react";
import { members, recentEntryRecords } from "../data/mockData";

interface EntryRecord {
  id: string;
  memberId: string;
  memberName: string;
  entryTime: Date;
  exitTime?: Date;
  status: "in" | "out";
}

const RoomEntry: React.FC = () => {
  const [cardNumber, setCardNumber] = useState("");
  const [entryRecords, setEntryRecords] = useState(recentEntryRecords);
  const [message, setMessage] = useState<{
    type: "success" | "error";
    text: string;
  } | null>(null);

  const handleCardSubmit = (e: React.FormEvent) => {
    e.preventDefault();

    if (!cardNumber.trim()) {
      setMessage({ type: "error", text: "Введите номер читательского билета" });
      return;
    }

    // Найти читателя по номеру билета (предполагаем, что номер билета = ID)
    const member = members.find(
      (m) =>
        m.id === cardNumber ||
        m.firstName.toLowerCase().includes(cardNumber.toLowerCase()),
    );

    if (!member) {
      setMessage({
        type: "error",
        text: "Читатель с таким номером билета не найден",
      });
      return;
    }

    // Проверить, находится ли читатель уже в зале
    const activeEntry = entryRecords.find(
      (r) => r.memberId === member.id && r.status === "in",
    );

    if (activeEntry) {
      // Выход из зала
      setEntryRecords((prev) =>
        prev.map((record) =>
          record.id === activeEntry.id
            ? { ...record, exitTime: new Date(), status: "out" as const }
            : record,
        ),
      );
      setMessage({
        type: "success",
        text: `${member.firstName} ${member.lastName} покинул(а) читальный зал`,
      });
    } else {
      // Вход в зал
      const newEntry: EntryRecord = {
        id: Date.now().toString(),
        memberId: member.id,
        memberName: `${member.firstName} ${member.lastName}`,
        entryTime: new Date(),
        status: "in",
      };
      setEntryRecords((prev) => [newEntry, ...prev]);
      setMessage({
        type: "success",
        text: `${member.firstName} ${member.lastName} вошёл(ла) в читальный зал`,
      });
    }

    setCardNumber("");
    setTimeout(() => setMessage(null), 3000);
  };

  const currentVisitors = entryRecords.filter((r) => r.status === "in").length;

  return (
    <div className="max-w-4xl mx-auto">
      <div className="mb-6">
        <h1 className="text-3xl font-serif font-bold text-gray-900">
          Вход в читальный зал
        </h1>
        <p className="text-gray-600 mt-1">
          Регистрация входа и выхода читателей в системе libr
        </p>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6 mb-8">
        {/* Форма ввода */}
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          className="card p-6"
        >
          <div className="flex items-center mb-4">
            <DoorOpen className="text-primary-600 mr-3" size={24} />
            <h2 className="text-xl font-semibold">Регистрация входа/выхода</h2>
          </div>

          <form onSubmit={handleCardSubmit} className="space-y-4">
            <div>
              <label
                htmlFor="cardNumber"
                className="block text-sm font-medium text-gray-700 mb-2"
              >
                Номер читательского билета
              </label>
              <input
                type="text"
                id="cardNumber"
                value={cardNumber}
                onChange={(e) => setCardNumber(e.target.value)}
                className="input w-full"
                placeholder="Введите номер билета..."
                autoFocus
              />
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
            <h2 className="text-xl font-semibold">Текущая статистика</h2>
          </div>

          <div className="space-y-4">
            <div className="bg-green-50 p-4 rounded-lg">
              <div className="flex items-center justify-between">
                <span className="text-green-800 font-medium">
                  Посетителей в зале
                </span>
                <span className="text-2xl font-bold text-green-600">
                  {currentVisitors}
                </span>
              </div>
            </div>

            <div className="bg-blue-50 p-4 rounded-lg">
              <div className="flex items-center justify-between">
                <span className="text-blue-800 font-medium">
                  Всего записей сегодня
                </span>
                <span className="text-2xl font-bold text-blue-600">
                  {entryRecords.length}
                </span>
              </div>
            </div>
          </div>
        </motion.div>
      </div>

      {/* Последние записи */}
      <motion.div
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ delay: 0.2 }}
        className="card"
      >
        <div className="p-4 bg-primary-500 text-white">
          <h2 className="text-xl font-semibold">Последние записи</h2>
        </div>
        <div className="p-4">
          {entryRecords.length > 0 ? (
            <div className="overflow-x-auto">
              <table className="min-w-full divide-y divide-gray-200">
                <thead>
                  <tr>
                    <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                      Читатель
                    </th>
                    <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                      Время входа
                    </th>
                    <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                      Время выхода
                    </th>
                    <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                      Статус
                    </th>
                  </tr>
                </thead>
                <tbody className="divide-y divide-gray-200">
                  {entryRecords.slice(0, 10).map((record) => (
                    <tr key={record.id} className="hover:bg-gray-50">
                      <td className="px-4 py-3 whitespace-nowrap">
                        <div className="font-medium text-gray-900">
                          {record.memberName}
                        </div>
                      </td>
                      <td className="px-4 py-3 whitespace-nowrap text-gray-700">
                        {record.entryTime.toLocaleTimeString("ru-RU")}
                      </td>
                      <td className="px-4 py-3 whitespace-nowrap text-gray-700">
                        {record.exitTime
                          ? record.exitTime.toLocaleTimeString("ru-RU")
                          : "—"}
                      </td>
                      <td className="px-4 py-3 whitespace-nowrap">
                        {record.status === "in" ? (
                          <span className="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-green-100 text-green-800">
                            <DoorOpen size={12} className="mr-1" />В зале
                          </span>
                        ) : (
                          <span className="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-gray-100 text-gray-800">
                            <DoorClosed size={12} className="mr-1" />
                            Покинул(а)
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
