package db

import (
	_ "database/sql"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var schema = `
CREATE TYPE thermometer AS ENUM ('meat', 'oven');

CREATE TABLE thermometer (
	id serial primary_key,
	thermometer_name text,
	thermometer_type thermometer,
);

CREATE TABLE measurement (
	id serial primary_key
	thermometers int references thermometer(id),
	temperature int,
    time timestamp
);

CREATE TABLE pork_session (
	id serial primary_key
	measurements int references measurement(id)
);`

type Person struct {
	FirstName string `db:"first_name"`
	LastName  string `db:"last_name"`
	Email     string
}

type Place struct {
	Country string
	TelCode int
}
type Dbconfig struct {
	User     string
	Password string
	Name     string
}

type persistence struct {
	db *sqlx.DB
}

func (p persistence) doThings() string {

	return ""
}
