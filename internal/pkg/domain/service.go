package domain

type Service struct {
	User   int32
	Forum  int32
	Thread int32
	Post   int64
}

type ServiceUsecase interface {
	Clear() error
	Status() (*Service, error)
}
