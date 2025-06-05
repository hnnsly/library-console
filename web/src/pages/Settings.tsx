import React, { useState } from "react";
import { motion } from "framer-motion";
import { Save, Bell, BookOpen } from "lucide-react";

const Settings: React.FC = () => {
  // Library settings
  const [libraryName, setLibraryName] = useState("Central City Library");
  const [checkoutDuration, setCheckoutDuration] = useState(14);
  const [maxCheckouts, setMaxCheckouts] = useState(5);
  const [dailyLateFee, setDailyLateFee] = useState(0.5);

  // Notification settings
  const [emailNotifications, setEmailNotifications] = useState(true);
  const [dueDateReminders, setDueDateReminders] = useState(true);
  const [overdueNotifications, setOverdueNotifications] = useState(true);
  const [reservationNotifications, setReservationNotifications] =
    useState(true);

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    // Save settings (would connect to an API in a real app)
    alert("Settings saved successfully!");
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
          Settings
        </h1>
        <p className="text-gray-600 mt-1">
          Configure your library system preferences.
        </p>
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
            Library Settings
          </div>
          <div className="p-6">
            <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
              <div>
                <label htmlFor="libraryName" className="label">
                  Library Name
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
                  Default Checkout Duration (days)
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
                  Maximum Books Per Member
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
                  Daily Late Fee ($)
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
            <Bell size={18} className="mr-2" />
            Notification Settings
          </div>
          <div className="p-6">
            <div className="mb-6">
              <div className="flex items-center justify-between py-3 border-b">
                <div>
                  <h3 className="font-medium text-gray-900">
                    Email Notifications
                  </h3>
                  <p className="text-sm text-gray-600">
                    Enable or disable all email notifications
                  </p>
                </div>
                <div className="relative inline-block w-12 h-6 transition duration-200 ease-in-out">
                  <input
                    type="checkbox"
                    id="emailNotifications"
                    className="opacity-0 absolute w-0 h-0"
                    checked={emailNotifications}
                    onChange={() => setEmailNotifications(!emailNotifications)}
                  />
                  <label
                    htmlFor="emailNotifications"
                    className={`block overflow-hidden cursor-pointer rounded-full h-6 ${
                      emailNotifications ? "bg-primary-500" : "bg-gray-300"
                    }`}
                  >
                    <span
                      className={`block h-6 w-6 rounded-full bg-white transform transition-transform duration-200 ease-in ${
                        emailNotifications ? "translate-x-6" : "translate-x-0"
                      }`}
                    />
                  </label>
                </div>
              </div>

              <div className="flex items-center justify-between py-3 border-b">
                <div>
                  <h3 className="font-medium text-gray-900">
                    Due Date Reminders
                  </h3>
                  <p className="text-sm text-gray-600">
                    Send reminder 2 days before due date
                  </p>
                </div>
                <div className="relative inline-block w-12 h-6 transition duration-200 ease-in-out">
                  <input
                    type="checkbox"
                    id="dueDateReminders"
                    className="opacity-0 absolute w-0 h-0"
                    checked={dueDateReminders}
                    onChange={() => setDueDateReminders(!dueDateReminders)}
                    disabled={!emailNotifications}
                  />
                  <label
                    htmlFor="dueDateReminders"
                    className={`block overflow-hidden cursor-pointer rounded-full h-6 ${
                      dueDateReminders && emailNotifications
                        ? "bg-primary-500"
                        : "bg-gray-300"
                    } ${!emailNotifications ? "opacity-50 cursor-not-allowed" : ""}`}
                  >
                    <span
                      className={`block h-6 w-6 rounded-full bg-white transform transition-transform duration-200 ease-in ${
                        dueDateReminders && emailNotifications
                          ? "translate-x-6"
                          : "translate-x-0"
                      }`}
                    />
                  </label>
                </div>
              </div>

              <div className="flex items-center justify-between py-3 border-b">
                <div>
                  <h3 className="font-medium text-gray-900">
                    Overdue Notifications
                  </h3>
                  <p className="text-sm text-gray-600">
                    Send notification when book is overdue
                  </p>
                </div>
                <div className="relative inline-block w-12 h-6 transition duration-200 ease-in-out">
                  <input
                    type="checkbox"
                    id="overdueNotifications"
                    className="opacity-0 absolute w-0 h-0"
                    checked={overdueNotifications}
                    onChange={() =>
                      setOverdueNotifications(!overdueNotifications)
                    }
                    disabled={!emailNotifications}
                  />
                  <label
                    htmlFor="overdueNotifications"
                    className={`block overflow-hidden cursor-pointer rounded-full h-6 ${
                      overdueNotifications && emailNotifications
                        ? "bg-primary-500"
                        : "bg-gray-300"
                    } ${!emailNotifications ? "opacity-50 cursor-not-allowed" : ""}`}
                  >
                    <span
                      className={`block h-6 w-6 rounded-full bg-white transform transition-transform duration-200 ease-in ${
                        overdueNotifications && emailNotifications
                          ? "translate-x-6"
                          : "translate-x-0"
                      }`}
                    />
                  </label>
                </div>
              </div>

              <div className="flex items-center justify-between py-3">
                <div>
                  <h3 className="font-medium text-gray-900">
                    Reservation Notifications
                  </h3>
                  <p className="text-sm text-gray-600">
                    Send notification when reserved book is available
                  </p>
                </div>
                <div className="relative inline-block w-12 h-6 transition duration-200 ease-in-out">
                  <input
                    type="checkbox"
                    id="reservationNotifications"
                    className="opacity-0 absolute w-0 h-0"
                    checked={reservationNotifications}
                    onChange={() =>
                      setReservationNotifications(!reservationNotifications)
                    }
                    disabled={!emailNotifications}
                  />
                  <label
                    htmlFor="reservationNotifications"
                    className={`block overflow-hidden cursor-pointer rounded-full h-6 ${
                      reservationNotifications && emailNotifications
                        ? "bg-primary-500"
                        : "bg-gray-300"
                    } ${!emailNotifications ? "opacity-50 cursor-not-allowed" : ""}`}
                  >
                    <span
                      className={`block h-6 w-6 rounded-full bg-white transform transition-transform duration-200 ease-in ${
                        reservationNotifications && emailNotifications
                          ? "translate-x-6"
                          : "translate-x-0"
                      }`}
                    />
                  </label>
                </div>
              </div>
            </div>
          </div>
        </motion.div>

        <div className="flex justify-end">
          <button type="submit" className="btn btn-primary flex items-center">
            <Save size={18} className="mr-2" /> Save Settings
          </button>
        </div>
      </motion.form>
    </div>
  );
};

export default Settings;
