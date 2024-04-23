<script>
import { mapActions } from "pinia";
import { useMainStore } from "../stores/mainStore";

export default {
  name: "RegisterPage",
  data: function () {
    return {
      username: "",
      email: "",
      password: "",
    };
  },
  methods: {
    ...mapActions(useMainStore, ["handleRegister"]),
  },
};
</script>

<template>
  <div class="flex items-center justify-center h-screen">
    <div class="w-1/3">
      <h2 class="mb-6 text-2xl text-center">Register</h2>
      <form
        @submit.prevent="handleRegister(username, email, password)"
        class="px-8 pt-6 pb-8 mb-4 bg-white rounded shadow-md"
      >
        <div class="mb-4">
          <label
            class="block mb-2 text-sm font-bold text-gray-700"
            for="username"
          >
            Username
          </label>
          <input
            type="text"
            placeholder="Username"
            required
            class="w-full px-3 py-2 leading-tight text-gray-700 border rounded appearance-none focus:outline-none focus:shadow-outline"
            v-model="username"
          />
        </div>
        <div class="mb-4">
          <label
            class="block mb-2 text-sm font-bold text-gray-700"
            for="username"
          >
            Email
          </label>
          <input
            type="text"
            placeholder="Email"
            required
            class="w-full px-3 py-2 leading-tight text-gray-700 border rounded appearance-none focus:outline-none focus:shadow-outline"
            v-model="email"
          />
        </div>
        <div class="mb-6">
          <label
            class="block mb-2 text-sm font-bold text-gray-700"
            for="password"
          >
            Password
          </label>
          <input
            type="password"
            placeholder="Password"
            required
            class="w-full px-3 py-2 leading-tight text-gray-700 border rounded appearance-none focus:outline-none focus:shadow-outline"
            v-model="password"
          />
        </div>
        <div class="flex items-center justify-center">
          <button
            type="submit"
            class="px-4 py-2 mr-2 font-bold text-white bg-green-500 rounded hover:bg-green-700 focus:outline-none focus:shadow-outline"
          >
            Register
          </button>
          <p>
            Sudah punya akun?
            <RouterLink to="login"
              ><a href="#" class="text-sm text-blue-500">Login!</a></RouterLink
            >
          </p>
        </div>
      </form>
    </div>
  </div>
</template>
