package tips

import (
	"github.com/forbole/juno/v5/node"

	"github.com/desmos-labs/athena/x/contracts"

	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	"github.com/forbole/juno/v5/types/config"
	"google.golang.org/grpc"

	contractsbase "github.com/desmos-labs/athena/x/contracts/base"
)

var (
	_ contracts.SmartContractModule = &Module{}
)

type Module struct {
	base *contractsbase.Module

	cfg        *Config
	db         Database
	node       node.Node
	wasmClient wasmtypes.QueryClient
}

// NewModule returns a new Module instance
func NewModule(junoCfg config.Config, node node.Node, grpcConnection *grpc.ClientConn, db Database) *Module {
	bz, err := junoCfg.GetBytes()
	if err != nil {
		panic(err)
	}

	cfg, err := ParseConfig(bz)
	if err != nil {
		panic(err)
	}

	if cfg == nil {
		return nil
	}

	wasmClient := wasmtypes.NewQueryClient(grpcConnection)
	return &Module{
		base:       contractsbase.NewModule(wasmClient, db),
		cfg:        cfg,
		db:         db,
		node:       node,
		wasmClient: wasmClient,
	}
}

// Name implements modules.Module
func (m *Module) Name() string {
	return "tips"
}
