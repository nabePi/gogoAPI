// typedef
package contler

import (
	"database/sql"
)

type AppContext struct {
	Status string
	Db     *sql.DB
}
