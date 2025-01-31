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
package get_stock

import (
	"fmt"
	"time"

	"github.com/proofrock/foodhubber/params"
	"github.com/proofrock/foodhubber/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/xuri/excelize/v2"
)

type stock4excel struct {
	Item  string `json:"item"`
	Stock int    `json:"stock"`
}

func GetStockExcel(c *fiber.Ctx) error {
	params.RWLock.RLock()
	defer params.RWLock.RUnlock()

	var stocks []stock4excel

	sql := `
		SELECT i.item || COALESCE('/' || i.subitem, '') as item,  s.quantity
		  FROM items i
		  JOIN stock s ON i.id = s.item_id 
	     WHERE i.active = 1
	     ORDER BY i.pos ASC`
	rows, err := params.Db.Query(sql)
	if err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "FHE001", "stock", &err)
	}
	defer rows.Close()
	for rows.Next() {
		var stock stock4excel
		err = rows.Scan(&stock.Item, &stock.Stock)
		if err != nil {
			return utils.SendError(c, fiber.StatusInternalServerError, "FHE001", "stock", &err)
		}
		stocks = append(stocks, stock)
	}
	if err = rows.Err(); err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "FHE004", "stock", &err)
	}

	f := excelize.NewFile()
	err = f.SetColWidth("Sheet1", "A", "C", 15)
	if err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "FHE010", "", &err)
	}

	// FIXME i18n (passed by the browser)
	f.SetCellValue("Sheet1", "A1", fmt.Sprintf("Report delle scorte, %s", time.Now().Format("02/01/2006, 15:04:05")))
	err = f.MergeCell("Sheet1", "A1", "C1")
	if err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "FHE010", "", &err)
	}

	styleHeader, err := f.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
		Font: &excelize.Font{
			Bold:  true,
			Color: "FFFFFF",
		},
		Fill: excelize.Fill{Type: "pattern", Color: []string{"000000"}, Pattern: 1},
	})
	if err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "FHE010", "", &err)
	}

	err = f.SetCellValue("Sheet1", "A3", "Articolo")
	if err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "FHE010", "", &err)
	}
	err = f.SetCellStyle("Sheet1", "A3", "A3", styleHeader)
	if err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "FHE010", "", &err)
	}
	err = f.SetCellValue("Sheet1", "B3", "Quantità")
	if err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "FHE010", "", &err)
	}
	err = f.SetCellStyle("Sheet1", "B3", "B3", styleHeader)
	if err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "FHE010", "", &err)
	}
	err = f.SetCellValue("Sheet1", "C3", "Unità")
	if err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "FHE010", "", &err)
	}
	err = f.SetCellStyle("Sheet1", "C3", "C3", styleHeader)
	if err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "FHE010", "", &err)
	}

	if len(stocks) > 0 {
		styleBordered, err := f.NewStyle(&excelize.Style{
			Border: []excelize.Border{
				{Type: "left", Color: "000000", Style: 1},
				{Type: "top", Color: "000000", Style: 1},
				{Type: "bottom", Color: "000000", Style: 1},
				{Type: "right", Color: "000000", Style: 1},
			},
		})
		if err != nil {
			return utils.SendError(c, fiber.StatusInternalServerError, "FHE010", "", &err)
		}

		for i := 0; i < len(stocks); i++ {
			err = f.SetCellValue("Sheet1", fmt.Sprintf("A%d", 4+i), stocks[i].Item)
			if err != nil {
				return utils.SendError(c, fiber.StatusInternalServerError, "FHE010", "", &err)
			}
			err = f.SetCellStyle("Sheet1", fmt.Sprintf("A%d", 4+i), fmt.Sprintf("A%d", 4+i), styleBordered)
			if err != nil {
				return utils.SendError(c, fiber.StatusInternalServerError, "FHE010", "", &err)
			}
			err = f.SetCellValue("Sheet1", fmt.Sprintf("B%d", 4+i), stocks[i].Stock)
			if err != nil {
				return utils.SendError(c, fiber.StatusInternalServerError, "FHE010", "", &err)
			}
			err = f.SetCellStyle("Sheet1", fmt.Sprintf("B%d", 4+i), fmt.Sprintf("B%d", 4+i), styleBordered)
			if err != nil {
				return utils.SendError(c, fiber.StatusInternalServerError, "FHE010", "", &err)
			}
			err = f.SetCellValue("Sheet1", fmt.Sprintf("C%d", 4+i), "pezzi")
			if err != nil {
				return utils.SendError(c, fiber.StatusInternalServerError, "FHE010", "", &err)
			}
			err = f.SetCellStyle("Sheet1", fmt.Sprintf("C%d", 4+i), fmt.Sprintf("C%d", 4+i), styleBordered)
			if err != nil {
				return utils.SendError(c, fiber.StatusInternalServerError, "FHE010", "", &err)
			}
		}
	}

	c.Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"stock_%s.xlsx\"", time.Now().Format("20060102_150405")))
	len, err := f.WriteTo(c.Response().BodyWriter(), excelize.Options{})
	if err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "FHE011", "", &err)
	}
	c.Set("Content-Length", fmt.Sprint(len))

	return c.SendStatus(fiber.StatusOK)
}
