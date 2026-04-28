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

---

## 16. Quy ước phát triển cho backend

Backend đi theo luồng xử lý chính:

Client
  ↓
Routes
  ↓
Middleware
  ↓
Handler
  ↓
Service
  ↓
Repository
  ↓
Database

### Vai trò từng layer

- `routes/`: đăng ký endpoint và gắn middleware.
- `middleware/`: xác thực JWT, lấy `user_id` và đưa vào request context.
- `handler/`: nhận request, decode JSON/query/path params, gọi service và trả response.
- `service/`: xử lý nghiệp vụ chính, kiểm tra quyền, validate logic, phối hợp nhiều repository.
- `repository/`: thao tác database, không chứa logic nghiệp vụ.
- `dto/`: định nghĩa request/response trả ra API.
- `models/`: ánh xạ dữ liệu gần với database.
- `util/`: hàm tiện ích như JWT.

### Quy tắc khi thêm feature mới

1. Tạo DTO request/response nếu API có input/output mới.
2. Thêm function repository nếu cần truy vấn database mới.
3. Viết logic nghiệp vụ trong service.
4. Handler chỉ gọi service, không viết SQL trong handler.
5. Đăng ký route trong `internal/routes/routes.go`.
6. Thêm Swagger annotation hoặc doc stub.
7. Chạy:
   ```bash
   gofmt -w .
   go test ./...
   swag init -g cmd/api/main.go -d .
---

## 17. CI/CD và Deploy Backend

Backend hiện đã được cấu hình CI/CD cơ bản với:

- **CI:** GitHub Actions
- **CD:** Render
- **Database production:** Neon PostgreSQL
- **Deploy runtime:** Docker

Luồng CI/CD hiện tại:

```text
Push code lên nhánh riêng
  ↓
Tạo Pull Request vào develop
  ↓
GitHub Actions chạy CI
  ↓
CI thành công
  ↓
Review code
  ↓
Merge vào develop
  ↓
Render tự động deploy backend
  ↓
Backend production được cập nhật
```

---

## 18. Branching workflow

Nhánh chính dùng để tích hợp code là:

```text
develop
```

Không push trực tiếp vào `develop`. Mọi thay đổi nên được làm trên nhánh riêng, sau đó tạo Pull Request vào `develop`.

### Quy ước đặt tên nhánh

```text
feature/<ten-chuc-nang>
fix/<ten-loi>
docs/<noi-dung-tai-lieu>
ci/<noi-dung-ci>
deploy/<noi-dung-deploy>
security/<noi-dung-bao-mat>
```

Ví dụ:

```bash
git checkout develop
git pull origin develop
git checkout -b feature/message-search
```

Sau khi code xong:

```bash
git add .
git commit -m "feat: add message search"
git push -u origin feature/message-search
```

Sau đó tạo Pull Request:

```text
feature/message-search → develop
```

---

## 19. Pull Request rule

Nhánh `develop` được cấu hình rule để hạn chế merge trực tiếp.

Quy trình chuẩn khi merge code:

```text
Tạo nhánh mới
  ↓
Commit code
  ↓
Push nhánh lên GitHub
  ↓
Tạo Pull Request vào develop
  ↓
Đợi CI chạy thành công
  ↓
Review code
  ↓
Merge vào develop
```

Hiện tại Pull Request nên đảm bảo:

- Không có file nhạy cảm bị commit.
- Code Go đã được format bằng `gofmt`.
- Test chạy thành công.
- Backend build thành công.
- Người review kiểm tra logic trước khi merge.

---

## 20. GitHub Actions CI

Backend CI được chạy bằng GitHub Actions.

File workflow:

```text
.github/workflows/be-ci.yml
```

CI hiện chạy khi có thay đổi liên quan đến backend hoặc workflow CI.

Các bước kiểm tra chính:

```text
1. Checkout source code
2. Check sensitive files
3. Setup Go theo version trong go.mod
4. Download dependencies
5. Check Go formatting bằng gofmt
6. Run tests bằng go test ./...
7. Build backend bằng go build
```

### Lệnh kiểm tra tương đương ở local

Trước khi push code, có thể tự chạy:

```bash
cd chattingapp_be
gofmt -w .
go test ./...
go build -o server ./cmd/api
```

Nếu các lệnh trên chạy ổn thì khả năng cao CI trên GitHub cũng sẽ pass.

---

## 21. Bảo vệ file nhạy cảm

Dự án có 3 lớp bảo vệ file nhạy cảm:

```text
1. .gitignore
2. GitHub Secret Protection / Push Protection
3. GitHub Actions CI check sensitive files
```

### Các file không được commit

Không commit các file chứa mật khẩu, token, private key, JWT secret hoặc database URL thật.

Các pattern bị chặn gồm:

```text
.env
.env.*
*.pem
*.key
*.p12
*.pfx
*.crt
*.cer
*.secret
*.secrets
secrets/
id_rsa
id_rsa.pub
```

Chỉ được commit file mẫu:

```text
.env.example
```

File `.env.example` chỉ dùng để mô tả các biến môi trường cần thiết, không chứa mật khẩu thật.

### Local env

Môi trường local dùng file:

```text
chattingapp_be/.env
```

File này không được push lên GitHub.

Ví dụ local `.env`:

```env
APP_ENV=development
APP_PORT=8080

