package db_ops

import (
	"fmt"
	"foodhubber/params"
	"os"
	"time"
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
	for range time.Tick((vacuum_period * time.Second)) {
		vacuum(false)
	}
}
