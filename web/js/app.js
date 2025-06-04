function app() {
  return {
    isAuthenticated: false,
    isLoading: true,
    user: null,
    activeView: "dashboard",

    async init() {
      // Check authentication status
      const isAuth = await authManager.checkAuth();
      this.isAuthenticated = isAuth;
      this.user = authManager.user;
      this.isLoading = false;

      // Listen for auth changes
      this.$watch("isAuthenticated", (value) => {
        if (value) {
          this.user = authManager.user;
          if (this.activeView === "login") {
            this.activeView = "dashboard";
          }
        }
      });

      // Listen for auth-changed events
      document.addEventListener("auth-changed", () => {
        this.isAuthenticated = authManager.isAuthenticated;
        this.user = authManager.user;
      });
    },

    setActiveView(view) {
      this.activeView = view;
    },

    async logout() {
      this.isLoading = true;
      await authManager.logout();
      this.isAuthenticated = false;
      this.user = null;
      this.activeView = "dashboard";
      this.isLoading = false;
    },
  };
}
