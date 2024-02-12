package pq

import (
	"database/sql"
	"database/sql/driver"
)

type Array func(a interface{}) interface {
	driver.Valuer
	sql.Scanner
}
