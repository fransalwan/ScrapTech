<script>
import Navbar from "../components/Navbar.vue";
import Footer from "../components/Footer.vue";
import { mapActions, mapState } from "pinia";
import { useMainStore } from "../stores/mainStore";
import { toast } from "vue-sonner";
import {
  Scale,
  Truck,
  TrendingUp,
  ShieldCheck,
  Building2,
  Copy,
  QrCode,
  AlertTriangle,
  CheckCircle2,
  Lock,
  ArrowUpRight,
  FileText,
  Coins,
  RefreshCw,
  Camera
} from "lucide-vue-next";

export default {
  name: "HomePage",
  components: {
    Navbar,
    Footer,
    Scale,
    Truck,
    TrendingUp,
    ShieldCheck,
    Building2,
    Copy,
    QrCode,
    AlertTriangle,
    CheckCircle2,
    Lock,
    ArrowUpRight,
    FileText,
    Coins,
    RefreshCw,
    Camera
  },
  data() {
    return {
      // Live Commodity Scrap Spot Ticker
      spotPrices: [
        { name: "HMS-1 (Baja Super)", price: 6350, change: "+1.2%", isUp: true },
        { name: "HMS-2 (Besi Campur)", price: 5850, change: "+0.8%", isUp: true },
        { name: "Plat Kapal / Boiler", price: 6750, change: "+2.1%", isUp: true },
        { name: "Besi Cor (Mesin)", price: 5900, change: "-0.4%", isUp: false },
        { name: "Tembaga Super Kupas", price: 96000, change: "+1.5%", isUp: true },
      ],
      // Mock Tenders
      tenders: [
        {
          id: "tnd-001",
          number: "TND-202610-881",
          company: "PT Krakatau Industrial Scrap",
          title: "Bongkaran Pabrik Plat Baja & HMS-1",
          type: "HMS_1",
          weight: 450000,
          reservePrice: 6200,
          currentBid: 6350,
          bidBond: 50000000,
          location: "Cilegon, Banten",
          deadline: "2 Hari Lagi",
          status: "BIDDING",
          vaNumber: "88081948201948",
        },
        {
          id: "tnd-002",
          number: "TND-202610-882",
          company: "PT Dok & Perkapalan Surabaya",
          title: "Scrap Potongan Plat Kapal & Wire Rope",
          type: "PLATE",
          weight: 280000,
          reservePrice: 6600,
          currentBid: 6750,
          bidBond: 35000000,
          location: "Tanjung Perak, Surabaya",
          deadline: "18 Jam Lagi",
          status: "BIDDING",
          vaNumber: "88082947192837",
        },
        {
          id: "tnd-003",
          number: "TND-202610-883",
          company: "PT Pindad Persero Material Afkir",
          title: "Besi Cor / Cast Iron Blok Mesin",
          type: "CAST_IRON",
          weight: 120000,
          reservePrice: 5800,
          currentBid: 5900,
          bidBond: 20000000,
          location: "Bandung, Jawa Barat",
          deadline: "4 Hari Lagi",
          status: "BIDDING",
          vaNumber: "88083827194821",
        },
      ],
      // Interactive Weighbridge Simulation State
      weighSim: {
        truckPlate: "B 9482 UYX",
        grossWeight: 28450,
        tareWeight: 11200,
        refractionPct: 2.0,
        pricePerKg: 6350,
        feePct: 0.5,
        settled: false,
        fraudAlert: false,
        fraudMessage: "",
        journal: null,
      },
      // Supply Chain Financing Simulator
      scfSim: {
        contractValue: 2857500000,
        lapakEquityPct: 30,
        tenorDays: 30,
        annualRate: 12.0,
      },
      // Live Ledger Journal Entries
      ledgerJournals: [
        {
          id: "JRN-ESCROW-001",
          type: "ESCROW_LOCK",
          ref: "TND-202610-881",
          desc: "Bid-bond locked from Juragan Besi Abadi",
          debitAcc: "1010 - Bank Custody",
          creditAcc: "2010 - Escrow Bid-Bond Payable",
          amount: 50000000,
          time: "10:14:22",
        },
        {
          id: "JRN-SCF-002",
          type: "SCF_DISBURSE",
          ref: "SPK-202610-091",
          desc: "Talangan tender SPK injected into milestone escrow",
          debitAcc: "1020 - SCF Loan Receivables",
          creditAcc: "2020 - Escrow Milestone Payable",
          amount: 1500000000,
          time: "09:30:10",
        },
      ],
      bidModal: {
        isOpen: false,
        tender: null,
        bidPrice: 0,
      },
      qrModal: {
        isOpen: false,
        tender: null,
      },
    };
  },
  computed: {
    ...mapState(useMainStore, ["isLogin", "user"]),
    netWeight() {
      return Math.max(0, this.weighSim.grossWeight - this.weighSim.tareWeight);
    },
    refractionKg() {
      return (this.netWeight * this.weighSim.refractionPct) / 100;
    },
    billableWeight() {
      return Math.max(0, this.netWeight - this.refractionKg);
    },
    grossPayable() {
      return Math.round(this.billableWeight * this.weighSim.pricePerKg);
    },
    platformFee() {
      return Math.round(this.grossPayable * (this.weighSim.feePct / 100));
    },
    netToSeller() {
      return this.grossPayable - this.platformFee;
    },
    scfLoanAmount() {
      return Math.round(this.scfSim.contractValue * ((100 - this.scfSim.lapakEquityPct) / 100));
    },
    scfInterestCost() {
      return Math.round(this.scfLoanAmount * (this.scfSim.annualRate / 100) * (this.scfSim.tenorDays / 365));
    },
  },
  methods: {
    formatRupiah(val) {
      return new Intl.NumberFormat("id-ID", {
        style: "currency",
        currency: "IDR",
        maximumFractionDigits: 0,
      }).format(val);
    },
    copyText(text, label) {
      navigator.clipboard.writeText(text);
      toast.success(`${label} Disalin!`, {
        description: text,
      });
    },
    openBidModal(tender) {
      this.bidModal.tender = tender;
      this.bidModal.bidPrice = tender.currentBid + 50;
      this.bidModal.isOpen = true;
    },
    openQRModal(tender) {
      this.qrModal.tender = tender;
      this.qrModal.isOpen = true;
    },
    confirmBid() {
      if (this.bidModal.bidPrice <= this.bidModal.tender.currentBid) {
        toast.error("Penawaran Tidak Valid", {
          description: "Harga penawaran harus lebih tinggi dari penawaran saat ini.",
        });
        return;
      }
      this.bidModal.tender.currentBid = this.bidModal.bidPrice;

      const newJournal = {
        id: `JRN-BID-${Date.now().toString().slice(-4)}`,
        type: "BID_BOND_DEPOSIT",
        ref: this.bidModal.tender.number,
        desc: `Bid-bond locked for ${this.bidModal.tender.title}`,
        debitAcc: "1010 - Bank Custody",
        creditAcc: "2010 - Escrow Bid-Bond Payable",
        amount: this.bidModal.tender.bidBond,
        time: new Date().toLocaleTimeString("id-ID"),
      };
      this.ledgerJournals.unshift(newJournal);

      toast.success("Bidding Berhasil Dicatat!", {
        description: `Penawaran ${this.formatRupiah(this.bidModal.bidPrice)}/Kg diterima. Jaminan ${this.formatRupiah(this.bidModal.tender.bidBond)} dikunci di Escrow Vault.`,
      });

      this.bidModal.isOpen = false;
    },
    simulateAIScan() {
      // Simulate OCR with slight tare anomaly to demonstrate AI credibility
      this.weighSim.tareWeight = 12450; // Anomaly > 11% deviation
      this.weighSim.fraudAlert = true;
      this.weighSim.fraudMessage = "CRITICAL ANOMALY: Deviasi berat tara 11.16% terdeteksi dari rata-rata historis truk (12.450 Kg vs 11.200 Kg). Tangki air atau modifikasi tara terindikasi.";
      toast.warning("AI Fraud Detection Alert!", {
        description: "Deviasi berat tara melebihi toleransi normal 3%.",
      });
    },
    executeWeighbridgeSettlement() {
      this.weighSim.settled = true;
      const journalId = `JRN-SETTLE-${Date.now().toString().slice(-4)}`;
      this.weighSim.journal = {
        id: journalId,
        ref: `WB-${this.weighSim.truckPlate.replace(/\s+/g, "")}`,
        debit: { account: "2020 - Escrow Milestone Payable", amount: this.grossPayable },
        creditSeller: { account: "2030 - Seller Withdrawable Balance", amount: this.netToSeller },
        creditFee: { account: "4010 - Escrow Fee Revenue", amount: this.platformFee },
      };

      this.ledgerJournals.unshift({
        id: journalId,
        type: "WEIGHBRIDGE_SETTLE",
        ref: `WB-${this.weighSim.truckPlate}`,
        desc: `Settlement ritase truk ${this.weighSim.truckPlate} (${this.billableWeight.toFixed(0)} Kg)`,
        debitAcc: "2020 - Escrow Milestone Payable",
        creditAcc: "2030 - Seller Balance",
        amount: this.grossPayable,
        time: new Date().toLocaleTimeString("id-ID"),
      });

      toast.success("Settlement Atomik Berhasil!", {
        description: `Dana ${this.formatRupiah(this.netToSeller)} dicairkan ke saldo pabrik via Double-Entry Ledger.`,
      });
    },
    resetWeighSim() {
      this.weighSim.settled = false;
      this.weighSim.fraudAlert = false;
      this.weighSim.journal = null;
      this.weighSim.tareWeight = 11200;
    },
  },
};
</script>

