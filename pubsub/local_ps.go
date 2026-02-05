package pubsub

import (
	"context"
	"log"
	"social-todo-list/middleware"
	"sync"
)

// localPubSub là implementation cục bộ của PubSub pattern
// Sử dụng channels để giao tiếp giữa publishers và subscribers
//
// Kiến trúc:
// Publisher → messageQueue (buffer 10000) → run() goroutine → mapChannel → Subscribers
//
// messageQueue: Hàng đợi trung tâm chứa tất cả messages từ mọi topics
// mapChannel: Map từ topic → danh sách channels của subscribers
// locker: RWMutex để đồng bộ truy cập vào mapChannel
type localPubSub struct {
	messageQueue chan *Message           // Queue trung tâm, buffer 10000 messages
	mapChannel   map[Topic][]chan *Message // Map: topic → list of subscriber channels
	locker       *sync.RWMutex            // Lock để bảo vệ mapChannel khỏi race condition
}

// NewPubSub tạo một PubSub instance mới và khởi động background goroutine
//
// Khởi tạo:
// 1. Tạo messageQueue với buffer 10000 (có thể chứa tối đa 10000 messages chưa xử lý)
// 2. Tạo mapChannel rỗng (sẽ được populate khi có subscribers)
// 3. Tạo RWMutex để đồng bộ hóa
// 4. Gọi run() để start goroutine phân phối messages
//
// Returns: PubSub instance đã sẵn sàng sử dụng
func NewPubSub() *localPubSub {
	pb := &localPubSub{
		messageQueue: make(chan *Message, 10000), // Buffer lớn để tránh block publisher
		mapChannel:   make(map[Topic][]chan *Message),
		locker:       new(sync.RWMutex),
	}

	// Start goroutine để phân phối messages từ queue đến subscribers
	pb.run()

	return pb
}

// Publish gửi một message vào topic
//
// Flow:
// 1. Set topic cho message
// 2. Tạo goroutine mới để gửi message (non-blocking cho caller)
// 3. Gửi message vào messageQueue
// 4. Log message đã được publish
//
// ⚠️ CHÚ Ý:
// - Method này KHÔNG BLOCK, trả về ngay lập tức
// - Message được gửi trong goroutine riêng
// - Nếu queue đầy (10000 messages), goroutine sẽ block ở bước 3
//
// Params:
//   - ctx: context (hiện tại chưa dùng)
//   - topic: tên topic để publish
//   - data: message cần gửi
//
// Returns: nil (luôn thành công vì xử lý async)
func (ps *localPubSub) Publish(ctx context.Context, topic Topic, data *Message) error {
	// Gán topic vào message
	data.SetChannel(topic)

	// Chạy trong goroutine riêng để không block API caller
	go func() {
		defer middleware.RecoverGoroutine() // Recovery để tránh crash

		// Gửi vào queue trung tâm
		// ⚠️ Có thể block nếu queue đầy (10000 messages)
		ps.messageQueue <- data

		log.Println("New message published", data.String())
	}()

	return nil // Trả về ngay, không đợi message được xử lý
}

