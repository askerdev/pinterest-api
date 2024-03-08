package errs

import (
	"encoding/json"
)

type FormError struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
	Errors  ErrMap `json:"errors"`
}

type ErrMap map[string]string

func (m ErrMap) Error() string {
	v, _ := json.Marshal(m)
	return string(v)
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