<template>
  <div class="min-h-screen bg-slate-950 text-slate-100 flex flex-col font-sans selection:bg-emerald-500 selection:text-slate-950 pb-16 md:pb-0">
    <Navbar />

    <!-- LIVE SPOT PRICE TICKER (COMMODITY BAR) -->
    <div class="bg-slate-900 border-b border-slate-800/80 py-2 px-4 overflow-x-auto no-scrollbar">
      <div class="max-w-7xl mx-auto flex items-center space-x-6 text-xs whitespace-nowrap">
        <div class="flex items-center space-x-1.5 text-emerald-400 font-bold">
          <span class="w-2 h-2 rounded-full bg-emerald-400 animate-pulse"></span>
          <span class="uppercase tracking-wider text-[11px]">Harga Besi Hari Ini:</span>
        </div>
        <div v-for="item in spotPrices" :key="item.name" class="flex items-center space-x-2 font-mono">
          <span class="text-slate-400">{{ item.name }}:</span>
          <span class="text-white font-bold">Rp {{ item.price.toLocaleString() }}/Kg</span>
          <span :class="['text-[10px] px-1 rounded', item.isUp ? 'text-emerald-400 bg-emerald-950/60' : 'text-rose-400 bg-rose-950/60']">
            {{ item.change }}
          </span>
        </div>
      </div>
    </div>

    <!-- HERO SECTION -->
    <header class="relative overflow-hidden border-b border-slate-800 bg-gradient-to-b from-slate-900 via-slate-950 to-slate-950 py-10 sm:py-16">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        
        <!-- Trust Seals & Badges -->
        <div class="flex flex-wrap items-center justify-center gap-2 mb-4">
          <span class="inline-flex items-center space-x-1 px-2.5 py-1 rounded-full bg-emerald-950/80 border border-emerald-800 text-emerald-400 text-[11px] font-semibold">
            <ShieldCheck class="w-3.5 h-3.5" />
            <span>OJK Sandbox Ready</span>
          </span>
          <span class="inline-flex items-center space-x-1 px-2.5 py-1 rounded-full bg-blue-950/80 border border-blue-800 text-blue-400 text-[11px] font-semibold">
            <Lock class="w-3.5 h-3.5" />
            <span>Bank Custody Vault</span>
          </span>
          <span class="inline-flex items-center space-x-1 px-2.5 py-1 rounded-full bg-amber-950/80 border border-amber-800 text-amber-400 text-[11px] font-semibold">
            <Scale class="w-3.5 h-3.5" />
            <span>ISO 17025 Calibrated</span>
          </span>
        </div>

        <div class="text-center max-w-3xl mx-auto space-y-3">
          <h1 class="text-3xl sm:text-5xl font-extrabold tracking-tight text-white leading-tight">
            Hubungkan Tender Besi Tua dengan <span class="text-emerald-400">Escrow & Financing</span>
          </h1>
          <p class="text-sm sm:text-base text-slate-400 leading-relaxed max-w-2xl mx-auto">
            Platform B2B resmi: <strong class="text-slate-200">Rekening Bersama Multi-Tahap</strong>, pencairan otomatis berbasis <strong class="text-slate-200">Jembatan Timbang Digital</strong>, dan talangan <strong class="text-slate-200">Supply Chain Financing</strong>.
          </p>
        </div>

        <!-- KEY METRICS CARDS -->
        <div class="grid grid-cols-2 lg:grid-cols-4 gap-3 sm:gap-4 mt-8 sm:mt-12">
          <div class="p-3.5 sm:p-4 rounded-xl bg-slate-900/90 border border-slate-800">
            <p class="text-[11px] text-slate-400 uppercase font-semibold">Dana Escrow Terkunci</p>
            <p class="text-xl sm:text-2xl font-black text-white mt-1">Rp 4.850.000.000</p>
            <p class="text-[10px] text-emerald-400 mt-1 flex items-center">
              <CheckCircle2 class="w-3 h-3 mr-1" /> Zero-Sum Verified
            </p>
          </div>
          <div class="p-3.5 sm:p-4 rounded-xl bg-slate-900/90 border border-slate-800">
            <p class="text-[11px] text-slate-400 uppercase font-semibold">Tender Terbuka Aktif</p>
            <p class="text-xl sm:text-2xl font-black text-white mt-1">12 Lot Pabrik</p>
            <p class="text-[10px] text-slate-400 mt-1">HMS, Plat, Cor, Kabel</p>
          </div>
          <div class="p-3.5 sm:p-4 rounded-xl bg-slate-900/90 border border-slate-800">
            <p class="text-[11px] text-slate-400 uppercase font-semibold">Tonase Selesai Timbang</p>
            <p class="text-xl sm:text-2xl font-black text-white mt-1">1.840 Ton</p>
            <p class="text-[10px] text-emerald-400 mt-1 flex items-center">
              <CheckCircle2 class="w-3 h-3 mr-1" /> Tiket Timbang Valid
            </p>
          </div>
          <div class="p-3.5 sm:p-4 rounded-xl bg-slate-900/90 border border-slate-800">
            <p class="text-[11px] text-slate-400 uppercase font-semibold">Plafon Talangan (SCF)</p>
            <p class="text-xl sm:text-2xl font-black text-white mt-1">Rp 2.150.000.000</p>
            <p class="text-[10px] text-slate-400 mt-1">Likuiditas SPK Tender</p>
          </div>
        </div>

      </div>
    </header>

    <!-- MAIN BODY -->
    <main class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8 sm:py-12 space-y-12 sm:space-y-16 flex-1 w-full">

      <!-- SECTION 1: TENDER BOARD -->
      <section id="tenders" class="space-y-4 sm:space-y-6">
        <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-2">
          <div>
            <h2 class="text-xl sm:text-2xl font-bold text-white tracking-tight flex items-center">
              <Building2 class="w-6 h-6 mr-2 text-emerald-400" />
              Papan Tender Scrap Aktif
            </h2>
            <p class="text-xs sm:text-sm text-slate-400">Pelelangan langsung dari pabrik & BUMN resmi tanpa calo perantara.</p>
          </div>
          <div class="flex items-center space-x-2">
            <span class="text-xs px-2.5 py-1 bg-slate-800 border border-slate-700 text-slate-300 rounded-lg font-mono">
              Golang Gin Core Engine
            </span>
          </div>
        </div>

        <div class="grid grid-cols-1 md:grid-cols-3 gap-4 sm:gap-6">
          <div
            v-for="tender in tenders"
            :key="tender.id"
            class="bg-slate-900/90 border border-slate-800 rounded-2xl p-4 sm:p-5 flex flex-col justify-between hover:border-slate-700 transition shadow-xl"
          >
            <div class="space-y-3">
              <div class="flex items-center justify-between">
                <span class="text-xs font-mono text-emerald-400 font-bold">{{ tender.number }}</span>
                <span class="text-[10px] px-2 py-0.5 rounded bg-emerald-950 text-emerald-300 border border-emerald-800 font-medium">
                  {{ tender.deadline }}
                </span>
              </div>
              
              <h3 class="text-base sm:text-lg font-bold text-white leading-snug">{{ tender.title }}</h3>
              <p class="text-xs text-slate-400">{{ tender.company }} • {{ tender.location }}</p>

              <div class="bg-slate-950/80 p-3 rounded-xl border border-slate-800/80 space-y-2 text-xs">
                <div class="flex justify-between">
                  <span class="text-slate-400">Estimasi Tonase:</span>
                  <strong class="text-slate-200">{{ (tender.weight / 1000).toLocaleString() }} Ton</strong>
                </div>
                <div class="flex justify-between">
                  <span class="text-slate-400">Harga Dasar:</span>
                  <span class="text-slate-200">Rp {{ tender.reservePrice.toLocaleString() }}/Kg</span>
                </div>
                <div class="flex justify-between">
                  <span class="text-slate-400">Penawaran Tertinggi:</span>
                  <strong class="text-emerald-400 text-sm">Rp {{ tender.currentBid.toLocaleString() }}/Kg</strong>
                </div>
                <div class="flex justify-between pt-1 border-t border-slate-800">
                  <span class="text-slate-400">Jaminan Bid-Bond:</span>
                  <strong class="text-amber-400">{{ formatRupiah(tender.bidBond) }}</strong>
                </div>
              </div>

              <!-- Quick VA Number Copy for Credibility -->
              <div class="p-2 bg-slate-900 border border-slate-800 rounded-lg flex items-center justify-between text-[11px]">
                <span class="text-slate-400">Virtual Account Escrow:</span>
                <button
                  @click="copyText(tender.vaNumber, 'Nomor VA')"
                  class="font-mono text-emerald-400 font-bold hover:underline flex items-center space-x-1"
                >
                  <span>{{ tender.vaNumber }}</span>
                  <Copy class="w-3 h-3 ml-1" />
                </button>
              </div>
            </div>

            <div class="mt-4 pt-3 border-t border-slate-800 flex items-center space-x-2">
              <button
                @click="openQRModal(tender)"
                class="p-2.5 bg-slate-800 hover:bg-slate-700 text-slate-300 rounded-xl transition"
                title="Lihat Surat Jalan Digital / QR"
              >
                <QrCode class="w-4 h-4" />
              </button>
              <button
                @click="openBidModal(tender)"
                class="flex-1 py-2.5 px-3 bg-emerald-500 hover:bg-emerald-400 text-slate-950 font-bold rounded-xl text-xs uppercase tracking-wide transition shadow-lg shadow-emerald-500/10 text-center"
              >
                Tawar (Lock Bid-Bond)
              </button>
            </div>
          </div>
        </div>
      </section>

      <!-- SECTION 2: DIGITAL WEIGHBRIDGE SIMULATOR -->
      <section id="weighbridge" class="bg-slate-900 border border-slate-800 rounded-2xl p-5 sm:p-8 space-y-6">
        <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-2 border-b border-slate-800 pb-4">
          <div>
            <h2 class="text-xl sm:text-2xl font-bold text-white flex items-center">
              <Scale class="w-6 h-6 mr-2 text-emerald-400" />
              Simulasi Jembatan Timbang & Settlement Otomatis
            </h2>
            <p class="text-xs sm:text-sm text-slate-400">Kalkulasi refraksi & pencairan dana langsung per ritase truk ke rekening bank pabrik.</p>
          </div>
          <button
            @click="simulateAIScan"
            class="inline-flex items-center space-x-1.5 px-3 py-1.5 bg-slate-800 hover:bg-slate-700 text-emerald-400 border border-slate-700 rounded-xl text-xs font-semibold transition"
          >
            <Camera class="w-3.5 h-3.5" />
            <span>Simulasi Scan OCR Slip</span>
          </button>
        </div>

        <!-- Fraud Anomaly Banner (If Triggered) -->
        <div v-if="weighSim.fraudAlert" class="p-3.5 bg-rose-950/50 border border-rose-800/80 rounded-xl text-xs text-rose-300 flex items-start space-x-2.5 animate-pulse">
          <AlertTriangle class="w-5 h-5 flex-shrink-0 text-rose-400 mt-0.5" />
          <div>
            <p class="font-bold">⚠️ Fraud Anomaly Detected oleh AI Engine:</p>
            <p class="mt-0.5 text-rose-200">{{ weighSim.fraudMessage }}</p>
          </div>
        </div>

        <div class="grid grid-cols-1 lg:grid-cols-2 gap-6 items-start">
          
          <!-- Form Input -->
          <div class="space-y-3.5 bg-slate-950 p-4 sm:p-5 rounded-xl border border-slate-800">
            <h3 class="text-xs font-semibold text-slate-300 uppercase tracking-wider">Input Tiket Timbangan Lapangan</h3>
            
            <div class="grid grid-cols-2 gap-3">
              <div>
                <label class="block text-xs text-slate-400 mb-1">Nomor Plat Truk</label>
                <input
                  v-model="weighSim.truckPlate"
                  type="text"
                  class="w-full bg-slate-900 border border-slate-700 rounded-lg px-3 py-2 text-xs sm:text-sm text-white font-mono focus:border-emerald-500 focus:outline-none"
                />
              </div>
              <div>
                <label class="block text-xs text-slate-400 mb-1">Harga Kontrak / Kg</label>
                <input
                  v-model.number="weighSim.pricePerKg"
                  type="number"
                  class="w-full bg-slate-900 border border-slate-700 rounded-lg px-3 py-2 text-xs sm:text-sm text-white focus:border-emerald-500 focus:outline-none"
                />
              </div>
              <div>
                <label class="block text-xs text-slate-400 mb-1">Berat Bruto (Truk + Muatan)</label>
                <div class="relative">
                  <input
                    v-model.number="weighSim.grossWeight"
                    type="number"
                    class="w-full bg-slate-900 border border-slate-700 rounded-lg px-3 py-2 text-xs sm:text-sm text-white focus:border-emerald-500 focus:outline-none"
                  />
                  <span class="absolute right-3 top-2 text-xs text-slate-500">Kg</span>
                </div>
              </div>
              <div>
                <label class="block text-xs text-slate-400 mb-1">Berat Tara (Truk Kosong)</label>
                <div class="relative">
                  <input
                    v-model.number="weighSim.tareWeight"
                    type="number"
                    class="w-full bg-slate-900 border border-slate-700 rounded-lg px-3 py-2 text-xs sm:text-sm text-white focus:border-emerald-500 focus:outline-none"
                  />
                  <span class="absolute right-3 top-2 text-xs text-slate-500">Kg</span>
                </div>
              </div>
            </div>

            <div>
              <label class="block text-xs text-slate-400 mb-1">Potongan Refraksi Kotoran / Karat (%)</label>
              <input
                v-model.number="weighSim.refractionPct"
                type="number"
                step="0.5"
                class="w-full bg-slate-900 border border-slate-700 rounded-lg px-3 py-2 text-xs sm:text-sm text-white focus:border-emerald-500 focus:outline-none"
              />
            </div>

            <div class="pt-2">
              <button
                v-if="!weighSim.settled"
                @click="executeWeighbridgeSettlement"
                class="w-full py-3 bg-emerald-500 hover:bg-emerald-400 text-slate-950 font-bold rounded-xl text-xs sm:text-sm transition shadow-lg shadow-emerald-500/20"
              >
                Eksekusi Pencairan Dana (Double-Entry Posting)
              </button>
              <button
                v-else
                @click="resetWeighSim"
                class="w-full py-2.5 bg-slate-800 hover:bg-slate-700 text-slate-300 font-semibold rounded-xl text-xs transition"
              >
                Reset Simulasi
              </button>
            </div>
          </div>

          <!-- Calculation Summary -->
          <div class="space-y-4">
            <div class="bg-slate-950 p-4 sm:p-5 rounded-xl border border-slate-800 space-y-3">
              <h3 class="text-xs font-semibold text-slate-300 uppercase tracking-wider">Hasil Kalkulasi Timbangan</h3>
              
              <div class="grid grid-cols-3 gap-2 text-center py-2 bg-slate-900 rounded-lg border border-slate-800 text-xs">
                <div>
                  <p class="text-slate-400">Netto Awal</p>
                  <strong class="text-white text-sm">{{ netWeight.toLocaleString() }} Kg</strong>
                </div>
                <div>
                  <p class="text-slate-400">Refraksi ({{ weighSim.refractionPct }}%)</p>
                  <strong class="text-rose-400 text-sm">-{{ refractionKg.toFixed(1) }} Kg</strong>
                </div>
                <div>
                  <p class="text-slate-400">Netto Tagihan</p>
                  <strong class="text-emerald-400 text-sm">{{ billableWeight.toFixed(1) }} Kg</strong>
                </div>
              </div>

              <div class="space-y-1.5 text-xs pt-1">
                <div class="flex justify-between">
                  <span class="text-slate-400">Nilai Bruto Tagihan:</span>
                  <span class="text-white font-mono">{{ formatRupiah(grossPayable) }}</span>
                </div>
                <div class="flex justify-between">
                  <span class="text-slate-400">Platform Escrow Fee (0.5%):</span>
                  <span class="text-amber-400 font-mono">{{ formatRupiah(platformFee) }}</span>
                </div>
                <div class="flex justify-between pt-2 border-t border-slate-800 text-sm">
                  <span class="text-slate-300 font-bold">Cair ke Rekening Pabrik:</span>
                  <strong class="text-emerald-400 font-mono">{{ formatRupiah(netToSeller) }}</strong>
                </div>
              </div>
            </div>

            <!-- Double-Entry Transaction Voucher -->
            <div v-if="weighSim.settled && weighSim.journal" class="bg-emerald-950/40 border border-emerald-800/80 p-4 rounded-xl space-y-2 text-xs">
              <div class="flex justify-between items-center">
                <span class="font-mono text-emerald-400 font-bold">VOUCHER JURNAL: {{ weighSim.journal.id }}</span>
                <span class="px-2 py-0.5 bg-emerald-900 text-emerald-300 rounded text-[10px] font-semibold">POSTED</span>
              </div>
              <div class="font-mono text-[11px] bg-slate-950 p-2.5 rounded border border-slate-800 space-y-1">
                <div class="text-rose-400">[DEBIT]  {{ weighSim.journal.debit.account }} : {{ formatRupiah(weighSim.journal.debit.amount) }}</div>
                <div class="text-emerald-400">[CREDIT] {{ weighSim.journal.creditSeller.account }} : {{ formatRupiah(weighSim.journal.creditSeller.amount) }}</div>
                <div class="text-amber-400">[CREDIT] {{ weighSim.journal.creditFee.account }} : {{ formatRupiah(weighSim.journal.creditFee.amount) }}</div>
              </div>
            </div>
          </div>

        </div>
      </section>

      <!-- SECTION 3: SCF (TALANGAN TENDER) CALCULATOR -->
      <section id="financing" class="bg-slate-900/80 border border-slate-800 rounded-2xl p-5 sm:p-8 space-y-6">
        <div>
          <h2 class="text-xl sm:text-2xl font-bold text-white flex items-center">
            <Coins class="w-6 h-6 mr-2 text-emerald-400" />
            Supply Chain Financing (Talangan SPK Tender)
          </h2>
          <p class="text-xs sm:text-sm text-slate-400">Fasilitas modal kerja bagi juragan lapak untuk menebus lelang bernilai miliaran.</p>
        </div>

        <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
          <div class="bg-slate-950 p-4 rounded-xl border border-slate-800 space-y-2">
            <label class="block text-xs text-slate-400">Nilai Kontrak SPK Tender</label>
            <input
              v-model.number="scfSim.contractValue"
              type="number"
              step="50000000"
              class="w-full bg-slate-900 border border-slate-700 rounded-lg px-3 py-2 text-xs sm:text-sm text-white focus:border-emerald-500 focus:outline-none"
            />
          </div>

          <div class="bg-slate-950 p-4 rounded-xl border border-slate-800 space-y-2">
            <label class="block text-xs text-slate-400">Modal Sendiri Lapak (DP {{ scfSim.lapakEquityPct }}%)</label>
            <input
              v-model.number="scfSim.lapakEquityPct"
              type="range"
              min="20"
              max="50"
              step="5"
              class="w-full accent-emerald-500"
            />
            <p class="text-xs text-emerald-400 font-bold">{{ formatRupiah(scfSim.contractValue * (scfSim.lapakEquityPct / 100)) }}</p>
          </div>

          <div class="bg-slate-950 p-4 rounded-xl border border-slate-800 space-y-2">
            <label class="block text-xs text-slate-400">Tenor Pembiayaan (Hari)</label>
            <select
              v-model.number="scfSim.tenorDays"
              class="w-full bg-slate-900 border border-slate-700 rounded-lg px-3 py-2 text-xs sm:text-sm text-white focus:border-emerald-500 focus:outline-none"
            >
              <option :value="14">14 Hari (Ritase Kilat)</option>
              <option :value="30">30 Hari (Standar Lelang)</option>
              <option :value="45">45 Hari (Bongkaran Pabrik)</option>
            </select>
          </div>
        </div>

        <div class="p-4 rounded-xl bg-emerald-950/20 border border-emerald-900/40 flex flex-col sm:flex-row sm:items-center justify-between gap-4">
          <div>
            <p class="text-[11px] text-emerald-400 uppercase font-bold">Plafon Talangan Yang Disediakan Platform</p>
            <p class="text-2xl sm:text-3xl font-black text-white mt-1">{{ formatRupiah(scfLoanAmount) }}</p>
            <p class="text-xs text-slate-400 mt-1">Estimasi Bunga ({{ scfSim.annualRate }}% p.a.): <strong class="text-amber-400">{{ formatRupiah(scfInterestCost) }}</strong></p>
          </div>
          <button
            @click="toast.success('Pengajuan Terkirim!', { description: 'Tim Underwriting SCF akan memverifikasi SPK dalam 24 jam.' })"
            class="px-5 py-2.5 bg-emerald-500 hover:bg-emerald-400 text-slate-950 font-bold rounded-xl text-xs uppercase tracking-wider transition"
          >
            Ajukan Fasilitas Talangan
          </button>
        </div>
      </section>

      <!-- SECTION 4: DOUBLE-ENTRY GENERAL LEDGER AUDIT TRAIL -->
      <section id="ledger" class="space-y-4">
        <div class="flex items-center justify-between">
          <div>
            <h2 class="text-xl sm:text-2xl font-bold text-white tracking-tight flex items-center">
              <ShieldCheck class="w-6 h-6 mr-2 text-emerald-400" />
              Audit Trail Buku Besar (Double-Entry Ledger)
            </h2>
            <p class="text-xs sm:text-sm text-slate-400">Mutasi keuangan immutable merekam pergerakan escrow, lelang, dan timbangan.</p>
          </div>
          <span class="text-xs px-2.5 py-1 bg-emerald-950 text-emerald-300 border border-emerald-800 rounded-md font-mono">
            Zero-Sum Validated
          </span>
        </div>

        <div class="overflow-x-auto rounded-xl border border-slate-800 shadow-xl">
          <table class="w-full text-left text-xs">
            <thead class="bg-slate-900 text-slate-400 uppercase tracking-wider border-b border-slate-800">
              <tr>
                <th class="px-4 py-3">Jurnal ID</th>
                <th class="px-4 py-3">Tipe Mutasi</th>
                <th class="px-4 py-3">Referensi</th>
                <th class="px-4 py-3">Deskripsi</th>
                <th class="px-4 py-3">Debit</th>
                <th class="px-4 py-3">Credit</th>
                <th class="px-4 py-3 text-right">Nominal (IDR)</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-slate-800/80 bg-slate-950/60 font-mono">
              <tr v-for="j in ledgerJournals" :key="j.id" class="hover:bg-slate-900/50 transition">
                <td class="px-4 py-3 text-emerald-400 font-bold">{{ j.id }}</td>
                <td class="px-4 py-3">
                  <span class="px-2 py-0.5 rounded text-[10px] bg-slate-800 text-slate-300">{{ j.type }}</span>
                </td>
                <td class="px-4 py-3 text-slate-400">{{ j.ref }}</td>
                <td class="px-4 py-3 text-slate-300 font-sans">{{ j.desc }}</td>
                <td class="px-4 py-3 text-rose-300">{{ j.debitAcc }}</td>
                <td class="px-4 py-3 text-emerald-300">{{ j.creditAcc }}</td>
                <td class="px-4 py-3 text-right font-bold text-white">{{ formatRupiah(j.amount) }}</td>
              </tr>
            </tbody>
          </table>
        </div>
      </section>

    </main>

    <!-- MOBILE STICKY BOTTOM NAVIGATION BAR -->
    <nav class="md:hidden fixed bottom-0 left-0 right-0 bg-slate-900/95 backdrop-blur-md border-t border-slate-800 py-2 px-4 z-40 flex items-center justify-around text-[10px]">
      <a href="#tenders" class="flex flex-col items-center text-slate-400 hover:text-emerald-400 transition">
        <Building2 class="w-5 h-5 mb-0.5" />
        <span>Tender</span>
      </a>
      <a href="#weighbridge" class="flex flex-col items-center text-slate-400 hover:text-emerald-400 transition">
        <Scale class="w-5 h-5 mb-0.5" />
        <span>Timbangan</span>
      </a>
      <a href="#financing" class="flex flex-col items-center text-slate-400 hover:text-emerald-400 transition">
        <Coins class="w-5 h-5 mb-0.5" />
        <span>SCF Talangan</span>
      </a>
      <a href="#ledger" class="flex flex-col items-center text-slate-400 hover:text-emerald-400 transition">
        <ShieldCheck class="w-5 h-5 mb-0.5" />
        <span>Audit</span>
      </a>
    </nav>

    <!-- BID MODAL -->
    <div
      v-if="bidModal.isOpen"
      class="fixed inset-0 bg-slate-950/80 backdrop-blur-sm z-50 flex items-center justify-center p-4"
    >
      <div class="bg-slate-900 border border-slate-800 rounded-2xl p-5 sm:p-6 max-w-md w-full space-y-4 shadow-2xl">
        <div class="flex justify-between items-start">
          <div>
            <h3 class="text-base sm:text-lg font-bold text-white">Kunci Penawaran & Bid-Bond</h3>
            <p class="text-xs text-slate-400">{{ bidModal.tender?.title }}</p>
          </div>
          <button @click="bidModal.isOpen = false" class="text-slate-500 hover:text-white">✕</button>
        </div>

        <div class="p-3 bg-amber-950/30 border border-amber-900/50 rounded-xl text-xs text-amber-300 space-y-1">
          <p class="font-bold">⚠️ Persyaratan Jaminan Escrow:</p>
          <p>Dengan menawar, jaminan bid-bond sebesar <strong>{{ formatRupiah(bidModal.tender?.bidBond) }}</strong> akan otomatis dikunci di rekening penampungan.</p>
        </div>

        <div>
          <label class="block text-xs text-slate-400 mb-1">Tawaran Harga per Kg (Min. Rp {{ (bidModal.tender?.currentBid + 1).toLocaleString() }})</label>
          <input
            v-model.number="bidModal.bidPrice"
            type="number"
            class="w-full bg-slate-950 border border-slate-700 rounded-xl px-3 py-2 text-sm text-white focus:border-emerald-500 focus:outline-none"
          />
        </div>

        <div class="flex space-x-3 pt-2">
          <button
            @click="bidModal.isOpen = false"
            class="w-1/2 py-2.5 bg-slate-800 hover:bg-slate-700 text-slate-300 rounded-xl text-xs font-semibold"
          >
            Batal
          </button>
          <button
            @click="confirmBid"
            class="w-1/2 py-2.5 bg-emerald-500 hover:bg-emerald-400 text-slate-950 rounded-xl text-xs font-bold shadow-lg shadow-emerald-500/20"
          >
            Konfirmasi Bidding
          </button>
        </div>
      </div>
    </div>

    <!-- QR CODE SURAT JALAN MODAL -->
    <div
      v-if="qrModal.isOpen"
      class="fixed inset-0 bg-slate-950/80 backdrop-blur-sm z-50 flex items-center justify-center p-4"
    >
      <div class="bg-slate-900 border border-slate-800 rounded-2xl p-6 max-w-sm w-full space-y-4 shadow-2xl text-center">
        <div class="flex justify-between items-center">
          <h3 class="text-sm font-bold text-white">Surat Jalan Digital & QR Timbang</h3>
          <button @click="qrModal.isOpen = false" class="text-slate-500 hover:text-white">✕</button>
        </div>

        <!-- Simulated QR Pattern -->
        <div class="w-44 h-44 mx-auto bg-white p-3 rounded-xl flex items-center justify-center shadow-lg">
          <div class="w-full h-full border-4 border-slate-950 flex flex-col items-center justify-center font-mono text-[9px] text-slate-950 font-bold p-1 leading-tight">
            <span>[SCRAPFLOW-SPK]</span>
            <span class="mt-2">{{ qrModal.tender?.number }}</span>
            <span class="mt-1">TRUCK-PASS-VALID</span>
            <div class="mt-2 w-16 h-1 bg-slate-950"></div>
          </div>
        </div>

        <div class="text-xs space-y-1">
          <p class="font-bold text-emerald-400">{{ qrModal.tender?.company }}</p>
          <p class="text-slate-400 text-[11px]">Tunjukkan QR ini ke pos satpam gerbang timbangan pabrik.</p>
        </div>

        <button
          @click="qrModal.isOpen = false"
          class="w-full py-2 bg-slate-800 hover:bg-slate-700 text-slate-300 rounded-xl text-xs font-semibold"
        >
          Tutup
        </button>
      </div>
    </div>

    <Footer />
  </div>
</template>
