package types

const (
	ContractTypeTips = "tips"
)

type Contract struct {
	Address  string
	Type     string
	ConfigBz []byte
	Height   int64
}

func NewContract(address string, contractType string, contactConfigBz []byte, height int64) Contract {
	return Contract{
		Address:  address,
		Type:     contractType,
		ConfigBz: contactConfigBz,
		Height:   height,
	}
}
