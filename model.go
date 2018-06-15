package main

import (
	"database/sql"
	"fmt"
)

type city struct {
	ID      int    `json:"ID"`
	Name    string `json:"name"`
	Pincode string `json:"pincode"`
}

func (c *city) getCity(db *sql.DB) error {
	statement := fmt.Sprintf("SELECT name, pincode FROM location_citys WHERE id=%d", c.ID)
	return db.QueryRow(statement).Scan(&c.Name, &c.Pincode)
}

func getCitys(db *sql.DB, start, count int) ([]city, error) {
	statement := fmt.Sprintf("SELECT id, name, pincode FROM location_citys LIMIT %d OFFSET %d", count, start)
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
