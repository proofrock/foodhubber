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
package flags

import (
	"flag"
	"foodhubber/params"
	"foodhubber/utils"
)

func Parse() {
	_db := flag.String("db", "./foodhubber.db", "The path of the sqlite database; defaults to './foodhubber.db'")
	_port := flag.Int("port", 31020, "Port; defaults to 31020")
	_forcedWeek := flag.Int("force-week", -1, "Forced week; for debug")

	flag.Parse()

	if !utils.FileExists(*_db) {
		utils.Abort("missing database file '%s'", *_db)
	}

	params.DbPath = *_db
	params.Port = *_port
	params.ForcedWeek = *_forcedWeek
}
