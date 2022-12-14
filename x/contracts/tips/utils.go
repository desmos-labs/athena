package tips

import (
	"encoding/json"
	"fmt"
)

// getContractConfig returns the configuration for the contract having the given address at the specified height.
// If the config at the specified height does not exist, returns the one stored inside the database
func (m *Module) getContractConfig(height int64, address string) (*configResponse, error) {
	res, err := m.base.QuerySmartContractState(height, address, &tipsContractQuery{
		Config: &configQuery{},
	})
	if err != nil {
		return nil, fmt.Errorf("error while querying contract state: %s", err)
	}

	if res == nil {
		contract, err := m.db.GetContract(address)
		if err != nil {
			return nil, nil
		}
		if contract == nil {
			return nil, fmt.Errorf("no contract found inside the database nor on the chain: %s", address)
		}
		res = contract.ConfigBz
	}

	var config configResponse
	err = json.Unmarshal(res, &config)
	if err != nil {
		return nil, fmt.Errorf("error while unmarshaling tips contract config response: %s", err)
	}
	return &config, nil
}
