package tests

import (
	"testing"

	"github.com/example/ag-directory-manager/internal/permissions"
)

func TestHasPermission(t *testing.T) {
	perms := []string{"user.view", "audit.view"}

	if !permissions.HasPermission(perms, "user.view") {
		t.Fatal("expected user.view to be present")
	}
	if permissions.HasPermission(perms, "user.delete") {
		t.Fatal("did not expect user.delete to be present")
	}
}
