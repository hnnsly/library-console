function readersComponent() {
  return {
    readers: [],
    searchQuery: "",
    isLoading: false,
    showCreateModal: false,

    async init() {
      this.isLoading = true;
      await this.loadReaders();
      this.isLoading = false;
    },

    async loadReaders() {
      try {
        this.readers = await api.getReaders({ page_limit: 50 });
      } catch (error) {
        console.error("Failed to load readers:", error);
        this.readers = [];
      }
    },

    async searchReaders() {
      if (!this.searchQuery.trim()) {
        return this.loadReaders();
      }

      this.isLoading = true;
      try {
        this.readers = await api.searchReaders({
          search_name: this.searchQuery,
          page_limit: 50,
        });
      } catch (error) {
        console.error("Failed to search readers:", error);
        this.readers = [];
      } finally {
        this.isLoading = false;
      }
    },

    getStatusColor(status) {
      const colors = {
        active: "bg-green-100 text-green-800",
        suspended: "bg-yellow-100 text-yellow-800",
        blocked: "bg-red-100 text-red-800",
        inactive: "bg-gray-100 text-gray-800",
      };
      return colors[status] || "bg-gray-100 text-gray-800";
    },

    getStatusText(status) {
      const texts = {
        active: "Активен",
        suspended: "Приостановлен",
        blocked: "Заблокирован",
        inactive: "Неактивен",
      };
      return texts[status] || "Неизвестно";
    },

    viewReader(reader) {
      // TODO: Implement reader details view
      console.log("View reader:", reader);
    },
  };
}
