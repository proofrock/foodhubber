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
package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/proofrock/foodhubber/params"

	"github.com/gofiber/fiber/v2"
)

func Abort(msg string, a ...any) {
	fmt.Fprintf(os.Stderr, "FATAL: %s\n", fmt.Sprintf(msg, a...))
	os.Exit(-1)
}

func FileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return !os.IsNotExist(err)
}

func WeekOfMonth(t time.Time) int {
	if params.ForcedWeek >= 0 {
		return params.ForcedWeek
	}

	year, month, _ := t.Date()
	firstOfMonth := time.Date(year, month, 1, 0, 0, 0, 0, t.Location())
	firstMonday := firstOfMonth
	for firstMonday.Weekday() != time.Monday {
		firstMonday = firstMonday.AddDate(0, 0, 1)
	}
	_, week := t.ISOWeek()
	_, firstWeek := firstMonday.ISOWeek()
	ret := week - firstWeek + 1
	if ret > 0 {
		return ret
	}
	// Allow to "break" into the following month
	return WeekOfMonth(t.AddDate(0, 0, -7)) + 1
}

func IsWeekValid(t time.Time) bool {
	weekNo := WeekOfMonth(t)
	return weekNo >= 1 && weekNo <= 4 // TODO not hard coded...
}

type errorr struct {
	Code   string  `json:"code"`
	Object string  `json:"object"`
	Error  *string `json:"error"`
}

func SendError(c *fiber.Ctx, status int, errCode string, obj string, err *error) error {
	var errString *string
	if err != nil {
		_errString := (*err).Error()
		errString = &_errString
	}
	e := errorr{
		Code:   errCode,
		Object: obj,
		Error:  errString,
	}

	str, _ := json.Marshal(e)
	fmt.Fprintf(os.Stderr, "%s\n", str)
	c.JSON(e)
	return c.SendStatus(status)
}

func Int2Bool(val int) bool {
	return val != 0
}
