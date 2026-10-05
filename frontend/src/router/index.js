import { createRouter, createWebHistory } from "vue-router";
import RegisterPage from "../views/RegisterPage.vue";
import LoginPage from "../views/LoginPage.vue";
import HomePage from "../views/HomePage.vue";
import MyOrderPage from "../views/MyOrderPage.vue";

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: "/",
      name: "home",
      component: HomePage,
    },
    {
      path: "/login",
      name: "login",
      component: LoginPage,
    },
    {
      path: "/register",
      name: "register",
      component: RegisterPage,
    },
    {
      path: "/myorder",
      name: "myorder",
      component: MyOrderPage,
    },
  ],
});

router.beforeEach((to, from, next) => {
  if (to.name === "login" && localStorage.getItem("access_token")) {
    next("/");
  } else if (to.name === "myorder" && !localStorage.getItem("access_token")) {
    next("/");
  } else {
    next();
  }
});

export default router;
