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
package del_order

import (
	"context"
	"foodhubber/db_ops"
	"foodhubber/params"
	"foodhubber/utils"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type response struct {
}

func DelOrder(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "FHE005", "id", &err)
	}

	defer func() { go db_ops.Backup() }()
	params.RWLock.Lock()
	defer params.RWLock.Unlock()

	tx, err := params.Db.BeginTx(context.Background(), nil)
	if err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "FHE007", "", &err)
	}
	defer tx.Rollback()

	query := "UPDATE orders SET active = 0 WHERE id = $1"
	if ret, err := tx.Exec(query, id); err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "FHE002", "orders", &err)
	} else {
		if ra, err := ret.RowsAffected(); err != nil {
			return utils.SendError(c, fiber.StatusInternalServerError, "FHE002", "orders", &err)
		} else if ra == 0 {
			return utils.SendError(c, fiber.StatusNotFound, "FHE102", "", nil)
		}
	}

	query = `
		UPDATE stock
		   SET quantity = stock.quantity + orw.quantity
		  FROM order_rows orw
		 WHERE stock.item_id = orw.item_id 
		   AND orw.order_id = $1`
	if _, err = tx.Exec(query, id); err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "FHE002", "stock", &err)
	}

	params.TouchOrdersGen()
	params.TouchStockGen()

	if err := tx.Commit(); err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "FHE008", "", &err)
	}

	c.JSON(response{})
	return c.SendStatus(fiber.StatusOK)
}
