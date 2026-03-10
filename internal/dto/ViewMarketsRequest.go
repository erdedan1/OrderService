package dto

type ViewMarketsRequest struct {
	UserRoles []string
}

func NewViewMarketsRequestFromRoles(roles []string) *ViewMarketsRequest {
	copied := make([]string, len(roles))
	copy(copied, roles)

	return &ViewMarketsRequest{UserRoles: copied}
}
