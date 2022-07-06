package types

import (
	"github.com/lib/pq"

	subspacestypes "github.com/desmos-labs/desmos/v4/x/subspaces/types"
)

// ConvertPermissions converts the given permissions into a pq.StringArray so that
// it can be properly inserted into the database
func ConvertPermissions(permissions subspacestypes.Permissions) pq.StringArray {
	return pq.StringArray(permissions)
}