// Subscribe đăng ký nhận messages từ một topic
//
// Flow:
// 1. Tạo channel mới cho subscriber này
// 2. Lock mapChannel để tránh race condition
// 3. Thêm channel vào danh sách subscribers của topic
// 4. Unlock
// 5. Trả về channel (để đọc messages) và hàm unsubscribe
//
// Cấu trúc mapChannel:
// {
//   "TopicUserLikeItem": [channel1, channel2, channel3],
//   "TopicUserDislikeItem": [channel4]
// }
//
// Params:
//   - ctx: context (hiện tại chưa dùng)
//   - topic: topic muốn subscribe
//
// Returns:
//   - ch: read-only channel để nhận messages
//   - unsubscribe: hàm để hủy subscription
func (ps *localPubSub) Subscribe(ctx context.Context, topic Topic) (ch <-chan *Message, unsubscribe func()) {
	// Tạo channel mới cho subscriber này (unbuffered - sẽ block nếu không đọc)
	c := make(chan *Message)

	// CRITICAL SECTION: Modify shared mapChannel
	ps.locker.Lock()

	// Kiểm tra topic đã tồn tại trong map chưa
	if val, ok := ps.mapChannel[topic]; ok {
		// Topic đã có subscribers → append channel mới vào list
		ps.mapChannel[topic] = append(val, c)
	} else {
		// Topic mới → tạo list mới với 1 channel
		ps.mapChannel[topic] = []chan *Message{c}
	}
	ps.locker.Unlock()

	// Trả về channel (read-only) và hàm unsubscribe
	return c, func() {
		// unsubscribe: hàm closure để hủy subscription
		// Khi gọi, sẽ:
		// 1. Xóa channel khỏi danh sách subscribers của topic
		// 2. Close channel để báo cho subscriber biết đã unsubscribe

		log.Println("Unsubscribe")

		// Tìm và xóa channel khỏi mapChannel[topic]
		if chans, ok := ps.mapChannel[topic]; ok {
			for i := range chans {
				// So sánh pointer để tìm đúng channel
				if chans[i] == c {
					// Xóa element tại index i bằng slice trick
					// [a, b, c, d] → xóa b → [a] + [c, d] = [a, c, d]
					chans = append(chans[:i], chans[i+1:]...)

					// CRITICAL SECTION: Update mapChannel
					ps.locker.Lock()
					ps.mapChannel[topic] = chans
					ps.locker.Unlock()

					// Close channel → subscriber sẽ nhận zero value và exit
					close(c)
					break
				}
			}
		}
	}
}

// run khởi động goroutine phân phối messages từ queue đến subscribers
//
// Đây là CORE của PubSub system - goroutine trung tâm phân phối messages
//
// Flow:
// 1. Đọc message từ messageQueue (BLOCK nếu queue trống)
// 2. Log message
// 3. Lock mapChannel (read lock) để tìm subscribers
// 4. Tìm tất cả subscribers của topic này
// 5. Gửi message đến TẤT CẢ subscribers (mỗi subscriber 1 goroutine riêng)
// 6. Unlock và quay lại bước 1
//
// ⚠️ QUAN TRỌNG:
// - Goroutine này chạy MÃI MÃI (infinite loop)
// - BLOCK tại dòng đọc messageQueue khi không có message
// - Mỗi subscriber nhận message trong goroutine RIÊNG (không block nhau)
// - Nếu 1 subscriber chậm/block → không ảnh hưởng subscribers khác
func (ps *localPubSub) run() error {
	// Khởi động goroutine phân phối message
	go func() {
		defer middleware.RecoverGoroutine() // Recovery để tránh crash

		// Vòng lặp vô tận để xử lý messages
		for {
			// BƯỚC 1: Đọc message từ queue (BLOCKING)
			// ⏸️ Goroutine ngủ ở đây nếu queue trống
			mess := <-ps.messageQueue
			log.Println("Message dequeue", mess.String())

			// BƯỚC 2: Tìm subscribers của topic này
			// Dùng RLock (read lock) vì chỉ đọc, không modify
			// → Nhiều goroutines có thể đọc cùng lúc
			ps.locker.RLock()

			// Lấy danh sách channels của subscribers cho topic này
			if subs, ok := ps.mapChannel[mess.Channel()]; ok {
				// BƯỚC 3: Gửi message đến TẤT CẢ subscribers
				// Mỗi subscriber chạy trong goroutine RIÊNG
				for i := range subs {
					// Goroutine riêng cho mỗi subscriber
					// → Nếu 1 subscriber chậm/block, không ảnh hưởng các subscribers khác
					go func(c chan *Message) {
						defer middleware.RecoverGoroutine()

						// Gửi message vào channel của subscriber
						// ⚠️ CÓ THỂ BLOCK nếu subscriber không đọc
						c <- mess
					}(subs[i])
				}
			}
			// Không có subscribers cho topic này → message bị drop

			ps.locker.RUnlock()
		}
	}()

	return nil
}
