class AuthManager {
  constructor() {
    this.user = null;
    this.isAuthenticated = false;
  }

  async checkAuth() {
    try {
      const response = await api.getMe();
      this.user = response;
      this.isAuthenticated = true;
      return true;
    } catch (error) {
      this.user = null;
      this.isAuthenticated = false;
      return false;
    }
  }

  async login(credentials) {
    try {
      const response = await api.login(credentials);
      this.user = response.user;
      this.isAuthenticated = true;
      return response;
    } catch (error) {
      this.user = null;
      this.isAuthenticated = false;
      throw error;
    }
  }

  async logout() {
    try {
      await api.logout();
    } catch (error) {
      console.error("Logout error:", error);
    } finally {
      this.user = null;
      this.isAuthenticated = false;
    }
  }

  hasRole(role) {
    if (!this.user) return false;

    const roleHierarchy = {
      reader: 0,
      librarian: 1,
      admin: 2,
      super_admin: 3,
    };

    const userLevel = roleHierarchy[this.user.role] || -1;
    const requiredLevel = roleHierarchy[role] || 999;

    return userLevel >= requiredLevel;
  }
}

window.authManager = new AuthManager();
