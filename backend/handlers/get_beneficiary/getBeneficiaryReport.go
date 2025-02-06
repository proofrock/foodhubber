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
	Date        string
	Beneficiary string
	Details     BeneficiarySituation
	Allowance   []allowanceForReport // Same as Details.Allowance but split into 2 columns
}

// TODO i18n
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
			<center><h5 class="center"><i>Report del {{.Date}}, profilo {{.Details.Profile}}</i></h5></center>
			<hr/>
			{{if .Details.TooManyOrdersInMonth}}
                <center><h3 class="center">Questo beneficiario ha effettuato più ritiri del consentito, questo mese.</h3></center>
			{{else if .Details.TooManyOrdersInWeek}}
			    <center><h3 class="center">Questo beneficiario ha effettuato un ritiro questa settimana</h3></center>
			{{else}}
				<table>
				<tr><th>Categoria</th><th>Quantità</th><th>Categoria</th><th>Quantità</th></tr>
				{{range .Allowance}}
					<tr><td>{{.Item1}}</td><td>{{.Residual1}}</td><td>{{.Item2}}</td><td>{{.Residual2}}</td></tr>
				{{end}}
			{{end}}
			</table>
		</section>
	</body>
</html>
`

func GetBeneficiaryReport(c *fiber.Ctx) error {
	id := c.Query("id", "")

	params.RWLock.RLock()
	defer params.RWLock.RUnlock()

	details, err := LoadBeneficiarySituation(id, true, c)
	if err != nil {
		return err
	}

	ret := responseForReport{
		Date:        time.Now().Format("02/01/2006"),
		Beneficiary: id,
		Details:     details,
	}

	templ, err := template.New("BeneficiaryReport").Parse(tpl)
	if err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "FHE012", "", &err)
	}

	origAllowance := ret.Details.Allowance
	ret.Allowance = make([]allowanceForReport, len(origAllowance)/2+len(origAllowance)%2) // hack: ceil(x/2) = x/2+x%2
	for i := 0; i < len(ret.Allowance); i++ {
		if i*2+1 < len(origAllowance) {
			ret.Allowance[i] = allowanceForReport{
				Item1: origAllowance[i*2].Item, Residual1: fmt.Sprint(origAllowance[i*2].Allowance),
				Item2: origAllowance[i*2+1].Item, Residual2: fmt.Sprint(origAllowance[i*2+1].Allowance),
			}
		} else {
			ret.Allowance[i] = allowanceForReport{
				Item1: origAllowance[i*2].Item, Residual1: fmt.Sprint(origAllowance[i*2].Allowance),
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
