import { defineStore } from "pinia";
import axios from "axios";
import { toast } from "vue-sonner";

const baseUrl = "http://localhost:8080/api/v1";

export const useMainStore = defineStore("main", {
  state: () => ({
    theme: localStorage.getItem("scrapflow_theme") || "dark",
    isLogin: localStorage.getItem("access_token") ? true : false,
    user: {
      email: localStorage.getItem("user_email") || "",
      fullName: localStorage.getItem("user_name") || "",
      role: localStorage.getItem("user_role") || "",
      companyName: localStorage.getItem("company_name") || "",
    },
    walletBalance: 0,
    tenders: [],
    loading: false,
  }),
  actions: {
    initTheme() {
      const saved = localStorage.getItem("scrapflow_theme") || "dark";
      this.theme = saved;
      this.applyTheme(saved);
    },

    toggleTheme() {
      const newTheme = this.theme === "dark" ? "light" : "dark";
      this.theme = newTheme;
      localStorage.setItem("scrapflow_theme", newTheme);
      this.applyTheme(newTheme);
      toast.success(newTheme === "dark" ? "Mode Gelap Aktif" : "Mode Terang Aktif", {
        description: `Tampilan dialihkan ke ${newTheme === "dark" ? "Dark Mode" : "Light Mode"}.`,
      });
    },

    applyTheme(theme) {
      this.theme = theme;
      localStorage.setItem("scrapflow_theme", theme);
      if (theme === "light") {
        document.documentElement.classList.add("light");
        document.documentElement.classList.remove("dark");
      } else {
        document.documentElement.classList.add("dark");
        document.documentElement.classList.remove("light");
      }
    },

    async handleLogin(email, password) {
      this.loading = true;
      try {
        const { data } = await axios.post(`${baseUrl}/auth/login`, {
          email,
          password,
        });

        if (data.success) {
          const userPayload = data.data;
          localStorage.setItem("access_token", userPayload.token);
          localStorage.setItem("user_email", userPayload.email);
          localStorage.setItem("user_name", userPayload.full_name);
          localStorage.setItem("user_role", userPayload.role);
          localStorage.setItem("company_name", userPayload.legal_name);

          this.isLogin = true;
          this.user = {
            email: userPayload.email,
            fullName: userPayload.full_name,
            role: userPayload.role,
            companyName: userPayload.legal_name,
          };

          toast.success(`Selamat datang, ${userPayload.full_name}!`, {
            description: `Login sebagai ${userPayload.legal_name} (${userPayload.role})`,
          });

          return true;
        }
        return false;
      } catch (err) {
        const errorMsg = err.response?.data?.message || "Email atau password tidak sesuai";
        toast.error("Gagal Masuk", {
          description: errorMsg,
        });
        return false;
      } finally {
        this.loading = false;
      }
    },

    async handleRegister(payload) {
      this.loading = true;
      try {
        const { data } = await axios.post(`${baseUrl}/auth/register`, payload);
        if (data.success) {
          toast.success("Registrasi Perusahaan Berhasil!", {
            description: "Silakan login dengan akun yang baru didaftarkan.",
          });
          return true;
        }
        return false;
      } catch (err) {
        toast.error("Registrasi Gagal", {
          description: err.response?.data?.message || "Periksa data legalitas & input Anda",
        });
        return false;
      } finally {
        this.loading = false;
      }
    },

    handleLogout() {
      localStorage.removeItem("access_token");
      localStorage.removeItem("user_email");
      localStorage.removeItem("user_name");
      localStorage.removeItem("user_role");
      localStorage.removeItem("company_name");

      this.isLogin = false;
      this.user = { email: "", fullName: "", role: "", companyName: "" };

      toast.info("Sesi Berakhir", {
        description: "Anda telah berhasil keluar dari sistem ScrapFlow.",
      });
      return true;
    },

    async fetchTenders() {
      try {
        const { data } = await axios.get(`${baseUrl}/tenders`);
        if (data.success) {
          this.tenders = data.data;
        }
      } catch (err) {
        console.error("Failed to fetch tenders:", err);
      }
    },
  },
});
