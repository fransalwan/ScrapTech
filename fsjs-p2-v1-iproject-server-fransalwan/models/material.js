"use strict";
const { Model } = require("sequelize");
module.exports = (sequelize, DataTypes) => {
  class Material extends Model {
    /**
     * Helper method for defining associations.
     * This method is not a part of Sequelize lifecycle.
     * The `models/index` file will call this method automatically.
     */
    static associate(models) {
      // define association here
      Material.hasMany(models.Category, { foreignKey: "CategoryId" });
      Material.hasMany(models.Order, { foreignKey: "MaterialId" });
    }
  }
  Material.init(
    {
      CategoryId: {
        type: DataTypes.INTEGER,
        allowNull: false,
        validate: {
          notEmpty: {
            msg: "Category id is required",
          },
          notNull: {
            msg: "Category id is required",
          },
        },
      },
      title: {
        type: DataTypes.STRING,
        allowNull: false,
        validate: {
          notEmpty: {
            msg: "Title is required",
          },
          notNull: {
            msg: "Title is required",
          },
        },
      },
      image: {
        type: DataTypes.STRING,
        allowNull: false,
        validate: {
          notEmpty: {
            msg: "Image is required",
          },
          notNull: {
            msg: "Image is required",
          },
        },
      },
      spec: {
        type: DataTypes.STRING,
        allowNull: false,
        validate: {
          notEmpty: {
            msg: "Specification is required",
          },
          notNull: {
            msg: "Specification is required",
          },
        },
      },
      price: {
        type: DataTypes.STRING,
        allowNull: false,
        validate: {
          notEmpty: {
            msg: "price is required",
          },
          notNull: {
            msg: "price is required",
          },
        },
      },
      status: {
        type: DataTypes.STRING,
        allowNull: false,
        defaultValue: "ready",
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
      modelName: "Material",
    }
  );
  return Material;
};
