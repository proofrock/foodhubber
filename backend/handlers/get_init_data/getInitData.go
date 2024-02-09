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
package get_init_data

import (
	"foodhubber/params"
	"foodhubber/utils"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type item struct {
	ID      int     `json:"id"`
	Color   string  `json:"color"`
	Item    string  `json:"item"`
	Subitem *string `json:"subitem"`
}

type checkout struct {
	ID                     string  `json:"id"`
	CanAccessOrderListPage bool    `json:"can_access_order_list_page"`
	CanAccessStatsPage     bool    `json:"can_access_stats_page"`
	CanAccessStockPage     bool    `json:"can_access_stock_page"`
	CanAccessConsolePage   bool    `json:"can_access_console_page"`
	Password               *string `json:"password"`
}

type response struct {
	Version      string     `json:"version"`
	PollingCycle int        `json:"polling_cycle"`
	WarnLimit    int        `json:"warn_limit"`
	Items        []item     `json:"items"`
	Checkouts    []checkout `json:"checkouts"`
}

func GetInitData(c *fiber.Ctx) error {
	ret := response{
		Version: params.VERSION,
	}

	params.RWLock.RLock()
	defer params.RWLock.RUnlock()

	var err error

	row := params.Db.QueryRow("SELECT value FROM configs WHERE key = $1", "polling_cycle")
	var val string
	if err = row.Scan(&val); err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "FHE001", "configs", &err)
	}
	if ret.PollingCycle, err = strconv.Atoi(val); err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "FHE005", "polling_cycle", &err)
	}

	row = params.Db.QueryRow("SELECT value FROM configs WHERE key = $1", "warn_limit")
	if err = row.Scan(&val); err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "FHE001", "configs", &err)
	}
	if ret.WarnLimit, err = strconv.Atoi(val); err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "FHE005", "warn_limit", &err)
	}

	sql := `
		SELECT id, color, item, subitem
		  FROM items
		 WHERE active = 1
		 ORDER BY pos ASC`
	rows, err := params.Db.Query(sql)
	if err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "FHE001", "items", &err)
	}
	defer rows.Close()
	for rows.Next() {
		var item item
		err = rows.Scan(&item.ID, &item.Color, &item.Item, &item.Subitem)
		if err != nil {
			return utils.SendError(c, fiber.StatusInternalServerError, "FHE001", "items", &err)
		}
		ret.Items = append(ret.Items, item)
	}
	if err = rows.Err(); err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "FHE004", "items", &err)
	}

	sql = `
		SELECT id, can_access_order_list_page, can_access_stats_page, can_access_stock_page, can_access_console_page,
		       password_hash
          FROM checkouts
		 WHERE active = 1
		 ORDER BY pos ASC`
	rows, err = params.Db.Query(sql)
	if err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "FHE001", "checkouts", &err)
	}
	defer rows.Close()
	for rows.Next() {
		var checkout checkout
		var v1, v2, v3, v4 int
		err = rows.Scan(&checkout.ID, &v1, &v2, &v3, &v4, &checkout.Password)
		if err != nil {
			return utils.SendError(c, fiber.StatusInternalServerError, "FHE001", "checkouts", &err)
		}
		checkout.CanAccessOrderListPage = utils.Int2Bool(v1)
		checkout.CanAccessStatsPage = utils.Int2Bool(v2)
		checkout.CanAccessStockPage = utils.Int2Bool(v3)
		checkout.CanAccessConsolePage = utils.Int2Bool(v4)
		ret.Checkouts = append(ret.Checkouts, checkout)
	}
	if err = rows.Err(); err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "FHE004", "checkouts", &err)
	}

	c.JSON(ret)
	return c.SendStatus(fiber.StatusOK)
}
