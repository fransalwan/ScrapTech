<script>
import { mapActions, mapState } from "pinia";
import { useMainStore } from "../stores/mainStore";
import { RouterLink } from "vue-router";
import { LogOut, User, Building2, ShieldCheck, Sun, Moon } from "lucide-vue-next";

export default {
  name: "Navbar",
  components: { RouterLink, LogOut, User, Building2, ShieldCheck, Sun, Moon },
  data() {
    return {
      highContrast: false,
    };
  },
  methods: {
    ...mapActions(useMainStore, ["handleLogout"]),
    toggleContrast() {
      this.highContrast = !this.highContrast;
      document.documentElement.classList.toggle("high-contrast-mode", this.highContrast);
    },
  },
  computed: {
    ...mapState(useMainStore, ["isLogin", "user"]),
  },
};
</script>

<template>
  <nav class="bg-slate-900 border-b border-slate-800 text-slate-100 sticky top-0 z-40 backdrop-blur-md bg-slate-900/95">
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
      <div class="flex items-center justify-between h-16">
        
        <!-- Logo & Branding -->
        <div class="flex items-center space-x-3">
          <RouterLink to="/" class="flex items-center space-x-2.5">
            <div class="w-8 h-8 sm:w-9 sm:h-9 rounded-xl bg-emerald-500 flex items-center justify-center font-bold text-slate-950 text-lg sm:text-xl shadow-lg shadow-emerald-500/20">
              ⚡
            </div>
            <div>
              <span class="text-lg sm:text-xl font-bold tracking-tight text-white">Scrap<span class="text-emerald-400">Flow</span></span>
              <span class="hidden sm:inline-block ml-2 px-2 py-0.5 text-[10px] font-semibold bg-emerald-950 text-emerald-400 border border-emerald-800 rounded-full uppercase tracking-wider">
                Fintech Core
              </span>
            </div>
          </RouterLink>

          <!-- Core Engine Badge -->
          <div class="hidden lg:flex items-center space-x-2 ml-4 pl-4 border-l border-slate-800 text-xs text-slate-400">
            <span class="w-2 h-2 rounded-full bg-emerald-400 animate-pulse"></span>
            <span>Golang Ledger: <strong class="text-slate-200">ACTIVE</strong></span>
          </div>
        </div>

        <!-- Desktop Navigation Links -->
        <div class="hidden md:flex items-center space-x-1 text-xs sm:text-sm font-medium">
          <a href="#tenders" class="text-slate-300 hover:text-white px-3 py-1.5 rounded-lg hover:bg-slate-800 transition">
            Tender Board
          </a>
          <a href="#weighbridge" class="text-slate-300 hover:text-white px-3 py-1.5 rounded-lg hover:bg-slate-800 transition">
            Jembatan Timbang
          </a>
          <a href="#financing" class="text-slate-300 hover:text-white px-3 py-1.5 rounded-lg hover:bg-slate-800 transition">
            Talangan SCF
          </a>
          <a href="#ledger" class="text-slate-300 hover:text-white px-3 py-1.5 rounded-lg hover:bg-slate-800 transition">
            Audit Ledger
          </a>
        </div>

        <!-- User Controls / Auth -->
        <div class="flex items-center space-x-2.5">
          
          <!-- High Contrast / Field Mode Toggle -->
          <button
            @click="toggleContrast"
            title="Mode Kontras Lapangan (Sinar Matahari)"
            class="p-2 text-slate-400 hover:text-amber-400 hover:bg-slate-800 rounded-lg transition"
          >
            <Sun v-if="!highContrast" class="w-4 h-4" />
            <Moon v-else class="w-4 h-4 text-amber-400" />
          </button>

          <!-- If Logged In -->
          <div v-if="isLogin" class="flex items-center space-x-2 sm:space-x-3">
            <div class="text-right hidden sm:block">
              <p class="text-xs font-bold text-white leading-tight">{{ user.fullName || "Pengguna Aktif" }}</p>
              <p class="text-[10px] text-emerald-400 font-medium">{{ user.companyName || "Perusahaan Rekanan" }}</p>
            </div>
            <span class="px-2 py-0.5 rounded text-[10px] font-bold bg-emerald-950 text-emerald-400 border border-emerald-800 uppercase">
              {{ user.role || "MEMBER" }}
            </span>
            <button
              @click="handleLogout"
              title="Logout"
              class="p-2 text-slate-400 hover:text-rose-400 hover:bg-slate-800 rounded-lg transition"
            >
              <LogOut class="w-4 h-4" />
            </button>
          </div>

          <!-- If Not Logged In -->
          <div v-else class="flex items-center space-x-2">
            <RouterLink
              to="/login"
              class="px-3 py-1.5 text-xs font-semibold text-slate-200 hover:text-white hover:bg-slate-800 border border-slate-700 rounded-lg transition"
            >
              Sign In
            </RouterLink>
            <RouterLink
              to="/register"
              class="px-3 py-1.5 text-xs font-semibold text-slate-950 bg-emerald-400 hover:bg-emerald-300 rounded-lg shadow-sm transition"
            >
              Register
            </RouterLink>
          </div>

        </div>

      </div>
    </div>
  </nav>
</template>
