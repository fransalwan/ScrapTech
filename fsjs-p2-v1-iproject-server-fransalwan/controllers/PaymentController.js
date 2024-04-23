const { Customer, Order, Material } = require("../models/index");
const midtransClient = require("midtrans-client");

class PaymentController {
  static async getAllMyOrder(req, res, next) {
    try {
      const result = await Order.findAll({
        include: Material,
        where: {
          CustomerId: req.additionalData.customerId,
          status: "pending",
        },
      });

      let total = 0; // Ubah const menjadi let
      result.forEach((order) => {
        const material = order.Material;
        const price = material.dataValues.price;
        total += price;
      });

      // console.log(total, "<<<< ini totalnya");

      if (!result) throw { name: "Error" };

      res.status(200).json({
        result,
        total,
      });
    } catch (err) {
      console.log(err, "<<<<");
    }
  }
  static async deleteOrderById(req, res, next) {
    try {
      const { id } = req.params;

      // console.log(id, "<<<ini id");

      await Order.destroy({
        where: {
          MaterialId: id,
          status: "pending",
        },
      });

      res.status(201).json({
        message: `order with ${id} has been delete`,
      });
    } catch (err) {
      console.log(err);
      next(err);
    }
  }
  static async addOrderById(req, res, next) {
    try {
      const { id } = req.params;
      // console.log(id, "<<<ini id");
      console.log(req.additionalData.customerId, "<<<ini id");

      const result = await Order.findOrCreate({
        where: {
          MaterialId: id,
          CustomerId: req.additionalData.customerId,
          status: "pending",
        },
        defaults: {
          MaterialId: id,
          CustomerId: req.additionalData.customerId,
          status: "pending",
        },
      });

      // console.log(result, "<<<< ini table conjuction");

      if (!result) throw { name: "Error" };

      console.log(result);

      // jika data tidak ditemukan kembalikan 404

      res.status(201).json({
        message: result,
      });
    } catch (err) {
      console.log(err);
      next(err);
    }
  }

  static async editStatusOrder(req, res, next) {
    try {
      // 1. jalankan fungsi midtrans di sini

      // 2. panggil semua order yang customerId sama dengan yang di authentication dan ubah status order yang sebelumnya pending menjadi paid
      const order = await Order.findAll({
        where: { CustomerId: req.additionalData.customerId },
      });

      order.forEach(async (el) => {
        if (order) {
          console.log("masuk sini");
          await Material.update(
            { status: "sold_out" },
            { where: { id: el.dataValues.MaterialId } }
          );
        } else {
          console.log("Order not found for the specified customer ID.");
        }
      });

      // edit status order
      await Order.update(
        {
          status: "paid",
        },
        {
          where: {
            CustomerId: req.additionalData.customerId,
          },
        }
      );

      res.status(200).json({
        message: "order has been paid",
      });
    } catch (err) {
      console.log(err);
      next(err);
    }
  }

  static async generateToken(req, res, next) {
    try {
      // console.log(req.body, "<<<");
      const { total } = req.body;

      // console.log(total);
      console.log("generateToken has been called");
      const findCustomer = await Customer.findByPk(
        req.additionalData.customerId
      );

      let snap = new midtransClient.Snap({
        // Set to true if you want Production Environment (accept real transaction).
        isProduction: false,
        serverKey: process.env.MID_TRANS_SERVER_KEY,
      });

      let parameter = {
        transaction_details: {
          order_id:
            "TRANSACTION_" + Math.floor(1000000 + Math.random() * 9000000), //harus unique
          gross_amount: total, // bisa tangkap dari client total harga order
        },
        credit_card: {
          secure: true,
        },
        customer_details: {
          //     first_name: "budi",
          //     last_name: "pratama",
          email: findCustomer.email,
          //     phone: "08111222333",
        },
      };

      const midtransToken = await snap.createTransaction(parameter);

      // console.log(midtransToken, "<<<<");
      res.status(201).json(midtransToken);
    } catch (err) {
      console.log(err);
      next(err);
    }
  }
}

module.exports = PaymentController;
