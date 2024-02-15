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
