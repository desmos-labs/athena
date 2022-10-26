package contracts

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	"github.com/forbole/juno/v3/node/remote"
	juno "github.com/forbole/juno/v3/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// QuerySmartContractState queries a generic smart contract state given its address and the query to perform
func (m *Module) QuerySmartContractState(height int64, address string, query interface{}) (wasmtypes.RawContractMessage, error) {
	queryBz, err := json.Marshal(query)
	if err != nil {
		return nil, fmt.Errorf("error while marshalling contract query request: %s", err)
	}

	res, err := m.wasmClient.SmartContractState(
		remote.GetHeightRequestContext(context.Background(), height),
		&wasmtypes.QuerySmartContractStateRequest{
			Address:   address,
			QueryData: queryBz,
		},
	)
	if err != nil {
		if status.Code(err) == codes.NotFound || strings.Contains(err.Error(), "not found") {
			return nil, nil
		}
		return nil, fmt.Errorf("error while querying contract state: %s", err)
	}

	return res.Data, nil
}

// ParseContractAddress returns the contract address that is instantiated with the provided transaction
func (m *Module) ParseContractAddress(tx *juno.Tx, index int) (string, error) {
	event, err := tx.FindEventByType(index, wasmtypes.EventTypeInstantiate)
	if err != nil {
		return "", fmt.Errorf("no event %s found", wasmtypes.EventTypeInstantiate)
	}
	address, err := tx.FindAttributeByKey(event, wasmtypes.AttributeKeyContractAddr)
	if err != nil {
		return "", fmt.Errorf("no %s attribute found", wasmtypes.AttributeKeyContractAddr)
	}
	return address, nil
}
