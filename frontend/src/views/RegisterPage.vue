<script>
import { mapActions, mapState } from "pinia";
import { useMainStore } from "../stores/mainStore";
import { RouterLink } from "vue-router";
import { Building2, ShieldCheck, Mail, Lock, User, FileText, ArrowRight } from "lucide-vue-next";

export default {
  name: "RegisterPage",
  components: {
    RouterLink,
    Building2,
    ShieldCheck,
    Mail,
    Lock,
    User,
    FileText,
    ArrowRight
  },
  data() {
    return {
      legalName: "",
      entityType: "LAPAK",
      npwp: "",
      fullName: "",
      email: "",
      password: "",
      role: "LAPAK_OWNER",
    };
  },
  computed: {
    ...mapState(useMainStore, ["loading"]),
  },
  methods: {
    ...mapActions(useMainStore, ["handleRegister"]),
    onEntityTypeChange() {
      if (this.entityType === "PABRIK") {
        this.role = "PABRIK_MANAGER";
      } else {
        this.role = "LAPAK_OWNER";
      }
    },
    async onSubmit() {
      await this.handleRegister({
        legal_name: this.legalName,
        entity_type: this.entityType,
        npwp: this.npwp,
        full_name: this.fullName,
        email: this.email,
        password: this.password,
        role: this.role,
      });
    },
  },
};
</script>

<template>
  <div class="min-h-screen bg-slate-950 text-slate-100 flex flex-col justify-center items-center px-4 py-8 sm:px-6 lg:px-8 font-sans selection:bg-emerald-500 selection:text-slate-950 relative overflow-hidden">
    
    <div class="w-full max-w-lg text-center space-y-3 z-10">
      <RouterLink to="/" class="inline-flex items-center space-x-2.5">
        <div class="w-10 h-10 rounded-xl bg-emerald-500 flex items-center justify-center font-bold text-slate-950 text-xl shadow-lg shadow-emerald-500/30">
          ⚡
        </div>
        <div class="text-left">
          <span class="text-2xl font-black tracking-tight text-white">Scrap<span class="text-emerald-400">Flow</span></span>
          <span class="block text-[10px] uppercase font-bold tracking-widest text-emerald-400/80">KYC Onboarding</span>
        </div>
      </RouterLink>
      <h1 class="text-xl sm:text-2xl font-bold tracking-tight text-slate-100">Registrasi Rekanan B2B</h1>
      <p class="text-xs text-slate-400">Daftarkan legalitas PT/CV untuk mengikuti lelang scrap & fasilitas talangan modal.</p>
    </div>

    <div class="w-full max-w-lg mt-6 z-10">
      <div class="bg-slate-900/90 border border-slate-800 rounded-2xl p-6 sm:p-8 shadow-2xl backdrop-blur-xl space-y-6">
        
        <form @submit.prevent="onSubmit" class="space-y-4">
          
          <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
            <div>
              <label class="block text-xs font-semibold text-slate-300 mb-1.5">Tipe Entitas</label>
              <select
                v-model="entityType"
                @change="onEntityTypeChange"
                class="w-full bg-slate-950 border border-slate-800 rounded-xl px-3.5 py-2.5 text-sm text-white focus:outline-none focus:border-emerald-500 transition"
              >
                <option value="LAPAK">Penampung / Juragan Lapak</option>
                <option value="PABRIK">Pabrik / Manufaktur / BUMN</option>
                <option value="SMELTER">Pabrik Peleburan (Offtaker)</option>
              </select>
            </div>
            <div>
              <label class="block text-xs font-semibold text-slate-300 mb-1.5">NPWP Badan Usaha</label>
              <input
                v-model="npwp"
                type="text"
                required
                placeholder="01.234.567.8-012.000"
                class="w-full bg-slate-950 border border-slate-800 rounded-xl px-3.5 py-2.5 text-sm text-white font-mono placeholder-slate-600 focus:outline-none focus:border-emerald-500 transition"
              />
            </div>
          </div>

          <div>
            <label class="block text-xs font-semibold text-slate-300 mb-1.5">Nama Legal PT / CV / UD</label>
            <div class="relative">
              <Building2 class="absolute left-3.5 top-3 w-4 h-4 text-slate-500" />
              <input
                v-model="legalName"
                type="text"
                required
                placeholder="PT Sumber Logam Lestari"
                class="w-full bg-slate-950 border border-slate-800 rounded-xl pl-10 pr-3.5 py-2.5 text-sm text-white placeholder-slate-600 focus:outline-none focus:border-emerald-500 transition"
              />
            </div>
          </div>

          <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
            <div>
              <label class="block text-xs font-semibold text-slate-300 mb-1.5">Nama PIC / Penanggung Jawab</label>
              <div class="relative">
                <User class="absolute left-3.5 top-3 w-4 h-4 text-slate-500" />
                <input
                  v-model="fullName"
                  type="text"
                  required
                  placeholder="Budi Setiawan"
                  class="w-full bg-slate-950 border border-slate-800 rounded-xl pl-10 pr-3.5 py-2.5 text-sm text-white placeholder-slate-600 focus:outline-none focus:border-emerald-500 transition"
                />
              </div>
            </div>
            <div>
              <label class="block text-xs font-semibold text-slate-300 mb-1.5">Email Akun</label>
              <div class="relative">
                <Mail class="absolute left-3.5 top-3 w-4 h-4 text-slate-500" />
                <input
                  v-model="email"
                  type="email"
                  required
                  placeholder="budi@perusahaan.com"
                  class="w-full bg-slate-950 border border-slate-800 rounded-xl pl-10 pr-3.5 py-2.5 text-sm text-white placeholder-slate-600 focus:outline-none focus:border-emerald-500 transition"
                />
              </div>
            </div>
          </div>

          <div>
            <label class="block text-xs font-semibold text-slate-300 mb-1.5">Kata Sandi</label>
            <div class="relative">
              <Lock class="absolute left-3.5 top-3 w-4 h-4 text-slate-500" />
              <input
                v-model="password"
                type="password"
                required
                minlength="6"
                placeholder="Minimal 6 karakter"
                class="w-full bg-slate-950 border border-slate-800 rounded-xl pl-10 pr-3.5 py-2.5 text-sm text-white placeholder-slate-600 focus:outline-none focus:border-emerald-500 transition"
              />
            </div>
          </div>

          <button
            type="submit"
            :disabled="loading"
            class="w-full py-3 bg-emerald-500 hover:bg-emerald-400 text-slate-950 font-bold rounded-xl text-sm transition shadow-lg shadow-emerald-500/20 flex items-center justify-center space-x-2 disabled:opacity-50"
          >
            <span v-if="!loading">Kirim Pengajuan KYC & Registrasi</span>
            <span v-else>Memproses...</span>
            <ArrowRight class="w-4 h-4" />
          </button>
        </form>

        <div class="text-center pt-2 border-t border-slate-800/80">
          <p class="text-xs text-slate-400">
            Sudah memiliki akun terverifikasi?
            <RouterLink to="/login" class="text-emerald-400 font-semibold hover:underline ml-1">
              Masuk di sini
            </RouterLink>
          </p>
        </div>

      </div>
    </div>

  </div>
</template>
