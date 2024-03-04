package structure

import (
	"errors"
	"gorm.io/gorm"
)

func NewRecordSkipErrNotFound[T any](err error) (*T, error) {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return new(T), nil
	}

	return nil, err
}

func RecordBuildSkipErrNotFound[M any, T any](m *M, transform func(m *M) *T, err error) (*T, error) {
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return new(T), nil
		}

		return nil, err
	}

	return transform(m), err
}
