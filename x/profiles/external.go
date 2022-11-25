package profiles

import (
	"context"
	"fmt"

	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/forbole/juno/v4/node/remote"

	profilestypes "github.com/desmos-labs/desmos/v4/x/profiles/types"

	"github.com/desmos-labs/djuno/v2/types"
)

// GetUserProfile queries the profile for the user having the given address, if any
func (m *Module) GetUserProfile(userAddress string) (*types.Profile, error) {
	height, err := m.node.LatestHeight()
	if err != nil {
		return nil, fmt.Errorf("error while getting latest height: %s", err)
	}

	profile, err := m.getProfile(height, userAddress)
	if err != nil {
		return nil, err
	}

	return profile, nil
}

// UpdateProfiles updates the profiles associated with the given addresses, if any
func (m *Module) UpdateProfiles(height int64, addresses []string) error {
	for _, address := range addresses {
		profile, err := m.getProfile(height, address)
		if err != nil {
			return err
		}

		err = m.db.SaveProfile(profile)
		if err != nil {
			return fmt.Errorf("error while saving profile: %s", err)
		}
	}

	return nil
}

func (m *Module) getProfile(height int64, address string) (*types.Profile, error) {
	res, err := m.client.Profile(
		remote.GetHeightRequestContext(context.Background(), height),
		profilestypes.NewQueryProfileRequest(address),
	)
	if err != nil {
		// If the profile was not found, just return nil
		if status.Code(err) == codes.NotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("error while getting profile from gRPC: %s", err)
	}

	var account authtypes.AccountI
	err = m.cdc.UnpackAny(res.Profile, &account)
	if err != nil {
		return nil, fmt.Errorf("error while unpacking profile: %s", err)
	}

	return types.NewProfile(account.(*profilestypes.Profile), height), nil
}

// --------------------------------------------------------------------------------------------------------------------

// RefreshChainLinks fetches and stores all the chain links present on the chain
func (m *Module) RefreshChainLinks(height int64) error {
	// Get the chain links
	chainLinks, err := m.queryAllChainLinks(height)
	if err != nil {
		return err
	}

	// Save the chain links
	for _, chainLink := range chainLinks {
		err = m.db.SaveChainLink(chainLink)
		if err != nil {
			return err
		}
	}

	// Get the default chain links
	defaultChainLinks, err := m.queryAllDefaultChainLinks(height)
	if err != nil {
		return err
	}

	// Save the default chain links
	for _, chainLink := range defaultChainLinks {
		err = m.db.SaveDefaultChainLink(chainLink)
		if err != nil {
			return err
		}
	}

	return nil
}

// queryAllChainLinks queries all the chain links stored inside the chain
func (m *Module) queryAllChainLinks(height int64) ([]types.ChainLink, error) {
	var chainLinks []types.ChainLink

	var nextKey []byte
	var stop = false
	for !stop {
		res, err := m.client.ChainLinks(
			remote.GetHeightRequestContext(context.Background(), height),
			&profilestypes.QueryChainLinksRequest{
				Pagination: &query.PageRequest{
					Key: nextKey,
				},
			},
		)
		if err != nil {
			return nil, err
		}

		for _, link := range res.Links {
			chainLinks = append(chainLinks, types.NewChainLink(link, height))
		}

		nextKey = res.Pagination.NextKey
		stop = nextKey == nil
	}

	return chainLinks, nil
}

// queryAllDefaultChainLinks queries all the default chain links stored inside the chain
func (m *Module) queryAllDefaultChainLinks(height int64) ([]types.ChainLink, error) {
	var chainLinks []types.ChainLink

	var nextKey []byte
	var stop = false
	for !stop {
		res, err := m.client.DefaultExternalAddresses(
			remote.GetHeightRequestContext(context.Background(), height),
			&profilestypes.QueryDefaultExternalAddressesRequest{
				Pagination: &query.PageRequest{
					Key: nextKey,
				},
			},
		)
		if err != nil {
			return nil, err
		}

		for _, link := range res.Links {
			chainLinks = append(chainLinks, types.NewChainLink(link, height))
		}

		nextKey = res.Pagination.NextKey
		stop = nextKey == nil
	}

	return chainLinks, nil
}

// --------------------------------------------------------------------------------------------------------------------

// RefreshApplicationLinks fetches and stores all the application links present on the chain
func (m *Module) RefreshApplicationLinks(height int64) error {
	// Get the chain links
	applicationLinks, err := m.queryAllApplicationLinks(height)
	if err != nil {
		return err
	}

	// Save the application links
	for _, chainLink := range applicationLinks {
		err = m.db.SaveApplicationLink(chainLink)
		if err != nil {
			return err
		}
	}

	return nil
}

// queryAllApplicationLinks queries all the application links stored inside the chain
func (m *Module) queryAllApplicationLinks(height int64) ([]types.ApplicationLink, error) {
	var chainLinks []types.ApplicationLink

	var nextKey []byte
	var stop = false
	for !stop {
		res, err := m.client.ApplicationLinks(
			remote.GetHeightRequestContext(context.Background(), height),
			&profilestypes.QueryApplicationLinksRequest{
				Pagination: &query.PageRequest{
					Key: nextKey,
				},
			},
		)
		if err != nil {
			return nil, err
		}

		for _, link := range res.Links {
			chainLinks = append(chainLinks, types.NewApplicationLink(link, height))
		}

		nextKey = res.Pagination.NextKey
		stop = nextKey == nil
	}

	return chainLinks, nil
}
