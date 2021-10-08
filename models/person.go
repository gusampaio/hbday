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
func CreateTable() error{
	_, err := DB.Query("CREATE TABLE IF NOT EXISTS \"people\" (\n    username TEXT PRIMARY KEY,\n    date_of_birth TEXT NOT NULL\n);")
	return err
}

func SetNewPerson(person Person) error{
	dbQuery := "INSERT INTO people (username, date_of_birth) VALUES ($1, $2);"
	_, err := DB.Exec(dbQuery, person.Username, person.DateOfBirth)
	return err
}

func GetPerson(username string) (Person, error){
	var person Person

	rows, err := DB.Query("SELECT * FROM people WHERE username = $1", username)
	if err != nil {
		return  person, err
	}

	defer rows.Close()
	// Loop through rows, using Scan to assign column data to struct fields.
	for rows.Next() {
		if err := rows.Scan(&person.Username, &person.DateOfBirth); err != nil {
			return person, err
		}
	}
	return person, nil
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