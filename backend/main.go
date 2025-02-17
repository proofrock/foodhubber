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
package main

import (
	"database/sql"
	"embed"
	"fmt"
	"net/http"

	"github.com/proofrock/foodhubber/db_ops"
	"github.com/proofrock/foodhubber/flags"
	"github.com/proofrock/foodhubber/handlers/del_order"
	"github.com/proofrock/foodhubber/handlers/del_session"
	"github.com/proofrock/foodhubber/handlers/do_cycle"
	"github.com/proofrock/foodhubber/handlers/get_beneficiary"
	"github.com/proofrock/foodhubber/handlers/get_init_data"
	"github.com/proofrock/foodhubber/handlers/get_orders"
	"github.com/proofrock/foodhubber/handlers/get_sessions"
	"github.com/proofrock/foodhubber/handlers/get_stats"
	"github.com/proofrock/foodhubber/handlers/get_stock"
	"github.com/proofrock/foodhubber/handlers/put_order"
	"github.com/proofrock/foodhubber/handlers/set_beneficiaries_excel"
	"github.com/proofrock/foodhubber/handlers/set_stock"
	"github.com/proofrock/foodhubber/params"
	"github.com/proofrock/foodhubber/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/fiber/v2/middleware/recover"
	_ "modernc.org/sqlite"
)

const DB_VERSION = 1

//go:embed static/*
var static embed.FS

func main() {
	flags.Parse()

	// FIXME don't open a new connection for each operation
	var err error
	params.Db, err = sql.Open("sqlite", params.DbPath)
	if err != nil {
		panic(err)
	}
	defer params.Db.Close()

	// check db version

	row := params.Db.QueryRow("SELECT version FROM vu_version")
	var dbVersion int
	if err := row.Scan(&dbVersion); err != nil {
		panic(err)
	}
	if dbVersion != DB_VERSION {
		utils.Abort("DB version is %d but should be %d. Please upgrade the database or the application.", dbVersion, DB_VERSION)
	}

	// VACUUM

	go db_ops.StartVacuum()

	// Backup

	db_ops.Backup()

	// server

	app := fiber.New(fiber.Config{ServerHeader: "foodhubber v." + params.VERSION, AppName: "foodhubber", DisableStartupMessage: true})

	app.Use(recover.New())

	app.Use("/", filesystem.New(filesystem.Config{
		Root:       http.FS(static),
		PathPrefix: "static",
	}))

	app.Delete("/api/delOrder", del_order.DelOrder)
	app.Delete("/api/delSession", del_session.DelSession)
	app.Post("/api/doCycle", do_cycle.DoCycle)
	app.Get("/api/getBeneficiary", get_beneficiary.GetBeneficiary)
	app.Get("/api/getBeneficiaryReport", get_beneficiary.GetBeneficiaryReport)
	app.Get("/api/getInitData", get_init_data.GetInitData)
	app.Get("/api/getOrders", get_orders.GetOrders)
	app.Get("/api/getSessions", get_sessions.GetSessions)
	app.Get("/api/getStats", get_stats.GetStats)
	app.Get("/api/getStock", get_stock.GetStock)
	app.Get("/api/getStockExcel", get_stock.GetStockExcel)
	app.Put("/api/putOrder", put_order.PutOrder)
	app.Post("/api/setBeneficiariesExcel", set_beneficiaries_excel.SetBeneficiariesExcel)
	app.Post("/api/setStock", set_stock.SetStock)

	fmt.Println("  - server on port", params.Port)
	fmt.Printf("  - all ok. Please open http://localhost:%d\n", params.Port)
	app.Listen(fmt.Sprintf(":%d", params.Port))
}
