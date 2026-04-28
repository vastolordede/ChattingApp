package repository

import (
	"chattingapp_be/internal/models"
	"context"
	"database/sql"
	"errors"
	"time"
)

type DeviceIdentityKeyRepository struct {
	db *sql.DB
}

func NewDeviceIdentityKeyRepository(db *sql.DB) *DeviceIdentityKeyRepository {
	return &DeviceIdentityKeyRepository{db: db}
}

func (r *DeviceIdentityKeyRepository) ReplaceActive(ctx context.Context, k *models.DeviceIdentityKey) (int64, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, `
		UPDATE device_identity_keys
		SET is_active = FALSE,
			revoked_at = $2
		WHERE device_id = $1
		  AND is_active = TRUE
	`, k.DeviceID, time.Now())
	if err != nil {
		return 0, err
	}

	query := `
		INSERT INTO device_identity_keys (
			device_id, public_key, algorithm, fingerprint,
			version, is_active, created_at, expired_at, revoked_at
		) VALUES (
			$1, $2, $3, $4,
			$5, TRUE, $6, $7, NULL
		)
		ON CONFLICT (device_id, fingerprint)
		DO UPDATE SET
			public_key = EXCLUDED.public_key,
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
		k.PublicKey,
		k.Algorithm,
		k.Fingerprint,
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

func (r *DeviceIdentityKeyRepository) GetActiveByDeviceID(ctx context.Context, deviceID int64) (*models.DeviceIdentityKey, error) {
	query := `
		SELECT id, device_id, public_key, algorithm, fingerprint,
		       version, is_active, created_at, expired_at, revoked_at
		FROM device_identity_keys
		WHERE device_id = $1
		  AND is_active = TRUE
		  AND revoked_at IS NULL
		ORDER BY created_at DESC
		LIMIT 1
	`

	var k models.DeviceIdentityKey
	err := r.db.QueryRowContext(ctx, query, deviceID).Scan(
		&k.ID,
		&k.DeviceID,
		&k.PublicKey,
		&k.Algorithm,
		&k.Fingerprint,
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
