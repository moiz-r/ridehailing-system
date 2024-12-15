package store

import (
	"database/sql"

	_ "github.com/lib/pq"
	"github.com/moiz-r/ridehailing-system/rides-service/configs"
	"github.com/moiz-r/ridehailing-system/rides-service/internal/model"
)

type RidesStore interface {
	UpdateRide(*model.Ride) (int, error)
}

type PostgresStore struct {
	db *sql.DB
}

func NewPostgresStore(cfg *configs.DatabaseConfig) (*PostgresStore, error) {
	db, err := sql.Open("postgres", "host="+cfg.Host+" port="+cfg.Port+" user="+cfg.User+" password="+cfg.Password+" dbname="+cfg.Name+" sslmode=disable")
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &PostgresStore{
		db: db,
	}, nil
}

func (ps *PostgresStore) Init() error {
	_, err := ps.db.Exec(`
		CREATE TABLE IF NOT EXISTS rides (
			ride_id SERIAL PRIMARY KEY,
			source TEXT,
			destination TEXT,
			distance NUMERIC(10, 2),
			cost NUMERIC(10, 2)
		)
	`)
	return err
}

func (ps *PostgresStore) UpdateRide(ride *model.Ride) (int, error) {
	stmt, err := ps.db.Prepare("UPDATE rides SET source=$1, destination=$2, distance=$3, cost=$4 WHERE ride_id=$5 returning ride_id")
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	err = stmt.QueryRow(ride.Source, ride.Destination, ride.Distance, ride.Cost, ride.ID).Scan(&ride.ID)
	if err != nil {
		return 0, err
	}

	return ride.ID, nil
}
