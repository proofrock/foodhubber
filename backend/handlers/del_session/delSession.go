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
package del_session

import (
	"github.com/proofrock/foodhubber/params"
	"github.com/proofrock/foodhubber/utils"

	"github.com/gofiber/fiber/v2"
)

type response struct {
}

func DelSession(c *fiber.Ctx) error {
	params.RWLock.Lock()
	defer params.RWLock.Unlock()

	id := c.Query("id")

	query := "DELETE FROM sessions WHERE checkout_id = $1"
	if _, err := params.Db.Exec(query, id); err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "FHE002", "sessions", &err)
	}

	c.JSON(response{})
	return c.SendStatus(fiber.StatusOK)
}
