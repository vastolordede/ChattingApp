# ChattingApp Frontend

Frontend của hệ thống ChattingApp, được xây dựng bằng Flutter. Phần frontend hiện đang ở giai đoạn khởi tạo project, chuẩn bị cấu trúc thư mục và kiểm tra kết nối cơ bản giữa Flutter app và Backend API thông qua endpoint `/ping`.

---

## 1. Tổng quan

Frontend hiện đã chuẩn bị:

- Flutter project cơ bản.
- Cấu trúc thư mục theo hướng tách `core`, `data`, `features`, `shared`.
- Thư mục assets cho icons và images.
- Màn hình test kết nối Frontend ↔ Backend.
- Gọi thử API backend qua endpoint `/ping`.
- Chuẩn bị các file/folder cho auth, chat, routing, theme, services và repositories.

Phần cấu hình chi tiết Android Studio, emulator, route, service, repository và UI chính sẽ được team frontend tiếp tục bổ sung sau.

---

## 2. Công nghệ sử dụng

### Bắt buộc đúng phiên bản

Frontend hiện đang được kiểm tra với:

- Flutter `3.41.5`
- Dart `3.11.3`

Thông tin từ môi trường hiện tại:

```text
Flutter 3.41.5
Dart 3.11.3
DevTools 2.54.2
```

Nên dùng đúng hoặc tương thích với các phiên bản này để tránh lỗi dependency hoặc lỗi build.

### Công nghệ chính

- Flutter
- Dart
- Material Design 3
- HTTP package để gọi REST API
- Android Studio để quản lý Android SDK, emulator và thiết bị Android
- VS Code để viết code Flutter

### Package chính trong `pubspec.yaml`

```yaml
dependencies:
  flutter:
    sdk: flutter
  cupertino_icons: ^1.0.8
  http: ^1.2.1

dev_dependencies:
  flutter_test:
    sdk: flutter
  flutter_lints: ^6.0.0
```

---

## 3. Cài đặt môi trường

### 3.1. Cài Flutter SDK

Cài Flutter SDK bản stable.

Sau khi cài, kiểm tra bằng lệnh:

```bash
flutter --version
```

Môi trường hiện tại đang dùng:

```text
Flutter 3.41.5
Dart 3.11.3
```

### 3.2. Cài Android Studio

Cài Android Studio bản stable mới.

Android Studio được dùng để:

- Cài Android SDK.
- Cài Android SDK Platform.
- Cài Android SDK Command-line Tools.
- Tạo và quản lý Android Emulator.
- Kiểm tra Android toolchain cho Flutter.

Trong Android Studio, nên cài thêm:

- Flutter plugin
- Dart plugin

### 3.3. Cài Android SDK

Trong Android Studio, mở:

```text
Settings > Languages & Frameworks > Android SDK
```

Hoặc:

```text
More Actions > SDK Manager
```

Nên kiểm tra các thành phần sau:

- Android SDK Platform
- Android SDK Build-Tools
- Android SDK Platform-Tools
- Android SDK Command-line Tools
- Android Emulator

Môi trường hiện tại đang nhận:

```text
Android SDK version 36.1.0
```

Không bắt buộc team phải dùng đúng duy nhất bản SDK này, nhưng nên dùng Android SDK tương thích với Flutter hiện tại và kiểm tra lại bằng `flutter doctor`.

### 3.4. Kiểm tra môi trường

Chạy:

```bash
flutter doctor
```

Môi trường hiện tại đang đạt trạng thái:

```text
[✓] Flutter
[✓] Windows Version
[✓] Android toolchain
[✓] Chrome
[✓] Visual Studio
[✓] Connected device
[✓] Network resources

No issues found!
```

Nếu `flutter doctor` báo thiếu Android toolchain hoặc Android Studio, cần mở Android Studio để cài SDK/Toolchain còn thiếu.

---

## 4. Cấu trúc thư mục

```text
chattingapp_fe/
├─ android/
├─ assets/
│  ├─ icons/
│  └─ images/
├─ build/
├─ ios/
├─ lib/
│  ├─ core/
│  │  ├─ constants/
│  │  │  └─ api_constants.dart
│  │  ├─ routes/
│  │  │  └─ app_routes.dart
│  │  ├─ theme/
│  │  └─ utils/
│  ├─ data/
│  │  ├─ models/
│  │  ├─ repositories/
│  │  └─ services/
│  │     ├─ auth_api_service.dart
│  │     ├─ chat_api_service.dart
│  │     └─ websocket_service.dart
│  ├─ features/
│  │  ├─ auth/
│  │  ├─ chat/
│  │  └─ splash/
│  ├─ shared/
│  │  └─ widgets/
│  │     └─ primary_button.dart
│  └─ main.dart
├─ linux/
├─ macos/
├─ test/
├─ web/
├─ windows/
├─ pubspec.yaml
├─ pubspec.lock
└─ README.md
```

---

## 5. Ý nghĩa các thư mục chính

### `lib/core`

Chứa các phần dùng chung toàn app.

Dự kiến gồm:

- `constants/`: cấu hình hằng số như API base URL.
- `routes/`: cấu hình route màn hình.
- `theme/`: cấu hình màu sắc, typography, theme.
- `utils/`: hàm tiện ích, validator, formatter.

### `lib/data`

Chứa phần xử lý dữ liệu.

Dự kiến gồm:

- `models/`: model dữ liệu như user, conversation, message.
- `services/`: gọi API backend và WebSocket.
- `repositories/`: lớp trung gian giữa controller/UI và service.

### `lib/features`

Chứa các tính năng chính của app.

Dự kiến gồm:

