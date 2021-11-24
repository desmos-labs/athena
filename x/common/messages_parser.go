package common

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	channeltypes "github.com/cosmos/ibc-go/modules/core/04-channel/types"
	profilestypes "github.com/desmos-labs/desmos/v2/x/profiles/types"
	poststypes "github.com/desmos-labs/desmos/v2/x/staging/posts/types"
	"github.com/forbole/juno/v2/modules/messages"
)

var MessagesParser = messages.JoinMessageParsers(
	ibcMessagesParser,
	messages.CosmosMessageAddressesParser,
	desmosMessagesParser,
)

var desmosMessagesParser = messages.JoinMessageParsers(
	postsMessagesParser,
	profilesMessagesParser,
)

func ibcMessagesParser(_ codec.Codec, cosmosMsg sdk.Msg) ([]string, error) {
	switch msg := cosmosMsg.(type) {
	case *channeltypes.MsgRecvPacket:
		return []string{msg.Signer}, nil
	case *channeltypes.MsgAcknowledgement:
		return []string{msg.Signer}, nil
	case *channeltypes.MsgTimeout:
		return []string{msg.Signer}, nil
	}

	return nil, messages.MessageNotSupported(cosmosMsg)
}

func postsMessagesParser(_ codec.Codec, cosmosMsg sdk.Msg) ([]string, error) {
	switch msg := cosmosMsg.(type) {
	case *poststypes.MsgCreatePost:
		return []string{msg.Creator}, nil

	case *poststypes.MsgEditPost:
		return []string{msg.Editor}, nil

	case *poststypes.MsgAddPostReaction:
		return []string{msg.User}, nil

	case *poststypes.MsgRemovePostReaction:
		return []string{msg.User}, nil

	case *poststypes.MsgRegisterReaction:
		return []string{msg.Creator}, nil

	case *poststypes.MsgAnswerPoll:
		return []string{msg.Answerer}, nil

	case *poststypes.MsgReportPost:
		return []string{msg.User}, nil
	}

	return nil, messages.MessageNotSupported(cosmosMsg)
}

func profilesMessagesParser(_ codec.Codec, cosmosMsg sdk.Msg) ([]string, error) {
	switch msg := cosmosMsg.(type) {
	case *profilestypes.MsgSaveProfile:
		return []string{msg.Creator}, nil

	case *profilestypes.MsgDeleteProfile:
		return []string{msg.Creator}, nil

	case *profilestypes.MsgRequestDTagTransfer:
		return []string{msg.Sender, msg.Receiver}, nil

	case *profilestypes.MsgCancelDTagTransferRequest:
		return []string{msg.Sender, msg.Receiver}, nil

	case *profilestypes.MsgAcceptDTagTransferRequest:
		return []string{msg.Sender, msg.Receiver}, nil

	case *profilestypes.MsgRefuseDTagTransferRequest:
		return []string{msg.Sender, msg.Receiver}, nil

	case *profilestypes.MsgCreateRelationship:
		return []string{msg.Sender, msg.Receiver}, nil

	case *profilestypes.MsgDeleteRelationship:
		return []string{msg.User, msg.Counterparty}, nil

	case *profilestypes.MsgBlockUser:
		return []string{msg.Blocker, msg.Blocked}, nil

	case *profilestypes.MsgUnblockUser:
		return []string{msg.Blocker, msg.Blocked}, nil

	case *profilestypes.MsgLinkChainAccount:
		return []string{msg.Signer}, nil

	case *profilestypes.MsgLinkApplication:
		return []string{msg.Sender}, nil
	}

	return nil, messages.MessageNotSupported(cosmosMsg)
}
