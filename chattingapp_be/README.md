# ChattingApp Backend

Backend API cho hệ thống ChattingApp, được xây dựng bằng Go. Backend phụ trách xác thực người dùng, quản lý bạn bè, cuộc trò chuyện, tin nhắn, thiết bị đăng nhập, refresh token, WebSocket realtime và hỗ trợ lớp E2EE cho tin nhắn mã hóa đầu cuối.

---

## 1. Tổng quan

Backend hiện hỗ trợ các nhóm chức năng chính:

- Xác thực người dùng bằng JWT.
- Đăng ký, đăng nhập, refresh token, logout.
- Cập nhật profile, đổi mật khẩu, quên mật khẩu và đặt lại mật khẩu.
- Quản lý lời mời kết bạn và quan hệ bạn bè.
- Tạo và quản lý cuộc trò chuyện trực tiếp.
- Gửi, sửa, xóa, thu hồi, tìm kiếm và forward tin nhắn.
- Reaction tin nhắn.
- Quản lý thiết bị đăng nhập.
- WebSocket realtime cho message event và typing event.
- Quản lý E2EE key.
- Gửi encrypted message và lưu ciphertext theo từng thiết bị nhận.

---

## 2. Công nghệ sử dụng

### Bắt buộc đúng phiên bản

- Go `1.26.1`

Phiên bản Go được khai báo trong `go.mod`, nên khi setup project nên dùng đúng hoặc tương thích với phiên bản này.

### Công nghệ chính

- Go
- `net/http`
- `http.ServeMux`
- PostgreSQL hoặc hệ quản trị cơ sở dữ liệu tương thích PostgreSQL
- JWT authentication
- Gorilla WebSocket
- Swagger / swaggo
- `github.com/lib/pq` để kết nối PostgreSQL
- `.env` để cấu hình môi trường

> Lưu ý: Backend hiện tại đăng ký route bằng `net/http` và `http.ServeMux`. Không dùng Gin làm router chính, dù trong `go.mod` có dependency liên quan đến Gin.

---

## 3. Cấu trúc thư mục

```text
chattingapp_be/
├─ cmd/
│  └─ api/
│     ├─ main.go
│     └─ swagger.go
│
├─ data/
│  └─ schema.sql
│
├─ docs/
│  ├─ docs.go
│  ├─ swagger.json
│  └─ swagger.yaml
│
├─ internal/
│  ├─ config/
│  │  └─ config.go
│  │
│  ├─ database/
│  │  └─ postgres.go
│  │
│  ├─ dto/
│  │  ├─ auth_dto.go
│  │  ├─ common_dto.go
│  │  ├─ conversation_dto.go
│  │  ├─ conversation_member_dto.go
│  │  ├─ conversation_query_dto.go
│  │  ├─ e2ee_key_dto.go
│  │  ├─ encrypted_message_dto.go
│  │  ├─ friend_request_dto.go
│  │  ├─ friendship_dto.go
│  │  ├─ message_attachment_create_dto.go
│  │  ├─ message_attachment_dto.go
│  │  ├─ message_dto.go
│  │  ├─ user_device_dto.go
│  │  └─ user_dto.go
│  │
│  ├─ handler/
│  │  ├─ auth_handler.go
│  │  ├─ conversation_handler.go
│  │  ├─ e2ee_key_handler.go
│  │  ├─ friend_handler.go
│  │  ├─ helper.go
│  │  ├─ message_handler.go
│  │  └─ user_device_handler.go
│  │
│  ├─ middleware/
│  │  └─ auth_middleware.go
│  │
│  ├─ models/
│  │  ├─ conversation.go
│  │  ├─ conversation_member.go
│  │  ├─ device_identity_key.go
│  │  ├─ device_one_time_prekey.go
│  │  ├─ device_signed_prekey.go
│  │  ├─ friend_request.go
│  │  ├─ friend_request_member.go
│  │  ├─ friendship.go
│  │  ├─ friendship_member.go
│  │  ├─ message.go
│  │  ├─ message_attachment.go
│  │  ├─ message_ciphertext.go
│  │  ├─ message_reaction.go
│  │  ├─ password_reset_token.go
│  │  ├─ user.go
│  │  ├─ user_device.go
│  │  └─ user_refresh_token.go
│  │
│  ├─ realtime/
│  │  └─ hub.go
│  │
│  ├─ repository/
│  │  ├─ conversation_member_repository.go
│  │  ├─ conversation_repository.go
│  │  ├─ device_identity_key_repository.go
│  │  ├─ device_one_time_prekey_repository.go
│  │  ├─ device_signed_prekey_repository.go
│  │  ├─ friend_request_repository.go
│  │  ├─ friendship_repository.go
│  │  ├─ message_attachment_repository.go
│  │  ├─ message_ciphertext_repository.go
│  │  ├─ message_repository.go
│  │  ├─ password_reset_token_repository.go
│  │  ├─ user_device_repository.go
│  │  ├─ user_refresh_token_repository.go
│  │  └─ user_repository.go
│  │
│  ├─ routes/
│  │  └─ routes.go
│  │
│  ├─ service/
│  │  ├─ auth_service.go
│  │  ├─ conversation_service.go
│  │  ├─ e2ee_key_service.go
│  │  ├─ friend_service.go
│  │  ├─ message_service.go
│  │  └─ user_device_service.go
│  │
│  ├─ swaggerdocs/
│  │  └─ api_docs.go
│  │
│  └─ util/
│     └─ jwt.go
│
├─ .env.example
├─ go.mod
└─ go.sum
```

