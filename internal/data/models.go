package data

import (
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict   = errors.New("edit conflict")
)

type Models struct {
	Instances InstanceModel
}

func NewModels(db *pgxpool.Pool) Models {
	return Models{
		Instances: InstanceModel{DB: db},
	}
}
