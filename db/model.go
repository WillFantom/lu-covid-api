package db

import "time"

// Rate represents the covid cases for a single day.
type Rate struct {
	Date   time.Time `db:"date"`
	Campus uint64    `db:"campus"`
	City   uint64    `db:"city"`
	Staff  uint64    `db:"staff"`
}
