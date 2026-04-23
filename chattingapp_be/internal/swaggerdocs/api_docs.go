package swaggerdocs

import "chattingapp_be/internal/dto"

// ============================================================
// SWAGGER GENERAL INFO
// ============================================================




// ============================================================
// DOC RESPONSE MODELS
// Không dùng APIResponse trực tiếp vì Data là interface{},
// swagger sẽ khó hiện đúng schema.
// ============================================================

type ErrorResponse struct {
	Success bool   `json:"success" example:"false"`
	Message string `json:"message" example:"request thất bại"`
	Error   string `json:"error,omitempty" example:"mô tả lỗi"`
}

type SuccessOnlyResponse struct {
	Success bool   `json:"success" example:"true"`
	Message string `json:"message" example:"thành công"`
}

type PingResponse struct {
	Message string `json:"message" example:"pong"`
}

type AuthResponseEnvelope struct {
	Success bool             `json:"success" example:"true"`
	Message string           `json:"message" example:"đăng nhập thành công"`
	Data    dto.AuthResponse `json:"data"`
}

type RefreshRequestDoc struct {
	RefreshToken string `json:"refresh_token" example:"your_refresh_token_here"`
}

type LogoutRequestDoc struct {
	RefreshToken string `json:"refresh_token" example:"your_refresh_token_here"`
}

type UserProfileResponseEnvelope struct {
	Success bool             `json:"success" example:"true"`
	Message string           `json:"message" example:"lấy thông tin cá nhân thành công"`
	Data    dto.UserResponse `json:"data"`
}

type FriendListResponseEnvelope struct {
	Success bool                         `json:"success" example:"true"`
	Message string                       `json:"message" example:"lấy danh sách bạn bè thành công"`
	Data    []dto.FriendListItemResponse `json:"data"`
}

type ConversationDetailResponseEnvelope struct {
	Success bool                           `json:"success" example:"true"`
	Message string                         `json:"message" example:"lấy chi tiết cuộc trò chuyện thành công"`
	Data    dto.ConversationDetailResponse `json:"data"`
}

type ConversationListData struct {
	Items []dto.ConversationListItemResponse `json:"items"`
}

type ConversationListResponseEnvelope struct {
	Success bool                 `json:"success" example:"true"`
	Message string               `json:"message" example:"lấy danh sách cuộc trò chuyện thành công"`
	Data    ConversationListData `json:"data"`
}

type MessageResponseEnvelope struct {
	Success bool                `json:"success" example:"true"`
	Message string              `json:"message" example:"gửi tin nhắn thành công"`
	Data    dto.MessageResponse `json:"data"`
}

type MessageListResponseEnvelope struct {
	Success bool                  `json:"success" example:"true"`
	Message string                `json:"message" example:"lấy danh sách tin nhắn thành công"`
	Data    []dto.MessageResponse `json:"data"`
}
type UserDeviceListResponseEnvelope struct {
	Success bool                     `json:"success" example:"true"`
	Message string                   `json:"message,omitempty" example:"Lấy danh sách thiết bị thành công"`
	Data    []dto.UserDeviceResponse `json:"data"`
}
// ============================================================
// DOC ENDPOINTS
// Chỉ là function rỗng để gắn annotation.
// Không liên quan handler thật.
// ============================================================

// PingDoc godoc
// @Summary Ping server
// @Description Kiểm tra server còn hoạt động hay không
// @Tags System
// @Produce json
// @Success 200 {object} PingResponse
// @Router /ping [get]
func PingDoc() {}

// RegisterDoc godoc
// @Summary Đăng ký tài khoản
// @Description Tạo tài khoản mới
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dto.RegisterRequest true "Thông tin đăng ký"
// @Success 201 {object} AuthResponseEnvelope
// @Failure 400 {object} ErrorResponse
// @Router /auth/register [post]
func RegisterDoc() {}

// LoginDoc godoc
// @Summary Đăng nhập
// @Description Đăng nhập bằng identifier và password
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dto.LoginRequest true "Thông tin đăng nhập"
// @Success 200 {object} AuthResponseEnvelope
// @Failure 400 {object} ErrorResponse
// @Router /auth/login [post]
func LoginDoc() {}

// RefreshDoc godoc
// @Summary Refresh token
// @Description Dùng refresh token để lấy access token mới và refresh token mới
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body swaggerdocs.RefreshRequestDoc true "Refresh token"
// @Success 200 {object} AuthResponseEnvelope
// @Failure 400 {object} ErrorResponse
// @Router /auth/refresh [post]
func RefreshDoc() {}

// LogoutDoc godoc
// @Summary Đăng xuất
// @Description Thu hồi refresh token hiện tại
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body swaggerdocs.LogoutRequestDoc true "Refresh token cần thu hồi"
// @Success 200 {object} SuccessOnlyResponse
// @Failure 400 {object} ErrorResponse
// @Router /auth/logout [post]
func LogoutDoc() {}

