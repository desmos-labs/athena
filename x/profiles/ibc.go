package profiles

import (
	"context"
	"fmt"

	juno "github.com/forbole/juno/v5/types"

	"github.com/forbole/juno/v5/node/remote"

	"github.com/desmos-labs/athena/v2/types"

	channeltypes "github.com/cosmos/ibc-go/v7/modules/core/04-channel/types"

	oracletypes "github.com/desmos-labs/desmos/v6/x/oracle/types"
	profilestypes "github.com/desmos-labs/desmos/v6/x/profiles/types"
)

// packetHandler defines a function that handles a packet.
// It returns true iff it was able to handle the packet, and an error if something goes wrong.
type packetHandler = func(height int64, packet channeltypes.Packet) (bool, error)

// handlePacket tries handling the given packet that was received at the given height
func (m *Module) handlePacket(tx *juno.Tx, packet channeltypes.Packet) error {
	// Try handling the packet
	handlers := []packetHandler{
		m.handleLinkChainAccountPacketData,
		m.handleOracleRequestPacketData,
		m.handleOracleResponsePacketData,
	}

	for _, handler := range handlers {
		handled, err := handler(tx.Height, packet)
		if handled {
			return err
		}
	}

	return nil
}

// handleLinkChainAccountPacketData tries handling the given packet as it contains a LinkChainAccountPacketData
// instance. This is done to store chain links that are created using IBC.
func (m *Module) handleLinkChainAccountPacketData(height int64, packet channeltypes.Packet) (bool, error) {
	// Try reading the packet data
	var packetData profilestypes.LinkChainAccountPacketData
	err := m.cdc.UnmarshalJSON(packet.GetData(), &packetData)
	if err != nil {
		return false, nil
	}

	var sourceAddr profilestypes.AddressData
	err = m.cdc.UnpackAny(packetData.SourceAddress, &sourceAddr)
	if err != nil {
		return true, fmt.Errorf("error while deserializing source address: %s", err)
	}

	// Get the link from the chain
	res, err := m.profilesClient.ChainLinks(
		remote.GetHeightRequestContext(context.Background(), height),
		&profilestypes.QueryChainLinksRequest{
			User:      packetData.DestinationAddress,
			ChainName: packetData.SourceChainConfig.Name,
			Target:    sourceAddr.GetValue(),
		},
	)
	if err != nil {
		return true, err
	}

	if len(res.Links) == 0 {
		return true, fmt.Errorf("chain link not found on chain")
	}

	if len(res.Links) > 1 {
		return true, fmt.Errorf("duplicated chain link found on chain")
	}

	// Save the chain link
	return true, m.db.SaveChainLink(types.NewChainLink(res.Links[0], height))
}

// handleOracleRequestPacketData tries handling the given packet as it contains a OracleRequestPacketData
// instance. This is done in order to update existing application links when their state changes after
// Band Protocol ends the verification process.
func (m *Module) handleOracleRequestPacketData(height int64, packet channeltypes.Packet) (bool, error) {
	var data oracletypes.OracleRequestPacketData
	if err := m.cdc.UnmarshalJSON(packet.GetData(), &data); err != nil {
		return false, nil
	}

	res, err := m.profilesClient.ApplicationLinkByClientID(
		remote.GetHeightRequestContext(context.Background(), height),
		profilestypes.NewQueryApplicationLinkByClientIDRequest(data.ClientID),
	)
	if err != nil {
		return true, fmt.Errorf("error while getting application link by client id: %s", err)
	}

	return true, m.db.SaveApplicationLink(types.NewApplicationLink(res.Link, height))
}

// handleOracleResponsePacketData tries handling the given packet as it contains a OracleResponsePacketData
// instance. This is done in order to update existing application links when their state changes after
// Band Protocol ends the verification process.
func (m *Module) handleOracleResponsePacketData(height int64, packet channeltypes.Packet) (bool, error) {
	var data oracletypes.OracleResponsePacketData
	if err := m.cdc.UnmarshalJSON(packet.GetData(), &data); err != nil {
		return false, nil
	}

	res, err := m.profilesClient.ApplicationLinkByClientID(
		remote.GetHeightRequestContext(context.Background(), height),
		profilestypes.NewQueryApplicationLinkByClientIDRequest(data.ClientID),
	)
	if err != nil {
		return true, fmt.Errorf("error while getting application link by client id: %s", err)
	}

	return true, m.db.SaveApplicationLink(types.NewApplicationLink(res.Link, height))
}
