package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http/httptest"
	"testing"

	"example.com/api/model"
	"example.com/shared/testutil"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

var (
	userColumns = []string{
		"id",
		"name",
	}
	user = model.User{
		ID:   "user-1",
		Name: "John Doe",
	}
)

func TestGetUser(t *testing.T) {
	assert := assert.New(t)

	t.Run("not found", func(t *testing.T) {
		db, mock := testutil.NewDB()
		redis := testutil.Redis{}
		h := New(nil, db, &redis)

		mock.ExpectQuery(`SELECT .+ FROM users`).
			WithArgs("user-1").
			WillReturnError(sql.ErrNoRows)

		app := fiber.New()
		h.Register(app)

		req := httptest.NewRequest("GET", "/v1/users/user-1", nil)
		resp, _ := app.Test(req, -1)
		assert.Equal(fiber.StatusNotFound, resp.StatusCode)
		body, _ := io.ReadAll(resp.Body)
		assert.Equal(`{"user":null,"error":{"code":"NOT_FOUND"}}`, string(body))
	})

	t.Run("success", func(t *testing.T) {
		db, mock := testutil.NewDB()
		redis := testutil.Redis{}
		h := New(nil, db, &redis)

		mock.ExpectQuery(`SELECT .+ FROM users`).
			WithArgs("user-1").
			WillReturnRows(sqlmock.NewRows(userColumns).AddRow(user.ID, user.Name))

		app := fiber.New()
		h.Register(app)

		req := httptest.NewRequest("GET", "/v1/users/user-1", nil)
		resp, _ := app.Test(req, -1)
		assert.Equal(fiber.StatusOK, resp.StatusCode)
		body, _ := io.ReadAll(resp.Body)
		buf, _ := json.Marshal(user)
		assert.Equal(fmt.Sprintf(`{"user":%s,"error":null}`, string(buf)), string(body))
	})
}

func TestPokeUser(t *testing.T) {
	redis := testutil.Redis{}
	h := New(nil, nil, &redis)

	app := fiber.New()
	h.Register(app)

	req := httptest.NewRequest("POST", "/v1/users/user-1/poke", nil)
	resp, _ := app.Test(req, -1)
	assert.Equal(t, "user.poke", redis.PublishChannel[0])
	assert.Equal(t, "user-1", redis.PublishMessage[0])
	assert.Equal(t, fiber.StatusAccepted, resp.StatusCode)
}
