function booksComponent() {
  return {
    books: [],
    categories: [],
    halls: [],
    searchQuery: "",
    selectedCategory: "",
    selectedHall: "",
    isLoading: false,
    showCreateModal: false,

    async init() {
      this.isLoading = true;
      await Promise.all([
        this.loadBooks(),
        this.loadCategories(),
        this.loadHalls(),
      ]);
      this.isLoading = false;
    },

    async loadBooks() {
      try {
        this.books = await api.getBooks({ limit: 50 });
      } catch (error) {
        console.error("Failed to load books:", error);
        this.books = [];
      }
    },

    async loadCategories() {
      try {
        this.categories = await api.getCategories();
      } catch (error) {
        console.error("Failed to load categories:", error);
        this.categories = [];
      }
    },

    async loadHalls() {
      try {
        this.halls = await api.getHalls();
      } catch (error) {
        console.error("Failed to load halls:", error);
        this.halls = [];
      }
    },

    async searchBooks() {
      if (
        !this.searchQuery.trim() &&
        !this.selectedCategory &&
        !this.selectedHall
      ) {
        return this.loadBooks();
      }

      this.isLoading = true;
      try {
        const searchData = {};

        if (this.searchQuery.trim()) {
          searchData.title = this.searchQuery;
          searchData.author = this.searchQuery;
        }
        if (this.selectedCategory) {
          searchData.category_id = parseInt(this.selectedCategory);
        }
        if (this.selectedHall) {
          searchData.hall_id = parseInt(this.selectedHall);
        }

        this.books = await api.searchBooks(searchData);
      } catch (error) {
        console.error("Failed to search books:", error);
        this.books = [];
      } finally {
        this.isLoading = false;
      }
    },

    editBook(book) {
      // TODO: Implement edit functionality
      console.log("Edit book:", book);
    },
  };
}
