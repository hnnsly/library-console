function initAuth() {
  return {
    isAuthenticated: false,
    user: null,
    loading: false,
    error: "",

    loginForm: {
      username: "",
      password: "",
    },

    async init() {
      const sessionId = localStorage.getItem("session_id");
      if (sessionId) {
        try {
          this.loading = true;
          this.user = await api.getCurrentUser();
          this.isAuthenticated = true;
        } catch (error) {
          console.log("Session expired or invalid");
          localStorage.removeItem("session_id");
          this.isAuthenticated = false;
        } finally {
          this.loading = false;
        }
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

        const response = await api.login({
          username: this.loginForm.username,
          password: this.loginForm.password,
        });

        if (response.session_id) {
          localStorage.setItem("session_id", response.session_id);
          this.user = response.user;
          this.isAuthenticated = true;

          // Reset form
          this.loginForm.username = "";
          this.loginForm.password = "";

          // Initialize app data
          await this.$nextTick();
          if (window.appInstance && window.appInstance.loadInitialData) {
            await window.appInstance.loadInitialData();
          }
        }
      } catch (error) {
        this.error = error.message || "Ошибка входа в систему";
        console.error("Login error:", error);
      } finally {
        this.loading = false;
      }
    },

    async logout() {
      try {
        await api.logout();
      } catch (error) {
        console.error("Logout error:", error);
      } finally {
        localStorage.removeItem("session_id");
        this.isAuthenticated = false;
        this.user = null;
        this.error = "";

        // Reset all forms and data
        if (window.appInstance && window.appInstance.resetData) {
          window.appInstance.resetData();
        }
      }
    },

    hasRole(role) {
      return this.user && this.user.role === role;
    },

    hasAnyRole(roles) {
      return this.user && roles.includes(this.user.role);
    },
  };
}
