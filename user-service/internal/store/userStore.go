package store

import (
	"database/sql"

	_ "github.com/lib/pq"
	"github.com/moiz-r/ridehailing-system/user-service/configs"
	"github.com/moiz-r/ridehailing-system/user-service/internal/model"
)

type UserStore interface {
	CreateUser(string) (*model.User, error)
	GetUser(int) (*model.User, error)
	DeleteUser(int) error
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
	_, err := ps.db.Exec("CREATE TABLE IF NOT EXISTS users (user_id SERIAL PRIMARY KEY, name VARCHAR(50))")
	return err
}
func (ps *PostgresStore) CreateUser(name string) (*model.User, error) {
	stmt, err := ps.db.Prepare("INSERT INTO users (name) VALUES ($1) RETURNING user_id")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var user model.User
	err = stmt.QueryRow(name).Scan(&user.UserID)
	if err != nil {
		return nil, err
	}
	user.Name = name
	return &user, nil
}

func (ps *PostgresStore) GetUser(id int) (*model.User, error) {
	stmt, err := ps.db.Prepare("SELECT user_id, name FROM users WHERE user_id = $1")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	row := stmt.QueryRow(id)
	user := &model.User{}
	err = row.Scan(&user.UserID, &user.Name)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (ps *PostgresStore) DeleteUser(id int) error {
	stmt, err := ps.db.Prepare("DELETE FROM users WHERE user_id = $1")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(id)
	return err
}
