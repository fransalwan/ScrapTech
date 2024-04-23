import { ref, computed } from "vue";
import { defineStore } from "pinia";
import axios from "axios";
import router from "../router";
import { event } from "vue-gtag";

const baseUrl = "http://localhost:3000";

export const useMainStore = defineStore("main", {
  state: () => ({
    isLogin: localStorage.getItem("access_token") ? true : false,
    isNotFound: false,
    myorder: [],
    materials: {
      data: [],
      pagination: {
        total_page: 0,
        per_page: 8,
      },
    },
  }),
  actions: {
    // since we rely on `this`, we cannot use an arrow function
    async handleRegister(username, email, password) {
      try {
        const { data } = await axios({
          method: "POST",
          url: `${baseUrl}/register`,
          data: {
            username,
            email,
            password,
          },
        });
        router.push("/login");
      } catch (err) {
        console.log(err);
      }
    },
    async handleLogin(email, password) {
      try {
        // console.log("<<<<fungsi nya terpanggil>>>>");
        const { data } = await axios({
          method: "POST",
          url: `${baseUrl}/login`,
          data: {
            email,
            password,
          },
        });
        event("user-login", {
          "event-label": "user-login",
          value: email,
        });

        localStorage.setItem("access_token", data.access_token);
        this.isLogin = true;
        router.push("/");
      } catch (err) {
        console.log(err);
      }
    },
    async handleOauth(res) {
      try {
        // console.log("<<< udah masuk method Oauth");
        // console.log(res.credential, "<<<<<ini resnya di method>>>>>");
        const { data } = await axios({
          method: "POST",
          url: `${baseUrl}/OAuthLogin`,
          headers: {
            google_token: res.credential,
          },
        });

        event("user-login", {
          "event-label": "user-login",
          value: data.username,
        });

        router.push("/");
        localStorage.setItem("access_token", data.access_token);
      } catch (err) {
        console.log(err);
      }
    },
    async fetchAllData(page, filter) {
      try {
        const { data } = await axios({
          method: "GET",
          url: `${baseUrl}/materials`,
          headers: {
            access_token: localStorage.getItem("access_token"),
          },
          params: {
            per_page: 8,
            page,
            title: filter,
          },
        });

        console.log(data, "<<< ini balikannya");
        if (data.message === "DataNotFound") {
          this.isNotFound = true;
        } else {
          this.isNotFound = false;
        }

        this.materials = data;
      } catch (err) {
        console.log(err);
      }
    },
    async handleDeleteOrder(id) {
      try {
        const { data } = await axios({
          method: "DELETE",
          url: `${baseUrl}/order/${id}`,
          headers: {
            access_token: localStorage.getItem("access_token"),
          },
        });
        this.handleMyOrder();
      } catch (err) {
        console.log(err);
      }
    },
    handleLogout() {
      console.log("method logout called");
      localStorage.clear();
      this.isLogin = false;
      router.push("/");
    },
    doFilter(filter) {
      // console.log(filter, "<<< ini filter nya");
      this.fetchAllData(1, filter);
    },
    pagingHandler(pageNumber) {
      // console.log(pageNumber, "<<<");
      this.fetchAllData(pageNumber);
    },
    async handleChangeStatus() {
      try {
        const { data } = await axios({
          method: "PATCH",
          url: `${baseUrl}/payment`,
          headers: {
            access_token: localStorage.getItem("access_token"),
          },
        });

        this.fetchAllData(1, {});
      } catch (err) {
        console.log(err);
      }
    },
    async handleOrder(total) {
      try {
        // console.log(total);
        const { data } = await axios({
          method: "POST",
          url: `${baseUrl}/generate-midtrans-token`,
          headers: {
            access_token: localStorage.getItem("access_token"),
          },
          data: {
            total,
          },
        });
        // console.log(data, "<<<<<");
        const cb = this.handleChangeStatus;

        window.snap.pay(data.token, {
          onSuccess: function (result) {
            /* You may add your own implementation here */
            console.log("payment success!");
            cb();
          },
        });
      } catch (err) {
        console.log(err);
      }
    },
    async handleMyOrder() {
      try {
        console.log("<<< method has been called");
        const { data } = await axios({
          method: "GET",
          url: `${baseUrl}/order`,
          headers: {
            access_token: localStorage.getItem("access_token"),
          },
        });

        console.log(data, "<<< ini data");

        this.myorder = data;
      } catch (err) {
        console.log(err);
      }
    },
    async handleAddOrder(id) {
      try {
        const { data } = await axios({
          method: "POST",
          url: `${baseUrl}/order/${id}`,
          headers: {
            access_token: localStorage.getItem("access_token"),
          },
        });

        event("buy-product", {
          "event-label": "user-buy-product",
          value: id,
        });

        router.push("/myorder");
      } catch (err) {
        console.log(err);
      }
    },
  },
});
