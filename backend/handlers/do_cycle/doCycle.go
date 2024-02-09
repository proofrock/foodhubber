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
package do_cycle

import (
	"foodhubber/params"
	"foodhubber/utils"
	"time"

	"github.com/gofiber/fiber/v2"
)

type response struct {
	Week      int   `json:"week"`
	RunID     int32 `json:"run_id"`
	GenStock  int   `json:"gen_stock"`
	GenOrders int   `json:"gen_orders"`
}

func DoCycle(c *fiber.Ctx) error {
	checkout := c.Query("pos", "")
	operator := c.Query("op", "")

	params.RWLock.Lock()
	defer params.RWLock.Unlock()

	if checkout != "" {
		query := `
			INSERT INTO sessions (checkout_id, operator, datetime)
				 VALUES ($1, $2, DATETIME('now', 'localtime'))
			ON CONFLICT (checkout_id) DO
				 UPDATE
					SET operator = EXCLUDED.operator,
						datetime = EXCLUDED.datetime`
		if _, err := params.Db.Exec(query, checkout, operator); err != nil {
			return utils.SendError(c, fiber.StatusInternalServerError, "FHE002", "sessions", &err)
		}
	}

	c.JSON(response{
		Week:      utils.WeekOfMonth(time.Now()),
		RunID:     params.RunID,
		GenOrders: params.GetOrderGen(),
		GenStock:  params.GetStockGen(),
	})
	return c.SendStatus(fiber.StatusOK)
}
