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
package put_order

import (
	"context"
	"fmt"
	"time"

	"github.com/proofrock/foodhubber/db_ops"
	"github.com/proofrock/foodhubber/params"
	"github.com/proofrock/foodhubber/utils"

	"github.com/gofiber/fiber/v2"
)

type row struct {
	Item     int `json:"item"`
	Quantity int `json:"quantity"`
}

type request struct {
	Checkout    string  `json:"checkout"`
	Operator    string  `json:"operator"`
	Beneficiary string  `json:"beneficiary"`
	Note        *string `json:"note"`
	Rows        []row   `json:"rows"`
}

type response struct {
	Id            int  `json:"id"`
	ExceededStock bool `json:"exceeded_stock"`
}

func PutOrder(c *fiber.Ctx) error {
	req := new(request)
	if err := c.BodyParser(req); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "FHE006", "body", &err)
	}

	if len(req.Rows) <= 0 {
		return utils.SendError(c, fiber.StatusBadRequest, "FHE100", "", nil)
	}

	if req.Note != nil && len(*req.Note) > 64 {
		return utils.SendError(c, fiber.StatusBadRequest, "FHE101", "", nil)
	}

	defer func() { go db_ops.Backup() }()
	params.RWLock.Lock()
	defer params.RWLock.Unlock()

	tx, err := params.Db.BeginTx(context.Background(), nil)
	if err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "FHE007", "", &err)
	}
	defer tx.Rollback()

	var res response

	query := `
		INSERT INTO orders (checkout_id, operator, beneficiary_id, note)
		     VALUES ($1, $2, $3, $4)
		  RETURNING id`
	row := tx.QueryRow(query, req.Checkout, req.Operator, req.Beneficiary, req.Note)
	if err := row.Scan(&res.Id); err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "FHE002", "orders", &err)
	}

	query = `
		INSERT INTO order_rows (order_id, item_id, quantity)
	         VALUES ($1, $2, $3);`
	for i := 0; i < len(req.Rows); i++ {
		if _, err = tx.Exec(query, res.Id, req.Rows[i].Item, req.Rows[i].Quantity); err != nil {
			return utils.SendError(c, fiber.StatusInternalServerError, "FHE002", "order_rows", &err)
		}
	}

	query = `
		UPDATE stock
	       SET quantity = stock.quantity - orw.quantity
          FROM order_rows orw
		 WHERE stock.item_id = orw.item_id 
           AND orw.order_id = $1`
	if _, err = tx.Exec(query, res.Id); err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "FHE002", "stock", &err)
	}

	weekNo := utils.WeekOfMonth(time.Now())
	query = fmt.Sprintf(`
			WITH ORDERED AS (
			  SELECT il1.item, SUM(orw.quantity) as quantity
			    FROM vu_items_lvl_1 il1
			    JOIN items i ON il1.item = i.item 
			    JOIN order_rows orw ON i.id = orw.item_id
			    JOIN orders o ON orw.order_id = o.id 
			   WHERE o.beneficiary_id = $1 
			     AND o.active = 1
			     AND o.datetime >= DATE(DATETIME('now', 'localtime'), 'weekday 1', '-7 days') || ' 00:00:00'
			   GROUP BY il1.item)
			, RESIDUAL AS (
			  SELECT r.item, r.quantity_w%d - COALESCE(o.quantity, 0) AS residual
			    FROM rules r
			    JOIN beneficiaries b ON r.profile = b.profile
			    LEFT JOIN ORDERED o ON r.item = o.item
			   WHERE b.id = $2)
			SELECT EXISTS (SELECT 1 FROM RESIDUAL WHERE residual < 0) AS err`, weekNo)
	row = tx.QueryRow(query, req.Beneficiary, req.Beneficiary)
	var isErr int
	if err := row.Scan(&isErr); err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "FHE001", "ops.RULES", &err)
	} else if utils.Int2Bool(isErr) {
		return utils.SendError(c, fiber.StatusBadRequest, "FHE103", "", nil)
	}

	params.TouchOrdersGen()
	params.TouchStockGen()

	// XXX stock can be < 0. The mere fact that the beneficiary carries an item with them is
	//     proof that there's a stock for it.
	query = "SELECT EXISTS (SELECT 1 FROM stock WHERE quantity < 0) AS exceeded"
	var exceededStock int
	row = tx.QueryRow(query)
	if err := row.Scan(&exceededStock); err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "FHE001", "stock", &err)
	}

	if err := tx.Commit(); err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "FHE008", "", &err)
	}

	res.ExceededStock = utils.Int2Bool(exceededStock)

	c.JSON(res)
	return c.SendStatus(fiber.StatusOK)
}
