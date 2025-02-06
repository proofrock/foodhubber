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
package get_beneficiary

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/proofrock/foodhubber/params"
	"github.com/proofrock/foodhubber/utils"

	"github.com/gofiber/fiber/v2"
)

type Order struct {
	ID   int    `json:"id"`
	Date string `json:"date"`
}

type Allowance struct {
	Item      string `json:"item"`
	Allowance int    `json:"allowance"`
}

type BeneficiarySituation struct {
	Profile              string      `json:"profile"`
	LastOrder            *Order      `json:"lastOrder"`
	OrdersInMonth        int         `json:"ordersInMonth"`        // orders since the first monday of this month
	TooManyOrdersInMonth bool        `json:"tooManyOrdersInMonth"` // are they more than allowed for profile?
	TooManyOrdersInWeek  bool        `json:"tooManyOrdersInWeek"`
	WeekIsOk             bool        `json:"weekIsOk"`
	Allowance            []Allowance `json:"allowance"`
}

func GetBeneficiary(c *fiber.Ctx) error {
	id := c.Query("id", "")

	params.RWLock.RLock()
	defer params.RWLock.RUnlock()

	ret, err := LoadBeneficiarySituation(id, true, c)
	if err != nil {
		return err
	}

	c.JSON(ret)
	return c.SendStatus(fiber.StatusOK)
}

// Must lock before it!
func LoadBeneficiarySituation(id string, loadAllowance bool, c *fiber.Ctx) (BeneficiarySituation, error) {
	ret := BeneficiarySituation{
		WeekIsOk:  utils.IsWeekValid(time.Now()),
		Allowance: make([]Allowance, 0),
	}

	monthlyOrdersAllowed := 0

	query := `
		SELECT b.profile, mop.num 
			FROM beneficiaries b
			INNER JOIN vu_monthly_orders_by_profile mop ON b.profile = mop.profile
			WHERE b.id = ?
			AND b.active = 1
		`
	row := params.Db.QueryRow(query, id)
	if err := row.Scan(&ret.Profile, &monthlyOrdersAllowed); err != nil && err != sql.ErrNoRows {
		return ret, utils.SendError(c, fiber.StatusInternalServerError, "FHE001", "vu_monthly_orders_by_profile", &err)
	} else if err != nil {
		return ret, utils.SendError(c, fiber.StatusNotFound, "FHE009", "", nil)
	}

	// Was the last order done this week?
	// ---
	// In the following query, DATE(DATETIME('now', 'localtime'), 'weekday 1', '-7 days')
	// is the date of last sunday, from yesterday. So if it's sunday, it's not today, but one
	// week ago. Given that sunday it's not a working day it's acceptable. Must be changed
	// if sunday IS a working day.
	query = `
		SELECT id, strftime('%Y%m%dT%H%M%S', datetime) AS datetime,
		       UNIXEPOCH(datetime) >= UNIXEPOCH(DATE(DATETIME('now', 'localtime'), 'weekday 0', '-7 days') || ' 00:00:00') AS inThisWeek
		  FROM orders 
		 WHERE beneficiary_id = $1 
		   AND active = 1
		 ORDER BY datetime DESC
		 LIMIT 1
		`
	row = params.Db.QueryRow(query, id)
	var lastOrder Order
	if err := row.Scan(&lastOrder.ID, &lastOrder.Date, &ret.TooManyOrdersInWeek); err != nil && err != sql.ErrNoRows {
		return ret, utils.SendError(c, fiber.StatusInternalServerError, "FHE001", "orders", &err)
	} else if err == nil {
		ret.LastOrder = &lastOrder
	}

	// How many orders were made this month, and are they too many for the profile?
	query = `
		SELECT COUNT(1) AS cnt
		  FROM orders
		 WHERE beneficiary_id = ?
		   AND UNIXEPOCH(datetime) > UNIXEPOCH(DATETIME('now', 'start of month', 'weekday 1'))
		`
	row = params.Db.QueryRow(query, id)
	if err := row.Scan(&ret.OrdersInMonth); err != nil {
		return ret, utils.SendError(c, fiber.StatusInternalServerError, "FHE001", "orders", &err)
	}
	ret.TooManyOrdersInMonth = ret.OrdersInMonth >= monthlyOrdersAllowed

	// Given this is the (ret.OrdersInMonth + 1)th order this month, what is the allowance?
	if loadAllowance {
		query = fmt.Sprintf(`
			SELECT r.item, r.quantity_o%d AS allowance
			  FROM rules r
			  JOIN vu_items_lvl_1 il1 ON r.item = il1.item
			 WHERE r.profile = $1
			 ORDER BY il1.pos ASC`, ret.OrdersInMonth+1)
		rows, err := params.Db.Query(query, ret.Profile)
		if err != nil {
			return ret, utils.SendError(c, fiber.StatusInternalServerError, "FHE001", "rules", &err)
		}
		defer rows.Close()
		for rows.Next() {
			var allowance Allowance
			err = rows.Scan(&allowance.Item, &allowance.Allowance)
			if err != nil {
				return ret, utils.SendError(c, fiber.StatusInternalServerError, "FHE001", "rules", &err)
			}
			ret.Allowance = append(ret.Allowance, allowance)
		}
		if err = rows.Err(); err != nil {
			return ret, utils.SendError(c, fiber.StatusInternalServerError, "FHE004", "rules", &err)
		}
	}

	return ret, nil
}
