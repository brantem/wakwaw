package constant

import (
	"errors"

	"github.com/gofiber/fiber/v2"
)

const AppID = "api"

const (
	StorageKeyPrefix = "api:"
)

var (
	ErrInternalServerError = errors.New("INTERNAL_SERVER_ERROR")
	ErrNotFound            = errors.New("NOT_FOUND")
)

var (
	RespInternalServerError = fiber.Map{"code": "INTERNAL_SERVER_ERROR"}
	RespNotFound            = fiber.Map{"code": "NOT_FOUND"}
)
