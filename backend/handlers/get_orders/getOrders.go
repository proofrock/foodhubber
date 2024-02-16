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
package get_orders

import (
	"database/sql"
	"foodhubber/params"
	"foodhubber/utils"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type roww struct {
	Item int `json:"item"`
	Qty  int `json:"qty"`
}

type order struct {
	ID          int     `json:"id"`
	Checkout    string  `json:"checkout"`
	Operator    string  `json:"operator"`
	Beneficiary string  `json:"beneficiary"`
	Note        *string `json:"note"`
	Datetime    string  `json:"datetime"`
	Rows        []roww  `json:"rows"`
}

type response struct {
	NumPages int     `json:"numPages"`
	Orders   []order `json:"orders"`
}

func GetOrders(c *fiber.Ctx) error {
	page, err := strconv.Atoi(c.Query("page", "1"))
	if err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "FHE005", "page", &err)
	}

	filter := c.Query("filter", "")
	// TODO? differentiate btw numeric filters and string ones, optimizations can be put
	// in places (avoid filtering on the numeric fields...)

	ret := response{
		Orders: make([]order, 0),
	}

	params.RWLock.RLock()
	defer params.RWLock.RUnlock()

	var row *sql.Row

	query := `
			SELECT MAX(1, CEIL(COUNT(1) / (SELECT value FROM configs WHERE key = 'order_list_page_size'))) AS pages
			  FROM orders
			 WHERE active = 1`
	if filter != "" {
		query += " AND (id = $1 OR beneficiary_id = $2 OR checkout_id LIKE '%' + $3 + '%' OR operator LIKE '%' + $4 + '%')"
		row = params.Db.QueryRow(query, filter, filter, filter, filter)
	} else {
		row = params.Db.QueryRow(query)
	}
	if err := row.Scan(&ret.NumPages); err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "FHE001", "configs", &err)
	}

	var rows *sql.Rows

	if filter != "" {
		query := `
			SELECT id, checkout_id, operator, beneficiary_id, note, strftime('%Y%m%dT%H%M%S', datetime) AS datetime
		      FROM orders
	         WHERE active = 1
			   AND (id = $1 OR beneficiary_id = $2
			        OR checkout_id LIKE '%' + $3 + '%' OR operator LIKE '%' + $4 + '%')
	         ORDER BY id DESC
             LIMIT (SELECT value FROM configs WHERE key = 'order_list_page_size')
            OFFSET ($5 - 1) * (SELECT value FROM configs WHERE key = 'order_list_page_size')`
		rows, err = params.Db.Query(query, filter, filter, filter, filter, page)
	} else {
		query := `
			SELECT id, checkout_id, operator, beneficiary_id, note, strftime('%Y%m%dT%H%M%S', datetime) AS datetime
		      FROM orders
	         WHERE active = 1
	         ORDER BY id DESC
             LIMIT (SELECT value FROM configs WHERE key = 'order_list_page_size')
            OFFSET ($1 - 1) * (SELECT value FROM configs WHERE key = 'order_list_page_size')`
		rows, err = params.Db.Query(query, page)
	}
	if err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "FHE001", "orders [2]", &err)
	}
	defer rows.Close()

	stmt, err := params.Db.Prepare("SELECT item_id, quantity FROM order_rows WHERE order_id = $1")
	if err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "FHE013", "order_rows", &err)
	}
	defer stmt.Close()

	for rows.Next() {
		var order order
		err = rows.Scan(&order.ID, &order.Checkout, &order.Operator, &order.Beneficiary, &order.Note, &order.Datetime)
		if err != nil {
			return utils.SendError(c, fiber.StatusInternalServerError, "FHE001", "orders [2]", &err)
		}

		// TODO? don't preload all rows but do it on demand
		rows2, err := stmt.Query(order.ID)
		if err != nil {
			return utils.SendError(c, fiber.StatusInternalServerError, "FHE001", "order_rows", &err)
		}
		defer rows2.Close()
		for rows2.Next() {
			var row roww
			err = rows2.Scan(&row.Item, &row.Qty)
			if err != nil {
				return utils.SendError(c, fiber.StatusInternalServerError, "FHE001", "order_rows", &err)
			}
			order.Rows = append(order.Rows, row)
		}
		if err = rows2.Err(); err != nil {
			return utils.SendError(c, fiber.StatusInternalServerError, "FHE004", "order_rows", &err)
		}

		ret.Orders = append(ret.Orders, order)
	}
	if err = rows.Err(); err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "FHE004", "orders", &err)
	}

	c.JSON(ret)
	return c.SendStatus(fiber.StatusOK)
}
