//const API_BASE_URL = "http://localhost:8080";
const API_BASE_URL = "https://lab.somerka.ru";

class ApiClient {
  constructor() {
    this.baseURL = API_BASE_URL;
  }

  async request(endpoint, options = {}) {
    const url = `${this.baseURL}${endpoint}`;
    const config = {
      credentials: "include",
      headers: {
        "Content-Type": "application/json",
        ...options.headers,
      },
      ...options,
    };

    if (options.body && typeof options.body === "object") {
      config.body = JSON.stringify(options.body);
    }

    try {
      const response = await fetch(url, config);

      if (!response.ok) {
        const error = await response
          .json()
          .catch(() => ({ error: "Network error" }));
        throw new Error(error.error || "Request failed");
      }

      return await response.json();
    } catch (error) {
      console.error("API request failed:", error);
      throw error;
    }
  }

  // Auth methods
  async login(credentials) {
    return this.request("/auth/login", {
      method: "POST",
      body: credentials,
    });
  }

  async logout() {
    return this.request("/auth/logout", {
      method: "POST",
    });
  }

  async getMe() {
    return this.request("/auth/me");
  }

  // Books methods
  async getBooks(params = {}) {
    const query = new URLSearchParams(params).toString();
    return this.request(`/api/books/${query ? "?" + query : ""}`);
  }

  async searchBooks(searchData) {
    return this.request("/api/books/search", {
      method: "POST",
      body: searchData,
    });
  }

  async createBook(bookData) {
    return this.request("/api/books/", {
      method: "POST",
      body: bookData,
    });
  }

  async getBookById(id) {
    return this.request(`/api/books/${id}`);
  }

  // Readers methods
  async getReaders(params = {}) {
    return this.request("/api/readers/", {
      method: "POST",
      body: params,
    });
  }

  async searchReaders(searchData) {
    return this.request("/api/readers/search", {
      method: "POST",
      body: searchData,
    });
  }

  async createReader(readerData) {
    return this.request("/api/readers/", {
      method: "POST",
      body: readerData,
    });
  }

  async getReaderById(id) {
    return this.request(`/api/readers/${id}`);
  }

  // Loans methods
  async createLoan(loanData) {
    return this.request("/api/loans/", {
      method: "POST",
      body: loanData,
    });
  }

  async getOverdueBooks(params = {}) {
    const query = new URLSearchParams(params).toString();
    return this.request(`/api/loans/overdue${query ? "?" + query : ""}`);
  }

  async getDueToday() {
    return this.request("/api/loans/due-today");
  }

  async returnBook(loanId, librarianId) {
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

  // Statistics methods
  async getDashboardStats() {
    return this.request("/api/statistics/dashboard");
  }

  async getRecentOperations(params = {}) {
    const query = new URLSearchParams(params).toString();
    return this.request(`/api/operations/recent${query ? "?" + query : ""}`);
  }

  // Categories and Halls
  async getCategories() {
    return this.request("/api/categories/");
  }

  async getHalls() {
    return this.request("/api/halls/");
  }
}

// Global API instance
window.api = new ApiClient();
