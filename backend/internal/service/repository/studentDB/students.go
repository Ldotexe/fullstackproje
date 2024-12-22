package studentDB

import (
	"context"
	"errors"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v4"
	"ss/internal/service/db"
	"ss/internal/service/repository"
)

type StudentsRepo struct {
	db db.DBops
}

func NewStudents(db db.DBops) *StudentsRepo {
	return &StudentsRepo{db: db}
}

func (r *StudentsRepo) Add(ctx context.Context, student *repository.Student) error {
	_, err := r.db.Exec(
		ctx, `INSERT INTO students(id,name,points) VALUES($1,$2,$3)`, student.ID, student.Name, student.Points,
	)
	var pgError *pgconn.PgError
	if errors.As(err, &pgError) && pgError.SQLState() == pgerrcode.UniqueViolation {
		return repository.ErrConflict
	}
	return err
}

func (r *StudentsRepo) GetByID(ctx context.Context, id int64) (*repository.Student, error) {
	var a repository.Student
	err := r.db.Get(ctx, &a, "SELECT id,name,points FROM students WHERE id=$1", id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repository.ErrObjectNotFound
		}
		return nil, err
	}
	return &a, nil
}

func (r *StudentsRepo) Update(ctx context.Context, student *repository.Student) error {
	res, err := r.db.Exec(
		ctx, "UPDATE students SET name=$1, points=$2 WHERE id=$3", student.Name, student.Points, student.ID,
	)
	if err != nil {
		return err
	}
	count := res.RowsAffected()
	if count == 0 {
		return repository.ErrObjectNotFound
	}
	return err
}

func (r *StudentsRepo) Delete(ctx context.Context, id int64) error {
	res, err := r.db.Exec(
		ctx, "DELETE FROM students WHERE id=$1", id,
	)
	if err != nil {
		return err
	}
	count := res.RowsAffected()
	if count == 0 {
		return repository.ErrObjectNotFound
	}
	return err
}
