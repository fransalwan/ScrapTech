const { Customer } = require("../models");
const { verifyToken } = require("../helpers/jwt");

const authentication = async (req, res, next) => {
  try {
    const { access_token } = req.headers;
    // console.log("authentication ni bos");
    if (!access_token) throw { name: "jwtNotFound" };

    const payload = verifyToken(access_token);

    // console.log(payload, "<<<<<<<<");

    if (!payload) throw { name: "JsonWebTokenError" };

    const customer = await Customer.findOne({
      where: {
        id: payload.id,
      },
    });

    // console.log(customer);

    if (!customer) throw { name: "NotFound" };

    req.additionalData = {
      customerId: customer.id,
      email: customer.email,
      role: customer.role,
    };
    next();
  } catch (err) {
    console.log(err);
    next(err);
  }
};

module.exports = authentication;
