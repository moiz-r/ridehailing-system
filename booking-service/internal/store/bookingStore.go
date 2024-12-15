package store

import (
	"database/sql"
	"time"

	_ "github.com/lib/pq"
	"github.com/moiz-r/ridehailing-system/booking-service/configs"
	"github.com/moiz-r/ridehailing-system/booking-service/internal/model"
)

type BookingStore interface {
	CreateBooking(user_id int, ride model.Ride) (*model.Booking, error)
	GetBooking(int) (*model.BookingDTO, error)
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
		CREATE TABLE IF NOT EXISTS bookings (
			booking_id SERIAL PRIMARY KEY,
			ride_id INT REFERENCES rides(ride_id),
			user_id INT REFERENCES users(user_id),
			time TIMESTAMP NOT NULL
		)
	`)
	return err
}

func (ps *PostgresStore) createRide(source, destination string, distance, cost float64) (int, error) {
	stmt, err := ps.db.Prepare("INSERT INTO rides (source, destination, distance, cost) VALUES ($1, $2, $3, $4) RETURNING ride_id")
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	var rideID int
	err = stmt.QueryRow(source, destination, distance, cost).Scan(&rideID)
	if err != nil {
		return 0, err
	}

	return rideID, nil
}

func (ps *PostgresStore) CreateBooking(userID int, ride model.Ride) (*model.Booking, error) {

	rideID, err := ps.createRide(ride.Source, ride.Destination, ride.Distance, ride.Cost)
	if err != nil {
		return nil, err
	}

	stmt, err := ps.db.Prepare("INSERT INTO bookings (ride_id, user_id, time) VALUES ($1, $2, $3) RETURNING booking_id")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var booking model.Booking
	err = stmt.QueryRow(rideID, userID, time.Now()).Scan(&booking.ID)
	if err != nil {
		return nil, err
	}
	booking.RideID = rideID
	booking.UserID = userID
	return &booking, nil
}

func (ps *PostgresStore) GetBooking(id int) (*model.BookingDTO, error) {
	stmt, err := ps.db.Prepare(`
		SELECT 
			b.booking_id, 
			b.ride_id, 
			b.user_id, 
			b.time, 
			r.source,
			r.destination,
			r.cost,
			r.distance 
		FROM bookings b
		INNER JOIN rides r ON b.ride_id = r.ride_id
		WHERE b.booking_id = $1
	`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var booking model.BookingDTO
	err = stmt.QueryRow(id).Scan(
		&booking.BookingID,
		&booking.RideID,
		&booking.UserID,
		&booking.Time,
		&booking.Source,
		&booking.Destination,
		&booking.Cost,
		&booking.Distance,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &booking, nil
}
