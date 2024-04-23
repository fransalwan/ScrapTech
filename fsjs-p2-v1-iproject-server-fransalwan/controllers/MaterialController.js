const { Op } = require("sequelize");
const { pagination } = require("../middlewares/pagination");
const { Material, Sequelize } = require("../models");

class MaterialController {
  static async getAllMaterials(req, res, next) {
    try {
      const where = {};

      const { title } = req.query;

      // pagination
      const page = req.query.page ? parseInt(req.query.page) : 1; // kalo tidak ada page atau perpage di kemabalikan 1 agar tidak ada error
      const per_page = req.query.per_page ? parseInt(req.query.per_page) : 1;

      if (page <= 0 || per_page <= 0)
        throw { name: "page minimal 1 or per_page minimal 1" };

      if (title) where.title = { [Sequelize.Op.iLike]: `%${title}%` };

      const { count, rows } = await Material.findAndCountAll({
        where,
        offset: (page - 1) * per_page,
        limit: per_page,
        distinct: true,
        order: [["status", "ASC"]],
      });

      // count mengembalikan jumlah data

      // console.log(count, "<<<");
      // console.log(rows, "<<<");

      const result = pagination({
        data: rows,
        count,
        page,
        per_page,
      });

      console.log(count, "<<<< count");

      // Penerapan infinity Scrool disini
      if (count <= 0) {
        res.status(200).json({
          message: "DataNotFound",
        });
      } else {
        res.status(200).json(result);
      }
    } catch (err) {
      console.log(err);
    }
  }
}

module.exports = MaterialController;
