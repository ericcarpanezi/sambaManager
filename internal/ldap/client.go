package ldap

import (
	"context"
	"crypto/tls"
	"fmt"
	"strings"
	"time"

	ldapv3 "github.com/go-ldap/ldap/v3"
	"github.com/google/uuid"
)

type DashboardCounters struct {
	UsersTotal        int64
	ComputersTotal    int64
	GroupsTotal       int64
	OUsTotal          int64
	LockedUsers       int64
	DisabledAccounts  int64
	InactiveComputers int64
}

type DirectoryUser struct {
	ID           string `json:"id"`
	DisplayName  string `json:"displayName"`
	Username     string `json:"username"`
	Email        string `json:"email"`
	Department   string `json:"department"`
	Title        string `json:"title"`
	OU           string `json:"ou"`
	Enabled      bool   `json:"enabled"`
	LastLogonRaw string `json:"lastLogonRaw"`
}

type Client interface {
	Authenticate(ctx context.Context, username, password string) error
	DashboardCounters(ctx context.Context, allowedOUs []string) (DashboardCounters, error)
	ListUsers(ctx context.Context, search string, allowedOUs []string, limit, offset int) ([]DirectoryUser, error)
	TestConnection(ctx context.Context) error
}

type Config struct {
	URL                string
	BaseDN             string
	BindDN             string
	BindPassword       string
	StartTLS           bool
	InsecureSkipVerify bool
}

type RealClient struct{ cfg Config }

func NewRealClient(cfg Config) *RealClient { return &RealClient{cfg: cfg} }

func (c *RealClient) Authenticate(ctx context.Context, username, password string) error {
	conn, err := c.dial(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	searchReq := ldapv3.NewSearchRequest(
		c.cfg.BaseDN,
		ldapv3.ScopeWholeSubtree,
		ldapv3.NeverDerefAliases,
		1,
		0,
		false,
		fmt.Sprintf("(&(objectClass=user)(sAMAccountName=%s))", ldapv3.EscapeFilter(username)),
		[]string{"dn"},
		nil,
	)

	res, err := conn.Search(searchReq)
	if err != nil {
		return err
	}
	if len(res.Entries) == 0 {
		return fmt.Errorf("ldap user not found")
	}
	userDN := res.Entries[0].DN
	if err := conn.Bind(userDN, password); err != nil {
		return fmt.Errorf("ldap invalid credentials")
	}
	return nil
}

func (c *RealClient) DashboardCounters(context.Context, []string) (DashboardCounters, error) {
	return DashboardCounters{}, fmt.Errorf("dashboard counters via real ldap not implemented yet")
}

func (c *RealClient) ListUsers(context.Context, string, []string, int, int) ([]DirectoryUser, error) {
	return nil, fmt.Errorf("list users via real ldap not implemented yet")
}

func (c *RealClient) TestConnection(ctx context.Context) error {
	conn, err := c.dial(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()
	return nil
}

func (c *RealClient) dial(ctx context.Context) (*ldapv3.Conn, error) {
	_ = ctx
	conn, err := ldapv3.DialURL(c.cfg.URL)
	if err != nil {
		return nil, err
	}
	if c.cfg.StartTLS {
		tlsCfg := &tls.Config{InsecureSkipVerify: c.cfg.InsecureSkipVerify}
		if err := conn.StartTLS(tlsCfg); err != nil {
			conn.Close()
			return nil, err
		}
	}
	if c.cfg.BindDN != "" {
		if err := conn.Bind(c.cfg.BindDN, c.cfg.BindPassword); err != nil {
			conn.Close()
			return nil, err
		}
	}
	return conn, nil
}

type DemoClient struct{}

func NewDemoClient() *DemoClient { return &DemoClient{} }

func (d *DemoClient) Authenticate(context.Context, string, string) error {
	return nil
}

func (d *DemoClient) DashboardCounters(context.Context, []string) (DashboardCounters, error) {
	return DashboardCounters{
		UsersTotal:        142,
		ComputersTotal:    89,
		GroupsTotal:       28,
		OUsTotal:          12,
		LockedUsers:       2,
		DisabledAccounts:  5,
		InactiveComputers: 7,
	}, nil
}

func (d *DemoClient) ListUsers(_ context.Context, search string, allowedOUs []string, limit, offset int) ([]DirectoryUser, error) {
	if limit <= 0 {
		limit = 20
	}
	users := []DirectoryUser{
		{ID: uuid.NewString(), DisplayName: "Joao Santos", Username: "joao.santos", Email: "joao.santos@dominio.local", Department: "Financeiro", Title: "Analista", OU: "OU=Financeiro,DC=dominio,DC=local", Enabled: true, LastLogonRaw: time.Now().Add(-2 * time.Hour).Format(time.RFC3339)},
		{ID: uuid.NewString(), DisplayName: "Maria Costa", Username: "maria.costa", Email: "maria.costa@dominio.local", Department: "RH", Title: "Coordenadora", OU: "OU=RH,DC=dominio,DC=local", Enabled: true, LastLogonRaw: time.Now().Add(-36 * time.Hour).Format(time.RFC3339)},
		{ID: uuid.NewString(), DisplayName: "Carlos Lima", Username: "carlos.lima", Email: "carlos.lima@dominio.local", Department: "TI", Title: "Administrador", OU: "OU=TI,DC=dominio,DC=local", Enabled: false, LastLogonRaw: time.Now().Add(-72 * time.Hour).Format(time.RFC3339)},
	}

	out := make([]DirectoryUser, 0, len(users))
	for _, u := range users {
		if len(allowedOUs) > 0 {
			ok := false
			for _, allowed := range allowedOUs {
				if strings.Contains(strings.ToLower(u.OU), strings.ToLower(allowed)) {
					ok = true
					break
				}
			}
			if !ok {
				continue
			}
		}
		if search != "" {
			s := strings.ToLower(search)
			if !strings.Contains(strings.ToLower(u.DisplayName), s) && !strings.Contains(strings.ToLower(u.Username), s) {
				continue
			}
		}
		out = append(out, u)
	}

	if offset >= len(out) {
		return []DirectoryUser{}, nil
	}
	end := offset + limit
	if end > len(out) {
		end = len(out)
	}
	return out[offset:end], nil
}

func (d *DemoClient) TestConnection(context.Context) error {
	return nil
}
