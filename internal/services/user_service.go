package services

import (
	"context"

	"github.com/example/ag-directory-manager/internal/ldap"
)

type UserService struct {
	directory ldap.Client
}

func NewUserService(directory ldap.Client) *UserService {
	return &UserService{directory: directory}
}

func (s *UserService) List(ctx context.Context, search string, allowedOUs []string, limit, offset int) ([]ldap.DirectoryUser, error) {
	return s.directory.ListUsers(ctx, search, allowedOUs, limit, offset)
}
