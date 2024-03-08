package errs

import (
	"encoding/json"
	"github.com/gofiber/fiber/v3"
)

type ErrMap map[string]string

func (m *ErrMap) ToFiberError(status int) error {
	v, _ := json.Marshal(*m)
	return fiber.NewError(status, string(v))
}

func (m *ErrMap) Set(key, value string) {
	if *m == nil {
		*m = ErrMap{}
	}

	(*m)[key] = value
}

func (m *ErrMap) Has(key string) bool {
	_, ok := (*m)[key]
	return ok
}

func (m *ErrMap) Remove(key string) {
	if !m.Has(key) {
		return
	}
	delete(*m, key)
}
