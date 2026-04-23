package generic

import (
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func ParseBody[T any](c *fiber.Ctx) (*T, error) {

	var dto T

	if err := c.BodyParser(&dto); err != nil {
		return nil, err
	}

	return &dto, nil

}

func ParseParam[T any](c *fiber.Ctx, key string) (T, error) {

	var zero T

	param := c.Params(key)

	switch any(zero).(type) {

	case int:
		v, err := strconv.Atoi(param)
		return any(v).(T), err

	case uint:
		v, err := strconv.ParseUint(param, 10, 64)
		return any(uint(v)).(T), err

	case string:
		return any(param).(T), nil

	default:
		return zero, fmt.Errorf("unsupported param type")

	}

}

// func ParseQuery(c *fiber.Ctx) (*repositories.Pagination, error) {

// 	var paginationFilter repositories.Pagination

// 	if err := c.QueryParser(&paginationFilter); err != nil {
// 		return nil, err
// 	}

// 	return &paginationFilter, nil

// }
