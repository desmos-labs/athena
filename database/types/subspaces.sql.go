package types

import (
	subspacestypes "github.com/desmos-labs/desmos/v4/x/subspaces/types"
	"github.com/lib/pq"
)

// ConvertPermissions converts the given permissions into a pq.StringArray so that
// it can be properly inserted into the database
func ConvertPermissions(permissions subspacestypes.Permissions) pq.StringArray {
	values := make([]string, len(permissions))
	for i, permission := range permissions {
		values[i] = permission
	}
	return values
}
