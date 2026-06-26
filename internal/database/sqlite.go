package database

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	"github.com/example/ag-directory-manager/internal/config"
	_ "modernc.org/sqlite"
)

const schema = `
CREATE TABLE IF NOT EXISTS roles (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  name TEXT NOT NULL UNIQUE,
  description TEXT NOT NULL,
  is_system BOOLEAN NOT NULL DEFAULT 1
);

CREATE TABLE IF NOT EXISTS permissions (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  code TEXT NOT NULL UNIQUE,
  description TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS role_permissions (
  role_id INTEGER NOT NULL,
  permission_id INTEGER NOT NULL,
  PRIMARY KEY(role_id, permission_id),
  FOREIGN KEY(role_id) REFERENCES roles(id),
  FOREIGN KEY(permission_id) REFERENCES permissions(id)
);

CREATE TABLE IF NOT EXISTS app_users (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  username TEXT NOT NULL UNIQUE,
  display_name TEXT NOT NULL,
  email TEXT,
  role_id INTEGER NOT NULL,
  is_active BOOLEAN NOT NULL DEFAULT 1,
  blocked_until DATETIME,
  last_login_at DATETIME,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY(role_id) REFERENCES roles(id)
);

CREATE TABLE IF NOT EXISTS ou_scopes (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  user_id INTEGER NOT NULL,
  ou_dn TEXT NOT NULL,
  FOREIGN KEY(user_id) REFERENCES app_users(id)
);

CREATE TABLE IF NOT EXISTS refresh_tokens (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  user_id INTEGER NOT NULL,
  token_hash TEXT NOT NULL,
  expires_at DATETIME NOT NULL,
  revoked_at DATETIME,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY(user_id) REFERENCES app_users(id)
);

CREATE TABLE IF NOT EXISTS audit_logs (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  actor_username TEXT NOT NULL,
  ip_address TEXT,
  user_agent TEXT,
  operation TEXT NOT NULL,
  object_type TEXT NOT NULL,
  object_id TEXT,
  before_payload TEXT,
  after_payload TEXT,
  result TEXT NOT NULL,
  duration_ms INTEGER,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS settings (
  key TEXT PRIMARY KEY,
  value TEXT NOT NULL,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);
`

var permissionSeed = []string{
	"user.view", "user.create", "user.edit", "user.delete", "user.move",
	"user.password.change", "user.password.reset", "user.enable", "user.disable", "user.unlock",
	"computer.view", "computer.move", "computer.rename", "computer.delete",
	"group.manage", "ou.manage", "audit.view", "settings.manage",
}

type roleSeed struct {
	Name        string
	Description string
	Permissions []string
}

var defaultRoles = []roleSeed{
	{Name: "Administrador", Description: "Acesso total ao sistema", Permissions: permissionSeed},
	{Name: "Supervisor", Description: "Supervisão de operações", Permissions: []string{"user.view", "user.edit", "user.password.reset", "computer.view", "group.manage", "ou.manage", "audit.view"}},
	{Name: "Help Desk", Description: "Suporte operacional", Permissions: []string{"user.view", "user.password.reset", "user.unlock", "computer.view", "computer.rename"}},
	{Name: "RH", Description: "Operações de recursos humanos", Permissions: []string{"user.view", "user.create", "user.edit", "user.move"}},
	{Name: "Auditor", Description: "Somente auditoria", Permissions: []string{"audit.view", "user.view", "computer.view"}},
	{Name: "Somente Leitura", Description: "Acesso apenas leitura", Permissions: []string{"user.view", "computer.view"}},
}

func OpenAndMigrate(cfg config.Config) (*sql.DB, error) {
	if err := os.MkdirAll(filepath.Dir(cfg.SQLitePath), 0o755); err != nil {
		return nil, fmt.Errorf("create sqlite directory: %w", err)
	}

	db, err := sql.Open("sqlite", cfg.SQLitePath)
	if err != nil {
		return nil, fmt.Errorf("open sqlite: %w", err)
	}

	if _, err = db.Exec(schema); err != nil {
		return nil, fmt.Errorf("run schema migration: %w", err)
	}

	if err = seedPermissions(db); err != nil {
		return nil, err
	}
	if err = seedRoles(db); err != nil {
		return nil, err
	}
	if err = seedAdminUser(db); err != nil {
		return nil, err
	}

	return db, nil
}

func seedPermissions(db *sql.DB) error {
	for _, code := range permissionSeed {
		if _, err := db.Exec(`INSERT OR IGNORE INTO permissions(code, description) VALUES (?, ?)`, code, code); err != nil {
			return fmt.Errorf("seed permission %s: %w", code, err)
		}
	}
	return nil
}

func seedRoles(db *sql.DB) error {
	for _, role := range defaultRoles {
		if _, err := db.Exec(`INSERT OR IGNORE INTO roles(name, description, is_system) VALUES (?, ?, 1)`, role.Name, role.Description); err != nil {
			return fmt.Errorf("seed role %s: %w", role.Name, err)
		}

		var roleID int64
		if err := db.QueryRow(`SELECT id FROM roles WHERE name = ?`, role.Name).Scan(&roleID); err != nil {
			return fmt.Errorf("fetch role id %s: %w", role.Name, err)
		}

		for _, permission := range role.Permissions {
			var permissionID int64
			if err := db.QueryRow(`SELECT id FROM permissions WHERE code = ?`, permission).Scan(&permissionID); err != nil {
				return fmt.Errorf("fetch permission id %s: %w", permission, err)
			}
			if _, err := db.Exec(`INSERT OR IGNORE INTO role_permissions(role_id, permission_id) VALUES (?, ?)`, roleID, permissionID); err != nil {
				return fmt.Errorf("seed role permission %s->%s: %w", role.Name, permission, err)
			}
		}
	}
	return nil
}

func seedAdminUser(db *sql.DB) error {
	var roleID int64
	if err := db.QueryRow(`SELECT id FROM roles WHERE name = 'Administrador'`).Scan(&roleID); err != nil {
		return fmt.Errorf("load administrator role: %w", err)
	}
	_, err := db.Exec(`
		INSERT OR IGNORE INTO app_users(username, display_name, email, role_id, is_active)
		VALUES ('operador.demo', 'Operador Demo', 'operador.demo@local', ?, 1)
	`, roleID)
	if err != nil {
		return fmt.Errorf("seed default app user: %w", err)
	}
	return nil
}
