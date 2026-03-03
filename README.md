# Social Todo List API

Backend RESTful API cho ứng dụng mạng xã hội quản lý công việc, được xây dựng bằng Golang và Gin Framework.

## 📋 Tổng quan

Dự án này là một hệ thống backend API cho phép người dùng:

- Quản lý danh sách công việc (Todo items)
- Tương tác xã hội (Like/Unlike items)
- Xác thực và phân quyền người dùng
- Upload hình ảnh
- Real-time updates sử dụng Pub/Sub pattern

## 🏗️ Kiến trúc

Dự án áp dụng **Clean Architecture** với các tầng rõ ràng:

```
┌─────────────────────────────────────────┐
│         Handler (HTTP Transport)        │  ← Gin Handlers
├─────────────────────────────────────────┤
│      Service (Application Logic)        │  ← Business orchestration
├─────────────────────────────────────────┤
│       UseCase (Business Logic)           │  ← Core business rules
├─────────────────────────────────────────┤
│       Storage (Data Access)              │  ← Database operations
└─────────────────────────────────────────┘
```

## 🚀 Tính năng chính

### 1. Authentication & Authorization

- JWT-based authentication
- Secure user registration and login
- Protected routes với middleware
- Token-based session management

### 2. User Management

- User registration
- User login
- Get user profile
- Secure password handling

### 3. Item Management (Todo)

- Create new items
- Get all items với pagination
- Get item by ID
- Update items
- Delete items
- Filter và search

### 4. Social Features

- Like/Unlike items
- Get users who liked an item
- Real-time like count updates
- Get liked items by user

### 5. File Upload

- Image upload functionality
- Static file serving
- Secure file handling

### 6. Async Processing

- Pub/Sub pattern cho event handling
- Async job processing
- Real-time like count synchronization

## 🛠️ Tech Stack

- **Language:** Go 1.25.5
- **Web Framework:** Gin
- **ORM:** GORM
- **Database:** MySQL
- **Authentication:** JWT (golang-jwt/jwt)
- **Architecture Pattern:** Clean Architecture, Repository Pattern
- **Other Libraries:**
  - go-resty (HTTP client)
  - UUID generation
  - Image processing

## 📁 Cấu trúc dự án

```
socialTodoList/
├── main.go                      # Entry point
├── common/                      # Shared utilities
│   ├── app_err.go              # Error handling
│   ├── app_response.go         # API response format
│   ├── paging.go               # Pagination helper
│   └── asyncjob/               # Async job utilities
├── components/                  # Reusable components
│   └── tokenprovider/          # JWT provider
├── middleware/                  # HTTP middlewares
│   ├── authorize.go            # Auth middleware
│   └── recover.go              # Recovery middleware
├── module/                      # Business modules
│   ├── item/                   # Todo item module
│   │   ├── handler/
│   │   ├── model/
│   │   ├── storage/
│   │   └── use_case/
│   ├── user/                   # User module
│   │   ├── handler/
│   │   ├── model/
│   │   ├── storage/
│   │   └── use_case/
│   └── userlikeitem/           # Like feature module
│       ├── handler/
│       ├── model/
│       ├── storage/
│       └── use_case/
├── pubsub/                      # Pub/Sub implementation
└── subscriber/                  # Event subscribers
```

## 🔧 Cài đặt và chạy

### Prerequisites

- Go 1.25+ installed
- MySQL server running
- Git

### 1. Clone repository

```bash
git clone https://github.com/[YOUR_USERNAME]/socialTodoList.git
cd socialTodoList
```

### 2. Cài đặt dependencies

```bash
go mod download
```

### 3. Setup database

```sql
CREATE DATABASE social_todo_list;
```

### 4. Cấu hình environment variables

```bash
export DB_CONN="user:password@tcp(localhost:3306)/social_todo_list?charset=utf8mb4&parseTime=True&loc=Local"
export SECRET_KEY="your-secret-key-here"
export ITEM_LIKE_SERVICE_URL="http://localhost:8001"  # Optional
```

### 5. Chạy ứng dụng

```bash
go run main.go
```

Server sẽ chạy tại `http://localhost:8080`

## 📚 API Documentation

### Authentication

#### Register

```http
POST /v1/auth/register
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "password123",
  "first_name": "John",
  "last_name": "Doe"
}
```

#### Login

```http
POST /v1/auth/login
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "password123"
}
```

#### Get Profile

```http
GET /v1/auth/profile
Authorization: Bearer <token>
```

### Items (Todo)

#### Get All Items

```http
GET /v1/items?page=1&limit=10
Authorization: Bearer <token>
```

#### Get Item by ID

```http
GET /v1/items/:id
Authorization: Bearer <token>
```

#### Create Item

```http
POST /v1/items
Authorization: Bearer <token>
Content-Type: application/json

{
  "title": "My Todo",
  "description": "Description here"
}
```

#### Update Item

```http
PATCH /v1/items/:id
Authorization: Bearer <token>
Content-Type: application/json

{
  "title": "Updated title",
  "status": "Doing"
}
```

#### Delete Item

```http
DELETE /v1/items/:id
Authorization: Bearer <token>
```

### Likes

#### Like Item

```http
POST /v1/items/:id/like
Authorization: Bearer <token>
```

#### Unlike Item

```http
DELETE /v1/items/:id/unlike
Authorization: Bearer <token>
```

#### Get Users who liked an item

```http
GET /v1/items/:id/liked-users
Authorization: Bearer <token>
```

### Upload

#### Upload Image

```http
PUT /v1/upload
Authorization: Bearer <token>
Content-Type: multipart/form-data

file: <image-file>
```

## 🔐 Security Features

- Password hashing
- JWT token authentication
- Request validation
- SQL injection prevention (GORM)
- CORS handling
- Recovery middleware for panic handling

## 📝 Design Patterns sử dụng

1. **Clean Architecture** - Phân tầng rõ ràng, dễ test và maintain
2. **Repository Pattern** - Abstraction cho data access layer
3. **Dependency Injection** - Loose coupling giữa các components
4. **Middleware Pattern** - Xử lý cross-cutting concerns
5. **Pub/Sub Pattern** - Event-driven architecture cho async processing

## 🧪 Testing

```bash
# Run all tests
go test -v ./...

# Run tests with coverage
go test -v -coverprofile=coverage.out ./...
```

## 🏃‍♂️ Build và Deploy

### Build

```bash
go build -v ./...
```

### Docker (nếu có Dockerfile)

```bash
docker build -t social-todo-list .
docker run -p 8080:8080 social-todo-list
```

## 📈 Roadmap

- [ ] Add unit tests coverage
- [ ] Implement Redis caching
- [ ] Add WebSocket support
- [ ] Implement rate limiting
- [ ] Add API documentation with Swagger
- [ ] Implement soft delete
- [ ] Add full-text search

## 👤 Author

[Your Name]

- GitHub: [@yourusername](https://github.com/yourusername)
- Email: your.email@example.com

## 📄 License

This project is for learning purposes.

## 🙏 Acknowledgments

- Gin Framework
- GORM
- Clean Architecture principles
