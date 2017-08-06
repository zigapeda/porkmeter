package db

import (
	_ "database/sql"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var schema = `
CREATE TYPE thermometer_type AS ENUM ('meat', 'oven');

CREATE TABLE thermometer (
	id SERIAL PRIMARY KEY,
	thermometer_name VARCHAR(50),
	thermometer_type thermometer_type
);

CREATE TABLE measurement (
	id SERIAL PRIMARY KEY,
	thermometers int REFERENCES thermometer(id),
	pork_session int REFERENCES pork_session(id),
	temperature INT,
    created timestamptz NOT NULL DEFAULT NOW()
);

CREATE TABLE pork_session (
	id SERIAL PRIMARY KEY,
	created timestamptz NOT NULL DEFAULT NOW()
);`

// Dbconfig contains Database login credentials
type Dbconfig struct {
	User     string
	Password string
	Name     string
}

// Persistence contains SQLX Database object
type Persistence struct {
	db *sqlx.DB
}

// PorkSession represents a pork_session table entry
type PorkSession struct {
	Id      int64
	Created time.Time
}

// Measurement represents a measurement table entry
type Measurement struct {
	Id           int64
	Message      string
	Thermometers int64
	PorkSession  int64 `db:"pork_session"`
	Temperature  int
	Created      time.Time
}

// Thermometer represents a thermometer table entry
type Thermometer struct {
	Id              int64
	ThermometerName string `db:"thermometer_name"`
	ThermometerType string `db:"thermometer_type"`
}

// MEAT represents a meat thermometer
const MEAT string = "MEAT"

// OVEN represents an oven thermometer
const OVEN string = "OVEN"

// Connect connects a Persistence instance using a given Dbconfig configuration
func (p *Persistence) Connect(config Dbconfig) bool {
	var err error
	p.db, err = sqlx.Connect("postgres", "user="+config.Name+" password="+config.Password+" dbname="+config.Name+" sslmode=disable")

	// TODO: proper error handling & returning instead of shitty boolean (what was actually used for the unittest only)
	if err != nil {
		log.Fatalln(err)
		return false
	}
	return true

}

// Disconnect disconnects a Persistence instance
func (p *Persistence) Disconnect() {

}

// CreateNewPorkSession creates a new pork_session database entry and returns its identifier
func (p *Persistence) CreateNewPorkSession() {

}

// CreateNewOvenThermometer creates a new oven thermometer
func (p *Persistence) CreateNewOvenThermometer(thermometerName string) {

}

// CreateNewOvenThermometer creates a new meat thermometer
func (p *Persistence) CreateNewMeatThermometer(thermometerName string) {

}

// CreateSchema creates necessary database tables & types
func (p *Persistence) CreateSchema() {
	p.db.MustExec(schema)
}
