package handler

import (
	"context"
	"social-todo-list/common"
	"social-todo-list/module/userlikeitem/model"
)

type LikeUseCase interface {
	LikeItem(ctx context.Context, data *model.Like) error
	UnLikeItem(ctx context.Context, userID, itemID int) error
	ListUsers(
		ctx context.Context,
		itemID int,
		paging *common.Paging,
	) ([]common.SimpleUser, error)
}

type LikeService struct {
	useCase LikeUseCase
}

func NewLikeService(useCase LikeUseCase) *LikeService {
	return &LikeService{
		useCase: useCase, //Context giữa UseCase và Service
	}
}

func (s *LikeService) LikeItem(ctx context.Context, data *model.Like) error {
	return s.useCase.LikeItem(ctx, data)
}

func (s *LikeService) UnLikeItem(ctx context.Context, userID, itemID int) error {
	return s.useCase.UnLikeItem(ctx, userID, itemID)
}

func (s *LikeService) ListUsers(
	ctx context.Context,
	itemID int,
	paging *common.Paging,
) ([]common.SimpleUser, error) {
	return s.useCase.ListUsers(ctx, itemID, paging)
}
