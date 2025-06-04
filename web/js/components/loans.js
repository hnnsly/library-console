function loansComponent() {
  return {
    loans: [],
    statusFilter: "",
    isLoading: false,
    showCreateModal: false,

    async init() {
      this.isLoading = true;
      await this.loadLoans();
      this.isLoading = false;
    },

    async loadLoans() {
      try {
        // TODO: Implement proper loans endpoint
        // For now, use overdue books as example
        this.loans = await api.getOverdueBooks({ limit: 50 });
      } catch (error) {
        console.error("Failed to load loans:", error);
        this.loans = [];
      }
    },

    async loadOverdueBooks() {
      this.isLoading = true;
      try {
        this.loans = await api.getOverdueBooks({ limit: 50 });
      } catch (error) {
        console.error("Failed to load overdue books:", error);
        this.loans = [];
      } finally {
        this.isLoading = false;
      }
    },

    async loadDueToday() {
      this.isLoading = true;
      try {
        this.loans = await api.getDueToday();
      } catch (error) {
        console.error("Failed to load books due today:", error);
        this.loans = [];
      } finally {
        this.isLoading = false;
      }
    },

    async returnBook(loan) {
      try {
        await api.returnBook(loan.id, authManager.user?.id);
        await this.loadLoans();
      } catch (error) {
        console.error("Failed to return book:", error);
        alert("Ошибка при возврате книги: " + error.message);
      }
    },

    async renewLoan(loan) {
      try {
        await api.renewLoan(loan.id);
        await this.loadLoans();
      } catch (error) {
        console.error("Failed to renew loan:", error);
        alert("Ошибка при продлении: " + error.message);
      }
    },

    getLoanStatusColor(loan) {
      if (loan.status === "returned") {
        return "bg-green-100 text-green-800";
      }
      if (loan.status === "lost") {
        return "bg-red-100 text-red-800";
      }
      if (this.isOverdue(loan)) {
        return "bg-red-100 text-red-800";
      }
      return "bg-blue-100 text-blue-800";
    },

    getLoanStatusText(loan) {
      if (loan.status === "returned") return "Возвращена";
      if (loan.status === "lost") return "Потеряна";
      if (this.isOverdue(loan)) return "Просрочена";
      return "Активна";
    },

    isOverdue(loan) {
      if (loan.status === "returned") return false;
      return new Date(loan.due_date) < new Date();
    },

    formatDate(dateString) {
      if (!dateString) return "";
      return new Date(dateString).toLocaleDateString("ru-RU");
    },
  };
}
