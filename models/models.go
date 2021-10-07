package models

import (
	"database/sql"
)

// Create an exported global variable to hold the database connection pool.
var DB *sql.DB

type Person struct {
	Username 	string `json:"username" gorm:"primaryKey"`
	DateOfBirth string `json:"date_of_birth"` // YYYY-MM-DD
}

func SetNewPerson(person Person) error{
	_, err := DB.Query("INSERT INTO people (username, date_of_birth) VALUES (%s, %s);", person.Username, person.DateOfBirth)
	return err
}

// AllBooks returns a slice of all books in the books table.
func GetAllPeople() ([]Person, error) {
	// Note that we are calling Query() on the global variable.
	rows, err := DB.Query("SELECT * FROM people")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ppls []Person

	for rows.Next() {
		var person Person

		err := rows.Scan(&person.Username, &person.DateOfBirth)
		if err != nil {
			return nil, err
		}

		ppls = append(ppls, person)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return ppls, nil
}