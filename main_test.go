package main

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"hbday/web-service-gin/models"
	"net/http"
	"net/http/httptest"
	"time"

	"testing"
)
const username = "userTestg"
var (
	validPerson = &models.User{DateOfBirth: "1990-01-01"}
	invalidBrithday = &models.User{DateOfBirth: "199-0-01"}
	router *gin.Engine
)

func init(){
	ConnectDb()
	deleteUser(username)
	// initialize the router
	router = gin.Default()
	router.PUT("/hello/:username", postPerson)
	router.GET("/hello/:username", getPersonByUsername)
}

func TestPutUserSuccess(t *testing.T) {
	// initialize http client
	w := httptest.NewRecorder()

	json, err := json.Marshal(validPerson)
	if err != nil {
		panic(err)
	}
	req, err := http.NewRequest(http.MethodPut, "/hello/"+username, bytes.NewBuffer(json))
	if err != nil {
		panic(err)
	}

	router.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Fail()
	}
}

func TestPutUserFailed(t *testing.T) {
	w := httptest.NewRecorder()

	json, err := json.Marshal(validPerson)
	if err != nil {
		panic(err)
	}
	req, err := http.NewRequest(http.MethodPut, "/hello/"+username, bytes.NewBuffer(json))
	if err != nil {
		panic(err)
	}

	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fail()
	}
}

func TestPostUserInvalidBirthday(t *testing.T) {
	// initialize http client
	w := httptest.NewRecorder()

	json, err := json.Marshal(invalidBrithday)
	if err != nil {
		panic(err)
	}
	req, err := http.NewRequest(http.MethodPut, "/hello/invalid"+username, bytes.NewBuffer(json))
	if err != nil {
		panic(err)
	}

	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fail()
	}
}

func TestGetExistingUser(t *testing.T) {
	// initialize http client
	w := httptest.NewRecorder()

	req, err := http.NewRequest(http.MethodGet, "/hello/"+username, nil)
	if err != nil {
		panic(err)
	}

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK{
		t.Fail()
	}
}

func TestGetNonExistingUser(t *testing.T) {
	// initialize http client
	w := httptest.NewRecorder()

	req, err := http.NewRequest(http.MethodGet, "/hello/random"+username, nil)
	if err != nil {
		panic(err)
	}

	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound{
		t.Fail()
	}
}

func TestDaysToBirthday(t *testing.T) {
	today := time.Now()
	days, _ := daysToBirthday(today.Format("2006-01-02"))
	if days != 0{
		t.Fail()
	}

	birthday := time.Date(today.Year(), today.Month(), today.Day()-1, 23, 59, 59, 0, time.UTC)
	days,_ = daysToBirthday(birthday.Format("2006-01-02"))
	if days < 364 {
		t.Fail()
	}
}

func deleteUser(user string){
	_, err := models.DB.Query("DELETE FROM people WHERE username = $1", user)
	if err != nil {
		panic(err)
	}
}
