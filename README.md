# ChattingApp

ChattingApp là hệ thống chat realtime gồm Backend API và Frontend Flutter. Dự án hỗ trợ các chức năng chính như xác thực người dùng, quản lý bạn bè, cuộc trò chuyện trực tiếp, tin nhắn, thiết bị đăng nhập, WebSocket realtime và backend support cho E2EE.

---

## 1. Tổng quan dự án

Dự án được chia thành 2 phần chính:

```text
ChattingApp/
├─ chattingapp_be/      # Backend API viết bằng Go
├─ chattingapp_fe/      # Frontend app viết bằng Flutter
├─ .gitignore
└─ README.md
```

### Backend

Backend phụ trách:

- Đăng ký, đăng nhập, refresh token, logout.
- Quản lý profile, đổi mật khẩu, quên mật khẩu.
- Quản lý bạn bè và lời mời kết bạn.
- Quản lý conversation trực tiếp.
- Gửi, sửa, xóa, thu hồi, tìm kiếm và forward tin nhắn.
- Reaction tin nhắn.
- Quản lý thiết bị đăng nhập.
- WebSocket realtime.
- Quản lý E2EE key.
- Lưu encrypted message theo từng thiết bị nhận.

Tài liệu chi tiết backend nằm tại:

```text
chattingapp_be/README.md
```

### Frontend

Frontend phụ trách:

- Giao diện người dùng.
- Khởi tạo app Flutter.
- Chuẩn bị cấu trúc thư mục theo `core`, `data`, `features`, `shared`.
- Test kết nối Frontend ↔ Backend thông qua API `/ping`.
- Chuẩn bị nền tảng cho màn hình auth, chat, routing, theme, service và repository.

Tài liệu chi tiết frontend nằm tại:

```text
chattingapp_fe/README.md
```

---

## 2. Công nghệ sử dụng

### Backend

Các công nghệ chính:

- Go
- `net/http`
- `http.ServeMux`
- PostgreSQL hoặc hệ quản trị cơ sở dữ liệu tương thích PostgreSQL
- JWT authentication
- Gorilla WebSocket
- Swagger / swaggo
- `github.com/lib/pq`
- `.env` để cấu hình môi trường

Phiên bản bắt buộc hoặc nên dùng đúng:

```text
Go 1.26.1
```

Backend hiện tại đăng ký route bằng `net/http` và `http.ServeMux`, không dùng Gin làm router chính.

### Frontend

Các công nghệ chính:

- Flutter
- Dart
- Material Design 3
- HTTP package
- Android Studio
- VS Code

Phiên bản đang được kiểm tra và nên dùng đúng hoặc tương thích:

```text
Flutter 3.41.5
Dart 3.11.3
```

Môi trường frontend hiện đang nhận:

```text
Android SDK version 36.1.0
```

Không bắt buộc mọi máy phải dùng đúng duy nhất Android SDK `36.1.0`, nhưng nên dùng Android SDK tương thích với Flutter hiện tại và kiểm tra bằng:

```bash
flutter doctor
```

---

## 3. Cấu trúc thư mục tổng quan

```text
ChattingApp/
├─ chattingapp_be/
│  ├─ cmd/
│  │  └─ api/
│  ├─ data/
│  │  └─ schema.sql
│  ├─ docs/
│  ├─ internal/
│  │  ├─ config/
│  │  ├─ database/
│  │  ├─ dto/
│  │  ├─ handler/
│  │  ├─ middleware/
│  │  ├─ models/
│  │  ├─ realtime/
│  │  ├─ repository/
│  │  ├─ routes/
│  │  ├─ service/
│  │  ├─ swaggerdocs/
│  │  └─ util/
│  ├─ .env.example
│  ├─ go.mod
│  ├─ go.sum
│  └─ README.md
│
├─ chattingapp_fe/
│  ├─ android/
│  ├─ assets/
│  │  ├─ icons/
│  │  └─ images/
│  ├─ ios/
│  ├─ lib/
│  │  ├─ core/
│  │  │  ├─ constants/
│  │  │  ├─ routes/
│  │  │  ├─ theme/
│  │  │  └─ utils/
│  │  ├─ data/
│  │  │  ├─ models/
│  │  │  ├─ repositories/
│  │  │  └─ services/
│  │  ├─ features/
│  │  │  ├─ auth/
│  │  │  ├─ chat/
│  │  │  └─ splash/
│  │  ├─ shared/
│  │  └─ main.dart
│  ├─ linux/
│  ├─ macos/
│  ├─ test/
│  ├─ web/
│  ├─ windows/
│  ├─ pubspec.yaml
│  ├─ pubspec.lock
│  └─ README.md
│
├─ .gitignore
└─ README.md
```

---

## 4. Cài đặt Backend

