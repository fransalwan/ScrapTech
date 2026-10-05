<script>
import { mapActions, mapState } from "pinia";
import { useMainStore } from "../stores/mainStore";
import { RouterLink } from "vue-router";
import { LogOut, User, Building2, ShieldCheck, Sun, Moon, RefreshCw, ChevronDown } from "lucide-vue-next";

export default {
  name: "Navbar",
  components: { RouterLink, LogOut, User, Building2, ShieldCheck, Sun, Moon, RefreshCw, ChevronDown },
  data() {
    return {
      highContrast: false,
    };
  },
  methods: {
    ...mapActions(useMainStore, ["handleLogout", "handleLogin"]),
    toggleContrast() {
      this.highContrast = !this.highContrast;
      document.documentElement.classList.toggle("high-contrast-mode", this.highContrast);
    },
    quickSwitchRole(role) {
      if (role === "PABRIK_MANAGER") {
        this.handleLogin("budi@krakatausteel.co.id", "password123");
      } else if (role === "LAPAK_OWNER") {
        this.handleLogin("haji.slamet@besiabadi.com", "password123");
      } else if (role === "WEIGH_OPERATOR") {
        this.handleLogin("operator@timbangan.com", "password123");
      }
    },
    getRoleLabel(role) {
      if (role === "PABRIK_MANAGER") return "Pabrik / Seller";
      if (role === "LAPAK_OWNER") return "Juragan Lapak";
      if (role === "WEIGH_OPERATOR") return "Petugas Timbangan";
      return role || "Tamu";
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
                Fintech B2B
              </span>
            </div>
          </RouterLink>

          <!-- Core Engine Badge -->
          <div class="hidden lg:flex items-center space-x-2 ml-4 pl-4 border-l border-slate-800 text-xs text-slate-400">
            <span class="w-2 h-2 rounded-full bg-emerald-400 animate-pulse"></span>
            <span>Buku Besar Golang: <strong class="text-slate-200">AKTIF</strong></span>
          </div>
        </div>

        <!-- Desktop Navigation Links -->
        <div class="hidden md:flex items-center space-x-1 text-xs sm:text-sm font-medium">
          <a href="#tenders" class="text-slate-300 hover:text-white px-3 py-1.5 rounded-lg hover:bg-slate-800 transition">
            Papan Lelang
          </a>
          <a href="#weighbridge" class="text-slate-300 hover:text-white px-3 py-1.5 rounded-lg hover:bg-slate-800 transition">
            Jembatan Timbang
          </a>
          <a href="#financing" class="text-slate-300 hover:text-white px-3 py-1.5 rounded-lg hover:bg-slate-800 transition">
            Talangan Modal (SCF)
          </a>
          <a href="#ledger" class="text-slate-300 hover:text-white px-3 py-1.5 rounded-lg hover:bg-slate-800 transition">
            Audit Buku Besar
          </a>
        </div>

        <!-- User Controls / Auth -->
        <div class="flex items-center space-x-2.5">
          
          <!-- Mode Kontras Lapangan Toggle -->
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
            
            <!-- Quick Role Switch Dropdown for Easy Testing -->
            <div class="hidden sm:flex items-center space-x-1 text-[11px] bg-slate-950 border border-slate-800 rounded-lg p-1">
              <span class="text-slate-500 px-1.5 font-medium">Ganti Peran:</span>
              <button
                @click="quickSwitchRole('PABRIK_MANAGER')"
                :class="['px-2 py-0.5 rounded font-semibold transition', user.role === 'PABRIK_MANAGER' ? 'bg-blue-600 text-white' : 'text-slate-400 hover:text-white']"
                title="Masuk sebagai Pabrik"
              >
                Pabrik
              </button>
              <button
                @click="quickSwitchRole('LAPAK_OWNER')"
                :class="['px-2 py-0.5 rounded font-semibold transition', user.role === 'LAPAK_OWNER' ? 'bg-emerald-600 text-white' : 'text-slate-400 hover:text-white']"
                title="Masuk sebagai Juragan Lapak"
              >
                Lapak
              </button>
              <button
                @click="quickSwitchRole('WEIGH_OPERATOR')"
                :class="['px-2 py-0.5 rounded font-semibold transition', user.role === 'WEIGH_OPERATOR' ? 'bg-amber-600 text-white' : 'text-slate-400 hover:text-white']"
                title="Masuk sebagai Operator Timbangan"
              >
                Timbangan
              </button>
            </div>

            <!-- Profile Info -->
            <div class="text-right hidden md:block">
              <p class="text-xs font-bold text-white leading-tight">{{ user.fullName || "Pengguna Aktif" }}</p>
              <p class="text-[10px] text-emerald-400 font-medium">{{ user.companyName || "Perusahaan Rekanan" }}</p>
            </div>

            <!-- Role Badge -->
            <span :class="[
              'px-2.5 py-1 rounded-lg text-[10px] font-bold border uppercase tracking-wider',
              user.role === 'PABRIK_MANAGER' ? 'bg-blue-950 text-blue-300 border-blue-800' :
              user.role === 'LAPAK_OWNER' ? 'bg-emerald-950 text-emerald-300 border-emerald-800' :
              'bg-amber-950 text-amber-300 border-amber-800'
            ]">
              {{ getRoleLabel(user.role) }}
            </span>

            <!-- Logout Button -->
            <button
              @click="handleLogout"
              title="Keluar Akun"
              class="p-2 text-slate-400 hover:text-rose-400 hover:bg-slate-800 rounded-lg transition"
            >
              <LogOut class="w-4 h-4" />
            </button>
          </div>

          <!-- If Not Logged In -->
          <div v-else class="flex items-center space-x-2">
            <RouterLink
              to="/login"
              class="px-3.5 py-1.5 text-xs font-semibold text-slate-200 hover:text-white hover:bg-slate-800 border border-slate-700 rounded-xl transition"
            >
              Masuk
            </RouterLink>
            <RouterLink
              to="/register"
              class="px-3.5 py-1.5 text-xs font-semibold text-slate-950 bg-emerald-400 hover:bg-emerald-300 rounded-xl shadow-sm transition"
            >
              Daftar Usaha
            </RouterLink>
          </div>

        </div>

      </div>
    </div>
  </nav>
</template>
