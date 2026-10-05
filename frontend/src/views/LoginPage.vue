<script>
import { mapActions, mapState } from "pinia";
import { useMainStore } from "../stores/mainStore";
import { RouterLink } from "vue-router";
import {
  ShieldCheck,
  Building2,
  Truck,
  Scale,
  Lock,
  Mail,
  ArrowRight,
  Sparkles,
  Zap
} from "lucide-vue-next";

export default {
  name: "LoginPage",
  components: {
    RouterLink,
    ShieldCheck,
    Building2,
    Truck,
    Scale,
    Lock,
    Mail,
    ArrowRight,
    Sparkles,
    Zap
  },
  data() {
    return {
      email: "",
      password: "",
      demoAccounts: [
        {
          role: "PABRIK_MANAGER",
          title: "Pabrik / Pemilik Scrap (PT Krakatau)",
          desc: "Buat tender lelang, pilih pemenang SPK, terima pencairan dana",
          email: "budi@krakatausteel.co.id",
          password: "password123",
          icon: "Building2",
          badge: "Penyedia Lelang",
          color: "border-blue-500/30 bg-blue-950/20 text-blue-400 hover:border-blue-500/60"
        },
        {
          role: "LAPAK_OWNER",
          title: "Juragan Lapak Besi Tua (CV Besi Abadi)",
          desc: "Tawar lelang, kunci jaminan bid-bond, ajukan talangan modal SCF",
          email: "haji.slamet@besiabadi.com",
          password: "password123",
          icon: "Truck",
          badge: "Peserta & Pembeli",
          color: "border-emerald-500/30 bg-emerald-950/20 text-emerald-400 hover:border-emerald-500/60"
        },
        {
          role: "WEIGH_OPERATOR",
          title: "Petugas Jembatan Timbang Lapangan",
          desc: "Input tiket timbang bruto/tara, deteksi kecurangan AI, pencairan per rit",
          email: "operator@timbangan.com",
          password: "password123",
          icon: "Scale",
          badge: "Operator Lapangan",
          color: "border-amber-500/30 bg-amber-950/20 text-amber-400 hover:border-amber-500/60"
        }
      ]
    };
  },
  computed: {
    ...mapState(useMainStore, ["loading"])
  },
  methods: {
    ...mapActions(useMainStore, ["handleLogin"]),
    async selectDemo(acc) {
      this.email = acc.email;
      this.password = acc.password;
      await this.onSubmit();
    },
    async onSubmit() {
      if (!this.email || !this.password) return;
      const ok = await this.handleLogin(this.email, this.password);
      if (ok) {
        this.$router.push("/");
      }
    }
  }
};
</script>

