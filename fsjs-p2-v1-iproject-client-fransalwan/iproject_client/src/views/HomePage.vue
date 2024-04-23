<script>
import Navbar from "../components/Navbar.vue";
import CardMaterial from "../components/CardMaterial.vue";
import Footer from "../components/Footer.vue";
import NotFound from "../components/NotFound.vue";
import { mapActions, mapState } from "pinia";
import { useMainStore } from "../stores/mainStore";

export default {
  name: "HomePage",
  components: { Navbar, CardMaterial, Footer, NotFound },
  data() {
    return {
      title: "",
    };
  },
  methods: {
    ...mapActions(useMainStore, ["fetchAllData", "pagingHandler", "doFilter"]),
  },
  created() {
    this.fetchAllData();
  },
  computed: {
    ...mapState(useMainStore, ["materials", "isNotFound", "isLogin"]),
  },
};
</script>

<template>
  <Navbar />

  <div class="container px-4 mx-auto mt-8 max-w-7xl">
    <!-- Search -->
    <div class="flex justify-start mt-4">
      <div class="mr-2 text-lg font-medium">Search:</div>
      <div class="flex items-center">
        <input
          @keyup.enter="doFilter(title)"
          type="text"
          placeholder="Search"
          class="px-2 py-1 border border-gray-300 rounded-md focus:outline-none focus:ring focus:border-blue-500"
          v-model="title"
        />
      </div>
    </div>
    <p v-if="!isLogin" class="text-rose-400">*Silahkan login untuk order</p>

    <!-- Card Section -->
    <div class="grid grid-cols-4 gap-4 mt-2">
      <!-- Card 1 -->
      <CardMaterial :isLogin="isLogin" :materials="materials" />
    </div>
    <NotFound v-if="isNotFound" />

    <!-- Pagination -->
    <div v-if="!isNotFound" class="flex justify-center mt-8 mb-8">
      <div
        v-for="pageNumber in materials.pagination.total_page"
        :key="pageNumber"
      >
        <a
          v-if="pageNumber === materials.pagination.current_page"
          class="px-4 py-2 mr-2 text-white bg-orange-500 rounded-md"
          >{{ pageNumber }}</a
        >
        <a
          @click="pagingHandler(pageNumber)"
          v-if="pageNumber !== materials.pagination.current_page"
          class="px-4 py-2 mr-2 text-white rounded-md bg-slate-500"
          >{{ pageNumber }}</a
        >
      </div>
    </div>
  </div>

  <!-- Footer -->
  <Footer />
</template>