---

## 4. Cấu hình môi trường

Tạo file `.env` trong thư mục `chattingapp_be`.

Ví dụ:

```env
APP_ENV=development
APP_PORT=8080

DB_HOST=localhost
DB_PORT=5432
DB_NAME=chatting_app
DB_USER=postgres
DB_PASSWORD=your_password
DB_SSLMODE=disable

JWT_SECRET=change_me_to_a_long_random_secret_key
JWT_EXPIRES_HOURS=72
REFRESH_EXPIRES_HOURS=168
```

Nên tạo thêm file mẫu:

```text
chattingapp_be/.env.example
```

Nội dung đề xuất cho `.env.example`:

```env
APP_ENV=development
APP_PORT=8080

DB_HOST=localhost
DB_PORT=5432
DB_NAME=chatting_app
DB_USER=postgres
DB_PASSWORD=your_password
DB_SSLMODE=disable

JWT_SECRET=change_me_to_a_long_random_secret_key
JWT_EXPIRES_HOURS=72
REFRESH_EXPIRES_HOURS=168
```

Không nên commit file `.env` thật lên Git vì file này có thể chứa mật khẩu database và JWT secret.

Trong `.gitignore`, nên có:

```gitignore
.env
.env.*
!*.example
```

---

## 5. Cài đặt database

Backend sử dụng PostgreSQL hoặc hệ quản trị cơ sở dữ liệu tương thích PostgreSQL.

### Bước 1: Tạo database

```sql
CREATE DATABASE chatting_app;
```

### Bước 2: Chạy schema

File schema nằm tại:

```text
chattingapp_be/data/schema.sql
```

Nếu đang đứng trong thư mục `chattingapp_be`, có thể chạy bằng `psql`:

```bash
psql -U postgres -d chatting_app -f data/schema.sql
```

Hoặc mở file `data/schema.sql` bằng DBeaver / pgAdmin và chạy trực tiếp trên database `chatting_app`.

---

## 6. Các nhóm bảng chính trong database

Schema hiện có các nhóm bảng chính:

### Người dùng và xác thực

- `users`
- `user_refresh_tokens`
- `password_reset_tokens`

### Bạn bè

- `friend_requests`
- `friend_request_members`
- `friendships`
- `friendship_members`
- `user_blocks`

### Thiết bị

- `user_devices`

### Cuộc trò chuyện

- `conversations`
- `conversation_members`

### Tin nhắn

- `messages`
- `message_attachments`
- `message_reactions`

### E2EE

- `device_identity_keys`
- `device_signed_prekeys`
- `device_one_time_prekeys`
- `message_ciphertexts`

---

## 7. Cài đặt dependency

Đứng trong thư mục backend:

```bash
cd chattingapp_be
```

Tải và đồng bộ dependency:

```bash
go mod tidy
```

Kiểm tra project có build được không:

```bash
go test ./...
```

---

## 8. Chạy backend

Đứng trong thư mục `chattingapp_be`:

```bash
go run ./cmd/api
```

Nếu chạy thành công, terminal sẽ hiển thị tương tự:

```text
server đang chạy tại http://localhost:8080
swagger đang chạy tại http://localhost:8080/swagger/index.html
```

Đường dẫn server mặc định:

```text
http://localhost:8080
```

---

## 9. Swagger API Docs

Backend có Swagger UI để xem và test API.

Sau khi chạy server, mở:

```text
http://localhost:8080/swagger/index.html
```

### Rebuild Swagger

Khi thay đổi annotation Swagger, chạy lại:

```bash
rmdir /s /q docs
swag init -g cmd/api/main.go -d .
```

Sau đó chạy lại server:

