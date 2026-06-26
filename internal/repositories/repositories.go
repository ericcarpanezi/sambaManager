package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/example/ag-directory-manager/internal/models"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository { return &UserRepository{db: db} }

func (r *UserRepository) GetByUsername(ctx context.Context, username string) (models.AppUser, error) {
	const q = `
	SELECT u.id, u.username, u.display_name, u.email, u.role_id, r.name, u.is_active, u.blocked_until, u.last_login_at
	FROM app_users u
	JOIN roles r ON r.id = u.role_id
	WHERE u.username = ?
	`
	row := r.db.QueryRowContext(ctx, q, username)

	var user models.AppUser
	var blockedUntil, lastLoginAt sql.NullTime
	if err := row.Scan(
		&user.ID,
		&user.Username,
		&user.DisplayName,
		&user.Email,
		&user.RoleID,
		&user.RoleName,
		&user.IsActive,
		&blockedUntil,
		&lastLoginAt,
	); err != nil {
		if err == sql.ErrNoRows {
			return models.AppUser{}, ErrNotFound
		}
		return models.AppUser{}, err
	}
	if blockedUntil.Valid {
		v := blockedUntil.Time
		user.BlockedUntil = &v
	}
	if lastLoginAt.Valid {
		v := lastLoginAt.Time
		user.LastLoginAt = &v
	}
	return user, nil
}

func (r *UserRepository) GetByID(ctx context.Context, userID int64) (models.AppUser, error) {
	const q = `
	SELECT u.id, u.username, u.display_name, u.email, u.role_id, r.name, u.is_active, u.blocked_until, u.last_login_at
	FROM app_users u
	JOIN roles r ON r.id = u.role_id
	WHERE u.id = ?
	`
	row := r.db.QueryRowContext(ctx, q, userID)

	var user models.AppUser
	var blockedUntil, lastLoginAt sql.NullTime
	if err := row.Scan(
		&user.ID,
		&user.Username,
		&user.DisplayName,
		&user.Email,
		&user.RoleID,
		&user.RoleName,
		&user.IsActive,
		&blockedUntil,
		&lastLoginAt,
	); err != nil {
		if err == sql.ErrNoRows {
			return models.AppUser{}, ErrNotFound
		}
		return models.AppUser{}, err
	}
	if blockedUntil.Valid {
		v := blockedUntil.Time
		user.BlockedUntil = &v
	}
	if lastLoginAt.Valid {
		v := lastLoginAt.Time
		user.LastLoginAt = &v
	}
	return user, nil
}

func (r *UserRepository) UpdateLastLogin(ctx context.Context, userID int64) error {
	_, err := r.db.ExecContext(ctx, `UPDATE app_users SET last_login_at = CURRENT_TIMESTAMP, updated_at = CURRENT_TIMESTAMP WHERE id = ?`, userID)
	return err
}

func (r *UserRepository) ListOUScopes(ctx context.Context, userID int64) ([]string, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT ou_dn FROM ou_scopes WHERE user_id = ? ORDER BY ou_dn ASC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []string
	for rows.Next() {
		var ou string
		if err := rows.Scan(&ou); err != nil {
			return nil, err
		}
		out = append(out, ou)
	}
	return out, rows.Err()
}

