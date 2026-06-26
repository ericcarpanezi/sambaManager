export interface DashboardSnapshot {
  usersTotal: number;
  computersTotal: number;
  groupsTotal: number;
  ousTotal: number;
  lockedUsers: number;
  disabledAccounts: number;
  inactiveComputers: number;
  recentChanges: number;
  recentEvents: number;
}

export interface DirectoryUser {
  id: string;
  displayName: string;
  username: string;
  email: string;
  department: string;
  title: string;
  ou: string;
  enabled: boolean;
  lastLogonRaw: string;
}
