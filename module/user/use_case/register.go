package usecase

import (
	"context"
	"social-todo-list/common"
	"social-todo-list/module/user/model"
)

type RegisterStorage interface {
	FindUser(ctx context.Context, conditions map[string]interface{}, moreInfor ...string) (*model.User, error)
	CreateUser(ctx context.Context, data *model.UserCreation) error
}

type Hasher interface {
	Hash(data string) string
}

type registerUseCase struct {
	registerStorage RegisterStorage
	hasher          Hasher
}

func NewRegisterUseCase(registerStorage RegisterStorage, hasher Hasher) *registerUseCase {
	return &registerUseCase{
		registerStorage: registerStorage,
		hasher:          hasher,
	}
}

func (useCase *registerUseCase) Register(ctx context.Context, data *model.UserCreation) error {
	user, _ := useCase.registerStorage.FindUser(ctx, map[string]interface{}{
		"email": data.Email,
	})

	if user != nil {
		return model.ErrEmailExisted
	}

	salt := common.GenSalt()

	data.Password = useCase.hasher.Hash(data.Password + salt)
	data.Salt = salt
	data.Role = "user"
	if err := useCase.registerStorage.CreateUser(ctx, data); err != nil {
		return common.ErrCannotCreateEntity(model.EntityName, err)
	}

	return nil
}
