package permissions

func HasPermission(permissions []string, wanted string) bool {
	for _, p := range permissions {
		if p == wanted {
			return true
		}
	}
	return false
}
