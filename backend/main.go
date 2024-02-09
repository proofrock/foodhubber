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
	"foodhubber/flags"
	"foodhubber/handlers/del_order"
	"foodhubber/handlers/del_session"
	"foodhubber/handlers/do_cycle"
	"foodhubber/handlers/get_beneficiary"
	"foodhubber/handlers/get_init_data"
	"foodhubber/handlers/get_orders"
	"foodhubber/handlers/get_sessions"
	"foodhubber/handlers/get_stats"
	"foodhubber/handlers/get_stock"
	"foodhubber/handlers/put_order"
	"foodhubber/handlers/set_stock"
	"foodhubber/params"
	"foodhubber/utils"
	"net/http"

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
	app.Post("/api/setStock", set_stock.SetStock)

	fmt.Println("  - server on port", params.Port)
	app.Listen(fmt.Sprintf(":%d", params.Port))
}
