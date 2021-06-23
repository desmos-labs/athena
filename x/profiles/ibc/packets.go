package ibc

import (
	"context"
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	channeltypes "github.com/cosmos/cosmos-sdk/x/ibc/core/04-channel/types"
	oracletypes "github.com/desmos-labs/desmos/x/oracle/types"
	profilestypes "github.com/desmos-labs/desmos/x/profiles/types"
	"github.com/desmos-labs/juno/client"

	desmosdb "github.com/desmos-labs/djuno/database"
	"github.com/desmos-labs/djuno/types"
)

// handleLinkChainAccountPacketData tries handling the given packet as it contains a LinkChainAccountPacketData
// instance. This is done to store chain links that are created using IBC.
func handleLinkChainAccountPacketData(
	height int64, packet channeltypes.Packet, profilesClient profilestypes.QueryClient, cdc codec.Marshaler, db *desmosdb.Db,
) (bool, error) {
	// Try reading the packet data
	var packetData profilestypes.LinkChainAccountPacketData
	err := cdc.UnmarshalJSON(packet.GetData(), &packetData)
	if err != nil {
		return false, nil
	}

	var sourceAddr profilestypes.AddressData
	err = cdc.UnpackAny(packetData.SourceAddress, &sourceAddr)
	if err != nil {
		return true, fmt.Errorf("error while deserializing source address: %s", err)
	}

	// Get the link from the chain
	res, err := profilesClient.UserChainLink(
		context.Background(),
		&profilestypes.QueryUserChainLinkRequest{
			User:      packetData.DestinationAddress,
			ChainName: packetData.SourceChainConfig.Name,
			Target:    sourceAddr.GetAddress(),
		},
		client.GetHeightRequestHeader(height),
	)
	if err != nil {
		return true, err
	}

	// Save the chain link
	return true, db.SaveChainLink(types.NewChainLink(res.Link, height))
}

// handleOracleRequestPacketData tries handling the given packet as it contains a OracleRequestPacketData
// instance. This is done in order to update existing application links when their state changes after
// Band Protocol ends the verification process.
func handleOracleRequestPacketData(
	height int64, packet channeltypes.Packet, profilesClient profilestypes.QueryClient, cdc codec.Marshaler, db *desmosdb.Db,
) (bool, error) {
	var data oracletypes.OracleRequestPacketData
	if err := cdc.UnmarshalJSON(packet.GetData(), &data); err != nil {
		return false, nil
	}

	res, err := profilesClient.ApplicationLinkByClientID(
		context.Background(),
		profilestypes.NewQueryApplicationLinkByClientIDRequest(data.ClientID),
		client.GetHeightRequestHeader(height),
	)
	if err != nil {
		return true, fmt.Errorf("error while getting application link by client id: %s", err)
	}

	return true, db.SaveApplicationLink(types.NewApplicationLink(res.Link, height))
}

// handleOracleResponsePacketData tries handling the given packet as it contains a OracleResponsePacketData
// instance. This is done in order to update existing application links when their state changes after
// Band Protocol ends the verification process.
func handleOracleResponsePacketData(
	height int64, packet channeltypes.Packet, profilesClient profilestypes.QueryClient, cdc codec.Marshaler, db *desmosdb.Db,
) (bool, error) {
	var data oracletypes.OracleResponsePacketData
	if err := cdc.UnmarshalJSON(packet.GetData(), &data); err != nil {
		return false, nil
	}

	res, err := profilesClient.ApplicationLinkByClientID(
		context.Background(),
		profilestypes.NewQueryApplicationLinkByClientIDRequest(data.ClientID),
		client.GetHeightRequestHeader(height),
	)
	if err != nil {
		return true, fmt.Errorf("error while getting application link by client id: %s", err)
	}

	return true, db.SaveApplicationLink(types.NewApplicationLink(res.Link, height))
}
