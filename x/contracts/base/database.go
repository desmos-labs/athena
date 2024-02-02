package contracts

import (
	"github.com/desmos-labs/athena/v2/types"
)

type Database interface {
	SaveContract(contract types.Contract) error
	GetContract(address string) (*types.Contract, error)
}
