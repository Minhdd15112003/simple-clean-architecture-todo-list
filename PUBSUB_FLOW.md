# 📡 PUBSUB SYSTEM - FLOW CHI TIẾT

## 🎯 Tổng Quan

Hệ thống PubSub (Publish-Subscribe) cho phép các thành phần giao tiếp với nhau một cách **bất đồng bộ** (async) thông qua **messages** và **topics**.

### Lợi ích:

- ✅ **Decoupling**: API không phụ thuộc trực tiếp vào business logic
- ✅ **Performance**: API trả về nhanh, xử lý chạy background
- ✅ **Scalability**: Dễ dàng thêm subscribers mới
- ✅ **Reliability**: Có retry mechanism từ asyncjob

---

## 🏗️ Kiến Trúc Tổng Thể

```
┌─────────────┐
│   API       │ User bấm like
│  Handler    │
└──────┬──────┘
       │ 1. Lưu DB
       │ 2. Publish message
       ▼
┌─────────────────────────────────────────────┐
│          LOCAL PUBSUB                       │
│  ┌──────────────────────────────────────┐  │
│  │  messageQueue (buffer 10000)         │  │
│  │  [ msg1 ] [ msg2 ] [ msg3 ] ...      │  │
│  └────────────┬─────────────────────────┘  │
│               │                             │
│               ▼                             │
│  ┌──────────────────────────────────────┐  │
│  │  run() goroutine (phân phối)         │  │
│  └────────────┬─────────────────────────┘  │
│               │                             │
│               ▼                             │
│  ┌──────────────────────────────────────┐  │
│  │  mapChannel                          │  │
│  │  {                                   │  │
│  │    "TopicLike": [ch1, ch2],         │  │
│  │    "TopicDislike": [ch3]            │  │
│  │  }                                   │  │
│  └──────────┬───────────────────────────┘  │
└─────────────┼───────────────────────────────┘
              │
              ├─────────────┬─────────────┐
              ▼             ▼             ▼
       ┌───────────┐ ┌───────────┐ ┌───────────┐
       │Subscriber │ │Subscriber │ │Subscriber │
       │    1      │ │    2      │ │    3      │
       └─────┬─────┘ └─────┬─────┘ └─────┬─────┘
             │             │             │
             ▼             ▼             ▼
       IncreaseLike  SendNotif    UpdateCache
```

---

## 🔄 FLOW HOẠT ĐỘNG CHI TIẾT

### 📤 PUBLISH FLOW (Khi User Like Item)

```
┌─────────────────────────────────────────────────────────────────┐
│ T=0ms: User bấm like                                            │
└─────────────────────────────────────────────────────────────────┘
           │
           ▼
┌─────────────────────────────────────────────────────────────────┐
│ API Handler: POST /items/123/like                              │
│ ├─ Parse request                                               │
│ ├─ Validate user                                               │
│ └─ Call: likeUseCase.LikeItem()                               │
└─────────────────────────────────────────────────────────────────┘
           │
           ▼
┌─────────────────────────────────────────────────────────────────┐
│ Use Case: LikeItem()                                            │
│ ├─ BƯỚC 1: Lưu vào DB user_like_items                          │
│ │   INSERT INTO user_like_items (user_id, item_id)            │
│ │                                                               │
│ ├─ BƯỚC 2: Publish message (KHÔNG BLOCK)                      │
│ │   ps.Publish(ctx, TopicUserLikeItem, message)              │
│ │                                                               │
│ └─ Return nil ✅ (Nhanh - chỉ mất vài ms)                      │
└─────────────────────────────────────────────────────────────────┘
           │
           ▼
┌─────────────────────────────────────────────────────────────────┐
│ T=5ms: API trả về SUCCESS cho user                             │
│ Response: 200 OK                                                │
└─────────────────────────────────────────────────────────────────┘

═══════════════════════════════════════════════════════════════════
        SAU ĐÓ - XỬ LÝ BẤT ĐỒNG BỘ (ASYNC)
═══════════════════════════════════════════════════════════════════

┌─────────────────────────────────────────────────────────────────┐
│ PubSub.Publish() - Chạy trong goroutine riêng                  │
│                                                                 │
│ go func() {                                                     │
│    // Gửi message vào messageQueue                             │
│    ps.messageQueue <- message                                  │
│                                                                 │
│    // Message: {                                               │
│    //   topic: "TopicUserLikeItem",                           │
│    //   data: {UserID: 1, ItemID: 123}                        │
│    // }                                                         │
│ }()                                                             │
└─────────────────────────────────────────────────────────────────┘
           │
           ▼
┌─────────────────────────────────────────────────────────────────┐
│ messageQueue                                                    │
│ [ Message{topic: TopicUserLikeItem, data: {1, 123}} ]         │
│                                                                 │
│ Buffer size: 10000                                             │
│ Current: 1 message                                             │
└─────────────────────────────────────────────────────────────────┘
           │
           ▼
┌─────────────────────────────────────────────────────────────────┐
│ PubSub.run() Goroutine (đang ngủ, chờ message)                │
│                                                                 │
│ for {                                                           │
│    mess := <-ps.messageQueue  ⏸️ BLOCKING - Được đánh thức!   │
│                                                                 │
│    log: "Message dequeue TopicUserLikeItem"                    │
│                                                                 │
│    // Tìm subscribers của topic này                            │
│    subs := mapChannel["TopicUserLikeItem"]                    │
│    // → [channel1, channel2]                                   │
│                                                                 │
│    // Gửi đến TẤT CẢ subscribers                              │
│    for each subscriber_channel {                               │
│       go func(ch) {                                            │
│          ch <- mess  // Gửi message                            │
│       }(subscriber_channel)                                    │
│    }                                                            │
│ }                                                               │
└─────────────────────────────────────────────────────────────────┘
           │
           ├────────────────────┬────────────────────┐
           ▼                    ▼                    ▼
     Subscriber 1          Subscriber 2        Subscriber N
     (Increase Like)       (Send Notif)        (Update Cache)
```

