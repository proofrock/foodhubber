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
	"time"

	"github.com/proofrock/foodhubber/db_ops"
	"github.com/proofrock/foodhubber/handlers/get_beneficiary"
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

	// Re-check all the prerequisites, as someone might circumvent them by modifying the WebUI
	// Paranoid, it's me
	if !utils.IsWeekValid(time.Now()) {
		return utils.SendError(c, fiber.StatusBadRequest, "FHE106", "", nil)
	}

	details, werr := get_beneficiary.LoadBeneficiarySituation(req.Beneficiary)
	if werr != nil {
		return utils.SendMadeError(c, *werr)
	}

	if details.TooManyOrdersInWeek {
		return utils.SendError(c, fiber.StatusBadRequest, "FHE104", "", nil)
	}

	if details.TooManyOrdersInMonth {
		return utils.SendError(c, fiber.StatusBadRequest, "FHE105", "", nil)
	}

	if ok, werr := checkAllowance(details.Allowance, req.Rows); !ok {
		return utils.SendError(c, fiber.StatusBadRequest, "FHE103", "", nil)
	} else if werr != nil {
		return utils.SendMadeError(c, *werr)
	}

	tx, err := params.Db.BeginTx(context.Background(), nil)
	if err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "FHE007", "", &err)
	}

	// If there are not errors in the future, it reaches the final commit,
	// then this rollback will be harmless. If it doesn't, rolls back.
	defer tx.Rollback()

	// All is good, save

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

// checks that the allowance (that is indicated by 'Item' column, aka the "group" of an item) doesn't go < 0
// after subtracting all the order rows
func checkAllowance(allowance []get_beneficiary.Allowance, orderRows []row) (bool, *utils.ErrorWrapper) {
	// builds a map of item => allowance
	allowanceMap := make(map[string]int)
	for _, a := range allowance {
		allowanceMap[a.Item] = a.Allowance
	}

	// retrieve a map of id => Item
	itemIdMap := make(map[int]string)
	sql := `
		SELECT id, item
		  FROM items
		 WHERE active = 1`
	rows, err := params.Db.Query(sql)
	if err != nil {
		return false, utils.MakeError(fiber.StatusInternalServerError, "FHE001", "items", &err)
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		var item string
		err = rows.Scan(&id, &item)
		if err != nil {
			return false, utils.MakeError(fiber.StatusInternalServerError, "FHE001", "items", &err)
		}
		itemIdMap[id] = item
	}
	if err = rows.Err(); err != nil {
		return false, utils.MakeError(fiber.StatusInternalServerError, "FHE004", "items", &err)
	}

	// subtracts each item in the order from the allowance. If it's < 0, error
	for _, row := range orderRows {
		item := itemIdMap[row.Item]
		allowanceMap[item] = allowanceMap[item] - row.Quantity
		if allowanceMap[item] < 0 {
			return false, nil
		}
	}

	return true, nil
}
