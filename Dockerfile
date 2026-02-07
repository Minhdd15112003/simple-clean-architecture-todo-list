
# Dockerfile đơn giản

# sử dụng hình ảnh cơ sở Alpine Linux nhẹ
# FROM alpine

# # Thiết lập thư mục làm việc trong container
# WORKDIR /app/

# #sao chép nội dung từ thư mục app trên máy chủ vào thư mục /app/ trong container
# ADD ./app /app/

# # Chạy ứng dụng khi container khởi động
# ENTRYPOINT ["./app"]


#Dockerfile-multi-stage
# state: builder
# tip: có thể  build trước một cache và sử dụng lại cache đó trong các lần build tiếp theo để tiết kiệm thời gian build
FROM golang:1.25.5 AS builder

RUN mkdir /app
ADD . /app/
WORKDIR /app
RUN CGO_ENABLED=0 GOOS=linux go build -o  demoApp .

# state: final
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/demoApp .
ENTRYPOINT ["./demoApp"]

