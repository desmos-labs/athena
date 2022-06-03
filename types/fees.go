package types

import (
	feestypes "github.com/desmos-labs/desmos/v3/x/fees/types"
)

type FeesParams struct {
	feestypes.Params
	Height int64
}

func NewFeesParams(params feestypes.Params, height int64) FeesParams {
	return FeesParams{
		Params: params,
		Height: height,
	}
}
