package subscriber

import (
	"context"
	"log"
	"social-todo-list/common"
	"social-todo-list/pubsub"
)

// consumerJob đại diện cho một công việc xử lý message từ pubsub
// Mỗi job có:
// - Title: tên job để dễ debug
// - Handler: hàm xử lý message khi nhận được từ topic
type consumerJob struct {
	Title   string
	Handler func(ctx context.Context, message *pubsub.Message) error
}

// consumerEngine là engine quản lý tất cả các subscribers
// Nó giữ reference đến pubsub để subscribe các topic
type consumerEngine struct {
	subTopic pubsub.PubSup // PubSub instance để subscribe topics
}

// NewEngine tạo một consumer engine mới
// Params:
//   - subTopic: pubsub instance để subscribe topics
func NewEngine(subTopic pubsub.PubSup) *consumerEngine {
	return &consumerEngine{subTopic: subTopic}
}

// Start khởi động tất cả subscribers cho các topics
// Đây là nơi đăng ký tất cả các job xử lý khi có message mới
// Params:
//   - itemStore: store có method IncreaseLikeCount và DecreaseLikeCount
func (ce *consumerEngine) Start(itemStore HasIncreaseLikeCount) error {
	// Subscribe topic "user like item" - khi user bấm like
	// isConcurrent = false: xử lý tuần tự, đảm bảo thứ tự
	ce.startSubTopic(common.TopicUserLikeItem, false, IncreaseLikeCountAfterUserLikeItem(context.Background(), itemStore))

	// Subscribe topic "user dislike item" - khi user bấm unlike
	// isConcurrent = false: xử lý tuần tự, đảm bảo thứ tự
	ce.startSubTopic(common.TopicUserDislikeItem, false, DecreaseLikeCountAfterUserUnlikeItem(context.Background(), itemStore))

	return nil
}

// startSubTopic subscribe một topic và setup xử lý message
// Đây là core function để lắng nghe message từ pubsub
// Params:
//   - topic: tên topic cần subscribe (vd: "TopicUserLikeItem")
//   - isConcurrent:
//   - true: xử lý message song song (concurrent) - nhanh nhưng có thể mất thứ tự
//   - false: xử lý message tuần tự (sequential) - chậm hơn nhưng đảm bảo thứ tự
//   - consumerJobs: danh sách các job cần chạy khi nhận được message
func (ce *consumerEngine) startSubTopic(
	topic pubsub.Topic,
	isConcurrent bool,
	consumerJobs ...consumerJob,
) error {
	// BƯỚC 1: Subscribe topic - tạo một channel để nhận message từ topic này
	// c là read-only channel (<-chan), chỉ dùng để đọc message
	c, _ := ce.subTopic.Subscribe(context.Background(), topic)

	// BƯỚC 2: In ra log các job đã được setup
	for _, item := range consumerJobs {
		log.Println("Setup consumer for:", item.Title)
	}

	// BƯỚC 3: Chạy goroutine để lắng nghe message liên tục
	// ⚠️ CHÚ Ý: Goroutine này (dòng 72) luôn tồn tại và chạy cho TẤT CẢ messages
	// Đây là goroutine CHÍNH để lắng nghe topic
	go func() {
		// Vòng for chạy mãi nhưng nó BỊ BLOCK (ngủ) ở dòng 87
		//
		// Flow hoạt động:
		// 1. Goroutine chạy đến dòng 87: msg := <-c
		// 2. DỪNG LẠI (block) - không tiêu tốn CPU - đợi message
		// 3. Khi có message mới → goroutine được "đánh thức"
		// 4. Xử lý message xong → quay lại dòng 87 → lại bị block (ngủ tiếp)
		//
		// 🎯 Đây là cơ chế "event-driven", KHÔNG phải "polling" (kiểm tra liên tục)
		for {
			// ⏸️ BLOCKING POINT: Goroutine ngủ tại đây, đợi message
			// Không tốn CPU, không chạy vòng lặp liên tục kiểm tra!
			// Chỉ được đánh thức khi channel c nhận được message mới
			msg := <-c
			log.Println("Received message from topic:", msg.Channel())

			// Xử lý message với TẤT CẢ các jobs đã đăng ký
			for i := range consumerJobs {
				job := consumerJobs[i]

				if isConcurrent {
					// Topic "UserLikeItem"
					// └─ 1 Goroutine chính (dòng 72) - lắng nghe message
					//    │
					//    Nhận Message 1:
					//    ├─ Goroutine con 1: Job A (chạy riêng)
					//    └─ Goroutine con 2: Job B (chạy riêng)

					//    Nhận Message 2:
					//    ├─ Goroutine con 3: Job A (chạy riêng)
					//    └─ Goroutine con 4: Job B (chạy riêng)
					go func(j consumerJob) {
						if err := j.Handler(context.Background(), msg); err != nil {
							log.Println("Error when running consumer job:", err)
						}
					}(job)
				} else {
					// Topic "UserLikeItem"
					// └─ 1 Goroutine chính (dòng 72)
					//    └─ Tất cả jobs chạy tuần tự TRONG goroutine này
					//       ├─ Message 1 → Job A xong → Job B xong
					//       ├─ Message 2 → Job A xong → Job B xong
					//       └─ Message 3 → ...
					if err := job.Handler(context.Background(), msg); err != nil {
						log.Println("Error when running consumer job:", err)
					}
				}
			}
		}
	}()

	return nil
}
