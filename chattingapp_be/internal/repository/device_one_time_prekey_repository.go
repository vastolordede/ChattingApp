package repository

import (
	"chattingapp_be/internal/models"
	"context"
	"database/sql"
	"errors"
	"time"
)

type DeviceOneTimePreKeyRepository struct {
	db *sql.DB
}

func NewDeviceOneTimePreKeyRepository(db *sql.DB) *DeviceOneTimePreKeyRepository {
	return &DeviceOneTimePreKeyRepository{db: db}
}

func (r *DeviceOneTimePreKeyRepository) CreateBatch(ctx context.Context, keys []models.DeviceOneTimePreKey) (int, error) {
	if len(keys) == 0 {
		return 0, nil
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	query := `
		INSERT INTO device_one_time_prekeys (
			device_id, key_id, public_key, algorithm,
			version, is_used, used_at, created_at, expired_at
		) VALUES (
			$1, $2, $3, $4,
			$5, FALSE, NULL, $6, $7
		)
		ON CONFLICT (device_id, key_id)
		DO UPDATE SET
			public_key = EXCLUDED.public_key,
			algorithm = EXCLUDED.algorithm,
			version = EXCLUDED.version,
			is_used = FALSE,
			used_at = NULL,
			expired_at = EXCLUDED.expired_at
	`

	count := 0
	for i := range keys {
		_, err := tx.ExecContext(
			ctx,
			query,
			keys[i].DeviceID,
			keys[i].KeyID,
			keys[i].PublicKey,
			keys[i].Algorithm,
			keys[i].Version,
			time.Now(),
			keys[i].ExpiredAt,
		)
		if err != nil {
			return 0, err
		}
		count++
	}

	if err := tx.Commit(); err != nil {
		return 0, err
	}

	return count, nil
}

func (r *DeviceOneTimePreKeyRepository) ConsumeOneByDeviceID(ctx context.Context, deviceID int64) (*models.DeviceOneTimePreKey, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	selectQuery := `
		SELECT id, device_id, key_id, public_key, algorithm,
		       version, is_used, used_at, created_at, expired_at
		FROM device_one_time_prekeys
		WHERE device_id = $1
		  AND is_used = FALSE
		  AND (expired_at IS NULL OR expired_at > NOW())
		ORDER BY key_id ASC
		LIMIT 1
		FOR UPDATE SKIP LOCKED
	`

	var k models.DeviceOneTimePreKey
	err = tx.QueryRowContext(ctx, selectQuery, deviceID).Scan(
		&k.ID,
		&k.DeviceID,
		&k.KeyID,
		&k.PublicKey,
		&k.Algorithm,
		&k.Version,
		&k.IsUsed,
		&k.UsedAt,
		&k.CreatedAt,
		&k.ExpiredAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	_, err = tx.ExecContext(ctx, `
		UPDATE device_one_time_prekeys
		SET is_used = TRUE,
			used_at = $2
		WHERE id = $1
	`, k.ID, time.Now())
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	k.IsUsed = true
	k.UsedAt = sql.NullTime{
		Time:  time.Now(),
		Valid: true,
	}

	return &k, nil
}

func (r *DeviceOneTimePreKeyRepository) CountUnusedByDeviceID(ctx context.Context, deviceID int64) (int, error) {
	query := `
		SELECT COUNT(*)
		FROM device_one_time_prekeys
		WHERE device_id = $1
		  AND is_used = FALSE
		  AND (expired_at IS NULL OR expired_at > NOW())
	`

	var total int
	err := r.db.QueryRowContext(ctx, query, deviceID).Scan(&total)
	if err != nil {
		return 0, err
	}

	return total, nil
}