function loginComponent() {
  return {
    credentials: {
      username: "",
      password: "",
    },
    error: "",
    isSubmitting: false,

    async login() {
      if (this.isSubmitting) return;

      this.isSubmitting = true;
      this.error = "";

      try {
        await authManager.login(this.credentials);
        // Trigger app state update
        this.$dispatch("auth-changed");
      } catch (error) {
        this.error = error.message || "Ошибка входа в систему";
      } finally {
        this.isSubmitting = false;
      }
    },
  };
}
