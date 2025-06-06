import React, { useState } from "react";
import { motion } from "framer-motion";
import {
  Save,
  BookOpen,
  Key,
  Lock,
  Plus,
  Edit2,
  Trash2,
  RotateCcw,
} from "lucide-react";

interface User {
  id: number;
  name: string;
  email: string;
  role: "admin" | "librarian";
  createdAt: string;
}

const Settings: React.FC = () => {
  // Current user role (в реальном приложении получать из контекста/стора)
  const [currentUserRole] = useState<"admin" | "librarian">("admin");

  // Library settings
  const [libraryName, setLibraryName] = useState("Центральная библиотека");
  const [checkoutDuration, setCheckoutDuration] = useState(14);
  const [maxCheckouts, setMaxCheckouts] = useState(5);
  const [dailyLateFee, setDailyLateFee] = useState(0.5);

  // Users management
  const [users, setUsers] = useState<User[]>([
    {
      id: 1,
      name: "Анна Петрова",
      email: "admin@libr.ru",
      role: "admin",
      createdAt: "2024-01-15",
    },
    {
      id: 2,
      name: "Михаил Сидоров",
      email: "librarian1@libr.ru",
      role: "librarian",
      createdAt: "2024-02-01",
    },
    {
      id: 3,
      name: "Елена Васильева",
      email: "librarian2@libr.ru",
      role: "librarian",
      createdAt: "2024-02-15",
    },
  ]);

  const [showUserModal, setShowUserModal] = useState(false);
  const [editingUser, setEditingUser] = useState<User | null>(null);
  const [newUser, setNewUser] = useState({
    name: "",
    email: "",
    role: "librarian" as "admin" | "librarian",
  });

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    alert("Настройки сохранены успешно!");
  };

  const handlePasswordChange = () => {
    alert("Функция смены пароля будет доступна в следующей версии libr");
  };

  const handlePasswordReset = () => {
    alert("Функция сброса пароля будет доступна в следующей версии libr");
  };

  const handleCreateUser = () => {
    setEditingUser(null);
    setNewUser({ name: "", email: "", role: "librarian" });
    setShowUserModal(true);
  };

  const handleEditUser = (user: User) => {
    setEditingUser(user);
    setNewUser({ name: user.name, email: user.email, role: user.role });
    setShowUserModal(true);
  };

  const handleDeleteUser = (userId: number) => {
    if (confirm("Вы уверены, что хотите удалить этого пользователя?")) {
      setUsers(users.filter((user) => user.id !== userId));
    }
  };

  const handleSaveUser = () => {
    if (editingUser) {
      setUsers(
        users.map((user) =>
          user.id === editingUser.id
            ? {
                ...user,
                name: newUser.name,
                email: newUser.email,
                role: newUser.role,
              }
            : user,
        ),
      );
    } else {
      const newId = Math.max(...users.map((u) => u.id)) + 1;
      setUsers([
        ...users,
        {
          id: newId,
          name: newUser.name,
          email: newUser.email,
          role: newUser.role,
          createdAt: new Date().toISOString().split("T")[0],
        },
      ]);
    }
    setShowUserModal(false);
  };

  const handleResetUserPassword = (userId: number, userName: string) => {
    if (confirm(`Сбросить пароль для пользователя ${userName}?`)) {
      alert("Новый пароль отправлен на email пользователя");
    }
  };

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
    <div className="max-w-5xl mx-auto">
      <div className="mb-6">
        <h1 className="text-3xl font-serif font-bold text-gray-900">
          Настройки
        </h1>
        <p className="text-gray-600 mt-1">Настройте параметры системы libr.</p>
      </div>

      <motion.form
        variants={containerVariants}
        initial="hidden"
        animate="visible"
        onSubmit={handleSubmit}
      >
        <motion.div
          variants={itemVariants}
          className="bg-white rounded-lg shadow-sm mb-6 overflow-hidden"
        >
          <div className="p-4 bg-primary-500 text-white font-medium flex items-center">
            <BookOpen size={18} className="mr-2" />
            Настройки библиотеки
          </div>
          <div className="p-6">
            <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
              <div>
                <label htmlFor="libraryName" className="label">
                  Название библиотеки
                </label>
                <input
                  id="libraryName"
                  type="text"
                  value={libraryName}
                  onChange={(e) => setLibraryName(e.target.value)}
                  className="input py-2"
                />
              </div>

              <div>
                <label htmlFor="checkoutDuration" className="label">
                  Срок выдачи по умолчанию (дни)
                </label>
                <input
                  id="checkoutDuration"
                  type="number"
                  min="1"
                  max="60"
                  value={checkoutDuration}
                  onChange={(e) =>
                    setCheckoutDuration(parseInt(e.target.value))
                  }
                  className="input py-2"
                />
              </div>

              <div>
                <label htmlFor="maxCheckouts" className="label">
                  Максимум книг на читателя
                </label>
                <input
                  id="maxCheckouts"
                  type="number"
                  min="1"
                  max="20"
                  value={maxCheckouts}
                  onChange={(e) => setMaxCheckouts(parseInt(e.target.value))}
                  className="input py-2"
                />
              </div>

              <div>
                <label htmlFor="dailyLateFee" className="label">
                  Штраф за день просрочки (₽)
                </label>
                <input
                  id="dailyLateFee"
                  type="number"
                  min="0"
                  step="0.01"
                  value={dailyLateFee}
                  onChange={(e) => setDailyLateFee(parseFloat(e.target.value))}
                  className="input py-2"
                />
              </div>
            </div>
          </div>
        </motion.div>

        <motion.div
          variants={itemVariants}
          className="bg-white rounded-lg shadow-sm mb-6 overflow-hidden"
        >
          <div className="p-4 bg-accent-500 text-white font-medium flex items-center">
            <Lock size={18} className="mr-2" />
            Безопасность
          </div>
          <div className="p-6">
            <div className="flex flex-col sm:flex-row gap-4">
              <button
                type="button"
                onClick={handlePasswordChange}
                className="btn btn-outline flex items-center justify-center"
              >
                <Key size={18} className="mr-2" />
                Сменить пароль
              </button>
              <button
                type="button"
                onClick={handlePasswordReset}
                className="btn btn-outline flex items-center justify-center"
              >
                <RotateCcw size={18} className="mr-2" />
                Сбросить пароль
              </button>
            </div>
          </div>
        </motion.div>

        {currentUserRole === "admin" && (
          <motion.div
            variants={itemVariants}
            className="bg-white rounded-lg shadow-sm mb-6 overflow-hidden"
          >
            <div className="p-4 bg-green-500 text-white font-medium flex items-center justify-between">
              <span className="flex items-center">
                <Edit2 size={18} className="mr-2" />
                Управление пользователями
              </span>
              <button
                type="button"
                onClick={handleCreateUser}
                className="bg-white text-green-500 px-3 py-1 rounded text-sm font-medium hover:bg-green-50 transition-colors flex items-center"
              >
                <Plus size={16} className="mr-1" />
                Добавить пользователя
              </button>
            </div>
            <div className="p-6">
              <div className="overflow-x-auto">
                <table className="w-full">
                  <thead>
                    <tr className="border-b border-gray-200">
                      <th className="text-left py-3 px-4 font-medium text-gray-900">
                        Имя
                      </th>
                      <th className="text-left py-3 px-4 font-medium text-gray-900">
                        Email
                      </th>
                      <th className="text-left py-3 px-4 font-medium text-gray-900">
                        Роль
                      </th>
                      <th className="text-left py-3 px-4 font-medium text-gray-900">
                        Дата создания
                      </th>
                      <th className="text-left py-3 px-4 font-medium text-gray-900">
                        Действия
                      </th>
                    </tr>
                  </thead>
                  <tbody>
                    {users.map((user) => (
                      <tr
                        key={user.id}
                        className="border-b border-gray-100 hover:bg-gray-50"
                      >
                        <td className="py-3 px-4">{user.name}</td>
                        <td className="py-3 px-4 text-gray-600">
                          {user.email}
                        </td>
                        <td className="py-3 px-4">
                          <span
                            className={`px-2 py-1 rounded-full text-xs font-medium ${
                              user.role === "admin"
                                ? "bg-red-100 text-red-800"
                                : "bg-blue-100 text-blue-800"
                            }`}
                          >
                            {user.role === "admin"
                              ? "Администратор"
                              : "Библиотекарь"}
                          </span>
                        </td>
                        <td className="py-3 px-4 text-gray-600">
                          {new Date(user.createdAt).toLocaleDateString("ru-RU")}
                        </td>
                        <td className="py-3 px-4">
                          <div className="flex items-center space-x-2">
                            <button
                              type="button"
                              onClick={() => handleEditUser(user)}
                              className="text-blue-600 hover:text-blue-800 transition-colors"
                              title="Редактировать"
                            >
                              <Edit2 size={16} />
                            </button>
                            <button
                              type="button"
                              onClick={() =>
                                handleResetUserPassword(user.id, user.name)
                              }
                              className="text-orange-600 hover:text-orange-800 transition-colors"
                              title="Сбросить пароль"
                            >
                              <RotateCcw size={16} />
                            </button>
                            <button
                              type="button"
                              onClick={() => handleDeleteUser(user.id)}
                              className="text-red-600 hover:text-red-800 transition-colors"
                              title="Удалить"
                            >
                              <Trash2 size={16} />
                            </button>
                          </div>
                        </td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>
            </div>
          </motion.div>
        )}

        <div className="flex justify-end">
          <button type="submit" className="btn btn-primary flex items-center">
            <Save size={18} className="mr-2" /> Сохранить настройки
          </button>
        </div>
      </motion.form>

      {/* Modal for creating/editing users */}
      {showUserModal && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
          <div className="bg-white rounded-lg shadow-xl p-6 w-full max-w-md mx-4">
            <h3 className="text-lg font-medium text-gray-900 mb-4">
              {editingUser
                ? "Редактировать пользователя"
                : "Создать пользователя"}
            </h3>
            <div className="space-y-4">
              <div>
                <label className="label">Имя</label>
                <input
                  type="text"
                  value={newUser.name}
                  onChange={(e) =>
                    setNewUser({ ...newUser, name: e.target.value })
                  }
                  className="input py-2 w-full"
                  placeholder="Введите имя пользователя"
                />
              </div>
              <div>
                <label className="label">Email</label>
                <input
                  type="email"
                  value={newUser.email}
                  onChange={(e) =>
                    setNewUser({ ...newUser, email: e.target.value })
                  }
                  className="input py-2 w-full"
                  placeholder="Введите email"
                />
              </div>
              <div>
                <label className="label">Роль</label>
                <select
                  value={newUser.role}
                  onChange={(e) =>
                    setNewUser({
                      ...newUser,
                      role: e.target.value as "admin" | "librarian",
                    })
                  }
                  className="input py-2 w-full"
                >
                  <option value="librarian">Библиотекарь</option>
                  <option value="admin">Администратор</option>
                </select>
              </div>
            </div>
            <div className="flex justify-end space-x-3 mt-6">
              <button
                type="button"
                onClick={() => setShowUserModal(false)}
                className="btn btn-outline"
              >
                Отмена
              </button>
              <button
                type="button"
                onClick={handleSaveUser}
                className="btn btn-primary"
                disabled={!newUser.name || !newUser.email}
              >
                {editingUser ? "Сохранить" : "Создать"}
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
};

export default Settings;
