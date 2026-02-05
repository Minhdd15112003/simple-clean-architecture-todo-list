package subscriber

import (
	"context"
	"social-todo-list/module/userlikeitem/model"
	"social-todo-list/pubsub"
)

// HasIncreaseLikeCount là interface định nghĩa các method cần thiết
// để tăng/giảm like count của item
// Store nào implement interface này đều có thể được dùng trong subscriber
type HasIncreaseLikeCount interface {
	IncreaseLikeCount(ctx context.Context, id int) error
	DecreaseLikeCount(ctx context.Context, id int) error
}

// IncreaseLikeCountAfterUserLikeItem tạo job xử lý khi user LIKE item
// Job này sẽ được chạy mỗi khi nhận được message từ topic "TopicUserLikeItem"
//
// Flow hoạt động:
// 1. User bấm like → API lưu vào DB user_like_items
// 2. API publish message lên topic "TopicUserLikeItem"
// 3. Subscriber nhận message → job này được trigger
// 4. Job gọi store.IncreaseLikeCount() để tăng liked_count trong bảng items
//
// Params:
//   - ctx: context
//   - store: store có method IncreaseLikeCount (thường là itemStore)
//
// Returns: consumerJob chứa handler để xử lý message
func IncreaseLikeCountAfterUserLikeItem(ctx context.Context, store HasIncreaseLikeCount) consumerJob {
	return consumerJob{
		Title: "Increase like count after user like item",
		Handler: func(ctx context.Context, message *pubsub.Message) error {
			// Parse message data thành model.Like
			// Message.Data() trả về interface{}, cần cast về đúng type
			likeData := message.Data().(*model.Like)

			// Gọi store để tăng like count của item
			return store.IncreaseLikeCount(ctx, likeData.ItemID)
		},
	}
}

// DecreaseLikeCountAfterUserUnlikeItem tạo job xử lý khi user UNLIKE item
// Job này sẽ được chạy mỗi khi nhận được message từ topic "TopicUserDislikeItem"
//
// Flow hoạt động:
// 1. User bấm unlike → API xóa khỏi DB user_like_items
// 2. API publish message lên topic "TopicUserDislikeItem"
// 3. Subscriber nhận message → job này được trigger
// 4. Job gọi store.DecreaseLikeCount() để giảm liked_count trong bảng items
//
// Params:
//   - ctx: context
//   - store: store có method DecreaseLikeCount (thường là itemStore)
//
// Returns: consumerJob chứa handler để xử lý message
func DecreaseLikeCountAfterUserUnlikeItem(ctx context.Context, store HasIncreaseLikeCount) consumerJob {
	return consumerJob{
		Title: "Decrease like count after user unlike item",
		Handler: func(ctx context.Context, message *pubsub.Message) error {
			// Parse message data thành model.Like
			likeData := message.Data().(*model.Like)

			// Gọi store để giảm like count của item
			return store.DecreaseLikeCount(ctx, likeData.ItemID)
		},
	}
}
