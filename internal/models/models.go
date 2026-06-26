package models

import "time"

type AppUser struct {
	ID           int64      `json:"id"`
	Username     string     `json:"username"`
	DisplayName  string     `json:"displayName"`
	Email        string     `json:"email"`
	RoleName     string     `json:"roleName"`
	RoleID       int64      `json:"roleId"`
	IsActive     bool       `json:"isActive"`
	BlockedUntil *time.Time `json:"blockedUntil,omitempty"`
	LastLoginAt  *time.Time `json:"lastLoginAt,omitempty"`
}

type AuditLog struct {
	ID            int64     `json:"id"`
	ActorUsername string    `json:"actorUsername"`
	IPAddress     string    `json:"ipAddress"`
	UserAgent     string    `json:"userAgent"`
	Operation     string    `json:"operation"`
	ObjectType    string    `json:"objectType"`
	ObjectID      string    `json:"objectId"`
	BeforePayload string    `json:"beforePayload"`
	AfterPayload  string    `json:"afterPayload"`
	Result        string    `json:"result"`
	DurationMS    int64     `json:"durationMs"`
	CreatedAt     time.Time `json:"createdAt"`
}

type DashboardSnapshot struct {
	UsersTotal        int64 `json:"usersTotal"`
	ComputersTotal    int64 `json:"computersTotal"`
	GroupsTotal       int64 `json:"groupsTotal"`
	OUsTotal          int64 `json:"ousTotal"`
	LockedUsers       int64 `json:"lockedUsers"`
	DisabledAccounts  int64 `json:"disabledAccounts"`
	InactiveComputers int64 `json:"inactiveComputers"`
	RecentChanges     int64 `json:"recentChanges"`
	RecentEvents      int64 `json:"recentEvents"`
}
