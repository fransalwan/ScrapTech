<script setup>
import { RouterView } from "vue-router";
import { Toaster } from "vue-sonner";
import { onMounted } from "vue";
import { useMainStore } from "./stores/mainStore";

const store = useMainStore();

onMounted(() => {
  store.initTheme();
});
</script>

<template>
  <Toaster richColors position="top-right" :theme="store.theme" closeButton />
  <RouterView />
</template>
