package usecase

import (
	"context"
	"log"
	"social-todo-list/common"
	"social-todo-list/module/userlikeitem/model"
	"social-todo-list/pubsub"
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
	store UserLikeItemStore
	// itemStore InDeCreaseItemStore
	ps pubsub.PubSup
}

func NewUserLikeItemUseCase(store UserLikeItemStore, ps pubsub.PubSup) *userLikeItemUseCase {
	return &userLikeItemUseCase{
		store: store,
		// itemStore: itemStore,
		ps: ps,
	}
}

func (usecase *userLikeItemUseCase) LikeItem(ctx context.Context, data *model.Like) error {
	if err := usecase.store.Create(ctx, data); err != nil {
		return model.ErrCannotLikeItem(err)
	}

	if err := usecase.ps.Publish(ctx, common.TopicUserLikeItem, pubsub.NewMessage(data)); err != nil {
		log.Println(err)
	}
	// job := asyncjob.NewJob(func(ctx context.Context) error {
	// 	if err := usecase.itemStore.IncreaseLikeCount(ctx, data.ItemID); err != nil {
	// 		return err
	// 	}
	// 	return nil
	// })
	// if err := asyncjob.NewGroup(true, job).Run(ctx); err != nil {
	// 	log.Println(err)
	// }

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

	data := &model.Like{
		UserID: userID,
		ItemID: itemID,
	}

	if err := usecase.ps.Publish(ctx, common.TopicUserDislikeItem, pubsub.NewMessage(data)); err != nil {
		log.Println(err)
	}

	// job := asyncjob.NewJob(func(ctx context.Context) error {
	// 	if err := usecase.itemStore.DecreaseLikeCount(ctx, itemID); err != nil {
	// 		return err
	// 	}
	// 	return nil
	// })
	// if err := asyncjob.NewGroup(true, job).Run(ctx); err != nil {
	// 	log.Println(err)
	// }

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
