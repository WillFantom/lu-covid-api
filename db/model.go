package db

import "time"

type Rate struct {
	ID     uint64    `db:"id"`
	Date   time.Time `db:"date"`
	Campus uint64    `db:"campus"`
	City   uint64    `db:"city"`
	Staff  uint64    `db:"staff"`
}
