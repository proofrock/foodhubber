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
package db_ops

import (
	"fmt"
	"os"
	"time"

	"github.com/proofrock/foodhubber/params"
)

const vacuum_period = 5 // min

func vacuum(allowToPanic bool) {
	// Execute non-concurrently
	params.RWLock.Lock()
	defer params.RWLock.Unlock()

	if _, err := params.Db.Exec("VACUUM"); err != nil {
		if allowToPanic {
			panic(err)
		} else {
			fmt.Fprintf(os.Stderr, "vacuum: %s\n", err.Error())
		}
	}
}

func StartVacuum() {
	vacuum(true)
	for range time.Tick((vacuum_period * time.Minute)) {
		vacuum(false)
	}
}
