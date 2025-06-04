window.utils = {
  formatDate(dateString, options = {}) {
    if (!dateString) return "";

    const defaultOptions = {
      year: "numeric",
      month: "short",
      day: "numeric",
    };

    return new Date(dateString).toLocaleDateString("ru-RU", {
      ...defaultOptions,
      ...options,
    });
  },

  formatDateTime(dateString) {
    return this.formatDate(dateString, {
      hour: "2-digit",
      minute: "2-digit",
    });
  },

  debounce(func, wait) {
    let timeout;
    return function executedFunction(...args) {
      const later = () => {
        clearTimeout(timeout);
        func(...args);
      };
      clearTimeout(timeout);
      timeout = setTimeout(later, wait);
    };
  },

  showNotification(message, type = "info") {
    // Simple notification system
    const notification = document.createElement("div");
    const bgColor =
      {
        info: "bg-blue-500",
        success: "bg-green-500",
        warning: "bg-yellow-500",
        error: "bg-red-500",
      }[type] || "bg-blue-500";

    notification.className = `fixed top-4 right-4 ${bgColor} text-white px-6 py-3 rounded-lg shadow-lg z-50 transition-all transform translate-x-full`;
    notification.textContent = message;

    document.body.appendChild(notification);

    // Slide in
    setTimeout(() => {
      notification.classList.remove("translate-x-full");
    }, 100);

    // Remove after 3 seconds
    setTimeout(() => {
      notification.classList.add("translate-x-full");
      setTimeout(() => {
        document.body.removeChild(notification);
      }, 300);
    }, 3000);
  },
};
