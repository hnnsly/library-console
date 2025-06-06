import React, { useState } from "react";
import { motion } from "framer-motion";
import {
  Search,
  Calendar,
  Check,
  AlertTriangle,
  BookOpen,
  Book,
} from "lucide-react";
import {
  booksWithDetails,
  readers,
  bookIssues,
  bookCopies,
  fines,
} from "../data/mockData";
import { format, isPast, parseISO } from "date-fns";
import { ru } from "date-fns/locale";
import type { BookIssue, Reader, BookWithDetails, BookCopy } from "../types";

interface ReturnData extends BookIssue {
  book?: BookWithDetails;
  reader?: Reader;
  bookCopy?: BookCopy;
}

const Returns: React.FC = () => {
  const [searchTerm, setSearchTerm] = useState("");
  const [selectedIssue, setSelectedIssue] = useState<ReturnData | null>(null);
  const [bookCondition, setBookCondition] = useState<"good" | "damaged">(
    "good",
  );
  const [fineAmount, setFineAmount] = useState<number>(0);
  const [returnComplete, setReturnComplete] = useState(false);

  // Получаем активные выдачи с информацией о книгах и читателях
  const activeIssues: ReturnData[] = bookIssues
    .filter((issue) => !issue.return_date)
    .map((issue) => {
      const bookCopy = bookCopies.find(
        (copy) => copy.id === issue.book_copy_id,
      );
      const book = booksWithDetails.find((b) => b.id === bookCopy?.book_id);
      const reader = readers.find((r) => r.id === issue.reader_id);

      return {
        ...issue,
        book,
        reader,
        bookCopy,
      };
    })
    .filter((issue) => issue.book && issue.reader); // Убираем записи с недостающими данными

  // Фильтруем выдачи по поисковому запросу
  const filteredIssues = activeIssues.filter((issue) => {
    const bookMatches =
      issue.book?.title.toLowerCase().includes(searchTerm.toLowerCase()) ||
      issue.book?.authors.some((author) =>
        author.full_name.toLowerCase().includes(searchTerm.toLowerCase()),
      );

    const readerMatches =
      issue.reader?.full_name
        .toLowerCase()
        .includes(searchTerm.toLowerCase()) ||
      issue.reader?.email?.toLowerCase().includes(searchTerm.toLowerCase()) ||
      issue.reader?.ticket_number.includes(searchTerm);

    const copyMatches = issue.bookCopy?.copy_code
      .toLowerCase()
      .includes(searchTerm.toLowerCase());

    return bookMatches || readerMatches || copyMatches;
  });

  const handleReturnSubmit = () => {
    if (selectedIssue) {
      // В реальном приложении это обновило бы базу данных
      setReturnComplete(true);

      // Сбрасываем форму через 3 секунды
      setTimeout(() => {
        setSelectedIssue(null);
        setBookCondition("good");
        setFineAmount(0);
        setReturnComplete(false);
      }, 3000);
    }
  };

  const calculateLateFee = (dueDate: string): number => {
    if (!isPast(parseISO(dueDate))) return 0;

    const today = new Date();
    const due = new Date(dueDate);
    const daysLate = Math.floor(
      (today.getTime() - due.getTime()) / (1000 * 60 * 60 * 24),
    );

    // 50 рублей за день просрочки
    return daysLate * 50;
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
        <h1 className="text-3xl font-serif font-bold text-gray-900">
          Возврат книг
        </h1>
        <p className="text-gray-600 mt-1">
          Обработка возврата книг и начисление штрафов в системе libr.
        </p>
      </div>

      <motion.div
        variants={containerVariants}
        initial="hidden"
        animate="visible"
        className="bg-white rounded-lg shadow-sm overflow-hidden"
      >
        {returnComplete ? (
          <motion.div
            initial={{ opacity: 0, scale: 0.9 }}
            animate={{ opacity: 1, scale: 1 }}
            className="p-8 text-center"
          >
            <div className="w-16 h-16 bg-green-100 rounded-full flex items-center justify-center mx-auto mb-4">
              <Check size={32} className="text-green-600" />
            </div>
            <h2 className="text-2xl font-bold text-gray-900 mb-2">
              Возврат завершен!
            </h2>
            <p className="text-gray-600 mb-2">
              <span className="font-medium">
                «{selectedIssue?.book?.title}»
              </span>{" "}
              (экземпляр {selectedIssue?.bookCopy?.copy_code}) возвращена
              читателем{" "}
              <span className="font-medium">
                {selectedIssue?.reader?.full_name}
              </span>
              .
            </p>
            {fineAmount > 0 && (
              <p className="text-orange-700 mb-6">
                Начислен штраф в размере {fineAmount.toFixed(0)} ₽.
              </p>
            )}
            <p className="text-gray-700">
              Состояние книги:{" "}
              <span className="font-medium">
                {bookCondition === "good" ? "Хорошее" : "Повреждена"}
              </span>
            </p>
          </motion.div>
        ) : (
          <div>
            <div className="p-6">
              <div className="mb-6">
                <label htmlFor="search" className="label">
                  Поиск книги, читателя или экземпляра
                </label>
                <div className="relative">
                  <Search
                    className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400"
                    size={18}
                  />
                  <input
                    id="search"
                    type="text"
                    placeholder="Введите название книги, автора, имя читателя, номер билета или код экземпляра..."
                    value={searchTerm}
                    onChange={(e) => setSearchTerm(e.target.value)}
                    className="input pl-10 py-2"
                  />
                </div>
              </div>

              <motion.div variants={itemVariants} className="mb-6">
                <h2 className="text-xl font-semibold mb-4 flex items-center">
                  <BookOpen size={20} className="mr-2 text-primary-500" />
                  Активные выдачи ({activeIssues.length})
                </h2>

                {filteredIssues.length > 0 ? (
                  <div className="overflow-x-auto border rounded-lg">
                    <table className="min-w-full divide-y divide-gray-200">
                      <thead className="bg-gray-50">
                        <tr>
                          <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                            Книга
                          </th>
                          <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                            Экземпляр
                          </th>
                          <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                            Читатель
                          </th>
                          <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                            Дата выдачи
                          </th>
                          <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                            Дата возврата
                          </th>
                          <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                            Статус
                          </th>
                          <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">
                            Действие
                          </th>
                        </tr>
                      </thead>
                      <tbody className="bg-white divide-y divide-gray-200">
                        {filteredIssues.map((issue) => {
                          const isOverdue = isPast(new Date(issue.due_date));

                          return (
                            <tr
                              key={issue.id}
                              className={`hover:bg-gray-50 cursor-pointer ${
                                selectedIssue?.id === issue.id
                                  ? "bg-primary-50"
                                  : ""
                              }`}
                              onClick={() => {
                                setSelectedIssue(issue);
                                const lateFee = calculateLateFee(
                                  typeof issue.due_date === "string"
                                    ? issue.due_date
                                    : issue.due_date.toISOString(),
                                );
                                setFineAmount(lateFee);
                                setBookCondition("good");
                              }}
                            >
                              <td className="px-6 py-4">
                                <div className="flex items-center">
                                  <div className="h-10 w-7 bg-gradient-to-br from-primary-50 to-primary-100 rounded overflow-hidden mr-3 flex-shrink-0 flex items-center justify-center">
                                    <Book
                                      size={14}
                                      className="text-primary-400"
                                    />
                                  </div>
                                  <div>
                                    <div className="font-medium text-gray-900">
                                      {issue.book?.title}
                                    </div>
                                    <div className="text-sm text-gray-500">
                                      {issue.book?.authors
                                        .map((author) => author.full_name)
                                        .join(", ")}
                                    </div>
                                  </div>
                                </div>
                              </td>
                              <td className="px-6 py-4 whitespace-nowrap">
                                <div className="text-sm font-medium text-gray-900">
                                  {issue.bookCopy?.copy_code}
                                </div>
                                {issue.bookCopy?.location_info && (
                                  <div className="text-xs text-gray-500">
                                    {issue.bookCopy.location_info}
                                  </div>
                                )}
                              </td>
                              <td className="px-6 py-4 whitespace-nowrap">
                                <div className="text-sm font-medium text-gray-900">
                                  {issue.reader?.full_name}
                                </div>
                                <div className="text-xs text-gray-500">
                                  Билет: {issue.reader?.ticket_number}
                                </div>
                              </td>
                              <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                                {format(
                                  new Date(issue.issue_date),
                                  "d MMM yyyy",
                                  { locale: ru },
                                )}
                              </td>
                              <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                                {format(
                                  new Date(issue.due_date),
                                  "d MMM yyyy",
                                  { locale: ru },
                                )}
                              </td>
                              <td className="px-6 py-4 whitespace-nowrap">
                                {isOverdue ? (
                                  <span className="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-red-100 text-red-800">
                                    <AlertTriangle size={12} className="mr-1" />
                                    Просрочена
                                  </span>
                                ) : (
                                  <span className="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-green-100 text-green-800">
                                    <Check size={12} className="mr-1" /> В срок
                                  </span>
                                )}
                              </td>
                              <td className="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
                                <button
                                  className="text-primary-600 hover:text-primary-900"
                                  onClick={(e) => {
                                    e.stopPropagation();
                                    setSelectedIssue(issue);
                                    const lateFee = calculateLateFee(
                                      typeof issue.due_date === "string"
                                        ? issue.due_date
                                        : issue.due_date.toISOString(),
                                    );
                                    setFineAmount(lateFee);
                                    setBookCondition("good");
                                  }}
                                >
                                  Оформить возврат
                                </button>
                              </td>
                            </tr>
                          );
                        })}
                      </tbody>
                    </table>
                  </div>
                ) : (
                  <div className="text-center py-8 border rounded-lg">
                    <p className="text-gray-500">
                      {searchTerm
                        ? "Подходящие выдачи не найдены."
                        : "Нет активных выдач для отображения."}
                    </p>
                  </div>
                )}
              </motion.div>

              {selectedIssue && (
                <motion.div
                  initial={{ opacity: 0, height: 0 }}
                  animate={{ opacity: 1, height: "auto" }}
                  transition={{ duration: 0.3 }}
                  className="border rounded-lg p-6 bg-gray-50"
                >
                  <h2 className="text-xl font-semibold mb-4">
                    Оформление возврата
                  </h2>

                  <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                    <div>
                      <h3 className="font-medium text-gray-800 mb-3">
                        Информация о книге
                      </h3>
                      <div className="flex items-start">
                        <div className="h-24 w-16 bg-gradient-to-br from-primary-100 to-primary-200 rounded overflow-hidden mr-4 flex-shrink-0 flex items-center justify-center">
                          <Book size={24} className="text-primary-500" />
                        </div>
                        <div>
                          <p className="font-medium text-gray-900">
                            {selectedIssue.book?.title}
                          </p>
                          <p className="text-sm text-gray-600 mb-1">
                            {selectedIssue.book?.authors
                              .map((author) => author.full_name)
                              .join(", ")}
                          </p>
                          {selectedIssue.book?.isbn && (
                            <p className="text-sm text-gray-600">
                              ISBN: {selectedIssue.book.isbn}
                            </p>
                          )}
                          <p className="text-sm text-gray-600 mt-2">
                            Экземпляр:{" "}
                            <span className="font-medium">
                              {selectedIssue.bookCopy?.copy_code}
                            </span>
                          </p>
                          <p className="text-sm text-gray-600">
                            Срок возврата:{" "}
                            {format(
                              new Date(selectedIssue.due_date),
                              "d MMMM yyyy",
                              { locale: ru },
                            )}
                          </p>
                          {isPast(new Date(selectedIssue.due_date)) && (
                            <p className="text-sm text-red-600 font-medium mt-1">
                              Книга просрочена!
                            </p>
                          )}
                        </div>
                      </div>

                      <div className="mt-4 p-3 bg-blue-50 rounded-lg">
                        <h4 className="font-medium text-gray-800 mb-2">
                          Читатель
                        </h4>
                        <p className="text-sm text-gray-700">
                          {selectedIssue.reader?.full_name}
                        </p>
                        <p className="text-xs text-gray-600">
                          Билет: {selectedIssue.reader?.ticket_number}
                        </p>
                        {selectedIssue.reader?.email && (
                          <p className="text-xs text-gray-600">
                            {selectedIssue.reader.email}
                          </p>
                        )}
                      </div>
                    </div>

                    <div>
                      <h3 className="font-medium text-gray-800 mb-3">
                        Детали возврата
                      </h3>

                      <div className="mb-4">
                        <label className="label">Состояние книги</label>
                        <div className="flex space-x-4">
                          <label className="flex items-center">
                            <input
                              type="radio"
                              name="condition"
                              value="good"
                              checked={bookCondition === "good"}
                              onChange={() => {
                                setBookCondition("good");
                                const lateFee = calculateLateFee(
                                  typeof selectedIssue.due_date === "string"
                                    ? selectedIssue.due_date
                                    : selectedIssue.due_date.toISOString(),
                                );
                                setFineAmount(lateFee);
                              }}
                              className="mr-2"
                            />
                            <span>Хорошее</span>
                          </label>
                          <label className="flex items-center">
                            <input
                              type="radio"
                              name="condition"
                              value="damaged"
                              checked={bookCondition === "damaged"}
                              onChange={() => {
                                setBookCondition("damaged");
                                // Добавляем штраф за повреждение
                                const lateFee = calculateLateFee(
                                  typeof selectedIssue.due_date === "string"
                                    ? selectedIssue.due_date
                                    : selectedIssue.due_date.toISOString(),
                                );
                                setFineAmount(lateFee + 500); // 500 рублей штраф за повреждение
                              }}
                              className="mr-2"
                            />
                            <span>Повреждена</span>
                          </label>
                        </div>
                      </div>

                      <div className="mb-4">
                        <label htmlFor="fineAmount" className="label">
                          Сумма штрафа (₽)
                        </label>
                        <input
                          id="fineAmount"
                          type="number"
                          step="1"
                          min="0"
                          value={fineAmount}
                          onChange={(e) =>
                            setFineAmount(parseFloat(e.target.value) || 0)
                          }
                          className="input py-2"
                        />
                        {isPast(new Date(selectedIssue.due_date)) && (
                          <p className="text-xs text-gray-500 mt-1">
                            Включая штраф за просрочку:{" "}
                            {calculateLateFee(
                              typeof selectedIssue.due_date === "string"
                                ? selectedIssue.due_date
                                : selectedIssue.due_date.toISOString(),
                            ).toFixed(0)}{" "}
                            ₽
                          </p>
                        )}
                        {bookCondition === "damaged" && (
                          <p className="text-xs text-gray-500 mt-1">
                            Включая штраф за повреждение: 500 ₽
                          </p>
                        )}
                      </div>
                    </div>
                  </div>
                </motion.div>
              )}
            </div>

            <div className="bg-gray-50 p-6 flex justify-end border-t">
              <button
                className={`btn flex items-center ${
                  selectedIssue
                    ? "bg-primary-500 text-white hover:bg-primary-600"
                    : "bg-gray-200 text-gray-500 cursor-not-allowed"
                }`}
                disabled={!selectedIssue}
                onClick={handleReturnSubmit}
              >
                Завершить возврат <Check size={18} className="ml-2" />
              </button>
            </div>
          </div>
        )}
      </motion.div>
    </div>
  );
};

export default Returns;