```bash
go run ./cmd/api
```

### Authorization trong Swagger

Các API protected yêu cầu Bearer token.

Sau khi login, copy `access_token`, bấm nút Authorize trong Swagger và nhập theo dạng:

```text
Bearer <access_token>
```

---

## 10. Các nhóm API chính

### System

| Method | Endpoint | Mô tả |
|---|---|---|
| GET | `/ping` | Kiểm tra server còn hoạt động |

---

### Auth

| Method | Endpoint | Mô tả |
|---|---|---|
| POST | `/auth/register` | Đăng ký tài khoản |
| POST | `/auth/login` | Đăng nhập |
| POST | `/auth/refresh` | Refresh token |
| POST | `/auth/logout` | Logout |
| POST | `/auth/forgot-password` | Gửi yêu cầu quên mật khẩu |
| POST | `/auth/reset-password` | Đặt lại mật khẩu |
| GET | `/auth/me` | Lấy thông tin user hiện tại |
| PATCH | `/auth/profile` | Cập nhật profile |
| PATCH | `/auth/change-password` | Đổi mật khẩu |

---

### Friends

| Method | Endpoint | Mô tả |
|---|---|---|
| POST | `/friends/requests` | Gửi lời mời kết bạn |
| PATCH | `/friends/requests/{id}` | Phản hồi lời mời kết bạn |
| GET | `/friends` | Lấy danh sách bạn bè |
| GET | `/friends/requests/incoming` | Lấy danh sách lời mời đến |
| GET | `/friends/requests/outgoing` | Lấy danh sách lời mời đi |
| DELETE | `/friends/requests/{id}` | Hủy lời mời kết bạn |
| DELETE | `/friends/{user_id}` | Hủy kết bạn |
| GET | `/friends/search` | Tìm kiếm user |
| GET | `/friends/{user_id}/mutual` | Lấy danh sách bạn chung |
| POST | `/friends/{user_id}/block` | Block user |

---

### Conversations

| Method | Endpoint | Mô tả |
|---|---|---|
| POST | `/conversations/direct` | Tạo cuộc trò chuyện trực tiếp |
| GET | `/conversations` | Lấy danh sách cuộc trò chuyện |
| GET | `/conversations/{id}` | Lấy chi tiết cuộc trò chuyện |
| PATCH | `/conversations/{id}/read` | Đánh dấu đã đọc |
| POST | `/conversations/{id}/typing` | Gửi sự kiện typing |
| PATCH | `/conversations/{id}/nickname` | Cập nhật nickname trong conversation |
| PATCH | `/conversations/{id}/mute` | Tắt thông báo conversation |
| PATCH | `/conversations/{id}/pin` | Ghim conversation |
| PATCH | `/conversations/{id}/archive` | Lưu trữ conversation |
| GET | `/conversations/unread-count` | Lấy số conversation chưa đọc |

---

### Messages

| Method | Endpoint | Mô tả |
|---|---|---|
| POST | `/messages` | Gửi tin nhắn thường |
| GET | `/messages` | Lấy danh sách tin nhắn |
| GET | `/messages/search` | Tìm kiếm tin nhắn |
| GET | `/messages/cursor` | Lấy tin nhắn theo cursor |
| PATCH | `/messages/{id}` | Sửa tin nhắn |
| DELETE | `/messages/{id}` | Xóa tin nhắn |
| PATCH | `/messages/{id}/recall` | Thu hồi tin nhắn |
| POST | `/messages/{id}/forward` | Forward tin nhắn |
| POST | `/messages/{id}/reactions` | Thêm hoặc cập nhật reaction |
| GET | `/messages/{id}/reactions` | Lấy danh sách reaction của tin nhắn |
| DELETE | `/messages/{id}/reactions` | Xóa reaction của user hiện tại |

---

### Encrypted Messages

| Method | Endpoint | Mô tả |
|---|---|---|
| POST | `/messages/encrypted` | Gửi encrypted message |
| GET | `/messages/ciphertexts` | Lấy ciphertext chưa delivered cho device hiện tại |
| PATCH | `/messages/ciphertexts/{id}/delivered` | Đánh dấu ciphertext đã delivered |

Encrypted message không lưu plaintext trong `messages.content`. Backend chỉ lưu metadata message và ciphertext theo từng thiết bị nhận trong bảng `message_ciphertexts`.

---

### Devices

