<script>
import { mapActions, mapState } from "pinia";
import { useMainStore } from "../stores/mainStore";
import { RouterLink } from "vue-router";

export default {
  name: "Navbar",
  components: { RouterLink },
  methods: {
    ...mapActions(useMainStore, ["handleLogout"]),
  },
  computed: {
    ...mapState(useMainStore, ["isLogin"]),
  },
};
</script>

<template>
  <nav class="bg-slate-900 border-b border-slate-800 text-slate-100 sticky top-0 z-50">
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
      <div class="flex items-center justify-between h-16">
        <!-- Logo & Branding -->
        <div class="flex items-center space-x-3">
          <RouterLink to="/" class="flex items-center space-x-2">
            <div class="w-9 h-9 rounded-lg bg-emerald-500 flex items-center justify-center font-bold text-slate-950 text-xl shadow-lg shadow-emerald-500/20">
              ⚡
            </div>
            <div>
              <span class="text-xl font-bold tracking-tight text-white">Scrap<span class="text-emerald-400">Flow</span></span>
              <span class="ml-2 px-2 py-0.5 text-[10px] font-semibold bg-emerald-950 text-emerald-400 border border-emerald-800 rounded-full uppercase tracking-wider">Fintech Core</span>
            </div>
          </RouterLink>

          <!-- System Status Badge -->
          <div class="hidden md:flex items-center space-x-2 ml-6 pl-6 border-l border-slate-800 text-xs text-slate-400">
            <span class="w-2 h-2 rounded-full bg-emerald-400 animate-pulse"></span>
            <span>Golang Ledger Engine: <strong class="text-slate-200">ACTIVE</strong></span>
          </div>
        </div>

        <!-- Navigation Links -->
        <div class="flex items-center space-x-3 text-sm">
          <a href="#tenders" class="text-slate-300 hover:text-white px-3 py-1.5 rounded-md hover:bg-slate-800 transition">
            Tender Board
          </a>
          <a href="#weighbridge" class="text-slate-300 hover:text-white px-3 py-1.5 rounded-md hover:bg-slate-800 transition">
            Jembatan Timbang
          </a>
          <a href="#financing" class="text-slate-300 hover:text-white px-3 py-1.5 rounded-md hover:bg-slate-800 transition">
            Talangan SCF
          </a>
          <a href="#ledger" class="text-slate-300 hover:text-white px-3 py-1.5 rounded-md hover:bg-slate-800 transition">
            Audit Ledger
          </a>

          <div class="pl-3 border-l border-slate-800 flex items-center space-x-2">
            <RouterLink
              v-if="!isLogin"
              to="/login"
              class="px-3.5 py-1.5 text-xs font-semibold text-slate-300 hover:text-white hover:bg-slate-800 border border-slate-700 rounded-lg transition"
            >
              Sign In
            </RouterLink>
            <RouterLink
              v-if="!isLogin"
              to="/register"
              class="px-3.5 py-1.5 text-xs font-semibold text-slate-950 bg-emerald-400 hover:bg-emerald-300 rounded-lg shadow-sm transition"
            >
              Register KYC
            </RouterLink>
            <button
              v-if="isLogin"
              @click="handleLogout()"
              class="px-3 py-1.5 text-xs font-medium text-rose-300 hover:bg-rose-950/50 border border-rose-900 rounded-lg transition"
            >
              Logout
            </button>
          </div>
        </div>
      </div>
    </div>
  </nav>
</template>
