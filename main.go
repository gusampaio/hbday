package main

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"hbday/web-service-gin/models"
	"log"
	"net/http"
	"time"
)

func main() {
	// Initalize the sql.DB connection pool and assign it to the models.DB
	// global variable.
	var err error

	models.DB, err = sql.Open("postgres", "postgres://gusampaio:gusampaio_pass@localhost/hbday_db")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("table company is created")

	router := gin.Default()
	router.GET("/hello/:username", getPersonByUsername)
	router.GET("/hello/all", getAll)
	router.PUT("/hello/:username", postPerson)
	router.Run("localhost:8080")
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
	//username := c.Param("username")

	// Loop through the list of people, looking for
	// an person whose username value matches the parameter.
/*	for _, a := range people {
		if a.Username == username {
			c.IndentedJSON(http.StatusOK, a)
			return
		}
	}*/
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "username not found"})
}

// postPerson adds a person from JSON received in the request body.
func postPerson(c *gin.Context) {
	var newPerson models.Person

	// Call BindJSON to bind the received JSON to
	// newPerson.
	if err := c.BindJSON(&newPerson); err != nil {
		return
	}

	err := isValidDoB(newPerson.DateOfBirth)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	// Including username param as not being passed in the JSON
	newPerson.Username = c.Param("username")

	err = models.SetNewPerson(newPerson)
	if err != nil {
		panic(err)
	}
	// Add the new person to the slice.
	c.IndentedJSON(http.StatusNoContent, newPerson)
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