DB_HOST=localhost
DB_PORT=5432
DB_NAME=chatting_app
DB_USER=postgres
DB_PASSWORD=your_local_password
DB_SSLMODE=disable

JWT_SECRET=change_me_to_a_long_random_secret_key
JWT_EXPIRES_HOURS=72
REFRESH_EXPIRES_HOURS=168
```

---

## 22. Docker

Backend có Dockerfile tại:

```text
chattingapp_be/Dockerfile
```

Render dùng Dockerfile này để build và chạy backend production.

### Build Docker image ở local

Nếu máy có Docker, có thể test bằng:

```bash
cd chattingapp_be
docker build -t chattingapp-be .
```

Chạy container:

```bash
docker run --env-file .env -p 8080:8080 chattingapp-be
```

Sau đó test:

```text
http://localhost:8080/ping
```

Nếu chưa dùng Docker local thì vẫn có thể deploy bằng Render vì Render sẽ tự build Docker image trên server của Render.

---

## 23. Continuous Deployment bằng Render

Backend được deploy trên Render.

Production URL:

```text
https://chattingapp-wxgj.onrender.com
```

Swagger UI production:

```text
https://chattingapp-wxgj.onrender.com/swagger/index.html
```

Health check:

```text
https://chattingapp-wxgj.onrender.com/ping
```

### Cấu hình Render

Cấu hình service trên Render:

```text
Runtime: Docker
Branch: develop
Root Directory: chattingapp_be
Region: Singapore
Instance Type: Free
```

Khi có code mới được merge vào `develop`, Render sẽ tự động deploy lại backend.

### Lưu ý với Render Free

Render Free có thể sleep sau một thời gian không có request.

Khi server sleep, request đầu tiên có thể chậm hơn bình thường vì Render cần khởi động lại service.

---

## 24. Production environment variables

Production không dùng file `.env` trong GitHub.

Các biến môi trường production được cấu hình trực tiếp trên Render.

Danh sách biến cần có:

```env
APP_ENV=production
APP_PORT=8080

DB_HOST=<neon-host>
DB_PORT=5432
DB_NAME=neondb
DB_USER=<neon-user>
DB_PASSWORD=<neon-password>
DB_SSLMODE=require

