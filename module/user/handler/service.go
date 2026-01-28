package handler

import (
	"context"
	"social-todo-list/components/tokenprovider"
	"social-todo-list/module/user/model"
)

type IUserUseCase interface {
	Register(ctx context.Context, data *model.UserCreation) error
	Login(ctx context.Context, data *model.UserLogin) (tokenprovider.Token, error)
	Profile(ctx context.Context) (*model.User, error)
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

func (s *IUserService) Login(ctx context.Context, data *model.UserLogin) (tokenprovider.Token, error) {
	return s.useCase.Login(ctx, data)
}

func (s *IUserService) Profile(ctx context.Context) (*model.User, error) {
	return s.useCase.Profile(ctx)
}