Xem hướng dẫn chi tiết tại:

```text
chattingapp_be/README.md
```

### 4.1. Tạo file môi trường

Tạo file:

```text
chattingapp_be/.env
```

Có thể copy từ:

```text
chattingapp_be/.env.example
```

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

Không commit file `.env` thật lên Git.

### 4.2. Tạo database

Backend dùng PostgreSQL hoặc hệ quản trị cơ sở dữ liệu tương thích PostgreSQL.

Tạo database:

```sql
CREATE DATABASE chatting_app;
```

Chạy schema:

```bash
cd chattingapp_be
psql -U postgres -d chatting_app -f data/schema.sql
```

Hoặc mở file sau bằng DBeaver / pgAdmin rồi chạy trực tiếp:

```text
chattingapp_be/data/schema.sql
```

### 4.3. Cài dependency backend

```bash
cd chattingapp_be
go mod tidy
```

### 4.4. Chạy backend

```bash
go run ./cmd/api
```

Nếu chạy thành công:

```text
server đang chạy tại http://localhost:8080
swagger đang chạy tại http://localhost:8080/swagger/index.html
```

Swagger UI:

```text
http://localhost:8080/swagger/index.html
```

---

## 5. Cài đặt Frontend

Xem hướng dẫn chi tiết tại:

```text
chattingapp_fe/README.md
```

### 5.1. Kiểm tra Flutter

```bash
flutter --version
```

Môi trường đang được dùng trong project:

```text
Flutter 3.41.5
Dart 3.11.3
```

### 5.2. Kiểm tra Android Studio và Android toolchain

```bash
flutter doctor
```

Yêu cầu cơ bản:

- Flutter SDK
- Dart SDK
- Android Studio
- Android SDK
- Android SDK Platform-Tools
- Android SDK Command-line Tools
- Android Emulator
- Flutter plugin trong Android Studio
- Dart plugin trong Android Studio

Phần cấu hình chi tiết Android Studio và emulator sẽ được team frontend bổ sung thêm trong `chattingapp_fe/README.md`.

### 5.3. Cài dependency frontend

```bash
cd chattingapp_fe
flutter pub get
```

### 5.4. Chạy frontend

Chạy trên Chrome:

```bash
flutter run -d chrome
```

Chạy trên Windows desktop:

```bash
flutter run -d windows
```

Chạy trên Android Emulator:

```bash
flutter devices
flutter run
```

---

## 6. Kết nối Frontend với Backend

Backend mặc định chạy tại:

```text
http://localhost:8080
```

Khi chạy Flutter bằng Android Emulator, không dùng `localhost` để gọi backend local. Thay vào đó dùng:

```text
http://10.0.2.2:8080
```

Ví dụ endpoint test hiện tại:

```text
http://10.0.2.2:8080/ping
```

Lý do:

- `localhost` trong Android Emulator là máy ảo Android.
- `10.0.2.2` là địa chỉ đặc biệt để emulator truy cập `localhost` của máy thật.

Nếu chạy Flutter Web bằng Chrome, có thể dùng:

```text
http://localhost:8080
```

Nếu chạy trên điện thoại thật, cần dùng IP LAN của máy chạy backend, ví dụ:

```text
http://192.168.1.x:8080
```

---

## 7. API Documentation

Backend có Swagger UI để xem và test API.

Sau khi chạy backend, mở:

```text
http://localhost:8080/swagger/index.html
```

Các API protected cần header:

```text
Authorization: Bearer <access_token>
```

Khi thay đổi Swagger annotation ở backend, rebuild Swagger:

```bash
cd chattingapp_be
rmdir /s /q docs
swag init -g cmd/api/main.go -d .
```

Sau đó chạy lại backend:

```bash
go run ./cmd/api
```

---

## 8. Các nhóm chức năng chính

### Auth

- Register
- Login
- Refresh token
- Logout
- Forgot password
- Reset password
- Get current user
- Update profile
- Change password

### Friends

- Send friend request
- Accept / reject friend request
- List friends
- Incoming friend requests
- Outgoing friend requests
- Cancel friend request
- Unfriend
- Search users
- Mutual friends
- Block user

### Conversations

- Create direct conversation
- List conversations
- Conversation detail
- Mark as read
- Typing event
- Update nickname
- Mute conversation
- Pin conversation
- Archive conversation
- Unread count

### Messages

- Send message
- List messages
- Search messages
- Cursor pagination
- Edit message
- Delete message
- Recall message
- Forward message
- React message
- List reactions
- Delete reaction

### Devices

- Register device
- List my devices
- Delete / disable device
- Logout device
- Trust device
- Untrust device

### Realtime

- WebSocket endpoint:

```text
/ws?user_id={id}
```

Ví dụ:

