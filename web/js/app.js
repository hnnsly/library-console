function app() {
  return {
    // Authentication state
    isAuthenticated: false,
    user: null,
    loading: false,
    successMessage: "",
    error: "",

    loginForm: {
      username: "",
      password: "",
    },

    // App state
    currentView: "dashboard",
    showModal: null,

    // Dashboard data
    dashboardStats: {
      totalBooks: 0,
      totalReaders: 0,
      overdueBooks: 0,
      debtors: 0,
    },
    booksDueToday: [],
    overdueBooks: [],

    // Books data
    books: [],
    bookSearchResults: [],
    searchFilters: {
      title: "",
      author: "",
      isbn: "",
      category_id: "",
      hall_id: "",
    },

    bookForm: {
      title: "",
      author: "",
      publication_year: new Date().getFullYear(),
      isbn: "",
      book_code: "",
      category_id: null,
      hall_id: 1,
      total_copies: 1,
      condition_status: "good",
      location_info: "",
    },

    // Readers data
    readers: [],
    readerSearchResults: [],
    readerSearchTerm: "",
    readerTicketSearch: "",

    readerForm: {
      full_name: "",
      ticket_number: "",
      birth_date: "",
      phone: "",
      email: "",
      education: "",
      hall_id: 1,
    },

    // Loans data
    loans: [],
    loanFilter: "all",
    loanForm: {
      bookSearch: "",
      readerSearch: "",
      selectedBook: null,
      selectedReader: null,
    },

    // Statistics data
    stats: {},
    monthlyStats: {},

    categories: [],
    halls: [],

    async init() {
      console.log("App initializing...");

      // Store global app instance
      window.appInstance = this;

      // Initialize auth
      await this.initAuth();

      // Load initial data if authenticated
      if (this.isAuthenticated) {
        await this.loadInitialData();
      }

      // Auto-hide messages
      this.$watch("successMessage", (value) => {
        if (value) {
          setTimeout(() => {
            this.successMessage = "";
          }, 5000);
        }
      });

      this.$watch("error", (value) => {
        if (value) {
          setTimeout(() => {
            this.error = "";
          }, 8000);
        }
      });

      // Watch current view changes
      this.$watch("currentView", async (newView) => {
        await this.loadViewData(newView);
      });

      // Watch loan filter changes
      this.$watch("loanFilter", async () => {
        if (this.currentView === "loans") {
          await this.loadLoans();
        }
      });

      console.log("App initialized successfully");
    },

    // Auth methods
    async initAuth() {
      try {
        this.loading = true;
        console.log("Checking authentication status...");

        // Сначала проверим debug endpoint
        try {
          const debugInfo = await api.debugAuth();
          console.log("Debug info:", debugInfo);
        } catch (debugError) {
          console.log("Debug endpoint not available:", debugError.message);
        }

        // Теперь пытаемся получить пользователя
        this.user = await api.getCurrentUser();
        this.isAuthenticated = true;
        console.log("User authenticated:", this.user);
      } catch (error) {
        console.log("Not authenticated:", error.message);
        this.isAuthenticated = false;
        this.user = null;
      } finally {
        this.loading = false;
      }
    },

    async login() {
      if (!this.loginForm.username || !this.loginForm.password) {
        this.error = "Пожалуйста, заполните все поля";
        return;
      }

      try {
        this.loading = true;
        this.error = "";

        console.log("Attempting login...");
        const response = await api.login({
          username: this.loginForm.username,
          password: this.loginForm.password,
        });

        console.log("Login response:", response);

        this.user = response.user;
        this.isAuthenticated = true;

        // Reset form
        this.loginForm.username = "";
        this.loginForm.password = "";

        // Load initial data
        await this.loadInitialData();

        console.log("Login successful, user:", this.user);
      } catch (error) {
        this.error = error.message || "Ошибка входа в систему";
        console.error("Login error:", error);
      } finally {
        this.loading = false;
      }
    },

    async logout() {
      try {
        console.log("Logging out...");
        await api.logout();
      } catch (error) {
        console.error("Logout error:", error);
      } finally {
        this.isAuthenticated = false;
        this.user = null;
        this.error = "";
        this.resetData();
        console.log("Logged out successfully");
      }
    },

    // Data loading methods
    async loadInitialData() {
      try {
        this.loading = true;
        console.log("Loading initial data...");
        await Promise.all([
          this.loadDashboardData(),
          this.loadCategories(),
          this.loadHalls(),
        ]);
        console.log("Initial data loaded successfully");
      } catch (error) {
        console.error("Error loading initial data:", error);
        this.showError("Ошибка загрузки данных: " + error.message);
      } finally {
        this.loading = false;
      }
    },

    async loadViewData(view) {
      this.loading = true;
      try {
        switch (view) {
          case "dashboard":
            await this.loadDashboardData();
            break;
          case "books":
            await this.loadBooks();
            break;
          case "readers":
            await this.loadReaders();
            break;
          case "loans":
            await this.loadLoans();
            break;
          case "statistics":
            await this.loadStatistics();
            break;
        }
      } catch (error) {
        this.showError(`Ошибка загрузки данных: ${error.message}`);
      } finally {
        this.loading = false;
      }
    },

    async loadDashboardData() {
      try {
        const [overview, dueToday, overdue] = await Promise.all([
          api.getLibraryOverview(),
          api.getBooksDueToday(),
          api.getOverdueBooks(20),
        ]);

        this.dashboardStats = {
          totalBooks:
            overview.hall_statistics?.reduce(
              (sum, hall) => sum + (hall.total_books || 0),
              0,
            ) || 0,
          totalReaders: overview.total_readers || 0,
          overdueBooks: overview.overdue_books_count || 0,
          debtors: overview.unpaid_fines_count || 0,
        };

        this.booksDueToday = dueToday || [];
        this.overdueBooks = overdue || [];
      } catch (error) {
        console.error("Error loading dashboard data:", error);
      }
    },

    async loadBooks() {
      try {
        this.loading = true;
        this.books = await api.getBooks();
      } catch (error) {
        this.showError("Ошибка загрузки книг: " + error.message);
        this.books = [];
      } finally {
        this.loading = false;
      }
    },

    async loadReaders() {
      try {
        this.loading = true;
        const params = {
          page_offset: 0,
          page_limit: 50,
        };
        this.readers = await api.getReaders(params);
      } catch (error) {
        this.showError("Ошибка загрузки читателей: " + error.message);
        this.readers = [];
      } finally {
        this.loading = false;
      }
    },

    async loadLoans() {
      try {
        switch (this.loanFilter) {
          case "overdue":
            this.loans = await api.getOverdueBooks();
            break;
          case "active":
            // Тут нужно будет добавить метод для получения активных займов
            this.loans = [];
            break;
          case "returned":
            // Тут нужно будет добавить метод для получения возвращенных книг
            this.loans = [];
            break;
          default:
            this.loans = [];
        }
      } catch (error) {
        console.error("Error loading loans:", error);
        this.loans = [];
      }
    },

    async loadStatistics() {
      try {
        const [overview, monthly] = await Promise.all([
          api.getLibraryOverview(),
          api.getMonthlyReport(),
        ]);

        this.stats = {
          totalReaders: overview.total_readers || 0,
          activeReaders: 0,
          debtorReaders: overview.unpaid_fines_count || 0,
          totalBooks:
            overview.hall_statistics?.reduce(
              (sum, hall) => sum + (hall.total_books || 0),
              0,
            ) || 0,
          availableBooks:
            overview.hall_statistics?.reduce(
              (sum, hall) => sum + (hall.available_books || 0),
              0,
            ) || 0,
          loanedBooks: 0,
        };

        this.monthlyStats = {
          loans: monthly?.[0]?.total_loans || 0,
          returns: monthly?.[0]?.total_returns || 0,
          newReaders: monthly?.[0]?.new_readers || 0,
        };
      } catch (error) {
        console.error("Error loading statistics:", error);
      }
    },

    async loadCategories() {
      try {
        this.categories = await api.getAllCategories();
      } catch (error) {
        console.error("Error loading categories:", error);
        this.categories = [];
      }
    },

    async loadHalls() {
      try {
        this.halls = await api.getAllHalls();
      } catch (error) {
        console.error("Error loading halls:", error);
        this.halls = [];
      }
    },

    // Book methods
    async searchBooks() {
      if (
        !this.searchFilters.title &&
        !this.searchFilters.author &&
        !this.searchFilters.isbn
      ) {
        await this.loadBooks();
        return;
      }

      try {
        this.loading = true;
        const searchData = {
          title: this.searchFilters.title,
          author: this.searchFilters.author,
          isbn: this.searchFilters.isbn,
          category_id: this.searchFilters.category_id
            ? parseInt(this.searchFilters.category_id)
            : 0,
          hall_id: this.searchFilters.hall_id
            ? parseInt(this.searchFilters.hall_id)
            : 0,
          page_offset: 0,
          page_limit: 50,
        };

        this.books = await api.searchBooks(searchData);
      } catch (error) {
        this.showError("Ошибка поиска книг: " + error.message);
      } finally {
        this.loading = false;
      }
    },

    async searchBooksForLoan() {
      if (!this.loanForm.bookSearch || this.loanForm.bookSearch.length < 2) {
        this.bookSearchResults = [];
        return;
      }

      try {
        const searchData = {
          title: this.loanForm.bookSearch,
          book_code: this.loanForm.bookSearch,
          page_offset: 0,
          page_limit: 10,
        };

        this.bookSearchResults = await api.searchBooks(searchData);
        // Filter only available books
        this.bookSearchResults = this.bookSearchResults.filter(
          (book) => book.available_copies > 0,
        );
      } catch (error) {
        console.error("Error searching books for loan:", error);
        this.bookSearchResults = [];
      }
    },

    selectBookForLoan(book) {
      this.loanForm.selectedBook = book;
      this.loanForm.bookSearch = `${book.title} (${book.book_code})`;
      this.bookSearchResults = [];
    },

    async addBook() {
      try {
        this.loading = true;

        if (
          !this.bookForm.title ||
          !this.bookForm.author ||
          !this.bookForm.book_code
        ) {
          throw new Error("Заполните все обязательные поля");
        }

        const bookData = {
          title: this.bookForm.title,
          author: this.bookForm.author,
          publication_year: parseInt(this.bookForm.publication_year),
          isbn: this.bookForm.isbn || null,
          book_code: this.bookForm.book_code,
          category_id: this.bookForm.category_id
            ? parseInt(this.bookForm.category_id)
            : null,
          hall_id: parseInt(this.bookForm.hall_id),
          total_copies: parseInt(this.bookForm.total_copies),
          condition_status: this.bookForm.condition_status,
          location_info: this.bookForm.location_info || null,
        };

        await api.createBook(bookData);

        this.showSuccess("Книга успешно добавлена");
        this.showModal = null;
        this.resetBookForm();
        await this.loadBooks();
      } catch (error) {
        this.showError("Ошибка добавления книги: " + error.message);
      } finally {
        this.loading = false;
      }
    },

    viewBook(book) {
      console.log("Viewing book:", book);
    },

    editBook(book) {
      console.log("Editing book:", book);
    },

    resetBookForm() {
      this.bookForm = {
        title: "",
        author: "",
        publication_year: new Date().getFullYear(),
        isbn: "",
        book_code: "",
        category_id: null,
        hall_id: 1,
        total_copies: 1,
        condition_status: "good",
        location_info: "",
      };
    },

    // Reader methods
    async searchReaders() {
      if (!this.readerSearchTerm || this.readerSearchTerm.length < 2) {
        await this.loadReaders();
        return;
      }

      try {
        this.loading = true;
        const searchData = {
          search_name: this.readerSearchTerm,
          page_offset: 0,
          page_limit: 50,
        };

        this.readers = await api.searchReadersByName(searchData);
      } catch (error) {
        this.showError("Ошибка поиска читателей: " + error.message);
      } finally {
        this.loading = false;
      }
    },

    async searchReaderByTicket() {
      if (!this.readerTicketSearch) {
        return;
      }

      try {
        this.loading = true;
        const reader = await api.getReaderByTicket(this.readerTicketSearch);
        this.readers = [reader];
      } catch (error) {
        if (error.message.includes("not found")) {
          this.readers = [];
        } else {
          this.showError("Ошибка поиска читателя: " + error.message);
        }
      } finally {
        this.loading = false;
      }
    },

    async searchReadersForLoan() {
      if (
        !this.loanForm.readerSearch ||
        this.loanForm.readerSearch.length < 2
      ) {
        this.readerSearchResults = [];
        return;
      }

      try {
        const searchData = {
          search_name: this.loanForm.readerSearch,
          page_offset: 0,
          page_limit: 10,
        };

        this.readerSearchResults = await api.searchReadersByName(searchData);

        // Also try to search by ticket number
        if (this.loanForm.readerSearch.match(/^\d+$/)) {
          try {
            const readerByTicket = await api.getReaderByTicket(
              this.loanForm.readerSearch,
            );
            // Add to results if not already present
            const exists = this.readerSearchResults.find(
              (r) => r.id === readerByTicket.id,
            );
            if (!exists) {
              this.readerSearchResults.unshift(readerByTicket);
            }
          } catch (ticketError) {
            // Ignore if not found by ticket
          }
        }
      } catch (error) {
        console.error("Error searching readers for loan:", error);
        this.readerSearchResults = [];
      }
    },

    selectReaderForLoan(reader) {
      this.loanForm.selectedReader = reader;
      this.loanForm.readerSearch = `${reader.full_name} (${reader.ticket_number})`;
      this.readerSearchResults = [];
    },

    async addReader() {
      try {
        this.loading = true;

        if (
          !this.readerForm.full_name ||
          !this.readerForm.ticket_number ||
          !this.readerForm.birth_date
        ) {
          throw new Error("Заполните все обязательные поля");
        }

        const readerData = {
          full_name: this.readerForm.full_name,
          ticket_number: this.readerForm.ticket_number,
          birth_date: this.readerForm.birth_date + "T00:00:00Z",
          phone: this.readerForm.phone || null,
          email: this.readerForm.email || null,
          education: this.readerForm.education || null,
          hall_id: parseInt(this.readerForm.hall_id),
        };

        await api.createReader(readerData);

        this.showSuccess("Читатель успешно добавлен");
        this.showModal = null;
        this.resetReaderForm();
        await this.loadReaders();
      } catch (error) {
        this.showError("Ошибка добавления читателя: " + error.message);
      } finally {
        this.loading = false;
      }
    },

    viewReader(reader) {
      console.log("Viewing reader:", reader);
    },

    editReader(reader) {
      console.log("Editing reader:", reader);
    },

    resetReaderForm() {
      this.readerForm = {
        full_name: "",
        ticket_number: "",
        birth_date: "",
        phone: "",
        email: "",
        education: "",
        hall_id: 1,
      };
    },

    // Loan methods
    async createLoan() {
      if (!this.loanForm.selectedBook || !this.loanForm.selectedReader) {
        this.showError("Выберите книгу и читателя");
        return;
      }

      try {
        this.loading = true;

        const loanData = {
          book_id: this.loanForm.selectedBook.id,
          reader_id: this.loanForm.selectedReader.id,
          librarian_id: this.user.id, // Assuming current user is librarian
        };

        await api.createLoan(loanData);

        this.showSuccess("Книга успешно выдана");
        this.showModal = null;
        this.resetLoanForm();

        // Reload data
        if (this.currentView === "loans") {
          await this.loadLoans();
        }
        if (this.currentView === "dashboard") {
          await this.loadDashboardData();
        }
      } catch (error) {
        this.showError("Ошибка выдачи книги: " + error.message);
      } finally {
        this.loading = false;
      }
    },

    resetLoanForm() {
      this.loanForm = {
        bookSearch: "",
        readerSearch: "",
        selectedBook: null,
        selectedReader: null,
      };
      this.bookSearchResults = [];
      this.readerSearchResults = [];
    },

    // Computed properties
    get filteredLoans() {
      if (!this.loans) return [];

      switch (this.loanFilter) {
        case "active":
          return this.loans.filter((loan) => loan.status === "active");
        case "overdue":
          return this.loans.filter(
            (loan) => loan.status === "overdue" || loan.days_overdue > 0,
          );
        case "returned":
          return this.loans.filter((loan) => loan.status === "returned");
        default:
          return this.loans;
      }
    },

    // Utility methods
    showSuccess(message) {
      this.successMessage = message;
      this.error = "";
    },

    showError(message) {
      this.error = message;
      this.successMessage = "";
    },

    formatDate(dateString) {
      if (!dateString) return "-";
      try {
        return new Date(dateString).toLocaleDateString("ru-RU");
      } catch {
        return dateString;
      }
    },

    getStatusText(status) {
      const statusMap = {
        active: "Активная",
        overdue: "Просрочена",
        returned: "Возвращена",
        lost: "Утеряна",
      };
      return statusMap[status] || status;
    },

    resetData() {
      this.currentView = "dashboard";
      this.books = [];
      this.readers = [];
      this.loans = [];
      this.dashboardStats = {
        totalBooks: 0,
        totalReaders: 0,
        overdueBooks: 0,
        debtors: 0,
      };
      this.booksDueToday = [];
      this.overdueBooks = [];
      this.stats = {};
      this.monthlyStats = {};
      this.showModal = null;
      this.resetBookForm();
      this.resetReaderForm();
      this.resetLoanForm();
    },
  };
}

// Register Alpine.js component
document.addEventListener("alpine:init", () => {
  Alpine.data("app", app);
});
