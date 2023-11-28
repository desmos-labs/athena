package profiles

import (
	"context"
	"fmt"

	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	authttypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/forbole/juno/v5/node/remote"

	profilestypes "github.com/desmos-labs/desmos/v6/x/profiles/types"

	"github.com/desmos-labs/athena/types"
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
	res, err := m.profilesClient.Profile(
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

	var account authttypes.AccountI
	err = m.cdc.UnpackAny(res.Profile, &account)
	if err != nil {
		return nil, fmt.Errorf("error while unpacking profile: %s", err)
	}

	return types.NewProfile(account.(*profilestypes.Profile), height), nil
}

// --------------------------------------------------------------------------------------------------------------------

// RefreshProfiles fetches and stores all the profiles present on the chain
func (m *Module) RefreshProfiles(height int64) error {
	profiles, err := m.queryAllProfiles(height)
	if err != nil {
		return fmt.Errorf("error while querying profiles: %s", err)
	}

	for _, profile := range profiles {
		log.Debug().Str("module", "profiles").Str("dTag", profile.DTag).Msg("saving profile")
		err = m.db.SaveProfile(profile)
		if err != nil {
			return err
		}
	}

	return nil
}

// queryAllProfiles queries all the profiles stored inside the chain
func (m *Module) queryAllProfiles(height int64) ([]*types.Profile, error) {
	var profiles []*types.Profile

	var nextKey []byte
	var stop = false
	for !stop {
		res, err := m.authClient.Accounts(
			remote.GetHeightRequestContext(context.Background(), height),
			&authttypes.QueryAccountsRequest{
				Pagination: &query.PageRequest{
					Key:        nextKey,
					Limit:      1000,
					CountTotal: true,
				},
			},
		)
		if err != nil {
			return nil, err
		}

		for _, acc := range res.Accounts {
			var account authttypes.AccountI
			err = m.cdc.UnpackAny(acc, &account)
			if err != nil {
				return nil, fmt.Errorf("error while unpacking account: %s", err)
			}

			if profile, ok := account.(*profilestypes.Profile); ok {
				profiles = append(profiles, types.NewProfile(profile, height))
			}
		}

		nextKey = res.Pagination.NextKey
		stop = nextKey == nil
	}

	return profiles, nil
}

// --------------------------------------------------------------------------------------------------------------------

// RefreshChainLinks fetches and stores all the chain links present on the chain
func (m *Module) RefreshChainLinks(height int64) error {
	// Delete the existing chain links
	err := m.db.DeleteAllChainLinks(height)
	if err != nil {
		return err
	}

	// Get the chain links
	chainLinks, err := m.queryAllChainLinks(height)
	if err != nil {
		return err
	}

	// Save the chain links
	for _, chainLink := range chainLinks {
		log.Debug().Str("module", "profiles").Str("user", chainLink.User).Str("chain", chainLink.ChainConfig.Name).Msg("saving chain link")
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

	// Delete the default chain links
	err = m.db.DeleteAllDefaultChainLinks(height)
	if err != nil {
		return err
	}

	// Save the default chain links
	for _, chainLink := range defaultChainLinks {
		log.Debug().Str("module", "profiles").Str("user", chainLink.User).Str("chain", chainLink.ChainConfig.Name).Msg("saving default chain link")
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
		res, err := m.profilesClient.ChainLinks(
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
		res, err := m.profilesClient.DefaultExternalAddresses(
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
	// Delete all the application links
	err := m.db.DeleteAllApplicationLinks(height)
	if err != nil {
		return err
	}

	// Get the chain links
	applicationLinks, err := m.queryAllApplicationLinks(height)
	if err != nil {
		return err
	}

	// Save the application links
	for _, applicationLink := range applicationLinks {
		log.Debug().Str("module", "applications").Str("user", applicationLink.User).
			Str("application", applicationLink.Data.Application).
			Str("username", applicationLink.Data.Username).
			Msg("saving application link")
		err = m.db.SaveApplicationLink(applicationLink)
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
		res, err := m.profilesClient.ApplicationLinks(
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
