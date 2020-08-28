package controllers

import "learning/unit-testing/database"

type Ctlr struct {
	DB database.Storage
}

func GetControllerDB() (Ctlr, error) {

	ctlr := Ctlr{}

	dbConnection, err := database.NewConnection()
	if err != nil {
		return ctlr, err
	}

	ctlr.DB = dbConnection

	return ctlr, nil
}

func GetControllerMockDB() Ctlr {

	dbConnection := database.NewMockConnection()

	return Ctlr{DB: dbConnection}
}
