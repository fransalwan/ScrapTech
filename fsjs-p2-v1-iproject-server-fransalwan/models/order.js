"use strict";
const { Model } = require("sequelize");
module.exports = (sequelize, DataTypes) => {
  class Order extends Model {
    /**
     * Helper method for defining associations.
     * This method is not a part of Sequelize lifecycle.
     * The `models/index` file will call this method automatically.
     */
    static associate(models) {
      // define association here
      Order.belongsTo(models.Customer, { foreignKey: "CustomerId" });
      Order.belongsTo(models.Material, { foreignKey: "MaterialId" });
    }
  }
  Order.init(
    {
      CustomerId: {
        type: DataTypes.INTEGER,
        allowNull: false,
        validate: {
          notEmpty: {
            msg: "Customer id is required",
          },
          notNull: {
            msg: "Customer id is required",
          },
        },
      },
      MaterialId: {
        type: DataTypes.INTEGER,
        allowNull: false,
        validate: {
          notEmpty: {
            msg: "Material id is required",
          },
          notNull: {
            msg: "Material id is required",
          },
        },
      },
      status: {
        type: DataTypes.STRING,
        allowNull: false,
        validate: {
          notEmpty: {
            msg: "Status is required",
          },
          notNull: {
            msg: "Status is required",
          },
        },
      },
    },
    {
      sequelize,
      modelName: "Order",
    }
  );
  return Order;
};
