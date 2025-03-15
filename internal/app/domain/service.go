package domain

type Service struct {
	Name   string
	Status Status
}

func NewService(
	name string,
	status Status,
) Service {
	return Service{
		Name:   name,
		Status: status,
	}
}
