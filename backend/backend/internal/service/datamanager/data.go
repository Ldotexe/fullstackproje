package datamanager

import (
	"context"
	"errors"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v4"
	"ss/internal/service/db"
)

type DataManager interface {
	AddUser(ctx context.Context, user *User) error
	AddMessage(ctx context.Context, message *Message) error
	GetUsers(ctx context.Context) ([]string, error)
	GetUser(ctx context.Context, user *User) (*User, error)
	GetMessages(ctx context.Context, message *Message) ([]Message, error)
}

type dataManager struct {
	db db.DBops
}

func New(db db.DBops) DataManager {
	return &dataManager{db: db}
}

func (d *dataManager) AddUser(ctx context.Context, user *User) error {
	_, err := d.db.Exec(
		ctx, `INSERT INTO users(login,password_hash, is_admin) VALUES($1,$2,$3)`, user.Login, user.Password, user.IsAdmin,
	)
	var pgError *pgconn.PgError
	if errors.As(err, &pgError) && pgError.SQLState() == pgerrcode.UniqueViolation {
		return ErrConflict
	}
	return err
}

func (d *dataManager) AddMessage(ctx context.Context, message *Message) error {
	_, err := d.db.Exec(
		ctx, `INSERT INTO messages(from_login, to_login, txt, ts) VALUES($1,$2,$3, $4)`, message.FromLogin, message.ToLogin, message.Text, message.Ts,
	)
	var pgError *pgconn.PgError
	if errors.As(err, &pgError) && pgError.SQLState() == pgerrcode.UniqueViolation {
		return ErrConflict
	}
	return err
}

func (d *dataManager) GetUsers(ctx context.Context) ([]string, error) {
	logins := make([]string, 0)
	err := d.db.Select(ctx, &logins, "SELECT login FROM users")
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrObjectNotFound
		}
		return nil, err
	}
	return logins, nil
}

func (d *dataManager) GetUser(ctx context.Context, user *User) (*User, error) {
	users := make([]User, 0)
	err := d.db.Select(ctx, &users, "SELECT * FROM users WHERE login=$1", user.Login) // What if we have more than one user with same login?
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrObjectNotFound
		}
		return nil, err
	}
	if len(users) == 0 {
		return nil, ErrObjectNotFound
	}
	if len(users) != 1 {
		return nil, ErrDuplication
	}
	return &users[0], nil
}

func (d *dataManager) GetMessages(ctx context.Context, message *Message) ([]Message, error) {
	messages := make([]Message, 0)
	err := d.db.Select(ctx, &messages, "SELECT * FROM messages WHERE (from_login=$1 AND to_login=$2) OR (from_login=$2 AND to_login=$1)", message.ToLogin, message.FromLogin)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrObjectNotFound
		}
		return nil, err
	}
	return messages, nil
}
