package db_ops

import (
	"errors"
	"fmt"
	"foodhubber/params"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

const bkpTimeFormat = "060102-150405"

var bkpTimeGlob = strings.Repeat("?", len(bkpTimeFormat))

const bkpFile = "foodhubber_%s.db"
const numFiles = 8

func Backup() {
	var bkpDir = filepath.Join(filepath.Dir(params.DbPath), "backups")
	var err error

	if _, err = os.Stat(bkpDir); errors.Is(err, os.ErrNotExist) {
		if err = os.Mkdir(bkpDir, 0755); err != nil {
			panic(err)
		}
	}

	// Execute non-concurrently
	params.RWLock.Lock()
	defer params.RWLock.Unlock()

	now := time.Now().Format(bkpTimeFormat)
	fname := fmt.Sprintf(filepath.Join(bkpDir, bkpFile), now)
	_, err = params.Db.Exec("VACUUM INTO ?", fname)
	if err != nil {
		fmt.Fprintf(os.Stderr, "backup: %s\n", err.Error())
		return
	}

	// delete the backup files, except for the last n
	list, err := filepath.Glob(fmt.Sprintf(filepath.Join(bkpDir, bkpFile), bkpTimeGlob))
	if err != nil {
		fmt.Fprintf(os.Stderr, "sched. task (pruning bkp files): %s\n", err.Error())
		return
	}

	sort.Strings(list)
	for i := 0; i < len(list)-numFiles; i++ {
		os.Remove(list[i])
	}
}