---

### 📥 SUBSCRIBE FLOW (Xử Lý Message)

```
┌─────────────────────────────────────────────────────────────────┐
│ SETUP PHASE (Khi server start)                                 │
└─────────────────────────────────────────────────────────────────┘
           │
           ▼
┌─────────────────────────────────────────────────────────────────┐
│ main.go                                                         │
│ ├─ ps := pubsub.NewPubSub()                                    │
│ ├─ itemStore := storage.NewSqlStore(db)                        │
│ │                                                               │
│ ├─ subEngine := subscriber.NewEngine(ps)                       │
│ └─ subEngine.Start(itemStore)                                  │
│    │                                                            │
│    ├─ startSubTopic(TopicUserLikeItem, ...)                   │
│    └─ startSubTopic(TopicUserDislikeItem, ...)                │
└─────────────────────────────────────────────────────────────────┘
           │
           ▼
┌─────────────────────────────────────────────────────────────────┐
│ startSubTopic(TopicUserLikeItem)                               │
│                                                                 │
│ // BƯỚC 1: Subscribe topic                                     │
│ channel, _ := ps.Subscribe(ctx, "TopicUserLikeItem")          │
│                                                                 │
│ // BƯỚC 2: Start goroutine lắng nghe                          │
│ go func() {                                                     │
│    for {                                                        │
│       msg := <-channel  ⏸️ NGỦ - Đợi message                   │
│                                                                 │
│       // ... xử lý message (xem bên dưới)                      │
│    }                                                            │
│ }()                                                             │
└─────────────────────────────────────────────────────────────────┘

═══════════════════════════════════════════════════════════════════
        KHI CÓ MESSAGE MỚI
═══════════════════════════════════════════════════════════════════

┌─────────────────────────────────────────────────────────────────┐
│ T=10ms: Subscriber Goroutine được đánh thức                    │
│                                                                 │
│ msg := <-channel  ✅ Nhận được message                         │
│ // msg.data = {UserID: 1, ItemID: 123}                        │
│                                                                 │
│ log: "Received message from topic: TopicUserLikeItem"         │
└─────────────────────────────────────────────────────────────────┘
           │
           ▼
┌─────────────────────────────────────────────────────────────────┐
│ BƯỚC 3: Tạo AsyncJobs từ Consumer Jobs                        │
│                                                                 │
│ jobs := []                                                      │
│                                                                 │
│ for each consumerJob {                                         │
│    job := asyncjob.NewJob(func(ctx) error {                   │
│       // Handler: IncreaseLikeCount                            │
│       return consumerJob.Handler(ctx, msg)                     │
│    }, asyncjob.WithName("Increase like count"))              │
│                                                                 │
│    jobs = append(jobs, job)                                    │
│ }                                                               │
└─────────────────────────────────────────────────────────────────┘
           │
           ▼
┌─────────────────────────────────────────────────────────────────┐
│ BƯỚC 4: Chạy Job Group                                         │
│                                                                 │
│ asyncjob.NewGroup(isConcurrent=false, jobs...).Run(ctx)       │
│                                                                 │
│ ├─ Job 1: IncreaseLikeCount(itemID=123)                      │
│ │  ├─ Execute()                                               │
│ │  │  └─ UPDATE items SET liked_count = liked_count + 1      │
│ │  │     WHERE id = 123                                       │
│ │  │                                                           │
│ │  ├─ Nếu LỖI → Retry (1s, 2s, 4s)                          │
│ │  └─ Nếu OK → Done ✅                                        │
│ │                                                             │
│ └─ Job 2: (nếu có job khác)                                   │
└─────────────────────────────────────────────────────────────────┘
           │
           ▼
┌─────────────────────────────────────────────────────────────────┐
│ T=50ms: Tất cả jobs hoàn thành                                │
│ ├─ like_count đã được tăng trong DB                           │
│ └─ Subscriber goroutine quay lại đợi message tiếp theo        │
│    msg := <-channel  ⏸️ NGỦ tiếp                               │
└─────────────────────────────────────────────────────────────────┘
```

