package usecase

import (
	"context"
	"social-todo-list/common"
	"social-todo-list/module/userlikeitem/model"
)

type UserLikeItemStore interface {
	ListUsers(
		ctx context.Context,
		itemID int,
		paging *common.Paging,
	) ([]common.SimpleUser, error)
	Find(ctx context.Context, userID, itemID int) (*model.Like, error)
	Create(ctx context.Context, data *model.Like) error
	Delete(ctx context.Context, userID, itemID int) error
}

type userLikeItemUseCase struct {
	store UserLikeItemStore
}

func NewUserLikeItemUseCase(store UserLikeItemStore) *userLikeItemUseCase {
	return &userLikeItemUseCase{
		store: store,
	}
}

func (usecase *userLikeItemUseCase) LikeItem(ctx context.Context, data *model.Like) error {
	if err := usecase.store.Create(ctx, data); err != nil {
		return model.ErrCannotLikeItem(err)
	}

	return nil

}

func (usecase *userLikeItemUseCase) UnLikeItem(ctx context.Context, userID, itemID int) error {
	_, err := usecase.store.Find(ctx, userID, itemID)

	if err == common.RecordNotFound {
		return model.ErrCannotDidNotLikeItem(err)
	}
	if err != nil {
		return model.ErrCannotUnLikeItem(err)
	}

	if err := usecase.store.Delete(ctx, userID, itemID); err != nil {
		return model.ErrCannotUnLikeItem(err)
	}

	return nil

}
func (usecase *userLikeItemUseCase) ListUsers(
	ctx context.Context,
	itemID int,
	paging *common.Paging,
) ([]common.SimpleUser, error) {
	data, err := usecase.store.ListUsers(ctx, itemID, paging)
	if err != nil {
		return nil, common.ErrCannotListEntity(model.EntityName, err)
	}

	return data, nil
}
