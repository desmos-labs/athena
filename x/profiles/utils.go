package profiles

import (
	"context"
	"fmt"

	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/forbole/juno/v3/node/remote"

	profilestypes "github.com/desmos-labs/desmos/v4/x/profiles/types"

	"github.com/desmos-labs/djuno/v2/types"
)

// GetDisplayName returns the name to be displayed for the user having the given address
func (m *Module) GetDisplayName(userAddress string) string {
	profile, err := m.GetUserProfile(userAddress)
	if err != nil || profile == nil {
		return fmt.Sprintf("%[1]s...%[2]s", userAddress[:9], userAddress[len(userAddress)-5:])
	}

	switch {
	case profile.Nickname != "":
		return fmt.Sprintf("%[1]s (@%[2]s)", profile.Nickname, profile.DTag)

	default:
		return fmt.Sprintf("@%[1]s", profile.DTag)
	}
}

// --------------------------------------------------------------------------------------------------------------------

// updateUserChainLinks updates the chain links for the given address, if any
func (m *Module) updateUserChainLinks(height int64, address string) error {
	chainLinks, err := m.queryAllUserChainLinks(height, address)
	if err != nil {
		return err
	}

	for _, chainLink := range chainLinks {
		err = m.db.SaveChainLink(chainLink)
		if err != nil {
			return err
		}
	}

	return err
}

// queryAllUserChainLinks queries all the chain links for the given address
func (m *Module) queryAllUserChainLinks(height int64, address string) ([]types.ChainLink, error) {
	var chainLinks []types.ChainLink

	var nextKey []byte
	var stop = false
	for !stop {
		res, err := m.client.ChainLinks(
			remote.GetHeightRequestContext(context.Background(), height),
			&profilestypes.QueryChainLinksRequest{
				User: address,
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

// updateUserDefaultChainLinks updates the default chain links associated with the given address, if any
func (m *Module) updateUserDefaultChainLinks(height int64, address string) error {
	chainLinks, err := m.queryAllUserDefaultChainLinks(height, address)
	if err != nil {
		return err
	}

	for _, chainLink := range chainLinks {
		err = m.db.SaveDefaultChainLink(chainLink)
		if err != nil {
			return err
		}
	}

	return nil
}

// queryAllUserDefaultChainLinks queries all the default chain links for the given address
func (m *Module) queryAllUserDefaultChainLinks(height int64, address string) ([]types.ChainLink, error) {
	var chainLinks []types.ChainLink

	var nextKey []byte
	var stop = false
	for !stop {
		res, err := m.client.DefaultExternalAddresses(
			remote.GetHeightRequestContext(context.Background(), height),
			&profilestypes.QueryDefaultExternalAddressesRequest{
				Owner: address,
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

// updateUserApplicationLinks updates the application links associated with the given address, if any
func (m *Module) updateUserApplicationLinks(height int64, address string) error {
	applicationLinks, err := m.queryAllUserApplicationLinks(height, address)
	if err != nil {
		return err
	}

	for _, applicationLink := range applicationLinks {
		err = m.db.SaveApplicationLink(applicationLink)
		if err != nil {
			return err
		}
	}

	return nil
}

// queryAllUserApplicationLinks queries all the application links for the given address
func (m *Module) queryAllUserApplicationLinks(height int64, address string) ([]types.ApplicationLink, error) {
	var chainLinks []types.ApplicationLink

	var nextKey []byte
	var stop = false
	for !stop {
		res, err := m.client.ApplicationLinks(
			remote.GetHeightRequestContext(context.Background(), height),
			&profilestypes.QueryApplicationLinksRequest{
				User: address,
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

// --------------------------------------------------------------------------------------------------------------------

// updateParams allows to update the profiles params by fetching them from the chain
func (m *Module) updateParams() error {
	height, err := m.node.LatestHeight()
	if err != nil {
		return fmt.Errorf("error while getting latest block height: %s", err)
	}

	res, err := m.client.Params(
		remote.GetHeightRequestContext(context.Background(), height),
		&profilestypes.QueryParamsRequest{},
	)
	if err != nil {
		return fmt.Errorf("error while getting params: %s", err)
	}

	err = m.db.SaveProfilesParams(types.NewProfilesParams(res.Params, height))
	if err != nil {
		return fmt.Errorf("error while storing profiles params: %s", err)
	}

	return nil
}