---

## 🔧 CÁC THÀNH PHẦN CHI TIẾT

### 1️⃣ **Message** (`pubsub/pubsub.go`)

```go
type Message struct {
    id        string        // ID duy nhất (timestamp nano)
    topic     Topic         // Topic name
    data      interface{}   // Payload (ví dụ: &Like{UserID, ItemID})
    CreatedAt time.Time     // Thời gian tạo
}
```

**Ví dụ:**

```go
msg := &Message{
    id: "1738774800000000000",
    topic: "TopicUserLikeItem",
    data: &model.Like{UserID: 1, ItemID: 123},
    CreatedAt: time.Now()
}
```

---

### 2️⃣ **localPubSub** (`pubsub/local_ps.go`)

```go
type localPubSub struct {
    messageQueue chan *Message              // Queue trung tâm
    mapChannel   map[Topic][]chan *Message  // Topic → Subscribers
    locker       *sync.RWMutex              // Đồng bộ hóa
}
```

#### **Cấu trúc mapChannel:**

```go
mapChannel = {
    "TopicUserLikeItem": [
        channel1,  // Subscriber 1: IncreaseLikeCount
        channel2,  // Subscriber 2: SendNotification
    ],
    "TopicUserDislikeItem": [
        channel3,  // Subscriber 3: DecreaseLikeCount
    ],
}
```

#### **Vai trò của run() goroutine:**

```
┌───────────────────────────────────────────────┐
│  messageQueue (FIFO Queue)                    │
│  [ msg1 ] [ msg2 ] [ msg3 ] [ msg4 ]          │
└──────────────┬────────────────────────────────┘
               │
               ▼
       ┌──────────────┐
       │ run()        │ ← Goroutine chạy mãi mãi
       │ goroutine    │   Đọc từ queue, phân phối
       └──────┬───────┘
              │
              ├─ Đọc msg1 → Topic A
              │  └─ Gửi đến subscribers của A
              │
              ├─ Đọc msg2 → Topic B
              │  └─ Gửi đến subscribers của B
              │
              └─ ...
```

---

### 3️⃣ **consumerEngine** (`subscriber/setup.go`)

```go
type consumerEngine struct {
    subTopic pubsub.PubSup  // Reference đến PubSub
}
```

**Start() method:**

```go
func (ce *consumerEngine) Start(store HasIncreaseLikeCount) {
    // Đăng ký 2 subscribers
    ce.startSubTopic(TopicUserLikeItem, false,
        IncreaseLikeCountAfterUserLikeItem(ctx, store))

    ce.startSubTopic(TopicUserDislikeItem, false,
        DecreaseLikeCountAfterUserUnlikeItem(ctx, store))
}
```

---

### 4️⃣ **AsyncJob Integration**

Mỗi consumer job được wrap trong asyncjob để có:

- ✅ **Retry mechanism**: Tự động retry khi lỗi
- ✅ **Error handling**: Xử lý lỗi tốt hơn
- ✅ **State management**: Theo dõi trạng thái job

```go
job := asyncjob.NewJob(func(ctx) error {
    return store.IncreaseLikeCount(ctx, itemID)
}, asyncjob.WithName("Increase like count"))

// Job sẽ tự động retry theo: 1s → 2s → 4s
// Nếu vẫn lỗi sau 3 lần → return error
```

---

## 🎭 2 CHẾ ĐỘ XỬ LÝ

