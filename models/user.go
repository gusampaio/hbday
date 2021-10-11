package models

import (
	"database/sql"
)

// Create an exported global variable to hold the database connection pool.
var DB *sql.DB

type User struct {
	Username 	string `json:"username" gorm:"primaryKey"`
	DateOfBirth string `json:"date_of_birth"` // YYYY-MM-DD
}

func CreateTable() error{
	_, err := DB.Query("CREATE TABLE IF NOT EXISTS \"people\" (\n    username TEXT PRIMARY KEY,\n    date_of_birth TEXT NOT NULL\n);")
	return err
}

func SetNewUser(user User) error{
	dbQuery := "INSERT INTO people (username, date_of_birth) VALUES ($1, $2);"
	_, err := DB.Exec(dbQuery, user.Username, user.DateOfBirth)
	return err
}

func GetUser(username string) (User, error){
	var person User

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

// GetAllPeople returns a slice of all person in the people table.
func GetAllPeople() ([]User, error) {
	// Note that we are calling Query() on the global variable.
	rows, err := DB.Query("SELECT * FROM people")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ppl []User

	for rows.Next() {
		var person User

		err := rows.Scan(&person.Username, &person.DateOfBirth)
		if err != nil {
			return nil, err
		}

		ppl = append(ppl, person)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return ppl, nil
}