package main

import (
	"encoding/json"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"io/ioutil"
	"log"
	"net/http"
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

// first speed
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
	return c.String(http.StatusOK, "we got your cat!")
}

// second speed
func addDog(c echo.Context) error {
	dog := Dog{}
	defer c.Request().Body.Close()
	err := json.NewDecoder(c.Request().Body).Decode(&dog)
	if err != nil {
		log.Printf("Failed reading request body for addDog %s", err)
		return c.String(http.StatusInternalServerError, "")
	}
	log.Printf("this is your dog %#v", dog)
	return c.String(http.StatusOK, "we got your dog!")
}

// third speed
func addFiker(c echo.Context) error {
	fiker := Fiker{}
	err := c.Bind(&fiker)
	if err != nil {
		log.Printf("Failed reading request body for addFiker %s", err)
		return c.String(http.StatusInternalServerError, "")
	}
	log.Printf("this is your fiker %#v", fiker)
	return c.String(http.StatusOK, "we got your fiker!")
}
func mainAdmin(c echo.Context) error {
	return c.String(http.StatusOK, " this from secret place of admin")
}

// //////////////////////middlewares///////////////////
func ServerHeader(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Response().Header().Set(echo.HeaderServer, "BluBot/1.0")
		c.Response().Header().Set("notReallyHeader", "thisHasNoMeaning")
		return next(c)
	}
}
func main() {
	e := echo.New()
	e.Use(ServerHeader)
	g := e.Group("/admin")
	g.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: `[${time_rfc3339}],${status},${method},${host},${path},${latency_human}` + "\n",
	}))
	g.Use(middleware.BasicAuth(func(username, password string, ctx echo.Context) (bool, error) {
		if username == "Abemelek" && password == "amen@rophi" {
			return true, nil
		}
		return false, nil
	}))
	g.GET("/main", mainAdmin)
	e.POST("/cats", addCat)
	e.POST("/dogs", addDog)
	e.POST("/fiker", addFiker)
	e.Logger.Fatal(e.Start(":1323"))
}
