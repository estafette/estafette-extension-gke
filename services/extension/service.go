package extension

import "context"

//go:generate mockgen -package=extension -destination ./mock.go -source=service.go
type Service interface {
	DryRun(ctx context.Context) (err error)
	Run(ctx context.Context) (err error)
}

// NewService returns a new extension.Service
func NewService(ctx context.Context) (Service, error) {
	return &service{}, nil
}

type service struct {
}

func (s *service) DryRun(ctx context.Context) (err error) {
	return nil
}

func (s *service) Run(ctx context.Context) (err error) {
	return nil
}