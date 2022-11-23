package tips

import (
	"github.com/desmos-labs/djuno/v2/types"
	contractsbase "github.com/desmos-labs/djuno/v2/x/contracts/base"
)

type Database interface {
	contractsbase.Database
	SaveTip(tip types.Tip) error
}
