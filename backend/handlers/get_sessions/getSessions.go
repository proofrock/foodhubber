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
package get_sessions

import (
	"foodhubber/params"
	"foodhubber/utils"

	"github.com/gofiber/fiber/v2"
)

type session struct {
	ID       string `json:"id"`
	Operator string `json:"operator"`
	DateTime string `json:"datetime"`
	Active   bool   `json:"active"`
}

type response struct {
	Checkouts []session `json:"checkouts"`
}

func GetSessions(c *fiber.Ctx) error {
	ret := response{Checkouts: make([]session, 0)}

	params.RWLock.RLock()
	defer params.RWLock.RUnlock()

	sql := `
		SELECT s.checkout_id, s.operator, s.active,
		       strftime('%Y%m%dT%H%M%S', datetime) AS datetime
          FROM vu_active_sessions s
          JOIN checkouts c ON s.checkout_id = c.id
         ORDER BY c.pos ASC`
	rows, err := params.Db.Query(sql)
	if err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "FHE001", "ops.VU_ACTIVE_SESSIONS", &err)
	}
	defer rows.Close()
	for rows.Next() {
		var session session
		err = rows.Scan(&session.ID, &session.Operator, &session.Active, &session.DateTime)
		if err != nil {
			return utils.SendError(c, fiber.StatusInternalServerError, "FHE001", "ops.VU_ACTIVE_SESSIONS", &err)
		}
		ret.Checkouts = append(ret.Checkouts, session)
	}
	if err = rows.Err(); err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "FHE004", "ops.VU_ACTIVE_SESSIONS", &err)
	}

	c.JSON(ret)
	return c.SendStatus(fiber.StatusOK)
}
