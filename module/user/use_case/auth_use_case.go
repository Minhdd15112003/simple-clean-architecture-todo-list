package usecase

import (
	"context"
	"social-todo-list/common"
	"social-todo-list/components/tokenprovider"
	"social-todo-list/module/user/model"
)

type AuthStorage interface {
	FindUser(ctx context.Context, conditions map[string]interface{}, moreInfor ...string) (*model.User, error)
	CreateUser(ctx context.Context, data *model.UserCreation) error
}

type authUseCase struct {
	authStorage   AuthStorage
	tokenprovider tokenprovider.Provider
	expiry        int
}

func NewAuthUseCase(authStorage AuthStorage, tokenprovider tokenprovider.Provider, expiry int) *authUseCase {
	return &authUseCase{
		authStorage:   authStorage,
		tokenprovider: tokenprovider,
		expiry:        expiry,
	}
}

func (useCase *authUseCase) Register(ctx context.Context, data *model.UserCreation) error {
	user, _ := useCase.authStorage.FindUser(ctx, map[string]interface{}{
		"email": data.Email,
	})

	if user != nil {
		return model.ErrEmailExisted
	}

	hashPassword, err := common.HashPassword(data.Password)
	if err != nil {
		return common.ErrCannotCreateEntity(model.EntityName, err)
	}

	data.Password = hashPassword
	data.Role = "user"
	if err := useCase.authStorage.CreateUser(ctx, data); err != nil {
		return common.ErrCannotCreateEntity(model.EntityName, err)
	}

	return nil
}

func (usecase *authUseCase) Login(ctx context.Context, data *model.UserLogin) (tokenprovider.Token, error) {
	user, err := usecase.authStorage.FindUser(ctx, map[string]interface{}{
		"email": data.Email,
	})
	if err != nil {
		return nil, model.ErrEmailOrPasswordInvalid
	}

	if err := common.CheckPassword(user.Password, data.Password); err != nil {
		return nil, common.ErrInternal(err)
	}

	payload := &common.TokenPayload{
		UId:   user.Id,
		URole: user.Role.String(),
	}

	accessToken, err := usecase.tokenprovider.Gernerate(payload, usecase.expiry)
	if err != nil {
		return nil, common.ErrInternal(err)
	}

	return accessToken, nil

}