// GetMyProfileDoc godoc
// @Summary Lấy thông tin tài khoản hiện tại
// @Description Lấy profile của user đang đăng nhập
// @Tags Auth
// @Security BearerAuth
// @Produce json
// @Success 200 {object} UserProfileResponseEnvelope
// @Failure 401 {object} ErrorResponse
// @Failure 400 {object} ErrorResponse
// @Router /auth/me [get]
func GetMyProfileDoc() {}

// SendFriendRequestDoc godoc
// @Summary Gửi lời mời kết bạn
// @Description Gửi lời mời kết bạn tới user khác
// @Tags Friends
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body dto.SendFriendRequestRequest true "Thông tin lời mời kết bạn"
// @Success 201 {object} SuccessOnlyResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /friends/requests [post]
func SendFriendRequestDoc() {}

// RespondFriendRequestDoc godoc
// @Summary Phản hồi lời mời kết bạn
// @Description accepted | rejected | cancelled
// @Tags Friends
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Friend Request ID"
// @Param request body dto.RespondFriendRequestRequest true "Hành động phản hồi"
// @Success 200 {object} SuccessOnlyResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /friends/requests/{id} [patch]
func RespondFriendRequestDoc() {}

// ListFriendsDoc godoc
// @Summary Lấy danh sách bạn bè
// @Description Trả về danh sách bạn bè của user hiện tại
// @Tags Friends
// @Security BearerAuth
// @Produce json
// @Success 200 {object} FriendListResponseEnvelope
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /friends [get]
func ListFriendsDoc() {}

// CreateDirectConversationDoc godoc
// @Summary Tạo direct conversation
// @Description Tạo cuộc trò chuyện 1-1 với user khác
// @Tags Conversations
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body dto.CreateDirectConversationRequest true "Thông tin tạo direct conversation"
// @Success 201 {object} ConversationDetailResponseEnvelope
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /conversations/direct [post]
func CreateDirectConversationDoc() {}

// ListMyConversationsDoc godoc
// @Summary Lấy danh sách cuộc trò chuyện
// @Description Có hỗ trợ phân trang với page và limit
// @Tags Conversations
// @Security BearerAuth
// @Produce json
// @Param page query int false "Trang hiện tại" default(1)
// @Param limit query int false "Số phần tử mỗi trang" default(20)
// @Success 200 {object} ConversationListResponseEnvelope
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /conversations [get]
func ListMyConversationsDoc() {}

// GetConversationDetailDoc godoc
// @Summary Lấy chi tiết cuộc trò chuyện
// @Description Lấy chi tiết conversation theo ID
// @Tags Conversations
// @Security BearerAuth
// @Produce json
// @Param id path int true "Conversation ID"
// @Success 200 {object} ConversationDetailResponseEnvelope
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /conversations/{id} [get]
func GetConversationDetailDoc() {}

// MarkConversationReadDoc godoc
// @Summary Đánh dấu đã đọc conversation
// @Description Cập nhật last_read_message_id cho thành viên hiện tại
// @Tags Conversations
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Conversation ID"
// @Param request body dto.MarkConversationReadRequest true "Dữ liệu đánh dấu đã đọc"
// @Success 200 {object} SuccessOnlyResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /conversations/{id}/read [patch]
func MarkConversationReadDoc() {}

// UpdateMyNicknameDoc godoc
// @Summary Cập nhật nickname trong conversation
// @Description Cập nhật nickname của user hiện tại trong conversation
// @Tags Conversations
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Conversation ID"
// @Param request body dto.UpdateConversationNicknameRequest true "Nickname mới"
// @Success 200 {object} SuccessOnlyResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /conversations/{id}/nickname [patch]
func UpdateMyNicknameDoc() {}

// MuteConversationDoc godoc
// @Summary Tắt thông báo conversation
// @Description Mute conversation đến thời điểm mute_until hoặc vô thời hạn tùy logic service
// @Tags Conversations
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Conversation ID"
// @Param request body dto.MuteConversationRequest true "Thông tin mute"
// @Success 200 {object} SuccessOnlyResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /conversations/{id}/mute [patch]
func MuteConversationDoc() {}

// SendMessageDoc godoc
// @Summary Gửi tin nhắn
// @Description Gửi một tin nhắn mới vào conversation
// @Tags Messages
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body dto.SendMessageRequest true "Thông tin tin nhắn"
// @Success 201 {object} MessageResponseEnvelope
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /messages [post]
func SendMessageDoc() {}

// EditMessageDoc godoc
// @Summary Sửa tin nhắn
// @Description Chỉ sửa nội dung message theo ID
// @Tags Messages
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Message ID"
// @Param request body dto.EditMessageRequest true "Nội dung sửa"
// @Success 200 {object} MessageResponseEnvelope
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /messages/{id} [patch]
func EditMessageDoc() {}

