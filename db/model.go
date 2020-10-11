package db

import "time"

type Rate struct {
	Date   time.Time `mapstructure:"date" db:"date"`
	Campus uint64    `db:"campus"`
	City   uint64    `db:"city"`
	Staff  uint64    `db:"staff"`
}
