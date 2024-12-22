package repository

type Student struct {
	ID     int64  `db:"id"`
	Name   string `db:"name"`
	Points int64  `db:"points"`
}
