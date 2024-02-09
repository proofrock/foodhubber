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
package set_stock

import (
	"context"
	"database/sql"
	"foodhubber/params"
	"foodhubber/utils"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type request struct {
	Item  int     `json:"item"`
	Stock *string `json:"stock"`
}

type response struct {
}

func SetStock(c *fiber.Ctx) error {
	req := new(request)
	if err := c.BodyParser(req); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "FHE006", "body", &err)
	}

	params.RWLock.Lock()
	defer params.RWLock.Unlock()

	tx, err := params.Db.BeginTx(context.Background(), nil)
	if err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "FHE007", "", &err)
	}
	defer tx.Rollback()

	if req.Stock == nil {
		query := "DELETE FROM stock WHERE item_id = $1"
		if _, err = tx.Exec(query, req.Item); err != nil {
			return utils.SendError(c, fiber.StatusInternalServerError, "FHE002", "stock", &err)
		}
	} else {
		query := "SELECT quantity FROM stock WHERE item_id = $1"
		row := tx.QueryRow(query, req.Item)

		var finalStock int

		stock := *req.Stock

		if stock[0] == '-' || stock[0] == '+' {
			variation, err := strconv.Atoi(stock[1:])
			if err != nil {
				return utils.SendError(c, fiber.StatusBadRequest, "FHE005", "stock", nil)
			}

			var originalStock int

			if err := row.Scan(&originalStock); err != nil && err != sql.ErrNoRows {
				return utils.SendError(c, fiber.StatusInternalServerError, "FHE001", "stock", &err)
			}

			if stock[0] == '-' {
				finalStock = originalStock - variation
			} else {
				finalStock = originalStock + variation
			}
		} else {
			quantity, err := strconv.Atoi(stock)
			if err != nil {
				return utils.SendError(c, fiber.StatusBadRequest, "FHE005", "stock", nil)
			}
			finalStock = quantity
		}

		query = `INSERT INTO stock (item_id, quantity)
				           VALUES ($1, $2)
			          ON CONFLICT (item_id) DO
				           UPDATE SET quantity = EXCLUDED.quantity`
		if _, err = tx.Exec(query, req.Item, finalStock); err != nil {
			return utils.SendError(c, fiber.StatusInternalServerError, "FHE002", "stock", &err)
		}
	}

	if err := tx.Commit(); err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "FHE008", "", &err)
	}

	c.JSON(response{})
	return c.SendStatus(fiber.StatusOK)
}
