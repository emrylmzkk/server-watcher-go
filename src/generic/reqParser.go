package generic

import (
	"fmt"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

func ParseBody[T any](c *fiber.Ctx) (*T, error) {

	var dto T

	// json --> struct

	if err := c.BodyParser(&dto); err != nil {
		return nil, err
	}

	// validator kutuphanesi ile struct validasyonu

	if err := validate.Struct(dto); err != nil {
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
