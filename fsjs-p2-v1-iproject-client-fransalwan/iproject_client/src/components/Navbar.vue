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
  <!-- {{ isLogin }} -->
  <nav class="bg-gray-800">
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
      <div class="flex items-center justify-between h-16">
        <div class="flex items-center">
          <div class="flex-shrink-0">
            <RouterLink
              to="/"
              class="text-gray-300 hover:bg-gray-700 hover:text-white px-3 py-2 rounded-md text-xl font-medium"
              >ScrapTech</RouterLink
            >
          </div>
          <div class="hidden md:block">
            <div class="ml-10 flex items-baseline space-x-4"></div>
          </div>
        </div>
        <div class="flex items-center ml-4">
          <RouterLink
            v-if="isLogin"
            to="/myorder"
            class="text-gray-300 hover:bg-gray-700 hover:text-white px-3 py-2 rounded-md text-sm font-medium"
            >My Order</RouterLink
          >
          <a
            v-if="isLogin"
            @click="handleLogout()"
            class="text-gray-300 hover:bg-gray-700 hover:text-white px-3 py-2 rounded-md text-sm font-medium"
            >Logout</a
          >
          <RouterLink
            v-if="!isLogin"
            to="login"
            class="px-4 py-2 text-sm font-medium text-gray-800 bg-white rounded-md"
            >Sign In</RouterLink
          >
          <RouterLink
            v-if="!isLogin"
            to="register"
            class="px-4 py-2 text-sm font-medium text-white bg-gray-900 rounded-md"
            >Sign Up</RouterLink
          >
        </div>
      </div>
    </div>
  </nav>
</template>
