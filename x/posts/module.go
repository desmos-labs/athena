package posts

import (
	"github.com/cosmos/cosmos-sdk/codec"

	"github.com/desmos-labs/djuno/database"

	"github.com/forbole/juno/v2/modules"
)

var _ modules.Module = &Module{}
var _ modules.GenesisModule = &Module{}
var _ modules.MessageModule = &Module{}

// Module represents the x/posts module handler
type Module struct {
	cdc            codec.Codec
	db             *database.Db
	profilesModule ProfilesModule
}

// NewModule allows to build a new Module instance
func NewModule(cdc codec.Codec, db *database.Db, profilesModule ProfilesModule) *Module {
	return &Module{
		cdc:            cdc,
		db:             db,
		profilesModule: profilesModule,
	}
}

// Name implements Module
func (m *Module) Name() string {
	return "posts"
}
