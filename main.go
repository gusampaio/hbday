package main

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"hbday/web-service-gin/models"
	"net/http"
	"os"
	"time"
	"unicode"
)

var (
	host     = os.Getenv("DATABASE_HOST")
	user     = os.Getenv("POSTGRES_USER")
	password = os.Getenv("POSTGRES_PASSWORD")
	dbname   = os.Getenv("POSTGRES_DB")
)

func main() {
	ConnectDb()
	// close database
	defer models.DB.Close()

	router := gin.Default()
	router.GET("/hello/:username", getPersonByUsername)
	router.GET("/hello/all", getAll)
	router.PUT("/hello/:username", postPerson)
	router.Run("0.0.0.0:8080") //nolint:errcheck
}

 func ConnectDb() {
	 var err error
	 // connection string
	 psqlconn := fmt.Sprintf("host=%s port=5432 dbname=%s user=%s  password=%s sslmode=disable", host, dbname, user, password)

	 // open database
	 models.DB, err = sql.Open("postgres", psqlconn)
	 if err != nil {
		 panic(err)
	 }

	 // check db
	 err = models.DB.Ping()
	 if err != nil {
		 panic(err)
	 }
	 fmt.Println("Connected!")
	 // creating db table if doesn't exist
	 models.CreateTable()
 }

func getAll(c *gin.Context) {
	all, err := models.GetAllPeople()
	if err != nil {
		panic(err)
	}
	c.IndentedJSON(http.StatusOK, all)
}

// getPersonByUsername locates the album whose ID value matches the id
// parameter sent by the client, then returns that album as a response.
func getPersonByUsername(c *gin.Context) {
	// checking if username already exist
	username := c.Param("username")
	messageSlice := make(map[string]string)
	if !exist(username){
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "username not found"})
		return
	}

	// Retrieving username in the database
	newP, err := models.GetUser(username)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	days, err := daysToBirthday(newP.DateOfBirth)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if days < 1 {
		messageSlice["message"] = "Hello " + newP.Username + "! Happy Birthday!"

	} else {
		messageSlice["message"] = fmt.Sprintf("Hello %s! Your Bithday is in %d days", newP.Username, days)
	}
	c.IndentedJSON(http.StatusOK, messageSlice)
	return
}

// postPerson adds a person from JSON received in the request body.
func postPerson(c *gin.Context) {
	var newPerson models.User

	// Call BindJSON to bind the received JSON to
	// newPerson.
	if err := c.BindJSON(&newPerson); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if !isLetter(c.Param("username")){
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "invalid username"})
		return
	}
	// Including username param as not being passed in the JSON
	newPerson.Username = c.Param("username");

	// User already exist
	if exist(newPerson.Username) {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "User already exist"})
		return
	}

	// Invalid Date of Birth
	err := isValidDoB(newPerson.DateOfBirth)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	// Add the new person to the db
	err = models.SetNewUser(newPerson)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
	}

	c.IndentedJSON(http.StatusNoContent, newPerson)
}

func exist(username string) bool{
	newP, _ := models.GetUser(username)
	return len(newP.Username) != 0
}

func isLetter(s string) bool {
	for _, r := range s {
		if !unicode.IsLetter(r) {
			return false
		}
	}
	return true
}

func daysToBirthday(dateOfBirth string) (int, error){
	today := time.Now()

	// parsing birthday the string into date
	date, err := parseDate(dateOfBirth)
	if err != nil {
		return -1, err
	}
	// creating new date based on user birthday in the current year
	// if the birthday is before the current date, change to next year
	birthday := time.Date(today.Year(), date.Month(), date.Day(), 23, 59, 59, 0, time.UTC)
	if birthday.Before(today) {
		birthday = time.Date(today.Year()+1, date.Month(), date.Day(), 23, 59, 59, 0, time.UTC)
	}
	// if the birthday still this year, calculate no need to increment one year
	daysToBirthday := int(birthday.Sub(today).Hours()/24)
	return daysToBirthday, err
}

// validate if date is in the correct format
func isValidDoB(dateOfBirth string) error {
	today := time.Now()

	// Validating if the provided date is defined as expected YYYY-MM-DD
	date, err := parseDate(dateOfBirth)
	if err != nil {
		return err
	}

	// If today date is not "after" the provided date, than error
	if !today.After(date) {
		return errors.New("Invalid Date of Birth, date must be not greater than today")
	}

	return err
}

func parseDate(dateOfBirth string) (time.Time, error){
	layout := "2006-01-02"

	// Validating if the provided date is defined as expected YYYY-MM-DD
	t, err := time.Parse(layout, dateOfBirth)
	if err != nil {
		return t, errors.New("Incorrect date format, expected YYYY-MM-DD")
	}
	return t, nil
}

