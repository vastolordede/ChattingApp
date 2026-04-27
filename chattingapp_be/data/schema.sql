-- Run this script on database: chatting_app

BEGIN;

-- =========================================================
-- PATH 1
-- Core social + chat + attachments + devices
-- =========================================================

-- =========================================================
-- 1) USERS
-- =========================================================
CREATE TABLE IF NOT EXISTS users (
    id              BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    username        VARCHAR(50)  NOT NULL UNIQUE,
    full_name       VARCHAR(100) NOT NULL,
    email           VARCHAR(100) NOT NULL UNIQUE,
    phone_number    VARCHAR(20)  NOT NULL UNIQUE,
    password_hash   VARCHAR(255) NOT NULL,
    avatar_url      TEXT NULL,
    bio             TEXT NULL,
    status          VARCHAR(20)  NOT NULL DEFAULT 'active',
    is_verified     BOOLEAN      NOT NULL DEFAULT FALSE,
    last_seen_at    TIMESTAMP NULL,
    created_at      TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT chk_users_status
        CHECK (status IN ('active', 'inactive', 'banned'))
);

-- =========================================================
-- 2) FRIEND REQUESTS
-- =========================================================
CREATE TABLE IF NOT EXISTS friend_requests (
    id              BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    status          VARCHAR(20) NOT NULL DEFAULT 'pending',
    message         TEXT NULL,
    created_at      TIMESTAMP   NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP   NOT NULL DEFAULT CURRENT_TIMESTAMP,
    responded_at    TIMESTAMP NULL,
    expired_at      TIMESTAMP NULL,
    CONSTRAINT chk_friend_requests_status
        CHECK (status IN ('pending', 'accepted', 'rejected', 'cancelled', 'expired'))
);

CREATE TABLE IF NOT EXISTS friend_request_members (
    id                  BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    friend_request_id   BIGINT      NOT NULL,
    user_id             BIGINT      NOT NULL,
    role                VARCHAR(20) NOT NULL,
    created_at          TIMESTAMP   NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_frm_friend_request
        FOREIGN KEY (friend_request_id) REFERENCES friend_requests(id) ON DELETE CASCADE,
    CONSTRAINT fk_frm_user
        FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT uq_frm_friend_request_user
        UNIQUE (friend_request_id, user_id),
    CONSTRAINT chk_frm_role
        CHECK (role IN ('sender', 'receiver'))
);

-- =========================================================
-- 3) FRIENDSHIPS
-- =========================================================
CREATE TABLE IF NOT EXISTS friendships (
    id              BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    status          VARCHAR(20) NOT NULL DEFAULT 'active',
    created_at      TIMESTAMP   NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP   NOT NULL DEFAULT CURRENT_TIMESTAMP,
    ended_at        TIMESTAMP NULL,
    CONSTRAINT chk_friendships_status
        CHECK (status IN ('active', 'blocked', 'deleted'))
);

CREATE TABLE IF NOT EXISTS friendship_members (
    id              BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    friendship_id   BIGINT    NOT NULL,
    user_id         BIGINT    NOT NULL,
    created_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_fsm_friendship
        FOREIGN KEY (friendship_id) REFERENCES friendships(id) ON DELETE CASCADE,
    CONSTRAINT fk_fsm_user
        FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT uq_fsm_friendship_user
        UNIQUE (friendship_id, user_id)
);

-- =========================================================
-- 4) USER DEVICES
-- Giữ từ Path 1 để sau này nối Path 2 luôn
-- =========================================================
CREATE TABLE IF NOT EXISTS user_devices (
    id              BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    user_id         BIGINT       NOT NULL,
    device_uuid     VARCHAR(100) NOT NULL UNIQUE,
    device_name     VARCHAR(100) NOT NULL,
    device_type     VARCHAR(20)  NOT NULL,
    platform        VARCHAR(20)  NOT NULL,
    app_version     VARCHAR(30) NULL,
    os_version      VARCHAR(30) NULL,
    push_token      TEXT NULL,
    is_trusted      BOOLEAN      NOT NULL DEFAULT FALSE,
    is_active       BOOLEAN      NOT NULL DEFAULT TRUE,
    last_seen_at    TIMESTAMP NULL,
    registered_at   TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_at      TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_user_devices_user
        FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT chk_user_devices_device_type
        CHECK (device_type IN ('phone', 'tablet', 'desktop', 'web')),
    CONSTRAINT chk_user_devices_platform
        CHECK (platform IN ('android', 'ios', 'windows', 'macos', 'linux', 'web'))
);

