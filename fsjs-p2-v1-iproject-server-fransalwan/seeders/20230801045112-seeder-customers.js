"use strict";

const data = require("../db/customers.json");
const { hashPassword } = require("../helpers/bcrypt");

data.forEach((el) => {
  delete el.id;
  el.password = hashPassword(el.password);
  el.createdAt = new Date();
  el.updatedAt = new Date();
});

/** @type {import('sequelize-cli').Migration} */
module.exports = {
  async up(queryInterface, Sequelize) {
    await queryInterface.bulkInsert("Customers", data);
    /**
     * Add seed commands here.
     *
     * Example:
     * await queryInterface.bulkInsert('People', [{
     *   name: 'John Doe',
     *   isBetaMember: false
     * }], {});
     */
  },

  async down(queryInterface, Sequelize) {
    await queryInterface.bulkInsert("Customers");

    /**
     * Add commands to revert seed here.
     *
     * Example:
     * await queryInterface.bulkDelete('People', null, {});
     */
  },
};
