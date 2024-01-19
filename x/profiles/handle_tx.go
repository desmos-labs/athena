package profiles

import (
	abci "github.com/cometbft/cometbft/abci/types"
	profilestypes "github.com/desmos-labs/desmos/v6/x/profiles/types"
	juno "github.com/forbole/juno/v5/types"

	"github.com/desmos-labs/athena/types"
	"github.com/desmos-labs/athena/utils"
)

// HandleTx implements modules.TransactionModule
func (m *Module) HandleTx(tx *juno.Tx) error {
	return utils.ParseTxEvents(tx, map[string]func(tx *juno.Tx, event abci.Event) error{
		profilestypes.EventTypeProfileSaved:            m.parseSaveProfileEvent,
		profilestypes.EventTypeProfileDeleted:          m.parseDeleteProfileEvent,
		profilestypes.EventTypeDTagTransferRequest:     m.parseRequestDTagTransferEvent,
		profilestypes.EventTypeDTagTransferAccept:      m.parseAcceptDTagTransferRequestEvent,
		profilestypes.EventTypeDTagTransferRefuse:      m.parseDeleteDTagTransferRequestEvent,
		profilestypes.EventTypeDTagTransferCancel:      m.parseDeleteDTagTransferRequestEvent,
		profilestypes.EventTypeLinkChainAccount:        m.parseLinkChainAccountEvent,
		profilestypes.EventTypeUnlinkChainAccount:      m.parseUnlinkChainAccountEvent,
		profilestypes.EventTypesApplicationLinkCreated: m.parseLinkApplicationEvent,
		profilestypes.EventTypeApplicationLinkDeleted:  m.parseUnlinkApplicationEvent,
	})
}

// -------------------------------------------------------------------------------------------------------------------

// parseSaveProfileEvent parses the save profile event
func (m *Module) parseSaveProfileEvent(tx *juno.Tx, event abci.Event) error {
	creator, err := GetCreatorAddressFromEvent(event)
	if err != nil {
		return err
	}

	return m.UpdateProfiles(tx.Height, []string{creator})
}

// parseDeleteProfileEvent parses the delete profile event
func (m *Module) parseDeleteProfileEvent(tx *juno.Tx, event abci.Event) error {
	creator, err := GetCreatorAddressFromEvent(event)
	if err != nil {
		return err
	}

	return m.db.DeleteProfile(creator, tx.Height)
}

// parseRequestDTagTransferEvent parses the request dtag transfer event
func (m *Module) parseRequestDTagTransferEvent(tx *juno.Tx, event abci.Event) error {
	sender, err := GetRequestSenderAddressFromEvent(event)
	if err != nil {
		return err
	}

	receiver, err := GetRequestReceiverAddressFromEvent(event)
	if err != nil {
		return err
	}

	dTagToTrade, err := juno.FindAttributeByKey(event, profilestypes.AttributeKeyDTagToTrade)
	if err != nil {
		return err
	}

	return m.db.SaveDTagTransferRequest(types.NewDTagTransferRequest(
		profilestypes.NewDTagTransferRequest(dTagToTrade.Value, sender, receiver),
		tx.Height,
	))
}

// parseAcceptDTagTransferRequestEvent parses the accept dtag transfer request event
func (m *Module) parseAcceptDTagTransferRequestEvent(tx *juno.Tx, event abci.Event) error {
	sender, err := GetRequestSenderAddressFromEvent(event)
	if err != nil {
		return err
	}

	receiver, err := GetRequestReceiverAddressFromEvent(event)
	if err != nil {
		return err
	}

	return m.UpdateProfiles(tx.Height, []string{sender, receiver})
}

// parseDeleteDTagTransferRequestEvent parses the refuse or cancel dtag transfer request events
func (m *Module) parseDeleteDTagTransferRequestEvent(tx *juno.Tx, event abci.Event) error {
	sender, err := GetRequestSenderAddressFromEvent(event)
	if err != nil {
		return err
	}

	receiver, err := GetRequestReceiverAddressFromEvent(event)
	if err != nil {
		return err
	}

	return m.db.DeleteDTagTransferRequest(types.NewDTagTransferRequest(
		profilestypes.NewDTagTransferRequest("", sender, receiver),
		tx.Height,
	))
}

// parseLinkChainAccountEvent parses the link chain account event
func (m *Module) parseLinkChainAccountEvent(tx *juno.Tx, event abci.Event) error {
	ownerAttr, err := juno.FindAttributeByKey(event, profilestypes.AttributeKeyChainLinkOwner)
	if err != nil {
		return err
	}
	owner := ownerAttr.Value

	// Save the chain links
	err = m.updateUserChainLinks(tx.Height, owner)
	if err != nil {
		return err
	}

	// Update the default chain links
	return m.updateUserDefaultChainLinks(tx.Height, owner)
}

// parseUnlinkChainAccountEvent parses the unlink chain account event
func (m *Module) parseUnlinkChainAccountEvent(tx *juno.Tx, event abci.Event) error {
	ownerAttr, err := juno.FindAttributeByKey(event, profilestypes.AttributeKeyChainLinkOwner)
	if err != nil {
		return err
	}
	owner := ownerAttr.Value

	targetAttr, err := juno.FindAttributeByKey(event, profilestypes.AttributeKeyChainLinkExternalAddress)
	if err != nil {
		return err
	}
	target := targetAttr.Value

	chainNameAttr, err := juno.FindAttributeByKey(event, profilestypes.AttributeKeyChainLinkChainName)
	if err != nil {
		return err
	}
	chainName := chainNameAttr.Value

	err = m.db.DeleteChainLink(owner, target, chainName, tx.Height)
	if err != nil {
		return err
	}

	// Update the default chain links
	return m.updateUserDefaultChainLinks(tx.Height, owner)
}

// parseLinkApplicationEvent parses the link application event
func (m *Module) parseLinkApplicationEvent(tx *juno.Tx, event abci.Event) error {
	ownerAttr, err := juno.FindAttributeByKey(event, profilestypes.AttributeKeyUser)
	if err != nil {
		return err
	}
	owner := ownerAttr.Value

	return m.updateUserApplicationLinks(tx.Height, owner)
}

// parseUnlinkApplicationEvent parses the unlink application event
func (m *Module) parseUnlinkApplicationEvent(tx *juno.Tx, event abci.Event) error {
	ownerAttr, err := juno.FindAttributeByKey(event, profilestypes.AttributeKeyUser)
	if err != nil {
		return err
	}
	owner := ownerAttr.Value

	applicationAttr, err := juno.FindAttributeByKey(event, profilestypes.AttributeKeyApplicationName)
	if err != nil {
		return err
	}
	application := applicationAttr.Value

	usernameAttr, err := juno.FindAttributeByKey(event, profilestypes.AttributeKeyApplicationUsername)
	if err != nil {
		return err
	}
	username := usernameAttr.Value

	return m.db.DeleteApplicationLink(owner, application, username, tx.Height)
}