-- =========================================================
-- 5) CONVERSATIONS
-- last_message_id add FK later because of cycle
-- =========================================================
CREATE TABLE IF NOT EXISTS conversations (
    id                  BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    conversation_type   VARCHAR(20) NOT NULL,
    title               TEXT NULL,
    avatar_url          TEXT NULL,
    created_by_user_id  BIGINT NULL,
    status              VARCHAR(20) NOT NULL DEFAULT 'active',
    last_message_id     BIGINT NULL,
    last_message_at     TIMESTAMP NULL,
    created_at          TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at          TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_conversations_created_by
        FOREIGN KEY (created_by_user_id) REFERENCES users(id) ON DELETE SET NULL,
    CONSTRAINT chk_conversations_type
        CHECK (conversation_type IN ('direct', 'group')),
    CONSTRAINT chk_conversations_status
        CHECK (status IN ('active', 'archived', 'deleted'))
);

-- =========================================================
-- 6) MESSAGES
-- =========================================================
CREATE TABLE IF NOT EXISTS messages (
    id                          BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    conversation_id             BIGINT      NOT NULL,
    sender_user_id              BIGINT      NOT NULL,
    message_type                VARCHAR(20) NOT NULL,
    content                     TEXT NULL,
    reply_to_message_id         BIGINT NULL,
    forwarded_from_message_id   BIGINT NULL,
    status                      VARCHAR(20) NOT NULL DEFAULT 'sent',
    is_edited                   BOOLEAN     NOT NULL DEFAULT FALSE,
    edited_at                   TIMESTAMP NULL,
    is_deleted                  BOOLEAN     NOT NULL DEFAULT FALSE,
    deleted_at                  TIMESTAMP NULL,
    client_message_id           VARCHAR(100) NULL UNIQUE,
    sent_at                     TIMESTAMP   NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_at                  TIMESTAMP   NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at                  TIMESTAMP   NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_messages_conversation
        FOREIGN KEY (conversation_id) REFERENCES conversations(id) ON DELETE CASCADE,
    CONSTRAINT fk_messages_sender
        FOREIGN KEY (sender_user_id) REFERENCES users(id) ON DELETE RESTRICT,
    CONSTRAINT fk_messages_reply_to
        FOREIGN KEY (reply_to_message_id) REFERENCES messages(id) ON DELETE SET NULL,
    CONSTRAINT fk_messages_forwarded_from
        FOREIGN KEY (forwarded_from_message_id) REFERENCES messages(id) ON DELETE SET NULL,
    CONSTRAINT chk_messages_type
        CHECK (message_type IN ('text', 'image', 'video', 'file', 'audio', 'link', 'system')),
    CONSTRAINT chk_messages_status
        CHECK (status IN ('sending', 'sent', 'delivered', 'read', 'failed'))
);

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1
        FROM pg_constraint
        WHERE conname = 'fk_conversations_last_message'
    ) THEN
        ALTER TABLE conversations
        ADD CONSTRAINT fk_conversations_last_message
        FOREIGN KEY (last_message_id) REFERENCES messages(id) ON DELETE SET NULL;
    END IF;
END $$;

-- =========================================================
-- 7) CONVERSATION MEMBERS
-- =========================================================
CREATE TABLE IF NOT EXISTS conversation_members (
    id                      BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    conversation_id         BIGINT      NOT NULL,
    user_id                 BIGINT      NOT NULL,
    role                    VARCHAR(20) NOT NULL DEFAULT 'member',
    joined_at               TIMESTAMP   NOT NULL DEFAULT CURRENT_TIMESTAMP,
    left_at                 TIMESTAMP NULL,
    is_active               BOOLEAN     NOT NULL DEFAULT TRUE,
    nickname                TEXT NULL,
    last_read_message_id    BIGINT NULL,
    last_read_at            TIMESTAMP NULL,
    mute_until              TIMESTAMP NULL,
    created_at              TIMESTAMP   NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at              TIMESTAMP   NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_cm_conversation
        FOREIGN KEY (conversation_id) REFERENCES conversations(id) ON DELETE CASCADE,
    CONSTRAINT fk_cm_user
        FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_cm_last_read_message
        FOREIGN KEY (last_read_message_id) REFERENCES messages(id) ON DELETE SET NULL,
    CONSTRAINT uq_cm_conversation_user
        UNIQUE (conversation_id, user_id),
    CONSTRAINT chk_cm_role
        CHECK (role IN ('owner', 'admin', 'member'))
);

