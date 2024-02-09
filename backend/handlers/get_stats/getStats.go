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
package get_stats

import (
	"foodhubber/params"
	"foodhubber/utils"

	"github.com/gofiber/fiber/v2"
)

type checkout struct {
	Checkout string `json:"checkout"`
	Orders   int    `json:"orders"`
	Items    int    `json:"items"`
}

type byItem struct {
	Item int `json:"item"`
	Qty  int `json:"qty"`
}

type response struct {
	From       string     `json:"from"`
	To         string     `json:"to"`
	Orders     int        `json:"orders"`
	Items      int        `json:"items"`
	ByCheckout []checkout `json:"byCheckout"`
	ByItem     []byItem   `json:"byItem"`
}

func fmt(in string) string {
	return in[0:4] + "-" + in[4:6] + "-" + in[6:8]
}

func GetStats(c *fiber.Ctx) error {
	from := fmt(c.Query("from")) + " 00:00:00"
	to := fmt(c.Query("to")) + " 23:59:59"

	ret := response{
		ByCheckout: make([]checkout, 0),
		ByItem:     make([]byItem, 0),
	}

	params.RWLock.RLock()
	defer params.RWLock.RUnlock()

	query := `
			SELECT COUNT(1) AS orders
			  FROM orders
			 WHERE active = 1
			   AND datetime BETWEEN $1 AND $2`
	row := params.Db.QueryRow(query, from, to)
	if err := row.Scan(&ret.Orders); err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "FHE001", "orders", &err)
	}

	query = `
		    SELECT COALESCE(SUM(orw.quantity), 0) as quantity
		      FROM order_rows orw
		      JOIN orders o ON orw.order_id = o.id
	         WHERE o.active = 1
		       AND o.datetime BETWEEN $1 AND $2`
	row = params.Db.QueryRow(query, from, to)
	if err := row.Scan(&ret.Items); err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "FHE001", "orders", &err)
	}

	query = `
			SELECT o.checkout_id, COUNT(DISTINCT o.id) as orders,
		           SUM(orw.quantity) as quantity
              FROM order_rows orw
              JOIN orders o ON orw.order_id = o.id
             WHERE o.active = 1
			   AND o.datetime BETWEEN $1 AND $2
             GROUP BY o.checkout_id`
	rows, err := params.Db.Query(query, from, to)
	if err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "FHE001", "orders", &err)
	}
	defer rows.Close()
	for rows.Next() {
		var checkout checkout
		err = rows.Scan(&checkout.Checkout, &checkout.Orders, &checkout.Items)
		if err != nil {
			return utils.SendError(c, fiber.StatusInternalServerError, "FHE001", "orders", &err)
		}
		ret.ByCheckout = append(ret.ByCheckout, checkout)
	}
	if err = rows.Err(); err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "FHE004", "orders", &err)
	}

	query = `
		    SELECT orw.item_id, SUM(orw.quantity) as quantity
		      FROM order_rows orw
		      JOIN orders o ON orw.order_id = o.id
	         WHERE o.active
		       AND o.datetime BETWEEN $1 AND $2
	         GROUP BY orw.item_id`
	rows, err = params.Db.Query(query, from, to)
	if err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "FHE001", "orders", &err)
	}
	defer rows.Close()
	for rows.Next() {
		var items byItem
		err = rows.Scan(&items.Item, &items.Qty)
		if err != nil {
			return utils.SendError(c, fiber.StatusInternalServerError, "FHE001", "orders", &err)
		}
		ret.ByItem = append(ret.ByItem, items)
	}
	if err = rows.Err(); err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "FHE004", "orders", &err)
	}

	c.JSON(ret)
	return c.SendStatus(fiber.StatusOK)
}
