function dashboardComponent() {
  return {
    stats: {},
    recentOperations: [],
    isLoading: false,

    async init() {
      this.isLoading = true;
      await Promise.all([this.loadStats(), this.loadRecentOperations()]);
      this.isLoading = false;
    },

    async loadStats() {
      try {
        this.stats = await api.getDashboardStats();
      } catch (error) {
        console.error("Failed to load dashboard stats:", error);
      }
    },

    async loadRecentOperations() {
      try {
        this.recentOperations = await api.getRecentOperations({ limit: 10 });
      } catch (error) {
        console.error("Failed to load recent operations:", error);
        this.recentOperations = [];
      }
    },

    formatDate(dateString) {
      if (!dateString) return "";
      return new Date(dateString).toLocaleString("ru-RU", {
        year: "numeric",
        month: "short",
        day: "numeric",
        hour: "2-digit",
        minute: "2-digit",
      });
    },
  };
}