-- =========================================================
-- 8) MESSAGE ATTACHMENTS
-- =========================================================
CREATE TABLE IF NOT EXISTS message_attachments (
    id                  BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    message_id          BIGINT       NOT NULL,
    attachment_type     VARCHAR(20)  NOT NULL,
    file_name           VARCHAR(255) NOT NULL,
    mime_type           VARCHAR(100) NOT NULL,
    file_size           BIGINT       NOT NULL,
    file_url            TEXT         NOT NULL,
    thumbnail_url       TEXT NULL,
    width               INTEGER NULL,
    height              INTEGER NULL,
    duration_seconds    INTEGER NULL,
    checksum            VARCHAR(255) NULL,
    encryption_key_hint VARCHAR(255) NULL,
    created_at          TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_message_attachments_message
        FOREIGN KEY (message_id) REFERENCES messages(id) ON DELETE CASCADE,
    CONSTRAINT chk_message_attachments_type
        CHECK (attachment_type IN ('image', 'video', 'file', 'audio', 'link'))
);

-- =========================================================
-- 9) INDEXES FOR PATH 1
-- =========================================================
CREATE INDEX IF NOT EXISTS idx_frm_user_id
    ON friend_request_members(user_id);

CREATE INDEX IF NOT EXISTS idx_fsm_user_id
    ON friendship_members(user_id);

CREATE INDEX IF NOT EXISTS idx_user_devices_user_id
    ON user_devices(user_id);

CREATE INDEX IF NOT EXISTS idx_conversations_created_by_user_id
    ON conversations(created_by_user_id);

CREATE INDEX IF NOT EXISTS idx_conversations_last_message_at
    ON conversations(last_message_at DESC);

CREATE INDEX IF NOT EXISTS idx_cm_user_id
    ON conversation_members(user_id);

CREATE INDEX IF NOT EXISTS idx_cm_last_read_message_id
    ON conversation_members(last_read_message_id);

CREATE INDEX IF NOT EXISTS idx_messages_conversation_id_sent_at
    ON messages(conversation_id, sent_at DESC);

CREATE INDEX IF NOT EXISTS idx_messages_sender_user_id
    ON messages(sender_user_id);

CREATE INDEX IF NOT EXISTS idx_messages_reply_to_message_id
    ON messages(reply_to_message_id);

CREATE INDEX IF NOT EXISTS idx_message_attachments_message_id
    ON message_attachments(message_id);

COMMIT;
BEGIN;

-- =========================================================
-- 10) USER REFRESH TOKENS
-- JWT refresh token theo user + device
-- =========================================================
CREATE TABLE IF NOT EXISTS user_refresh_tokens (
    id                  BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    user_id             BIGINT NOT NULL,
    user_device_id      BIGINT NULL,
    token_hash          VARCHAR(255) NOT NULL,
    expires_at          TIMESTAMP NOT NULL,
    revoked_at          TIMESTAMP NULL,
    last_used_at        TIMESTAMP NULL,
    created_at          TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at          TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_urt_user
        FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_urt_user_device
        FOREIGN KEY (user_device_id) REFERENCES user_devices(id) ON DELETE SET NULL
);

CREATE INDEX IF NOT EXISTS idx_urt_user_id
    ON user_refresh_tokens(user_id);

CREATE INDEX IF NOT EXISTS idx_urt_user_device_id
    ON user_refresh_tokens(user_device_id);

CREATE INDEX IF NOT EXISTS idx_urt_expires_at
    ON user_refresh_tokens(expires_at);

CREATE UNIQUE INDEX IF NOT EXISTS uq_urt_token_hash
    ON user_refresh_tokens(token_hash);

