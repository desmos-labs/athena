package builder

import (
	"github.com/forbole/juno/v3/node"
	"github.com/forbole/juno/v3/types/config"
	"google.golang.org/grpc"

	"github.com/desmos-labs/djuno/v2/database"
	"github.com/desmos-labs/djuno/v2/x/contracts"
	"github.com/desmos-labs/djuno/v2/x/contracts/tips"
)

func BuildModule(junoCfg config.Config, node node.Node, grpcConnection *grpc.ClientConn, db *database.Db) *contracts.Module {
	return contracts.NewModule([]contracts.SmartContractModule{
		tips.NewModule(junoCfg, node, grpcConnection, db),
	})
}
