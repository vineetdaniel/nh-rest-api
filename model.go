package main

import (
	"database/sql"
	"fmt"
	"log"
)

type city struct {
	ID      int    `json:"ID"`
	Name    string `json:"name"`
	Pincode string `json:"pincode"`
}

type location struct {
	ID   int    `json:"ID"`
	Name string `json:"name"`
	// Address     string `json:"address"`
	Pincode int `json:"pincode"`
	// Landline_no int    `json:"landline_no"`
	// Mobile_no   int    `json:"mobile_no"`
	// ContactName string `json:"contact_name"`
	// City_ID     string `json:"city_id"`
	Image       string `json:"image"`
	Thumbnail   string `json:"thumbnail"`
	Description string `json:"description"`
}

func (c *city) getCity(db *sql.DB) error {
	statement := fmt.Sprintf("SELECT name, pincode FROM locations_city WHERE id=%d", c.ID)
	return db.QueryRow(statement).Scan(&c.Name, &c.Pincode)
}

//get the list of cities
func getCities(db *sql.DB, start, count int) ([]city, error) {
	log.Print("in db query")
	statement := fmt.Sprintf("SELECT id, name, pincode FROM locations_city")
	rows, err := db.Query(statement)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	citys := []city{}

	for rows.Next() {
		var c city
		if err := rows.Scan(&c.ID, &c.Name, &c.Pincode); err != nil {
			return nil, err
		}
		citys = append(citys, c)
	}
	return citys, nil
}

///get a location
func (l *location) getLocation(db *sql.DB) error {
	statement := fmt.Sprintf("SELECT name, pincode FROM locations_locations WHERE id=%d", l.ID)
	return db.QueryRow(statement).Scan(&l.Name, &l.Pincode)
}

///get locations
func getLocations(db *sql.DB, start, count int) ([]location, error) {
	log.Print("in db query")
	statement := fmt.Sprintf("SELECT id, name, pincode, image, thumbnail, description FROM locations_locations")
	rows, err := db.Query(statement)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	locations := []location{}

	for rows.Next() {
		var l location
		if err := rows.Scan(&l.ID, &l.Name, &l.Pincode, &l.Image, &l.Thumbnail, &l.Description); err != nil {
			return nil, err
		}
		locations = append(locations, l)
	}
	return locations, nil
}
