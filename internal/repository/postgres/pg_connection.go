package postgres

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/rinuccia/travels-api/config"
	"os"
)

var (
	schema = `
	CREATE TABLE IF NOT EXISTS users
	(
		user_id serial not null unique ,
		email varchar(100) not null unique,
		first_name varchar(50) not null,
		last_name varchar(50) not null,
		gender varchar(1) not null
	);
	
	CREATE TABLE IF NOT EXISTS locations
	(
		location_id serial not null unique,
		place text not null,
		country varchar(50) not null
	);
	
	CREATE TABLE IF NOT EXISTS visits
	(
		visit_id serial not null unique,
		location_id int not null references locations(location_id),
		user_id int not null references users(user_id),
		visited_at varchar(10) not null,
		mark int not null
	)`
)

func NewPostgresClient(cfg *config.Config) (*sqlx.DB, error) {
	url := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.DB.Username, os.Getenv("DB_PASSWORD"), cfg.DB.Host, cfg.DB.Port, cfg.DB.DBName)
	db, err := sqlx.Open("postgres", url)
	if err != nil {
		return nil, err
	}

	db.SetMaxIdleConns(20)
	db.SetMaxOpenConns(20)

	if err = db.Ping(); err != nil {
		return nil, err
	}

	db.MustExec(schema)

	return db, nil
}
