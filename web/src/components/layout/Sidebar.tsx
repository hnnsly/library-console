import React, { useState } from "react";
import { NavLink } from "react-router-dom";
import { motion } from "framer-motion";
import {
  BookOpen,
  Users,
  Home,
  Settings,
  LogOut,
  Menu,
  X,
  BookPlus,
  BookCheck,
  DoorOpen,
  BarChart3,
  User,
} from "lucide-react";

const Sidebar: React.FC = () => {
  const [isMobileMenuOpen, setIsMobileMenuOpen] = useState(false);

  const toggleMobileMenu = () => {
    setIsMobileMenuOpen(!isMobileMenuOpen);
  };

  const navItems = [
    { to: "/", icon: <Home size={20} />, label: "Главная" },
    { to: "/books", icon: <BookOpen size={20} />, label: "Книги" },
    { to: "/readers", icon: <Users size={20} />, label: "Читатели" },
    { to: "/checkout", icon: <BookPlus size={20} />, label: "Выдача книг" },
    { to: "/returns", icon: <BookCheck size={20} />, label: "Возврат книг" },
    {
      to: "/room-entry",
      icon: <DoorOpen size={20} />,
      label: "Посещения залов",
    },
    {
      to: "/room-overview",
      icon: <BarChart3 size={20} />,
      label: "Обзор залов",
    },
    { to: "/settings", icon: <Settings size={20} />, label: "Настройки" },
  ];

  const sidebarVariants = {
    open: { x: 0 },
    closed: { x: "-100%" },
  };

  const itemVariants = {
    open: {
      opacity: 1,
      y: 0,
      transition: { type: "spring", stiffness: 300, damping: 24 },
    },
    closed: { opacity: 0, y: 20, transition: { duration: 0.2 } },
  };

  return (
    <>
      <button
        className="md:hidden fixed top-4 left-4 z-50 bg-primary-500 text-white p-2 rounded-md"
        onClick={toggleMobileMenu}
      >
        {isMobileMenuOpen ? <X size={24} /> : <Menu size={24} />}
      </button>

      <motion.nav
        className="md:hidden fixed inset-0 z-40 bg-primary-500 text-white w-64 p-5 shadow-lg"
        variants={sidebarVariants}
        initial="closed"
        animate={isMobileMenuOpen ? "open" : "closed"}
        transition={{ type: "spring", stiffness: 300, damping: 30 }}
      >
        <div className="pt-16">
          <div className="mb-8">
            <h2 className="text-2xl font-serif font-bold text-white">libr</h2>
            <p className="text-sm text-white/70 mt-1">
              Система управления библиотекой
            </p>
          </div>
          <ul className="space-y-2">
            {navItems.map((item, index) => (
              <motion.li key={index} variants={itemVariants}>
                <NavLink
                  to={item.to}
                  className={({ isActive }) =>
                    `flex items-center p-3 rounded-md transition-all ${
                      isActive
                        ? "bg-primary-600 text-white"
                        : "text-white/80 hover:bg-primary-600 hover:text-white"
                    }`
                  }
                  onClick={() => setIsMobileMenuOpen(false)}
                >
                  <span className="mr-3">{item.icon}</span>
                  <span>{item.label}</span>
                </NavLink>
              </motion.li>
            ))}
          </ul>
        </div>
      </motion.nav>

      <nav className="hidden md:flex flex-col w-64 bg-primary-500 text-white p-5 shadow-lg">
        <div className="mb-8">
          <h2 className="text-2xl font-serif font-bold text-white">libr</h2>
          <p className="text-sm text-white/70 mt-1">
            Система управления библиотекой
          </p>
        </div>
        <ul className="space-y-2">
          {navItems.map((item, index) => (
            <li key={index}>
              <NavLink
                to={item.to}
                className={({ isActive }) =>
                  `flex items-center p-3 rounded-md transition-all ${
                    isActive
                      ? "bg-primary-600 text-white"
                      : "text-white/80 hover:bg-primary-600 hover:text-white"
                  }`
                }
              >
                <span className="mr-3">{item.icon}</span>
                <span>{item.label}</span>
              </NavLink>
            </li>
          ))}
        </ul>
        <div className="mt-auto">
          <button className="flex items-center p-3 w-full text-white/80 hover:bg-primary-600 hover:text-white rounded-md transition-all">
            <LogOut size={20} className="mr-3" />
            <span>Выход</span>
          </button>
        </div>
      </nav>
    </>
  );
};

export default Sidebar;