// DeleteMessageDoc godoc
// @Summary Xóa tin nhắn
// @Description Xóa mềm hoặc theo logic service
// @Tags Messages
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Message ID"
// @Param request body dto.DeleteMessageRequest false "Tùy chọn xóa"
// @Success 200 {object} SuccessOnlyResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /messages/{id} [delete]
func DeleteMessageDoc() {}

// ListMessagesDoc godoc
// @Summary Lấy danh sách tin nhắn
// @Description Lấy danh sách message theo conversation_id, có hỗ trợ page và limit
// @Tags Messages
// @Security BearerAuth
// @Produce json
// @Param conversation_id query int true "Conversation ID"
// @Param page query int false "Trang hiện tại" default(1)
// @Param limit query int false "Số phần tử mỗi trang" default(20)
// @Success 200 {object} MessageListResponseEnvelope
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /messages [get]
func ListMessagesDoc() {}

// RegisterDeviceDoc godoc
// @Summary Đăng ký hoặc cập nhật device
// @Description Đăng ký device cho user hiện tại
// @Tags Devices
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body dto.RegisterDeviceRequest true "Thông tin device"
// @Success 200 {object} SuccessOnlyResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /devices/register [post]
func RegisterDeviceDoc() {}
// UpdateMyProfileDoc godoc
// @Summary Cập nhật profile của user hiện tại
// @Description Cập nhật full_name, avatar_url, bio
// @Tags Auth
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body dto.UpdateProfileRequest true "Thông tin profile mới"
// @Success 200 {object} UserProfileResponseEnvelope
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /auth/profile [patch]
func UpdateMyProfileDoc() {}

// ChangePasswordDoc godoc
// @Summary Đổi mật khẩu
// @Description Yêu cầu mật khẩu cũ và nhập lại mật khẩu mới 2 lần
// @Tags Auth
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body dto.ChangePasswordRequest true "Thông tin đổi mật khẩu"
// @Success 200 {object} SuccessOnlyResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /auth/change-password [patch]
func ChangePasswordDoc() {}
// ForgotPasswordDoc godoc
// @Summary Quên mật khẩu
// @Description Gửi email chứa hướng dẫn đặt lại mật khẩu
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dto.ForgotPasswordRequest true "Email cần khôi phục"
// @Success 200 {object} SuccessOnlyResponse
// @Failure 400 {object} ErrorResponse
// @Router /auth/forgot-password [post]
func ForgotPasswordDoc() {}

// ResetPasswordDoc godoc
// @Summary Đặt lại mật khẩu
// @Description Dùng token reset để đặt lại mật khẩu mới
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dto.ResetPasswordRequest true "Token và mật khẩu mới"
// @Success 200 {object} SuccessOnlyResponse
// @Failure 400 {object} ErrorResponse
// @Router /auth/reset-password [post]
func ResetPasswordDoc() {}
// ListMyDevicesDoc godoc
// @Summary Danh sách thiết bị của tôi
// @Description Lấy danh sách thiết bị đã đăng ký của user hiện tại
// @Tags Devices
// @Security BearerAuth
// @Accept json
// @Produce json
// @Success 200 {object} UserDeviceListResponseEnvelope
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /devices [get]
func ListMyDevicesDoc() {}
// DeleteDeviceDoc godoc
// @Summary Xóa thiết bị
// @Description Disable thiết bị và revoke toàn bộ refresh token của thiết bị đó
// @Tags Devices
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param uuid path string true "Device UUID"
// @Success 200 {object} SuccessOnlyResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /devices/{uuid} [delete]
func DeleteDeviceDoc() {}
// LogoutDeviceDoc godoc
// @Summary Logout theo thiết bị
// @Description Revoke toàn bộ refresh token của thiết bị theo UUID
// @Tags Devices
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param uuid path string true "Device UUID"
// @Success 200 {object} SuccessOnlyResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /devices/{uuid}/logout [post]
func LogoutDeviceDoc() {}
// DisableDeviceDoc godoc
// @Summary Disable thiết bị
// @Description Vô hiệu hóa thiết bị và revoke toàn bộ refresh token của thiết bị đó
// @Tags Devices
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param uuid path string true "Device UUID"
// @Success 200 {object} SuccessOnlyResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /devices/{uuid}/disable [patch]
func DisableDeviceDoc() {}

// TrustDeviceDoc godoc
// @Summary Trust thiết bị
// @Description Đánh dấu thiết bị là đáng tin cậy
// @Tags Devices
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param uuid path string true "Device UUID"
// @Success 200 {object} SuccessOnlyResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /devices/{uuid}/trust [patch]
func TrustDeviceDoc() {}

// UntrustDeviceDoc godoc
// @Summary Untrust thiết bị
// @Description Bỏ đánh dấu thiết bị là đáng tin cậy
// @Tags Devices
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param uuid path string true "Device UUID"
// @Success 200 {object} SuccessOnlyResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /devices/{uuid}/untrust [patch]
func UntrustDeviceDoc() {}