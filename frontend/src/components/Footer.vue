<script>
export default {
  name: "Footer",
};
</script>

<template>
  <footer class="py-8 bg-slate-900 border-t border-slate-800 text-slate-400 text-xs">
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 space-y-4">
      <div class="flex flex-col sm:flex-row items-center justify-between gap-4">
        <div class="flex items-center space-x-2">
          <div class="w-6 h-6 rounded-lg bg-emerald-500 flex items-center justify-center font-bold text-slate-950 text-xs">
            ⚡
          </div>
          <span class="text-sm font-bold text-white tracking-tight">Scrap<span class="text-emerald-400">Flow</span> Technologies</span>
        </div>
        <p class="text-slate-500 text-center sm:text-right text-[11px]">
          Platform B2B Commodity Auction, Milestone Escrow Vault, & Supply Chain Financing.
        </p>
      </div>
      <div class="pt-4 border-t border-slate-800/80 flex flex-col sm:flex-row items-center justify-between text-[11px] text-slate-500 gap-2">
        <p>&copy; 2026 ScrapFlow Core Engine. All rights reserved.</p>
        <div class="flex space-x-4">
          <span>OJK Regulatory Sandbox Ready</span>
          <span>•</span>
          <span>Double-Entry Zero-Sum Ledger</span>
          <span>•</span>
          <span>ISO 17025 Calibrated</span>
        </div>
      </div>
    </div>
  </footer>
</template>
