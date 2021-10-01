package ibc

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	channeltypes "github.com/cosmos/ibc-go/modules/core/04-channel/types"
	profilestypes "github.com/desmos-labs/desmos/v2/x/profiles/types"

	desmosdb "github.com/desmos-labs/djuno/database"
)

var (
	handlers = []packetHandler{
		handleLinkChainAccountPacketData,
		handleOracleRequestPacketData,
		handleOracleResponsePacketData,
	}
)

// packetHandler defines a function that handles a packet.
// It returns true iff it was able to handle the packet, and an error if something goes wrong.
type packetHandler = func(
	height int64, packet channeltypes.Packet,
	profilesClient profilestypes.QueryClient, cdc codec.Codec, db *desmosdb.Db,
) (bool, error)

// HandlePacket tries handling the given packet that was received at the given height
func HandlePacket(
	height int64, packet channeltypes.Packet, client profilestypes.QueryClient, cdc codec.Codec, db *desmosdb.Db,
) error {
	// Try handling the packet
	for _, handler := range handlers {
		handled, err := handler(height, packet, client, cdc, db)
		if handled {
			return err
		}
	}

	return fmt.Errorf("cannot handle packet directed to port %s and channel %s", packet.DestinationPort, packet.DestinationChannel)
}