JWT_SECRET=<production-secret>
JWT_EXPIRES_HOURS=72
REFRESH_EXPIRES_HOURS=168
```

Trong đó:

- `DB_HOST`, `DB_NAME`, `DB_USER`, `DB_PASSWORD` lấy từ Neon PostgreSQL.
- `DB_SSLMODE=require` là bắt buộc khi kết nối Neon.
- `JWT_SECRET` phải là chuỗi dài, random, không dùng lại secret test.

Có thể tạo JWT secret bằng PowerShell:

```powershell
[guid]::NewGuid().ToString() + [guid]::NewGuid().ToString()
```

---

## 25. Database production trên Neon

Production database đang dùng Neon PostgreSQL.

Local development vẫn dùng PostgreSQL local.

### Phân biệt database local và production

Local:

```text
DB_HOST=localhost
DB_NAME=chatting_app
DB_SSLMODE=disable
```

Production:

```text
DB_HOST=<neon-host>
DB_NAME=neondb
DB_SSLMODE=require
```

Không dùng database local cho production.

Không dùng Neon production cho việc test lung tung nếu không cần thiết.

### Import database lên Neon

Có thể export database local bằng DBeaver hoặc `pg_dump`, sau đó import vào Neon.

Sau khi import xong, cần kiểm tra trên Neon:

- Có đủ bảng.
- Có dữ liệu cần thiết.
- Các foreign key không lỗi.
- Các sequence/id vẫn hoạt động đúng.

Test nhanh bằng cách gọi API production:

```text
POST /auth/register
POST /auth/login
GET /auth/me
```

Nếu user mới được tạo trong Neon thì backend production đã kết nối database đúng.

---

## 26. Kiểm tra production sau deploy

Sau mỗi lần deploy, nên kiểm tra nhanh các endpoint sau:

### Ping server

```text
GET /ping
```

URL:

```text
https://chattingapp-wxgj.onrender.com/ping
```

Kết quả mong muốn:

```json
{
  "message": "pong"
}
```

### Swagger

```text
https://chattingapp-wxgj.onrender.com/swagger/index.html
```

### Auth flow

Test các API:

```text
POST /auth/register
POST /auth/login
GET /auth/me
```

Sau khi login, copy `access_token` và dùng trong Swagger:

```text
Bearer <access_token>
```

---

## 27. Quy trình làm việc sau khi đã có CI/CD

Khi phát triển feature mới:

```bash
git checkout develop
git pull origin develop
git checkout -b feature/<ten-feature>
```

Sau khi sửa code:

```bash
cd chattingapp_be
gofmt -w .
go test ./...
go build -o server ./cmd/api
```

Commit:

```bash
git add .
git commit -m "feat: mo ta ngan gon thay doi"
git push -u origin feature/<ten-feature>
```

Tạo Pull Request vào:

```text
develop
```

Đợi CI pass, review code rồi merge.

Sau khi merge, Render sẽ tự deploy backend.

---

## 28. Khi CI fail thì xử lý thế nào

Nếu CI fail, vào tab **Actions** hoặc tab **Checks** trong Pull Request để xem log.

Một số lỗi thường gặp:

### Lỗi gofmt

Nếu CI báo file chưa format:

```bash
cd chattingapp_be
gofmt -w .
```

Sau đó commit lại:

```bash
git add .
git commit -m "style: format go code"
git push
```

### Lỗi test

Nếu `go test ./...` fail, đọc log để biết package hoặc function nào lỗi.

Chạy lại ở local:

```bash
cd chattingapp_be
go test ./...
```

### Lỗi build

Nếu `go build` fail, thường là do:

- Import sai package.
- Function đổi tên nhưng chỗ khác chưa sửa.
- Thiếu dependency.
- Sai đường dẫn package.
- Code compile không qua.

Chạy lại ở local:

```bash
cd chattingapp_be
go build -o server ./cmd/api
```

### Lỗi sensitive files

Nếu CI báo có file nhạy cảm bị track, cần remove khỏi Git index:

```bash
git rm --cached <ten-file>
```

Ví dụ:

```bash
git rm --cached chattingapp_be/.env
```

Sau đó commit lại:

```bash
git add .
git commit -m "chore: remove sensitive file from git tracking"
git push
```

---

## 29. Khi Render deploy fail thì xử lý thế nào

Nếu Render deploy fail, vào Render Dashboard và mở tab **Logs**.

Các lỗi thường gặp:

### Sai env production

Kiểm tra lại Render Environment Variables:

```text
APP_ENV
APP_PORT
DB_HOST
DB_PORT
DB_NAME
DB_USER
DB_PASSWORD
DB_SSLMODE
JWT_SECRET
JWT_EXPIRES_HOURS
REFRESH_EXPIRES_HOURS
```

### Sai database host

`DB_HOST` chỉ nên là host, không bao gồm:

```text
postgresql://
```

Không để nguyên cả connection string vào `DB_HOST`.

### Sai SSL mode

Với Neon, cần:

```text
DB_SSLMODE=require
```

### Sai port

Render cần app listen theo port được cấu hình.

Hiện backend dùng:

```text
APP_PORT=8080
```

### Docker build fail

Kiểm tra lại:

```text
Root Directory: chattingapp_be
Dockerfile: Dockerfile
```

Dockerfile phải nằm tại:

```text
chattingapp_be/Dockerfile
```

---

## 30. Ghi chú bảo mật

- Không commit `.env` thật.
- Không commit database password.
- Không commit JWT secret.
- Không commit private key.
- Không commit file certificate thật.
- Không ghi password trực tiếp vào README.
- Không gửi production secret qua chat nhóm.
- Nếu lỡ leak secret, cần reset ngay trên Neon hoặc Render.
- Nếu lỡ commit secret lên GitHub, chỉ xóa file ở commit mới là chưa đủ. Cần rotate secret vì secret cũ đã bị lộ.

---

## 31. Trạng thái hiện tại

Backend hiện đã có:

```text
- API backend Go
- PostgreSQL schema
- Swagger UI
- JWT authentication
- Refresh token
- Friend system
- Direct conversation
- Message system
- Reaction
- Forward message
- Recall message
- Device management
- WebSocket realtime
- E2EE key management
- Encrypted message storage
- GitHub Actions CI
- Sensitive file protection
- Dockerfile
- Render deployment
- Neon PostgreSQL production database
```

Backend production hiện chạy tại:

```text
https://chattingapp-wxgj.onrender.com
```

Swagger production:

```text
https://chattingapp-wxgj.onrender.com/swagger/index.html
```

Health check:

```text
https://chattingapp-wxgj.onrender.com/ping
```