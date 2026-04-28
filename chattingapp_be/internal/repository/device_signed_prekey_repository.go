package repository

import (
	"chattingapp_be/internal/models"
	"context"
	"database/sql"
	"errors"
	"time"
)

type DeviceSignedPreKeyRepository struct {
	db *sql.DB
}

func NewDeviceSignedPreKeyRepository(db *sql.DB) *DeviceSignedPreKeyRepository {
	return &DeviceSignedPreKeyRepository{db: db}
}

func (r *DeviceSignedPreKeyRepository) ReplaceActive(ctx context.Context, k *models.DeviceSignedPreKey) (int64, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, `
		UPDATE device_signed_prekeys
		SET is_active = FALSE,
			revoked_at = $2
		WHERE device_id = $1
		  AND is_active = TRUE
	`, k.DeviceID, time.Now())
	if err != nil {
		return 0, err
	}

	query := `
		INSERT INTO device_signed_prekeys (
			device_id, key_id, public_key, signature,
			algorithm, version, is_active, created_at,
			expired_at, revoked_at
		) VALUES (
			$1, $2, $3, $4,
			$5, $6, TRUE, $7,
			$8, NULL
		)
		ON CONFLICT (device_id, key_id)
		DO UPDATE SET
			public_key = EXCLUDED.public_key,
			signature = EXCLUDED.signature,
			algorithm = EXCLUDED.algorithm,
			version = EXCLUDED.version,
			is_active = TRUE,
			expired_at = EXCLUDED.expired_at,
			revoked_at = NULL
		RETURNING id
	`

	var id int64
	err = tx.QueryRowContext(
		ctx,
		query,
		k.DeviceID,
		k.KeyID,
		k.PublicKey,
		k.Signature,
		k.Algorithm,
		k.Version,
		time.Now(),
		k.ExpiredAt,
	).Scan(&id)
	if err != nil {
		return 0, err
	}

	if err := tx.Commit(); err != nil {
		return 0, err
	}

	return id, nil
}

func (r *DeviceSignedPreKeyRepository) GetActiveByDeviceID(ctx context.Context, deviceID int64) (*models.DeviceSignedPreKey, error) {
	query := `
		SELECT id, device_id, key_id, public_key, signature,
		       algorithm, version, is_active, created_at, expired_at, revoked_at
		FROM device_signed_prekeys
		WHERE device_id = $1
		  AND is_active = TRUE
		  AND revoked_at IS NULL
		ORDER BY created_at DESC
		LIMIT 1
	`

	var k models.DeviceSignedPreKey
	err := r.db.QueryRowContext(ctx, query, deviceID).Scan(
		&k.ID,
		&k.DeviceID,
		&k.KeyID,
		&k.PublicKey,
		&k.Signature,
		&k.Algorithm,
		&k.Version,
		&k.IsActive,
		&k.CreatedAt,
		&k.ExpiredAt,
		&k.RevokedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &k, nil
}