func (r *UserRepository) ListPermissions(ctx context.Context, userID int64) ([]string, error) {
	const q = `
	SELECT p.code
	FROM app_users u
	JOIN role_permissions rp ON rp.role_id = u.role_id
	JOIN permissions p ON p.id = rp.permission_id
	WHERE u.id = ?
	ORDER BY p.code ASC
	`
	rows, err := r.db.QueryContext(ctx, q, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	permissions := make([]string, 0, 8)
	for rows.Next() {
		var code string
		if err := rows.Scan(&code); err != nil {
			return nil, err
		}
		permissions = append(permissions, code)
	}
	return permissions, rows.Err()
}

type AuditRepository struct{ db *sql.DB }

func NewAuditRepository(db *sql.DB) *AuditRepository { return &AuditRepository{db: db} }

func (r *AuditRepository) Insert(ctx context.Context, log models.AuditLog) error {
	const q = `
	INSERT INTO audit_logs(actor_username, ip_address, user_agent, operation, object_type, object_id, before_payload, after_payload, result, duration_ms)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	_, err := r.db.ExecContext(ctx, q,
		log.ActorUsername,
		log.IPAddress,
		log.UserAgent,
		log.Operation,
		log.ObjectType,
		log.ObjectID,
		log.BeforePayload,
		log.AfterPayload,
		log.Result,
		log.DurationMS,
	)
	return err
}

func (r *AuditRepository) List(ctx context.Context, limit int) ([]models.AuditLog, error) {
	if limit <= 0 || limit > 500 {
		limit = 100
	}
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, actor_username, ip_address, user_agent, operation, object_type, object_id, before_payload, after_payload, result, duration_ms, created_at
		FROM audit_logs
		ORDER BY id DESC
		LIMIT ?
	`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	logs := make([]models.AuditLog, 0, limit)
	for rows.Next() {
		var l models.AuditLog
		if err := rows.Scan(&l.ID, &l.ActorUsername, &l.IPAddress, &l.UserAgent, &l.Operation, &l.ObjectType, &l.ObjectID, &l.BeforePayload, &l.AfterPayload, &l.Result, &l.DurationMS, &l.CreatedAt); err != nil {
			return nil, err
		}
		logs = append(logs, l)
	}
	return logs, rows.Err()
}

type TokenRepository struct{ db *sql.DB }

func NewTokenRepository(db *sql.DB) *TokenRepository { return &TokenRepository{db: db} }

func (r *TokenRepository) StoreRefreshToken(ctx context.Context, userID int64, tokenHash string, expiresAt time.Time) error {
	_, err := r.db.ExecContext(ctx, `INSERT INTO refresh_tokens(user_id, token_hash, expires_at) VALUES (?, ?, ?)`, userID, tokenHash, expiresAt.UTC())
	return err
}

func (r *TokenRepository) ValidateRefreshToken(ctx context.Context, userID int64, tokenHash string) error {
	var count int
	err := r.db.QueryRowContext(ctx, `
		SELECT COUNT(1)
		FROM refresh_tokens
		WHERE user_id = ? AND token_hash = ? AND revoked_at IS NULL AND expires_at > CURRENT_TIMESTAMP
	`, userID, tokenHash).Scan(&count)
	if err != nil {
		return err
	}
	if count == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *TokenRepository) RevokeRefreshToken(ctx context.Context, userID int64, tokenHash string) error {
	res, err := r.db.ExecContext(ctx, `UPDATE refresh_tokens SET revoked_at = CURRENT_TIMESTAMP WHERE user_id = ? AND token_hash = ?`, userID, tokenHash)
	if err != nil {
		return err
	}
	if rows, _ := res.RowsAffected(); rows == 0 {
		return ErrNotFound
	}
	return nil
}

type DashboardRepository struct{ db *sql.DB }

func NewDashboardRepository(db *sql.DB) *DashboardRepository { return &DashboardRepository{db: db} }

func (r *DashboardRepository) CountRecentChanges(ctx context.Context) (int64, error) {
	var n int64
	err := r.db.QueryRowContext(ctx, `SELECT COUNT(1) FROM audit_logs WHERE created_at > DATETIME('now', '-24 hour')`).Scan(&n)
	return n, err
}

func (r *DashboardRepository) CountRecentEvents(ctx context.Context) (int64, error) {
	var n int64
	err := r.db.QueryRowContext(ctx, `SELECT COUNT(1) FROM audit_logs WHERE created_at > DATETIME('now', '-1 hour')`).Scan(&n)
	return n, err
}

var ErrNotFound = fmt.Errorf("not found")
