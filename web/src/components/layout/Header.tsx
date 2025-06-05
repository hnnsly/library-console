import React, { useState, useEffect } from "react";
import { Bell, Search } from "lucide-react";
import { motion } from "framer-motion";

const Header: React.FC = () => {
  const [isScrolled, setIsScrolled] = useState(false);
  const [searchQuery, setSearchQuery] = useState("");
  const [notifications, setNotifications] = useState([
    {
      id: 1,
      message: "Книга 'Война и мир' должна быть возвращена завтра",
      isNew: true,
    },
    { id: 2, message: "Новый читатель: Анна Иванова", isNew: true },
    {
      id: 3,
      message: "Срок бронирования истек: 'Преступление и наказание'",
      isNew: false,
    },
  ]);
  const [showNotifications, setShowNotifications] = useState(false);

  const newNotificationsCount = notifications.filter((n) => n.isNew).length;

  useEffect(() => {
    const handleScroll = () => {
      setIsScrolled(window.scrollY > 10);
    };

    window.addEventListener("scroll", handleScroll);
    return () => window.removeEventListener("scroll", handleScroll);
  }, []);

  const handleSearchChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setSearchQuery(e.target.value);
  };

  const toggleNotifications = () => {
    setShowNotifications(!showNotifications);
    if (!showNotifications) {
      setNotifications(notifications.map((n) => ({ ...n, isNew: false })));
    }
  };

  return (
    <header
      className={`sticky top-0 z-30 transition-all duration-300 px-4 md:px-6 py-3 flex items-center justify-between ${
        isScrolled ? "bg-white shadow-md" : "bg-transparent"
      }`}
    >
      <div className="relative md:w-64 lg:w-96">
        <Search
          size={18}
          className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400"
        />
        <input
          type="text"
          placeholder="Поиск книг, читателей..."
          value={searchQuery}
          onChange={handleSearchChange}
          className="input pl-10 py-2"
        />
      </div>

      <div className="flex items-center space-x-4">
        <div className="relative">
          <button
            className="p-2 rounded-full hover:bg-gray-100 transition-colors relative"
            onClick={toggleNotifications}
          >
            <Bell size={20} />
            {newNotificationsCount > 0 && (
              <span className="absolute top-0 right-0 bg-error-500 text-white text-xs rounded-full w-4 h-4 flex items-center justify-center">
                {newNotificationsCount}
              </span>
            )}
          </button>

          {showNotifications && (
            <motion.div
              className="absolute right-0 mt-2 w-80 bg-white rounded-lg shadow-lg overflow-hidden z-50"
              initial={{ opacity: 0, y: -10 }}
              animate={{ opacity: 1, y: 0 }}
              exit={{ opacity: 0, y: -10 }}
              transition={{ duration: 0.2 }}
            >
              <div className="p-3 bg-primary-500 text-white font-medium">
                Уведомления
              </div>
              <div className="max-h-80 overflow-y-auto">
                {notifications.length > 0 ? (
                  <ul>
                    {notifications.map((notification) => (
                      <li
                        key={notification.id}
                        className={`p-3 border-b border-gray-100 hover:bg-gray-50 transition-colors ${
                          notification.isNew ? "bg-blue-50" : ""
                        }`}
                      >
                        <p className="text-sm text-gray-700">
                          {notification.message}
                        </p>
                      </li>
                    ))}
                  </ul>
                ) : (
                  <div className="p-4 text-center text-gray-500">
                    Нет уведомлений
                  </div>
                )}
              </div>
              <div className="p-2 bg-gray-50 text-center">
                <button className="text-sm text-primary-500 hover:underline">
                  Показать все уведомления
                </button>
              </div>
            </motion.div>
          )}
        </div>

        <div className="flex items-center">
          <div className="w-8 h-8 rounded-full bg-primary-500 text-white flex items-center justify-center mr-2">
            А
          </div>
          <span className="hidden md:block font-medium">Администратор</span>
        </div>
      </div>
    </header>
  );
};

export default Header;