### **Sequential (isConcurrent = false)** ✅ Đang dùng

```
Message 1 arrives:
  ├─ Job A: IncreaseLikeCount() → xong
  └─ Job B: SendNotification() → xong

Message 2 arrives:
  ├─ Job A: IncreaseLikeCount() → xong
  └─ Job B: SendNotification() → xong
```

**Đảm bảo thứ tự, an toàn**

---

### **Concurrent (isConcurrent = true)**

```
Message 1 arrives:
  ├─ Job A: IncreaseLikeCount() ╮
  └─ Job B: SendNotification()  ├─ Chạy song song
                                 ╯
Message 2 arrives (chưa đợi Msg1 xong):
  ├─ Job A: IncreaseLikeCount() ╮
  └─ Job B: SendNotification()  ├─ Chạy song song
                                 ╯
```

**Nhanh hơn nhưng không đảm bảo thứ tự**

---

## 📊 GOROUTINES BREAKDOWN

### Trong hệ thống hiện tại (2 topics):

```
main() goroutine
│
├─ PubSub.run() goroutine (1)
│  └─ Phân phối messages từ queue
│
├─ PubSub.Publish() goroutines (N)
│  └─ Mỗi lần publish → 1 goroutine mới
│
├─ Subscriber 1: TopicUserLikeItem
│  └─ startSubTopic() goroutine (1)
│     └─ Lắng nghe messages
│
└─ Subscriber 2: TopicUserDislikeItem
   └─ startSubTopic() goroutine (1)
      └─ Lắng nghe messages
```

**Tổng goroutines tối thiểu: 4**

- 1 main
- 1 run()
- 2 subscribers

---

## 🚀 PERFORMANCE & SCALABILITY

### **Bottlenecks:**

1. **messageQueue buffer (10000)**
   - Nếu vượt quá → Publish bị block
   - Giải pháp: Tăng buffer hoặc dùng external queue (Redis, Kafka)

2. **Subscriber channel (unbuffered)**
   - Nếu subscriber chậm → run() goroutine bị block
   - Giải pháp: Buffer channels hoặc drop messages

### **Best Practices:**

✅ Giữ handler nhanh (< 100ms)
✅ Dùng asyncjob cho retry
✅ Monitor queue size
✅ Graceful shutdown để xử lý hết messages

---

## 🐛 ERROR HANDLING

### **Các tầng error handling:**

```
1. AsyncJob Level:
   ├─ Execute() fail → Retry (1s, 2s, 4s)
   └─ Retry hết → Return error

2. Job Group Level:
   ├─ Sequential: Job fail → Stop group
   └─ Concurrent: Job fail → Continue others

3. Subscriber Level:
   └─ Log error, không crash goroutine

4. PubSub Level:
   └─ Recover panic trong goroutines
```

---

## 📝 EXAMPLE: Complete Flow

### Kịch bản: User 1 like item 123

```
T=0ms
├─ User click "Like" button
├─ POST /items/123/like
├─ Handler → Use Case → Store
├─ INSERT INTO user_like_items
└─ ps.Publish(TopicUserLikeItem, {UserID: 1, ItemID: 123})

T=2ms
├─ API returns 200 OK ✅
└─ User sees "Liked!" animation

═══ Background Processing ═══

T=3ms
├─ Message → messageQueue
└─ run() goroutine picks up message

T=4ms
├─ run() → finds subscribers of TopicUserLikeItem
└─ Sends message to subscriber channels

T=5ms
├─ Subscriber goroutine wakes up
├─ Creates AsyncJob: IncreaseLikeCount
└─ Job Group starts

T=6ms
├─ Job executes: UPDATE items SET liked_count = liked_count + 1
└─ Success ✅

T=10ms
└─ liked_count updated in database

User doesn't wait for any of this! 🚀
```

---

## 🎯 KẾT LUẬN

Hệ thống PubSub giúp:

- ✅ API nhanh (< 10ms)
- ✅ Business logic tách biệt
- ✅ Dễ mở rộng (thêm subscribers)
- ✅ Reliable (có retry)
- ✅ Maintainable (code rõ ràng)

**Trade-offs:**

- ❌ Eventual consistency (like count update có độ trễ)
- ❌ Phức tạp hơn (nhiều goroutines)
- ❌ Khó debug (async flow)

**Khi nào dùng:**

- ✅ Actions không cần kết quả ngay (notifications, analytics)
- ✅ Heavy operations (email, image processing)
- ❌ Critical operations cần kết quả ngay (payment, authentication)
