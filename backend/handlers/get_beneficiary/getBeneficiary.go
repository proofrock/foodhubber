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

type order struct {
	ID       int    `json:"id"`
	Date     string `json:"date"`
	ThisWeek bool   `json:"thisWeek"`
}

type allowance struct {
	Item      string `json:"item"`
	Allowance int    `json:"allowance"`
}

	WeekIsOk             bool        `json:"weekIsOk"`
}

func GetBeneficiary(c *fiber.Ctx) error {
	id := c.Query("id", "")

	weekNo := utils.WeekOfMonth(time.Now())

	ret := response{
		Allowance: make([]allowance, 0),
	}

	params.RWLock.RLock()
	defer params.RWLock.RUnlock()

	query := `
		SELECT id, strftime('%Y%m%dT%H%M%S', datetime) AS datetime,
		       UNIXEPOCH(datetime) >= UNIXEPOCH(DATE(DATETIME('now', 'localtime'), 'weekday 0', '-7 days') || ' 00:00:00') AS inThisWeek
		  FROM orders 
		 WHERE beneficiary_id = $1 
		   AND active = 1
		 ORDER BY datetime DESC
		 LIMIT 1`
	row := params.Db.QueryRow(query, id)
	var lastOrder order
	if err := row.Scan(&lastOrder.ID, &lastOrder.Date, &lastOrder.ThisWeek); err != nil && err != sql.ErrNoRows {
		return utils.SendError(c, fiber.StatusInternalServerError, "FHE001", "orders", &err)
	} else if err == nil {
		ret.LastOrder = &lastOrder
	}

	row = params.Db.QueryRow("SELECT profile FROM beneficiaries WHERE id = $1 AND active = 1", id)
	if err := row.Scan(&ret.Profile); err != nil && err != sql.ErrNoRows {
		return utils.SendError(c, fiber.StatusInternalServerError, "FHE001", "beneficiaries", &err)
	} else if err != nil {
		return utils.SendError(c, fiber.StatusNotFound, "FHE003", "", nil)
	}

	if weekNo < 1 || weekNo > 4 {
		ret.EnabledForWeek = false
	} else {
		query = fmt.Sprintf("SELECT enabled_w%d FROM vu_enabled_weeks WHERE profile = $1", weekNo)
		var enabledForWeek int
		row = params.Db.QueryRow(query, ret.Profile)
		if err := row.Scan(&enabledForWeek); err != nil && err != sql.ErrNoRows {
			return utils.SendError(c, fiber.StatusInternalServerError, "FHE001", "vu_enabled_weeks", &err)
		} else if err != nil {
			return utils.SendError(c, fiber.StatusNotFound, "FHE009", "", nil)
		}
		ret.EnabledForWeek = utils.Int2Bool(enabledForWeek)
	}

	if ret.EnabledForWeek {
		query = fmt.Sprintf(`
			SELECT r.item, r.quantity_w%d AS allowance
			  FROM rules r
			  JOIN vu_items_lvl_1 il1 ON r.item = il1.item
			 WHERE r.profile = $2
			 ORDER BY il1.pos ASC`, weekNo)
		rows, err := params.Db.Query(query, id, ret.Profile)
		if err != nil {
			return utils.SendError(c, fiber.StatusInternalServerError, "FHE001", "rules", &err)
		}
		defer rows.Close()
		for rows.Next() {
			var allowance allowance
			err = rows.Scan(&allowance.Item, &allowance.Allowance)
			if err != nil {
				return utils.SendError(c, fiber.StatusInternalServerError, "FHE001", "rules", &err)
			}
			ret.Allowance = append(ret.Allowance, allowance)
		}
		if err = rows.Err(); err != nil {
			return utils.SendError(c, fiber.StatusInternalServerError, "FHE004", "rules", &err)
		}
	}

	c.JSON(ret)
	return c.SendStatus(fiber.StatusOK)
}
