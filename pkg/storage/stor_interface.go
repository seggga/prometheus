package storage

import "github.com/seggga/prometheus/pkg/storage/model"

// CropURLStorage describes storage necessary methods to work with the application.
type CropURLStorage interface {
	// function New() lets make it easy to create an instance of the storage from specified package

	Close()                                   // close connection to the storage (database / file / ....)
	IsSet(string) bool                        // checks if the short ID is in the database
	AddURL(*model.LinkData) error             // creates a new redirect link
	Resolve(string) (string, error)           // retrieves long URL from database to produce redirect
	ViewStat(string) (*model.LinkData, error) // retrieves all data that corresponds to the short ID
	Delete(string) error                      // deletes data about specified short ID
}