- `auth/`: login, register, auth controller.
- `chat/`: danh sách conversation, chat screen, message input.
- `splash/`: màn hình splash hoặc kiểm tra trạng thái đăng nhập ban đầu.

### `lib/shared`

Chứa widget dùng chung.

Ví dụ:

- Button dùng chung.
- Input dùng chung.
- Loading widget.
- Error widget.

### `assets`

Chứa tài nguyên tĩnh.

Hiện đã chuẩn bị:

```text
assets/
├─ icons/
└─ images/
```

Nếu sau này sử dụng assets trong Flutter, cần khai báo trong `pubspec.yaml`.

Ví dụ:

```yaml
flutter:
  uses-material-design: true

  assets:
    - assets/images/
    - assets/icons/
```

---

## 6. Cài dependency

Đứng trong thư mục frontend:

```bash
cd chattingapp_fe
```

Tải dependency:

```bash
flutter pub get
```

Kiểm tra project:

```bash
flutter analyze
```

Chạy test nếu có:

```bash
flutter test
```

---

## 7. Chạy frontend

### 7.1. Chạy trên Chrome

```bash
flutter run -d chrome
```

### 7.2. Chạy trên Windows desktop

Nếu máy đã hỗ trợ Windows desktop:

```bash
flutter run -d windows
```

### 7.3. Chạy trên Android Emulator

Mở Android Studio, tạo hoặc mở emulator.

Sau đó kiểm tra thiết bị:

```bash
flutter devices
```

Chạy app:

```bash
flutter run
```

Hoặc chọn device cụ thể:

```bash
flutter run -d <device_id>
```

---

## 8. Kết nối Frontend với Backend

Hiện tại file `lib/main.dart` đang dùng màn hình test kết nối Backend.

API đang được gọi:

```text
GET /ping
```

Khi chạy bằng Android Emulator, frontend gọi backend local thông qua địa chỉ:

```text
http://10.0.2.2:8080/ping
```

Lý do:

- `localhost` trong Android Emulator là máy ảo Android.
- `10.0.2.2` là địa chỉ đặc biệt để emulator truy cập `localhost` của máy thật.
- Backend cần chạy ở máy thật tại port `8080`.

Backend cần chạy trước:

```bash
cd chattingapp_be
go run ./cmd/api
```

Sau đó frontend gọi:

```text
http://10.0.2.2:8080/ping
```

Nếu backend trả về thành công, màn hình sẽ hiển thị message từ API.

---

## 9. Ghi chú API Base URL

Hiện tại API URL đang được test trực tiếp trong `main.dart`:

```dart
Uri.parse('http://10.0.2.2:8080/ping')
```

Trong quá trình phát triển tiếp theo, nên chuyển cấu hình API base URL vào file:

```text
lib/core/constants/api_constants.dart
```

Ví dụ đề xuất:

```dart
class ApiConstants {
  static const String emulatorBaseUrl = 'http://10.0.2.2:8080';
  static const String webBaseUrl = 'http://localhost:8080';
}
```

Sau này team frontend có thể chỉnh lại theo môi trường:

- Android Emulator: `http://10.0.2.2:8080`
- Web Chrome: `http://localhost:8080`
- Physical device: dùng IP LAN của máy chạy backend, ví dụ `http://192.168.1.x:8080`

---

## 10. Android Studio

Android Studio hiện dùng để hỗ trợ Flutter chạy trên Android.

Cần cài:

- Android Studio stable
- Flutter plugin
- Dart plugin
- Android SDK
- Android SDK Platform-Tools
- Android SDK Command-line Tools
- Android Emulator

Kiểm tra bằng:

```bash
flutter doctor
```

Nếu chưa có emulator, tạo trong Android Studio:

```text
Tools > Device Manager > Create Virtual Device
```

Phần cấu hình chi tiết emulator, Android SDK path và Android Studio sẽ được bổ sung sau bởi team frontend.

---

## 11. File `main.dart` hiện tại

Hiện tại `main.dart` đang là màn hình test kết nối backend.

Chức năng chính:

- Khởi tạo Flutter app.
- Bật chế độ full screen immersive.
- Khóa hướng portrait.
- Gọi API `/ping`.
- Hiển thị kết quả trả về từ backend.
- Có nút gọi lại API.

Mục đích của file hiện tại là kiểm tra Frontend có gọi được Backend hay không.

Sau này khi hoàn thiện app, `main.dart` có thể được chỉnh lại để chạy routing, splash screen, auth flow và chat flow.

---

## 12. Lệnh thường dùng

### Kiểm tra Flutter version

```bash
flutter --version
```

### Kiểm tra môi trường

```bash
flutter doctor
```

### Tải dependency

```bash
flutter pub get
```

### Chạy app

```bash
flutter run
```

### Chạy app trên Chrome

```bash
flutter run -d chrome
```

### Chạy app trên Windows

```bash
flutter run -d windows
```

### Phân tích code

```bash
flutter analyze
```

### Chạy test

```bash
flutter test
```

### Xóa cache build

```bash
flutter clean
```

Sau khi clean, chạy lại:

```bash
flutter pub get
```

---

## 13. Ghi chú phát triển

- Frontend hiện đang ở giai đoạn khởi tạo và kiểm tra kết nối với backend.
- Backend cần chạy trước khi test màn hình gọi API.
- Android Emulator dùng `10.0.2.2` để gọi backend local.
- Nếu chạy Flutter Web, có thể dùng `localhost:8080`.
- Nếu chạy trên điện thoại thật, cần dùng IP LAN của máy chạy backend.
- Cấu hình route, API service, repository, state management và UI chính sẽ được bổ sung tiếp.
- Phần cấu hình chi tiết Android Studio và emulator sẽ được team frontend cập nhật sau.