package tips

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/forbole/juno/v3/node"

	"github.com/desmos-labs/djuno/v2/x/contracts"

	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	"github.com/forbole/juno/v3/types/config"
	"google.golang.org/grpc"

	"github.com/desmos-labs/djuno/v2/database"
	contractsbase "github.com/desmos-labs/djuno/v2/x/contracts/base"
)

var (
	_ contracts.SmartContractModule = &Module{}
)

type Module struct {
	base *contractsbase.Module

	cfg        *Config
	db         *database.Db
	node       node.Node
	wasmClient wasmtypes.QueryClient
}

// NewModule returns a new Module instance
func NewModule(junoCfg config.Config, node node.Node, grpcConnection *grpc.ClientConn, db *database.Db) *Module {
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

// getContractConfig returns the configuration for the contract having the given address
func (m *Module) getContractConfig(address string) (*configResponse, error) {
	queryDataBz, err := json.Marshal(&ContractQuery{
		Config: &configQuery{},
	})
	if err != nil {
		return nil, fmt.Errorf("error while marshalling contract query: %s", err)
	}

	res, err := m.wasmClient.SmartContractState(context.Background(), &wasmtypes.QuerySmartContractStateRequest{
		Address:   address,
		QueryData: queryDataBz,
	})
	if err != nil {
		return nil, fmt.Errorf("error while querying contract state: %s", err)
	}

	var configRes configResponse
	err = json.Unmarshal(res.Data, &configRes)
	if err != nil {
		return nil, fmt.Errorf("error while unmarshalling contract config response: %s", err)
	}

	return &configRes, nil
}
