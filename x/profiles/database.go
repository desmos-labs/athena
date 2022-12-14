package profiles

import (
	profilestypes "github.com/desmos-labs/desmos/v4/x/profiles/types"

	"github.com/desmos-labs/djuno/v2/types"
)

type Database interface {
	SaveProfilesParams(params types.ProfilesParams) error
	SaveUserIfNotExisting(address string, height int64) error
	GetUserByAddress(address string) (*profilestypes.Profile, error)
	SaveProfile(profile *types.Profile) error
	DeleteProfile(address string, height int64) error
	GetProfilesAddresses() ([]string, error)
	SaveDTagTransferRequest(request types.DTagTransferRequest) error
	DeleteDTagTransferRequest(request types.DTagTransferRequest) error
	SaveChainLink(link types.ChainLink) error
	SaveDefaultChainLink(chainLink types.ChainLink) error
	DeleteChainLink(user string, externalAddress string, chainName string, height int64) error
	DeleteProfileChainLinks(user string) error
	SaveApplicationLink(link types.ApplicationLink) error
	GetApplicationLinkInfos() ([]types.ApplicationLinkInfo, error)
	DeleteApplicationLink(user, application, username string, height int64) error
}
