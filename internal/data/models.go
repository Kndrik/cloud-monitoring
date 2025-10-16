package data

import "github.com/jackc/pgx/v5/pgxpool"

type Models struct {
	Instances InstanceModel
}

func NewModels(db *pgxpool.Pool) Models {
	return Models{
		Instances: InstanceModel{DB: db},
	}
}
