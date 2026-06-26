package services

import (
	"context"

	"github.com/example/ag-directory-manager/internal/ldap"
	"github.com/example/ag-directory-manager/internal/models"
	"github.com/example/ag-directory-manager/internal/repositories"
)

type DashboardService struct {
	directory ldap.Client
	dashboard *repositories.DashboardRepository
}

func NewDashboardService(directory ldap.Client, dashboard *repositories.DashboardRepository) *DashboardService {
	return &DashboardService{directory: directory, dashboard: dashboard}
}

func (s *DashboardService) Snapshot(ctx context.Context, allowedOUs []string) (models.DashboardSnapshot, error) {
	counters, err := s.directory.DashboardCounters(ctx, allowedOUs)
	if err != nil {
		return models.DashboardSnapshot{}, err
	}

	recentChanges, err := s.dashboard.CountRecentChanges(ctx)
	if err != nil {
		return models.DashboardSnapshot{}, err
	}
	recentEvents, err := s.dashboard.CountRecentEvents(ctx)
	if err != nil {
		return models.DashboardSnapshot{}, err
	}

	return models.DashboardSnapshot{
		UsersTotal:        counters.UsersTotal,
		ComputersTotal:    counters.ComputersTotal,
		GroupsTotal:       counters.GroupsTotal,
		OUsTotal:          counters.OUsTotal,
		LockedUsers:       counters.LockedUsers,
		DisabledAccounts:  counters.DisabledAccounts,
		InactiveComputers: counters.InactiveComputers,
		RecentChanges:     recentChanges,
		RecentEvents:      recentEvents,
	}, nil
}
