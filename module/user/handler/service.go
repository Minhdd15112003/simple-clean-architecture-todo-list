package handler

import (
	"context"
	"social-todo-list/module/user/model"
)

type IUserUseCase interface {
	Register(ctx context.Context, data *model.UserCreation) error
}

type IUserService struct {
	useCase IUserUseCase
}

func NewUserService(useCase IUserUseCase) *IUserService {
	return &IUserService{
		useCase: useCase,
	}
}

func (s *IUserService) Register(ctx context.Context, data *model.UserCreation) error {
	return s.useCase.Register(ctx, data)
}