```text
ws://localhost:8080/ws?user_id=1
```

### E2EE

Backend hỗ trợ E2EE ở mức lưu và phân phối dữ liệu cần thiết:

- Upload identity key
- Upload signed prekey
- Upload one-time prekeys
- Get receiver key bundle
- Send encrypted message
- Store ciphertext per target device
- Get undelivered ciphertexts for current device
- Mark ciphertext delivered

Backend không mã hóa hoặc giải mã plaintext. Việc mã hóa và giải mã thuộc phía client.

---

## 9. E2EE Notes

Encrypted message được thiết kế theo nguyên tắc backend không biết plaintext.

Với encrypted message:

- Client mã hóa nội dung trước khi gửi.
- Backend tạo message với `message_type = encrypted`.
- Backend không lưu plaintext trong `messages.content`.
- Ciphertext được lưu trong bảng `message_ciphertexts`.
- Mỗi target device có một ciphertext riêng.
- WebSocket chỉ gửi metadata, không gửi plaintext hoặc ciphertext.
- Receiver device lấy ciphertext qua API `/messages/ciphertexts`.

Luồng tổng quát:

1. User B upload identity key, signed prekey và one-time prekeys cho device.
2. User A lấy key bundle của User B.
3. User A mã hóa message ở phía client.
4. User A gửi ciphertext lên backend.
5. Backend lưu metadata message và ciphertext theo từng target device.
6. User B gọi API để lấy ciphertext của device hiện tại.
7. User B decrypt ở phía client.

---

## 10. Git Ignore

Dự án có `.gitignore` ở root để bỏ qua các file build, cache, môi trường và file sinh tự động.

Các nhóm chính bị ignore:

### VS Code

```gitignore
.vscode/
```

### Flutter

```gitignore
chattingapp_fe/build/
chattingapp_fe/.dart_tool/
chattingapp_fe/.idea/
chattingapp_fe/.flutter-plugins
chattingapp_fe/.flutter-plugins-dependencies
chattingapp_fe/.pub/
chattingapp_fe/coverage/
```

### Android build files

```gitignore
chattingapp_fe/android/.gradle/
chattingapp_fe/android/local.properties
chattingapp_fe/android/app/debug/
chattingapp_fe/android/app/profile/
chattingapp_fe/android/app/release/
```

### iOS build files

```gitignore
chattingapp_fe/ios/Flutter/.last_build_id
chattingapp_fe/ios/Pods/
chattingapp_fe/ios/.symlinks/
```

### Backend

```gitignore
chattingapp_be/.env
chattingapp_be/.env.*
chattingapp_be/bin/
chattingapp_be/build/
chattingapp_be/tmp/
chattingapp_be/.cache/
```

Lưu ý: nếu muốn commit file mẫu môi trường backend, cần đảm bảo `.env.example` không bị ignore.

Nên dùng rule dạng:

```gitignore
chattingapp_be/.env
chattingapp_be/.env.*
!chattingapp_be/.env.example
```

Không nên dùng:

```gitignore
chattingapp_be/!*.example
```

vì dòng này không đúng ý nghĩa để unignore `.env.example`.

---

## 11. Lệnh thường dùng

### Backend

Chạy backend:

```bash
cd chattingapp_be
go run ./cmd/api
```

Test backend:

```bash
cd chattingapp_be
go test ./...
```

Rebuild Swagger:

```bash
cd chattingapp_be
rmdir /s /q docs
swag init -g cmd/api/main.go -d .
```

### Frontend

Tải dependency:

```bash
cd chattingapp_fe
flutter pub get
```

Chạy frontend:

```bash
cd chattingapp_fe
flutter run
```

Chạy frontend trên Chrome:

```bash
cd chattingapp_fe
flutter run -d chrome
```

Kiểm tra môi trường Flutter:

```bash
flutter doctor
```

Kiểm tra Flutter version:

```bash
flutter --version
```

---

## 12. Tài liệu chi tiết

Backend README:

```text
chattingapp_be/README.md
```

Frontend README:

```text
chattingapp_fe/README.md
```

Database schema:

```text
chattingapp_be/data/schema.sql
```

Swagger UI sau khi chạy backend:

```text
http://localhost:8080/swagger/index.html
```

---

## 13. Ghi chú phát triển

- Backend cần chạy trước nếu muốn frontend test API.
- Frontend Android Emulator gọi backend local bằng `10.0.2.2`.
- Không commit file `.env` thật.
- Nên commit `.env.example`.
- Không commit thư mục build/cache.
- Backend dùng `net/http` và `http.ServeMux`.
- Frontend dùng Flutter.
- E2EE hiện được backend hỗ trợ ở mức lưu key, key bundle và ciphertext; mã hóa/giải mã thật sự nằm phía client.