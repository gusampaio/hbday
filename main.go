package main

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"hbday/web-service-gin/models"
	"net/http"
	"time"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "gusampaio"
	password = "gusampaio_pass"
	dbname   = "hbday_db"
)

func main() {
	connectDb()
	// close database
	defer models.DB.Close()

	router := gin.Default()
	router.GET("/hello/:username", getPersonByUsername)
	router.GET("/hello/all", getAll)
	router.PUT("/hello/:username", postPerson)
	router.Run("localhost:8080") //nolint:errcheck
}
 func connectDb() {
	 var err error
	 // connection string
	 psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

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
	username := c.Param("username")
	if !exist(username){
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "username not found"})
		return
	}

	newP, err := models.GetPerson(username)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, newP)
}

// postPerson adds a person from JSON received in the request body.
func postPerson(c *gin.Context) {
	var newPerson models.Person

	// Call BindJSON to bind the received JSON to
	// newPerson.
	if err := c.BindJSON(&newPerson); err != nil {
		return
	}

	// Including username param as not being passed in the JSON
	newPerson.Username = c.Param("username")

	if exist(newPerson.Username) {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "User already exist"})
		return
	}

	err := isValidDoB(newPerson.DateOfBirth)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	err = models.SetNewPerson(newPerson)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
	}

	// Add the new person to the slice.
	c.IndentedJSON(http.StatusNoContent, newPerson)
}

func exist(username string) bool{
	newP, _ := models.GetPerson(username)
	return len(newP.Username) != 0
}


func isValidDoB(dateOfBirthday string) error {
	layout := "2006-01-02"
	today := time.Now()

	// Validating if the provided date is defined as expected YYYY-MM-DD
	t, err := time.Parse(layout, dateOfBirthday)
	if err != nil {
		return errors.New("Incorrect date format, expected YYYY-MM-DD")
	}

	fmt.Println(t)
	fmt.Println(today)
	// If today date is not "after" the provided date, than error
	if !today.After(t) {
		return errors.New("Invalid Date of Birth, date must be not greater than today")
	}

	return err
}


