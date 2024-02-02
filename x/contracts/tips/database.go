package tips

import (
	"github.com/desmos-labs/athena/v2/types"
	contractsbase "github.com/desmos-labs/athena/v2/x/contracts/base"
)

type Database interface {
	contractsbase.Database
	HasPost(height int64, subspaceID uint64, postID uint64) (bool, error)
	SaveTip(tip types.Tip) error
}
