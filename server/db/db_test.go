package db

import "testing"

/*
	Create postgres db "testing" owned by user "testing" with password "testing" before
*/

var dropSchema = `
DROP TABLE IF EXISTS thermometer, measurement, pork_session;
DROP TYPE IF EXISTS thermometer_type;`

func TestConnect(t *testing.T) {
	testConfig := Dbconfig{"testing", "testing", "testing"}
	p := Persistence{}

	if !p.Connect(testConfig) {
		t.Errorf("Database Login failed")
	}
}

func TestCreateSchema(t *testing.T) {
	testConfig := Dbconfig{"testing", "testing", "testing"}
	p := Persistence{}

	if !p.Connect(testConfig) {
		t.Errorf("Database Login failed")
	}

	p.CreateSchema()
	p.db.MustExec(dropSchema)
}
