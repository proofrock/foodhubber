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
package params

import (
	"database/sql"
	"fmt"
	"math/rand"
	"sync"
)

const VERSION = "v0.0.0"

// https://manytools.org/hacker-tools/ascii-banner/, profile "Slant"
const banner = `    ______                ____  __      __    __
   / ____/___  ____  ____/ / / / /_  __/ /_  / /_  ___  _____
  / /_  / __ \/ __ \/ __  / /_/ / / / / __ \/ __ \/ _ \/ ___/
 / __/ / /_/ / /_/ / /_/ / __  / /_/ / /_/ / /_/ /  __/ /
/_/    \____/\____/\__,_/_/ /_/\__,_/_.___/_.___/\___/_/ `

var RunID int32

var RWLock sync.RWMutex

var Db *sql.DB

func init() {
	fmt.Println(banner, VERSION)
	fmt.Println()
	RunID = rand.Int31()
	fmt.Println("  - run ID", RunID)
}
