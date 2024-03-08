package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
)
type (
	Cat struct {
		Name string `json:"name"`
		Type string `json:"type"`
	}
	Dog struct {
		Name string `json:"name"`
		Type string `json:"type"`
	}
	Fiker struct {
		Name string `json:"name"`
		Type string `json:"type"`
	}
)
//first speed
func addCat(c echo.Context) error {
	cat := Cat{}
	defer c.Request().Body.Close()
	b, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
	  log.Printf("Failed reading request body for addcart %s", err)
	  return c.String(http.StatusInternalServerError, "")
	}
	err = json.Unmarshal(b, &cat)
	if err != nil {
	  log.Printf("Failed unmarshaling json data for addcart %s", err)
	  return c.String(http.StatusInternalServerError, "")
	}
	log.Printf("this is your cat %#v", cat)
	return c.String(http.StatusOK,"we got your cat!")
  }
  //second speed
  func addDog(c echo.Context)error{
	dog := Dog{}
	defer c.Request().Body.Close()
	err := json.NewDecoder(c.Request().Body).Decode(&dog)
	if err != nil{
	  log.Printf("Failed reading request body for addDog %s", err)
	  return c.String(http.StatusInternalServerError, "")
	}
	log.Printf("this is your dog %#v", dog)
	return c.String(http.StatusOK,"we got your dog!")
  }
  //third speed
  func addFiker(c echo.Context)error{
	fiker := Fiker{}
	err := c.Bind(&fiker)
	if err != nil{
	  log.Printf("Failed reading request body for addFiker %s", err)
	  return c.String(http.StatusInternalServerError, "")
	}
	log.Printf("this is your fiker %#v", fiker)
	return c.String(http.StatusOK,"we got your fiker!")
  }

func main() {
	e := echo.New()
	e.POST("/cats",addCat)
	e.POST("/dogs",addDog)
	e.POST("/fiker",addFiker)
	e.Logger.Fatal(e.Start(":1323"))
}
