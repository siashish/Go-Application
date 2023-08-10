package main

import (
	"net/http/httptest"
	"siashish/application/routes"
	"strings"
	"testing"

	"github.com/go-playground/assert/v2"
	"github.com/gofiber/fiber/v2"
)

func TestHome(t *testing.T) {
	app := fiber.New()

	//routes
	routes.UserRoute(app)

	r := httptest.NewRequest("GET", "/", nil)
	resp, _ := app.Test(r, -1)
	assert.Equal(t, 200, resp.StatusCode)
}

func TestCreateUser(t *testing.T) {
	app := fiber.New()

	//routes
	routes.UserRoute(app)

	payload := strings.NewReader(`{
		"username": "singh1",
		"expiry_date": 1648867200,
		"outputs": ["hls"],
		"password": "mypassword"
		}`)

	r := httptest.NewRequest("POST", "/user", payload)
	r.Header.Add("Content-Type", "application/json")
	resp, _ := app.Test(r, -1)
	assert.Equal(t, 201, resp.StatusCode)
}

func TestGetAUser(t *testing.T) {
	app := fiber.New()

	//routes
	routes.UserRoute(app)

	// success
	r := httptest.NewRequest("GET", "/user/ashish", nil)
	resp, _ := app.Test(r, -1)
	assert.Equal(t, 200, resp.StatusCode)
}

func TestEditAUser(t *testing.T) {
	app := fiber.New()

	//routes
	routes.UserRoute(app)

	payload := strings.NewReader(`{
		"username": "ashsingh"
	}`)

	r := httptest.NewRequest("PATCH", "/user/singh", payload)
	r.Header.Add("Content-Type", "application/json")
	resp, _ := app.Test(r, -1)
	assert.Equal(t, 200, resp.StatusCode)
}

func TestDeleteAUser(t *testing.T) {
	app := fiber.New()

	//routes
	routes.UserRoute(app)

	// success
	r := httptest.NewRequest("DELETE", "/users/ashish3", nil)
	resp, _ := app.Test(r, -1)
	assert.Equal(t, 200, resp.StatusCode)
}

func TestGetAllUsers(t *testing.T) {
	app := fiber.New()

	//routes
	routes.UserRoute(app)

	// success
	r := httptest.NewRequest("GET", "/users", nil)
	resp, _ := app.Test(r, -1)
	assert.Equal(t, 200, resp.StatusCode)
}
