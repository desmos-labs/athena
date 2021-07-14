package utils

import (
	"context"
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	profilestypes "github.com/desmos-labs/desmos/x/profiles/types"
	"github.com/desmos-labs/juno/client"

	desmosdb "github.com/desmos-labs/djuno/database"
	"github.com/desmos-labs/djuno/types"
)

// UpdateProfiles updates the profiles associated with the given addresses, if any.
func UpdateProfiles(
	height int64, addresses []string, profilesClient profilestypes.QueryClient, cdc codec.Marshaler, db *desmosdb.Db,
) error {
	for _, address := range addresses {
		res, err := profilesClient.Profile(
			context.Background(),
			profilestypes.NewQueryProfileRequest(address),
			client.GetHeightRequestHeader(height),
		)
		if err != nil {
			return fmt.Errorf("error while getting profile from gRPC: %s", err)
		}

		if res.Profile != nil {
			var p profilestypes.Profile
			err = cdc.UnpackAny(res.Profile, &p)
			if err != nil {
				return fmt.Errorf("error while unpacking profile: %s", err)
			}

			err = db.SaveProfile(types.NewProfile(&p, height))
			if err != nil {
				return fmt.Errorf("error while saving profile: %s", err)
			}
		}
	}

	return nil
}