<template>
  <div class="min-h-screen bg-slate-950 text-slate-100 flex flex-col justify-center items-center px-4 py-8 sm:px-6 lg:px-8 font-sans selection:bg-emerald-500 selection:text-slate-950 relative overflow-hidden">
    
    <!-- Background Glow -->
    <div class="absolute top-1/4 left-1/2 -translate-x-1/2 -translate-y-1/2 w-96 h-96 bg-emerald-500/10 rounded-full blur-3xl pointer-events-none"></div>

    <!-- Header & Branding -->
    <div class="w-full max-w-md text-center space-y-3 z-10">
      <RouterLink to="/" class="inline-flex items-center space-x-2.5">
        <div class="w-10 h-10 rounded-xl bg-emerald-500 flex items-center justify-center font-bold text-slate-950 text-xl shadow-lg shadow-emerald-500/30">
          ⚡
        </div>
        <div class="text-left">
          <span class="text-2xl font-black tracking-tight text-white">Scrap<span class="text-emerald-400">Flow</span></span>
          <span class="block text-[10px] uppercase font-bold tracking-widest text-emerald-400/80">Platform Tender & Pembiayaan Besi Tua</span>
        </div>
      </RouterLink>
      <h1 class="text-xl sm:text-2xl font-bold tracking-tight text-slate-100">Masuk ke Akun Anda</h1>
      <p class="text-xs text-slate-400">Akses lelang terverifikasi, rekening bersama (escrow), dan pencairan timbangan.</p>
    </div>

    <!-- Main Card Form -->
    <div class="w-full max-w-md mt-6 z-10">
      <div class="bg-slate-900/90 border border-slate-800 rounded-2xl p-6 sm:p-8 shadow-2xl backdrop-blur-xl space-y-6">
        
        <!-- Standard Form -->
        <form @submit.prevent="onSubmit" class="space-y-4">
          <div>
            <label class="block text-xs font-semibold text-slate-300 mb-1.5">Alamat Email Terdaftar</label>
            <div class="relative">
              <Mail class="absolute left-3.5 top-3 w-4 h-4 text-slate-500" />
              <input
                v-model="email"
                type="email"
                required
                placeholder="nama@perusahaan.co.id"
                class="w-full bg-slate-950 border border-slate-800 rounded-xl pl-10 pr-3.5 py-2.5 text-sm text-white placeholder-slate-600 focus:outline-none focus:border-emerald-500 focus:ring-1 focus:ring-emerald-500 transition font-sans"
              />
            </div>
          </div>

          <div>
            <div class="flex items-center justify-between mb-1.5">
              <label class="block text-xs font-semibold text-slate-300">Kata Sandi</label>
              <a href="#" class="text-[11px] text-emerald-400 hover:underline">Lupa kata sandi?</a>
            </div>
            <div class="relative">
              <Lock class="absolute left-3.5 top-3 w-4 h-4 text-slate-500" />
              <input
                v-model="password"
                type="password"
                required
                placeholder="••••••••"
                class="w-full bg-slate-950 border border-slate-800 rounded-xl pl-10 pr-3.5 py-2.5 text-sm text-white placeholder-slate-600 focus:outline-none focus:border-emerald-500 focus:ring-1 focus:ring-emerald-500 transition"
              />
            </div>
          </div>

          <button
            type="submit"
            :disabled="loading"
            class="w-full py-3 bg-emerald-500 hover:bg-emerald-400 text-slate-950 font-bold rounded-xl text-sm transition shadow-lg shadow-emerald-500/20 flex items-center justify-center space-x-2 disabled:opacity-50"
          >
            <span v-if="!loading">Masuk ke Sistem</span>
            <span v-else>Memverifikasi Akun...</span>
            <ArrowRight class="w-4 h-4" />
          </button>
        </form>

        <!-- Divider -->
        <div class="relative flex items-center justify-center">
          <div class="border-t border-slate-800 w-full"></div>
          <span class="bg-slate-900 px-3 text-[11px] uppercase tracking-wider text-slate-400 font-semibold absolute">
            Pilih Peran Pengguna (1-Klik Demo)
          </span>
        </div>

        <!-- 1-Click Demo Account Switcher -->
        <div class="space-y-2.5">
          <p class="text-[11px] text-slate-400 text-center font-medium">Klik profil di bawah untuk menguji hak akses masing-masing peran:</p>
          
          <div class="grid grid-cols-1 gap-2.5">
            <button
              v-for="acc in demoAccounts"
              :key="acc.role"
              type="button"
              @click="selectDemo(acc)"
              :class="['w-full p-3 rounded-xl border flex flex-col sm:flex-row sm:items-center justify-between text-left transition group gap-2', acc.color]"
            >
              <div class="flex items-start space-x-2.5">
                <component :is="acc.icon" class="w-4 h-4 flex-shrink-0 mt-0.5" />
                <div>
                  <div class="flex items-center space-x-2">
                    <p class="text-xs font-bold text-slate-100 group-hover:text-emerald-300 transition">{{ acc.title }}</p>
                  </div>
                  <p class="text-[10px] text-slate-400 leading-tight mt-0.5">{{ acc.desc }}</p>
                  <p class="text-[10px] text-slate-500 font-mono mt-1">{{ acc.email }}</p>
                </div>
              </div>
              <span class="self-start sm:self-center text-[10px] px-2 py-0.5 rounded-md bg-slate-900/90 border border-slate-700 font-semibold whitespace-nowrap">
                {{ acc.badge }}
              </span>
            </button>
          </div>
        </div>

        <!-- Register Link -->
        <div class="text-center pt-2 border-t border-slate-800/80">
          <p class="text-xs text-slate-400">
            Perusahaan Anda belum terdaftar?
            <RouterLink to="/register" class="text-emerald-400 font-semibold hover:underline ml-1">
              Daftar Profil Usaha & Verifikasi KYC
            </RouterLink>
          </p>
        </div>

      </div>
    </div>

    <!-- Credibility & Trust Badges Footer -->
    <div class="w-full max-w-md mt-6 grid grid-cols-3 gap-2 text-center text-[10px] text-slate-400 z-10">
      <div class="flex flex-col items-center space-y-1 p-2 rounded-lg bg-slate-900/50 border border-slate-800">
        <ShieldCheck class="w-4 h-4 text-emerald-400" />
        <span class="font-medium text-slate-200">Sandbox OJK</span>
        <span class="text-[9px] text-slate-500">Escrow Komoditas</span>
      </div>
      <div class="flex flex-col items-center space-y-1 p-2 rounded-lg bg-slate-900/50 border border-slate-800">
        <Lock class="w-4 h-4 text-emerald-400" />
        <span class="font-medium text-slate-200">Buku Besar Imbang</span>
        <span class="text-[9px] text-slate-500">Audit Zero-Sum</span>
      </div>
      <div class="flex flex-col items-center space-y-1 p-2 rounded-lg bg-slate-900/50 border border-slate-800">
        <Scale class="w-4 h-4 text-emerald-400" />
        <span class="font-medium text-slate-200">ISO 17025</span>
        <span class="text-[9px] text-slate-500">Tera Jembatan Timbang</span>
      </div>
    </div>

  </div>
</template>