-- =========================================================
-- 11) USER BLOCKS
-- Block 1 chiều: A block B
-- =========================================================
CREATE TABLE IF NOT EXISTS user_blocks (
    id                  BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    blocker_user_id     BIGINT NOT NULL,
    blocked_user_id     BIGINT NOT NULL,
    created_at          TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_user_blocks_blocker
        FOREIGN KEY (blocker_user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_user_blocks_blocked
        FOREIGN KEY (blocked_user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT uq_user_blocks_pair
        UNIQUE (blocker_user_id, blocked_user_id),
    CONSTRAINT chk_user_blocks_not_self
        CHECK (blocker_user_id <> blocked_user_id)
);

CREATE INDEX IF NOT EXISTS idx_user_blocks_blocker
    ON user_blocks(blocker_user_id);

CREATE INDEX IF NOT EXISTS idx_user_blocks_blocked
    ON user_blocks(blocked_user_id);

-- =========================================================
-- 12) MESSAGE REACTIONS
-- Mỗi user chỉ có 1 reaction trên 1 message
-- =========================================================
CREATE TABLE IF NOT EXISTS message_reactions (
    id                  BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    message_id          BIGINT NOT NULL,
    user_id             BIGINT NOT NULL,
    reaction_type       VARCHAR(30) NOT NULL,
    created_at          TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at          TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_message_reactions_message
        FOREIGN KEY (message_id) REFERENCES messages(id) ON DELETE CASCADE,
    CONSTRAINT fk_message_reactions_user
        FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT uq_message_reactions_message_user
        UNIQUE (message_id, user_id),
    CONSTRAINT chk_message_reactions_type
        CHECK (reaction_type IN ('like', 'love', 'haha', 'wow', 'sad', 'angry'))
);

CREATE INDEX IF NOT EXISTS idx_message_reactions_message_id
    ON message_reactions(message_id);

CREATE INDEX IF NOT EXISTS idx_message_reactions_user_id
    ON message_reactions(user_id);

COMMIT;
CREATE TABLE public.password_reset_tokens (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    token_hash VARCHAR(255) NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    used_at TIMESTAMP NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_password_reset_tokens_user
        FOREIGN KEY (user_id)
        REFERENCES public.users(id)
        ON DELETE CASCADE
);
CREATE INDEX idx_password_reset_tokens_user_id
    ON public.password_reset_tokens(user_id);

CREATE INDEX idx_password_reset_tokens_token_hash
    ON public.password_reset_tokens(token_hash);

CREATE INDEX idx_password_reset_tokens_expires_at
    ON public.password_reset_tokens(expires_at);
ALTER TABLE public.password_reset_tokens
ADD CONSTRAINT uq_password_reset_tokens_token_hash UNIQUE (token_hash);
INSERT INTO user_devices (
   user_id,
   device_uuid,
   device_name,
   device_type,
   platform,
   app_version,
   os_version,
   push_token,
   is_trusted,
   is_active,
   last_seen_at,
   registered_at,
   created_at,
   updated_at
) VALUES (
   1,                                -- nhớ phải có user_id tồn tại trong bảng users
   'device-001',
   'My iPhone',
   'phone',                          -- ⚠️ QUAN TRỌNG
   'ios',                            -- ⚠️ QUAN TRỌNG
   '1.0',
   '16.0',
   'test_token_123',
   false,
   true,
   NOW(),
   NOW(),
   NOW(),
   NOW()
);

ALTER TABLE conversation_members
ADD COLUMN is_pinned BOOLEAN DEFAULT FALSE,
ADD COLUMN is_archived BOOLEAN DEFAULT FALSE;
ALTER TABLE messages
ADD COLUMN IF NOT EXISTS is_recalled BOOLEAN NOT NULL DEFAULT false;

ALTER TABLE messages
ADD COLUMN IF NOT EXISTS recalled_at TIMESTAMP NULL;
-- Run this script on database: chatting_app

BEGIN;

-- =========================================================
-- PART 2: E2EE LAYER
-- Requires Part 1 tables:
-- users, user_devices, conversations, conversation_members, messages
-- =========================================================

-- =========================================================
-- 1) DEVICE IDENTITY KEYS
-- One active identity public key per user device version/fingerprint
-- =========================================================
CREATE TABLE IF NOT EXISTS device_identity_keys (
    id              BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    device_id       BIGINT       NOT NULL,
    public_key      TEXT         NOT NULL,
    algorithm       VARCHAR(30)  NOT NULL DEFAULT 'X25519',
    fingerprint     VARCHAR(255) NOT NULL,
    version         INTEGER      NOT NULL DEFAULT 1,
    is_active       BOOLEAN      NOT NULL DEFAULT TRUE,
    created_at      TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    expired_at      TIMESTAMP NULL,
    revoked_at      TIMESTAMP NULL,

    CONSTRAINT fk_device_identity_keys_device
        FOREIGN KEY (device_id)
        REFERENCES user_devices(id)
        ON DELETE CASCADE,

    CONSTRAINT uq_device_identity_keys_device_fingerprint
        UNIQUE (device_id, fingerprint)
);

CREATE INDEX IF NOT EXISTS idx_device_identity_keys_device_id
    ON device_identity_keys(device_id);

CREATE INDEX IF NOT EXISTS idx_device_identity_keys_device_active
    ON device_identity_keys(device_id, is_active);


