package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
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
	JwtClaims struct {
		Name string `json:"name"`
		jwt.StandardClaims
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

func mainCookie(c echo.Context) error {
	return c.String(http.StatusOK, "this from cookie main page")
}
func login(c echo.Context) error {
	username := c.QueryParam("username")
	password := c.QueryParam("password")
	if username == "Abemelek" && password == "change@be" {
		cookie := &http.Cookie{}
		cookie.Name = "authorization_cookie"
		cookie.Value = "some_string"
		cookie.Expires = time.Now().Add(48 * time.Hour)
		c.SetCookie(cookie)
		// TODO : create jwt token
		token, err := createJwtToken()
		if err != nil {
			log.Println("Error while creating token", err)
			return c.String(http.StatusInternalServerError, "something went wrong")
		}
		return c.JSON(http.StatusOK, map[string]string{
			"message": "you were login",
			"token":   token,
		})
	}
	return c.String(http.StatusUnauthorized, "your username or password were wrong!")
}

func createJwtToken() (string, error) {
	claims := JwtClaims{
		"Abemelek",
		jwt.StandardClaims{
			Id:        "main_user_id",
			ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
		},
	}
	rawToken := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	token, err := rawToken.SignedString([]byte("RXreUwEbzNBDfszcUIasj7kOYtUorEfD"))
	if err != nil {
		return "", err
	}
	return token, nil
}
func mainJwt(c echo.Context) error {
	// if user, ok := c.Get("user").(*jwt.Token); ok {
	// 	// user is now a pointer to a jwt.Token, proceed with claims access
	// 	claims := user.Claims.(jwt.MapClaims)
	// 	// ... rest of your code
	// 	log.Println("User name : ", claims["name"], "User ID : ", claims["jti"])
	// 	} else {
	// 		// Handle the case where user is not a jwt.Token
	// 		return echo.ErrUnauthorized
	// 	}
	return c.String(http.StatusOK, "your on top secret jwt page")
}

// //////////////////////middlewares///////////////////
func ServerHeader(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Response().Header().Set(echo.HeaderServer, "BluBot/1.0")
		c.Response().Header().Set("notReallyHeader", "thisHasNoMeaning")
		return next(c)
	}
}
func checkCookie(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cookie, err := c.Cookie("authorization_cookie")
		if err != nil {
			log.Println(err)
			return err
		}
		if cookie.Value == "some_string" {
			return next(c)
		}
		return c.String(http.StatusUnauthorized, "you don't have the right cookie , cookie")
	}
}
func main() {
	e := echo.New()
	e.Use(ServerHeader)
	adminGroup := e.Group("/admin")
	cookieGroup := e.Group("/cookie")
	jwtGroup := e.Group("/jwt")
	adminGroup.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: `[${time_rfc3339}],${status},${method},${host},${path},${latency_human}` + "\n",
	}))
	adminGroup.Use(middleware.BasicAuth(func(username, password string, ctx echo.Context) (bool, error) {
		if username == "Abemelek" && password == "amen@rophi" {
			return true, nil
		}
		return false, nil
	}))
	jwtGroup.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey:    []byte("RXreUwEbzNBDfszcUIasj7kOYtUorEfD"),
		SigningMethod: "HS512",
	}))

	cookieGroup.Use(checkCookie)
	cookieGroup.GET("/main", mainCookie)

	adminGroup.GET("/main", mainAdmin)

	jwtGroup.GET("/main", mainJwt)
	e.GET("/login", login)
	e.POST("/cats", addCat)
	e.POST("/dogs", addDog)
	e.POST("/fiker", addFiker)
	e.Logger.Fatal(e.Start(":1323"))
}
