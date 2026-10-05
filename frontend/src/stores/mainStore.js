import { defineStore } from "pinia";
import axios from "axios";
import router from "../router";
import { toast } from "vue-sonner";

const baseUrl = "http://localhost:8080/api/v1";

export const useMainStore = defineStore("main", {
  state: () => ({
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

          router.push("/");
          return true;
        }
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
          router.push("/login");
        }
      } catch (err) {
        toast.error("Registrasi Gagal", {
          description: err.response?.data?.message || "Periksa data legalitas & input Anda",
        });
      } finally {
        this.loading = false;
      }
    },

    handleLogout() {
      localStorage.clear();
      this.isLogin = false;
      this.user = { email: "", fullName: "", role: "", companyName: "" };
      toast.info("Anda telah keluar dari sistem.");
      router.push("/login");
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
