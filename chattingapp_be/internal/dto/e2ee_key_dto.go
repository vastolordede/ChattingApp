package dto

type UploadIdentityKeyRequest struct {
	DeviceUUID  string `json:"device_uuid"`
	PublicKey   string `json:"public_key"`
	Algorithm   string `json:"algorithm,omitempty"`
	Fingerprint string `json:"fingerprint"`
	Version     int    `json:"version,omitempty"`
}

type UploadSignedPreKeyRequest struct {
	DeviceUUID string `json:"device_uuid"`
	KeyID      int    `json:"key_id"`
	PublicKey  string `json:"public_key"`
	Signature  string `json:"signature"`
	Algorithm  string `json:"algorithm,omitempty"`
	Version    int    `json:"version,omitempty"`
}

type UploadOneTimePreKeyItem struct {
	KeyID     int    `json:"key_id"`
	PublicKey string `json:"public_key"`
	Algorithm string `json:"algorithm,omitempty"`
	Version   int    `json:"version,omitempty"`
}

type UploadOneTimePreKeysRequest struct {
	DeviceUUID string                    `json:"device_uuid"`
	PreKeys    []UploadOneTimePreKeyItem `json:"prekeys"`
}

type E2EEKeyUploadResponse struct {
	Message string `json:"message"`
	ID      int64  `json:"id,omitempty"`
	Count   int    `json:"count,omitempty"`
}

type IdentityKeyResponse struct {
	KeyID       int64  `json:"key_id"`
	DeviceID    int64  `json:"device_id"`
	PublicKey   string `json:"public_key"`
	Algorithm   string `json:"algorithm"`
	Fingerprint string `json:"fingerprint"`
	Version     int    `json:"version"`
}

type SignedPreKeyResponse struct {
	KeyID     int64  `json:"key_id"`
	DeviceID  int64  `json:"device_id"`
	PreKeyID  int    `json:"prekey_id"`
	PublicKey string `json:"public_key"`
	Signature string `json:"signature"`
	Algorithm string `json:"algorithm"`
	Version   int    `json:"version"`
}

type OneTimePreKeyResponse struct {
	KeyID     int64  `json:"key_id"`
	DeviceID  int64  `json:"device_id"`
	PreKeyID  int    `json:"prekey_id"`
	PublicKey string `json:"public_key"`
	Algorithm string `json:"algorithm"`
	Version   int    `json:"version"`
}

type DeviceKeyBundleResponse struct {
	DeviceID         int64                  `json:"device_id"`
	DeviceUUID       string                 `json:"device_uuid"`
	IdentityKey      IdentityKeyResponse    `json:"identity_key"`
	SignedPreKey     SignedPreKeyResponse   `json:"signed_prekey"`
	OneTimePreKey    *OneTimePreKeyResponse `json:"one_time_prekey,omitempty"`
	HasOneTimePreKey bool                   `json:"has_one_time_prekey"`
}

type UserKeyBundleResponse struct {
	UserID  int64                     `json:"user_id"`
	Devices []DeviceKeyBundleResponse `json:"devices"`
}