| Method | Endpoint | Mô tả |
|---|---|---|
| POST | `/devices/register` | Đăng ký hoặc cập nhật thiết bị |
| GET | `/devices` | Lấy danh sách thiết bị của user hiện tại |
| DELETE | `/devices/{uuid}` | Xóa / disable thiết bị |
| POST | `/devices/{uuid}/logout` | Logout theo thiết bị |
| PATCH | `/devices/{uuid}/disable` | Disable thiết bị |
| PATCH | `/devices/{uuid}/trust` | Đánh dấu thiết bị tin cậy |
| PATCH | `/devices/{uuid}/untrust` | Bỏ đánh dấu thiết bị tin cậy |

---

### E2EE Key Management

| Method | Endpoint | Mô tả |
|---|---|---|
| POST | `/e2ee/keys/identity` | Upload identity key |
| POST | `/e2ee/keys/signed-prekey` | Upload signed prekey |
| POST | `/e2ee/keys/one-time-prekeys` | Upload one-time prekeys |
| GET | `/e2ee/users/{user_id}/key-bundle` | Lấy key bundle của user nhận |

---

### Realtime WebSocket

| Method | Endpoint | Mô tả |
|---|---|---|
| GET | `/ws?user_id={id}` | Kết nối WebSocket realtime |

Ví dụ:

```text
ws://localhost:8080/ws?user_id=1
```

WebSocket dùng để nhận realtime event như message mới và typing event.

Với encrypted message, WebSocket chỉ gửi metadata, không gửi plaintext hoặc ciphertext. Client cần gọi API `/messages/ciphertexts?device_uuid=...` để lấy ciphertext tương ứng với device hiện tại.

---

## 11. Luồng xác thực cơ bản

### Bước 1: Đăng ký

Gọi:

```text
POST /auth/register
```

### Bước 2: Đăng nhập

Gọi:

```text
POST /auth/login
```

Sau khi đăng nhập thành công, backend trả về `access_token` và `refresh_token`.

### Bước 3: Gọi API protected

Thêm header:

```text
Authorization: Bearer <access_token>
```

### Bước 4: Refresh token

Khi access token hết hạn, gọi:

```text
POST /auth/refresh
```

---

## 12. Luồng device cơ bản

Sau khi login, client nên đăng ký device:

```text
POST /devices/register
```

Device được dùng cho:

- Quản lý phiên đăng nhập theo thiết bị.
- Refresh token theo thiết bị.
- E2EE key theo từng thiết bị.
- Gửi encrypted message theo từng target device.

---

## 13. Tổng quan E2EE

E2EE trong backend được thiết kế theo nguyên tắc: backend không biết plaintext.

Backend chỉ hỗ trợ:

- Lưu public key của từng device.
- Lưu signed prekey.
- Lưu one-time prekey.
- Trả key bundle cho sender.
- Lưu ciphertext theo từng target device.
- Đánh dấu ciphertext đã delivered.

Backend không thực hiện mã hóa hoặc giải mã nội dung tin nhắn. Việc mã hóa và giải mã thuộc phía client.

### Luồng gửi encrypted message

1. User B upload identity key, signed prekey và one-time prekeys cho device của mình.
2. User A gọi API lấy key bundle của user B.
3. User A dùng key bundle để mã hóa nội dung tin nhắn ở phía client.
4. User A gửi ciphertext lên backend qua `/messages/encrypted`.
5. Backend tạo message với `message_type = encrypted`.
6. Backend để `messages.content = NULL`.
7. Backend lưu ciphertext vào `message_ciphertexts`.
8. User B gọi `/messages/ciphertexts?device_uuid=...` để lấy ciphertext của device hiện tại.
9. Sau khi nhận, client có thể gọi `/messages/ciphertexts/{id}/delivered`.

---

## 14. Một số lệnh thường dùng

### Chạy backend

```bash
go run ./cmd/api
```

### Test toàn bộ project

```bash
go test ./...
```

### Format code

```bash
gofmt -w .
```

Hoặc format một file cụ thể:

```bash
gofmt -w internal/service/message_service.go
```

### Rebuild Swagger

```bash
rmdir /s /q docs
swag init -g cmd/api/main.go -d .
```

### Kiểm tra Git trước khi commit

```bash
git status
```

---

## 15. Ghi chú phát triển

- Backend hiện sử dụng `net/http` và `http.ServeMux`.
- Các API protected yêu cầu Bearer token.
- File schema chính nằm ở `data/schema.sql`.
- Swagger docs được generate vào folder `docs/`.
- Không commit file `.env` thật.
- Nên dùng `.env.example` để mô tả các biến môi trường cần thiết.
- E2EE hiện ở mức backend hỗ trợ lưu key/ciphertext và bảo đảm không lưu plaintext.
- Việc mã hóa và giải mã thật sự nằm ở phía client.