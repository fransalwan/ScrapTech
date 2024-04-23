require("dotenv").config();

const express = require("express");
const app = express();
const port = process.env.PORT || 3000;
const cors = require("cors");
const MaterialController = require("./controllers/MaterialController");
const CategoryController = require("./controllers/CategoryController");
const authentication = require("./middlewares/authentication");
const PaymentController = require("./controllers/PaymentController");
const CustomerController = require("./controllers/CustomerController");

app.use(cors());
app.use(express.urlencoded({ extended: false }));
app.use(express.json());

app.post("/register", CustomerController.register);
app.post("/login", CustomerController.login);
app.post("/OAuthLogin", CustomerController.OAuthLogin);
app.get("/materials", MaterialController.getAllMaterials);
app.get("/categories", CategoryController.getAllCategories);

app.use(authentication);

app.get("/order", PaymentController.getAllMyOrder);
app.post("/order/:id", PaymentController.addOrderById);
app.delete("/order/:id", PaymentController.deleteOrderById);
app.patch("/payment", PaymentController.editStatusOrder);
app.post("/generate-midtrans-token", PaymentController.generateToken);

app.listen(port, () => {
  console.log(`Example app listening on port ${port}`);
});
