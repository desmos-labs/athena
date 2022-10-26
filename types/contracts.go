package types

const (
	ContractTypeTips = "tips"
)

type Contract struct {
	Address string
	Type    string
	Height  int64
}

func NewContract(address string, contractType string, height int64) Contract {
	return Contract{
		Address: address,
		Type:    contractType,
		Height:  height,
	}
}

type ContractConfig struct {
	Address      string
	ConfigJSONBz []byte
	Height       int64
}

func NewContractConfig(address string, configJSONBz []byte, height int64) ContractConfig {
	return ContractConfig{
		Address:      address,
		ConfigJSONBz: configJSONBz,
		Height:       height,
	}
}
