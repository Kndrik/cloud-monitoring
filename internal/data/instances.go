package data

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/Kndrik/cloud-monitoring/internal/validator"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Instance struct {
	Id          int           `json:"id"`
	CreatedAt   time.Time     `json:"-"`
	Name        string        `json:"name"`
	Ip          string        `json:"ip"`
	RefreshRate time.Duration `json:"refresh_rate"`
	Version     int           `json:"version"`
}

func ValidateInstance(v *validator.Validator, instance *Instance) {
	v.Check(instance.Name != "", "name", "must be provided")
	v.Check(len(instance.Name) <= 500, "name", "name must not be longer than 500 bytes")
	v.Check(instance.Ip != "", "ip", "must be provided")
	v.Check(instance.RefreshRate >= 1*time.Minute, "refresh_rate", "refresh rate must be at least one minute")
	v.Check(instance.RefreshRate <= 24*time.Hour, "refresh_rate", "refresh rate must be less than 24 hours")
}

type InstanceModel struct {
	DB *pgxpool.Pool
}

func (m *InstanceModel) Insert(instance *Instance) error {
	query := `
		INSERT INTO instances (name, ip, refresh_rate)
		VALUES ($1, $2, $3)
		RETURNING id, created_at, version`

	args := []any{instance.Name, instance.Ip, instance.RefreshRate}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return m.DB.QueryRow(ctx, query, args...).Scan(&instance.Id, &instance.CreatedAt, &instance.Version)
}

func (m *InstanceModel) GetAll() ([]*Instance, error) {
	query := "SELECT * FROM instances"

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	instances := []*Instance{}
	for rows.Next() {
		var instance Instance
		err := rows.Scan(
			&instance.Id,
			&instance.CreatedAt,
			&instance.Name,
			&instance.Ip,
			&instance.RefreshRate,
			&instance.Version,
		)
		if err != nil {
			return nil, err
		}
		instances = append(instances, &instance)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return instances, nil
}

func (m *InstanceModel) Get(id int64) (*Instance, error) {
	query := `
	SELECT id, created_at, name, ip, refresh_rate, version
	FROM instances
	WHERE id = $1`

	var instance Instance

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRow(ctx, query, id).Scan(
		&instance.Id,
		&instance.CreatedAt,
		&instance.Name,
		&instance.Ip,
		&instance.RefreshRate,
		&instance.Version,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &instance, nil
}

func (m *InstanceModel) Count() (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM instances`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRow(ctx, query).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (m *InstanceModel) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}

	query := `
		DELETE FROM instances
		WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := m.DB.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return ErrRecordNotFound
	}

	return nil
}

func (m *InstanceModel) Update(instance *Instance) error {
	query := `
	UPDATE instances
	SET name = $1, ip = $2, refresh_rate = $3, version = version + 1
	WHERE id = $4 AND version = $5
	RETURNING version`

	args := []any{
		instance.Name,
		instance.Ip,
		instance.RefreshRate,
		instance.Id,
		instance.Version,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRow(ctx, query, args...).Scan(&instance.Version)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}

	return nil
}
