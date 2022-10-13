package authz

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/authz"
	juno "github.com/forbole/juno/v3/types"

	"github.com/desmos-labs/djuno/v2/types"
)

// HandleMsg implements modules.MessageModule
func (m *Module) HandleMsg(index int, msg sdk.Msg, tx *juno.Tx) error {
	if len(tx.Logs) == 0 {
		return nil
	}

	switch desmosMsg := msg.(type) {
	case *authz.MsgGrant:
		return m.handleMsgGrant(tx, desmosMsg)
	case *authz.MsgRevoke:
		return m.handleMsgRevoke(tx, desmosMsg)
	}

	return nil
}

// handleMsgGrant handles the parsing of a single MsgGrant message
func (m *Module) handleMsgGrant(tx *juno.Tx, msg *authz.MsgGrant) error {
	// Unpack the interfaces
	err := msg.Grant.UnpackInterfaces(m.cdc)
	if err != nil {
		return fmt.Errorf("error when unpacking authorization: %s", err)
	}

	return m.db.SaveAuthzGrant(types.NewAuthzGrant(
		msg.Granter,
		msg.Grantee,
		msg.Grant.GetAuthorization(),
		msg.Grant.Expiration,
		tx.Height,
	))
}

// handleMsgGrant handles the parsing of a single MsgRevoke message
func (m *Module) handleMsgRevoke(tx *juno.Tx, msg *authz.MsgRevoke) error {
	return m.db.DeleteAuthzGrant(msg.Granter, msg.Grantee, msg.MsgTypeUrl, tx.Height)
}
