package usecase

import (
	"context"
	"log"
	"social-todo-list/common"
	"social-todo-list/middleware"
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
	GetLikeItem(ctx context.Context, ids []int) (map[int]int, error)
}

type InDeCreaseItemStore interface {
	IncreaseLikeCount(ctx context.Context, id int) error
	DecreaseLikeCount(ctx context.Context, id int) error
}

type userLikeItemUseCase struct {
	store     UserLikeItemStore
	itemStore InDeCreaseItemStore
}

func NewUserLikeItemUseCase(store UserLikeItemStore, itemStore InDeCreaseItemStore) *userLikeItemUseCase {
	return &userLikeItemUseCase{
		store:     store,
		itemStore: itemStore,
	}
}

func (usecase *userLikeItemUseCase) LikeItem(ctx context.Context, data *model.Like) error {
	if err := usecase.store.Create(ctx, data); err != nil {
		return model.ErrCannotLikeItem(err)
	}

	go func() {
		defer middleware.RecoverGoroutine()
		if err := usecase.itemStore.IncreaseLikeCount(ctx, data.ItemID); err != nil {
			log.Println(err)
		}
	}()

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

	go func() {
		defer middleware.RecoverGoroutine()
		if err := usecase.itemStore.DecreaseLikeCount(ctx, itemID); err != nil {
			log.Println(err)
		}
	}()

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
