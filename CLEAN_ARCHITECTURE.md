# Clean Architecture trong Go - Todo List Project

## 📚 Mục lục

- [Tổng quan](#tổng-quan)
- [Kiến trúc các tầng](#kiến-trúc-các-tầng)
- [Chi tiết từng tầng](#chi-tiết-từng-tầng)
- [Dependency Injection](#dependency-injection)
- [Ví dụ thực tế](#ví-dụ-thực-tế)
- [Lợi ích](#lợi-ích)
- [Quy tắc vàng](#quy-tắc-vàng)

---

## 🎯 Tổng quan

Clean Architecture là một kiến trúc phần mềm giúp tách biệt các tầng của ứng dụng, làm cho code dễ maintain, test và mở rộng.

**Nguyên tắc cốt lõi**: Dependencies luôn đi từ **NGOÀI → TRONG**

```
Client → Handler → Service → UseCase → Storage → Database
```

---

## 🏗️ Kiến trúc các tầng

```
┌────────────────────────────────────────────────────┐
│                   CLIENT (Browser/Mobile)          │
└────────────────────────┬───────────────────────────┘
                         │ HTTP Request
                         ▼
┌────────────────────────────────────────────────────┐
│  LAYER 1: HANDLER (Delivery/Transport Layer)      │
│  📍 module/item/handler/gin_item/item_handler.go   │
│                                                    │
│  • Nhận HTTP request từ client                     │
│  • Parse request body, query params, URL params    │
│  • Validate input format                           │
│  • Gọi Service/UseCase                             │
│  • Trả về HTTP response                            │
└────────────────────────┬───────────────────────────┘
                         │
                         ▼
┌────────────────────────────────────────────────────┐
│  LAYER 2: SERVICE (Optional Orchestration Layer)  │
│  📍 module/item/handler/service.go                 │
│                                                    │
│  • Điều phối giữa Handler và UseCase               │
│  • Tổng hợp nhiều UseCase                          │
│  • Chuyển đổi data giữa các layer                  │
└────────────────────────┬───────────────────────────┘
                         │
                         ▼
┌────────────────────────────────────────────────────┐
│  LAYER 3: USE CASE (Business Logic Layer)         │
│  📍 module/item/use_case/item_use_case.go          │
│                                                    │
│  • Chứa TOÀN BỘ business logic                     │
│  • Quy tắc nghiệp vụ (Business Rules)              │
│  • Validation nghiệp vụ                            │
│  • Điều phối workflow                              │
└────────────────────────┬───────────────────────────┘
                         │
                         ▼
┌────────────────────────────────────────────────────┐
│  LAYER 4: STORAGE/REPOSITORY (Data Access Layer)  │
│  📍 module/item/storage/itemStorage.go             │
│                                                    │
│  • Thực hiện truy vấn database (CRUD)              │
│  • Xử lý SQL queries                               │
│  • Map data từ DB sang struct                      │
└────────────────────────┬───────────────────────────┘
                         │
                         ▼
┌────────────────────────────────────────────────────┐
│              DATABASE (MySQL/PostgreSQL)           │
└────────────────────────────────────────────────────┘
```

---

## 📝 Chi tiết từng tầng

### 1️⃣ Handler Layer (Delivery/Transport)

**File**: `module/item/handler/gin_item/item_handler.go`

**Trách nhiệm**:

- ✅ Nhận và parse HTTP request (JSON, query params, URL params)
- ✅ Validate format của input (không phải business validation)
- ✅ Gọi Service/UseCase
- ✅ Chuyển đổi kết quả thành HTTP response
- ✅ Xử lý HTTP status codes

**KHÔNG làm**:

- ❌ Business logic
- ❌ Database access
- ❌ SQL queries

**Ví dụ**:

```go
func (h *GinItemHandler) DeleteItem(ctx *gin.Context) {
    // Parse URL parameter
    id, err := strconv.Atoi(ctx.Param("id"))
    if err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
        return
    }

    // Gọi service
    if err := h.service.DeleteItem(ctx.Request.Context(), id); err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    ctx.JSON(http.StatusOK, common.SimpleSuccessResponse(true))
}
```

---

### 2️⃣ Service Layer (Optional)

**File**: `module/item/handler/service.go`

**Trách nhiệm**:

- ✅ Điều phối giữa Handler và UseCase
- ✅ Tổng hợp nhiều UseCase khi cần
- ✅ Chuyển đổi data giữa các layer

**KHÔNG làm**:

- ❌ HTTP logic
- ❌ Database access
- ❌ Business logic (để cho UseCase)

**Ví dụ đơn giản** (forward xuống UseCase):

```go
func (s *ItemService) DeleteItem(ctx context.Context, id int) error {
    return s.useCase.DeleteItem(ctx, id)
}
```

**Ví dụ phức tạp** (tổng hợp nhiều UseCase):

```go
func (s *ItemService) DeleteItemWithNotification(ctx context.Context, id int) error {
    // Xóa item
    if err := s.itemUseCase.DeleteItem(ctx, id); err != nil {
        return err
    }

    // Gửi notification
    if err := s.notificationUseCase.SendDeletedNotification(ctx, id); err != nil {
        log.Println("Failed to send notification:", err)
    }

    return nil
}
```

---

### 3️⃣ UseCase Layer (Business Logic) ⭐ QUAN TRỌNG NHẤT

**File**: `module/item/use_case/item_use_case.go`

**Trách nhiệm**:

- ✅ Chứa TOÀN BỘ business logic của ứng dụng
- ✅ Quy tắc nghiệp vụ (Business Rules)
- ✅ Validation nghiệp vụ
- ✅ Điều phối workflow
- ✅ Quyết định logic xử lý

**KHÔNG làm**:

- ❌ HTTP logic (status codes, request/response)
- ❌ SQL queries trực tiếp
- ❌ Framework-specific code

**Ví dụ**:

```go
func (useCase *itemUseCase) DeleteItem(ctx context.Context, id int) error {
    // Business Logic 1: Lấy item để kiểm tra
    itemData, err := useCase.store.GetItem(ctx, map[string]interface{}{"id": id})
    if err != nil {
        return err
    }

    // Business Logic 2: Kiểm tra đã xóa chưa
    if itemData.Status == "Deleted" {
        return model.ErrItemIsDeleted  // Không cho xóa 2 lần
    }

    // Business Logic 3: Thực hiện soft delete
    if err := useCase.store.DeleteItem(ctx, map[string]interface{}{"id": id}); err != nil {
        return err
    }

    return nil
}
```

**Business Rules trong code**:

- ✅ Item đã xóa không được xóa lại
- ✅ Title không được rỗng khi tạo item
- ✅ Sử dụng soft delete (không xóa vật lý)

---

### 4️⃣ Storage Layer (Data Access)

**File**: `module/item/storage/itemStorage.go`

**Trách nhiệm**:

- ✅ Thực hiện CRUD operations với database
- ✅ Viết và thực thi SQL queries
- ✅ Map data giữa DB và Go structs
- ✅ Caching (nếu có)

**KHÔNG làm**:

- ❌ Business logic
- ❌ Validation nghiệp vụ

**Ví dụ**:

```go
func (s *sqlStore) DeleteItem(ctx context.Context, cond map[string]interface{}) error {
    const Deleted = "Deleted"

    // Chỉ thực hiện SQL UPDATE
    if err := s.db.Table(model.TodoItem{}.TableName()).
        Where(cond).
        Updates(map[string]interface{}{
            "status": Deleted,
        }).Error; err != nil {
        return err
    }

    return nil
}
```

---

## 💉 Dependency Injection

**1. Handler nhận dependencies qua constructor**:

```go
type GinItemHandler struct {
    service *handler.ItemService  // Đã được tạo sẵn từ ngoài
}

func NewGinItemHandler(service *handler.ItemService) *GinItemHandler {
    return &GinItemHandler{
        service: service,  // Nhận vào, KHÔNG tự tạo
    }
}

func (h *GinItemHandler) GetItems(ctx *gin.Context) {
    // Chỉ dùng service đã có sẵn
    itemData, err := h.service.GetItems(...)
}
```

**2. Tạo dependencies ở main.go (ngoài cùng)**:

```go
func main() {
    db := gorm.Open(...)

    // Tạo TẤT CẢ dependencies 1 LẦN
    itemStorage := storage.NewSqlStore(db)              // Layer 1
    itemUseCase := usecase.NewItemUseCase(itemStorage)  // Layer 2
    itemService := handler.NewItemService(itemUseCase)  // Layer 3
    itemHandler := ginitem.NewGinItemHandler(itemService) // Layer 4

    // Dùng handler đã được "inject đầy đủ"
    router.GET("/items", itemHandler.GetItems)
}
```

**Lợi ích**:

- ✅ Tạo dependencies 1 lần, dùng lại nhiều lần
- ✅ Dễ test với mock objects
- ✅ Loose coupling
- ✅ Dễ thay đổi implementation

---

## 🔄 Ví dụ thực tế: Request "DELETE Item ID=5"

### Flow xử lý request:

```
1. CLIENT
   ↓ DELETE /v1/items/5

2. HANDLER (item_handler.go)
   • Parse id từ URL: "5" → 5
   • Validate id là số hợp lệ
   • Gọi service.DeleteItem(5)
   ↓

3. SERVICE (service.go)
   • Forward xuống useCase.DeleteItem(5)
   ↓

4. USE CASE (item_use_case.go)
   • Lấy item từ storage để kiểm tra
   • Kiểm tra item.Status != "Deleted" (Business Rule)
   • Gọi storage.DeleteItem(id: 5)
   ↓

5. STORAGE (itemStorage.go)
   • Thực thi SQL: UPDATE todo_items SET status='Deleted' WHERE id=5
   ↓

6. DATABASE
   • Cập nhật record
```

### Code chi tiết:

**Handler**:

```go
func (h *GinItemHandler) DeleteItem(ctx *gin.Context) {
    id, err := strconv.Atoi(ctx.Param("id"))  // Parse "5" → 5
    if err != nil {
        ctx.JSON(400, gin.H{"error": "Invalid ID"})
        return
    }

    if err := h.service.DeleteItem(ctx.Request.Context(), id); err != nil {
        ctx.JSON(400, gin.H{"error": err.Error()})
        return
    }

    ctx.JSON(200, gin.H{"data": true})
}
```

**Service**:

```go
func (s *ItemService) DeleteItem(ctx context.Context, id int) error {
    return s.useCase.DeleteItem(ctx, id)
}
```

**UseCase**:

```go
func (uc *itemUseCase) DeleteItem(ctx context.Context, id int) error {
    // Business Logic: Kiểm tra item có tồn tại
    item, err := uc.store.GetItem(ctx, map[string]interface{}{"id": id})
    if err != nil {
        return err
    }

    // Business Logic: Không cho xóa item đã xóa
    if item.Status == "Deleted" {
        return model.ErrItemIsDeleted
    }

    // Gọi storage để xóa
    return uc.store.DeleteItem(ctx, map[string]interface{}{"id": id})
}
```

**Storage**:

```go
func (s *sqlStore) DeleteItem(ctx context.Context, cond map[string]interface{}) error {
    return s.db.Table("todo_items").
        Where(cond).
        Updates(map[string]interface{}{"status": "Deleted"}).
        Error
}
```

---

## ✨ Lợi ích của Clean Architecture

### 1. Dễ thay đổi từng phần

```go
// Thay Gin → Echo framework
// ✅ Chỉ viết lại Handler layer
// ✅ UseCase, Storage giữ nguyên

// Thay MySQL → PostgreSQL
// ✅ Chỉ viết lại Storage layer
// ✅ UseCase, Handler giữ nguyên
```

### 2. Dễ test

```go
// Test UseCase KHÔNG cần database thật
func TestDeleteItem(t *testing.T) {
    // Mock Storage
    mockStore := &MockStorage{
        GetItemFunc: func(...) (*model.TodoItem, error) {
            return &model.TodoItem{Status: "Active"}, nil
        },
        DeleteItemFunc: func(...) error {
            return nil
        },
    }

    // Test UseCase với mock
    uc := NewItemUseCase(mockStore)
    err := uc.DeleteItem(ctx, 5)

    assert.NoError(t, err)
}
```

### 3. Tái sử dụng code

```go
// Cùng 1 UseCase, dùng cho nhiều delivery mechanisms
itemUseCase := usecase.NewItemUseCase(storage)

// REST API
restHandler := ginitem.NewGinItemHandler(itemUseCase)

// gRPC API
grpcHandler := grpcitem.NewGrpcItemHandler(itemUseCase)

// CLI Command
cliHandler := cliitem.NewCliItemHandler(itemUseCase)
```

### 4. Dễ maintain và scale

- Code được tổ chức rõ ràng theo từng tầng
- Mỗi tầng có trách nhiệm riêng biệt
- Dễ tìm và sửa bugs
- Dễ thêm features mới

---

## 🔑 Quy tắc vàng

### 1. Phân chia trách nhiệm rõ ràng

| Layer       | Chỉ biết                 | KHÔNG biết           |
| ----------- | ------------------------ | -------------------- |
| **Handler** | HTTP, JSON, Query params | Business logic, SQL  |
| **Service** | Điều phối UseCase        | HTTP, SQL            |
| **UseCase** | Business logic           | HTTP, SQL, Framework |
| **Storage** | SQL, Database            | Business logic       |

### 2. Dependencies luôn đi từ NGOÀI → TRONG

```
Handler → Service → UseCase → Storage → Database
```

**KHÔNG BAO GIỜ ngược lại**: UseCase không được import Handler

### 3. UseCase là trung tâm

- **UseCase chứa business logic**
- Handler và Storage là "chi tiết" có thể thay đổi
- UseCase KHÔNG phụ thuộc vào chi tiết

### 4. Sử dụng Interfaces

```go
// UseCase phụ thuộc vào interface, KHÔNG phụ thuộc vào concrete implementation
type ItemStorage interface {
    GetItem(...) (*model.TodoItem, error)
    DeleteItem(...) error
}

type itemUseCase struct {
    store ItemStorage  // Interface, không phải *sqlStore
}
```

### 5. Dependency Injection

- Tạo dependencies ở ngoài cùng (main.go)
- Inject vào thông qua constructor
- KHÔNG tạo dependencies bên trong

---

## 📁 Cấu trúc thư mục

```
module/item/
├── handler/
│   ├── service.go              # Service layer
│   └── gin_item/
│       └── item_handler.go     # Gin HTTP handlers
├── use_case/
│   └── item_use_case.go        # Business logic
├── storage/
│   ├── itemStorage.go          # Interface implementation
│   └── sql.go                  # SQL store
└── model/
    ├── item.go                 # Domain models
    └── filter.go               # Filter models
```

---

## 🚀 Getting Started

### Setup dependencies trong main.go

```go
func main() {
    // 1. Setup database
    db, err := gorm.Open(mysql.Open(os.Getenv("DB_CONN")), &gorm.Config{})
    if err != nil {
        log.Fatal(err)
    }

    // 2. Dependency Injection (từ trong ra ngoài)
    itemStorage := storage.NewSqlStore(db)
    itemUseCase := usecase.NewItemUseCase(itemStorage)
    itemService := handler.NewItemService(itemUseCase)
    itemHandler := ginitem.NewGinItemHandler(itemService)

    // 3. Setup routes
    router := gin.Default()
    v1 := router.Group("/v1")
    {
        items := v1.Group("/items")
        {
            items.GET("", itemHandler.GetItems)
            items.GET("/:id", itemHandler.GetItem)
            items.POST("", itemHandler.CreateItem)
            items.PATCH("/:id", itemHandler.UpdateItem)
            items.DELETE("/:id", itemHandler.DeleteItem)
        }
    }

    // 4. Start server
    router.Run(":8000")
}
```

---

## 📚 Tài liệu tham khảo

- [The Clean Architecture - Uncle Bob](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [Go Clean Architecture Example](https://github.com/bxcodec/go-clean-arch)
- [Dependency Injection in Go](https://github.com/google/wire)

---

## ❓ FAQ

**Q: Service layer có bắt buộc không?**
A: Không. Nếu ứng dụng đơn giản, Handler có thể gọi trực tiếp UseCase. Service layer hữu ích khi cần tổng hợp nhiều UseCase.

**Q: Khi nào nên tạo UseCase mới?**
A: Mỗi business operation nên có 1 method trong UseCase. Nếu quá nhiều operations, tách thành nhiều UseCase files.

**Q: Interface nên đặt ở đâu?**
A: Interface nên đặt ở layer sử dụng nó. Ví dụ: `ItemStorage` interface đặt trong `use_case` vì UseCase sử dụng nó.

**Q: Có nên dùng global variables không?**
A: Không. Luôn dùng Dependency Injection thay vì global variables.

---

**Created**: January 2026
**Author**: Clean Architecture Pattern Implementation
**Project**: Social Todo List API
