package usecase

import (
	"github.com/jackc/pgx"
	"technopark-dbms/internal/pkg/domain"
)

type serviceUsecase struct {
	DB *pgx.ConnPool
}

func (s *serviceUsecase) Clear() error {
	panic("implement me")
}

func (s *serviceUsecase) Status() (*domain.Service, error) {
	panic("implement me")
}

func NewServiceUsecase(db *pgx.ConnPool) domain.ServiceUsecase {
	return &serviceUsecase{
		DB: db,
	}
}
