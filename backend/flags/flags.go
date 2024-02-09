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
