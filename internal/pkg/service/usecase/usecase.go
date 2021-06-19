package usecase

import (
	"technopark-dbms/internal/pkg/domain"
	"database/sql"
)

type serviceUsecase struct {
	DB *sql.DB
}

func (s *serviceUsecase) Clear() error {
	panic("implement me")
}

func (s *serviceUsecase) Status() (*domain.Service, error) {
	panic("implement me")
}

func NewServiceUsecase(db *sql.DB) domain.ServiceUsecase {
	return &serviceUsecase{
		DB: db,
	}
}

