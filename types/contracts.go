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
