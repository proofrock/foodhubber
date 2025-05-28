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

	"github.com/gofiber/fiber/v2"
	"github.com/proofrock/foodhubber/params"
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
	if !params.NoDateErrors {
		return true
	}
	weekNo := WeekOfMonth(t)
	return weekNo >= 1 && weekNo <= 4 // TODO not hard coded...
}

type ErrorWrapper struct {
	Status       int
	Code         string  `json:"code"`
	Object       string  `json:"object"`
	WrappedError *string `json:"error"`
}

func (e ErrorWrapper) Error() string {
	str, _ := json.Marshal(e)
	return string(str)
}

func MakeError(status int, errCode string, obj string, err *error) *ErrorWrapper {
	var errString *string
	if err != nil {
		_errString := (*err).Error()
		errString = &_errString
	}
	ret := ErrorWrapper{
		Status:       status,
		Code:         errCode,
		Object:       obj,
		WrappedError: errString,
	}
	return &ret
}

func SendMadeError(c *fiber.Ctx, wrappedErr ErrorWrapper) error {
	str := wrappedErr.Error()
	fmt.Fprintf(os.Stderr, "%s\n", str)
	c.JSON(wrappedErr) // FIXME marshalling is done two times
	return c.SendStatus(wrappedErr.Status)
}

func SendError(c *fiber.Ctx, status int, errCode string, obj string, err *error) error {
	wrappedErr := MakeError(status, errCode, obj, err)
	str, _ := json.Marshal(*wrappedErr)
	fmt.Fprintf(os.Stderr, "%s\n", str)
	c.JSON(*wrappedErr) // FIXME marshalling is done two times
	return c.SendStatus((*wrappedErr).Status)
}

func Int2Bool(val int) bool {
	return val != 0
}
