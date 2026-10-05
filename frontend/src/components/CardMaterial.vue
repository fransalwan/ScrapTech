<script>
import { mapActions } from "pinia";
import { useMainStore } from "../stores/mainStore";

export default {
  name: "CardMaterial",
  props: ["materials", "isLogin"],
  methods: {
    ...mapActions(useMainStore, ["handleAddOrder"]),
  },
};
</script>

<template>
  <!-- Card Section -->

  <div v-for="material in materials.data" class="grid justify-between gap-4">
    <!-- Card -->
    <div class="max-w-sm p-4 bg-white rounded-lg shadow-md">
      <img
        class="object-cover w-full h-48 mb-4 rounded-md"
        alt="Card Image"
        :src="material.image"
      />

      <a
        v-if="material.status === 'ready'"
        class="px-4 py-0 text-sm text-green-600 border-2 border-green-600 rounded-full shadow-sm font-poppins focus-visible:outline-orange-600"
      >
        {{ material.status }}
      </a>
      <a
        v-if="material.status === 'sold_out'"
        class="px-4 py-0 text-sm text-red-600 border-2 border-red-600 rounded-full shadow-sm font-poppins focus-visible:outline-orange-600"
      >
        {{ material.status }}
      </a>
      <h2 class="mt-4 text-xl text-orange-400 font-poppins">
        {{ material.title }}
      </h2>

      <p class="mb-2 text-base font-poppins text-slate-700">
        Rp. {{ material.price }}
      </p>
      <!-- <p class="mb-2 text-gray-700">{{ material.spec }}</p> -->

      <button
        @click="handleAddOrder(material.id)"
        v-if="isLogin"
        class="px-1 py-1 mr-2 font-medium text-orange-400 no-underline border-2 border-orange-400 rounded-sm hover:bg-orange-400 hover:text-white"
      >
        Order
      </button>
    </div>
  </div>
</template>
