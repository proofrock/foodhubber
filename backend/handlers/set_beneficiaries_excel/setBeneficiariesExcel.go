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
package set_beneficiaries_excel

import (
	"context"
	"slices"

	"github.com/proofrock/foodhubber/db_ops"
	"github.com/proofrock/foodhubber/params"
	"github.com/proofrock/foodhubber/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/xuri/excelize/v2"
)

type response struct {
}

// TODO not hard coded
const headerForKey = "id_paziente"
const headerForProfile1 = "codice settimana"
const headerForProfile2 = "codice_settimana"

func SetBeneficiariesExcel(c *fiber.Ctx) error {
	file, err := c.FormFile("file")
	if err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "FHE014", "body", &err)
	}
	f, err := file.Open()
	if err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "FHE014", "body", &err)
	}
	defer f.Close()

	xlsx, err := excelize.OpenReader(f)
	if err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "FHE015", "body", &err)
	}
	defer xlsx.Close()

	cells, err := xlsx.GetRows(xlsx.GetSheetName(0))
	if err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "FHE015", "body", &err)
	}

	keysColumn := slices.Index[[]string, string](cells[0], headerForKey)
	profilesColumn := slices.Index[[]string, string](cells[0], headerForProfile1)
	if profilesColumn < 0 {
		profilesColumn = slices.Index[[]string, string](cells[0], headerForProfile2)
	}
	if keysColumn < 0 || profilesColumn < 0 {
		return utils.SendError(c, fiber.StatusBadRequest, "FHE200", "", nil)
	}

	defer func() { go db_ops.Backup() }()
	params.RWLock.Lock()
	defer params.RWLock.Unlock()

	tx, err := params.Db.BeginTx(context.Background(), nil)
	if err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "FHE007", "", &err)
	}
	defer tx.Rollback()

	query := "UPDATE beneficiaries SET active = 0"
	if _, err := tx.Exec(query); err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "FHE002", "beneficiaries", &err)
	}

	query = `INSERT INTO beneficiaries (id, profile) VALUES (?, ?)
	             ON CONFLICT(id) DO UPDATE SET profile = excluded.profile, active = 1`
	stmt, err := tx.Prepare(query)
	if err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "FHE013", "beneficiaries", &err)
	}
	defer stmt.Close()

	for i := 1; i < len(cells); i++ {
		_, err := stmt.Exec(cells[i][keysColumn], cells[i][profilesColumn])
		if err != nil {
			return utils.SendError(c, fiber.StatusInternalServerError, "FHE002", "beneficiaries", &err)
		}
	}

	if err := tx.Commit(); err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "FHE008", "", &err)
	}

	c.JSON(response{})
	return c.SendStatus(fiber.StatusOK)
}
