package fees

import "github.com/desmos-labs/djuno/v2/types"

type Database interface {
	SaveFeesParams(params types.FeesParams) error
}
