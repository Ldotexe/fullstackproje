package repository

import "context"

type StudentRepo interface {
	Add(ctx context.Context, student *Student) error
	GetByID(ctx context.Context, id int64) (*Student, error)
	Delete(ctx context.Context, id int64) error
	Update(ctx context.Context, student *Student) error
}
