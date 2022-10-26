package tips

import (
	"encoding/json"
	"fmt"
)

// getContractConfig returns the configuration for the contract having the given address at the specified height
func (m *Module) getContractConfig(height int64, address string) (*configResponse, error) {
	res, err := m.base.QuerySmartContractState(height, address, &tipsContractQuery{
		Config: &configQuery{},
	})
	if err != nil {
		return nil, fmt.Errorf("error while querying contract state: %s", err)
	}

	var config configResponse
	err = json.Unmarshal(res, &config)
	if err != nil {
		return nil, fmt.Errorf("error while unmarshaling tips contract config response: %s", err)
	}
	return &config, nil
}
