package builder

import (
	"github.com/forbole/juno/v5/node"
	"github.com/forbole/juno/v5/types/config"
	"google.golang.org/grpc"

	"github.com/desmos-labs/athena/x/contracts"
	"github.com/desmos-labs/athena/x/contracts/tips"
)

func BuildModule(junoCfg config.Config, node node.Node, grpcConnection *grpc.ClientConn, db tips.Database) *contracts.Module {
	return contracts.NewModule([]contracts.SmartContractModule{
		tips.NewModule(junoCfg, node, grpcConnection, db),
	})
}
