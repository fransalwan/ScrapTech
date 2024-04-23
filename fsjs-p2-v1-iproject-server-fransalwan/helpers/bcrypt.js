const bcrypt = require("bcryptjs");
const salt = bcrypt.genSaltSync(10);

const hashPassword = (plainPassword) => {
  return bcrypt.hashSync(plainPassword, salt);
};

const comparePassword = (plainPasswod, hashPassword) => {
  return bcrypt.compareSync(plainPasswod, hashPassword);
};

module.exports = { hashPassword, comparePassword };
