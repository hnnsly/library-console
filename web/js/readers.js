function initReaders() {
  return {
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

    selectedReader: null,

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
      } finally {
        this.loading = false;
      }
    },

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

        // Validation
        if (
          !this.readerForm.full_name ||
          !this.readerForm.ticket_number ||
          !this.readerForm.birth_date
        ) {
          throw new Error("Заполните все обязательные поля");
        }

        // Validate birth date
        const birthDate = new Date(this.readerForm.birth_date);
        const today = new Date();
        const age = Math.floor(
          (today - birthDate) / (365.25 * 24 * 60 * 60 * 1000),
        );

        if (age < 6) {
          throw new Error("Возраст читателя должен быть не менее 6 лет");
        }

        if (birthDate > today) {
          throw new Error("Дата рождения не может быть в будущем");
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

    async viewReader(reader) {
      try {
        this.selectedReader = await api.getReaderById(reader.id);
        this.showModal = "viewReader";
      } catch (error) {
        this.showError(
          "Ошибка загрузки информации о читателе: " + error.message,
        );
      }
    },

    editReader(reader) {
      this.readerForm = {
        ...reader,
        id: reader.id,
        birth_date: reader.birth_date ? reader.birth_date.split("T")[0] : "",
      };
      this.showModal = "editReader";
    },

    async updateReader() {
      try {
        this.loading = true;

        const readerData = {
          full_name: this.readerForm.full_name,
          phone: this.readerForm.phone || null,
          email: this.readerForm.email || null,
          education: this.readerForm.education || null,
          hall_id: parseInt(this.readerForm.hall_id),
        };

        await api.updateReader(this.readerForm.id, readerData);

        this.showSuccess("Информация о читателе обновлена");
        this.showModal = null;
        this.resetReaderForm();
        await this.loadReaders();
      } catch (error) {
        this.showError("Ошибка обновления читателя: " + error.message);
      } finally {
        this.loading = false;
      }
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
  };
}
