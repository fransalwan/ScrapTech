import "./assets/main.css";

import { createApp } from "vue";
import { createPinia } from "pinia";
import VueGtag from "vue-gtag";

import App from "./App.vue";
import router from "./router";
import vue3GoogleLogin from "vue3-google-login";

const app = createApp(App);
app.use(vue3GoogleLogin, {
  clientId:
    "168733778258-3u3d7t3qu4nh8k24c9j4oqv15m7nco4f.apps.googleusercontent.com",
});

app.use(createPinia());
app.use(router);
app.use(
  VueGtag,
  {
    appName: "ScrapTech",
    pageTrackerScreenviewEnabled: true,
    config: { id: "G-YV5WBN14YC" },
  },
  router
);

app.mount("#app");
