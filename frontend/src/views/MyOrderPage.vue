<script>
import { mapActions, mapState } from "pinia";
import { useMainStore } from "../stores/mainStore";
import Navbar from "../components/Navbar.vue";

export default {
  name: "MyOrderPage",
  components: { Navbar },
  methods: {
    ...mapActions(useMainStore, [
      "handleMyOrder",
      "handleOrder",
      "handleDeleteOrder",
    ]),
  },
  created() {
    this.handleMyOrder();
  },
  computed: {
    ...mapState(useMainStore, ["myorder"]),
  },
};
</script>

<template>
  <Navbar />
  <div class="py-10 px-28">
    <h1 class="mb-4 text-2xl font-bold">My Order</h1>

    <!-- Cart items list -->
    <div class="bg-white shadow sm:rounded-md">
      <ul class="divide-y divide-gray-200">
        <!-- Cart item 1 -->
        <li v-for="item in myorder.result" class="px-4 py-4 sm:px-6">
          <div class="flex items-center justify-between">
            <div class="flex items-center">
              <img
                class="object-cover w-16 h-16 rounded-full"
                :src="item.Material.image"
                alt="Product 1"
              />
              <div class="ml-4">
                <h2 class="text-lg font-medium text-gray-900">
                  {{ item.Material.title }}
                </h2>
                <p class="text-gray-500">Rp: {{ item.Material.price }}</p>
              </div>
            </div>
            <button
              @click="handleDeleteOrder(item.Material.id)"
              class="px-4 py-2 text-white bg-red-500 rounded hover:bg-red-600"
            >
              Remove
            </button>
          </div>
        </li>
      </ul>
    </div>

    <!-- Cart total -->
    <div class="mt-4 text-right">
      <p class="text-lg font-medium">Total: Rp. {{ myorder.total }}</p>
    </div>

    <!-- Checkout button -->
    <div class="mt-4 text-right">
      <button
        @click="handleOrder(myorder.total)"
        class="px-4 py-2 text-white bg-orange-400 rounded hover:bg-orange-600"
      >
        Checkout
      </button>
    </div>
  </div>

  <Footer />
</template>