-- =========================================================
-- 2) DEVICE SIGNED PREKEYS
-- Signed prekey for each device, signed by device identity key
-- =========================================================
CREATE TABLE IF NOT EXISTS device_signed_prekeys (
    id              BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    device_id       BIGINT      NOT NULL,
    key_id          INTEGER     NOT NULL,
    public_key      TEXT        NOT NULL,
    signature       TEXT        NOT NULL,
    algorithm       VARCHAR(30) NOT NULL DEFAULT 'X25519',
    version         INTEGER     NOT NULL DEFAULT 1,
    is_active       BOOLEAN     NOT NULL DEFAULT TRUE,
    created_at      TIMESTAMP   NOT NULL DEFAULT CURRENT_TIMESTAMP,
    expired_at      TIMESTAMP NULL,
    revoked_at      TIMESTAMP NULL,

    CONSTRAINT fk_device_signed_prekeys_device
        FOREIGN KEY (device_id)
        REFERENCES user_devices(id)
        ON DELETE CASCADE,

    CONSTRAINT uq_device_signed_prekeys_device_key
        UNIQUE (device_id, key_id)
);

CREATE INDEX IF NOT EXISTS idx_device_signed_prekeys_device_id
    ON device_signed_prekeys(device_id);

CREATE INDEX IF NOT EXISTS idx_device_signed_prekeys_device_active
    ON device_signed_prekeys(device_id, is_active);


-- =========================================================
-- 3) DEVICE ONE-TIME PREKEYS
-- One-time prekeys are consumed once during session setup
-- =========================================================
CREATE TABLE IF NOT EXISTS device_one_time_prekeys (
    id              BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    device_id       BIGINT      NOT NULL,
    key_id          INTEGER     NOT NULL,
    public_key      TEXT        NOT NULL,
    algorithm       VARCHAR(30) NOT NULL DEFAULT 'X25519',
    version         INTEGER     NOT NULL DEFAULT 1,
    is_used         BOOLEAN     NOT NULL DEFAULT FALSE,
    used_at         TIMESTAMP NULL,
    created_at      TIMESTAMP   NOT NULL DEFAULT CURRENT_TIMESTAMP,
    expired_at      TIMESTAMP NULL,

    CONSTRAINT fk_device_one_time_prekeys_device
        FOREIGN KEY (device_id)
        REFERENCES user_devices(id)
        ON DELETE CASCADE,

    CONSTRAINT uq_device_one_time_prekeys_device_key
        UNIQUE (device_id, key_id)
);

CREATE INDEX IF NOT EXISTS idx_device_one_time_prekeys_device_unused
    ON device_one_time_prekeys(device_id, is_used);

CREATE INDEX IF NOT EXISTS idx_device_one_time_prekeys_device_key
    ON device_one_time_prekeys(device_id, key_id);


-- =========================================================
-- 4) MESSAGE CIPHERTEXTS
-- Store encrypted payload per target device.
-- messages.content should be NULL for encrypted messages.
-- =========================================================
CREATE TABLE IF NOT EXISTS message_ciphertexts (
    id                  BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    message_id          BIGINT      NOT NULL,
    target_device_id    BIGINT      NOT NULL,
    sender_device_id    BIGINT NULL,
    ciphertext          TEXT        NOT NULL,
    encryption_header   TEXT NULL,
    nonce               TEXT NULL,
    algorithm           VARCHAR(30) NOT NULL DEFAULT 'XCHACHA20_POLY1305',
    message_version     INTEGER     NOT NULL DEFAULT 1,
    is_delivered        BOOLEAN     NOT NULL DEFAULT FALSE,
    delivered_at        TIMESTAMP NULL,
    created_at          TIMESTAMP   NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_message_ciphertexts_message
        FOREIGN KEY (message_id)
        REFERENCES messages(id)
        ON DELETE CASCADE,

    CONSTRAINT fk_message_ciphertexts_target_device
        FOREIGN KEY (target_device_id)
        REFERENCES user_devices(id)
        ON DELETE CASCADE,

    CONSTRAINT fk_message_ciphertexts_sender_device
        FOREIGN KEY (sender_device_id)
        REFERENCES user_devices(id)
        ON DELETE SET NULL,

    CONSTRAINT uq_message_ciphertexts_message_target
        UNIQUE (message_id, target_device_id)
);

CREATE INDEX IF NOT EXISTS idx_message_ciphertexts_message_id
    ON message_ciphertexts(message_id);

CREATE INDEX IF NOT EXISTS idx_message_ciphertexts_target_device_id
    ON message_ciphertexts(target_device_id);

CREATE INDEX IF NOT EXISTS idx_message_ciphertexts_sender_device_id
    ON message_ciphertexts(sender_device_id);

CREATE INDEX IF NOT EXISTS idx_message_ciphertexts_target_delivered
    ON message_ciphertexts(target_device_id, is_delivered);

COMMIT;

