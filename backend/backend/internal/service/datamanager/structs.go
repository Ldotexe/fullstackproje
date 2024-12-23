package datamanager

type Message struct {
	FromLogin string `db:"from_login"`
	ToLogin   string `db:"to_login"`
	Text      string `db:"txt"`
	Ts        string `db:"ts"`
}

type User struct {
	ID       uint64 `db:"id"`
	Login    string `db:"login"`
	Password string `db:"password_hash"`
	IsAdmin  bool   `db:"is_admin"`
}
