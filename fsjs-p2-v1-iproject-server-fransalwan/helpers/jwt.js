const jwt = require("jsonwebtoken");
const SECRET_JWT = process.env.JWT_SECRET;

const signToken = (payload) => {
  return jwt.sign(payload, SECRET_JWT);
};

const verifyToken = (token) => {
  return jwt.verify(token, SECRET_JWT);
};

module.exports = { signToken, verifyToken };
