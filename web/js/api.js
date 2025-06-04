const API_BASE_URL = "https://lab.somerka.ru";
//const API_BASE_URL = "http://localhost:8080";

class ApiClient {
  constructor() {
    // Убираем работу с localStorage для session_id
    // Куки будут обрабатываться автоматически браузером
  }

  async request(endpoint, options = {}) {
    const url = `${API_BASE_URL}${endpoint}`;
    const config = {
      headers: {
        "Content-Type": "application/json",
        ...options.headers,
      },
      credentials: "same-origin", // Это критически важно!
      ...options,
    };

    if (options.body && typeof options.body === "object") {
      config.body = JSON.stringify(options.body);
    }

    try {
      console.log(`Making ${config.method || "GET"} request to:`, url);
      console.log("Request config:", config);

      const response = await fetch(url, config);

      console.log("Response status:", response.status);
      console.log("Response headers:", response.headers);

      if (!response.ok) {
        const errorData = await response.json().catch(() => ({}));
        throw new Error(errorData.message || `HTTP ${response.status}`);
      }

      return await response.json();
    } catch (error) {
      console.error("API request failed:", error);
      throw error;
    }
  }

  // Auth endpoints
  async login(credentials) {
    const response = await this.request("/auth/login", {
      method: "POST",
      body: credentials,
    });

    // НЕ сохраняем session_id в localStorage
    // Куки устанавливаются автоматически сервером
    console.log("Login successful:", response);
    return response;
  }

  async logout() {
    const response = await this.request("/auth/logout", {
      method: "POST",
    });

    // НЕ удаляем из localStorage, куки очищаются сервером
    return response;
  }

  async getCurrentUser() {
    return this.request("/auth/me");
  }

  // Добавляем отладочный метод
  async debugAuth() {
    return this.request("/auth/debug");
  }

  // Остальные методы остаются без изменений...
  async getBooks(params = {}) {
    const searchParams = new URLSearchParams();
    Object.entries(params).forEach(([key, value]) => {
      if (value) searchParams.append(key, value);
    });

    return this.request(`/api/books/available?${searchParams}`);
  }

  async searchBooks(searchData) {
    return this.request("/api/books/search", {
      method: "POST",
      body: searchData,
    });
  }

  async createBook(bookData) {
    return this.request("/api/books", {
      method: "POST",
      body: bookData,
    });
  }

  async getBookById(id) {
    return this.request(`/api/books/${id}`);
  }

  async updateBookAvailability(id, data) {
    return this.request(`/api/books/${id}/availability`, {
      method: "PUT",
      body: data,
    });
  }

  async getPopularBooks(limit = 10) {
    return this.request(`/api/books/popular?limit=${limit}`);
  }

  async getTopRatedBooks(limit = 10) {
    return this.request(`/api/books/top-rated?limit=${limit}`);
  }

  async getReaders(params = {}) {
    return this.request("/api/readers", {
      method: "POST",
      body: params,
    });
  }

  async searchReadersByName(searchData) {
    return this.request("/api/readers/search", {
      method: "POST",
      body: searchData,
    });
  }

  async createReader(readerData) {
    return this.request("/api/readers", {
      method: "POST",
      body: readerData,
    });
  }

  async getReaderById(id) {
    return this.request(`/api/readers/${id}`);
  }

  async getReaderByTicket(ticket) {
    return this.request(`/api/readers/ticket/${ticket}`);
  }

  async updateReader(id, readerData) {
    return this.request(`/api/readers/${id}`, {
      method: "PUT",
      body: readerData,
    });
  }

  async getActiveReaders(limit = 50) {
    return this.request(`/api/readers/active?limit=${limit}`);
  }

  async getReadersCount() {
    return this.request("/api/readers/count");
  }

  async createLoan(loanData) {
    return this.request("/api/loans", {
      method: "POST",
      body: loanData,
    });
  }

  async getLoanById(id) {
    return this.request(`/api/loans/${id}`);
  }

  async getReaderCurrentLoans(readerId) {
    return this.request(`/api/readers/${readerId}/loans`);
  }

  async getOverdueBooks(limit = 50) {
    return this.request(`/api/loans/overdue?limit=${limit}`);
  }

  async getBooksDueToday() {
    return this.request("/api/loans/due-today");
  }

  async returnBook(loanId, librarianId = null) {
    return this.request(`/api/loans/${loanId}/return`, {
      method: "PUT",
      body: { librarian_id: librarianId },
    });
  }

  async renewLoan(loanId) {
    return this.request(`/api/loans/${loanId}/renew`, {
      method: "PUT",
    });
  }

  async checkLoanEligibility(readerId, bookId) {
    return this.request(
      `/api/loans/check-eligibility?reader_id=${readerId}&book_id=${bookId}`,
    );
  }

  async getDashboardStats() {
    return this.request("/api/statistics/dashboard");
  }

  async getLibraryOverview() {
    return this.request("/api/statistics/overview");
  }

  async getMonthlyReport() {
    return this.request("/api/statistics/monthly");
  }

  async getLoanStatusStatistics(daysBack = 30) {
    return this.request(`/api/statistics/loans?days_back=${daysBack}`);
  }

  async getInventoryReport() {
    return this.request("/api/statistics/inventory");
  }

  async getUnpaidFines() {
    return this.request("/api/fines/unpaid");
  }

  async getReaderFines(readerId) {
    return this.request(`/api/readers/${readerId}/fines`);
  }

  async payFine(fineId) {
    return this.request(`/api/fines/${fineId}/pay`, {
      method: "PUT",
    });
  }

  async createFine(fineData) {
    return this.request("/api/fines", {
      method: "POST",
      body: fineData,
    });
  }

  async getAllCategories() {
    return this.request("/api/categories");
  }

  async getCategoryStatistics() {
    return this.request("/api/categories/statistics");
  }

  async getAllHalls() {
    return this.request("/api/halls");
  }

  async getHallStatistics() {
    return this.request("/api/halls/statistics");
  }

  async createReservation(reservationData) {
    return this.request("/api/reservations", {
      method: "POST",
      body: reservationData,
    });
  }

  async getReaderReservations(readerId) {
    return this.request(`/api/readers/${readerId}/reservations`);
  }

  async getExpiredReservations() {
    return this.request("/api/reservations/expired");
  }

  async cancelReservation(reservationId) {
    return this.request(`/api/reservations/${reservationId}/cancel`, {
      method: "PUT",
    });
  }

  async globalSearch(searchTerm) {
    return this.request(
      `/api/search/global?q=${encodeURIComponent(searchTerm)}`,
    );
  }
}

// Create global API instance
window.api = new ApiClient();
