function initBooks() {
  return {
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

    selectedBook: null,
    categories: [],
    halls: [],

    async loadBooks() {
      try {
        this.loading = true;
        this.books = await api.getBooks();
      } catch (error) {
        this.showError("Ошибка загрузки книг: " + error.message);
      } finally {
        this.loading = false;
      }
    },

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

        // Validation
        if (
          !this.bookForm.title ||
          !this.bookForm.author ||
          !this.bookForm.book_code
        ) {
          throw new Error("Заполните все обязательные поля");
        }

        if (this.bookForm.total_copies < 1) {
          throw new Error("Количество экземпляров должно быть больше 0");
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

    async viewBook(book) {
      try {
        this.selectedBook = await api.getBookById(book.id);
        this.showModal = "viewBook";
      } catch (error) {
        this.showError("Ошибка загрузки информации о книге: " + error.message);
      }
    },

    editBook(book) {
      this.bookForm = {
        ...book,
        id: book.id,
      };
      this.showModal = "editBook";
    },

    async updateBook() {
      try {
        this.loading = true;

        const bookData = {
          book_id: this.bookForm.id,
          total_copies: parseInt(this.bookForm.total_copies),
          available_copies: parseInt(
            this.bookForm.available_copies || this.bookForm.total_copies,
          ),
        };

        await api.updateBookAvailability(this.bookForm.id, bookData);

        this.showSuccess("Информация о книге обновлена");
        this.showModal = null;
        this.resetBookForm();
        await this.loadBooks();
      } catch (error) {
        this.showError("Ошибка обновления книги: " + error.message);
      } finally {
        this.loading = false;
      }
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

    async loadCategories() {
      try {
        this.categories = await api.getAllCategories();
      } catch (error) {
        console.error("Error loading categories:", error);
      }
    },

    async loadHalls() {
      try {
        this.halls = await api.getAllHalls();
      } catch (error) {
        console.error("Error loading halls:", error);
      }
    },
  };
}
