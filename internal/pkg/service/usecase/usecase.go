package usecase

import (
	"github.com/jackc/pgx"
	"technopark-dbms/internal/pkg/domain"
)

type serviceUsecase struct {
	DB *pgx.ConnPool
}

func (s *serviceUsecase) Clear() error {
	query := "truncate forums, users, f_u, posts, threads, votes;"
	_, err := s.DB.Exec(query)
	if err != nil {
		return err
	}
	return nil
}

func (s *serviceUsecase) Status() (*domain.Service, error) {
	query := "select (select count(*) from users), (select count(*) from forums),  (select count(*) from threads), (select count(*) from posts);"
	res := &domain.Service{}
	err := s.DB.QueryRow(query).Scan(&res.User, &res.Forum, &res.Thread, &res.Post)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func NewServiceUsecase(db *pgx.ConnPool) domain.ServiceUsecase {
	return &serviceUsecase{
		DB: db,
	}
}
