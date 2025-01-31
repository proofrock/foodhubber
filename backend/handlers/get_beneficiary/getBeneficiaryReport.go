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
	"html/template"
	"time"

	"github.com/proofrock/foodhubber/params"
	"github.com/proofrock/foodhubber/utils"

	"github.com/gofiber/fiber/v2"
)

type allowanceForReport struct {
	Item1     string
	Residual1 string
	Item2     string
	Residual2 string
}

type responseForReport struct {
	Date           string
	Beneficiary    string
	EnabledForWeek bool
	Week           int
	Profile        string
	Allowance      []allowanceForReport
}

const tpl = `
<!DOCTYPE html>
<html lang="it">
	<head>
		<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
		<link rel="stylesheet" href="/css/bonsai-base.min.css">
		<link rel="stylesheet" href="/css/paper.min.css">
		<style>@page { size: A5 }</style>
	</head>
	<body class="A5" onload="window.setTimeout(window.print, 250);">
		<section class="sheet padding-10mm">
			<center><h3 class="center">Scheda beneficiario {{.Beneficiary}}</h3></center>
			<center><h5 class="center"><i>Report del {{.Date}}, settimana {{.Week}}, profilo {{.Profile}}</i></h5></center>
			<hr/>
			{{if .EnabledForWeek}}
				<table>
				<tr><th>Categoria</th><th>Quantità</th><th>Categoria</th><th>Quantità</th></tr>
				{{range .Allowance}}
					<tr><td>{{.Item1}}</td><td>{{.Residual1}}</td><td>{{.Item2}}</td><td>{{.Residual2}}</td></tr>
				{{end}}
			{{else}}
				<center><h3 class="center">Profilo non abilitato per questa settimana.</h3></center>
			{{end}}
			</table>
		</section>
	</body>
</html>
`

func GetBeneficiaryReport(c *fiber.Ctx) error {
	id := c.Query("id", "")

	weekNo := utils.WeekOfMonth(time.Now())

	ret := responseForReport{
		Date:        time.Now().Format("02/01/2006"),
		Beneficiary: id,
		Week:        utils.WeekOfMonth(time.Now()),
	}

	templ, err := template.New("BeneficiaryReport").Parse(tpl)
	if err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "FHE012", "", &err)
	}

	params.RWLock.RLock()
	defer params.RWLock.RUnlock()

	row := params.Db.QueryRow("SELECT profile FROM beneficiaries WHERE id = $1 AND active = 1", id)
	if err := row.Scan(&ret.Profile); err != nil && err != sql.ErrNoRows {
		return utils.SendError(c, fiber.StatusInternalServerError, "FHE001", "beneficiaries", &err)
	} else if err != nil {
		return utils.SendError(c, fiber.StatusNotFound, "FHE003", "", nil)
	}

	if weekNo < 1 || weekNo > 4 {
		ret.EnabledForWeek = false
	} else {
		query := fmt.Sprintf("SELECT enabled_w%d FROM vu_enabled_weeks WHERE profile = $1", weekNo)
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
		allowances := make([]allowance, 0)
		query := fmt.Sprintf(`
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
			allowances = append(allowances, allowance)
		}
		if err = rows.Err(); err != nil {
			return utils.SendError(c, fiber.StatusInternalServerError, "FHE004", "rules", &err)
		}

		ret.Allowance = make([]allowanceForReport, len(allowances)/2+len(allowances)%2) // hack: ceil(x/2) = x/2+x%2
		for i := 0; i < len(ret.Allowance); i++ {
			if i*2+1 < len(allowances) {
				ret.Allowance[i] = allowanceForReport{
					Item1: allowances[i*2].Item, Residual1: fmt.Sprint(allowances[i*2].Allowance),
					Item2: allowances[i*2+1].Item, Residual2: fmt.Sprint(allowances[i*2+1].Allowance),
				}
			} else {
				ret.Allowance[i] = allowanceForReport{
					Item1: allowances[i*2].Item, Residual1: fmt.Sprint(allowances[i*2].Allowance),
				}
			}
		}
	}

	c.Set("Content-Type", "text/html")
	err = templ.Execute(c.Response().BodyWriter(), ret)
	if err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "FHE012", "", &err)
	}

	return c.SendStatus(fiber.StatusOK)
}
