/*
 * Copyright (C) 2024- Germano Rizzo
 *
 * This file is part of FoodHubber.
 *
 * FoodHubber is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * FoodHubber is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with FoodHubber.  If not, see <http://www.gnu.org/licenses/>.
 */
package get_stock

import (
	"foodhubber/params"
	"foodhubber/utils"

	"github.com/gofiber/fiber/v2"
)

type stock struct {
	Item  int `json:"item"`
	Stock int `json:"stock"`
}

type response struct {
	Stock []stock `json:"stock"`
}

func GetStock(c *fiber.Ctx) error {
	ret := response{Stock: make([]stock, 0)}

	params.RWLock.RLock()
	defer params.RWLock.RUnlock()

	sql := "SELECT item_id, quantity FROM stock"
	rows, err := params.Db.Query(sql)
	if err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "FHE001", "stock", &err)
	}
	defer rows.Close()
	for rows.Next() {
		var stock stock
		err = rows.Scan(&stock.Item, &stock.Stock)
		if err != nil {
			return utils.SendError(c, fiber.StatusInternalServerError, "FHE001", "stock", &err)
		}
		ret.Stock = append(ret.Stock, stock)
	}
	if err = rows.Err(); err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "FHE004", "stock", &err)
	}

	c.JSON(ret)
	return c.SendStatus(fiber.StatusOK)
}
