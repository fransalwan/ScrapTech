const { comparePassword } = require("../helpers/bcrypt");
const { signToken } = require("../helpers/jwt");
const { Customer } = require("../models/index");
const { OAuth2Client } = require("google-auth-library");
const client = new OAuth2Client();

class CustomerController {
  static async register(req, res) {
    try {
      const { username, email, password } = req.body;

      const created = await Customer.create({
        username,
        email,
        password,
      });

      res.status(201).json({
        id: created.id,
        email: created.email,
      });
    } catch (err) {
      console.log(err);
    }
  }

  static async login(req, res, next) {
    try {
      const { email, password } = req.body;

      console.log(email, password);

      if (!email || !password) {
        throw { name: "loginError" };
      }

      const customer = await Customer.findOne({
        where: {
          email,
        },
      });

      if (!customer) {
        throw { name: "loginError" };
      }

      if (!comparePassword(password, customer.password)) {
        throw { name: "loginError" };
      }

      const payload = {
        id: customer.id,
        email: customer.email,
        role: customer.role,
      };

      const token = signToken(payload);

      res.status(200).json({
        access_token: token,
      });
    } catch (err) {
      console.log(err);
    }
  }

  static async OAuthLogin(req, res, next) {
    try {
      // console.log("fungsi ke panggil");
      // console.log(req.headers.google_token, "<<<<<<<<<<");

      const ticket = await client.verifyIdToken({
        idToken: req.headers.google_token,
        audience: process.env.GOOGLE_CLIENT_ID,
      });

      // console.log(ticket, "<<<<<<<<<<<");

      const payload = ticket.getPayload();
      // console.log(payload, "<<<<<<< payload");

      // console.log("dibawah harunya ada user <<<<<<<");
      const [customer, isCreated] = await Customer.findOrCreate({
        where: {
          email: payload.email,
        },
        defaults: {
          username: payload.name,
          role: "customer",
          password: "google oauth",
        },
        hooks: false, // agar tidak dihash
      });

      // console.log(customer.dataValues.username, "<<<<< ini customer google");

      const access_token = signToken({
        id: customer.id,
        email: customer.email,
      });
      // console.log(access_token, "acces token bos");
      // console.log(customer.username, "<<< ini customer");
      res.status(200).json({
        username: customer.username,
        access_token,
      });
    } catch (err) {
      next(err);
    }
  }
}

module.exports = CustomerController;
